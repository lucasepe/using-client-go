package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	cfg, err := configLoader.ClientConfig()
	if err != nil {
		panic(err)
	}

	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	// utility function to create a int64 pointer
	i64Ptr := func(i int64) *int64 { return &i }

	// watch options
	opts := metav1.ListOptions{
		TimeoutSeconds: i64Ptr(120),
		Watch:          true,
	}

	// reference to Namespaces resources
	nsr := dc.Resource(gvr)

	// create the watcher
	watcher, err := nsr.Watch(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	// the patch data, just add a custom label
	pd := []byte(`{"metadata":{"labels":{"modified-by":"lucasepe"}}}`)

	// the patch type
	pt := types.MergePatchType

	// who did this patch?
	po := metav1.PatchOptions{
		FieldManager: "my-cool-app",
	}

	// iterate all the events
	for event := range watcher.ResultChan() {
		// retrieve the Namespace
		item := event.Object.(*unstructured.Unstructured)

		switch event.Type {

		// when a namespace is deleted...
		case watch.Deleted:
			// let's say hello!
			fmt.Printf("- '%s' %v ...bye bye\n", item.GetName(), event.Type)

		// when a namespace is added...
		case watch.Added:
			fmt.Printf("+ '%s' %v  ", item.GetName(), event.Type)

			// try to patch it!
			_, err = nsr.Patch(context.TODO(), item.GetName(), pt, pd, po)
			if err != nil {
				panic(err)
			}

			fmt.Println(" ...patched!")
		}
	}
}
