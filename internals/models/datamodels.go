package models

type OCILoginCreds struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ComplianceFieldOps struct {
	NullValues    int64 `json:"null_values"`
	FirstRecordAt int64 `json:"first_record_at"`
	LastRecordAt  int64 `json:"last_record_at"`
	TotalRecords  int64 `json:"total_records"`
}

type VDCLocalGovernanceModel struct {
	AccessScope         string `json:"access"`
	GovLabel            string `json:"govLabel"`
	DataProductId       string `json:"dataProductId"`
	Organization        string `json:"Organization"`
	DataProduct         string `json:"dataProduct"`
	UserId              string `json:"userId"`
	AccountId           string `json:"accountId"`
	AccessSelectorLabel string `json:"accessSelectorLabel"`
}

// Not db model
type DataWorkerSecretOpts struct {
	Extraction      *DataSourceSecretModel `validate:"required" json:"extraction" yaml:"extraction" toml:"extraction"`
	Loading         *DataSourceSecretModel `validate:"required" json:"loading" yaml:"loading" toml:"loading"`
	DataWorkerStore *DataSourceSecretModel `json:"data_worker_store" yaml:"data_worker_store" toml:"data_worker_store"`
	Destination     *DataSourceSecretModel `json:"destination" yaml:"destination" toml:"destination"`
}

// Not db model
type DataWorkerOpts struct {
	RequestData    []byte                `json:"request_data" yaml:"request_data" toml:"request_data"`
	DataWorkerType string                `json:"data_worker_type" yaml:"data_worker_type" toml:"data_worker_type"`
	InputFormat    string                `json:"input_format" yaml:"input_format" toml:"input_format"`
	DataWorkerOpts *DataWorkerSecretOpts `json:"agent_opts" yaml:"agent_opts" toml:"agent_opts"`
}

// Not db model
type DataSets struct {
	// Name of the collection, in case of SQL it will be table name and in case of NoSQL it will be collection name. I
	// if there are more than one table or collection then add it with comma separated values.
	DataTable []string         `json:"data_table" yaml:"data_table" toml:"data_table"`
	DataSet   []map[string]any `json:"data_set" yaml:"data_set" toml:"data_set"`
}

// Not db model
type WorkerDataset struct {
	DataSets map[string]*DataSets `json:"data_sets" yaml:"data_sets" toml:"data_sets"`
}

// NON-DB
type DataSourceSecrets struct {
	*GenericCredentialModel `json:"genericCredentialObj,omitempty" yaml:"genericCredentialObj,omitempty" toml:"genericCredentialObj,omitempty"`
	DB                      string                 `json:"db,omitempty" yaml:"db,omitempty" toml:"db,omitempty"`
	URL                     string                 `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	Port                    int64                  `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	Version                 string                 `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`
	Params                  map[string]interface{} `json:"params" yaml:"params" toml:"params"`
}

// NON-DB
type DataSourceCredsParams struct {
	DataSourceCreds       *DataSourceSecrets `validate:"required" json:"dataSourceCreds" yaml:"dataSourceCreds" toml:"dataSourceCreds"`
	DataSourceEngine      string             `validate:"required" json:"dataSourceEngine" yaml:"dataSourceEngine" toml:"dataSourceEngine"`
	DataSourceType        string             `validate:"required" json:"dataSourceType" yaml:"dataSourceType" toml:"dataSourceType"`
	DataSourceService     string             `json:"dataSourceService" yaml:"dataSourceService" toml:"dataSourceService"`
	Params                map[string]any     `json:"params" yaml:"params" toml:"params"`
	DataSourceSvcProvider string             `json:"dataSourceSvcProvider" yaml:"dataSourceSvcProvider" toml:"dataSourceSvcProvider"`
	DataStorelist         []string           `json:"dataStorelist" yaml:"dataStorelist" toml:"dataStorelist"`
	DataStorePrefixList   []string           `json:"dataStorePrefixList" yaml:"dataStorePrefixList" toml:"dataStorePrefixList"`
	Version               string             `json:"version" yaml:"version" toml:"version"`
	DataSourceId          string             `json:"dataSourceId" yaml:"dataSourceId" toml:"dataSourceId"`
}

func (x *DataSourceCredsParams) GetSourceCredentials() (bool, *GenericCredentialModel) {
	if x == nil {
		return false, nil
	}
	if x.DataSourceCreds == nil && x.DataSourceCreds.GenericCredentialModel != nil {
		return false, x.DataSourceCreds.GenericCredentialModel
	}
	return false, nil
}

// Non DB model
type DataSourceSecretModel struct {
	*DataSourceCredsParams
	DataSourceId string `json:"dataSourceId" yaml:"dataSourceId" toml:"dataSourceId"`
}

type CountModel struct {
	Count int64 `json:"count"`
}

type DataStoreSecription struct {
	DataStore        string           `json:"dataStore"`
	Tables           []map[string]any `json:"tables"`
	ComplianceFields []map[string]any `json:"complianceFields"`
}
