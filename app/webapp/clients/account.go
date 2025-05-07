package clients

import (
	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"

	"github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
)

func (s *GrpcClient) GetAccountInfo(eCtx echo.Context) (*pb.AccountResponse, error) {
	return pkgs.VapusSvcInternalClientManager.PlConn.AccountGetter(s.SetAuthCtx(eCtx), &mpb.EmptyRequest{})
}

func (x *GrpcClient) GetUserInfo(eCtx echo.Context, userId string) (*mpb.User, map[string]string, error) {
	req := &pb.UserGetterRequest{
		Action: pb.UserGetterActions_GET_USER,
	}
	if userId != "" {
		req.UserId = userId
	}
	result, err := x.UserConn.UserGetter(x.SetAuthCtx(eCtx), req)

	if err != nil || result.Output == nil || len(result.Output.Users) == 0 {
		x.logger.Err(err).Msg("error while getting user info")
		return nil, nil, err
	}
	return result.Output.Users[0], result.OrganizationMap, nil
}

func (x *GrpcClient) GetMyOrganizationUsers(eCtx echo.Context) []*mpb.User {
	result, err := x.UserConn.UserGetter(x.SetAuthCtx(eCtx), &pb.UserGetterRequest{
		Action: pb.UserGetterActions_LIST_USERS,
	})

	if err != nil || result.Output == nil || len(result.Output.Users) == 0 {
		x.logger.Err(err).Msg("error while getting user info")
		return []*mpb.User{}
	}
	return result.Output.Users
}

func (x *GrpcClient) GetPlatformUsers(eCtx echo.Context) []*mpb.User {
	result, err := x.UserConn.UserGetter(x.SetAuthCtx(eCtx), &pb.UserGetterRequest{
		Action: pb.UserGetterActions_LIST_PLATFORM_USERS,
	})

	if err != nil || result.Output == nil || len(result.Output.Users) == 0 {
		x.logger.Err(err).Msg("error while getting user info")
		return []*mpb.User{}
	}
	return result.Output.Users
}

func (x *GrpcClient) GetDashboard(eCtx echo.Context) *mpb.OrganizationDashboard {
	result, err := x.OrganizationConn.Dashboard(x.SetAuthCtx(eCtx), &mpb.EmptyRequest{})

	if err != nil || result.Output == nil {
		x.logger.Err(err).Msg("error while getting dashboard info")
		return &mpb.OrganizationDashboard{}
	}
	return result.Output
}

func (x *GrpcClient) ListPlugins(eCtx echo.Context) []*mpb.Plugin {
	result, err := x.PluginServiceClient.List(x.SetAuthCtx(eCtx), &pb.PluginGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin list")
		return []*mpb.Plugin{}
	}
	return result.Output
}

func (x *GrpcClient) GetPlugin(eCtx echo.Context, id string) *mpb.Plugin {
	result, err := x.PluginServiceClient.Get(x.SetAuthCtx(eCtx), &pb.PluginGetterRequest{
		PluginId: id,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting plugin info")
		return nil
	}
	return result.Output[0]
}

func (x *GrpcClient) ResourceGetter(eCtx echo.Context) []*pb.PluginTypeMap {
	result, err := x.PlConn.ResourceGetter(x.SetAuthCtx(eCtx), &pb.ResourceGetterRequest{})
	if err != nil || len(result.PluginTypeMap) == 0 {
		x.logger.Err(err).Msg("error while getting plugin info")
		return nil
	}
	return result.PluginTypeMap
}

func (x *GrpcClient) PlatformPublicInfo(eCtx echo.Context, id string) *pb.PlatformPublicInfoResponse {
	result, err := x.PlConn.PlatformPublicInfo(x.SetAuthCtx(eCtx), &mpb.EmptyRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting platform public info")
		return result
	}
	return result
}

func (x *GrpcClient) GetSecretServiceList(eCtx echo.Context) []*mpb.SecretStore {
	result, err := x.SecretServiceClient.List(x.SetAuthCtx(eCtx), &pb.SecretGetterRequest{})
	if err != nil {
		x.logger.Err(err).Msg("error while getting the secret service list")
		return []*mpb.SecretStore{}
	}
	return result.Output
}

func (x *GrpcClient) SecretServiceDetails(eCtx echo.Context, name string) *mpb.SecretStore {
	result, err := x.SecretServiceClient.Get(x.SetAuthCtx(eCtx), &pb.SecretGetterRequest{
		Name: name,
	})
	if err != nil || len(result.Output) == 0 {
		x.logger.Err(err).Msg("error while getting Secret Service Details")
		return nil
	}
	return result.Output[0]
}
