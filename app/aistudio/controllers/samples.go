package dmcontrollers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog"
	grpccodes "google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"

	faker "github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"

	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
)

var aiModelNode *mpb.AIModelNode = &mpb.AIModelNode{
	Attributes: &mpb.AIModelNodeAttributes{
		NetworkParams: &mpb.AIModelNodeNetworkParams{
			Credentials: &mpb.GenericCredentialObj{},
		},
	},
}

var dataSourceCreds *mpb.GenericCredentialObj = &mpb.GenericCredentialObj{
	AwsCreds:   &mpb.AWSCreds{},
	GcpCreds:   &mpb.GCPCreds{},
	AzureCreds: &mpb.AzureCreds{},
}

var organizationObj *mpb.Organization = &mpb.Organization{
	SecretPasscode: &mpb.CredentialSalt{},
	Attributes: &mpb.OrganizationAttributes{
		AuthnJwtParams: &mpb.JWTParams{},
	},
	BackendSecretStorage: &mpb.BackendStorages{},
	ArtifactStorage:      &mpb.BackendStorages{},
}

var dataSourceObj *mpb.DataSource = &mpb.DataSource{
	NetParams: &mpb.DataSourceNetParams{
		DsCreds: []*mpb.DataSourceCreds{
			{
				Credentials: &mpb.GenericCredentialObj{
					AwsCreds:   &mpb.AWSCreds{},
					GcpCreds:   &mpb.GCPCreds{},
					AzureCreds: &mpb.AzureCreds{},
				},
			},
		},
	},
	Attributes:    &mpb.DataSourceAttributes{},
	Tags:          []*mpb.Mapper{{}},
	SharingParams: &mpb.DataSourceSharingParams{},
}

func (dmc *VapusDataController) GetSampleResourceConfiguration(ctx context.Context, request *pb.SampleResourceConfigurationOptions) (*pb.SampleResourceConfiguration, error) {
	output := []*pb.SampleResourceConfiguration_ResourceConfigs{}
	switch request.RequestObj {

	case mpb.Resources_AIMODELS:
		obj, err := generateAIModelNodeSpec(request.GetFormat().String(), request.PopulateFakeData, dmc.Logger)
		if err != nil {
			return returnSampleResourceConfigurationErr(request, err)
		}
		return &pb.SampleResourceConfiguration{
			Output: []*pb.SampleResourceConfiguration_ResourceConfigs{returnSampleResourceConfiguration(request.Format, obj, request.GetRequestObj())},
		}, nil

	case mpb.Resources_ALL:
		aiModelNodeSpec, _ := generateAIModelNodeSpec(request.GetFormat().String(), request.PopulateFakeData, dmc.Logger)
		output = append(output, returnSampleResourceConfiguration(request.Format, aiModelNodeSpec, mpb.Resources_AIMODELS))
		return &pb.SampleResourceConfiguration{
			Output: output,
		}, nil
	default:
		dmc.Logger.Error().Msgf("unsupported request object - %v", request.GetRequestObj())
		return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrInvalidRequestObj, nil), grpccodes.InvalidArgument)
	}
}

func returnSampleResourceConfigurationErr(req *pb.SampleResourceConfigurationOptions, err error) (*pb.SampleResourceConfiguration, error) {
	errs := dmerrors.DMError(err, nil)
	return &pb.SampleResourceConfiguration{
		Output: []*pb.SampleResourceConfiguration_ResourceConfigs{{
			FileContent: "",
			RequestObj:  req.GetRequestObj(),
			Format:      req.GetFormat(),
		}},
	}, pbtools.HandleGrpcError(errs, grpccodes.Internal)
}

func returnSampleResourceConfiguration(format mpb.ContentFormats, spec []byte, requestObj mpb.Resources) *pb.SampleResourceConfiguration_ResourceConfigs {
	return &pb.SampleResourceConfiguration_ResourceConfigs{
		FileContent: string(spec),
		RequestObj:  requestObj,
		Format:      format,
	}
}

func generateDataSourceSpec(format string, hasFakeData bool, logger zerolog.Logger) ([]byte, error) {
	if hasFakeData {
		err := faker.FakeData(dataSourceObj, options.WithTagName(strings.ToLower(format)), options.WithRecursionMaxDepth(1))
		if err != nil {
			logger.Err(err).Msg("error while generating empty request file format for dataSourceObj with fake data")
			return nil, apperr.ErrGeneratingRequestFiles
		}
	}
	switch format {
	case mpb.ContentFormats_YAML.String():
		bytesData, err := yaml.Marshal(dataSourceObj)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for dataSourceObj")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	case mpb.ContentFormats_JSON.String():
		bytesData, err := json.Marshal(dataSourceObj)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for dataSourceObj")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	// case mpb.ContentFormats_TOML.String():
	default:
		logger.Error().Msgf("unsupported file format - %v", format)
		return nil, apperr.ErrInvalidFileFormat
	}
}

func generateOrganizationSpec(format string, hasFakeData bool, logger zerolog.Logger) ([]byte, error) {
	if hasFakeData {
		err := faker.FakeData(organizationObj, options.WithTagName(strings.ToLower(format)), options.WithRecursionMaxDepth(1))
		if err != nil {
			logger.Err(err).Msg("error while generating empty request file format for organizationObj with fake data")
			return nil, apperr.ErrGeneratingRequestFiles
		}
	}
	switch format {
	case mpb.ContentFormats_YAML.String():
		bytesData, err := yaml.Marshal(organizationObj)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for organizationObj")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	case mpb.ContentFormats_JSON.String():
		bytesData, err := json.Marshal(organizationObj)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for organizationObj")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	// case mpb.ContentFormats_TOML.String():
	default:
		logger.Error().Msgf("unsupported file format - %v", format)
		return nil, apperr.ErrInvalidFileFormat
	}
}

func generateDataSourceCredsSpec(format string, hasFakeData bool, logger zerolog.Logger) ([]byte, error) {
	if hasFakeData {
		err := faker.FakeData(dataSourceCreds, options.WithTagName(strings.ToLower(format)), options.WithRecursionMaxDepth(1))
		if err != nil {
			logger.Err(err).Msg("error while generating empty request file format for dataSourceCreds with fake data")
			return nil, apperr.ErrGeneratingRequestFiles
		}
	}
	switch format {
	case mpb.ContentFormats_YAML.String():
		bytesData, err := yaml.Marshal(dataSourceCreds)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for dataSourceCreds")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	case mpb.ContentFormats_JSON.String():
		bytesData, err := json.Marshal(dataSourceCreds)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for dataSourceCreds")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	// case mpb.ContentFormats_TOML.String():
	default:
		logger.Error().Msgf("unsupported file format - %v", format)
		return nil, apperr.ErrInvalidFileFormat
	}
}

func generateAIModelNodeSpec(format string, hasFakeData bool, logger zerolog.Logger) ([]byte, error) {
	if hasFakeData {
		err := faker.FakeData(aiModelNode, options.WithTagName(strings.ToLower(format)), options.WithRecursionMaxDepth(1))
		if err != nil {
			logger.Err(err).Msg("error while generating empty request file format for aiModelNode with fake data")
			return nil, apperr.ErrGeneratingRequestFiles
		}
	}
	switch format {
	case mpb.ContentFormats_YAML.String():
		bytesData, err := yaml.Marshal(aiModelNode)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for aiModelNode")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	case mpb.ContentFormats_JSON.String():
		bytesData, err := json.Marshal(aiModelNode)
		if err != nil {
			logger.Err(err).Msg("error while marshelling empty request file format for aiModelNode")
			return nil, apperr.ErrGeneratingRequestFiles
		}
		return bytesData, nil
	// case mpb.ContentFormats_TOML.String():
	default:
		logger.Error().Msgf("unsupported file format - %v", format)
		return nil, apperr.ErrInvalidFileFormat
	}
}
