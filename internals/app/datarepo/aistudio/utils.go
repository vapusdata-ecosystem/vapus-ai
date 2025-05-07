package aidmstore

import (
	"context"

	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func (ds *AIStudioDMStore) GetAIModelNodeNetworkParams(ctx context.Context, aiModelNode *models.AIModelNode) (*models.AIModelNodeNetworkParams, error) {
	secrets, err := apppkgs.ReadCredentialFromStore(ctx, aiModelNode.NetworkParams.SecretName, ds.VapusStore, ds.logger)
	if err != nil {
		ds.logger.Err(err).Msg("error while reading credentials from store")
		return nil, err
	}

	return &models.AIModelNodeNetworkParams{
		Url:                 aiModelNode.NetworkParams.GetUrl(),
		Credentials:         secrets,
		ApiVersion:          aiModelNode.NetworkParams.GetApiVersion(),
		LocalPath:           aiModelNode.NetworkParams.GetLocalPath(),
		SecretName:          aiModelNode.NetworkParams.SecretName,
		IsAlreadyInSecretBs: aiModelNode.NetworkParams.IsAlreadyInSecretBs,
	}, nil
}
