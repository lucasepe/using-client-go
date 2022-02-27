package main

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// $ kubectl apply -f using-dynamic-interface/get-and-update-crds/pizza_crd.yaml
// $ kubectl apply -f using-dynamic-interface/get-and-update-crds/margherita.yaml
// $ kubectl get piz
// go run using-dynamic-interface/get-and-update-crds/main.go
// $ kubectl get piz
func main() {
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	namespace, _, err := configLoader.Namespace()
	if err != nil {
		panic(err)
	}

	cfg, err := configLoader.ClientConfig()
	if err != nil {
		panic(err)
	}

	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	// identify out custom resource
	gvr := schema.GroupVersionResource{
		Group:    "bella.napoli.it",
		Version:  "v1alpha1",
		Resource: "pizzas",
	}
	// retrieve the resource of kind Pizza named 'margherita'
	res, err := dc.Resource(gvr).
		Namespace(namespace).
		Get(context.TODO(), "margherita", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return
		}
		panic(err)
	}

	// grab the status if exists
	status, ok := res.Object["status"]
	if !ok {
		// otherwise create it
		status = make(map[string]interface{})
	}

	// change the 'margherita' price
	status.(map[string]interface{})["cost"] = 6.50
	res.Object["status"] = status

	// update the 'margherita' custom resource with the new price
	_, err = dc.Resource(gvr).Namespace(namespace).Update(context.TODO(), res, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}
}
