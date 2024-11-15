package client

const (
	NodeGPUInfoType           = "nodeGPUInfo"
	NodeGPUInfoFieldNodeName  = "nodeName"
	NodeGPUInfoFieldTotalGPU  = "totalGPU"
	NodeGPUInfoFieldUnusedGPU = "unusedGPU"
	NodeGPUInfoFieldUsedGPU   = "usedGPU"
)

type NodeGPUInfo struct {
	NodeName  string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	TotalGPU  int64  `json:"totalGPU,omitempty" yaml:"totalGPU,omitempty"`
	UnusedGPU int64  `json:"unusedGPU,omitempty" yaml:"unusedGPU,omitempty"`
	UsedGPU   int64  `json:"usedGPU,omitempty" yaml:"usedGPU,omitempty"`
}
