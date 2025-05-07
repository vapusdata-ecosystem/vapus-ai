package gcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	"k8s.io/client-go/tools/clientcmd/api"
)

func GetGkeKubeConfig(ctx context.Context, opts *GcpConfig, clusterName string, logger zerolog.Logger) (*api.Config, error) {
	client, err := container.NewService(ctx, option.WithCredentialsJSON(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, err
	}
	log.Println(opts.Zone)
	log.Println(clusterName)
	log.Println(opts.GetGcpResource(opts.Zone, "clusters", clusterName), "-------------------??????????????")
	cluster, err := client.Projects.Locations.Clusters.Get(opts.GetGcpResource(opts.Zone, "clusters", clusterName)).Context(ctx).Do()
	if err != nil {
		logger.Err(err).Msgf("Error while getting cluster details -- %v", err)
		return nil, err
	}

	cacert, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		logger.Err(err).Msgf("Error while decoding cluster certificate -- %v", err)
		return nil, err
	}

	accessToken, err := retrieveAccessToken(opts.ServiceAccountKey, logger)
	if err != nil {
		logger.Err(err).Msgf("Error while getting access token from gcp for GKE - %v", err)
		return nil, err
	}

	return &api.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: clusterName,
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server:                   fmt.Sprintf("https://%s", cluster.Endpoint),
				CertificateAuthorityData: cacert,
			},
		},
		Contexts: map[string]*api.Context{
			clusterName: {
				Cluster:  clusterName,
				AuthInfo: clusterName,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			clusterName: {
				Token: accessToken,
			},
		},
	}, nil

}

func retrieveAccessToken(serviceAccountKey []byte, logger zerolog.Logger) (string, error) {
	jwtConfig, err := google.JWTConfigFromJSON(serviceAccountKey, container.CloudPlatformScope)
	if err != nil {
		logger.Err(err).Msgf("Error while creating jwt config in GCP GKE for token source -- %v", err)
		return "", err
	}
	tokenSource := jwtConfig.TokenSource(context.Background())

	token, err := tokenSource.Token()
	if err != nil {
		logger.Err(err).Msgf("Error while getting access token for GKE-- %v", err)
	}
	return token.AccessToken, nil
}
