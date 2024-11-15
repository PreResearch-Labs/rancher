package client

const (
	NodeGPUInfoType              = "nodeGPUInfo"
	NodeGPUInfoFieldNodeHostName = "nodeHostName"
	NodeGPUInfoFieldNodeId       = "nodeId"
	NodeGPUInfoFieldNodeName     = "nodeName"
	NodeGPUInfoFieldTotalGPU     = "totalGPU"
	NodeGPUInfoFieldUnusedGPU    = "unusedGPU"
	NodeGPUInfoFieldUsedGPU      = "usedGPU"
)

type NodeGPUInfo struct {
	NodeHostName string `json:"nodeHostName,omitempty" yaml:"nodeHostName,omitempty"`
	NodeId       string `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	NodeName     string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	TotalGPU     int64  `json:"totalGPU,omitempty" yaml:"totalGPU,omitempty"`
	UnusedGPU    int64  `json:"unusedGPU,omitempty" yaml:"unusedGPU,omitempty"`
	UsedGPU      int64  `json:"usedGPU,omitempty" yaml:"usedGPU,omitempty"`
}
