package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	namespace := flag.String("namespace", metav1.NamespaceDefault, "create the deployment in this namespace")

	flag.Parse()

	// build the config from the specified kubeconfig filepath
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// creates a new Clientset for the given config
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// the list of Pods in the namespace
	res, err := cs.CoreV1().Pods(*namespace).List(context.TODO(), metav1.ListOptions{})
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
