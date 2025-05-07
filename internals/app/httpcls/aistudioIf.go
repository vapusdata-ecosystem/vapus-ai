package httpcls

import (
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type AIStudio struct {
	httpOpts   *HttpRequestGeneric
	httpClient *VapusHttpClient
}

func NewAIStudioClient(httpClient *VapusHttpClient, httpOpts *HttpRequestGeneric) *AIStudio {
	return &AIStudio{
		httpClient: httpClient,
	}
}

func (x *AIStudio) Chat(httpOpts *HttpRequestGeneric) (*pb.ChatResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/chat"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.ChatResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIStudio) ChatStream(httpOpts *HttpRequestGeneric) (*pb.ChatResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/chat-stream"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.ChatResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *AIStudio) GenerateEmbeddings(httpOpts *HttpRequestGeneric) (*pb.EmbeddingsResponse, error) {
	httpOpts.Method = types.POST
	httpOpts.Uri = "/api/v1alpha1/aistudio/embeddings"
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.EmbeddingsResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
