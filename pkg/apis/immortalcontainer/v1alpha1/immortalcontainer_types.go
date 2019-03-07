package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ImmortalContainerSpec defines the desired state of ImmortalContainer
// +k8s:openapi-gen=true
type ImmortalContainerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
}

// ImmortalContainerStatus defines the observed state of ImmortalContainer
// +k8s:openapi-gen=true
type ImmortalContainerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	CurrentPod string `json:"currentPod,omitempty"`
	StartTimes int `json:"startTimes,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImmortalContainer is the Schema for the immortalcontainers API
// +k8s:openapi-gen=true
type ImmortalContainer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImmortalContainerSpec   `json:"spec,omitempty"`
	Status ImmortalContainerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImmortalContainerList contains a list of ImmortalContainer
type ImmortalContainerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImmortalContainer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImmortalContainer{}, &ImmortalContainerList{})
}
