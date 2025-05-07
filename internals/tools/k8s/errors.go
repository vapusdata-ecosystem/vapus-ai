package k8s

import (
	"errors"
)

var (
	ErrK8SInfraParamsNil          = errors.New("k8s infra params is nil")
	ErrUnsupportedServiceProvider = errors.New("unsupported service provider")
	ErrInvalidK8SServiceType      = errors.New("invalid k8s service type")
	ErrK8SCouldNotConnect         = errors.New("k8s client could not connect")
	ErrKubeConfigNotFound         = errors.New("kubeconfig not found")
)
