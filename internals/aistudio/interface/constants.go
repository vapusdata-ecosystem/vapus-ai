package aiface

type AIInterfaceType string

const (
	Chat  AIInterfaceType = "chat"
	Audio AIInterfaceType = "audio"
	Video AIInterfaceType = "video"
)

func (a AIInterfaceType) String() string {
	return string(a)
}
