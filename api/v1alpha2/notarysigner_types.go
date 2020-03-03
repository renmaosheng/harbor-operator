package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// NotarySigner is the Schema for the notariesigners API
type NotarySigner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NotarySignerSpec `json:"spec,omitempty"`

	// Most recently observed status of the Harbor.
	Status ComponentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NotarySignerList contains a list of NotarySigner
type NotarySignerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NotarySigner `json:"items"`
}

// NotarySignerSpec defines the desired state of NotarySigner
type NotarySignerSpec struct {
	ComponentSpec         `json:",inline"`
	NotarySignerComponent `json:",inline"`

	// The url exposed to clients to access notary
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*$"
	PublicURL string `json:"publicURL"`

	// +kubebuilder:validation:Required
	DatabaseSecret string `json:"databaseSecret"`

	// +kubebuilder:validation:Required
	CertificateSecret string `json:"certificateSecret"`
}

// nolint:gochecknoinits
func init() {
	SchemeBuilder.Register(&NotarySigner{}, &NotarySignerList{})
}
