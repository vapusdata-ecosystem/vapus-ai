package httpcls

import (
	"fmt"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type AIAgents struct {
	Params     *pb.AgentGetterRequest
	httpClient *VapusHttpClient
}

func NewAIAgentsClient(httpClient *VapusHttpClient, params *pb.AgentGetterRequest) *AIAgents {
	return &AIAgents{
		Params:     params,
		httpClient: httpClient,
	}
}

func (x *AIAgents) Create(httpOpts *HttpRequestGeneric) (*pb.AgentResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/agents"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AgentResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIAgents) Update(httpOpts *HttpRequestGeneric) (*pb.AgentResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/aistudio/agents"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AgentResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIAgents) Getter(httpOpts *HttpRequestGeneric, params *pb.AgentGetterRequest) (*pb.AgentResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/aistudio/agents"
	urlParams := url.Values{}
	if params.GetVapusAgentId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetVapusAgentId())
	} else {
		urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
		httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AgentResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIAgents) Archive(httpOpts *HttpRequestGeneric, params *pb.AgentGetterRequest) (*pb.AgentResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/aistudio/agents"
	if params.GetVapusAgentId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetVapusAgentId())
	} else {
		return nil, fmt.Errorf("aiAgentId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AgentResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
