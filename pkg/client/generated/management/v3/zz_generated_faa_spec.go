package client

const (
	FaaSpecType      = "faaSpec"
	FaaSpecFieldBar1 = "bar1"
)

type FaaSpec struct {
	Bar1 *bool `json:"bar1,omitempty" yaml:"bar1,omitempty"`
}
