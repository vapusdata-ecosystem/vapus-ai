package pkgs

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
)

func BuildVDPAClaim(userObj *models.Users, organizationId string, validTill time.Time) (*encryption.VapusDataPlatformAccessClaims, error) {
	var roleScope string
	var organizationRoles string
	log.Println(organizationId, userObj.UserId, userObj.OwnerAccount, validTill, "==========================----------------")
	if organizationId != "" {
		roleScope = encryption.JwtOrganizationScope
		log.Println("organizationId", organizationId, "userObj.UserId", userObj.UserId, "userObj.OwnerAccount", userObj.OwnerAccount, "validTill", validTill)
		r := userObj.GetOrganizationRole(organizationId)
		if len(r) < 1 && organizationId != "" {
			return nil, dmerrors.ErrUserORGANIZATION404
		} else {
			for _, val := range r[0].RoleArns {
				if organizationRoles == "" {
					organizationRoles = val
				} else {
					organizationRoles = organizationRoles + encryption.JwtClaimRoleSeparator + val
				}
			}
		}
	} else {
		roleScope = encryption.JwtPlatformScope
	}
	log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++")
	log.Println("roleScope", roleScope, "organizationRoles", organizationRoles, "userObj.UserId", userObj.UserId, "organizationId", organizationId, "userObj.OwnerAccount", userObj.OwnerAccount, "validTill", validTill)
	return &encryption.VapusDataPlatformAccessClaims{
		Scope: &encryption.PlatformScope{
			UserId:           userObj.UserId,
			OrganizationId:   organizationId,
			RoleScope:        roleScope,
			AccountId:        userObj.OwnerAccount,
			OrganizationRole: organizationRoles,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   encryption.VapusPlatformTokenSubject,
			Audience:  []string{NetworkConfigManager.ExternalURL},
			ExpiresAt: jwt.NewNumericDate(validTill), // configurable tokens
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 1)),
			Issuer:    NetworkConfigManager.ExternalURL,
		},
	}, nil
}

func BuildVDPRTClaim(userObj *models.Users, organizationId string, validTill time.Time) (*encryption.VapusDataPlatformRefreshTokenClaims, error) {
	return &encryption.VapusDataPlatformRefreshTokenClaims{
		UserId:         userObj.UserId,
		OrganizationId: organizationId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   encryption.VapusPlatformTokenSubject,
			Audience:  []string{NetworkConfigManager.ExternalURL},
			ExpiresAt: jwt.NewNumericDate(validTill), // configurable tokens
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 1)),
			Issuer:    NetworkConfigManager.ExternalURL,
		},
	}, nil
}
