package client

const (
	GPUSpecType          = "gpuSpec"
	GPUSpecFieldNodeGPUs = "nodeGPUs"
)

type GPUSpec struct {
	NodeGPUs map[string]int64 `json:"nodeGPUs,omitempty" yaml:"nodeGPUs,omitempty"`
}
