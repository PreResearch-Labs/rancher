package client

const (
	EchoActionInputType      = "echoActionInput"
	EchoActionInputFieldAaaa = "aaaa"
	EchoActionInputFieldEcho = "echo"
)

type EchoActionInput struct {
	Aaaa bool `json:"aaaa,omitempty" yaml:"aaaa,omitempty"`
	Echo bool `json:"echo,omitempty" yaml:"echo,omitempty"`
}
