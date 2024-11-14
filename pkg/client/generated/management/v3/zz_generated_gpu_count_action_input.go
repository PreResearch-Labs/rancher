package client

const (
	GpuCountActionInputType       = "gpuCountActionInput"
	GpuCountActionInputFieldCount = "count"
	GpuCountActionInputFieldDebug = "debug"
)

type GpuCountActionInput struct {
	Count bool `json:"count,omitempty" yaml:"count,omitempty"`
	Debug bool `json:"debug,omitempty" yaml:"debug,omitempty"`
}
