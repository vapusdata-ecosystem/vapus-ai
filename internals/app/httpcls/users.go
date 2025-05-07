package httpcls

import (
	"log"
	"net/url"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserManagement struct {
	httpOpts    *HttpRequestGeneric
	httpClient  *VapusHttpClient
	params      *pb.UserGetterRequest
	authzParams *pb.AuthzGetterRequest
}

func NewUserManagementClient(httpClient *VapusHttpClient, params *pb.UserGetterRequest, authzParams *pb.AuthzGetterRequest) *UserManagement {
	return &UserManagement{
		httpClient: httpClient,
	}
}

func (x *UserManagement) AccessTokenInterface(httpOpts *HttpRequestGeneric) (*pb.AccessTokenResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/auth"
	httpOpts.Method = types.POST
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AccessTokenResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *UserManagement) UserManager(httpOpts *HttpRequestGeneric) (*pb.UserResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/users"
	httpOpts.Method = types.POST
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.UserResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *UserManagement) UserGetter(httpOpts *HttpRequestGeneric, params *pb.UserGetterRequest) (*pb.UserResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/users"
	httpOpts.Method = types.GET
	urlParams := url.Values{}
	if params.GetUserId() != "" {
		httpOpts.Uri = appendIdToUrl(httpOpts.Uri, params.GetUserId())
	}
	urlParams = ConvertSearchParamToUrlValues(params.GetSearchParam(), urlParams)
	httpOpts.Uri = httpOpts.Uri + "/?" + urlParams.Encode()
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.UserResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *UserManagement) AuthzManager(httpOpts *HttpRequestGeneric) (*pb.AuthzResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/authz"
	httpOpts.Method = types.POST
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AuthzResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *UserManagement) AuthzGetter(httpOpts *HttpRequestGeneric, params *pb.AuthzGetterRequest) (*pb.AuthzResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/authz"
	httpOpts.Method = types.POST
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AuthzResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (x *UserManagement) LoginHandler(httpOpts *HttpRequestGeneric) (*pb.LoginHandlerResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/login"
	httpOpts.Method = types.GET
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	log.Println("LoginHandler Response: ", string(resp))
	response := &pb.LoginHandlerResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	log.Println("LoginHandler Response: ", response)
	return response, nil
}

func (x *UserManagement) LoginCallback(httpOpts *HttpRequestGeneric) (*pb.AccessTokenResponse, error) {
	httpOpts.Uri = "/api/v1alpha1/login/callback"
	httpOpts.Method = types.POST
	resp, err := x.httpClient.VapusHttpClient(httpOpts)
	if err != nil {
		return nil, err
	}
	response := &pb.AccessTokenResponse{}
	err = protojson.Unmarshal(resp, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
