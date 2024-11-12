package client

const (
	FaaStatusType      = "faaStatus"
	FaaStatusFieldMsg1 = "msg1"
)

type FaaStatus struct {
	Msg1 string `json:"msg1,omitempty" yaml:"msg1,omitempty"`
}
