# Vault-monkey

Vault-monkey is an application that extracts secrets from a [vault](https://vaultproject.io).
It is designed to be used in a micro-service environment, running autonomously.

Vault-monkey formats the extracted secrets as individual files, or as key-value pairs
formatted into an environment file.

All extract functions use a 2 step server login, that is designed to simplify management of
large clusters of machines, where access policies are organized per cluster-job pair
instead of machine-job pair. See [authentication scheme](#authentication-scheme).

## Usage

### Extracting secrets

#### Extract one or more secrets as environment variables

`vault-monkey extract env --target <environment-file-path> <key>=<path>[#<field]...`

Example:

`vault-monkey extract env --target /tmp/mysecrets KEY1=/secret/somekey#myfield KEY2=/secret/otherkey`

This command results in a file in `/tmp/mysecrets` containing:

```
KEY1=content of 'myfield' field under '/secret/somekey' path
KEY2=content of 'value' field under '/secret/otherkey' path
```

#### Extract a secret as a file

`vault-monkey extract file --target <file-path> <path>[#<field]`

Example:

`vault-monkey extract file --target /tmp/myfile /secret/somekey#myfield`

This command results in a file in `/tmp/myfile` containing the content
of the 'myfield' under '/secret/somekey' path.


### Operational commands

Operations can use vault-monkey to prepare the vault for the 2 step authentication using several
`cluster` and `job` commands.

To create a new cluster, use:

```
vault-monkey cluster create -G <github-token> --cluster-id <cluster-id>
```

This will automatically create a policy needed for step 1 of the authentication scheme.

To add a machine to a cluster, use:

```
vault-monkey cluster add -G <github-token> --cluster-id <cluster-id> --machine-id <machine-id>
```

To remove a machine from a cluster, use:

```
vault-monkey cluster remove -G <github-token> --cluster-id <cluster-id> --machine-id <machine-id>
```

To create a new job, use:

```
vault-monkey job create -G <github-token> --job-id <cluster-id> --policy <policy-name>
```

To allow a cluster to access secrets for a job, use:

```
vault-monkey job allow -G <github-token> --job-id <job-id> --cluster-id <cluster-id>
```

To deny a cluster to access secrets for a job, use:

```
vault-monkey job deny -G <github-token> --job-id <job-id> --cluster-id <cluster-id>
```

To remove a job, use:

```
vault-monkey job delete -G <github-token> --job-id <cluster-id>
```

Note that deleting a job does not remove all cluster grants.

To show the seal status of all instances of a vault, use:

```
vault-monkey seal-status -G <github-token>
```

To seal a vault, use:

```
vault-monkey seal -G <github-token>
```

To unseal a vault, use:

```
vault-monkey unseal -G <github-token> <script-to-fetch-a-key> [script argument]...
```

The script to fetch an unseal key will be executed several times (how many depends on the unseal threshold).
The arguments of the script will be processed as a go template with `{{.Key}}` as the number of the key
to extract. This value can be 1..N where N is the unseal threshold.

E.g. if you use [`pass`](http://www.passwordstore.org/) to store your unseal keys, use something like this:

```
vault-monkey unseal -G <github-token> pass show MyVault/UnsealKey{{.Key}}
```

This will fetch keys from your password-store with path `MyVaultUnsealKey1`, `MyVaultUnsealKey2` etc.
Note that vault-monkey will shuffle the keys, so if your vault has 5 unseal keys with a threshold of 3
if may ask for key3, key1, key5.

## Authentication Scheme

Vault-monkey is designed to function in an environment with lots of servers, running lots of different
jobs, without fixed constraints about which job run of which server(s).

With lots of changing servers, it is not nice to configure something per server/job pair.
If that would be the case then adding/removing one server would result in changing a lot of these
pairs. The same is true for adding/removing a single job.

To avoid this, vault-monkey is build around a 2 step authentication process.

It assumes that all servers in a cluster are allowed to access data for all jobs that are intended
to run on that cluster.

### Step 1: Cluster membership

The first step during authentication is to establish cluster membership. It does so by trying to
login with a cluster-id combined with a machine-id.
The cluster-id is pass to the machine during provisioning and must be the same for all machines
in the cluster.
The machine-id is created during the first-boot of the machine and must remain the same throughout
the lifetime of the machine.

It uses the app-id authentication for this, where the cluster-id becomes the app-id and the  
machine-id becomes the user-id.

If thirst first login in successful, vault-monkey will read a user-id which is specific per
cluster/job pair.

This pair must be stored under:

- Path: `/secret/cluster-auth/{cluster-id}/job/{job-id}`
- Field: `user-id`

### Step 2: Job specific login

Once the cluster/job specific user-id is fetched, vault-monkey will perform a second app-id login
using this user-id combined with the job-id (as app-id).

With the token obtained from this second login, vault-monkey will fetch the intended secrets and write
them to file.

### Security notes

#### Note 1

It may be possible that a machine stores the `user-id` it fetches in step 1 longer than it should.
In that case this machine will be able to access secrets for the configured jobs even after it has been
removed from the cluster.

If that is the case, replace the `user-id` by running `vault-monkey job allow ...` again.

#### Note 2

The primary use of vault-monkey is to extract secrets from the vault.
This will result in files in your filesystem. To make sure these secrets do not survive a reboot,
use a directory that is mounted on non-persistent storage.

## Vault policies

Vault-monkey will automatically create a policy for step 1 of the authentication.
To allow your operations team to execute all the [operational commands](#operational-commands)
use a policy like this:

```
// Allow operations to seal the vault
path "sys/seal" {
    policy = "sudo"
}

// Allow operations to configure app-id's
path "auth/app-id/*" {
    policy = "write"
}

// Allow operations to create 2 step cluster authentication policies
path "sys/policy/cluster_auth_*" {
    policy = "write"
}

// Allow operations to access all normal secrets
// This is not needed for vault-monkey, but it likely to be convenient.
path "secret/*" {
    policy = "write"
}
```

## Building

To build vault-monkey, run:

```
make
```

This will setup a local `GOPATH` and run a docker container to build vault-monkey.
