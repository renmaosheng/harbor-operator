package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// NotaryServer is the Schema for the notarieservers API
type NotaryServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NotaryServerSpec `json:"spec,omitempty"`

	// Most recently observed status of the Harbor.
	Status ComponentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NotaryServerList contains a list of NotaryServer
type NotaryServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NotaryServer `json:"items"`
}

// NotaryServerSpec defines the desired state of NotaryServer
type NotaryServerSpec struct {
	ComponentSpec         `json:",inline"`
	NotaryServerComponent `json:",inline"`

	// The url exposed to clients to access notary
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*$"
	PublicURL string `json:"publicURL"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*$"
	SignerURL string `json:"signerURL"`

	// +kubebuilder:validation:Required
	DatabaseSecret string `json:"databaseSecret"`

	// +kubebuilder:validation:Required
	CertificateSecret string `json:"certificateSecret"`

	// +kubebuilder:validation:Required
	TokenSecret string `json:"tokenSecret"`
}

// nolint:gochecknoinits
func init() {
	SchemeBuilder.Register(&NotaryServer{}, &NotaryServerList{})
}
