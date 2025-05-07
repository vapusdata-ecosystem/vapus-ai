package clients

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	runnerutils "github.com/vapusdata-ecosystem/vapusdata/core/nabrunners/utils"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

type NabRunnerClient struct {
	redisParams     *models.DataSourceSecrets
	logger          zerolog.Logger
	ConcurrentCalls int32
	RedisOpt        asynq.RedisClientOpt
	Client          *asynq.Client
	Inspector       *asynq.Inspector
	Scheduler       *asynq.Scheduler
}

type NabRunnerClientOpts func(*NabRunnerClient)

func WithClientRedisOpt(obj *models.DataSourceSecrets) NabRunnerClientOpts {
	return func(s *NabRunnerClient) {
		s.redisParams = obj
	}
}

func NewAsynqNabRunnerClient(debugMode bool, opts ...NabRunnerClientOpts) *NabRunnerClient {
	cl := &NabRunnerClient{}
	for _, opt := range opts {
		opt(cl)
	}
	cl.logger = dmlogger.GetDMLogger(debugMode, true, "")
	cl.RedisOpt = runnerutils.GetRedisOpts(cl.redisParams, cl.logger)
	cl.Client = asynq.NewClient(cl.RedisOpt)
	cl.Inspector = asynq.NewInspector(cl.RedisOpt)
	cl.Scheduler = asynq.NewScheduler(cl.RedisOpt, &asynq.SchedulerOpts{})
	return cl
}

func (s *NabRunnerClient) Close() {
	s.Client.Close()
}
