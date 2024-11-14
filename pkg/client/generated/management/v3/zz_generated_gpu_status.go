package client

const (
	GPUStatusType         = "gpuStatus"
	GPUStatusFieldMessage = "message"
)

type GPUStatus struct {
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}
