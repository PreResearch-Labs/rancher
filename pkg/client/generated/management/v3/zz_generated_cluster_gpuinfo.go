package client

const (
	ClusterGPUInfoType               = "clusterGPUInfo"
	ClusterGPUInfoFieldClusterId     = "clusterId"
	ClusterGPUInfoFieldNodeGPUInfo   = "nodeGPUInfo"
	ClusterGPUInfoFieldTotalGPUCount = "totalGPUCount"
)

type ClusterGPUInfo struct {
	ClusterId     string        `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	NodeGPUInfo   []NodeGPUInfo `json:"nodeGPUInfo,omitempty" yaml:"nodeGPUInfo,omitempty"`
	TotalGPUCount int64         `json:"totalGPUCount,omitempty" yaml:"totalGPUCount,omitempty"`
}
