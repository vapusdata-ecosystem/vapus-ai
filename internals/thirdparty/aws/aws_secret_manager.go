package vaws

import (
	"context"
	json "encoding/json"

	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"

	awsg "github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// TODO: Add aws smithy error checks

// AWSSecretManager is an interface that contains the methods for AWS Secret Manager.
type AWSecretManager interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
	UpdateSecret(ctx context.Context, credData any, secretId string) error
}

type AwsSmStore struct {
	Client                                   *sm.Client
	secretNameTemplate, prefix, kmsSecretKey string
}

// NewAwsSmClient creates a new instance of AwsSmStore using the provided AWSConfig options.
// It returns a pointer to AwsSmStore and an error, if any.
func NewAwsSmClient(ctx context.Context, opts *AWSConfig) (*AwsSmStore, error) {
	configCl, err := opts.getAwsCLientConfig(ctx)
	if err != nil {
		return nil, dmerrors.DMError(ErrAwsConfigLoading, nil)
	}

	return &AwsSmStore{
		Client:       sm.NewFromConfig(configCl),
		kmsSecretKey: opts.KMSKey,
	}, nil
}

// WriteSecret writes the secret data to AWS Secret Manager.
// It takes a context, the credential data as a map, and the name of the secret as input.
// It returns an error if there was a problem writing the secret.
func (dmaws *AwsSmStore) WriteSecret(ctx context.Context, credData any, name string) error {
	cdBytes, err := json.Marshal(credData)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonUnMarshel, err)
	}

	inp := &sm.CreateSecretInput{
		Name:         awsg.String(name),
		SecretBinary: cdBytes,
	}
	if dmaws.kmsSecretKey != "" {
		inp.KmsKeyId = awsg.String(dmaws.kmsSecretKey)

	}
	if _, err = dmaws.Client.CreateSecret(ctx, inp); err != nil {
		return dmerrors.DMError(ErrAwsSecretWrite, err)
	}

	return nil
}

// ReadSecret reads the secret value from AWS Secrets Manager using the provided secret ID.
// It returns the secret value as a JSON string and unmarshals it into the provided 'val' parameter.
// If the secret is not found or an error occurs during the retrieval, it returns an error.
func (dmaws *AwsSmStore) ReadSecret(ctx context.Context, secretId string) (any, error) {
	secretValue, err := dmaws.Client.GetSecretValue(ctx, &sm.GetSecretValueInput{
		SecretId: awsg.String(secretId),
	})
	if err != nil {
		if err.Error() == "ResourceNotFoundException: Secrets Manager can't find the specified secret." {
			return nil, dmerrors.DMError(ErrAwsSecretRead, ErrAwsSecretRead)
		}
		return nil, dmerrors.DMError(ErrAwsSecretRead, err)
	}

	return secretValue.SecretBinary, nil
}

func (dmaws *AwsSmStore) DeleteSecret(ctx context.Context, secretId string) error {
	if _, err := dmaws.Client.DeleteSecret(ctx, &sm.DeleteSecretInput{
		SecretId: awsg.String(secretId),
	}); err != nil {
		return dmerrors.DMError(ErrAwsSecretDelete, err)
	}
	return nil
}

func (dmaws *AwsSmStore) UpdateSecret(ctx context.Context, credData any, secretId string) error {
	cdBytes, err := json.Marshal(credData)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonUnMarshel, err)
	}
	inp := &sm.UpdateSecretInput{
		SecretId:     awsg.String(secretId),
		SecretBinary: cdBytes,
	}
	if dmaws.kmsSecretKey != "" {
		inp.KmsKeyId = awsg.String(dmaws.kmsSecretKey)
	}
	if _, err := dmaws.Client.UpdateSecret(ctx, inp); err != nil {
		return dmerrors.DMError(ErrAwsSecretDelete, err)
	}
	return nil
}

func (dmaws *AwsSmStore) Close() {
	if dmaws.Client != nil {
		dmaws.Client = nil
	}
}
