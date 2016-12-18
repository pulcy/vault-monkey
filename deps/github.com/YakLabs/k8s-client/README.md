k8s-client [![GoDoc](https://godoc.org/github.com/YakLabs/k8s-client?status.svg)](https://godoc.org/github.com/YakLabs/k8s-client)
=======

k8s-client is a Kubernetes client for Go.

## Status

k8s-client should be considered *alpha* quality.  The API is subject
to change.

k8s-client does not have 100% coverage of the Kubenretes API, but
supports the common operations.  Submit an issue or PR for support of
others.

## Motivation

The "official" Kubernetes client for Go is part of the Kubernetes
project proper and lives within the main Kubernetes repository. It
also requires a rather large number of dependencies.

k8s-client is a small simple API client.

## License

[Apache 2 License](./LICENSE)

k8s-client also includes parts of the Kubernetes code that is also
under the Apache 2 License.

## Usage

See
[![GoDoc](https://godoc.org/github.com/YakLabs/k8s-client?status.svg)](https://godoc.org/github.com/YakLabs/k8s-client)
for documentation.

## Testing

To test the included [Kubernetes http client](./http/), I use
[minikube](https://github.com/kubernetes/minikube) to start a local
cluster. Then I run `kubectl proxy` to proxy to the local cluster
without authentication (this makes testing easier). Then:

```
cd http
go test -v
```

The tests will try to connect to Kubernetes at
`http://127.0.0.1:8001`. This can be overriden using the `K8S_SERVER`
environment value.

## TODO

- [ ] Mock client for testing
- [ ] Better docs/examples
- [ ] Support all the Kubernetes types and operations
- [ ] Support watches
