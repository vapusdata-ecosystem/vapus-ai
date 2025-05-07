package pkgs

import (
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type DataTablesOpts struct {
	Query        string
	StructScheme interface{}
	MapsScheme   map[string]any
	Name         string
}

type InsertDataRequest struct {
	DataTable      string           `json:"data_tables"`
	TableSchema    string           `json:"table_schema"`
	DataSets       []map[string]any `json:"data_sets"`
	DataSet        map[string]any   `json:"data_set"`
	BatchSize      int64            `json:"batch_size"`
	BatchWait      int64            `json:"batch_wait"`
	StructDataSets []interface{}    `json:"struct_data_sets"`
	StructDataSet  interface{}      `json:"struct_data_set"`
}

type InsertDataResponse struct {
	DataTable       string           `json:"data_tables"`
	RecordsInserted int64            `json:"records_inserted"`
	RecordsFailed   int64            `json:"records_failed"`
	FailedDataSet   []map[string]any `json:"failed_dataSet"`
}

type SQLBaseModel struct {
	TableType string `gorm:"column:TABLE_TYPE"`
	TableName string `gorm:"column:TABLE_NAME"`
}

type SQLConstraints struct {
	ConstraintName string `gorm:"column:CONSTRAINT_NAME"`
	TableName      string `gorm:"column:TABLE_NAME"`
	ConstraintType string `gorm:"column:CONSTRAINT_TYPE"`
	Enforced       string `gorm:"column:ENFORCED"`
	TargetTable    string `gorm:"column:TARGET_TABLE"`
}

type SQLTableInfo struct {
	TableName    string `gorm:"column:TABLE_NAME"`
	TableType    string `gorm:"column:TABLE_TYPE"`
	Engine       string `gorm:"column:ENGINE"`
	TableRows    uint64 `gorm:"column:TABLE_ROWS"`
	DataLength   uint64 `gorm:"column:DATA_LENGTH"`
	IndexLength  uint64 `gorm:"column:INDEX_LENGTH"`
	Version      string `gorm:"column:VERSION"`
	AvgRowLength uint64 `gorm:"column:AVG_ROW_LENGTH"`
	CreatedAt    string `gorm:"column:CREATE_TIME"`
	UpdatedAt    string `gorm:"column:UPDATE_TIME"`
}

type VDCLocalGovernanceModel struct {
	AccessScope         string `json:"access"`
	GovLabel            string `json:"govLabel"`
	DataProductId       string `json:"dataProductId"`
	ORGANIZATION        string `json:"organization"`
	DataProduct         string `json:"dataProduct"`
	UserId              string `json:"userId"`
	AccountId           string `json:"accountId"`
	AccessSelectorLabel string `json:"accessSelectorLabel"`
}

type Filter struct {
	Name               string `json:"name"`
	FieldName          string `json:"fieldName"`
	Comparator         string `json:"comparator"`
	Value              string `json:"value"`
	IsBoolField        bool   `json:"boolField"`
	NextFilterOperator string `json:"nextFilter"`
}

type QueryOpts struct {
	DistinctFields      string                   `json:"distinctColumns"`
	DataCollection      string                   `json:"dataCollections"`
	IncludeFields       []string                 `json:"fields"`
	ExcludeFields       []string                 `json:"excludeFields"`
	QueryString         string                   `json:"queryString,omitempty"`
	Filters             []*Filter                `json:"filters"`
	GroupBy             []string                 `json:"groupBy"`
	OrderBy             *OrderByParam            `json:"orderBy"`
	JoinDataCollections []string                 `json:"joinDataCollections"`
	RawQuery            string                   `json:"rawQuery"`
	Limit               int64                    `json:"limit"`
	BatchSize           int64                    `json:"batchSize"`
	CountRecords        bool                     `json:"countRecords"`
	QueryFilters        map[string][]interface{} `json:"queryFilters"`
}

type OrderType string

var ASC OrderType = "asc"
var DESC OrderType = "desc"

func (o OrderType) String() string {
	return string(o)
}

type OrderByParam struct {
	Field string    `json:"field"`
	Order OrderType `json:"order"`
}

var PostgresGinTemplate = "CREATE INDEX IF NOT EXISTS %s ON %s USING GIN(%s)"

type PostgresIndexOpts struct {
	TableName string
	FieldName string
	Indexname string
	IndexType string
	IndexAlgo string
}

type AIQueryGeneratorparams struct {
	DatabaseType  string
	Schema        string
	TextQuery     string
	AimodelParams *mpb.AIModelLocalParams
	Logger        zerolog.Logger
}
