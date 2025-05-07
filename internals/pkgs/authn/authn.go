package authn

import (
	"strings"
	"time"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
)

type Authenticator struct {
	*OIDCAuthenticator
	AuthnMethod string
}

func NewWithConfig(path string) (*Authenticator, error) {
	authnConfig, err := LoadAuthnSecrets(path)
	if err != nil {
		return nil, err
	}
	return New(authnConfig.OIDCSecrets, authnConfig.AuthnMethod)
}

func New(oidcOpts *OIDCSecrets, authnMethod string) (*Authenticator, error) {
	switch authnMethod = strings.TrimSpace(authnMethod); authnMethod {
	case mpb.AuthnMethod_OIDC.String():
		oidcAuthn, err := NewOIDCAuthenticator(oidcOpts)
		if err != nil {
			return nil, err

		}
		return &Authenticator{
			OIDCAuthenticator: oidcAuthn,
			AuthnMethod:       mpb.AuthnMethod_OIDC.String(),
		}, nil
	default:
		oidcAuthn, err := NewOIDCAuthenticator(oidcOpts)
		if err != nil {
			return nil, err

		}
		return &Authenticator{
			OIDCAuthenticator: oidcAuthn,
			AuthnMethod:       mpb.AuthnMethod_OIDC.String(),
		}, nil
	}
}

func ValidateOIDCAuth(token string, logger zerolog.Logger) (map[string]any, error) {
	// Use OIDC JWT token to validate the user
	claims, err := encryption.ParseUnValidatedJWT(token)
	logger.Info().Msgf("Claims -- %v", claims)
	if err != nil {
		logger.Err(err).Msg("error while parsing unvalidated claim from OIDC provider")
		return nil, err
	}

	expVal, ok := claims["exp"]
	if !ok {
		logger.Err(err).Msg("invalid token, user is not authenticated")
		return nil, ErrUnAuthenticated
	}

	logger.Info().Msgf("exp -- %v", expVal.(float64))

	if int64(expVal.(float64)) < time.Now().Unix() {
		logger.Err(err).Msg("invalid token, authentication token is expired")
		return nil, ErrTokenExpired
	}

	return claims, nil
}
