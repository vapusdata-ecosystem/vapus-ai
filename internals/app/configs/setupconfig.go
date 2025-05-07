package appconfigs

import (
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
)

type VapusInstallerConfig struct {
	App struct {
		Name         string `yaml:"name"`
		Namespace    string `yaml:"namespace"`
		Organization string `yaml:"organization"`
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
		Datamarketplace struct {
			Name    string `yaml:"name"`
			Creator string `yaml:"creator"`
		} `yaml:"datamarketplace"`
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
	Vault   *Vault `yaml:"vault"`
	Trino   *Trino `yaml:"trino"`
	TLSCert struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	} `yaml:"tlsCert"`
	CreateDatabase bool                          `yaml:"createDatabase"`
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

type Trino struct {
	FullnameOverride        string `yaml:"fullnameOverride"`
	CoordinatorNameOverride string `yaml:"coordinatorNameOverride"`
	WorkerNameOverride      string `yaml:"workerNameOverride"`
	Service                 struct {
		Port int `yaml:"port"`
	} `yaml:"service"`
	Coordinator struct {
		Resources struct {
			Requests struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"requests"`
		} `yaml:"resources"`
		Jvm struct {
			MaxHeapSize string `yaml:"maxHeapSize"`
			Gc          string `yaml:"gc"`
		} `yaml:"jvm"`
	} `yaml:"coordinator"`
	Worker struct {
		Replicas  int `yaml:"replicas"`
		Resources struct {
			Requests struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"requests"`
		} `yaml:"resources"`
	} `yaml:"worker"`
	Config struct {
		Coordinator struct {
			LogLevel string `yaml:"logLevel"`
		} `yaml:"coordinator"`
		Properties struct {
			QueryMaxMemory        string `yaml:"query.max-memory"`
			QueryMaxMemoryPerNode string `yaml:"query.max-memory-per-node"`
		} `yaml:"properties"`
	} `yaml:"config"`
}

type VapusSecretInstallerConfig struct {
	SecretStore       *models.DataSourceCredsParams `yaml:"secretStore"`
	DevSecretStore    *models.DataSourceCredsParams `yaml:"devSecretStore"`
	BackendDataStore  *models.DataSourceCredsParams `yaml:"backendDataStore"`
	BackendCacheStore *models.DataSourceCredsParams `yaml:"backendCacheStore"`
	FileStore         *models.DataSourceCredsParams `yaml:"fileStore"`
	JWTAuthnSecrets   *encryption.JWTAuthn          `yaml:"JWTAuthnSecrets"`
	AuthnSecrets      *authn.AuthnSecrets           `yaml:"authnSecrets"`
	ArtifactStore     *models.DataSourceCredsParams `yaml:"artifactStore"`
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
	FileStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"fileStore"`
	JWTAuthnSecrets struct {
		Secret string `yaml:"secret"`
	} `yaml:"JWTAuthnSecrets"`
	AuthnSecrets struct {
		Secret string `yaml:"secret"`
	} `yaml:"authnSecrets"`
	ArtifactStore struct {
		Secret string `yaml:"secret"`
	} `yaml:"artifactStore"`
}
