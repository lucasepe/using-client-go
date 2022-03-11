package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

	// the base API path "/apis"
	cfg.APIPath = "apis"
	// the Deployment group and version "/apps/v1"
	cfg.GroupVersion = &appsv1.SchemeGroupVersion
	// specify the serializer
	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	// create a RESTClient instance, using the the
	// configuration object as input parameter
	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err.Error())
	}

	// operation result
	res := &metav1.Status{}

	// fluent interface to setup and perform the request
	// DELETE /apis/apps/v1/namespaces/{namespace}/deployments/{name}
	err = rc.Delete().
		Namespace(namespace).
		Resource("deployments").
		Name("nginx").
		Do(context.TODO()).
		Into(res)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("%s\n", err.Error())
			return
		}
		panic(err.Error())
	}

	fmt.Printf("deployment.apps \"nginx\" delete: %s\n", res.Status)
}
