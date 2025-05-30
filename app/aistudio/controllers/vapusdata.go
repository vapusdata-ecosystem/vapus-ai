package dmcontrollers

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/rs/zerolog"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"

	types "github.com/vapusdata-ecosystem/vapusai/core/types"

	dmstores "github.com/vapusdata-ecosystem/vapusai/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	dmsvc "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	utils "github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type VapusDataController struct {
	pb.UnimplementedVapusdataServiceServer
	validator  *dmutils.DMValidator
	DMServices *dmsvc.AIStudioServices
	Logger     zerolog.Logger
}

var VapusDataControllerManager *VapusDataController

func NewVapusDataController() *VapusDataController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusDataController")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusDataController Controller initialized")
	return &VapusDataController{
		validator:  validator,
		Logger:     l,
		DMServices: dmsvc.AIStudioServiceManager,
	}
}

func InitVapusDataController() {
	if VapusDataControllerManager == nil {
		VapusDataControllerManager = NewVapusDataController()
	}
}

func (dmc *VapusDataController) PlatformPublicInfo(ctx context.Context, request *mpb.EmptyRequest) (*pb.PlatformPublicInfoResponse, error) {
	accountInfo := dmstores.DMStoreManager.Account
	fmt.Println("Account Info", dmstores.DMStoreManager.Account.Name)
	if accountInfo == nil {
		return &pb.PlatformPublicInfoResponse{}, grpcstatus.Error(grpccodes.NotFound, "Account details not found")
	}
	return &pb.PlatformPublicInfoResponse{
		Logo:        accountInfo.Profile.Logo,
		AccountName: accountInfo.Name,
		Favicon:     accountInfo.Profile.Favicon,
	}, nil
}

func (dmc *VapusDataController) AccountManager(ctx context.Context, request *pb.AccountManagerRequest) (*pb.AccountResponse, error) {
	agent, err := dmc.DMServices.NewAccountAgent(ctx, request)
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while initializing AccountManager")
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal)
	}
	err = agent.Act(ctx, "")
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while performing the action on AccountManager")
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal)
	}
	agent.LogAgent()
	return &pb.AccountResponse{
		Output: agent.GetResponse().ConvertToPb(),
		DmResp: pbtools.HandleDMResponse(ctx, utils.ACCOUNT_CREATED, "200"),
	}, nil
}

func (dmc *VapusDataController) AccountGetter(ctx context.Context, request *mpb.EmptyRequest) (*pb.AccountResponse, error) {
	accountInfo, err := dmc.DMServices.GetAccount(ctx)
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while getting account info")
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal)
	}
	return &pb.AccountResponse{
		Output: accountInfo.ConvertToPb(),
		DmResp: pbtools.HandleDMResponse(ctx, utils.ACCOUNT_CREATED, "200"),
	}, nil
}

func (x *VapusDataController) ResourceGetter(ctx context.Context, request *pb.ResourceGetterRequest) (*pb.ResourceGetterResponse, error) {
	// Not required to call the services
	specMap := appconfigs.SpecMap
	resourceActionsMap := appconfigs.ResourceActionsMap
	result := &pb.ResourceGetterResponse{}

	// Resource Response
	for key, val := range specMap {
		resourceResponse := &pb.ResourceGetterRequest{}
		resourceResponse.Resourse = key
		resourceResponse.YamlSpec = val
		actions := resourceActionsMap[key]
		resourceResponse.Actions = actions
		result.ResourceResponse = append(result.ResourceResponse, resourceResponse)
	}

	// Enum Response
	enumSpecs := appconfigs.EnumSpecs
	for key, val := range enumSpecs {
		enumRequest := &pb.EnumGetterRequest{}
		enumRequest.Name = key
		enumRequest.Value = slices.Collect(maps.Keys(val))
		result.EnumResponse = append(result.EnumResponse, enumRequest)
	}

	for key, val := range types.DataSourceTypeMap {
		dataSourceTypeMap := &pb.DataSourceTypeMap{}
		dataSourceTypeMap.Service = key
		dataSourceTypeMap.SourceType = val
		result.DataSourceTypeMap = append(result.DataSourceTypeMap, dataSourceTypeMap)
	}
	for key, val := range types.StorageEngineMap {
		storageEngineMap := &pb.StorageEngineMap{}
		storageEngineMap.Service = key
		storageEngineMap.Engine = val.String()
		result.StorageEngineMap = append(result.StorageEngineMap, storageEngineMap)
	}

	for key, val := range types.StorageEngineLogoMap {
		storageEngineLogoMap := &pb.StorageEngineLogoMap{}
		storageEngineLogoMap.Engine = key.String()
		storageEngineLogoMap.Url = val
		result.StorageEngineLogoMap = append(result.StorageEngineLogoMap, storageEngineLogoMap)
	}

	for key, val := range types.DataSourceServicesLogoMap {
		dataSourceServicesLogoMap := &pb.DataSourceServicesLogoMap{}
		dataSourceServicesLogoMap.Service = key
		dataSourceServicesLogoMap.Url = val
		result.DataSourceServicesLogoMap = append(result.DataSourceServicesLogoMap, dataSourceServicesLogoMap)
	}

	for key, val := range types.ServiceProviderLogoMap {
		serviceProviderLogoMap := &pb.ServiceProviderLogoMap{}
		serviceProviderLogoMap.ServiceProvider = key
		serviceProviderLogoMap.Url = val
		result.ServiceProviderLogoMap = append(result.ServiceProviderLogoMap, serviceProviderLogoMap)
	}

	for key, val := range appconfigs.PluginTypes {
		pluginTypeMap := &pb.PluginTypeMap{}
		pluginTypeMap.PluginTypes = key
		pluginTypeMap.Services = map[string]string{}
		for k, v := range val {
			pluginTypeMap.Services[k.String()] = v
		}
		result.PluginTypeMap = append(result.PluginTypeMap, pluginTypeMap)
	}

	thirdPartyGuardrailList := pb.ThirdPartyGuardrailList{
		Pangea:  []string{},
		Mistral: []string{},
		Bedrock: []string{},
		Vapus:   []string{},
	}
	// for _, val := range types.PanegaGuardrailList {
	// 	thirdPartyGuardrailList.Pangea = append(thirdPartyGuardrailList.Pangea, val.String())
	// }
	// for _, val := range types.MistralGuardrailList {
	// 	thirdPartyGuardrailList.Mistral = append(thirdPartyGuardrailList.Mistral, val.String())
	// }
	result.GuardrailTypes = &thirdPartyGuardrailList
	return result, nil
}
