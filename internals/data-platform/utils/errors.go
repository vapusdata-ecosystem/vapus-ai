package utils

import (
	"errors"
)

var (
	ErrNoCredentialFoundForDataSource = errors.New("no credential found for data node")
	ErrNoDataSourceFound              = errors.New("no data node found")
	ErrListingArtifacts               = errors.New("error while listing artifacts")
	ErrNoNameFoundForPackage          = errors.New("no name found for package")
	ErrDataSourceAttributesNotFound   = errors.New("error: no attributes found in data node")
	ErrDecodingCredential             = errors.New("error while decoding node credential")
	ErrNoPackagesFoundInECR           = errors.New("no packages found in ecr")
	ErrConnDataSource                 = errors.New("error while connecting to data node")

	ErrNoMappings = errors.New("no mappings found in ElasticSearch for this index")
	ErrNoProps    = errors.New("no properties found in ElasticSearch for this index")

	ErrInvalidDataSourceType = errors.New("invalid data source")

	ErrTokenExpired    = errors.New("token expired")
	ErrUnAuthenticated = errors.New("user is not authenticated")
)

var (
	ErrInvalidDataProductFormat = errors.New("invalid format for data product spec")
	ErrInvalidNabhikWorkflow    = errors.New("invalid workflow requested from nabhik engine")

	ErrInvalidNameSpaceOperation = errors.New("invalid action for namespace operations")
	ErrProductDeploymentFailed   = errors.New("data product deployment failed in k8s")

	ErrInvalidNabhikAgent     = errors.New("invalid configuration for nabhik agent")
	ErrInvalidNabhikOperation = errors.New("invalid action for nabhik agent")
	ErrNoDPContainers         = errors.New("error: no data workers are configured in this data product")

	ErrMissingORGANIZATIONArtifactStore = errors.New("no artifact stores are present for current ORGANIZATION/platform")
	ErrImagePullSecretOperation         = errors.New("error while fetching or creating ORGANIZATION pull secrets")

	ErrInvalidDataProductPublisherAgent = errors.New("invalid configuration for data product publisher agent")
	ErrORGANIZATIONJwtSecretFailed      = errors.New("error while fetching/creating ORGANIZATION jwt secret")
	ErrInvalidDataProductAgentOperation = errors.New("invalid action for data product agent")
	ErrBuildingDataProduct              = errors.New("error while building data product")

	ErrInvalidDataWorkerAgentOperation = errors.New("invalid action for data worker agent")
	ErrDataWorkerConfig404             = errors.New("data worker config not found")
)

var (
	ErrInvalidDataStorageEngine    = errors.New("invalid data storage engine")
	ErrDataStoreConn               = errors.New("error creating data store connection")
	ErrScanDestinationPtr          = errors.New("destination should be a pointer")
	ErrScanDestinationNil          = errors.New("destination should not be nil for scanning the result")
	ErrInvalidDestinationType      = errors.New("invalid destination type")
	ErrInvalidDataStoreEngine      = errors.New("invalid data store engine")
	ErrInvalidSecretData           = errors.New("invalid secret data")
	ErrDataStoreParams404          = errors.New("data store params not found")
	ErrInvalidBlobStoreSvc         = errors.New("invalid blob store service, this service is not supported yet")
	ErrInvalidEmailStoreSvc        = errors.New("invalid email store service, this service is not supported yet")
	ErrInvalidDataStoreCredentials = errors.New("invalid data store credentials, please check the credentials provided")
	ErrInvalidArtifactStoreSvc     = errors.New("invalid artifact store service, this service is not supported yet")
)
