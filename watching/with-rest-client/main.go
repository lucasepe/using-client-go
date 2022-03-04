package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
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

	cfg.APIPath = "/api"
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err.Error())
	}

	// utility function to create a int64 pointer
	i64Ptr := func(i int64) *int64 { return &i }

	opts := metav1.ListOptions{
		TimeoutSeconds: i64Ptr(120),
		Watch:          true,
	}

	watcher, err := rc.Get().Resource("namespaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(time.Duration(*opts.TimeoutSeconds)).
		Watch(context.TODO())
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
		item := event.Object.(*corev1.Namespace)

		switch event.Type {

		// when a namespace is deleted...
		case watch.Deleted:
			// let's say hello!
			fmt.Printf("- '%s' %v ...bye bye\n", item.GetName(), event.Type)

		// when a namespace is added...
		case watch.Added:
			fmt.Printf("+ '%s' %v  ", item.GetName(), event.Type)

			// try to patch it!
			err = rc.Patch(pt).Resource("namespaces").
				Name(item.Name).
				VersionedParams(&po, scheme.ParameterCodec).
				Body(pd).
				Do(context.TODO()).
				Error()
			if err != nil {
				panic(err)
			}

			fmt.Println(" ...patched!")
		}
	}
}
