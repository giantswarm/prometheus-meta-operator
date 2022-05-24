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
	Spec              RemoteWriteSpec `json:"spec"`
}

type RemoteWriteSpec struct {
	RemotWrite      pov1.RemoteWriteSpec `json:"remoteWrite"`
	ClusterSelector metav1.LabelSelector `json:"clusterSelector"`
}

type RemoteWriteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RemoteWrite `json:"items"`
}
