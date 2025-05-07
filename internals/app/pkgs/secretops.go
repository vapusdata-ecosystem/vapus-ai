package apppkgs

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

type DataWorkerSecretsRequest struct {
	ESource, LSource, DSource *models.DataSource
	EDB, LDB, FLDB            string
}

func GetDataSourceCreds(ctx context.Context, dataSource *models.DataSource, store *VapusStore, logger zerolog.Logger) (*models.DataSourceCredsParams, error) {
	secrets := &models.DataSourceSecrets{}
	creds, err := dataSource.GetCredentials("", true, "")
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("Error while getting the credentials for the node")
		return nil, dmerrors.DMError(ErrDataSourceCredsNotFound, err)
	}

	genericCreds, err := ReadCredentialFromStore(ctx, creds.SecretName, store, logger)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("Error while reading the secret from the vault")
		return nil, dmerrors.DMError(ErrDataSourceCredsSecretGet, err)
	}

	secrets.DB = creds.DB
	secrets.Port = int64(dataSource.NetParams.Port)
	secrets.URL = dataSource.NetParams.Address
	secrets.GenericCredentialModel = genericCreds
	secrets.Version = dataSource.NetParams.Version

	return &models.DataSourceCredsParams{
		DataSourceCreds:   secrets,
		DataSourceEngine:  dataSource.StorageEngine,
		DataSourceType:    dataSource.DataSourceType,
		DataSourceService: dataSource.ServiceName,
	}, nil
}

func ReadCredentialFromStore(ctx context.Context, secretName string, store *VapusStore, logger zerolog.Logger) (*models.GenericCredentialModel, error) {
	secrets := &models.GenericCredentialModel{}
	origVal, err := store.ReadSecret(ctx, secretName)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("Error while reading the secret from the vault")
		return nil, dmerrors.DMError(ErrDataSourceCredsSecretGet, err)
	}

	err = json.Unmarshal([]byte(dmutils.AnyToStr(origVal)), secrets)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while unmarshelling creds from secret store.")
		return nil, dmerrors.DMError(ErrDataSourceCredsSecretGet, err)
	}

	return secrets, nil
}

func GetDataCredsFromSecret(ctx context.Context, secretName string, store *VapusStore, logger zerolog.Logger) (*models.DataSourceCreds, error) {
	logger.Debug().Ctx(ctx).Msgf("Getting secret for %v", secretName)
	gCred, err := ReadCredentialFromStore(ctx, secretName, store, logger)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while reading the secrets from secret store")
		return nil, err
	}
	return &models.DataSourceCreds{
		SecretName:          secretName,
		Credentials:         gCred,
		IsAlreadyInSecretBs: true,
		Name:                secretName,
	}, nil
}

func SaveCredentialsCreds(ctx context.Context, secretName string, creds *models.GenericCredentialModel, store *VapusStore, logger zerolog.Logger) error {
	result, err := dmutils.StructToMap(creds)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while converting struct to map")
		return err
	}

	err = store.SecretStore.WriteSecret(ctx, result, secretName)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while writing secret %v", secretName)
		return err
	}
	return nil
}

func BuildDataWorkerSecrets(ctx context.Context, params *DataWorkerSecretsRequest, store *VapusStore, logger zerolog.Logger) (*models.DataWorkerSecretOpts, error) {
	var err error
	eCreds, err := GetDataCredsFromSecret(ctx, params.ESource.GetNetParams().DsCreds[0].SecretName, store, logger)
	if err != nil {
		return nil, err
	}
	lCreds, err := GetDataCredsFromSecret(ctx, params.ESource.GetNetParams().DsCreds[0].SecretName, store, logger)
	if err != nil {
		return nil, err
	}

	dCreds, err := GetDataCredsFromSecret(ctx, params.DSource.GetNetParams().DsCreds[0].SecretName, store, logger)
	if err != nil {
		return nil, err
	}

	extraction := &models.DataSourceSecretModel{
		DataSourceId: params.ESource.GetDataSourceId(),
		DataSourceCredsParams: &models.DataSourceCredsParams{
			DataSourceEngine:      params.ESource.StorageEngine,
			DataSourceType:        params.ESource.DataSourceType,
			DataSourceService:     params.ESource.ServiceName,
			DataSourceSvcProvider: params.ESource.ServiceProvider,
			DataSourceCreds: &models.DataSourceSecrets{
				GenericCredentialModel: eCreds.Credentials,
				URL:                    params.ESource.GetNetParams().Address,
				DB:                     params.EDB,
				Port:                   int64(params.ESource.GetNetParams().Port),
			},
		},
	}

	loading := &models.DataSourceSecretModel{
		DataSourceId: params.ESource.GetDataSourceId(),
		DataSourceCredsParams: &models.DataSourceCredsParams{
			DataSourceEngine: params.LSource.StorageEngine,
			DataSourceType:   params.LSource.DataSourceType,
			DataSourceCreds: &models.DataSourceSecrets{
				GenericCredentialModel: lCreds.Credentials,
				URL:                    params.LSource.GetNetParams().Address,
				DB:                     params.LDB,
				Port:                   int64(params.LSource.GetNetParams().Port),
			},
		},
	}

	destination := &models.DataSourceSecretModel{
		DataSourceId: params.ESource.GetDataSourceId(),
		DataSourceCredsParams: &models.DataSourceCredsParams{
			DataSourceEngine: params.DSource.StorageEngine,
			DataSourceType:   params.DSource.DataSourceType,
			DataSourceCreds: &models.DataSourceSecrets{
				GenericCredentialModel: dCreds.Credentials,
				URL:                    params.DSource.GetNetParams().Address,
				DB:                     params.FLDB,
				Port:                   int64(params.DSource.GetNetParams().Port),
			},
		},
	}

	return &models.DataWorkerSecretOpts{
		Extraction:  extraction,
		Loading:     loading,
		Destination: destination,
	}, nil
}
