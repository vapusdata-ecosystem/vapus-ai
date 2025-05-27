package booter

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	apppdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type PlatformSetup struct {
	bootConfig          *appconfigs.PlatformBootConfig
	dbManager           *apppkgs.VapusStore
	accountId           string
	ownerPolicies       []string
	accountOrganization string
	marketplaceId       string
	fakeCtxClaim        map[string]string
	logger              zerolog.Logger
	svcPackages         *apppkgs.VapusSvcPackages
	svcPkgsParams       *apppkgs.VapusSvcPackageParams
	artifactStoreCred   *models.DataSourceCredsParams
}

func NewPlatformSetup(bc *appconfigs.PlatformBootConfig, dbm *apppkgs.VapusStore, svcPkgs *apppkgs.VapusSvcPackages, svcPkgsParams *apppkgs.VapusSvcPackageParams, logger zerolog.Logger) *PlatformSetup {
	return &PlatformSetup{
		ownerPolicies: func() []string {
			if svcPkgs != nil {
				if svcPkgs.PlatformRBACManager != nil {
					for _, role := range svcPkgs.PlatformRBACManager.Roles {
						if role.Name == "platformOwners" {
							return role.Policies
						}
					}
				}
			}
			return []string{}
		}(),
		bootConfig:    bc,
		dbManager:     dbm,
		logger:        logger,
		svcPackages:   svcPkgs,
		svcPkgsParams: svcPkgsParams,
	}
}

func (p *PlatformSetup) Clean() {
	p.fakeCtxClaim = nil
}

func (p *PlatformSetup) AddVapusDataPlatformOwners(ctx context.Context) error {
	count, err := p.dbManager.Db.CountRows(ctx, &datasvcpkgs.QueryOpts{
		DataCollection: apppkgs.UsersTable,
	}, p.logger)
	if err != nil {
		p.logger.Fatal().Msgf("error while fetching the existing val of account created. - %v", err)
	}
	if count == 0 {
		for _, user := range p.bootConfig.PlatformOwners {
			log.Println("Creating platform owner for account ", p.accountId)
			user, err := apppdrepo.CreateUser(ctx, &apppdrepo.LocalUserM{
				Email:             user,
				AccountId:         p.accountId,
				Organization:      p.accountOrganization,
				OrganizationRoles: []string{mpb.UserRoles_SERVICE_OWNER.String()},
			}, []string{mpb.UserRoles_SERVICE_OWNER.String()}, nil, map[string]string{encryption.ClaimAccountKey: p.accountId, encryption.ClaimUserIdKey: user},
				p.dbManager, p.logger)
			if err != nil {
				p.logger.Fatal().Err(err).Msgf("error while creating platform owner with email '%v'", user)
			} else {
				p.logger.Info().Msgf("User '%v' is created successfully.......", user)
			}
		}
	}
	return nil
}

func (p *PlatformSetup) AddVapusDataPlatformAccount(ctx context.Context) error {
	count, err := p.dbManager.Db.CountRows(ctx, &datasvcpkgs.QueryOpts{
		DataCollection: apppkgs.AccountsTable,
	}, p.logger)
	if err != nil {
		p.logger.Fatal().Msgf("error while fetching the existing val of account created. - %v", err)
	}
	p.logger.Info().Msgf("Authn Method - %v", p.svcPkgsParams.AuthnParams.AuthnMethod)
	p.logger.Info().Msgf("Authn Method - %v", mpb.AuthnMethod(mpb.AuthnMethod_value[strings.ToUpper(p.svcPkgsParams.AuthnParams.AuthnMethod)]))
	if count == 0 {
		account := &models.Account{
			Name:         p.bootConfig.PlatformAccount.Name,
			Status:       mpb.CommonStatus_ACTIVE.String(), // TO:DO add logic for activating and deactivating account
			AuthnMethod:  p.svcPkgsParams.AuthnParams.AuthnMethod,
			Profile:      &models.AccountProfile{},
			AIAttributes: &models.AccountAIAttributes{},
		}
		account.SetAccountId()
		account.PreSaveCreate(p.bootConfig.PlatformAccount.Creator)
		account.BackendDataStorage = p.getBeDbStorage(ctx, account.Name)
		account.BackendSecretStorage = p.getBeSecretStorage(ctx, account.Name)
		// account.ArtifactStorage = p.getArtifactStorage(ctx, account.Name)
		account.DmAccessJwtKeys = p.getJwtAccessKeys(ctx, account.Name)

		_, err := apppdrepo.CreateAccount(ctx, account, p.dbManager, p.logger)
		if err != nil {
			p.logger.Fatal().Msgf("error while booting the VapusData account creation for account %v . error: %v", p.bootConfig.PlatformAccount.Name, err)
		}
		p.logger.Info().Msgf("Account successfully created for '%v' with creator set as '%v'.......", p.bootConfig.PlatformAccount.Name, p.bootConfig.PlatformAccount.Creator)
		p.accountId = account.VapusID
		p.logger.Info().Msgf("Account ID for this setup is '%v'", p.accountId)
		return nil
	}
	p.logger.Info().Msg("Account for this setup already exists")
	return nil
}

func (p *PlatformSetup) AddVapusDataPlatformOwnerOrganization(ctx context.Context) error {
	count, err := p.dbManager.Db.CountRows(ctx, &datasvcpkgs.QueryOpts{
		DataCollection: apppkgs.OrganizationsTable,
	}, p.logger)
	if err != nil {
		p.logger.Fatal().Msgf("error while fetching the existing val of account created. - %v", err)
	}
	var OrganizationName string
	if p.bootConfig.PlatformAccountOrganization.Name == "" {
		OrganizationName = p.bootConfig.PlatformAccount.Name + "-platformOrganization"
	} else {
		OrganizationName = p.bootConfig.PlatformAccountOrganization.Name
	}

	if count == 0 && len(p.bootConfig.PlatformOwners) > 0 {
		OrganizationObj := &models.Organization{
			Name:  OrganizationName,
			Users: []string{},
		}
		OrganizationObj.Users = append(OrganizationObj.Users, p.bootConfig.PlatformOwners...)
		OrganizationObj.OrganizationType = mpb.OrganizationType_SERVICE_ORGANIZATION.String()
		OrganizationObj.PreSaveCreate(map[string]string{encryption.ClaimAccountKey: p.accountId,
			encryption.ClaimUserIdKey: p.bootConfig.PlatformOwners[0]})
		OrganizationObj.SetOrganizationId()
		OrganizationObj.SecretPasscode = "Organization_" + dmutils.SlugifyBase(OrganizationObj.VapusID)
		OrganizationObj.Status = mpb.CommonStatus_ACTIVE.String()

		OrganizationObj.SetAccountId(p.accountId)
		err = apppdrepo.ConfigureOrganization(ctx, OrganizationObj, p.fakeCtxClaim, p.dbManager, p.logger)
		if err != nil {
			p.logger.Fatal().Msgf("error while booting the account Organization %v . error: %v", OrganizationObj.Name, err)
		}

		p.logger.Info().Msgf("Organization '%v' successfully created with Organization owner set as '%v'.......", OrganizationObj.Name, OrganizationObj.CreatedBy)
		p.accountOrganization = OrganizationObj.VapusID
	} else {
		p.logger.Info().Msg("Organization for this setup already exists")
	}
	return nil
}

func (p *PlatformSetup) getArtifactStorage(ctx context.Context, accountName string) *models.BackendStorages {
	secName := fmt.Sprintf("%d-artifactstore", dmutils.GetEpochTime())
	log.Println("Artifact Storage Manager: ", p.artifactStoreCred.DataSourceCreds.GenericCredentialModel)
	err := p.dbManager.SecretStore.WriteSecret(ctx,
		p.artifactStoreCred.DataSourceCreds.GenericCredentialModel, secName)
	if err != nil {
		if strings.Contains(err.Error(), "code = AlreadyExists") {
			p.logger.Info().Msgf("Secret %v already exists", secName)
		} else {
			p.logger.Fatal().Err(err).Msgf("error while storing artifact creds info in secret storage.")
			return nil
		}
	}
	result := &models.BackendStorages{
		BesType:       mpb.DataSourceType_ARTIFACT.String(),
		BesOnboarding: mpb.BackendStorageOnboarding_BE_DEFAULT_PLATFORM.String(),
		BesService:    p.artifactStoreCred.DataSourceService,
		NetParams: &models.DataSourceNetParams{
			Address: p.artifactStoreCred.DataSourceCreds.URL,
			DsCreds: []*models.DataSourceCreds{
				{
					SecretName: secName,
					Name:       accountName + " artifact store",
				},
			},
		},
		Status: mpb.CommonStatus_ACTIVE.String(),
	}
	return result
}

// Make below logic more generic
func (p *PlatformSetup) getBeDbStorage(ctx context.Context, accountName string) *models.BackendStorages {
	secName := fmt.Sprintf("%d-dbstore", dmutils.GetEpochTime())
	err := p.dbManager.SecretStore.WriteSecret(ctx,
		p.dbManager.GetDbStoreParams().DataSourceCreds.GenericCredentialModel,
		secName)
	// err := p.dbManager.SecretStore.WriteSecret(ctx, &models.GenericCredentialModel{
	// 	ApiToken:     p.dbManager.GetDbStoreParams().DataSourceCreds.ApiToken,
	// 	ApiTokenType: mpb.ApiTokenType_APIKEY.String(),
	// }, secName)
	if err != nil {
		if strings.Contains(err.Error(), "code = AlreadyExists") {
			p.logger.Info().Msgf("Secret %v already exists", secName)
		} else {
			p.logger.Fatal().Err(err).Msgf("error while storing db store info in secret storage.")
			return nil
		}
	}
	result := &models.BackendStorages{
		BesType:       mpb.DataSourceType_DATABASE.String(),
		BesOnboarding: mpb.BackendStorageOnboarding_BE_DEFAULT_PLATFORM.String(),
		BesService:    p.dbManager.GetDbStoreParams().DataSourceService,
		NetParams: &models.DataSourceNetParams{
			Address: p.dbManager.GetDbStoreParams().DataSourceCreds.URL,
			DsCreds: []*models.DataSourceCreds{
				{
					SecretName: secName,
					Name:       accountName + " db store",
				},
			},
		},
		Status: mpb.CommonStatus_ACTIVE.String(),
	}
	return result
}

// TODO: Make below logic more generic for diff storage types
func (p *PlatformSetup) getBeSecretStorage(ctx context.Context, accountName string) *models.BackendStorages {
	secName := fmt.Sprintf("%d-secretstore", dmutils.GetEpochTime())
	log.Println("Secret Store Manager: ", p.dbManager.GetCreds().DataSourceCreds.GenericCredentialModel)
	// err := p.dbManager.SecretStore.WriteSecret(ctx, &models.GenericCredentialModel{
	// 	ApiToken:     p.dbManager.GetDbStoreParams().DataSourceCreds.ApiToken,
	// 	ApiTokenType: mpb.ApiTokenType_APIKEY.String(),
	// }, secName)
	err := p.dbManager.SecretStore.WriteSecret(ctx,
		p.dbManager.GetCreds().DataSourceCreds.GenericCredentialModel,
		secName)
	if err != nil {
		if strings.Contains(err.Error(), "code = AlreadyExists") {
			p.logger.Info().Msgf("Secret %v already exists", secName)
		} else {
			p.logger.Fatal().Err(err).Msgf("error while storing secret store db info in secret storage.")
			return nil
		}
	}
	result := &models.BackendStorages{
		BesType:       mpb.DataSourceType_DATABASE.String(),
		BesOnboarding: mpb.BackendStorageOnboarding_BE_DEFAULT_PLATFORM.String(),
		BesService:    p.dbManager.GetCreds().DataSourceService,
		NetParams: &models.DataSourceNetParams{
			Address: p.dbManager.GetCreds().DataSourceCreds.URL,
			DsCreds: []*models.DataSourceCreds{
				{
					SecretName: secName,
					Name:       accountName + " secret store",
				},
			},
		},
		Status: mpb.CommonStatus_ACTIVE.String(),
	}
	return result
}

func (p *PlatformSetup) getJwtAccessKeys(ctx context.Context, accountName string) *models.JWTParams {
	secretName := fmt.Sprintf("%d-jwtaccessStore", dmutils.GetEpochTime())
	err := p.dbManager.SecretStore.WriteSecret(ctx, &encryption.JWTAuthn{
		PrivateJWTKey:    p.svcPackages.VapusJwtAuth.Opts.PrivateJWTKey,
		PublicJWTKey:     p.svcPackages.VapusJwtAuth.Opts.PublicJWTKey,
		SigningAlgorithm: p.svcPackages.VapusJwtAuth.Opts.SigningAlgorithm,
	}, secretName)
	if err != nil {
		if strings.Contains(err.Error(), "code = AlreadyExists") {
			p.logger.Info().Msgf("Secret %v already exists", secretName)
		} else {
			p.logger.Fatal().Err(err).Msgf("error while storing jwt info in secret storage.")
			return nil
		}
	}
	return &models.JWTParams{
		Name:                secretName,
		Status:              mpb.CommonStatus_ACTIVE.String(),
		VId:                 secretName,
		SigningAlgorithm:    p.svcPackages.VapusJwtAuth.Opts.SigningAlgorithm,
		IsAlreadyInSecretBs: true,
	}
}
