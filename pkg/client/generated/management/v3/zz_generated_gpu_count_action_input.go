package client

const (
	GpuCountActionInputType               = "gpuCountActionInput"
	GpuCountActionInputFieldDebug         = "debug"
	GpuCountActionInputFieldNodeName      = "nodeName"
	GpuCountActionInputFieldStatDimension = "statDimension"
)

type GpuCountActionInput struct {
	Debug         bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
	NodeName      string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	StatDimension string `json:"statDimension,omitempty" yaml:"statDimension,omitempty"`
}
