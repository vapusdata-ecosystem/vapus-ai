package httpcls

import (
	"fmt"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type AIModels struct {
	Params     *pb.AIModelNodeGetterRequest
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewAIModelsClient(httpClient *VapusHttpClient, params *pb.AIModelNodeGetterRequest) *AIModels {
	return &AIModels{
		Params:     params,
		httpClient: httpClient,
	}
}

func (x *AIModels) Create(httpOpts *HttpRequestGeneric) (*pb.AIModelNodeResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/models"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AIModelNodeResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIModels) Update(httpOpts *HttpRequestGeneric) (*pb.AIModelNodeResponse, error) {
	httpOpts.Method = types.PUT
	httpOpts.Uri = "/api/v1alpha1/aistudio/models"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AIModelNodeResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIModels) Getter(httpOpts *HttpRequestGeneric, params *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	httpOpts.Method = types.GET
	httpOpts.Uri = "/api/v1alpha1/aistudio/models"
	urlParams := url.Values{}
	if params.GetAiModelNodeId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetAiModelNodeId())
	} else {
		urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
		httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AIModelNodeResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIModels) Archive(httpOpts *HttpRequestGeneric, params *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/aistudio/models"
	if params.GetAiModelNodeId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetAiModelNodeId())
	} else {
		return nil, fmt.Errorf("ResourceId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AIModelNodeResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIModels) Sync(httpOpts *HttpRequestGeneric, params *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	httpOpts.Method = types.DELETE
	httpOpts.Uri = "/api/v1alpha1/aistudio/models"
	if params.GetAiModelNodeId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetAiModelNodeId()) + "/sync"
	} else {
		return nil, fmt.Errorf("aiModelNodeId is required")
	}
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AIModelNodeResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
