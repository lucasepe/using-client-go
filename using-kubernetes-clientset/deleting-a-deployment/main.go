package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	namespace := flag.String("namespace", metav1.NamespaceDefault, "create the deployment in this namespace")

	flag.Parse()

	// build the config from the specified kubeconfig filepath
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// creates a new Clientset for the given config
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// delete the 'nginx' deployment in the namespace
	err = cs.AppsV1().Deployments(*namespace).
		Delete(context.TODO(), "nginx", metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return
		}
		panic(err.Error())
	}

	fmt.Println("deployment.apps \"nginx\" deleted")
}
