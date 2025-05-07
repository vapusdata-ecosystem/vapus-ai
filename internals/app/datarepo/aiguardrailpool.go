package datarepo

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	gdrl "github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

var GuardrailPoolManager *GuardrailPool

type GuardrailPool struct {
	// ModelGuardRails    map[string][]string
	// GuardrailClientMap map[string]*gdrl.GuardRailClient
	// ModelGuardrailMap  map[string][]string
	ModelGuardRails    sync.Map
	GuardrailClientMap sync.Map
	ModelGuardrailMap  sync.Map
	AccountGuardRails  []string
	dmstores           *apppkgs.VapusStore
	logger             zerolog.Logger
	nodePool           *AIModelNodeConnectionPool
}

type guardrailPoolOpts func(*GuardrailPool)

func WithGpStore(dmstores *apppkgs.VapusStore) guardrailPoolOpts {
	return func(gp *GuardrailPool) {
		gp.dmstores = dmstores
	}
}

func WithGpLogger(logger zerolog.Logger) guardrailPoolOpts {
	return func(gp *GuardrailPool) {
		gp.logger = logger
	}
}

func InitGuardrailPool(ctx context.Context, nodePool *AIModelNodeConnectionPool, opts ...guardrailPoolOpts) (*GuardrailPool, error) {
	obj := &GuardrailPool{
		ModelGuardRails:    sync.Map{},
		AccountGuardRails:  make([]string, 0),
		GuardrailClientMap: sync.Map{},
		nodePool:           nodePool,
	}

	for _, opt := range opts {
		opt(obj)
	}

	modelNodes, err := ListAIModelNodes(ctx, obj.dmstores, obj.logger, "status = 'ACTIVE' AND deleted_at IS NULL ORDER BY created_at DESC", nil)
	if err != nil {
		return nil, err
	}
	// guardrails, err := ListAIGuardrails(ctx, obj.dmstores, obj.logger, "", nil)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, guardrail := range guardrails {
	// 	nodePool := make([]*gdrl.GuardModelNodePool, 0)
	// 	attr := dmstores.Account.AIAttributes
	// 	if attr.GuardrailModelNode == "" || attr.GuardrailModel == "" {
	// 		continue
	// 	}
	// 	nodeConn := AIModelNodeConnectionPoolManager.GetConnectionById(attr.GuardrailModelNode)

	// 	nodePool = append(nodePool, &gdrl.GuardModelNodePool{
	// 		Connection: nodeConn,
	// 		IsAccount:  true,
	// 		Model:      attr.GuardrailModel,
	// 	})
	// 	gd := gdrl.New(
	// 		gdrl.WithSpec(guardrail),
	// 		gdrl.WithModelPool(nodePool),
	// 	)
	// 	obj.GuardrailClientMap[guardrail.VapusID] = gd
	// }
	for _, modelNode := range modelNodes {
		if modelNode.SecurityGuardrails != nil {
			for _, guardrail := range modelNode.SecurityGuardrails.Guardrails {
				if val, ok := obj.GuardrailClientMap.Load(guardrail); ok {
					gd, valid := val.(*gdrl.GuardRailClient)
					if !valid {
						continue
					}
					obj.AddGuardrails(modelNode.VapusID, gd.Guardrail.VapusID)
				}
			}
		}
	}
	return obj, nil
}

func (gp *GuardrailPool) AddGuardrails(modelsNode string, guardrailId string) {
	if modelsNode == "" {
		gp.AccountGuardRails = append(gp.AccountGuardRails, guardrailId)
		return
	} else {
		if val, ok := gp.ModelGuardRails.Load(modelsNode); !ok {
			gp.ModelGuardRails.Store(modelsNode, make([]string, 0))
		} else {
			gdList, valid := val.([]string)
			if !valid {
				return
			}
			gdList = append(gdList, guardrailId)
			gp.ModelGuardRails.Store(modelsNode, gdList)
		}
	}
}

func (gp *GuardrailPool) GetGuardrail(guardrailId string) *gdrl.GuardRailClient {
	if guardrail, ok := gp.GuardrailClientMap.Load(guardrailId); ok {
		gdObj, valid := guardrail.(*gdrl.GuardRailClient)
		if valid {
			return gdObj
		}
		return nil
	}
	return nil
}

func (gp *GuardrailPool) UpdateGuardrailPool(guardrail *models.AIGuardrails) {
	nodePool := make([]*gdrl.GuardModelNodePool, 0)
	nodeConn := gp.nodePool.GetConnectionById(guardrail.GuardModel.ModelNodeID)
	if nodeConn == nil {
		return
	}
	nodePool = append(nodePool, &gdrl.GuardModelNodePool{
		Connection: nodeConn,
		IsAccount:  true,
		Model:      guardrail.GuardModel.ModelID,
	})
	gd := gdrl.New(
		gdrl.WithSpec(guardrail),
		gdrl.WithModelPool(nodePool),
	)
	gp.GuardrailClientMap.Store(guardrail.VapusID, gd)
	return
}

func (gp *GuardrailPool) RemoveGuardrail(guardrailId string) {
	for i, guardrail := range gp.AccountGuardRails {
		if guardrail == guardrailId {
			gp.AccountGuardRails = append(gp.AccountGuardRails[:i], gp.AccountGuardRails[i+1:]...)
		}
	}
	gp.ModelGuardRails.Range(func(key, value interface{}) bool {
		if guardrails, ok := value.([]string); ok {
			for i, guardrail := range guardrails {
				if guardrail == guardrailId {
					newGuardrails := append(guardrails[:i], guardrails[i+1:]...)
					gp.ModelGuardRails.Store(key, newGuardrails)
					break
				}
			}
		}
		return true
	})
	// for _, guardrails := range gp.ModelGuardRails.Range() {
	// 	for i, guardrail := range guardrails {
	// 		if guardrail == guardrailId {
	// 			guardrails = append(guardrails[:i], guardrails[i+1:]...)
	// 		}
	// 	}
	// }
}
