package models

import (
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

type BaseParams struct {
	Schedule         *models.VapusSchedule `json:"schedule"`
	Scheduled        bool                  `json:"scheduled"`
	EntryId          string                `json:"entryId"`
	ScheduledEntryId string                `json:"scheduledEntryId"`
	MaxRetries       int                   `json:"maxRetries"`
	InitiatedBy      string                `json:"initiatedBy"`
}

type CancellationServiceParams struct {
	TaskId           string `json:"task_id"`
	IsRecurring      bool   `json:"isRecurring"`
	ScheduledEntryId string `json:"scheduledEntryId"`
	Queue            string `json:"queue"`
	TaskType         string `json:"taskType"`
	CancelledBy      string `json:"cancelledBy"`
}
