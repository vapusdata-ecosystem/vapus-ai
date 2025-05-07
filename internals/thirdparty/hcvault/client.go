package hcvault

import (
	"context"
	"encoding/json"
	"fmt"

	vaultApi "github.com/hashicorp/vault/api"
	approle "github.com/hashicorp/vault/api/auth/approle"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
)

type HCVault interface {
	WriteSecret(ctx context.Context, data map[string]any, name string) error
	ReadSecret(ctx context.Context, secretId string) (any, error)
	DeleteSecret(ctx context.Context, secretId string) error
	UpdateSecret(ctx context.Context, data any, name string) error
}

// VaultStore struct to store vault client
type HcVaultStore struct {
	Client               *vaultApi.Client
	SecretEngine, prefix string
}

// Connect to vault server and return the client  object
func NewHcVaultManager(ctx context.Context, conf *Vault) (*HcVaultStore, error) {
	config := vaultApi.DefaultConfig()

	config.Address = conf.URL

	client, err := vaultApi.NewClient(config)
	if err != nil {
		return nil, dmerrors.DMError(ErrVaultConnection, err)
	}

	// Authenticate with AppRole if flag is enabled
	if conf.AuthAppRole {
		token, err := loginWithAppRole(ctx, conf, client)
		if err != nil {
			return nil, err
		}
		client.SetToken(token)
	} else {
		// Authenticate with token
		client.SetToken(conf.Token)

	}

	return &HcVaultStore{
		Client:       client,
		SecretEngine: conf.SecretEngine,
	}, nil
}

// loginWithAppRole logs in to vault using AppRole and returns the token
func loginWithAppRole(ctx context.Context, conf *Vault, client *vaultApi.Client) (string, error) {
	// If its not file then no n eed to read from file
	approleSecretID := &approle.SecretID{
		FromFile: conf.ApproleSecretID,
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		conf.ApproleRoleID,
		approleSecretID,
		approle.WithWrappingToken(), // Needed token wrapped for approle login
	)
	if err != nil {
		return "", dmerrors.DMError(ErrVaultAppRoleAuth, err)
	}

	authInfo, err := client.Auth().Login(ctx, appRoleAuth)

	if err != nil {
		return "", dmerrors.DMError(ErrVaultAppRoleAuthLogin, err)
	}

	if authInfo == nil {
		return "", dmerrors.DMError(ErrInvalidVaultAppRole, fmt.Errorf("no auth info found"))
	}

	return authInfo.Auth.ClientToken, nil
}

// StoreKVCredentials stores the credentials in the vault server at the given path and returns error
func (vlt *HcVaultStore) WriteSecret(ctx context.Context, data any, name string) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(ErrVaultWrite, err)
	}
	obj := make(map[string]any)
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		return dmerrors.DMError(ErrVaultWrite, err)
	}

	_, err = vlt.Client.KVv2(vlt.SecretEngine).Put(ctx, name, obj)

	if err != nil {
		return dmerrors.DMError(ErrVaultWrite, err)
	}
	return nil
}

func (vlt *HcVaultStore) ReadSecret(ctx context.Context, secretId string) (any, error) {
	secret, err := vlt.Client.KVv2(vlt.SecretEngine).Get(ctx, secretId)
	if err != nil {
		return nil, dmerrors.DMError(ErrVaultRead, err)
	}
	return json.Marshal(secret.Data)
}

func (vlt *HcVaultStore) DeleteSecret(ctx context.Context, secretId string) error {
	err := vlt.Client.KVv2(vlt.SecretEngine).Delete(ctx, secretId)
	if err != nil {
		return dmerrors.DMError(ErrVaultDelete, err)
	}
	return nil
}

func (vlt *HcVaultStore) UpdateSecret(ctx context.Context, data any, name string) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return dmerrors.DMError(ErrVaultWrite, err)
	}
	obj := make(map[string]any)
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		return dmerrors.DMError(ErrVaultWrite, err)
	}

	_, err = vlt.Client.KVv2(vlt.SecretEngine).Put(ctx, name, obj)

	if err != nil {
		return dmerrors.DMError(ErrVaultWrite, err)
	}
	return nil
}

func (vlt *HcVaultStore) Close() {
	if vlt.Client != nil {
		vlt.Client = nil
	}
}
