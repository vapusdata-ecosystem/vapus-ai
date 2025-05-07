package httpcls

import (
	"fmt"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type AIGuardrails struct {
	Params     *pb.GuardrailsGetterRequest
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewAIGuardrailsClient(httpClient *VapusHttpClient, params *pb.GuardrailsGetterRequest) *AIGuardrails {
	return &AIGuardrails{
		Params:     params,
		httpClient: httpClient,
	}
}

func (x *AIGuardrails) Create(httpOpts *HttpRequestGeneric) (*pb.GuardrailsResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/guardrails"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.GuardrailsResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIGuardrails) Update(httpOpts *HttpRequestGeneric) (*pb.GuardrailsResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/aistudio/guardrails"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.GuardrailsResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIGuardrails) Getter(httpOpts *HttpRequestGeneric, params *pb.GuardrailsGetterRequest) (*pb.GuardrailsResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/aistudio/guardrails"
	urlParams := url.Values{}
	if params.GetGuardrailId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetGuardrailId())
	} else {
		urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
		httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.GuardrailsResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIGuardrails) Archive(httpOpts *HttpRequestGeneric, params *pb.GuardrailsGetterRequest) (*pb.GuardrailsResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/aistudio/guardrails"
	if params.GetGuardrailId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetGuardrailId())
	} else {
		return nil, fmt.Errorf("ResourceId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.GuardrailsResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
