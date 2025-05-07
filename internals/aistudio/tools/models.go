package tools

import (
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aiinteface "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/interface"
	appdrepo "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo"
	aidmstore "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo/aistudio"
	appcl "github.com/vapusdata-ecosystem/vapusdata/core/app/grpcclients"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
)

type ToolPropertiesType string

func (t ToolPropertiesType) String() string {
	return string(t)
}

const (
	FuncParamType                    = "object"
	StringProp    ToolPropertiesType = "string"
)

type Options struct {
	Input               string
	Sysmess             string
	Tools               []*mpb.ToolCall
	ActionAnalyzer      bool
	ParamAnalyzer       bool
	FetchData           bool
	RetryCount          int `default:"0"`
	CtxClaim            map[string]string
	AssistantMessages   []string
	MaxCompletionTokens int32
	Temperature         float32
	ToolChoiceParam     string
	ToolChoice          *mpb.ToolChoice
	// Stream
}

type ToolCaller struct {
	AIModel       string
	AIModelNode   string
	Client        *appcl.VapusSvcInternalClients
	Logger        zerolog.Logger
	TcDbStore     *apppkgs.VapusStore
	ModelConnPool *appdrepo.AIModelNodeConnectionPool
	GuardrailPool *appdrepo.GuardrailPool
	AIDmStore     *aidmstore.AIStudioDMStore
	Agent         *aiinteface.AIChatGateway
}
