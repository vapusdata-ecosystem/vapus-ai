package models

import (
	"encoding/json"
	"strings"
	"time"

	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"gopkg.in/yaml.v3"
)

type DataWorkerRun struct {
	WorkerRunId             string     `bun:"worker_run_id" json:"workerRunId"`
	SyncStartAt             int64      `bun:"sync_start_at" json:"syncStartAt"`
	SyncStatus              string     `bun:"sync_status" json:"syncStatus"`
	SyncEndAt               int64      `bun:"sync_end_at" json:"syncEndAt"`
	Status                  string     `bun:"status" json:"status"`
	SyncId                  string     `bun:"sync_id" json:"syncId"`
	Logs                    []*Logs    `bun:"logs,type:jsonb" json:"logs"`
	SpecDigest              *DigestVal `bun:"spec_digest,type:jsonb" json:"specDigest"`
	DataWorkerType          string     `bun:"data_worker_type" json:"dataWorkerType"`
	DataWorkerId            string     `bun:"data_worker_id" json:"dataWorkerId"`
	DataTable               string     `bun:"data_table" json:"dataTable"`
	TotalRecordsExtracted   int64      `bun:"total_records_extracted" json:"totalRecordsExtracted"`
	RecordFailedToLoaded    int64      `bun:"record_failed_to_loaded" json:"recordFailedToLoad"`
	TotalRecordsLoaded      int64      `bun:"total_records_loaded" json:"totalRecordsLoaded"`
	TotalRecordsTransformed int64      `bun:"total_records_transformed" json:"totalRecordsTransformed"`
	TotalFieldsTransformed  int64      `bun:"total_fields_transformed" json:"totalFieldsTransformed"`
}

func (w *DataWorkerRun) SetWorkerRunId(dataWorkerId string) {
	w.WorkerRunId = dmutils.SetIds(types.IdSeparator, "dataWorkerRun", dataWorkerId, guuid.New().String())
}
func (wl *DataWorkerRun) SaveWorkerRunLog(logType string, message ...string) {
	if wl == nil {
		return
	}
	message = append([]string{wl.WorkerRunId}, message...)
	if len(wl.Logs) > 0 {
		wl.Logs = append(wl.Logs, &Logs{
			Time:    time.Now().Unix(),
			Message: strings.Join(message, "||"),
			LogType: logType,
		})
	} else {
		wl.Logs = []*Logs{{
			Time:    time.Now().Unix(),
			Message: strings.Join(message, "||"),
			LogType: logType,
		}}
	}
}

func (x *DataWorkerRun) GetFormatedWorkerRun(fileFormat string) string {
	if x != nil {
		switch strings.ToLower(fileFormat) {
		case strings.ToLower(mpb.ContentFormats_YAML.String()):
			yamlData, err := yaml.Marshal(x)
			if err != nil {
				return ""
			}
			return string(yamlData)
		case strings.ToLower(mpb.ContentFormats_JSON.String()):
			jsonData, err := json.Marshal(x)
			if err != nil {
				return ""
			}
			return string(jsonData)
		default:
			return ""
		}
	} else {
		return ""
	}
}

type Logs struct {
	Time    int64  `json:"time,omitempty" yaml:"time"`
	HTime   string `json:"hTime,omitempty" yaml:"hTime"`
	LogType string `json:"logType,omitempty" yaml:"logType"`
	Message string `json:"message,omitempty" yaml:"message"`
}
