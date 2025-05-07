package nabrunners

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	runnerutils "github.com/vapusdata-ecosystem/vapusai/core/nabrunners/utils"
)

type RunnerScheduller struct {
	Scheduler   *asynq.Scheduler
	redisParams *models.DataSourceSecrets
	logger      zerolog.Logger
}

type SchedullerOpts func(*RunnerScheduller)

func SchedulerWithRedisClientOpt(obj *models.DataSourceSecrets) SchedullerOpts {
	return func(s *RunnerScheduller) {
		s.redisParams = obj
	}
}

func NewAsynqScheduller(debugMode bool, opts ...SchedullerOpts) *RunnerScheduller {
	sc := &RunnerScheduller{}
	for _, opt := range opts {
		opt(sc)
	}
	redOpts := runnerutils.GetRedisOpts(sc.redisParams, sc.logger)
	sc.Scheduler = asynq.NewScheduler(redOpts, &asynq.SchedulerOpts{})
	return sc
}
