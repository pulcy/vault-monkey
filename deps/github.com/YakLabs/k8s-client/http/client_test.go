package http_test

import (
	"os"
	"testing"

	"github.com/YakLabs/k8s-client/http"
	"github.com/stretchr/testify/require"
)

// create a test client based on env variables.
func testClient(t *testing.T) *http.Client {
	server := os.Getenv("K8S_SERVER")

	if server == "" {
		server = "http://127.0.0.1:8001"
	}

	opts := []http.OptionsFunc{
		http.SetServer(server),
	}

	if caFile := os.Getenv("K8S_CAFILE"); caFile != "" {
		opts = append(opts, http.SetCAFromFile(caFile))
	}

	if clientCert := os.Getenv("K8S_CLIENTCERT"); clientCert != "" {
		opts = append(opts, http.SetClientCertFromFile(clientCert))
	}

	if clientKey := os.Getenv("K8S_CLIENTKEY"); clientKey != "" {
		opts = append(opts, http.SetClientKeyFromFile(clientKey))
	}

	c, err := http.New(opts...)
	require.Nil(t, err)

	return c
}
