package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// Tracer implements http.RoundTripper.  It prints each request and
// response/error to os.Stderr.  WARNING: this may output sensitive information
// including bearer tokens.
type Tracer struct {
	http.RoundTripper
}

// RoundTrip calls the nested RoundTripper while printing each request and
// response/error to os.Stderr on either side of the nested call.  WARNING: this
// may output sensitive information including bearer tokens.
func (t *Tracer) RoundTrip(req *http.Request) (*http.Response, error) {
	// Dump the request to os.Stderr.
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	os.Stderr.Write(b)
	os.Stderr.Write([]byte{'\n'})

	// Call the nested RoundTripper.
	resp, err := t.RoundTripper.RoundTrip(req)

	// If an error was returned, dump it to os.Stderr.
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return resp, err
	}

	// Dump the response to os.Stderr.
	b, err = httputil.DumpResponse(resp, req.URL.Query().Get("watch") != "true")
	if err != nil {
		return nil, err
	}
	os.Stderr.Write(b)
	os.Stderr.Write([]byte{'\n'})

	return resp, err
}

// go run main.go -namespace kube-system
func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	namespace := flag.String("namespace", metav1.NamespaceDefault,
		"create the deployment in this namespace")

	verbose := flag.Bool("verbose", false, "display HTTP calls")

	flag.Parse()

	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	if *verbose {
		// wrap the default RoundTripper with an instance of Tracer.
		rc.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return &Tracer{rt}
		}
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
