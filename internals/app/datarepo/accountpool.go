package datarepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func BootAccountCache(dbStore *apppkgs.VapusStore, logger zerolog.Logger) map[string]*models.Account {
	var result []*models.Account
	var AccountPool = make(map[string]*models.Account)
	query := fmt.Sprintf("SELECT * FROM %v", apppkgs.AccountsTable)
	err := dbStore.Db.PostgresClient.SelectInApp(context.TODO(), &query, &result)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while fetching accounts info from datastore")
	}
	if len(result) == 0 {
		logger.Fatal().Msg("No account found in the system")
	}
	for _, account := range result {
		AccountPool[account.VapusID] = account
	}
	return AccountPool
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
}
