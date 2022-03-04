package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiWatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/watch"
)

// sentinel is an object that knows how to
// start a watch on namespaces resources
//
// this is our implementation of `cache.Watcher`
type sentinel struct {
	client      *kubernetes.Clientset
	timeoutSecs int64
}

// newSentinel returns a new `sentinel` object that implements `cache.Watcher`
func newSentinel(cs *kubernetes.Clientset, timeout int64) cache.Watcher {
	return &sentinel{cs, timeout}
}

// Watch begin a watch on namespaces resources
func (s *sentinel) Watch(options metav1.ListOptions) (apiWatch.Interface, error) {
	return s.client.CoreV1().Namespaces().
		Watch(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &s.timeoutSecs,
		})
}

// just to be sure that `cache.Watcher` interface
// is being implemented by our `sentinel` struct type
var _ cache.Watcher = (*sentinel)(nil)

func main() {
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	cfg, err := configLoader.ClientConfig()
	if err != nil {
		panic(err)
	}

	// create a new `Clientset`` for the given config
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("---- Start watching namespaces ----\n")

	// create a `cache.Watcher` implementation using the `ClientSet``
	watcher := newSentinel(cs, 50)
	// create a `RetryWatcher` using initial
	// version "1" and our specialized watcher
	rw, err := watch.NewRetryWatcher("1", watcher)
	if err != nil {
		panic(err)
	}

	// process incoming event notifications
	for {
		// grab the event object
		event, ok := <-rw.ResultChan()
		if !ok {
			panic(fmt.Errorf("closed channel"))
		}

		// cast to namespace
		ns, ok := event.Object.(*corev1.Namespace)
		if !ok {
			panic(fmt.Errorf("invalid type '%T'", event.Object))
		}

		// skip events older then five minutes
		creationTime := ns.GetCreationTimestamp().Time
		fiveMinsAgo := time.Now().Add(-5 * time.Minute)
		if event.Type == apiWatch.Added && creationTime.Before(fiveMinsAgo) {
			// fmt.Printf(">> skip older events (creationTime: %s, currentTime: %s)\n",
			//	creationTime.Format(time.RFC3339), time.Now().Format(time.RFC3339))
			continue
		}

		// print some info about the event
		fmt.Printf("%s %s (createdAt: %s, phase: %s)\n",
			event.Type, ns.Name, creationTime.Format(time.RFC3339), ns.Status.Phase)

		// sleep a bit
		time.Sleep(5 * time.Second)
	}
}
