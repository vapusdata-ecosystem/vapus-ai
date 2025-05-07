package pkgs

import (
	"context"

	validator "github.com/go-playground/validator/v10"
	appcl "github.com/vapusdata-ecosystem/vapusdata/core/app/grpcclients"
)

var DmValidator = validator.New()

const (
	IDEN    = "Identifier"
	CNTRLR  = "controller"
	SVCS    = "services"
	DSTORES = "dmstores"
)

var VapusSvcInternalClientManager *appcl.VapusSvcInternalClients

func InitVapusSvcInternalClients() {
	res, err := appcl.SetupVapusSvcInternalClients(context.Background(), NetworkConfigManager, NetworkConfigManager.AIStudioSvc.ServiceName, pkgLogger)
	if err != nil {
		pkgLogger.Err(err).Msg("error while initializing vapus svc internal clients")
	}
	VapusSvcInternalClientManager = res
}

var AIAgentSystemMessage = `You are Agent steps formatter that takes imput in to plain text and then extract the info and put it tool response based on relevant fields.If data product or query is present in user input then in query field, just provide the descriptive text query, don't convert it into SQL or any other format because query generation based on data engine will be the next step in further calls.All fields are optional so in case of missing fields provide empty values of the respective fields. If data product query is not provided by user then eliminate that property. Data product STRICTLY is an uuid, if not found then pass empty or eliminate.`
