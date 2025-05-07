package models

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type K8SInfraParams struct {
	ServiceProvider  string                  `json:"serviceProvider" yaml:"serviceProvider"`
	InfraService     string                  `json:"infraService" yaml:"infraService"`
	Name             string                  `json:"name" yaml:"name"`
	Credentials      *GenericCredentialModel `json:"credentials" yaml:"credentials"`
	DisplayName      string                  `json:"displayName" yaml:"displayName"`
	InfraId          string                  `json:"infraId" yaml:"infraId"`
	KubeConfig       string                  `json:"kubeConfig" yaml:"kubeConfig"`
	SecretName       string                  `json:"secretName" yaml:"secretName"`
	IsDefault        bool                    `json:"isDefault" yaml:"isDefault"`
	KubeConfigFormat string                  `json:"kubeConfigFormat" yaml:"kubeConfigFormat"`
}

func (k *K8SInfraParams) GetServiceProvider() string {
	return k.ServiceProvider
}

func (k *K8SInfraParams) GetInfraService() string {
	return k.InfraService
}

func (k *K8SInfraParams) GetName() string {
	return k.Name
}

func (k *K8SInfraParams) GetCredentials() *GenericCredentialModel {
	return k.Credentials
}

func (k *K8SInfraParams) GetDisplayName() string {
	return k.DisplayName
}

func (k *K8SInfraParams) GetInfraId() string {
	return k.InfraId
}

func (k *K8SInfraParams) GetKubeConfig() string {
	return k.KubeConfig
}

func (k *K8SInfraParams) GetSecretName() string {
	return k.SecretName
}

func (k *K8SInfraParams) GetIsDefault() bool {
	return k.IsDefault
}

func (k *K8SInfraParams) GetKubeConfigFormat() string {
	return k.KubeConfigFormat
}

func (k *K8SInfraParams) ConvertToPb() *mpb.K8SInfraParams {
	if k != nil {
		return &mpb.K8SInfraParams{
			ServiceProvider:  mpb.ServiceProvider(mpb.ServiceProvider_value[k.ServiceProvider]),
			InfraService:     mpb.InfraService(mpb.InfraService_value[k.InfraService]),
			Name:             k.Name,
			Credentials:      k.Credentials.ConvertToPb(),
			DisplayName:      k.DisplayName,
			InfraId:          k.InfraId,
			KubeConfig:       k.KubeConfig,
			SecretName:       k.SecretName,
			IsDefault:        k.IsDefault,
			KubeConfigFormat: mpb.ContentFormats(mpb.ContentFormats_value[k.KubeConfigFormat]),
		}
	}
	return nil
}

func (k *K8SInfraParams) ConvertFromPb(pb *mpb.K8SInfraParams) *K8SInfraParams {
	if pb != nil {
		return &K8SInfraParams{
			ServiceProvider:  pb.GetServiceProvider().String(),
			InfraService:     pb.GetInfraService().String(),
			Name:             pb.GetName(),
			Credentials:      (&GenericCredentialModel{}).ConvertFromPb(pb.GetCredentials()),
			DisplayName:      pb.GetDisplayName(),
			InfraId:          pb.GetInfraId(),
			KubeConfig:       pb.GetKubeConfig(),
			SecretName:       pb.GetSecretName(),
			IsDefault:        pb.GetIsDefault(),
			KubeConfigFormat: pb.GetKubeConfigFormat().String(),
		}
	}
	return nil
}
