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

	// 集群中所有节点的 GPU 总量
	TotalGPUCount int `json:"totalGPUCount"`

	// 每个节点的 GPU 信息
	NodeGPUInfo []NodeGPUInfo `json:"nodeGPUInfo"`
}

type GPUSpec struct {
	// norman:"nullable" norman 规范
	Enabled bool `json:"enabled" norman:"nullable"`
}

type GPUStatus struct {
	Message string `json:"message"`
}

// 表示每个节点的 GPU 信息
type NodeGPUInfo struct {
	NodeName  string `json:"nodeName"`
	TotalGPU  int    `json:"totalGPU"`
	UsedGPU   int    `json:"usedGPU"`
	UnusedGPU int    `json:"unusedGPU"`
}

// action 入参结构体
type GpuCountActionInput struct {
	// 节点名称，用于筛选特定节点
	NodeName string `json:"nodeName"`
	// 统计维度：total, used, unused
	StatDimension string `json:"statDimension"`
	// 是否启用调试模式
	Debug bool `json:"debug"`
}
