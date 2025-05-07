package secretstore

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	tpaws "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/aws"
	tpazure "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/azure"
	tpgcp "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/gcp"
	tphcvault "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/hcvault"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type SecretStore interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
	UpdateSecret(ctx context.Context, credData any, secretId string) error
	Close()
}

type SecretStoreClient struct {
	DataStoreParams *models.DataSourceCredsParams
	Logger          zerolog.Logger
	Debug           bool
	SecretStore
}

func (d *SecretStoreClient) Close() {
	// TODO: Implement this
}

type Options func(*SecretStoreClient)

func WithDebug(debug bool) Options {
	return func(d *SecretStoreClient) {
		d.Debug = debug
	}
}

func WithLogger(log zerolog.Logger) Options {
	return func(d *SecretStoreClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) Options {
	return func(d *SecretStoreClient) {
		d.DataStoreParams = params
	}
}

var SecretParentKey = "secret_value"

func New(ctx context.Context, opts ...Options) (*SecretStoreClient, error) {
	resultCl := &SecretStoreClient{}
	for _, opt := range opts {
		opt(resultCl)
	}
	if opts == nil || resultCl.DataStoreParams.DataSourceEngine == "" || resultCl.DataStoreParams.DataSourceCreds == nil {
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
	switch resultCl.DataStoreParams.DataSourceService {
	case mpb.DataSourceServices_HASHICORP_VAULT.String():
		var se string
		if val, ok := resultCl.DataStoreParams.Params[types.SECRETENGINE]; ok {
			se = val.(string)
		}
		client, err := tphcvault.NewHcVaultManager(ctx, &tphcvault.Vault{
			URL: resultCl.DataStoreParams.DataSourceCreds.URL,
			// AuthAppRole:     conf.HashicorpVault.AppRoleAuthnEnabled,
			// ApproleRoleID:   conf.HashicorpVault.AppRoleID,
			// ApproleSecretID: conf.HashicorpVault.AppRoleSecret,
			Token:        resultCl.DataStoreParams.DataSourceCreds.ApiToken,
			SecretEngine: se,
		})
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error creating vault client")
			return nil, err
		}
		resultCl.SecretStore = client
		return resultCl, nil
	case mpb.DataSourceServices_AWS_SECRET_MANAGER.String():
		client, err := tpaws.NewAwsSmClient(ctx, &tpaws.AWSConfig{
			Region:          resultCl.DataStoreParams.DataSourceCreds.AwsCreds.Region,
			AccessKeyId:     resultCl.DataStoreParams.DataSourceCreds.AwsCreds.AccessKeyId,
			SecretAccessKey: resultCl.DataStoreParams.DataSourceCreds.AwsCreds.SecretAccessKey,
		})
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error creating aws secret manager client")
			return nil, err
		}
		resultCl.SecretStore = client
		return resultCl, nil
	case mpb.DataSourceServices_GCP_SECRET_MANAGER.String():
		decodeData, err := base64.StdEncoding.DecodeString(resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ServiceAccountKey)
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		client, err := tpgcp.NewGcpSMStore(ctx, &tpgcp.GcpConfig{
			ServiceAccountKey: []byte(decodeData),
			ProjectID:         resultCl.DataStoreParams.DataSourceCreds.GcpCreds.ProjectId,
			Region:            resultCl.DataStoreParams.DataSourceCreds.GcpCreds.Region,
		})
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error creating gcp secret manager client")
			return nil, err
		}
		resultCl.SecretStore = client
		return resultCl, nil
	case mpb.DataSourceServices_AZURE_SECRET_MANAGER.String():
		client, err := tpazure.NewAzureKeyVault(ctx, &tpazure.AzureConfig{
			TenantID:     resultCl.DataStoreParams.DataSourceCreds.AzureCreds.TenantId,
			ClientID:     resultCl.DataStoreParams.DataSourceCreds.AzureCreds.ClientId,
			ClientSecret: resultCl.DataStoreParams.DataSourceCreds.AzureCreds.ClientSecret,
		}, resultCl.DataStoreParams.DataSourceCreds.URL)
		if err != nil {
			resultCl.Logger.Err(err).Msg("Error creating azure key vault client")
			return nil, err
		}
		resultCl.SecretStore = client
		return resultCl, nil
	default:
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
}

// WriteSecret writes the secret to the secret store
func (be *SecretStoreClient) WriteSecret(ctx context.Context, secrdData any, name string) error {
	bytes, err := json.Marshal(secrdData)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonMarshel, err)
	}
	return be.SecretStore.WriteSecret(ctx, map[string]any{SecretParentKey: string(bytes)}, name)
}

// ReadSecret reads the secret from the secret store
func (be *SecretStoreClient) ReadSecret(ctx context.Context, secretId string) (any, error) {
	secVal, err := be.SecretStore.ReadSecret(ctx, secretId)
	if err != nil {
		return "", err
	}
	origVal := dmutils.AnyToStr(secVal)
	result := map[string]any{}
	err = json.Unmarshal([]byte(origVal), &result)
	if err != nil {
		be.Logger.Err(err).Msg("Error unmarshalling secret data")
		return "", dmerrors.DMError(dmerrors.ErrJsonUnMarshel, err)
	}
	val, ok := result[SecretParentKey]
	if !ok {
		be.Logger.Err(err).Msg("Error reading secret data, secret key not found")
		return origVal, nil
	}
	_, ok = val.(string)
	if !ok {
		be.Logger.Err(err).Msg("Error reading secret data, secret value not found")
		return "", dmerrors.DMError(dputils.ErrInvalidSecretData, nil)
	}
	return val, nil
}

// DeleteSecret deletes the secret from the secret store
func (be *SecretStoreClient) DeleteSecret(ctx context.Context, secretId string) error {
	return be.SecretStore.DeleteSecret(ctx, secretId)
}

func (be *SecretStoreClient) UpdateSecret(ctx context.Context, secrdData any, secretId string) error {
	return be.SecretStore.UpdateSecret(ctx, secrdData, secretId)
}
