package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create a new dynamic client using the rest.Config
	dc, err := dynamic.NewForConfig(rc)
	if err != nil {
		panic(err.Error())
	}

	// identify pods resource
	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}

	// list all pods in the specified namespace
	res, err := dc.Resource(gvr).
		Namespace(*namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			panic(err)
		}
	}

	// for each pod, print just the name
	for _, el := range res.Items {
		fmt.Printf("%v\n", el.GetName())
	}
}
