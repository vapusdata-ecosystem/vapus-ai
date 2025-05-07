package gcp

import "errors"

var (
	ErrCreatingGcpArClient      = errors.New("error creating gcp artifact registry client")
	ErrCreatingGcpSMClient      = errors.New("error while creating GCP Secret Manager client")
	ErrReadingGcpSecret         = errors.New("error while reading GCP Secret Manager")
	ErrDeletingGcpSecret        = errors.New("error while deleting GCP Secret Manager")
	ErrCreatingGcpSecret        = errors.New("error while creating secret in GCP Secret Manager")
	ErrParsingGAR               = errors.New("error while parsing GCP Artifact Registry URL")
	ErrParsingGARRegion         = errors.New("error while parsing GCP Artifact Registry Region")
	ErrParsingGARHost           = errors.New("error while parsing GCP Artifact Registry Host")
	ErrListingGARPackages       = errors.New("error while listing GCP Artifact Registry packages")
	ErrCreatingGcpSecretVersion = errors.New("error while creating secret version in GCP Secret Manager")
	ErrCreatingBucketClient     = errors.New("error while creating GCP bucket client")
	ErrDeletingBucket           = errors.New("error while deleting GCP bucket")
	ErrCreatingBucket           = errors.New("error while creating GCP bucket")
	ErrDownloadingObject        = errors.New("error while downloading object from GCP bucket")
	ErrDeletingObject           = errors.New("error while deleting object from GCP bucket")
	ErrUploadingObject          = errors.New("error while uploading object to GCP bucket")
	ErrGettingObject            = errors.New("error while getting object from GCP bucket")
	ErrGetBucketAttrs           = errors.New("error while getting bucket attributes")
	ErrGetBucket                = errors.New("error while getting bucket")
	ErrUpdatingGcpSecret        = errors.New("error while updating secret in GCP Secret Manager")
)
