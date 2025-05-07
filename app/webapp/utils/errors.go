package utils

import "errors"

var (
	ErrAuthenticatorInitFailed     = errors.New("authenticator initialization failed")
	ErrInvalidLogin                = errors.New("invalid login url")
	ErrGeneratingState             = errors.New("error generating state for login")
	ErrTokenExchangeFailed         = errors.New("error exchanging token with OIDC provider based on recieved code")
	ErrIDTokenVerificationFailed   = errors.New("error verifying ID token from recieved token")
	ErrIDTokenClaimFailed          = errors.New("error extracting claims from ID token")
	ErrUnauthorized                = errors.New("unauthorized access")
	ErrAccessDMTokenCreationFailed = errors.New("error creating Organization jwt access token")
	ErrInvalidNetworkConfig        = errors.New("invalid network configuration")
	ErrInvalidURLResourceName      = errors.New("invalid URL resource name")
	ErrInvalidURLResourceId        = errors.New("invalid URL resource id")
)
