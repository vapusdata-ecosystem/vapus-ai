package datarepo

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func CreateAccount(ctx context.Context, obj *models.Account, ds *apppkgs.VapusStore, logger zerolog.Logger) (*models.Account, error) {
	logger.Info().Msgf("Creating account : %v", obj)
	_, err := ds.Db.PostgresClient.DB.NewInsert().ModelTableExpr(apppkgs.AccountsTable).Model(obj).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving account in datastore")
		return nil, err
	}
	return obj, nil
}

func ConfigureOrganization(ctx context.Context, Organization *models.Organization, ctxClaim map[string]string, ds *apppkgs.VapusStore, logger zerolog.Logger) error {
	Organization.SetAccountId(ctxClaim[encryption.ClaimAccountKey])
	if ds.Cacher != nil {
		_, err := ds.BeDataStore.Cacher.RedisClient.WrtiteData(ctx, Organization.VapusID, types.EMPTYSTR, Organization)
		if err != nil {
			logger.Err(err).Ctx(ctx).Msg(apperr.ErrOrganizationInitialization.Error())
			return err
		}
	}
	_, err := ds.Db.PostgresClient.DB.NewInsert().ModelTableExpr(apppkgs.OrganizationsTable).Model(Organization).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving datamarketplace in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   Organization.VapusID,
			ResourceName: mpb.Resources_ORGANIZATIONS.String(),
			VapusBase: models.VapusBase{
				Editors: Organization.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func CreateUser(ctx context.Context, lu *LocalUserM, role []string, uo *models.Users, ctxClaim map[string]string, ds *apppkgs.VapusStore, logger zerolog.Logger) (*models.Users, error) {
	userObj := &models.Users{}
	if uo == nil {
		logger.Info().Msgf("Creating user object for '%v'", lu.Email)
		userObj.Email = lu.Email
		if userObj.Status == "" {
			userObj.Status = mpb.CommonStatus_ACTIVE.String()
		}
		userObj.FirstName = lu.FirstName
		userObj.DisplayName = lu.DisplayName
		userObj.LastName = lu.LastName
		userObj.InvitedType = mpb.UserInviteType_INVITE_ACCESS.String()
		userObj.InvitedOn = dmutils.GetEpochTime()
		userObj.OwnerAccount = ctxClaim[encryption.ClaimAccountKey]
		if lu.Organization != "" {
			userObj.Roles = []*models.UserOrganizationRole{
				{
					OrganizationId: lu.Organization,
					RoleArns:       lu.OrganizationRoles,
				},
			}
		}
		userObj.Profile = &models.UserProfile{}
		userObj.SetUserId()
		userObj.PreSaveCreate(ctxClaim)
		logger.Info().Msgf("New User object for '%v'", userObj)
	} else {
		userObj = uo
		logger.Info().Msgf("User object for '%v'", userObj)
	}
	userObj.SetAccountId(ctxClaim[encryption.ClaimAccountKey])
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(userObj).ModelTableExpr(apppkgs.UsersTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving datamarketplace in datastore")
		return nil, err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   userObj.VapusID,
			ResourceName: "USER",
			VapusBase: models.VapusBase{
				Editors: []string{userObj.UserId},
			},
		}, logger, ctxClaim)
	}()
	return userObj, nil
}
