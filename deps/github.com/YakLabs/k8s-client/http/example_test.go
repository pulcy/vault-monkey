package http_test

import (
	"fmt"
	"os"

	"github.com/YakLabs/k8s-client/http"
)

func ExampleNew() {

	// get server from environment
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

	if token := os.Getenv("K8S_TOKEN"); token != "" {
		opts = append(opts, http.SetToken(token))
	}

	// create a new client using the options
	c, err := http.New(opts...)

	if err != nil {
		// handle the error
	}

	// get a list of all the name spaces
	list, err := c.ListNamespaces(nil)
	if err != nil {
		// handle the error
	}

	for _, v := range list.Items {
		fmt.Println(v.Name)
	}
}

func ExampleNewInCluster() {
	// NewInCluster will use the environment and the service account files to create a new client
	c, err := http.NewInCluster()

	// get a list of all the name spaces
	list, err := c.ListNamespaces(nil)
	if err != nil {
		// handle the error
	}

	for _, v := range list.Items {
		fmt.Println(v.Name)
	}
}
