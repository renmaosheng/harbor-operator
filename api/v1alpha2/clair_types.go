package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// Clair is the Schema for the clairs API
type Clair struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClairSpec `json:"spec,omitempty"`

	// Most recently observed status of the Harbor.
	Status ComponentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClairList contains a list of Clair
type ClairList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Clair `json:"items"`
}

// ClairSpec defines the desired state of Clair
type ClairSpec struct {
	ComponentSpec  `json:",inline"`
	ClairComponent `json:",inline"`

	// +kubebuilder:validation:Required
	DatabaseSecret string `json:"databaseSecret"`

	// +kubebuilder:validation:Optional
	VulnerabilitySources []string `json:"vulnerabilitySources"`

	// +kubebuilder:validation:Required
	Adapter Adapter `json:"adapter"`
}

type Adapter struct {
	// +kubebuilder:validation:Required
	RedisSecret string `json:"redisSecret"`
}

// nolint:gochecknoinits
func init() {
	SchemeBuilder.Register(&Clair{}, &ClairList{})
}
