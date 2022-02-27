package main

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
)

// Writing Kubernetes client apps using Go
//
// Creating a `rest.Config` using default kubeconfig rules.
//
func main() {
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Determine the namespace referenced by the
	// current context in the kubeconfig file.
	namespace, _, err := configLoader.Namespace()
	if err != nil {
		panic(err)
	}

	// `rest.Config`` holds the common attributes that can
	// be passed to a Kubernetes client on initialization.
	cfg, err := configLoader.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Just print some `config` struct variables
	fmt.Printf("Kubernetes Host: %s (current context namespace: %s)\n", cfg.Host, namespace)
}
