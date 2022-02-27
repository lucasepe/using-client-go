package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	keySupportBy = "app.kubernetes.io/support-by"
	keyPartOf    = "app.kubernetes.io/part-of"
)

func main() {
	// requirements contains values, key, and an operator that relates the key and values
	supportExists, err := labels.NewRequirement(keySupportBy, selection.Exists, []string{})
	if err != nil {
		panic(err)
	}

	supportTeam, err := labels.NewRequirement(keySupportBy, selection.In, []string{"team_1", "team_2"})
	if err != nil {
		panic(err)
	}

	partOfReq, err := labels.NewRequirement(keyPartOf, selection.Equals, []string{"payment-system"})
	if err != nil {
		panic(err)
	}

	//  create a label selector
	selector := labels.NewSelector()

	// adds requirements to the selector
	selector = selector.Add(*supportExists, *supportTeam, *partOfReq)

	// print the selector (it will be rendered as in `kubectl label...`)
	fmt.Printf("here your selector: %s\n", selector.String())
}
