package hcvault

type Vault struct {
	// VaultConfig is the configuration for the Hashicorp Vault
	URL, Token, Path, ApproleRoleID, ApproleSecretID, SecretEngine string
	AuthAppRole                                                    bool
}
