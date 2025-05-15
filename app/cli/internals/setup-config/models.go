package setupconfig

import (
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
)

type VapusInstallerConfig struct {
	App struct {
		Name         string `yaml:"name"`
		Namespace    string `yaml:"namespace"`
		Organization string `yaml:"domain"`
		Address      string `yaml:"address"`
		Dev          bool   `yaml:"dev"`
	} `yaml:"app"`
	AccountBootstrap struct {
		PlatformOwners  []string `yaml:"platformOwners"`
		PlatformAccount struct {
			Name    string `yaml:"name"`
			Creator string `yaml:"creator"`
		} `yaml:"platformAccount"`
		PlatformAccountOrganization struct {
			Name string `yaml:"name"`
		} `yaml:"platformAccountOrganization"`
	} `yaml:"accountBootstrap"`
	Secrets    *VapusSecretsMap `yaml:"secrets"`
	Postgresql struct {
		FullnameOverride string `yaml:"fullnameOverride"`
		Auth             struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"auth"`
	} `yaml:"postgresql"`
	Redis struct {
		FullnameOverride string `yaml:"fullnameOverride"`
		Auth             struct {
			Password string `yaml:"password"`
		} `yaml:"auth"`
	} `yaml:"postgresql"`
	Vault   *Vault `yaml:"vault"`
	TLSCert struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tlsCert"`
	SecretStore    *models.DataSourceCredsParams `yaml:"secretStore"`
	DevSecretStore *models.DataSourceCredsParams `yaml:"devSecretStore"`
}

type Vault struct {
	FullnameOverride string `yaml:"fullnameOverride"`
	Server           struct {
		Standalone struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"standalone"`
	} `yaml:"server"`
}

type VapusSecretInstallerConfig struct {
	SecretStore       *models.DataSourceCredsParams `yaml:"secretStore"`
	DevSecretStore    *models.DataSourceCredsParams `yaml:"devSecretStore"`
	BackendDataStore  *models.DataSourceCredsParams `yaml:"backendDataStore"`
	BackendCacheStore *models.DataSourceCredsParams `yaml:"backendCacheStore"`
	JWTAuthnSecrets   *encryption.JWTAuthn          `yaml:"JWTAuthnSecrets"`
	FileStore         *models.DataSourceCredsParams `yaml:"fileStore"`
	AuthnSecrets      *authn.AuthnSecrets           `yaml:"authnSecrets"`
	CreateDatabase    bool                          `yaml:"createDatabase"`
}

type VapusSecretsMap struct {
	BackendSecretStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"backendSecretStore"`
	BackendDataStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"backendDataStore"`
	BackendCacheStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"backendCacheStore"`
	JWTAuthnSecrets struct {
		Secret string `yaml:"secret"`
	} `yaml:"JWTAuthnSecrets"`
	AuthnSecrets struct {
		Secret string `yaml:"secret"`
	} `yaml:"authnSecrets"`
	FileStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"fileStore"`
}
