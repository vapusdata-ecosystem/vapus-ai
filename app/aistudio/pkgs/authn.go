package pkgs

import (
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/authn"
)

type AuthnService struct {
	*authn.Authenticator
	Auth, Callback string
}

var AuthnManager *AuthnService

var AuthnParams *authn.AuthnSecrets

func InitAuthnManager(params *authn.AuthnSecrets) {
	if AuthnManager == nil {
		AuthnManager = &AuthnService{
			Authenticator: SvcPackageManager.AuthnManager,
		}
	}
}
