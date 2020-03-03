package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ChartMuseumStorageKindKey = "kind"
)

const (
	ChartMuseumCacheURLKey = "url"
)

const (
	ChartMuseumBasicAuthKey = "BASIC_AUTH_PASS"
)

// +kubebuilder:object:root=true

// ChartMuseum is the Schema for the chartmuseums API
type ChartMuseum struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ChartMuseumSpec `json:"spec,omitempty"`

	// Most recently observed status of the Harbor.
	Status ComponentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChartMuseumList contains a list of ChartMuseum
type ChartMuseumList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChartMuseum `json:"items"`
}

// ChartMuseumSpec defines the desired state of ChartMuseum
type ChartMuseumSpec struct {
	ComponentSpec        `json:",inline"`
	ChartMuseumComponent `json:",inline"`

	StorageSecret string `json:"storageSecret,omitempty"`

	CacheSecret string `json:"cacheSecret,omitempty"`

	SecretName string `json:"secret,omitempty"`

	// The url exposed to clients to access ChartMuseum (probably https://.../chartrepo)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*$"
	PublicURL string `json:"publicURL"`
}

// nolint:gochecknoinits
func init() {
	SchemeBuilder.Register(&ChartMuseum{}, &ChartMuseumList{})
}
