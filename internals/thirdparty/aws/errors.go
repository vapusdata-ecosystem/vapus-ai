package vaws

import "errors"

var (

	// Error constants for AWS Config
	ErrInvalidAwsConfig = errors.New("invalid AWS config")
	ErrAwsConfigLoading = errors.New("error while loading AWS config")

	// Error constants for AWS Secret Manager
	ErrAwsSecretWrite  = errors.New("error while writting credentials in the aws secret manager")
	ErrAwsSecretRead   = errors.New("error while writting credentials in the aws secret manager")
	ErrAwsSecretDelete = errors.New("error while writting credentials in the aws secret manager")
	ErrAwsSecret404    = errors.New("aws Secrets Manager can't find the specified secret")

	//ECR Manager Errors
	ErrListingECRPackages  = errors.New("error while listing ECR packages")
	ErrParsingEcrAddress   = errors.New("error while parsing the ECR address, its not a valid URL - registries must be valid RFC 3986 URI authorities")
	ErrGettingEcrAuthToken = errors.New("error while getting the ECR auth token, please check the credentials provided")
	ErrECRNoAuthTokenFound = errors.New("no auth token found for the ECR repository")
	ErrParsingEcrId        = errors.New("error while getting ECR Id")
)
