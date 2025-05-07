package azure

import "errors"

var (
	ErrCreatingAzureKeyVaultClient = errors.New("error while creating Azure Key Vault client")
	ErrCreatingAzureCredential     = errors.New("error while creating Azure credentials")
	ErrReadingAzureSecret          = errors.New("error while reading Azure Key Vault")
	ErrDeletingAzureSecret         = errors.New("error while deleting Azure Key Vault")
	ErrCreatingAzureSecret         = errors.New("error while creating secret in Azure Key Vault")
)
