package client

const (
	GPUStatusType           = "gpuStatus"
	GPUStatusFieldTotalGPUs = "totalGPUs"
)

type GPUStatus struct {
	TotalGPUs int64 `json:"totalGPUs,omitempty" yaml:"totalGPUs,omitempty"`
}
