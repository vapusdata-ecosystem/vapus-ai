package services

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmstores "github.com/vapusdata-ecosystem/vapusai/aistudio/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/options"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

// TODO Repalce below 2 logic to directly store from account obj, and load account obj while boot

func setOrganizationArtifactBEStore(ctx context.Context, organization *models.Organization, dbstore *aidmstore.AIStudioDMStore) (*models.BackendStorages, error) {
	var err error
	secName := dmutils.GetSecretName("organization", organization.VapusID, "artifactStoreSecret")
	if organization.ArtifactStorage != nil && organization.ArtifactStorage.NetParams != nil && len(organization.ArtifactStorage.NetParams.DsCreds) > 0 {
		if organization.ArtifactStorage.NetParams.DsCreds[0].Credentials == nil || organization.ArtifactStorage.NetParams.DsCreds[0].IsAlreadyInSecretBs {
			return nil, nil
		}
		err = dbstore.SecretStore.WriteSecret(ctx, organization.ArtifactStorage.NetParams.DsCreds[0].Credentials, secName)
		if err != nil {
			helperLogger.Err(err).Msgf("error while storing artifact store secrets in secret storage.")
			return nil, err
		}
		netParams := &models.DataSourceNetParams{
			DsCreds: []*models.DataSourceCreds{
				{
					IsAlreadyInSecretBs: true,
					SecretName:          secName,
				},
			},
			Address: organization.ArtifactStorage.NetParams.Address,
		}
		organization.ArtifactStorage.NetParams = netParams
		organization.ArtifactStorage.Status = mpb.CommonStatus_ACTIVE.String()
		return organization.ArtifactStorage, nil
	} else {
		// return dmstores.AccountPool.ArtifactStorage, nil
		return &models.BackendStorages{
			BesType:       dmstores.DMStoreManager.Account.ArtifactStorage.BesType,
			BesOnboarding: dmstores.DMStoreManager.Account.ArtifactStorage.BesOnboarding,
			BesService:    dmstores.DMStoreManager.Account.ArtifactStorage.BesService,
			NetParams:     dmstores.DMStoreManager.Account.ArtifactStorage.NetParams,
			Status:        dmstores.DMStoreManager.Account.ArtifactStorage.Status,
		}, nil

	}
}

func setOrganizationDPK8sInfra(ctx context.Context, organization *models.Organization, infraList []*models.K8SInfraParams, dbstore *aidmstore.AIStudioDMStore) error {
	var err error
	if len(infraList) > 0 {
		for _, infra := range infraList {
			if infra.SecretName == "" {
				secName := dmutils.GetSecretName("k8sInfraParams", organization.VapusID, dmutils.GetStrEpochTime())
				err = dbstore.SecretStore.WriteSecret(ctx, infra.Credentials, secName)
				if err != nil {
					helperLogger.Err(err).Msgf("error while k8s infra params in secret storage for organization %v.", organization.VapusID)
					return err
				}
				infra.SecretName = secName
				infra.InfraId = dmutils.GetUUID()
			}
			infra.Credentials = nil

			// organization.DataProductInfraPlatform = append(organization.DataProductInfraPlatform, infra)
		}
	}
	return nil
}

func getPlatformBESecretStore(ctx context.Context, organization *models.Organization, dbstore *aidmstore.AIStudioDMStore) (*models.BackendStorages, error) {
	am := pkgs.VapusArtifactStorageManager.Spec
	var err error
	secName := dmutils.GetSecretName("organization", organization.VapusID, dmutils.GetStrEpochTime())
	err = dbstore.SecretStore.WriteSecret(ctx, am.DataSourceCreds.GenericCredentialModel, secName)
	if err != nil {
		helperLogger.Err(err).Msgf("error while storing artifact store secrets in secret storage.")
		return nil, err
	}
	netParams := &models.DataSourceNetParams{
		DsCreds: []*models.DataSourceCreds{
			{
				IsAlreadyInSecretBs: true,
				SecretName:          secName,
			},
		},
		Address: am.DataSourceCreds.URL,
	}

	return &models.BackendStorages{
		BesType:       am.DataSourceType,
		BesOnboarding: mpb.BackendStorageOnboarding_BE_DEFAULT_PLATFORM.String(),
		BesService:    am.DataSourceService,
		NetParams:     netParams,
		Status:        mpb.CommonStatus_ACTIVE.String(),
	}, nil
}

func organizationConfigureTool(ctx context.Context, organization *models.Organization, dbStore *aidmstore.AIStudioDMStore, logger zerolog.Logger, ctxClaim map[string]string) (*models.Organization, error) {
	var err error
	organization.VapusID = ""
	organization.SetOrganizationId()
	organization.Status = mpb.CommonStatus_ACTIVE.String()
	organization.PreSaveCreate(ctxClaim)
	jwtSecretName := dmutils.GetSecretName("organization", organization.VapusID, "authJwtParams")
	if organization.AuthnJwtParams != nil {
		if !organization.AuthnJwtParams.IsAlreadyInSecretBs {
			jwtParam, err := setJWTAuthzParams(ctx, jwtSecretName, dbStore.SecretStore, false, organization.AuthnJwtParams)
			if err != nil {
				return nil, dmerrors.DMError(apperr.ErrSavingOrganizationAuthJwt, err)
			}
			organization.AuthnJwtParams.Reset()
			organization.AuthnJwtParams = jwtParam
		}
	} else {
		jwtParam, err := setJWTAuthzParams(ctx, jwtSecretName, dbStore.SecretStore, true, nil)
		if err != nil {
			return nil, dmerrors.DMError(apperr.ErrSavingOrganizationAuthJwt, err)
		}
		organization.AuthnJwtParams = jwtParam
	}

	organization.SecretPasscode = dmutils.GenerateRandomString(16)

	if organization.BackendSecretStorage == nil {
		// TO:DO add logic to store account bestorage
		// secName := getSecretName("organization", organization.VapusID, "organizationBeSecretStore")
	}

	resp, err := setOrganizationArtifactBEStore(ctx, organization, dbStore)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while setting organization artifact store %v", organization)
		return nil, dmerrors.DMError(apperr.ErrSettingOrganizationArtifactStore, err) //nolint:wrapcheck
	}
	organization.ArtifactStorage = resp
	var errCh = make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = dbStore.BlobStore.CreateBucket(ctx, &options.BlobOpsParams{
			BucketName: organization.VapusID,
		})
	}()
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	organization.Users = []string{ctxClaim[encryption.ClaimUserIdKey]}
	organization.Editors = organization.Users
	err = dbStore.ConfigureOrganization(ctx, organization, ctxClaim)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while configuring organization %v", organization)
		return nil, dmerrors.DMError(apperr.ErrCreateOrganization, err) //nolint:wrapcheck
	}
	return organization, nil
}
