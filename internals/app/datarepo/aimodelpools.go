package datarepo

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	aimodels "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type AIModelNodeConnectionPool struct {
	connectionPool sync.Map //map[string]aimodels.AIModelNodeInterface
	logger         zerolog.Logger
	dmStore        *apppkgs.VapusStore
	errPool        map[string]error
	objectPool     sync.Map //map[string]*models.AIModelNode
}

type aiPoolOpts func(*AIModelNodeConnectionPool)

func WithMpLogger(logger zerolog.Logger) aiPoolOpts {
	return func(a *AIModelNodeConnectionPool) {
		a.logger = logger
	}
}

func WithMpDMStore(dmStore *apppkgs.VapusStore) aiPoolOpts {
	return func(a *AIModelNodeConnectionPool) {
		a.dmStore = dmStore
	}
}

func InitAIModelNodeConnectionPool(obj *AIModelNodeConnectionPool, opts ...aiPoolOpts) *AIModelNodeConnectionPool {
	if obj != nil {
		return obj
	}
	obj = &AIModelNodeConnectionPool{}
	for _, opt := range opts {
		opt(obj)
	}
	// obj.connectionPool = map[string]aimodels.AIModelNodeInterface{}
	// obj.objectPool = map[string]*models.AIModelNode{}
	obj.connectionPool = sync.Map{}
	obj.objectPool = sync.Map{}
	obj.errPool = map[string]error{}
	ctx, cancel := context.WithCancel(context.Background())
	obj.bootConnectionPool(ctx)
	defer cancel()
	if len(obj.errPool) > 0 {
		obj.logger.Error().Msg("error while booting connection pool")
	}
	return obj
}

func (a *AIModelNodeConnectionPool) AddConnection(model *models.AIModelNode, connection aimodels.AIModelNodeInterface) {
	a.connectionPool.Store(model.VapusID, connection)
}

func (a *AIModelNodeConnectionPool) GetConnectionById(nodeId string) aimodels.AIModelNodeInterface {
	val, ok := a.connectionPool.Load(nodeId)
	if !ok {
		a.logger.Info().Msgf("Connection not found in pool for %s", nodeId)
		return nil
	} else {
		intVal, valid := val.(aimodels.AIModelNodeInterface)
		if !valid {
			a.logger.Info().Msgf("Connection not found in pool for %s", nodeId)
			return nil
		}
		return intVal
	}
}

func (a *AIModelNodeConnectionPool) GetorSetNodeObject(nodeId string, nodeObject *models.AIModelNode, add bool) (*models.AIModelNode, error) {
	fmt.Println("I am in get or set node: ", nodeId)
	// fmt.Println("value is: ", val)
	if add {
		a.logger.Info().Msgf("Force update for %s", nodeId)
		if nodeObject != nil {
			a.objectPool.Store(nodeId, nodeObject)
			return nodeObject, nil
		} else {
			a.logger.Error().Msgf("Node object is nil for %s", nodeId)
			return nil, dmerrors.DMError(apperr.ErrAIModelNode404, nil)
		}
	}
	val, ok := a.objectPool.Load(nodeId)
	if !ok {
		if nodeObject == nil {
			a.logger.Info().Msgf("Connection not found in local pool for %s , getting it from cache", nodeId)
			nodeObject = a.getObjectCache(nodeId)
			if nodeObject != nil {
				return nodeObject, nil
			}
		} else {
			a.addObjectCache(nodeObject)
			a.objectPool.Store(nodeId, nodeObject)
			return nodeObject, nil
		}
	} else {
		modelNode, valid := val.(*models.AIModelNode)
		// fmt.Println("ModelNode: ", modelNode.Name)
		if !valid {
			return nil, dmerrors.DMError(apperr.ErrAIModelNode404, nil)
		}
		return modelNode, nil
	}
	return nil, dmerrors.DMError(apperr.ErrAIModelNode404, nil)
}

func (a *AIModelNodeConnectionPool) addObjectCache(model *models.AIModelNode) {
	mBytes, err := dmutils.Marshall(model)
	if err != nil {
		a.logger.Err(err).Msg("error while marshalling model")
		a.errPool[model.VapusID] = dmerrors.DMError(apperr.ErrAIModelNode404, err)
		return
	}
	_, err = a.dmStore.Cacher.RedisClient.WriteKV(context.TODO(), model.VapusID, string(mBytes))
	if err != nil {
		a.logger.Err(err).Msg("error while writing model to cache")
		a.errPool[model.VapusID] = dmerrors.DMError(apperr.ErrAIModelNode404, err)
		return
	}
}

func (a *AIModelNodeConnectionPool) getObjectCache(nodeId string) *models.AIModelNode {
	mBytes, err := a.dmStore.Cacher.RedisClient.ReadKV(context.TODO(), nodeId)
	if err != nil {
		a.logger.Err(err).Msg("error while writing model to cache")
		return nil
	}
	model := &models.AIModelNode{}
	err = dmutils.Unmarshall([]byte(mBytes.(string)), model)
	if err != nil {
		a.logger.Err(err).Msg("error while unmarshalling model from cache")
		return nil
	}
	if model.VapusID == "" {
		a.logger.Error().Msg("model VapusID is empty")
		return nil
	}
	return model
}

func (a *AIModelNodeConnectionPool) GetorSetConnection(model *models.AIModelNode, addIfNotPresent, forceupdate bool) (aimodels.AIModelNodeInterface, error) {
	val, ok := a.connectionPool.Load(model.VapusID)
	if !ok {
		a.logger.Info().Msgf("Connection not found in pool for %s , creating new connection", model.VapusID)
		if !addIfNotPresent {
			return nil, dmerrors.DMError(apperr.ErrAIModelConn, nil)
		} else {
			a.logger.Info().Msg("++++++++++++++++++++++ GET OR SET NODE OBJECT")
			a.logger.Info().Msg(model.Name)
			a.logger.Info().Msg(string(model.ID))
			a.logger.Info().Msg(model.VapusID)

			a.createModelConnection(context.Background(), model)
			val, ok = a.connectionPool.Load(model.VapusID)
			if !ok {
				return nil, dmerrors.DMError(apperr.ErrAIModelConn, nil)
			}
			nodeInterface, valid := val.(aimodels.AIModelNodeInterface)
			if !valid {
				return nil, dmerrors.DMError(apperr.ErrAIModelConn, nil)
			}
			return nodeInterface, nil
		}
	}
	if forceupdate {
		a.logger.Info().Msgf("Force update for %s", model.VapusID)
		a.createModelConnection(context.Background(), model)
		val, ok = a.connectionPool.Load(model.VapusID)
		if !ok {
			return nil, dmerrors.DMError(apperr.ErrAIModelConn, nil)
		}
	}
	nodeInterface, valid := val.(aimodels.AIModelNodeInterface)
	if !valid {
		return nil, dmerrors.DMError(apperr.ErrAIModelConn, nil)
	}
	return nodeInterface, nil
}

func (a *AIModelNodeConnectionPool) RemoveConnection(model *models.AIModelNode) {
	a.connectionPool.Delete(model.VapusID)
}

func (a *AIModelNodeConnectionPool) RemoveNodeObject(nodeId string) {
	a.objectPool.Delete(nodeId)
}

func (a *AIModelNodeConnectionPool) bootConnectionPool(ctx context.Context) error {
	result, err := ListAIModelNodes(ctx, a.dmStore, a.logger, "status = 'ACTIVE' AND deleted_at IS NULL ORDER BY created_at DESC", nil)
	if err != nil {
		a.logger.Error().Err(err).Msg("error while fetching models from datastore")
		return err
	}
	var wg sync.WaitGroup
	for _, model := range result {
		wg.Add(1)
		go func(model *models.AIModelNode) {
			defer wg.Done()
			ctx := context.Background()
			a.objectPool.Store(model.VapusID, model)
			a.createModelConnection(ctx, model)
		}(model)
	}
	wg.Wait()
	return nil
}

func (a *AIModelNodeConnectionPool) createModelConnection(ctx context.Context, model *models.AIModelNode) error {
	if model.NetworkParams.SecretName != "" {
		secrets, err := apppkgs.ReadCredentialFromStore(ctx, model.NetworkParams.SecretName, a.dmStore, a.logger)
		if err != nil {
			a.logger.Err(err).Msg("error while reading credentials from store")
			a.errPool[model.VapusID] = dmerrors.DMError(apperr.ErrGetAIModelNetParams, err)
			return err
		}

		netParam := &models.AIModelNodeNetworkParams{
			Url:                 model.NetworkParams.GetUrl(),
			Credentials:         secrets,
			ApiVersion:          model.NetworkParams.GetApiVersion(),
			LocalPath:           model.NetworkParams.GetLocalPath(),
			SecretName:          model.NetworkParams.SecretName,
			IsAlreadyInSecretBs: model.NetworkParams.IsAlreadyInSecretBs,
		}

		model.NetworkParams = netParam
		conn, err := aimodels.NewAIModelNode(aimodels.WithAIModelNode(model), aimodels.WithLogger(a.logger))
		if err != nil {
			a.logger.Err(err).Ctx(ctx).Msg("error while creating model connection")
			a.errPool[model.VapusID] = dmerrors.DMError(apperr.ErrAIModelConn, err)
		}
		a.logger.Info().Ctx(ctx).Msgf("Connection created for AI model %s", model.VapusID)
		a.AddConnection(model, conn)
	} else {
		a.logger.Info().Ctx(ctx).Msgf("No secret found for AI model %s", model.VapusID)
		a.RemoveConnection(model)
	}
	return nil
}
