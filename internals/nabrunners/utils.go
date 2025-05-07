package nabrunners

import (
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func GetRedisOpts(redisOps *models.DataSourceSecrets, logger zerolog.Logger) asynq.RedisClientOpt {
	var err error
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", redisOps.URL, redisOps.Port),
	}
	if redisOps.Password != "" {
		redisOpt.Password = redisOps.Password
	}
	if redisOps.DB != "" {
		redisOpt.DB, err = strconv.Atoi(redisOps.DB)
		if err != nil {
			logger.Err(err).Msgf("Error converting db to int: %v", err)
			redisOpt.DB = 0
		}
	}
	if redisOps.Username != "" {
		redisOpt.Username = redisOps.Username
	}
	return redisOpt
}
