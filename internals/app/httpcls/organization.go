package httpcls

import (
	"errors"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type Organization struct {
	Params     *pb.OrganizationGetterRequest
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewOrganizationClient(httpClient *VapusHttpClient, params *pb.OrganizationGetterRequest) *Organization {
	return &Organization{
		Params:     params,
		httpClient: httpClient,
	}
}

func (x *Organization) Create(httpOpts *HttpRequestGeneric) (*pb.OrganizationResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/Organizations"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.OrganizationResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *Organization) Update(httpOpts *HttpRequestGeneric) (*pb.OrganizationResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/Organizations"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.OrganizationResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *Organization) Getter(httpOpts *HttpRequestGeneric, params *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/Organizations"
	urlParams := url.Values{}
	if params.GetOrganizationId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetOrganizationId())
	} else {
		urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
		httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.OrganizationResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *Organization) Archive(httpOpts *HttpRequestGeneric, params *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/Organizations"
	if params.GetOrganizationId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetOrganizationId())
	} else {
		return nil, errors.New("Organization ID is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.OrganizationResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *Organization) UpgradeOS(httpOpts *HttpRequestGeneric, params *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/Organizations"
	if params.GetOrganizationId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetOrganizationId()) + "/upgrade-os"
	} else {
		return nil, errors.New("Organization ID is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.OrganizationResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *Organization) AddUsers(httpOpts *HttpRequestGeneric, params *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/Organizations"
	if params.GetOrganizationId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetOrganizationId()) + "/users"
	} else {
		return nil, errors.New("Organization ID is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.OrganizationResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
