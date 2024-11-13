package client

const (
	GpuActionInputType          = "gpuActionInput"
	GpuActionInputFieldGPUCount = "gpuCount"
	GpuActionInputFieldNodeName = "nodeName"
)

type GpuActionInput struct {
	GPUCount int64  `json:"gpuCount,omitempty" yaml:"gpuCount,omitempty"`
	NodeName string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
}
