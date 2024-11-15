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
	// norman:"nullable" norman 规范
	Enabled bool `json:"enabled" norman:"nullable"`
}

type GPUStatus struct {
	Message string `json:"message"`
	// 所有集群中所有节点的 GPU 总量
	TotalGPUCount int `json:"totalGPUCount"`

	// 每个集群的 GPU 信息
	ClusterGPUInfo []ClusterGPUInfo `json:"clusterGPUInfo"`
}

// 表示每个集群的 GPU 信息
type ClusterGPUInfo struct {
	// 集群 ID
	ClusterId string `json:"clusterId"`
	// 当前集群中所有节点的 GPU 总量
	TotalGPUCount int `json:"totalGPUCount"`
	// 当前集群中所有节点的 GPU 信息
	NodeGPUInfo []NodeGPUInfo `json:"nodeGPUInfo"`
}

// 表示每个节点的 GPU 信息
type NodeGPUInfo struct {
	// 节点 ID
	NodeId string `json:"nodeId"`
	// 节点主机名
	NodeHostName string `json:"nodeHostName"`
	// 节点名称
	NodeName string `json:"nodeName"`
	// 总 GPU 数量
	TotalGPU int `json:"totalGPU"`
	// 已使用的 GPU 数量
	UsedGPU int `json:"usedGPU"`
	// 未使用的 GPU 数量
	UnusedGPU int `json:"unusedGPU"`
}

// action 入参结构体
type GpuCountActionInput struct {
	// 节点主机名，用于筛选特定节点
	NodeHostName string `json:"nodeHostName"`
	// 节点名称，用于筛选特定节点
	NodeName string `json:"nodeName"`
	// 节点 ID，用于筛选特定节点
	NodeId string `json:"nodeId"`
	// 集群 ID，用于筛选特定集群
	ClusterId string `json:"clusterId"`
	// 统计维度：total, used, unused
	StatDimension string `json:"statDimension"`
	// 是否启用调试模式
	Debug bool `json:"debug"`
}
