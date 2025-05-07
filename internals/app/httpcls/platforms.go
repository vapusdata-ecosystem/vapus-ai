package httpcls

import (
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type VapusdataService struct {
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewVapusdataServiceClient(httpClient *VapusHttpClient, httpOpts *HttpRequestGeneric) *VapusdataService {
	return &VapusdataService{
		httpClient: httpClient,
	}
}

func (x *VapusdataService) AccountManager(httpOpts *HttpRequestGeneric) (*pb.AccountResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/platform"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AccountResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *VapusdataService) AccountGetter(httpOpts *HttpRequestGeneric) (*pb.AccountResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/platform"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AccountResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *VapusdataService) PlatformServicesInfo(httpOpts *HttpRequestGeneric) (*pb.VapusdataServicesResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/services"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.VapusdataServicesResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *VapusdataService) GetSampleResourceConfiguration(httpOpts *HttpRequestGeneric) (*pb.SampleResourceConfiguration, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/sample-resources"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.SampleResourceConfiguration{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
