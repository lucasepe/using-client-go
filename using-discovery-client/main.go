package main

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	rc, err := configLoader.ClientConfig()
	if err != nil {
		panic(err)
	}

	// create a new DiscoveryClient using the given config
	// this client will be used to discover supported resources in the API server
	dc, err := discovery.NewDiscoveryClientForConfig(rc)
	if err != nil {
		panic(err)
	}

	// storage for errors
	errs := []error{}

	// retrieve the supported resources with the version
	// preferred by the server
	lists, err := dc.ServerPreferredResources()
	if err != nil {
		errs = append(errs, err)
	}

	// utility struct holding information to print
	type info struct {
		Kind       string   `json:"kind"`
		APIVersion string   `json:"apiVersion"`
		Name       string   `json:"name"`
		Verbs      []string `json:"verbs"`
	}

	// iterate all the APIResource collections
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}

		// grab the API resource info
		for _, el := range list.APIResources {
			if len(el.Verbs) == 0 {
				continue
			}

			tmp := info{el.Kind, list.GroupVersion, el.Name, el.Verbs}
			// convert to json...
			res, err := json.Marshal(&tmp)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			//..and print
			fmt.Printf("%s\n", res)
		}
	}

	// if there has been an error
	// print it on the screen
	if len(errs) > 0 {
		panic(errors.NewAggregate(errs))
	}
}
