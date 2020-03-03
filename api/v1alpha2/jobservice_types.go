package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	JobServiceSecretKey = "JOBSERVICE_SECRET"
)

const (
	// ipaddress:port[,weight,password,database_index]
	JobServiceRedisURLKey       = "JOB_SERVICE_POOL_REDIS_URL"
	JobServiceRedisNamespaceKey = "JOB_SERVICE_POOL_REDIS_NAMESPACE"
)

const (
	JobServiceRegistryControllerURLKey = "REGISTRY_CONTROLLER_URL"
)

// +kubebuilder:object:root=true

// JobService is the Schema for the jobservices API
type JobService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec JobServiceSpec `json:"spec,omitempty"`

	// Most recently observed status of the Harbor.
	Status ComponentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JobServiceList contains a list of JobService
type JobServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JobService `json:"items"`
}

// JobServiceSpec defines the desired state of JobService
type JobServiceSpec struct {
	ComponentSpec    `json:",inline"`
	JobServiceConfig `json:",inline"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="INFO,DEBUG,WARNING,ERROR,FATAL"
	LogLevel string `json:"logLevel,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="INFO,DEBUG,WARNING,ERROR,FATAL"
	PublicLogLevel string `json:"publicLogLevel"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*$"
	CoreURL string `json:"coreURL"`
}

type JobServiceConfig struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	WorkerCount int32 `json:"workerCount"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	WebHookMaxRetry int32 `json:"webHookMaxRetry"`

	// +kubebuilder:validation:Required
	CoreSecret string `json:"coreSecret"`

	// +kubebuilder:validation:Required
	RedisSecret string `json:"redisSecret"`
}

// nolint:gochecknoinits
func init() {
	SchemeBuilder.Register(&JobService{}, &JobServiceList{})
}
