package pkgs

import (
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	trinocl "github.com/vapusdata-ecosystem/vapusdata/core/tools/trino"
)

var TrinoClient *trinocl.TrinoClient

func newTrinoClient() error {
	if ServiceConfigManager == nil || ServiceConfigManager.TrinoSpecs == nil {
		pkgLogger.Error().Msg("Trino client not initialized, service config not found")
		return apperr.ErrServiceConfigNotInitialized
	}
	if ServiceConfigManager.TrinoSpecs.TrinoCatalog == "" {
		ServiceConfigManager.TrinoSpecs.TrinoCatalog = trinocl.TrinoCatalogDefault
	}
	TrinoClient = trinocl.New(DmLogger,
		trinocl.WithCatalog(ServiceConfigManager.TrinoSpecs.TrinoCatalog),
		trinocl.WithCatalogSecretsMountName(trinocl.TrinoCatalogSecretsMountName),
		trinocl.WithCatalogMountName(trinocl.TrinoCatalogMountName),
		trinocl.WithCatalogSecretsMountPath(trinocl.TrinoCatalogSecretsMountpath),
		trinocl.WithCatalogMountPath(trinocl.TrinoCatalogMountPath),
		trinocl.WithCatalogSecrets(trinocl.TrinoCatalogSecrets),
		trinocl.WithDeploymentSpec(ServiceConfigManager.TrinoSpecs),
		trinocl.WithReaderSpec(false),
	)
	return nil
}

func InitTrinoClient() error {
	if TrinoClient == nil {
		if err := newTrinoClient(); err != nil {
			return err
		}
		return nil
	}
	return nil
}
