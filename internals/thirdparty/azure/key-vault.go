package azure

import (
	"context"
	"encoding/json"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type AzureKVManager interface {
	WriteSecret(ctx context.Context, data any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
	UpdateSecret(ctx context.Context, data any, secretName string) error
}

type AzureKeyVault struct {
	client                           *azsecrets.Client
	secretPrefix, secretNameTemplate string
}

func NewAzureKeyVault(ctx context.Context, opts *AzureConfig, valultURI string) (*AzureKeyVault, error) {
	credential, err := azidentity.NewClientSecretCredential(opts.TenantID, opts.ClientID, opts.ClientSecret, nil)
	if err != nil {
		return nil, dmerrors.DMError(ErrCreatingAzureCredential, err)
	}

	// Establish a connection to the Key Vault client
	client, err := azsecrets.NewClient(valultURI, credential, nil)
	if err != nil {
		return nil, dmerrors.DMError(ErrCreatingAzureKeyVaultClient, err)
	}
	return &AzureKeyVault{
		client: client,
	}, nil
}

func (akv *AzureKeyVault) WriteSecret(ctx context.Context, data any, secretName string) error {
	// Convert the secret value to a byte array
	secretValue, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonMarshel, err)
	}
	_, err = akv.client.SetSecret(ctx, secretName, azsecrets.SetSecretParameters{
		Value: dmutils.Str2Ptr(string(secretValue)),
	}, nil)
	if err != nil {
		return dmerrors.DMError(ErrCreatingAzureSecret, err)
	}

	return nil
}

func (akv *AzureKeyVault) ReadSecret(ctx context.Context, secretName string) (any, error) {
	resp, err := akv.client.GetSecret(ctx, secretName, "", nil)
	// TO:DO check error for 404 or other using error.As
	if err != nil {
		return nil, dmerrors.DMError(ErrReadingAzureSecret, err)
	}

	return json.Marshal([]byte(*resp.Value))
}

func (akv *AzureKeyVault) DeleteSecret(ctx context.Context, secretName string) error {
	_, err := akv.client.DeleteSecret(ctx, secretName, nil)
	if err != nil {
		return dmerrors.DMError(ErrDeletingAzureSecret, err)
	}
	return nil
}

func (akv *AzureKeyVault) UpdateSecret(ctx context.Context, data any, secretName string) error {
	// Convert the secret value to a byte array
	secretValue, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(dmerrors.ErrJsonMarshel, err)
	}
	_, err = akv.client.SetSecret(ctx, secretName, azsecrets.SetSecretParameters{
		Value: dmutils.Str2Ptr(string(secretValue)),
	}, nil)
	if err != nil {
		return dmerrors.DMError(ErrCreatingAzureSecret, err)
	}

	return nil
}

func (akv *AzureKeyVault) Close() {
	if akv.client != nil {
		akv.client = nil
	}
}
