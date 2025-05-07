package encryption

import (
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
)

// type VapusDataPlatformAccessClaims struct {
// 	Scope   string                 `json:"scope"`
// 	Profile map[string]any `json:"profile"`
// 	jwt.RegisteredClaims
// }

type VapusDataResources struct {
	Name        string   `json:"name"`
	Identifiers []string `json:"identifiers"`
	Role        []string `json:"role"`
}

// TODO : add validators for below 2 structs

type PlatformScope struct {
	UserId           string `validate:"required" json:"userId"`
	AccountId        string `validate:"required" json:"accountId"`
	OrganizationId   string `validate:"required" json:"OrganizationId"`
	OrganizationRole string `validate:"required" json:"OrganizationRole"`
	RoleScope        string `validate:"required" json:"roleScope"`
	PlatformRole     string `validate:"required" json:"platformRole"`
	// MarketplaceId 	 string `validate:"required" json:"marketplaceId"`
}

type VapusDataPlatformAccessClaims struct {
	jwt.RegisteredClaims
	Scope *PlatformScope `validate:"required" json:"scope"`
}

func (x *PlatformScope) Validate() error {
	validator := validator.New()
	return validator.Struct(x)
}

type VapusDataPlatformRefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId         string `validate:"required" json:"userId"`
	OrganizationId string `validate:"required" json:"OrganizationId"`
}

func (x *VapusDataPlatformRefreshTokenClaims) Validate() error {
	validator := validator.New()
	return validator.Struct(x)
}
