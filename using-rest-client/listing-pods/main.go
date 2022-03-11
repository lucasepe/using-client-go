package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	corev1 "k8s.io/api/core/v1"
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

	// the base API path "/api" (legacy resource)
	cfg.APIPath = "api"
	// the Pod group and version "/v1" (group name is empty for legacy resources)
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	// specify the serializer
	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	// create a RESTClient instance, using the the
	// configuration object as input parameter
	rc, err := rest.RESTClientFor(cfg)
	if err != nil {
		panic(err.Error())
	}

	// the list of Pods (the result)
	res := &corev1.PodList{}

	// fluent interface to setup and perform
	// the GET /api/v1/pods request
	err = rc.Get().
		Namespace(namespace).
		Resource("pods").
		Do(context.TODO()).
		Into(res)
	if err != nil {
		panic(err.Error())
	}

	// print the results on the terminal in the form of a table
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', 0)

	dorow := func(cells []string) {
		fmt.Fprintln(w, strings.Join(cells, "\t"))
	}

	dorow([]string{"NAME", "STATUS", "AGE"})

	for _, p := range res.Items {
		age := time.Since(p.CreationTimestamp.Time).Round(time.Second)
		dorow([]string{p.Name, string(p.Status.Phase), fmt.Sprintf("%dm", int(age.Minutes()))})
	}

	w.Flush()
}
