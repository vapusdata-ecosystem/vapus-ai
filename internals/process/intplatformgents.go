package processes

import (
	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

// Non DB model
type AgentLiv struct {
	Operation          string           `json:"operation" yaml:"operation" toml:"operation"`
	AgentId            string           `json:"agentId" yaml:"agentId" toml:"agentId"`
	AgentType          types.AgentType  `json:"agentType" yaml:"agentType" toml:"agentType"`
	AgentStatus        mpb.CommonStatus `json:"agentStatus" yaml:"agentStatus" toml:"agentStatus"`
	InitializedAt      int64            `json:"initializedAt" yaml:"initializedAt" toml:"initializedAt"`
	OperationStartTime int64            `json:"operationStartTime" yaml:"operationStartTime" toml:"operationStartTime"`
	OperationEndTime   int64            `json:"operationEndTime" yaml:"operationEndTime" toml:"operationEndTime"`
	OperationStatus    mpb.CommonStatus `json:"operationStatus" yaml:"operationStatus" toml:"operationStatus"`
	OperationErr       string           `json:"operationErr" yaml:"operationErr" toml:"operationErr"`
	InvokedBy          string           `json:"invokedBy" yaml:"invokedBy" toml:"invokedBy"`
	FuncsCalled        []string         `json:"funcsCalled" yaml:"funcsCalled" toml:"funcsCalled"`
	Parents            []string         `json:"parentAgents" yaml:"parentAgents" toml:"parentAgents"`
	Logs               []string         `json:"logs" yaml:"logs" toml:"logs"`
	ErrorLogs          []error          `json:"errorLogs" yaml:"errorLogs" toml:"errorLogs"`
}

// Non DB model
func HBDAgent(ops string) *AgentLiv {
	return &AgentLiv{
		AgentId:       guuid.New().String(),
		Logs:          make([]string, 0),
		ErrorLogs:     make([]error, 0),
		Operation:     ops,
		AgentStatus:   mpb.CommonStatus_ACTIVE,
		InitializedAt: dmutils.GetEpochTime(),
	}
}
