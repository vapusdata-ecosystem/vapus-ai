package models

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type FabricDataQuerySummary struct {
	Dataproducts []string
	Query        string
	ResultLength int64
	DataFields   []string
	ResultMap    []map[string]any
}

type FabricOwnerMetricsChatLog struct {
	VapusBase    `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	ChatId       string                                    `bun:"chat_id" json:"chatId,omitempty" yaml:"chatId,omitempty" toml:"chatId,omitempty"`
	ThreadEnded  bool                                      `bun:"thread_ended" json:"threadEnded,omitempty" yaml:"threadEnded,omitempty" toml:"threadEnded,omitempty"`
	Messages     map[string]*FabricOwnerMetricsChatMessage `bun:"message,type:jsonb" json:"message,omitempty" yaml:"message,omitempty" toml:"message,omitempty"`
	MessageQueue []string                                  `bun:"message_queue,array" json:"messageQueue,omitempty" yaml:"messageQueue,omitempty" toml:"messageQueue,omitempty"`
}

type FabricOwnerMetricsChatMessage struct {
	MessageId          string           `json:"messageId,omitempty" yaml:"messageId,omitempty" toml:"messageId,omitempty"`
	Input              string           `json:"input,omitempty" yaml:"input,omitempty" toml:"input,omitempty"`
	InputDataProducts  []string         `json:"inputDataProducts,omitempty" yaml:"inputDataProducts,omitempty" toml:"inputDataProducts,omitempty"`
	OutputDataProducts []string         `json:"outputDataProducts,omitempty" yaml:"outputDataProducts,omitempty" toml:"outputDataProducts,omitempty"`
	InputFiles         []*mpb.FileData  `json:"inputFiles,omitempty" yaml:"inputFiles,omitempty" toml:"inputFiles,omitempty"`
	OutputFiles        []*mpb.FileData  `json:"outputFiles,omitempty" yaml:"outputFiles,omitempty" toml:"outputFiles,omitempty"`
	OutputSummary      []map[string]any `json:"outputSummary,omitempty" yaml:"outputSummary,omitempty" toml:"outputSummary,omitempty"`
	Error              string           `json:"error,omitempty" yaml:"error,omitempty" toml:"error,omitempty"`
	Timestamp          int64            `json:"timestamp,omitempty" yaml:"timestamp,omitempty" toml:"timestamp,omitempty"`
}

func (dmn *FabricOwnerMetricsChatLog) PreSaveCreate(authzClaim map[string]string) {
	if dmn == nil {
		return
	}
	if dmn.CreatedBy == types.EMPTYSTR {
		dmn.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dmn.CreatedAt == 0 {
		dmn.CreatedAt = dmutils.GetEpochTime()
	}
	dmn.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
}

func (dn *FabricOwnerMetricsChatLog) EndThread() {
	if dn == nil {
		return
	}
	dn.DeletedAt = dmutils.GetEpochTime()
	dn.DeletedBy = dn.CreatedBy
	dn.Status = mpb.CommonStatus_EXPIRED.String()
	dn.ThreadEnded = true
}
