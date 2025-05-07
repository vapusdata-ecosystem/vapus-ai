package services

import (
	"time"

	guuid "github.com/google/uuid"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	apppdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	authn "github.com/vapusdata-ecosystem/vapusai/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
)

func claimToLocal(claims map[string]interface{}) (*apppdrepo.LocalUserM, error) {
	email, ok := claims["email"]
	if !ok {
		helperLogger.Error().Msgf("error: email not found in authenticated idtoken")
		return nil, encryption.ErrInvalidUserAuthentication
	}
	authSubject, err := authn.GetAuthnConnectionSub(claims, helperLogger)
	if err != nil {
		helperLogger.Err(err).Msg(err.Error())
		return nil, err
	}
	displayName, err := authSubject.GetDisplayName(claims)
	if err != nil {
		helperLogger.Error().Msgf("error: displayName not found in authenticated idtoken")
		return nil, encryption.ErrInvalidUserAuthentication
	}
	fName, err := authSubject.GetFirstName(claims)
	if err != nil {
		helperLogger.Error().Msgf("error: first name not found in authenticated idtoken")
		return nil, encryption.ErrInvalidUserAuthentication
	}
	lName, err := authSubject.GetLastName(claims)
	if err != nil {
		helperLogger.Error().Msgf("error: last name not found in authenticated idtoken")
		return nil, encryption.ErrInvalidUserAuthentication
	}
	picture, err := authSubject.GetProfilePic(claims)
	if err != nil {
		helperLogger.Error().Msgf("error: picture not found in authenticated idtoken")
		return nil, encryption.ErrInvalidUserAuthentication
	}
	return &apppdrepo.LocalUserM{
		Email:        email.(string),
		DisplayName:  displayName,
		FirstName:    fName,
		LastName:     lName,
		ProfileImage: picture,
	}, nil
}

func generateOrganizationAccessToken(user *models.Users, organization string, validTill time.Time) (string, string, error) {
	vapusClaim, err := pkgs.BuildVDPAClaim(user, organization, validTill)
	if err != nil {
		helperLogger.Err(err).Msg(err.Error())
		return "", "0", err
	}
	tokenId := guuid.New().String() // make it more unique
	vapusClaim.RegisteredClaims.ID = tokenId
	token, err := pkgs.SvcPackageManager.VapusJwtAuth.GenerateVDPAJWT(vapusClaim)
	if err != nil {
		helperLogger.Err(err).Msg("error while generating VDPA jwt token")
		return "", "", err
	}
	return tokenId, token, err
}

func generateRefreshToken(user *models.Users, organization string, validTill time.Time, isOwner bool) (string, string, error) {
	vapusClaim, err := pkgs.BuildVDPRTClaim(user, organization, validTill)
	if err != nil {
		helperLogger.Err(err).Msg(err.Error())
		return "", "0", err
	}
	tokenId := guuid.New().String() // make it more unique
	vapusClaim.RegisteredClaims.ID = tokenId
	// vapusClaim.RegisteredClaims.Is = isOwner
	token, err := pkgs.SvcPackageManager.VapusJwtAuth.GenerateVDPARefreshJWT(vapusClaim)
	if err != nil {
		helperLogger.Err(err).Msg("error while generating VDPA refresh token")
		return "", "", err
	}
	return tokenId, token, err
}
