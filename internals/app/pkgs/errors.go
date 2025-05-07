package apppkgs

import "errors"

var (
	ErrAuthenticatorInitFailed  = errors.New("error while initializing authenticator for the service based on the provided configuration")
	ErrAuthenticatorParamsNil   = errors.New("error while initializing authenticator for the service, authn params are nil")
	ErrJwtParamsNil             = errors.New("error while initializing jwt authn for the service, jwt params are nil")
	ErrJwtAuthInitFailed        = errors.New("error while initializing jwt authn for the service based on the provided configuration")
	ErrValidatorInitFailed      = errors.New("error while initializing validator for the service")
	ErrPbacConfigPathEmpty      = errors.New("error while initializing pbac config for the service, pbac config path is empty")
	ErrPbacConfigInitFailed     = errors.New("error while initializing pbac config for the service based on the provided configuration")
	ErrDataSourceCredsSecretGet = errors.New("error while reading the secret from the vault")
	ErrDataSourceCredsNotFound  = errors.New("error while getting the credentials for the node")
)
