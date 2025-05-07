package apppkgs

import (
	"log"

	validator "github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	pbac "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbac"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/tools/k8s"
	"k8s.io/client-go/kubernetes"
)

type VapusSvcPackageParams struct {
	JwtParams      *encryption.JWTAuthn
	AuthnParams    *authn.AuthnSecrets
	PbacConfigPath string
	NeedAuthn      bool `default:"false"`
	NeedPbac       bool `default:"false"`
}

type VapusSvcPackages struct {
	VapusJwtAuth         *encryption.VapusDataJwtAuthn
	AuthnManager         *authn.Authenticator
	PlatformRBACManager  *pbac.PbacConfig
	GrpcRequestValidator *dmutils.DMValidator
	ModelValidator       *validator.Validate
	ValidEnums           map[string]map[string]int32
	HostK8SClient        *kubernetes.Clientset
}

type VapusSvcPkgOpts func(*VapusSvcPackageParams)

func WithJwtParams(params *encryption.JWTAuthn) VapusSvcPkgOpts {
	return func(p *VapusSvcPackageParams) {
		p.JwtParams = params
	}
}

func WithAuthnParams(params *authn.AuthnSecrets) VapusSvcPkgOpts {
	return func(p *VapusSvcPackageParams) {
		p.AuthnParams = params
		p.NeedAuthn = true
	}
}

func WithPbacConfigPath(path string) VapusSvcPkgOpts {
	return func(p *VapusSvcPackageParams) {
		p.PbacConfigPath = path
		p.NeedPbac = true
	}
}

func InitSvcPackages(params *VapusSvcPackageParams, VapusSvcPackageManager *VapusSvcPackages, logger zerolog.Logger, opts ...VapusSvcPkgOpts) (*VapusSvcPackageParams, *VapusSvcPackages, error) {
	var err error
	params = &VapusSvcPackageParams{}
	for _, opt := range opts {
		opt(params)
	}
	log.Println(params, "========================>>>>>>>>>>>>>>>>>>>>>>>>")
	if VapusSvcPackageManager == nil {
		VapusSvcPackageManager = &VapusSvcPackages{}
	}
	VapusSvcPackageManager.ModelValidator = validator.New()
	if VapusSvcPackageManager.GrpcRequestValidator == nil {
		VapusSvcPackageManager.GrpcRequestValidator, err = dmutils.NewDMValidator()
		if err != nil {
			logger.Err(err).Msg("Error while loading validator")
			return nil, nil, ErrValidatorInitFailed
		}
		logger.Info().Msg("GRPC request Validator initialized successfully")
	}

	if params.NeedAuthn {
		if VapusSvcPackageManager.AuthnManager == nil {
			if params.AuthnParams == nil {
				logger.Err(ErrAuthenticatorParamsNil).Msg("Error while loading authn config")
				return nil, nil, ErrAuthenticatorParamsNil
			}
			log.Println("params.AuthnParams - ", params.AuthnParams)
			VapusSvcPackageManager.AuthnManager, err = authn.New(params.AuthnParams.OIDCSecrets, params.AuthnParams.AuthnMethod)
			if err != nil {
				logger.Err(err).Msg("Error while initializing authenticator")
				return nil, nil, ErrAuthenticatorInitFailed
			}
			logger.Info().Msg("Authenticator initialized successfully")
		}
	}
	if VapusSvcPackageManager.VapusJwtAuth == nil {
		if params.JwtParams == nil {
			logger.Err(ErrJwtParamsNil).Msg("Error while loading jwt config")
			return nil, nil, ErrJwtParamsNil
		}
		log.Println("params.JwtParams - ", params.JwtParams)
		VapusSvcPackageManager.VapusJwtAuth, err = encryption.NewVapusDataJwtAuthn(params.JwtParams)
		if err != nil {
			logger.Err(err).Msg("Error while initializing jwt authn")
			return nil, nil, ErrJwtAuthInitFailed
		}
		if err := VapusSvcPackageManager.ModelValidator.Struct(VapusSvcPackageManager.VapusJwtAuth.Opts); err != nil {
			logger.Err(err).Msg("Error while validating jwt config")
			return nil, nil, err
		}
		logger.Info().Msg("JWT Authn initialized successfully")
	}

	if params.NeedPbac {
		if VapusSvcPackageManager.PlatformRBACManager == nil {
			log.Println(params.PbacConfigPath)
			if params.PbacConfigPath == "" {
				logger.Err(ErrPbacConfigPathEmpty).Msg("Error while loading pbac config")
				return nil, nil, ErrPbacConfigPathEmpty
			}
			VapusSvcPackageManager.PlatformRBACManager, err = pbac.LoadPbacConfig(params.PbacConfigPath)
			if err != nil {
				logger.Err(err).Msg("Error while initializing pbac config")
				return params, VapusSvcPackageManager, ErrPbacConfigInitFailed
			}
			logger.Info().Msg("PBAC config initialized successfully")
		}
	}
	VapusSvcPackageManager.ValidEnums = appconfigs.GetValidEnums()
	k8sConfig, err := k8s.GetK8sClusteAPI(logger, nil, "")
	if err != nil {
		logger.Err(err).Msg("Error while getting k8s config")
		// return nil, nil, err
	}
	if k8sConfig != nil {
		VapusSvcPackageManager.HostK8SClient, err = k8s.GetK8SClientSet(k8sConfig)
	}
	logger.Info().Msg("Service packages initialized successfully")
	log.Println(VapusSvcPackageManager)
	return params, VapusSvcPackageManager, nil
}
