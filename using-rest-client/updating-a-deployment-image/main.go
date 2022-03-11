package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
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

	// create a RESTClient instance, using the the configuration object as input parameter
	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err.Error())
	}

	// store the operation result here
	res := &appsv1.Deployment{}

	// fluent interface to setup and perform the request
	// GET /apis/apps/v1/namespaces/{namespace}/deployments/{name}
	err = rc.Get().
		Namespace(namespace).
		Resource("deployments").
		Name("nginx").
		Do(context.TODO()).
		Into(res) // store the result into `res` object
	if err != nil {
		panic(err.Error())
	}

	// print the current image (we know that there is only one container)
	fmt.Printf("before patching: deployment.apps/%s image is %s\n",
		res.Name, res.Spec.Template.Spec.Containers[0].Image)

	// the JSON payload to partially update the specified deployment
	patch := []byte(`{"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"nginx:1.14.2"}]}}}}`)

	// fluent interface to setup and perform the request
	// PATCH /apis/apps/v1/namespaces/{namespace}/deployments/{name}
	err = rc.Patch(types.StrategicMergePatchType).
		Namespace(namespace).
		Resource("deployments").
		Name("nginx").
		Body(patch).
		Do(context.TODO()).
		Into(res) // store the result into `res` object
	if err != nil {
		panic(err.Error())
	}

	// print the deployment image after the patch
	fmt.Printf("after  patching: deployment.apps/%s image is %s\n",
		res.Name, res.Spec.Template.Spec.Containers[0].Image)
}
