package nabrunners

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	runnerutils "github.com/vapusdata-ecosystem/vapusdata/core/nabrunners/utils"
)

type RunnerInspector struct {
	Inspector   *asynq.Inspector
	redisParams *models.DataSourceSecrets
	logger      zerolog.Logger
}

type InspectorOpts func(*RunnerInspector)

func InspectorWithRedisClientOpt(obj *models.DataSourceSecrets) InspectorOpts {
	return func(s *RunnerInspector) {
		s.redisParams = obj
	}
}

func NewAsynqInspector(debugMode bool, opts ...InspectorOpts) *RunnerInspector {
	sc := &RunnerInspector{}
	for _, opt := range opts {
		opt(sc)
	}
	redOpts := runnerutils.GetRedisOpts(sc.redisParams, sc.logger)
	sc.Inspector = asynq.NewInspector(redOpts)
	return sc
}

func (x *RunnerInspector) GetNextTaskRunAt(queue, taskId string, logger zerolog.Logger) int64 {
	tasks, err := x.Inspector.ListScheduledTasks(queue, taskId)
	if err != nil {
		logger.Err(err).Msgf("Error getting task info: %v", taskId)
		return 0
	}
	for _, task := range tasks {
		if task.ID == taskId {
			logger.Info().Msgf("Task: %v", task.ID)
			return task.NextProcessAt.Unix()
		}
	}
	return 0
}
