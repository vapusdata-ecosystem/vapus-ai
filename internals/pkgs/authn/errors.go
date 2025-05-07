package authn

import "errors"

var (
	ErrOIDCProviderInit = errors.New("error while initializing OIDC provider")
	ErrOIDCInvalidToken = errors.New("no id_token field in oauth2 token")
	ErrUnAuthenticated  = errors.New("user is not authenticated")
	ErrTokenExpired     = errors.New("token expired")

	ErrAuthenticatorInitFailed      = errors.New("authenticator initialization failed")
	ErrInvalidLogin                 = errors.New("invalid login url")
	ErrGeneratingState              = errors.New("error generating state for login")
	ErrTokenExchangeFailed          = errors.New("error exchanging token with OIDC provider based on recieved code")
	ErrIDTokenVerificationFailed    = errors.New("error verifying ID token from recieved token")
	ErrIDTokenClaimFailed           = errors.New("error extracting claims from ID token")
	ErrUnauthorized                 = errors.New("unauthorized access")
	ErrAccessDMTokenCreationFailed  = errors.New("error creating ORGANIZATION jwt access token")
	ErrInvalidUserAuthenticationSub = errors.New("invalid user authentication sub")
	ErrClaimFirstNameNotFound       = errors.New("error: first name not found in authenticated idtoken")
	ErrClaimEmailNotFound           = errors.New("error: email not found in authenticated idtoken")
	ErrClaimDisplayNameNotFound     = errors.New("error: displayName not found in authenticated idtoken")
	ErrClaimLastNameNotFound        = errors.New("error: last name not found in authenticated idtoken")
)
