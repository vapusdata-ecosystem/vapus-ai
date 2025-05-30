package aidmstore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apppdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

func (ds *AIStudioDMStore) GetOrUpdateUser(ctx context.Context, lu *apppdrepo.LocalUserM, createIfInvited bool, useDefaultOrganization bool, ctxClaim map[string]string) (*models.Users, error) {
	result := []*models.Users{}
	var err error
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = '%s'", apppkgs.UsersTable, lu.Email)
	err = ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Error().Msgf("User Not found with userEmail - %v, please get an invite.", lu.Email)
		return nil, dmerrors.DMError(apperr.ErrUser404, err)
	}
	userObj := result[0]
	if userObj.Status == mpb.CommonStatus_INVITED.String() && createIfInvited && userObj.InviteExpiresOn > time.Now().Unix() {
		validOrganization := false
		organizationId := ""
		if !useDefaultOrganization {
			for _, val := range userObj.Roles {
				if val.OrganizationId == lu.Organization {
					validOrganization = true
					organizationId = val.OrganizationId
				}
			}
		} else {
			if len(userObj.Roles) > 0 {
				validOrganization = true
				organizationId = userObj.Roles[0].OrganizationId
			}
		}

		if !validOrganization {
			logger.Error().Msgf("error: organization %v is not attached to user %v", lu.Organization, lu.Email)
			return nil, apperr.ErrUserOrganization404
		}
		userObj.Status = mpb.CommonStatus_ACTIVE.String()
		userObj.FirstName = lu.FirstName
		userObj.LastName = lu.LastName
		userObj.DisplayName = lu.DisplayName
		if userObj.Profile == nil {
			userObj.Profile = &models.UserProfile{}
		}
		userObj.Profile.Avatar = lu.ProfileImage
		userObj.SetUserId()
		userObj.PreSaveCreate(ctxClaim)
		userObj.SetDefaultOrganization(organizationId)
		err = ds.PutUser(ctx, userObj, ctxClaim)
		if err != nil {
			logger.Error().Msgf("error: user %v is not attached to any organization", lu.Email)
			return nil, err
		}
		return userObj, nil
	}
	if len(userObj.Roles) == 0 {
		logger.Error().Msgf("error: user %v is not attached to any organization", lu.Email)
		return nil, apperr.ErrUserOrganization404
	}
	if userObj.Status == mpb.CommonStatus_ACTIVE.String() {
		return userObj, nil
	}
	return nil, apperr.ErrUser404
}

func (ds *AIStudioDMStore) CreateUser(ctx context.Context, lu *apppdrepo.LocalUserM, uo *models.Users, ctxClaim map[string]string) (*models.Users, error) {
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
		}, ds.logger, ctxClaim)
	}()
	return userObj, nil
}

func (ds *AIStudioDMStore) UserInviteExists(ctx context.Context, userId string, ctxClaim map[string]string) bool {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = '%s'", apppkgs.UsersTable, userId)
	var user *models.Users
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, user)
	if err != nil || user == nil {
		logger.Info().Msgf("User Not found with userEmail - %v, please get an invite.", userId)
		return false
	}
	return true
}

func (ds *AIStudioDMStore) LogPlatformRTinfo(ctx context.Context, obj *models.RefreshTokenLog, ctxClaim map[string]string) error {
	obj.SetAccountId(ctxClaim[encryption.ClaimAccountKey])
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.RefreshTokenLogsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving refresh token log in datastore")
		return err
	}
	ctx.Done()
	return nil
}

func (ds *AIStudioDMStore) LogPlatformJwtinfo(ctx context.Context, obj *models.JwtLog, ctxClaim map[string]string) error {
	obj.SetAccountId(ctxClaim[encryption.ClaimAccountKey])
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.JwtLogsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving jwt log in datastore")
		return err
	}
	ctx.Done()
	return nil
}

func (ds *AIStudioDMStore) GetPlatformRTinfo(ctx context.Context, token string, ctxClaim map[string]string) (*models.RefreshTokenLog, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE token_hash = '%s'", apppkgs.RefreshTokenLogsTable, token)
	var user *models.RefreshTokenLog
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, user)
	if err != nil {
		logger.Info().Msg("error while getting refresh token log from datastore")
		return nil, err
	}
	return user, nil
}

func (ds *AIStudioDMStore) PatchUser(ctx context.Context, userId string, data, conditions map[string]interface{}, ctxClaim map[string]string) error {
	// Convert the script query to JSON
	pq := ds.Db.PostgresClient.DB.NewUpdate().Model(&data).ModelTableExpr(apppkgs.UsersTable)

	for key, value := range conditions {
		pq = pq.Where(key, value)
	}
	_, err := pq.Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while patching user info in datastore")
		return err
	}
	logger.Debug().Ctx(ctx).Msgf("User '%v' successfully updated", userId)
	return nil
}

func (ds *AIStudioDMStore) PutUser(ctx context.Context, obj *models.Users, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.UsersTable).Where("user_id = ?", obj.UserId).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating user in datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) GetOrganizationUsers(ctx context.Context, organization string, ctxClaim map[string]string) ([]*models.Users, error) {
	result := []*models.Users{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE organization_roles @> '[{"organizationId": "%s"}]'`, apppkgs.UsersTable, organization)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting users from datastore")
		return result, err
	}
	return result, nil
}

func (ds *AIStudioDMStore) ListUsers(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.Users, error) {
	result := make([]*models.Users, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.UsersTable, apppkgs.GetAccountFilter(ctxClaim, condition))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting users from datastore")
		return nil, err
	}
	return result, err
}

func (ds *AIStudioDMStore) CountUsers(ctx context.Context, condition string, ctxClaim map[string]string) (int64, error) {
	var result int64
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", apppkgs.UsersTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting users count from datastore")
		return 0, err
	}
	return result, err
}

func (ds *AIStudioDMStore) CustomListUsers(ctx context.Context, fieldQuery, condition, postFilterForamtting string, ctxClaim map[string]string) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s %s", fieldQuery, apppkgs.UsersTable, condition, postFilterForamtting)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting dataproducts from datastore")
		return nil, err
	}
	return result, err
}

func (ds *AIStudioDMStore) GetUser(ctx context.Context, userId string, ctxClaim map[string]string) (*models.Users, error) {
	result := make([]*models.Users, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.UsersTable, apppkgs.GetByIdFilter("user_id", userId, ctxClaim))
	log.Println("Query to get user:", query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting users from datastore")
		return nil, apperr.ErrUser404
	}
	return result[0], err
}
