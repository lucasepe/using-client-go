package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

/*
kubectl create namespace be-system
kubectl create namespace fe-system

kubectl label namespaces fe-system team=frontend
kubectl label namespaces be-system team=backend

$ go run labels-and-selectors/list-pods-using-label-selectors/main.go
be-system
default
fe-system
kube-node-lease
kube-public
kube-system
local-path-storage

$ go run labels-and-selectors/list-pods-using-label-selectors/main.go team=frontend

$ go run labels-and-selectors/list-pods-using-label-selectors/main.go 'team notin (frontend,backend)'
*/
func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	flag.Parse()

	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	dc, err := dynamic.NewForConfig(rc)
	if err != nil {
		panic(err)
	}

	var filter string
	if len(flag.Args()) > 0 {
		sel, err := labels.Parse(strings.Join(flag.Args(), " "))
		if err != nil {
			panic(err)
		}
		filter = sel.String()
	}

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	res, err := dc.Resource(gvr).
		List(context.TODO(), metav1.ListOptions{LabelSelector: filter})
	if err != nil {
		panic(err)
	}

	/*
		var list corev1.NamespaceList
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(res.UnstructuredContent(), &list)
		if err != nil {
			panic(err)
		}
	*/
	for _, el := range res.Items {
		fmt.Println(el.GetName())
	}
}
