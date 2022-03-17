package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/PaesslerAG/gval"

	expressionV1alpha1Api "github.com/lucasepe/using-client-go/using-codegen/pkg/apis/expression/v1alpha1"
	expressionV1alpha1Clientset "github.com/lucasepe/using-client-go/using-codegen/pkg/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
		defaultKubeconfig, "absolute path to the kubeconfig file")

	namespace := flag.String("namespace", metav1.NamespaceDefault, "create the deployment in this namespace")

	flag.Usage = func() {
		name := os.Args[0]
		if strings.Contains(name, "go-build") {
			name = "go run main.go"
		}

		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage:  %s <Expression Custom Resource Name>\n\n", name)

		fmt.Fprintf(w, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// build the config from the specified kubeconfig filepath
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	cs, err := expressionV1alpha1Clientset.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	res, err := cs.ExampleV1alpha1().Expressions(*namespace).
		Get(context.TODO(), flag.Args()[0], metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	val, err := evalExpression(res)
	if err != nil {
		panic(err)
	}
	res.Status.Result = strval(val)

	_, err = cs.ExampleV1alpha1().Expressions(*namespace).
		UpdateStatus(context.TODO(), res, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Expression evaluated! Type 'kubectl get exp %s' to check the result.\n", flag.Args()[0])
}

func evalExpression(src *expressionV1alpha1Api.Expression) (string, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(src.Spec.Data), &data)
	if err != nil {
		return err.Error(), err
	}

	val, err := gval.Evaluate(src.Spec.Body, data)
	if err != nil {
		return err.Error(), err
	}

	return strval(val), nil
}

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
