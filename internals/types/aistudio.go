package types

type GuardrailsProvider string

const (
	Bedrock GuardrailsProvider = "bedrock"
	Mistral GuardrailsProvider = "mistral"
	Nemo    GuardrailsProvider = "nemo"
	Pangea  GuardrailsProvider = "pangea"
	Vapus   GuardrailsProvider = "vapus"
)

func (x GuardrailsProvider) String() string {
	return string(x)
}
