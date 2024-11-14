package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GPU struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GPUSpec   `json:"spec"`
	Status GPUStatus `json:"status"`
}

type GPUSpec struct {
	// Mock data for GPU availability
	NodeGPUs map[string]int `json:"nodeGPUs"`
}

type GPUStatus struct {
	TotalGPUs int `json:"totalGPUs"`
}

// Action input structure
type GpuActionInput struct {
	NodeName string `json:"nodeName"`
	GPUCount int    `json:"gpuCount"`
}
