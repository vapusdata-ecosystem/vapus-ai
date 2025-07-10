package datarepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func BootAccountCache(dbStore *apppkgs.VapusStore, logger zerolog.Logger) *models.Account {
	var accountPool *models.Account
	// if redis client is already present then we can read form there...
	if dbStore.Cacher != nil && dbStore.Cacher.RedisClient != nil {
		cachedData, err := dbStore.Cacher.RedisClient.ReadKV(context.TODO(), types.AccountCacheKey.String())
		if err == nil && cachedData != nil {
			err = dmutils.Unmarshall([]byte(cachedData.(string)), accountPool)
			if err != nil {

			} else {
				return accountPool
			}
		}
	}
	// else I have to fetch the accounts details. And then add it to the redis
	var result []*models.Account
	query := fmt.Sprintf("SELECT * FROM %v", apppkgs.AccountsTable)
	err := dbStore.Db.PostgresClient.SelectInApp(context.TODO(), &query, &result)
	if err != nil || len(result) == 0 {
		logger.Fatal().Err(err).Msg("error while fetching accounts info from datastore")
		return &models.Account{}
	}
	data, err := dmutils.Marshall(result[0])
	if err != nil {
		logger.Err(err).Msg("error while marshalling account data")
		return result[0]
	}
	_, err = dbStore.Cacher.RedisClient.WriteKV(context.TODO(), types.AccountCacheKey.String(), string(data))
	if err != nil {
		logger.Err(err).Msg("error while writing account cache to redis")
		return result[0]
	}
	return result[0]
}

func UpdateBootAccountCache(dbStore *apppkgs.VapusStore, accountIds []string, pool map[string]*models.Account, logger zerolog.Logger) {
	var result []*models.Account
	query := fmt.Sprintf("SELECT * FROM %v where vapus_id in (%s)", apppkgs.AccountsTable, "'"+strings.Join(accountIds, "','")+"'")
	err := dbStore.Db.PostgresClient.SelectInApp(context.TODO(), &query, &result)
	if err != nil {
		logger.Err(err).Msg("error while fetching accounts info from datastore")
	}
	if len(result) == 0 {
		logger.Error().Msg("No account found in the system")
	}
	for _, account := range result {
		pool[account.VapusID] = account
	}
	if pool != nil {
		data, err := dmutils.Marshall(pool)
		if err != nil {
			logger.Err(err).Msg("error while marshalling account data")
			return
		}
		_, err = dbStore.Cacher.RedisClient.WriteKV(context.TODO(), types.AccountCacheKey.String(), string(data))
		if err != nil {
			logger.Err(err).Msg("error while writing account cache to redis")
			return
		}
		logger.Info().Msg("Account cache updated successfully")
	} else {
		logger.Error().Msg("No account found to update in cache")
	}
}
