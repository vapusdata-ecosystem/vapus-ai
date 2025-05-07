package httpcls

import (
	"fmt"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type DataSource struct {
	Params     *pb.DataSourceGetterRequest
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewDataSourceClient(httpClient *VapusHttpClient, params *pb.DataSourceGetterRequest) *DataSource {
	return &DataSource{
		Params:     params,
		httpClient: httpClient,
	}
}

func (x *DataSource) Create(httpOpts *HttpRequestGeneric) (*pb.DataSourceResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/datasources"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.DataSourceResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *DataSource) Update(httpOpts *HttpRequestGeneric) (*pb.DataSourceResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/datasources"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.DataSourceResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *DataSource) Getter(httpOpts *HttpRequestGeneric, params *pb.DataSourceGetterRequest) (*pb.DataSourceResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/datasources"
	urlParams := url.Values{}
	if params.GetDataSourceId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetDataSourceId())
	} else {
		urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
		httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.DataSourceResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *DataSource) Archive(httpOpts *HttpRequestGeneric, params *pb.DataSourceGetterRequest) (*pb.DataSourceResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/datasources"
	if params.GetDataSourceId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetDataSourceId())
	} else {
		return nil, fmt.Errorf("ResourceId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.DataSourceResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *DataSource) Sync(httpOpts *HttpRequestGeneric, params *pb.DataSourceGetterRequest) (*pb.DataSourceResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/datasources"
	if params.GetDataSourceId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetDataSourceId()) + "/sync"
	} else {
		return nil, fmt.Errorf("ResourceId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.DataSourceResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
