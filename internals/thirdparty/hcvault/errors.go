package hcvault

import "errors"

var (
	//Vault Error Constants
	ErrVaultConnection       = errors.New("error while connecting to vault")
	ErrVaultWrite            = errors.New("error while writting secrets in vault")
	ErrVaultRead             = errors.New("error while reading secrets in vault")
	ErrVaultDelete           = errors.New("error while deleting secrets in vault")
	ErrVaultAppRoleAuth      = errors.New("error while creating AppRole Auth")
	ErrVaultAppRoleAuthLogin = errors.New("error while logging in with AppRole Auth")
	ErrInvalidVaultAppRole   = errors.New("invalid AppRole config for vault login")
)
