package v1alpha1

import (
	pov1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:resource:categories=common;giantswarm
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// RemoteWrite represents schema for managed RemoteWrites in Prometheus. Reconciled by prometheus-meta-operator.
type RemoteWrite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              RemoteWriteSpec   `json:"spec"`
	Status            RemoteWriteStatus `json:"status"`
}

type RemoteWriteSpec struct {
	RemoteWrite pov1.RemoteWriteSpec `json:"remoteWrite"`
	// +immutable
	ClusterSelector metav1.LabelSelector `json:"clusterSelector"`

	// Secrets data to be created along with the configured Prometheus resource.
	// This provides the data for any v1.SecretKeySelector used in the subsequent RemoteWrite field.
	// Provided name and keys should match values in v1.SecretKeySelector fields.
	// +optional
	// +immutable
	Secrets []RemoteWriteSecretSpec `json:"secrets,omitempty"`
}

type RemoteWriteSecretSpec struct {
	Name string            `json:"name"`
	Data map[string][]byte `json:"data,omitempty"`
}

type RemoteWriteStatus struct {
	ConfiguredPrometheuses []RemoteWriteStatusConfiguredPrometheus `json:"configuredPrometheuses"`
}

type RemoteWriteStatusConfiguredPrometheus struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RemoteWriteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RemoteWrite `json:"items"`
}
