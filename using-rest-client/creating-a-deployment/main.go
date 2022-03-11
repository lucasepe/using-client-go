package main

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	// Determine the namespace referenced by the
	// current context in the kubeconfig file.
	namespace, _, err := configLoader.Namespace()
	if err != nil {
		panic(err)
	}

	cfg, err := configLoader.ClientConfig()
	if err != nil {
		panic(err)
	}

	// POST /apis/apps/v1/namespaces/{namespace}/deployments

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

	// utility function to create a int32 pointer
	i32Ptr := func(i int32) *int32 { return &i }

	// the required request body (a deployment object)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: i32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.21.6",
						},
					},
				},
			},
		},
	}

	// encode the request payload as JSON
	body, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}

	// store the operation result here
	res := &appsv1.Deployment{}

	// fluent interface to setup and perform the
	// POST /apis/apps/v1/namespaces/{namespace}/deployments request
	err = rc.Post().
		Namespace(namespace).
		Resource("deployments").
		Body(body).
		Do(context.TODO()).
		Into(res) // store the result into `res` object
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("deployment.apps/%s created\n", res.Name)
}
