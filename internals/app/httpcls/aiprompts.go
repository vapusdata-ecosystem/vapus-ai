package httpcls

import (
	"fmt"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type AIPrompts struct {
	Params     *pb.PromptGetterRequest
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewAIPromptsClient(httpClient *VapusHttpClient, params *pb.PromptGetterRequest) *AIPrompts {
	return &AIPrompts{
		Params:     params,
		httpClient: httpClient,
	}
}

func (x *AIPrompts) Create(httpOpts *HttpRequestGeneric) (*pb.PromptResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/prompts"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.PromptResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIPrompts) Update(httpOpts *HttpRequestGeneric) (*pb.PromptResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/aistudio/prompts"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.PromptResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIPrompts) Getter(httpOpts *HttpRequestGeneric, params *pb.PromptGetterRequest) (*pb.PromptResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/aistudio/prompts"
	urlParams := url.Values{}
	if params.GetPromptId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetPromptId())
	} else {
		urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
		httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.PromptResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIPrompts) Archive(httpOpts *HttpRequestGeneric, params *pb.PromptGetterRequest) (*pb.PromptResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/aistudio/prompts"
	if params.GetPromptId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetPromptId())
	} else {
		return nil, fmt.Errorf("ResourceId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.PromptResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
