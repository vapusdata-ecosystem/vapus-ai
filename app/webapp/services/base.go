package services

import (
	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/webapp/clients"
	"github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
)

type WebappService struct {
	logger        zerolog.Logger
	grpcClients   *clients.GrpcClient
	templateDir   string
	SvcPkgManager *apppkgs.VapusSvcPackages
}

var WebappServiceManager *WebappService

func NewWebappSvc() *WebappService {
	svc := &WebappService{}
	svc.logger = pkgs.GetSubDMLogger("webapp", "services")
	svc.templateDir = "datamarketplace"
	clients.InitGrpcClient()
	svc.grpcClients = clients.GrpcClientManager
	svc.SvcPkgManager = pkgs.SvcPackageManager
	return svc
}

func InitWebappSvc() {
	if WebappServiceManager == nil {
		WebappServiceManager = NewWebappSvc()
	}
}
