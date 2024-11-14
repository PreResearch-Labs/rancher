package client

const (
	GPUSpecType         = "gpuSpec"
	GPUSpecFieldEnabled = "enabled"
)

type GPUSpec struct {
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}
