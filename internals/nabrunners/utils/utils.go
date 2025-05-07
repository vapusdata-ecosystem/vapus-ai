package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func GetRedisOpts(redisOps *models.DataSourceSecrets, logger zerolog.Logger) asynq.RedisClientOpt {
	var err error
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", redisOps.URL, redisOps.Port),
	}
	log.Println("Redis URL: >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ", redisOpt.Addr)
	if redisOps.Password != "" {
		redisOpt.Password = redisOps.Password
	}
	log.Println("redisOps.DB: >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ", redisOps.DB)
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
	log.Println("redisOpt: >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ", redisOpt)
	log.Println("redisOpt: >>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ", redisOpt.DB)
	return redisOpt
}
