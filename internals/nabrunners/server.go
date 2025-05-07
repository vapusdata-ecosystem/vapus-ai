package nabrunners

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	nabtypes "github.com/vapusdata-ecosystem/vapusdata/core/nabrunners/types"
	runnerutils "github.com/vapusdata-ecosystem/vapusdata/core/nabrunners/utils"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

const (
	ConcurrentCalls = 20 // Configurable via configmap
)

type NabRunnerServer struct {
	MuxServer       *asynq.ServeMux
	redisParams     *models.DataSourceSecrets
	Server          *asynq.Server
	logger          zerolog.Logger
	ConcurrentCalls int32
	Scheduler       *asynq.Scheduler
	Inspector       *asynq.Inspector
}

type NabRunnerOpts func(*NabRunnerServer)

func WithRedisClientOpt(obj *models.DataSourceSecrets) NabRunnerOpts {
	return func(s *NabRunnerServer) {
		s.redisParams = obj
	}
}

func WithConcurrentCalls(obj int32) NabRunnerOpts {
	return func(s *NabRunnerServer) {
		s.ConcurrentCalls = obj
	}
}

func WithInspector(obj *asynq.Inspector) NabRunnerOpts {
	return func(s *NabRunnerServer) {
		s.Inspector = obj
	}
}

func WithScheduler(obj *asynq.Scheduler) NabRunnerOpts {
	return func(s *NabRunnerServer) {
		s.Scheduler = obj
	}
}

func NewAsynqNabRunner(debugMode bool, opts ...NabRunnerOpts) *NabRunnerServer {
	server := &NabRunnerServer{}
	for _, opt := range opts {
		opt(server)
	}
	server.logger = dmlogger.GetDMLogger(debugMode, true, "")
	redisOpt := runnerutils.GetRedisOpts(server.redisParams, server.logger)

	server.Server = asynq.NewServer(redisOpt, asynq.Config{
		Queues: nabtypes.NabRunnerQueues,
	})
	server.MuxServer = asynq.NewServeMux()
	return server
}

func (s *NabRunnerServer) Run() {
	s.logger.Info().Msg("Starting asynq server...")
	for {
		if err := s.Server.Run(s.MuxServer); err != nil {
			s.logger.Fatal().Msgf("asynq server start error error: %v", err)
		}
		time.Sleep(5 * time.Second)
		s.logger.Info().Msg("Restarting asynq server...")
	}
}

func (s *NabRunnerServer) Stop() {
	s.logger.Info().Msg("Stopping asynq server...")
	s.Server.Shutdown()
}

func (s *NabRunnerServer) StopScheduler() {
	s.logger.Info().Msg("Stopping asynq scheduler...")
	s.Scheduler.Shutdown()
}

func (s *NabRunnerServer) RunScheduler() {
	s.logger.Info().Msg("Starting asynq scheduler...")
	for {
		if err := s.Scheduler.Run(); err != nil {
			s.logger.Fatal().Msgf("asynq server start error error: %v", err)
		}
		time.Sleep(5 * time.Second)
		s.logger.Info().Msg("Restarting asynq server...")
	}
}
