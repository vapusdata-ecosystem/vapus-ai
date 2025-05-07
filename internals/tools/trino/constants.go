package trinocl

import "fmt"

type TrinoPropertiesString string

var Secretproperttemplate = "${%s}"

var (
	UsernameTemplate          TrinoPropertiesString = "%s.username=%s\n"
	PasswordTemplate          TrinoPropertiesString = "%s.password=%s\n"
	AWSAccessKeyTemplate      TrinoPropertiesString = "%s.aws-access-key-id=%s\n"
	AWSSecretKeyTemplate      TrinoPropertiesString = "%s.aws-secret-access-key=%s\n"
	GCPAccessKeyTemplate      TrinoPropertiesString = "%s.google-access-key-id=%s\n"
	GCPSecretKeyTemplate      TrinoPropertiesString = "%s.google-secret-access-key=%s\n"
	ApiTokenTemplate          TrinoPropertiesString = "%s.api-token=%s\n"
	ApiKeyTemplate            TrinoPropertiesString = "%s.api-key=%s\n"
	AzureClientIdTemplate     TrinoPropertiesString = "%s.azure-client-id=%s\n"
	AzureClientSecretTemplate TrinoPropertiesString = "%s.azure-client-secret=%s\n"
)

const (
	TrinoCatalogSecretsMountName = "vapusdata-trino-catalog-secrets-mount"
	TrinoCatalogMountPath        = "/etc/trino/catalog"
	TrinoCatalogMountName        = "vapusdata-trino-catalog-mount"
	TrinoCatalogDefault          = "my-trino-trino-catalog"
	TrinoCatalogSecrets          = "vapusdata-trino-catalog-secrets"
	TrinoCatalogSecretsMountpath = "/etc/trino/vapusdata/catalog/secrets"
)

func (cc TrinoPropertiesString) Render(tag, value string) string {
	return fmt.Sprintf(string(cc), tag, value)
}

func (cc TrinoPropertiesString) String() string {
	return string(cc)
}

const (
	DBName                               = "db"
	VapusTrinoDBQueryUri                 = "vapus.trino.db.query.uri"
	DataSourceId                         = "datasource.id"
	TrinoCredentialProviderType_FILE     = "FILE"
	TrinoCredentialProviderType_INLINE   = "INLINE"
	TrinoCredentialProviderType_KEYSTORE = "KEYSTORE"
	DEFAULT_USERNAME                     = "vapusdata"
)
