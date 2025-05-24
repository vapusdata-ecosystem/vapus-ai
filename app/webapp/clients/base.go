package clients

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	appcl "github.com/vapusdata-ecosystem/vapusai/core/app/grpcclients"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/webapp/utils"
)

type GrpcClient struct {
	*appcl.VapusSvcInternalClients
	logger zerolog.Logger
}

var GrpcClientManager *GrpcClient

func NewGrpcClient() *GrpcClient {
	logger := pkgs.GetSubDMLogger("webapp", "grpcClients")
	err := appcl.SvcUpTimeCheck(context.Background(), pkgs.NetworkConfigManager, "", logger, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while checking service uptime.")
	} else {
		logger.Info().Msg("service is up and running.")
	}
	cl, err := appcl.SetupVapusSvcInternalClients(context.Background(), pkgs.NetworkConfigManager, "", logger)
	if err != nil {
		logger.Err(err).Msg("error while initializing vapus svc internal clients.")
	}
	return &GrpcClient{
		VapusSvcInternalClients: cl,
		logger:                  logger,
	}
}

func InitGrpcClient() {
	if GrpcClientManager == nil {
		GrpcClientManager = NewGrpcClient()
	}
}

func (s *GrpcClient) Close() {
	s.VapusSvcInternalClients.Close()
}

func (s *GrpcClient) SetAuthCtx(eCtx echo.Context) context.Context {
	token, err := utils.GetCookie(eCtx, types.ACCESS_TOKEN)
	if err != nil || token == "" {
		return eCtx.Request().Context()
	}
	return utils.GetBearerCtx(context.Background(), token)
}
