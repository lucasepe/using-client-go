package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// retrieve the value of the KUBECONFIG environment variable
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	// if KUBECONFIG is empty
	if len(defaultKubeconfig) == 0 {
		// look for file $HOME/.kube/config
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	flag.Parse()

	// build the config from the specified kubeconfig filepath
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// Just print some `config` struct variables
	fmt.Printf("Kubernetes Host: %v\n", cfg.Host)
}
