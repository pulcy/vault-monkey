package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
	b, _ := Factory(logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadCreds(t, b, "web"),
		},
	})

}

func TestBackend_roleCrud(t *testing.T) {
	b, _ := Factory(logical.TestBackendConfig())

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepRole(t),
			testAccStepReadRole(t, "web", testRole),
			testAccStepDeleteRole(t, "web"),
			testAccStepReadRole(t, "web", ""),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PG_URL"); v == "" {
		t.Fatal("PG_URL must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"value": os.Getenv("PG_URL"),
		},
	}
}

func testAccStepRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/web",
		Data: map[string]interface{}{
			"sql": testRole,
		},
	}
}

func testAccStepDeleteRole(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + n,
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated credentials: %v", d)

			conn, err := pq.ParseURL(os.Getenv("PG_URL"))
			if err != nil {
				t.Fatal(err)
			}

			conn += " timezone=utc"

			db, err := sql.Open("postgres", conn)
			if err != nil {
				t.Fatal(err)
			}

			returnedRows := func() int {
				stmt, err := db.Prepare(fmt.Sprintf(
					"SELECT DISTINCT schemaname FROM pg_tables WHERE has_table_privilege('%s', 'information_schema.role_column_grants', 'select');",
					d.Username))
				if err != nil {
					return -1
				}
				defer stmt.Close()

				rows, err := stmt.Query()
				if err != nil {
					return -1
				}
				defer rows.Close()

				i := 0
				for rows.Next() {
					i++
				}
				return i
			}

			userRows := returnedRows()
			if userRows != 2 {
				t.Fatalf("did not get expected number of rows, got %d", userRows)
			}

			resp, err = b.HandleRequest(&logical.Request{
				Operation: logical.RevokeOperation,
				Secret: &logical.Secret{
					InternalData: map[string]interface{}{
						"secret_type": "creds",
						"username":    d.Username,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp != nil {
				if resp.IsError() {
					return fmt.Errorf("Error on resp: %#v", *resp)
				}
			}

			userRows = returnedRows()
			// User shouldn't exist so returnedRows() should encounter an error and exit with -1
			if userRows != -1 {
				t.Fatalf("did not get expected number of rows, got %d", userRows)
			}

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name string, sql string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if sql == "" {
					return nil
				}

				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				SQL string `mapstructure:"sql"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.SQL != sql {
				return fmt.Errorf("bad: %#v", resp)
			}

			return nil
		},
	}
}

const testRole = `
CREATE ROLE "{{name}}" WITH
  LOGIN
  PASSWORD '{{password}}'
  VALID UNTIL '{{expiration}}';
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";
`
