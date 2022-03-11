package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

	// get the 'nginx' deployment in the namespace
	res, err := cs.AppsV1().Deployments(*namespace).
		Get(context.TODO(), "nginx", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// print the current image (we know that there is only one container)
	fmt.Printf("before patching: deployment.apps/%s image is %s\n",
		res.Name, res.Spec.Template.Spec.Containers[0].Image)

	// the JSON payload to partially update the specified deployment
	patch := []byte(`{"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"nginx:1.20.2"}]}}}}`)

	// apply the patch
	res, err = cs.AppsV1().Deployments(*namespace).
		Patch(context.TODO(), "nginx", types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		panic(err.Error())
	}

	// print the deployment image after the patch
	fmt.Printf("after  patching: deployment.apps/%s image is %s\n",
		res.Name, res.Spec.Template.Spec.Containers[0].Image)
}
