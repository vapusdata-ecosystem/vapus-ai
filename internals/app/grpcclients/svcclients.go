package appcl

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rs/zerolog"
	atpb "github.com/vapusdata-ecosystem/apis/protos/vapus-aiutilities/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	svcconfig "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type VapusSvcInternalClients struct {
	Host                string
	PlConn              pb.VapusdataServiceClient
	DatasourceConn      pb.DatasourceServiceClient
	UserConn            pb.UserManagementServiceClient
	OrganizationConn    pb.OrganizationServiceClient
	UtilityConn         pb.UtilityServiceClient
	AIStudioConn        pb.AIStudioClient
	AgentServiceClient  pb.AgentServiceClient
	AgentStudioClient   pb.AgentStudioClient
	AIPromptClient      pb.AIPromptsClient
	AIModelClient       pb.AIModelsClient
	AIGurdrailsClient   pb.AIGuardrailsClient
	platformGrpcClient  *pbtools.GrpcClient
	aiStudioGrpcClient  *pbtools.GrpcClient
	PluginServiceClient pb.PluginServiceClient
	NetworkConfig       *svcconfig.NetworkConfig
	// PlDns                 string
	// NabhikServerDns       string
	AIStudioDns           string
	DataServerClient      *pbtools.GrpcClient
	SecretServiceClient   pb.SecretServiceClient
	AIUtilityServerClient atpb.AIUtilityClient
	AiUtilityGrpcClient   *pbtools.GrpcClient
	AiUtilityServerDns    string
}

func SvcUpTimeCheck(ctx context.Context, networkConfig *svcconfig.NetworkConfig, self string, logger zerolog.Logger, counter int64) error {
	if counter > 1 {
		time.Sleep(15 * time.Second)
	}
	counter++
	logger.Info().Msg("Checking if all services are up........")
	aiStudioSvcDns := fmt.Sprintf("%s:%d", networkConfig.AIStudioSvc.ServiceName, networkConfig.AIStudioSvc.ServicePort)
	aiUtilityDns := fmt.Sprintf("%s:%d", networkConfig.NabhikServer.ServiceName, networkConfig.AIUtility.ServicePort)

	err := dmutils.Telnet("tcp", aiStudioSvcDns)
	if err != nil {
		logger.Error().Err(err).Msg("AI Studio service is not up yet")
		if counter > 6 {
			return err
		} else {
			return SvcUpTimeCheck(ctx, networkConfig, self, logger, counter)
		}
	}

	err = dmutils.Telnet("tcp", aiUtilityDns)
	if err != nil {
		logger.Error().Err(err).Msg("Ai Utility Server service is not up yet")
		if counter > 6 {
			return err
		} else {
			return SvcUpTimeCheck(ctx, networkConfig, self, logger, counter)
		}
	}
	return nil
}

func SetupVapusSvcInternalClients(ctx context.Context, networkConfig *svcconfig.NetworkConfig, self string, logger zerolog.Logger) (*VapusSvcInternalClients, error) {
	var err error
	client := &VapusSvcInternalClients{
		AIStudioDns: fmt.Sprintf("%s:%d", networkConfig.AIStudioSvc.ServiceName, networkConfig.AIStudioSvc.ServicePort),
		//  here I need to do something
		// PlDns: fmt.Sprintf("%s:%d", networkConfig.PlatformSvc.ServiceName, networkConfig.PlatformSvc.ServicePort),
		// NabhikServerDns:    fmt.Sprintf("%s:%d", networkConfig.NabhikServer.ServiceName, networkConfig.NabhikServer.ServicePort),
		AiUtilityServerDns: fmt.Sprintf("%s:%d", networkConfig.AIUtility.ServiceName, networkConfig.AIUtility.ServicePort),
	}
	logger.Info().Msg("Setting up VapusSvcInternalClients........")
	// logger.Info().Msgf("PlatformSvcDns: %s", client.PlDns)
	logger.Info().Msgf("AIStudioSvcDns: %s", client.AIStudioDns)
	// logger.Info().Msgf("DataproductServerDns: %s", client.NabhikServerDns)
	logger.Info().Msgf("AiutilityServerDns: %s", client.AiUtilityServerDns)
	// err = dmutils.Telnet("tcp", client.PlDns)
	// if err != nil && self != "" && self != networkConfig.PlatformSvc.ServiceName {
	// 	logger.Error().Err(err).Msg("Platform service is not up yet")
	// 	client.platformGrpcClient = nil
	// }
	logger.Info().Msg("Setting up VapusSvcInternalClients........")
	log.Println("client.platformGrpcClient: ", client.platformGrpcClient)
	// if client.platformGrpcClient == nil {
	// 	client.platformGrpcClient = pbtools.NewGrpcClient(logger,
	// 		pbtools.ClientWithInsecure(true),
	// 		pbtools.ClientWithServiceAddress(client.PlDns))
	// 	// platformGrpcClient => removed
	// }

	err = dmutils.Telnet("tcp", client.AIStudioDns)
	if err != nil && self != "" && self != networkConfig.AIStudioSvc.ServiceName {
		logger.Error().Err(err).Msg("AI Studio service is not up yet")
		client.aiStudioGrpcClient = nil
	}
	log.Println("client.aiStudioGrpcClient: ", client.aiStudioGrpcClient)
	if client.aiStudioGrpcClient == nil {
		client.aiStudioGrpcClient = pbtools.NewGrpcClient(logger,
			pbtools.ClientWithInsecure(true),
			pbtools.ClientWithServiceAddress(client.AIStudioDns))
		client.AIStudioConn = pb.NewAIStudioClient(client.aiStudioGrpcClient.Connection)
		client.AIPromptClient = pb.NewAIPromptsClient(client.aiStudioGrpcClient.Connection)
		client.AIModelClient = pb.NewAIModelsClient(client.aiStudioGrpcClient.Connection)
		client.AIGurdrailsClient = pb.NewAIGuardrailsClient(client.aiStudioGrpcClient.Connection)
		client.AgentServiceClient = pb.NewAgentServiceClient(client.aiStudioGrpcClient.Connection)
		client.AgentStudioClient = pb.NewAgentStudioClient(client.aiStudioGrpcClient.Connection)
		client.OrganizationConn = pb.NewOrganizationServiceClient(client.aiStudioGrpcClient.Connection)
		client.PlConn = pb.NewVapusdataServiceClient(client.aiStudioGrpcClient.Connection)
		client.UserConn = pb.NewUserManagementServiceClient(client.aiStudioGrpcClient.Connection)
		client.PluginServiceClient = pb.NewPluginServiceClient(client.aiStudioGrpcClient.Connection)
		client.UtilityConn = pb.NewUtilityServiceClient(client.aiStudioGrpcClient.Connection)
		client.DatasourceConn = pb.NewDatasourceServiceClient(client.aiStudioGrpcClient.Connection)
		client.SecretServiceClient = pb.NewSecretServiceClient(client.aiStudioGrpcClient.Connection)
	}
	client.NetworkConfig = networkConfig

	// err = dmutils.Telnet("tcp", client.NabhikServerDns)
	// if err != nil && self != "" && self != networkConfig.NabhikServer.ServiceName {
	// 	logger.Error().Err(err).Msg("NabhikServer service is not up yet")
	// 	client.DataServerClient = nil
	// }
	log.Println("client.NewDataServerClient: ", client.DataServerClient)
	// if client.DataServerClient == nil {
	// 	client.DataServerClient = pbtools.NewGrpcClient(logger,
	// 		pbtools.ClientWithInsecure(true),
	// 		pbtools.ClientWithServiceAddress(client.NabhikServerDns))
	// }
	if client.AiUtilityGrpcClient == nil {
		client.AiUtilityGrpcClient = pbtools.NewGrpcClient(logger,
			pbtools.ClientWithInsecure(true),
			pbtools.ClientWithServiceAddress(client.AiUtilityServerDns))
		client.AIUtilityServerClient = atpb.NewAIUtilityClient(client.AiUtilityGrpcClient.Connection)
	}
	client.NetworkConfig = networkConfig
	log.Println("platformGrpcClient: ", client.platformGrpcClient)
	log.Println("aiStudioGrpcClient: ", client.aiStudioGrpcClient)
	log.Println("NabhikServer: ", client.DataServerClient)
	log.Println("AiUtilityServerClient: ", client.AIUtilityServerClient)
	return client, nil
}

func SetupAIStudioClient(ctx context.Context, dns string, self string, logger zerolog.Logger) (pb.AIStudioClient, error) {
	// dns = "localhost:9013"
	telnet2, err := net.DialTimeout("tcp", dns, 1*time.Second)
	if err != nil {
		if self != "ai-studio" {
			return nil, err
		}
	}
	if telnet2 != nil {
		defer telnet2.Close()
	}
	grpcClient := pbtools.NewGrpcClient(logger,
		pbtools.ClientWithInsecure(true),
		pbtools.ClientWithServiceAddress(dns))
	return pb.NewAIStudioClient(grpcClient.Connection), nil
}

func (x *VapusSvcInternalClients) Close() {
	x.PlConn = nil
	x.UserConn = nil
	x.OrganizationConn = nil
	x.AIStudioConn = nil
	x.aiStudioGrpcClient = nil
	x.platformGrpcClient = nil
	x.DataServerClient = nil
	x.AgentServiceClient = nil
	x.AgentStudioClient = nil
}

func (x *VapusSvcInternalClients) PingTestAndReconnect(ctx context.Context, dns string, logger zerolog.Logger) error {
	err := dmutils.Telnet("tcp", dns)
	if err == nil {
		w, err := SetupVapusSvcInternalClients(ctx, x.NetworkConfig, "", logger)
		if err != nil {
			return err
		}
		x = w
		return nil
	}
	return err
}
