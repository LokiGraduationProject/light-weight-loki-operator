package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// CanarySpec defines the desired state of Canary
type CanarySpec struct {
	// Name is the name of the Canary deployment
	Name string `json:"name,omitempty"`

	// Image is the container image to use for the Canary deployment
	Image string `json:"image,omitempty"`

	// Addr is the address of the Loki service
	Addr string `json:"addr,omitempty"`

	// Port is the port for the service
	Port int32 `json:"port,omitempty"`

	// Labels are additional labels to apply to the Canary resources
	DaemonSetLabels []Label `json:"daemonSetLabels,omitempty"`

	// PodLabels are additional labels to apply to the Canary pods
	PodLabels []Label `json:"podLabels,omitempty"`

	// PodAnnotations are additional annotations to apply to the Canary pods
	PodAnnotations []Label `json:"podAnnotations,omitempty"`

	// Number of buckets in the response_latency histogram (default 10)
	Buckets int32 `json:"buckets,omitempty"`
	// Number of concurrent queries to run (default 1)
	TenantID string `json:"tenantId,omitempty"`

	// Logging parameters
	LabelName         string `json:"labelName,omitempty"`
	LabelValue        string `json:"labelValue,omitempty"`
	StreamName        string `json:"streamName,omitempty"`
	StreamValue       string `json:"streamValue,omitempty"`
	Size              int32  `json:"size,omitempty"`
	OutOfOrderMax     string `json:"outOfOrderMax,omitempty"`
	OutOfOrderMin     string `json:"outOfOrderMin,omitempty"`
	OutOfOrderPercent int32  `json:"outOfOrderPercentage,omitempty"`

	// Timing and interval configurations
	Interval             string `json:"interval,omitempty"`
	MaxWait              string `json:"maxWait,omitempty"`
	MetricTestInterval   string `json:"metricTestInterval,omitempty"`
	MetricTestRange      string `json:"metricTestRange,omitempty"`
	QueryTimeout         string `json:"queryTimeout,omitempty"`
	SpotCheckInitialWait string `json:"spotCheckInitialWait,omitempty"`
	SpotCheckInterval    string `json:"spotCheckInterval,omitempty"`
	SpotCheckMax         string `json:"spotCheckMax,omitempty"`
	SpotCheckQueryRate   string `json:"spotCheckQueryRate,omitempty"`
	PruneInterval        string `json:"pruneInterval,omitempty"`
	WaitDuration         string `json:"waitDuration,omitempty"`
	WriteMaxBackoff      string `json:"writeMaxBackoff,omitempty"`
	WriteMinBackoff      string `json:"writeMinBackoff,omitempty"`
	WriteTimeout         string `json:"writeTimeout,omitempty"`
	WriteMaxRetries      int32  `json:"writeMaxRetries,omitempty"`
	Push                 bool   `json:"push,omitempty"`
}

// CanaryStatus defines the observed state of Canary
type CanaryStatus struct {
	// Define observed state fields here
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Canary is the Schema for the canaries API
type Canary struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CanarySpec   `json:"spec,omitempty"`
	Status CanaryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CanaryList contains a list of Canary
type CanaryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Canary `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Canary{}, &CanaryList{})
}
