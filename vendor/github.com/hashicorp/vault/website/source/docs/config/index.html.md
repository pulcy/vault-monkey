---
layout: "docs"
page_title: "Server Configuration"
sidebar_current: "docs-config"
description: |-
  Vault server configuration reference.
---

# Server Configuration

Outside of development mode, Vault servers are configured using a file.
The format of this file is [HCL](https://github.com/hashicorp/hcl) or JSON.
An example configuration is shown below:

```javascript
backend "consul" {
  address = "127.0.0.1:8500"
  path = "vault"
}

listener "tcp" {
  address = "127.0.0.1:8200"
  tls_disable = 1
}

telemetry {
  statsite_address = "127.0.0.1:8125"
  disable_hostname = true
}
```

After the configuration is written, use the `-config` flag with `vault server`
to specify where the configuration is.

## Reference

* `backend` (required) - Configures the storage backend where Vault data
  is stored. There are multiple options available for storage backends,
  and they're documented below.

* `ha_backend` (optional) - Configures the storage backend where Vault HA
  coordination will take place. Must be an HA-supporting backend using the
  configuration options as documented below. If not set, HA will be attempted
  on the backend given in the `backend` parameter.

* `listener` (required) - Configures how Vault is listening for API requests.
  "tcp" is currently the only option available. A full reference for the
   inner syntax is below.

* `disable_cache` (optional) - A boolean. If true, this will disable the
  read cache used by the physical storage subsystem. This will very
  significantly impact performance.

* `disable_mlock` (optional) - A boolean. If true, this will disable the
  server from executing the `mlock` syscall to prevent memory from being
  swapped to disk. This is not recommended in production (see below).

* `telemetry` (optional)  - Configures the telemetry reporting system
  (see below).

* `default_lease_ttl` (optional) - Configures the default lease duration
  for tokens and secrets. This is a string value using a suffix, e.g. "720h".
  Default value is 30 days. This value cannot be larger than `max_lease_ttl`.

* `max_lease_ttl` (optional) - Configures the maximum possible
  lease duration for tokens and secrets. This is a string value using a suffix,
  e.g. "720h". Default value is 30 days.

In production, you should only consider setting the `disable_mlock` option
on Linux systems that only use encrypted swap or do not use swap at all.
Vault does not currently support memory locking on Mac OS X and Windows
and so the feature is automatically disabled on those platforms.  To give
the Vault executable access to the `mlock` syscall on Linux systems:

```shell
sudo setcap cap_ipc_lock=+ep $(readlink -f $(which vault))
```

## Listener Reference

For the `listener` section, the only supported listener currently
is "tcp". Regardless of future plans, this is the recommended listener,
since it allows for HA mode.

The supported options are:

  * `address` (optional) - The address to bind to for listening. This
      defaults to "127.0.0.1:8200".

  * `tls_disable` (optional) - If true, then TLS will be disabled.
      This will parse as boolean value, and can be set to "0", "no",
      "false", "1", "yes", or "true". This is an opt-in; Vault assumes
      by default that TLS will be used.

  * `tls_cert_file` (required unless disabled) - The path to the certificate
      for TLS.

  * `tls_key_file` (required unless disabled) - The path to the private key
      for the certificate.

  * `tls_min_version` (optional) - **(Vault > 0.2)** If provided, specifies
      the minimum supported version of TLS. Accepted values are "tls10", "tls11"
      or "tls12". This defaults to "tls12". WARNING: TLS 1.1 and lower
      are generally considered less secure; avoid using these if
      possible.

## Telemetry Reference

For the `telemetry` section, there is no resource name. All configuration
is within the object itself.

* `statsite_address` (optional) - An address to a [Statsite](https://github.com/armon/statsite)
  instances for metrics. This is highly recommended for production usage.

* `statsd_address` (optional) - This is the same as `statsite_address` but
  for StatsD.

* `disable_hostname` (optional) - Whether or not to prepend runtime telemetry
  with the machines hostname. This is a global option. Defaults to false.

## Backend Reference

For the `backend` section, the supported physical backends are shown below.
Vault requires that the backend itself will be responsible for backups,
durability, etc.

__*Please note*__: The only physical backends actively maintained by HashiCorp
are `consul`, `inmem`, and `file`. The other backends are community-derived and
community-supported. We include them in the hope that they will be useful to
those users that wish to utilize them, but they receive minimal validation and
testing from HashiCorp, and HashiCorp staff may not be knowledgeable about the
data store being utilized. If you encounter problems with them, we will attempt
to help you, but may refer you to the backend author.

  * `consul` - Store data within [Consul](https://www.consul.io). This
    backend supports HA. It is the most recommended backend for Vault and has
    been shown to work at high scale under heavy load.

  * `etcd` - Store data within [etcd](https://coreos.com/etcd/).
    This backend supports HA. This is a community-supported backend.

  * `zookeeper` - Store data within [Zookeeper](https://zookeeper.apache.org/).
    This backend supports HA. This is a community-supported backend.

  * `dynamodb` - Store data in a [DynamoDB](https://aws.amazon.com/dynamodb/) table.
    This backend supports HA. This is a community-supported backend.

  * `s3` - Store data within an S3 bucket [S3](https://aws.amazon.com/s3/).
    This backend does not support HA. This is a community-supported backend.

  * `mysql` - Store data within MySQL. This backend does not support HA. This
    is a community-supported backend.

  * `inmem` - Store data in-memory. This is only really useful for
    development and experimentation. Data is lost whenever Vault is
    restarted.

  * `file` - Store data on the filesystem using a directory structure.
    This backend does not support HA.


#### Common Backend Options

All backends support the following options:

  * `advertise_addr` (optional) - For backends that support HA, this
      is the address to advertise to other Vault servers in the cluster
      for request forwarding. Most HA backends will attempt to determine
      the advertise address if not provided. This can also be set via
      the `VAULT_ADVERTISE_ADDR` environment variable.

#### Backend Reference: Consul

For Consul, the following options are supported:

  * `path` (optional) - The path within Consul where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address of the Consul agent to talk to.
      Defaults to the local agent address, if available.

  * `scheme` (optional) - "http" or "https" for talking to Consul.

  * `token` (optional) - An access token to use to write data to Consul.

  * `max_parallel` (optional) - The maximum number of connections to Consul;
      defaults to "128".

  * `tls_skip_verify` (optional) - If non-empty, then TLS host verification
      will be disabled for Consul communication.
      Defaults to false.

  The following settings should be set according to your [Consul encryption settings](https://www.consul.io/docs/agent/encryption.html):

  * `tls_ca_file` (optional) - The path to the CA certificate used for Consul communication.
      Defaults to system bundle if not specified.
      Set accordingly to the [ca_file](https://www.consul.io/docs/agent/options.html#ca_file) setting in Consul.

  * `tls_cert_file` (optional) - The path to the certificate for Consul communication.
      Set accordingly to the [cert_file](https://www.consul.io/docs/agent/options.html#cert_file) setting in Consul.

  * `tls_key_file` (optional) - The path to the private key for Consul communication.
      Set accordingly to the [key_file](https://www.consul.io/docs/agent/options.html#key_file) setting in Consul.

#### Backend Reference: etcd (Community-Supported)

For etcd, the following options are supported:

  * `path` (optional) - The path within etcd where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address(es) of the etcd instance(s) to talk to.
      Can be comma separated list (protocol://host:port) of many etcd instances.
      Defaults to "http://localhost:2379" if not specified.

  * `sync` (optional) - Should we synchronize the list of available etcd
      servers on startup?  This is a **string** value to allow for auto-sync to
      be implemented later. It can be set to "0", "no", "n", "false", "1",
      "yes", "y", or "true".  Defaults to on.  Set to false if your etcd
      cluster is behind a proxy server and syncing causes Vault to fail.

  * `username` (optional) - Username to use when authenticating with the etcd
      server.  May also be specified via the ETCD_USERNAME environment variable.

  * `password` (optional) - Password to use when authenticating with the etcd
      server.  May also be specified via the ETCD_PASSWORD environment variable.

  * `tls_ca_file` (optional) - The path to the CA certificate used for etcd communication.
      Defaults to system bundle if not specified.

  * `tls_cert_file` (optional) - The path to the certificate for etcd communication.

  * `tls_key_file` (optional) - The path to the private key for etcd communication.

#### Backend Reference: Zookeeper (Community-Supported)

For Zookeeper, the following options are supported:

  * `path` (optional) - The path within Zookeeper where data will be stored.
      Defaults to "vault/".

  * `address` (optional) - The address(es) of the Zookeeper instance(s) to talk to.
      Can be comma separated list (host:port) of many Zookeeper instances.
      Defaults to "localhost:2181" if not specified.

#### Backend Reference: DynamoDB (Community-Supported)

The DynamoDB backend has the following options:

  * `table` (optional) - The name of the DynamoDB table to store data in. The default table name is `vault-dynamodb-backend`. This option can also be provided via the environment variable `AWS_DYNAMODB_TABLE`. If the specified table does not yet exist, it will be created during initialization.

  * `read_capacity` (optional) - The read capacity to provision when creating the DynamoDB table. This is the maximum number of reads consumed per second on the table. The default value is 5. This option can also be provided via the environment variable `AWS_DYNAMODB_READ_CAPACITY`.

  * `write_capacity` (optional) - The write capacity to provision when creating the DynamoDB table. This is the maximum number of writes performed per second on the table. The default value is 5. This option can also be provided via the environment variable `AWS_DYNAMODB_WRITE_CAPACITY`.

  * `access_key` - (required) The AWS access key. It must be provided, but it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable.

  * `secret_key` - (required) The AWS secret key. It must be provided, but it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable.

  * `session_token` - (optional) The AWS session token. It can also be sourced from the `AWS_SESSION_TOKEN` environment variable.

  * `endpoint` - (optional) An alternative (AWS compatible) DynamoDB endpoint to use. It can also be sourced from the `AWS_DYNAMODB_ENDPOINT` environment variable.

  * `region` (optional) - The AWS region. It can be sourced from the `AWS_DEFAULT_REGION` environment variable and will default to "us-east-1" if not specified.

  * `recovery_mode` (optional) - When the Vault leader crashes or is killed without being able to shut down properly, no other node can become the new leader because the DynamoDB table still holds the old leader's lock record. To recover from this situation, one can start a single Vault node with this option set to `1` and the node will remove the old lock from DynamoDB. It is important that only one node is running in recovery mode! After this node has become the leader, other nodes can be started with regular configuration.
    This option can also be provided via the environment variable `RECOVERY_MODE`.

For more information about the read/write capacity of DynamoDB tables, see the [official AWS DynamoDB docs](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithTables.html#ProvisionedThroughput).
If you are running your Vault server on an EC2 instance, you can also make use
of the EC2 instance profile service to provide the credentials Vault will use to
make DynamoDB API calls. Leaving the `access_key` and `secret_key` fields empty
will cause Vault to attempt to retrieve credentials from the metadata service.

#### Backend Reference: S3 (Community-Supported)

For S3, the following options are supported:

  * `bucket` (required) - The name of the S3 bucket to use. It must be provided, but it can also be sourced from the `AWS_S3_BUCKET` environment variable.

  * `access_key` - (required) The AWS access key. It must be provided, but it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable.

  * `secret_key` - (required) The AWS secret key. It must be provided, but it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable.

  * `session_token` - (optional) The AWS session token. It can also be sourced from the `AWS_SESSION_TOKEN` environment variable.

  * `endpoint` - (optional) An alternative (AWS compatible) S3 endpoint to use. It can also be sourced from the `AWS_S3_ENDPOINT` environment variable.

  * `region` (optional) - The AWS region. It can be sourced from the `AWS_DEFAULT_REGION` environment variable and will default to "us-east-1" if not specified.

If you are running your Vault server on an EC2 instance, you can also make use
of the EC2 instance profile service to provide the credentials Vault will use to
make S3 API calls.  Leaving the `access_key` and `secret_key` fields empty
will cause Vault to attempt to retrieve credentials from the metadata service.
You are responsible for ensuring your instance is launched with the appropriate
profile enabled. Vault will handle renewing profile credentials as they rotate.

#### Backend Reference: MySQL (Community-Supported)

The MySQL backend has the following options:

  * `username` (required) - The MySQL username to connect with.

  * `password` (required) - The MySQL password to connect with.

  * `address` (optional) - The address of the MySQL host. Defaults to
    "127.0.0.1:3306.

  * `database` (optional) - The name of the database to use. Defaults to "vault".

  * `table` (optional) - The name of the table to use. Defaults to "vault".

  * `tls_ca_file` (optional) - The path to the CA certificate to connect using TLS

#### Backend Reference: PostgreSQL (Community-Supported)

The PostgreSQL backend has the following options:

  * `connection_url` (required) - The connection string used to connect to PostgreSQL.

    Examples:

    * postgres://username:password@localhost:5432/database?sslmode=disabled

    * postgres://username:password@localhost:5432/database?sslmode=verify-full

    A list of all supported parameters can be found in [the pq library documentation](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters).

  * `table` (optional) - The name of the table to write vault data to. Defaults
    to "vault_kv_store".

Add the following table and index to a new or existing PostgreSQL database:

```sql
CREATE TABLE vault_kv_store (
  parent_path TEXT COLLATE "C" NOT NULL,
  path        TEXT COLLATE "C",
  key         TEXT COLLATE "C",
  value       BYTEA,
  CONSTRAINT pkey PRIMARY KEY (path, key)
);

CREATE INDEX parent_path_idx ON vault_kv_store (parent_path);
```

If you're using a version of PostgreSQL prior to 9.5, create the following
function:

```sql
CREATE FUNCTION vault_kv_put(_parent_path TEXT, _path TEXT, _key TEXT, _value BYTEA) RETURNS VOID AS
$$
BEGIN
    LOOP
        -- first try to update the key
        UPDATE vault_kv_store
          SET (parent_path, path, key, value) = (_parent_path, _path, _key, _value)
          WHERE _path = path AND key = _key;
        IF found THEN
            RETURN;
        END IF;
        -- not there, so try to insert the key
        -- if someone else inserts the same key concurrently,
        -- we could get a unique-key failure
        BEGIN
            INSERT INTO vault_kv_store (parent_path, path, key, value)
              VALUES (_parent_path, _path, _key, _value);
            RETURN;
        EXCEPTION WHEN unique_violation THEN
            -- Do nothing, and loop to try the UPDATE again.
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;
```

More info can be found in the [PostgreSQL documentation](http://www.postgresql.org/docs/9.4/static/plpgsql-control-structures.html#PLPGSQL-UPSERT-EXAMPLE):

#### Backend Reference: Inmem

The in-memory backend has no configuration options.

#### Backend Reference: File

The file backend has the following options:

  * `path` (required) - The path on disk to a directory where the
      data will be stored.
