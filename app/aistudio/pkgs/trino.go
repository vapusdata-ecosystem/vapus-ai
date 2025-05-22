package pkgs

import (
	sqlops "github.com/vapusdata-ecosystem/vapusai/core/tools/sqlops"
	trinocl "github.com/vapusdata-ecosystem/vapusai/core/tools/trino"
)

var TrinoClient *trinocl.TrinoClient

var SqlOps *sqlops.SQLOperator

// const (
// 	TrinoCatalogSecretsMountName = "vapusdata-trino-catalog-secrets-mount"
// 	TrinoCatalogMountPath        = "/etc/trino/catalog"
// 	TrinoCatalogMountName        = "vapusdata-trino-catalog-mount"
// 	TrinoCatalog                 = "my-trino-trino-catalog"
// 	TrinoCatalogSecrets          = "vapusdata-trino-catalog-secrets"
// 	TrinoCatalogSecretsMountpath = "/etc/trino/vapusdata/catalog/secrets"
// )

// func newTrinoClient() error {
// 	if ServiceConfigManager == nil || ServiceConfigManager.TrinoSpecs == nil {
// 		pkgLogger.Error().Msg("Trino client not initialized, service config not found")
// 		return apperr.ErrServiceConfigNotInitialized
// 	}
// 	if ServiceConfigManager.TrinoSpecs.TrinoCatalog == "" {
// 		ServiceConfigManager.TrinoSpecs.TrinoCatalog = TrinoCatalog
// 	}
// 	TrinoClient = trinocl.New(DmLogger,
// 		trinocl.WithCatalog(ServiceConfigManager.TrinoSpecs.TrinoCatalog),
// 		trinocl.WithCatalogSecretsMountName(TrinoCatalogSecretsMountName),
// 		trinocl.WithCatalogMountName(TrinoCatalogMountName),
// 		trinocl.WithCatalogSecretsMountPath(TrinoCatalogSecretsMountpath),
// 		trinocl.WithCatalogMountPath(TrinoCatalogMountPath),
// 		trinocl.WithCatalogSecrets(TrinoCatalogSecrets),
// 		trinocl.WithDeploymentSpec(ServiceConfigManager.TrinoSpecs),
// 		trinocl.WithReaderSpec(true),
// 	)
// 	return nil
// }

// func InitTrinoClient() error {
// 	if TrinoClient == nil {
// 		if err := newTrinoClient(); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	return nil
// }

func InitSqlOps() error {
	var err error
	if SqlOps == nil {
		SqlOps, err = sqlops.New(DmLogger)
		if err != nil {
			return err
		}

	}
	return nil
}
