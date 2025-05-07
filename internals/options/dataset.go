package options

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"google.golang.org/protobuf/types/known/structpb"
)

type DataSetSummary struct {
	Dataproducts   []string
	Query          string
	ResultLength   int64
	DataFields     []string
	ResultMap      []map[string]any
	TimeStamp      int64
	Description    string
	DataSources    []string
	ColumnsArray   [][]string
	IsColumnerData bool
}

type InternalDataQueryOpts struct {
	EndStream              bool
	SendStreamData         bool `default:"false"`
	MessageId              string
	StoreInFile            bool `default:"false"`
	GenerateQueryOnly      bool `default:"false"`
	GenerateFederatedQuery bool `default:"false"`
	GeneratedQuery         string
	ResponseMetadata       []*mpb.Mapper
	Response               []*structpb.Struct
}
