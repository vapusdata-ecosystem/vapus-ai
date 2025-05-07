package plclient

import (
	"context"

	zerolog "github.com/rs/zerolog"
)

var MasterGlobals *GlobalsPersists
var MasterCommonFlags *CtlCommonFlags

type GlobalsPersists struct {
	CurrentContext    string
	Logger            zerolog.Logger
	CfgFile           string
	CfgDir            string
	DebugLogFlag      bool
	AgentsActions     map[string][]interface{}
	AgentsUtilities   map[string][]interface{}
	AgentsReflexes    map[string][]interface{}
	AgentInterfaceMap map[string]string

	*VapusCtlClient
	CurrentIdToken, CurrentAccessToken string
	Ctx                                context.Context
	CommonFlags                        *CtlCommonFlags
}

type CtlCommonFlags struct {
	La, DebugLogFlag                      bool
	CfgFile                               string
	Action, File, SearchQuery, OutputFile string
}
