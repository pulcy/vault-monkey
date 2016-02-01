# Vault-monkey

Vault-monkey is an application that extracts secrets from a [vault](https://vaultproject.io).
It is designed to be used in a micro-service environment running autonomously.

## Usage

### Extract one or more secrets as environment variables

`vault-monkey extract env --target <environment-file-path> <key>=<path>[#<field]...`

Example:

`vault-monkey extract env --target /tmp/mysecrets KEY1=/generic/somekey#myfield KEY2=/generic/otherkey`

This command results in a file in `/tmp/mysecrets` containing:

```
KEY1="content of 'myfield' field under '/generic/somekey' path"
KEY2="content of 'value' field under '/generic/otherkey' path"
```


### Extract a secret as a file

`vault-monkey extract file --target <environment-file-path> <path>[#<field]`

Example:

`vault-monkey extract file --target /tmp/myfile /generic/somekey#myfield`

This command results in a file in `/tmp/myfile` containing the content
of the 'myfield' under '/generic/somekey' path.

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

- Path: `/generic/cluster-auth/{cluster-id}/job/{job-id}`
- Field: `user-id`

### Step 2: Job specific login

Once the cluster/job specific user-id is fetched, vault-monkey will perform a second app-id login
using this user-id combined with the job-id (as app-id).

With the token obtained from this second login, vault-monkey will fetch the intended secrets and write
them to file.
