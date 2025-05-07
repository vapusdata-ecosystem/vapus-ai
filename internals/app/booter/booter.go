package booter

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	datasvc "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices"
	datasvcpkg "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
)

func BootPlatform(ctx context.Context, conf *appconfigs.PlatformBootConfig, dbm *apppkgs.VapusStore, svcPkgs *apppkgs.VapusSvcPackages, svcPkgsParams *apppkgs.VapusSvcPackageParams, logger zerolog.Logger) error {
	var err error
	if dbm.Db == nil {
		logger.Fatal().Msg("error while booting VapusData platform")
	}
	logger.Info().Msgf("Platform Boot Config - %v", conf)
	pBooter := NewPlatformSetup(conf, dbm, svcPkgs, svcPkgsParams, logger)

	err = pBooter.AddVapusDataPlatformAccount(ctx)
	if err != nil {
		logger.Fatal().Msgf("error while booting platform account. error: %v", err)
	}

	err = pBooter.AddVapusDataPlatformOwnerOrganization(ctx)
	if err != nil {
		logger.Fatal().Msgf("error while booting platform owner Organization. error: %v", err)
	}

	err = pBooter.AddVapusDataPlatformOwners(ctx)
	if err != nil {
		logger.Fatal().Msgf("error while booting platform owners. error: %v", err)
	}
	pBooter.Clean()
	return nil
}

func BootDataTables(ctx context.Context, cl *apppkgs.VapusStore, logger zerolog.Logger) {
	if cl.Db == nil {
		logger.Fatal().Msg("error while booting ES indexes")
	}
	errChan := make(chan error, len(datasvc.INDEX_LIST))
	var wg sync.WaitGroup
	// dmstores.DMStoreManager.UpdateArns(ctx)
	for name, obj := range apppkgs.DBTablesMap {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err := cl.Db.CreateDataTables(ctx, &datasvcpkg.DataTablesOpts{
				StructScheme: obj,
				Name:         name,
			})
			if err != nil {
				errChan <- err
			} else {
				logger.Info().Msgf("Data table created for %v", name)
			}

		}(&wg)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			logger.Err(err).Msg("error while booting data tables ")
		}
	}
}
