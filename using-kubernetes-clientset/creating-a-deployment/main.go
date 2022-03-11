package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	// create the deployment in the specified namespace
	res, err := cs.AppsV1().Deployments(*namespace).
		Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("deployment.apps/%s created\n", res.Name)
}
