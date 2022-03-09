package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	flag.Parse()

	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create a client set from config
	clientSet, err := kubernetes.NewForConfig(rc)
	if err != nil {
		panic(err.Error())
	}

	// create a new instance of sharedInformerFactory for all namespaces
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Minute*1)

	// using this factory create an informer for `secret` resources
	secretsInformer := informerFactory.Core().V1().Secrets()

	// adds an event handler to the shared informer
	secretsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			item := obj.(*corev1.Secret)
			fmt.Printf("secret added (ns=%s): %s\n", item.GetNamespace(), item.GetName())
		},

		UpdateFunc: func(old, new interface{}) {
			item := old.(*corev1.Secret)
			fmt.Printf("secret updated (ns=%s): %s\n", item.GetNamespace(), item.GetName())
		},

		DeleteFunc: func(obj interface{}) {
			item := obj.(*corev1.Secret)
			fmt.Printf("secret deleted (ns=%s): %s\n", item.GetNamespace(), item.GetName())
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	// starts the shared informers that have been created by the factory
	informerFactory.Start(stopCh)

	// wait for the initial synchronization of the local cache
	if !cache.WaitForCacheSync(stopCh, secretsInformer.Informer().HasSynced) {
		panic("failed to sync")
	}

	// causes the goroutine to block (hit CTRL+C to exit)
	select {}
}
