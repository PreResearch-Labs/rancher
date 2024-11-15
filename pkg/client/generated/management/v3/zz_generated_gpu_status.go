package client

const (
	GPUStatusType                = "gpuStatus"
	GPUStatusFieldClusterGPUInfo = "clusterGPUInfo"
	GPUStatusFieldMessage        = "message"
	GPUStatusFieldTotalGPUCount  = "totalGPUCount"
)

type GPUStatus struct {
	ClusterGPUInfo []ClusterGPUInfo `json:"clusterGPUInfo,omitempty" yaml:"clusterGPUInfo,omitempty"`
	Message        string           `json:"message,omitempty" yaml:"message,omitempty"`
	TotalGPUCount  int64            `json:"totalGPUCount,omitempty" yaml:"totalGPUCount,omitempty"`
}
