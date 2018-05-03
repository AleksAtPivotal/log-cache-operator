package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterPhase string

const (
	ClusterPhaseInitial ClusterPhase = ""
	ClusterPhaseRunning              = "Running"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LogCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []LogCache `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCache defines the CRD LogCache object
type LogCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              LogCacheSpec   `json:"spec"`
	Status            LogCacheStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCacheSpec contains the spec for Log Cache CRD Object
type LogCacheSpec struct {
	LogCachePod       LogCachePodSpec       `json:"logcache"`
	LogCacheNozzle    LogCacheNozzleSpec    `json:"logcachenozzle"`
	LogCacheScheduler LogCacheSchedulerSpec `json:"logcachescheduler"`
	// TLS policy of  nodes
	TLS *TLSPolicy `json:"TLS,omitempty"`
}

// LogCachePodSpec contains the specs for LogCache Pod
type LogCachePodSpec struct {
	// Number of nodes to deploy for a LogCache deployment.
	// Default: 1.
	Nodes int32 `json:"nodes,omitempty"`

	// Base image to use for a Log Cache deployment.
	BaseImage string `json:"baseImage"`

	// Version of Log Cache Nozzle to be deployed.
	Version string `json:"version"`

	// Base image to use for a Log Cache deployment.
	GatewayBaseImage string `json:"GatewayBaseImage"`

	// Version of Log Cache Nozzle to be deployed.
	GatewayVersion string `json:"Gatewayversion"`
}

// LogCacheNozzleSpec contains the specs for LogCache Nozzle Deployment
type LogCacheNozzleSpec struct {
	// Number of nodes to deploy for a LogCache Nozzledeployment.
	// Default: 1.
	Nodes int32 `json:"nodes,omitempty"`

	// Base image to use for a Log Cache Nozzle deployment.
	BaseImage string `json:"baseImage"`

	// Version of Log Cache Nozzle to be deployed.
	Version string `json:"version"`
}

// LogCacheSchedulerSpec contains the specs for LogCache Scheduler Deployment
type LogCacheSchedulerSpec struct {
	// Number of nodes to deploy for a LogCache Schedulerdeployment.
	// Default: 1.
	Nodes int32 `json:"nodes,omitempty"`

	// Base image to use for a Log Cache Scheduler deployment.
	BaseImage string `json:"baseImage"`

	// Version of Log Cache Scheduler to be deployed.
	Version string `json:"version"`
}

// LogCacheStatus contains the status of Log Cache CRD Object
type LogCacheStatus struct {
	// Phase indicates the state this Vault cluster jumps in.
	// Phase goes as one way as below:
	//   Initial -> Running
	Phase ClusterPhase `json:"phase"`

	// Initialized indicates if the Vault service is initialized.
	Initialized bool `json:"initialized"`

	// ServiceName is the LB service for accessing vault nodes.
	ServiceName string `json:"serviceName,omitempty"`

	State string `json:"state"`
}
