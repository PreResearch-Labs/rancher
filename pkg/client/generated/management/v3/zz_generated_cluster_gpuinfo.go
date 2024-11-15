package client

const (
	ClusterGPUInfoType               = "clusterGPUInfo"
	ClusterGPUInfoFieldClusterName   = "clusterName"
	ClusterGPUInfoFieldNodeGPUInfo   = "nodeGPUInfo"
	ClusterGPUInfoFieldTotalGPUCount = "totalGPUCount"
)

type ClusterGPUInfo struct {
	ClusterName   string        `json:"clusterName,omitempty" yaml:"clusterName,omitempty"`
	NodeGPUInfo   []NodeGPUInfo `json:"nodeGPUInfo,omitempty" yaml:"nodeGPUInfo,omitempty"`
	TotalGPUCount int64         `json:"totalGPUCount,omitempty" yaml:"totalGPUCount,omitempty"`
}
