package setupconfig

import (
	"github.com/rs/zerolog"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
)

type VapusInstallerConfig struct {
	App struct {
		Name      string `yaml:"name" validate:"required"`
		Namespace string `yaml:"namespace" validate:"required"`
		Domain    string `yaml:"domain"`
		Address   string `yaml:"address"`
		Dev       bool   `yaml:"dev"`
	} `yaml:"app" validate:"required"`
	AccountBootstrap *AppBootConfig   `yaml:"accountBootstrap" validate:"required"`
	Secrets          *VapusSecretsMap `yaml:"secrets"`
	Postgresql       struct {
		FullnameOverride string `yaml:"fullnameOverride"`
		Auth             struct {
			Username string `yaml:"username" validate:"required"`
			Password string `yaml:"password" validate:"required"`
			Database string `yaml:"database" validate:"required"`
		} `yaml:"auth" validate:"required"`
	} `yaml:"postgresql"`
	Redis struct {
		FullnameOverride string `yaml:"fullnameOverride"`
		Auth             struct {
			Enabled  bool   `yaml:"enabled"`
			Password string `yaml:"password" validate:"required"`
		} `yaml:"auth" validate:"required"`
		Master struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"master"`
		Sentinal struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"sentinal"`
	} `yaml:"redis"`
	Vault          *Vault                        `yaml:"vault" validate:"required"`
	TLSCert        *TLSCert                      `yaml:"tlsCert" validate:"required"`
	SecretStore    *models.DataSourceCredsParams `yaml:"secretStore"`
	DevSecretStore *models.DataSourceCredsParams `yaml:"devSecretStore"`
}

func (x *VapusInstallerConfig) Validate(logger zerolog.Logger) error {
	validator := GetValidator()
	validator.Struct(x)
	err := validator.Struct(x)
	if err != nil {
		HandleValiationError(err, logger)
		return err
	}
	return nil
}

type TLSCert struct {
	Cert         string `yaml:"cert"`
	Key          string `yaml:"key"`
	CertFile     string `yaml:"certFile"`
	KeyFile      string `yaml:"keyFile"`
	AutoGenerate bool   `yaml:"autoGenerate"`
}
type Vault struct {
	FullnameOverride string `yaml:"fullnameOverride"`
	Server           struct {
		Standalone struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"standalone"`
	} `yaml:"server"`
}

type AppBootConfig struct {
	PlatformOwners  []string `yaml:"platformOwners"`
	PlatformAccount struct {
		Name    string `yaml:"name"`
		Creator string `yaml:"creator"`
	} `yaml:"platformAccount"`
	PlatformAccountOrganization struct {
		Name string `yaml:"name"`
	} `yaml:"platformAccountOrganization"`
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
	DevMode           bool                          `yaml:"devMode"`
	TLSCert           *TLSCert                      `yaml:"tlsCert"`
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
