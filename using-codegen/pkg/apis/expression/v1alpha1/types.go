package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Expression is our custom resource type. A client is created for it.
type Expression struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ExpressionSpec `json:"spec"`

	// +optional
	Status ExpressionStatus `json:"status,omitempty"`
}

// ExpressionSpec defines the specs for our custom resource
type ExpressionSpec struct {
	// Body is the expression.
	Body string `json:"body"`

	// Data is the JSON with all the expression variables values-
	Data string `json:"data,omitempty"`
}

// ExpressionStatus define the status of our custom resource
type ExpressionStatus struct {
	// Result contains the result of expression evaluation
	Result string `json:"result"`
}

// HasChanged is an utility function to check if two expression are equals.
func (e *Expression) HasChanged(other *Expression) bool {
	return e.Spec.Body != other.Spec.Body || e.Spec.Data != other.Spec.Data
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExpressionList is a top-level list type. The client methods for lists are automatically created.
// You are not supposed to create a separated client for this one.
type ExpressionList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Expression `json:"items"`
}
