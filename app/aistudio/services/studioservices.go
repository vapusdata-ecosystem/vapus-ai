package services

import (
	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	aidmstore "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo/aistudio"
)

// AIStudioServices is a struct that contains the DMStore.
type AIStudioServices struct {
	DMStore *aidmstore.AIStudioDMStore
	Logger  zerolog.Logger
}

// AIStudioServicesManager is the global variable for AIStudioServices struct.
var AIStudioServiceManager *AIStudioServices

// newAIStudioServices creates a new object for AIStudioServices struct.
func newAIStudioServices(dmstore *aidmstore.AIStudioDMStore) *AIStudioServices {
	return &AIStudioServices{
		DMStore: dmstore,
		Logger:  pkgs.GetSubDMLogger(pkgs.SVCS, "AIStudioServices"),
	}
}

// InitAIStudioServices initializes the data mesh services.
func InitAIStudioServices(dmstore *aidmstore.AIStudioDMStore) {
	if AIStudioServiceManager == nil {
		AIStudioServiceManager = newAIStudioServices(dmstore)
	}
}
