package client

const (
	GpuCountActionInputType               = "gpuCountActionInput"
	GpuCountActionInputFieldClusterName   = "clusterName"
	GpuCountActionInputFieldDebug         = "debug"
	GpuCountActionInputFieldNodeHostName  = "nodeHostName"
	GpuCountActionInputFieldNodeId        = "nodeId"
	GpuCountActionInputFieldNodeName      = "nodeName"
	GpuCountActionInputFieldStatDimension = "statDimension"
)

type GpuCountActionInput struct {
	ClusterName   string `json:"clusterName,omitempty" yaml:"clusterName,omitempty"`
	Debug         bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
	NodeHostName  string `json:"nodeHostName,omitempty" yaml:"nodeHostName,omitempty"`
	NodeId        string `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	NodeName      string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	StatDimension string `json:"statDimension,omitempty" yaml:"statDimension,omitempty"`
}
