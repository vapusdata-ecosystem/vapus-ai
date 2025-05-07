package models

import (
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type NabrunnerLog struct {
	VapusBase        `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	TaskId           string                `bun:"task_id" json:"taskId,omitempty" yaml:"taskId,omitempty"`
	SchedulerEntryId string                `bun:"scheduler_entry_id" json:"schedulerEntryId,omitempty" yaml:"schedulerEntryId,omitempty"`
	Schedule         *NabrunnerLogSchedule `bun:"schedule,type:jsonb" json:"schedule,omitempty" yaml:"schedule,omitempty"`
	IsRecurring      bool                  `bun:"is_recurring" json:"isRecurring,omitempty" yaml:"isRecurring,omitempty"`
	Resource         string                `bun:"resource" json:"resource,omitempty" yaml:"resource,omitempty"`
	ResourceId       string                `bun:"resource_id" json:"resourceId,omitempty" yaml:"resourceId,omitempty"`
	Canceled         bool                  `bun:"canceled" json:"canceled,omitempty" yaml:"canceled,omitempty"`
	Queue            string                `bun:"queue" json:"queue,omitempty" yaml:"queue,omitempty"`
	TaskType         string                `bun:"task_type" json:"taskType,omitempty" yaml:"taskType,omitempty"`
	Params           string                `bun:"params" json:"params,omitempty" yaml:"params,omitempty"`
	MaxRetries       int                   `bun:"max_retries" json:"maxRetries,omitempty" yaml:"maxRetries,omitempty"`
}

type NabrunnerLogSchedule struct {
	CronTab  string `json:"cronTab,omitempty" yaml:"cronTab,omitempty"`
	RunAt    int64  `json:"runAt,omitempty" yaml:"runAt,omitempty"`
	RunAfter int64  `json:"runAfter,omitempty" yaml:"runAfter,omitempty"`
	Limit    int64  `json:"limit,omitempty" yaml:"limit,omitempty"`
}

func (dn *NabrunnerLog) PreSaveCreate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreSaveVapusBase(authzClaim)
}

func (dn *NabrunnerLog) PreSaveUpdate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = authzClaim[encryption.ClaimUserIdKey]
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *NabrunnerLog) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = authzClaim[encryption.ClaimUserIdKey]
	dn.UpdatedAt = dmutils.GetEpochTime()
}
