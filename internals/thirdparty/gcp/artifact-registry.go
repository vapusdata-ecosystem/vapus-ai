package gcp

import (
	"context"
	"net/url"
	"strings"

	ar "cloud.google.com/go/artifactregistry/apiv1"
	arPb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	option "google.golang.org/api/option"
)

type GcpArManager struct {
	Cl             *ar.Client
	address        string
	repository     string
	opts           *GcpConfig
	repositoryType string
	parent         string
	logger         zerolog.Logger
}

func NewGcpArManager(ctx context.Context, opts *GcpConfig, address string, l zerolog.Logger) (*GcpArManager, error) {
	creds, err := google.CredentialsFromJSON(ctx, opts.ServiceAccountKey)
	if err != nil {
		// TODO: handle error.
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, dmerrors.DMError(ErrCreatingGcpArClient, err)
	}
	cli, err := ar.NewClient(context.TODO(), option.WithCredentials(creds))

	if err != nil {
		return nil, dmerrors.DMError(ErrCreatingGcpArClient, err)
	}

	obj := &GcpArManager{
		Cl:      cli,
		opts:    opts,
		address: address,
		logger:  l,
	}

	err = obj.Render()
	if err != nil {
		logger.Err(err).Msg("Error while rendering GCP details")
		return nil, err
	}

	return obj, nil
}

// Render renders the GcpArManager object by parsing the dnNode and setting the repositoryType, opts.Region, opts.ProjectID, and repository fields.
// It returns an error if there is an issue parsing the dnNode or if the parsing results in an invalid format.
func (gcn *GcpArManager) Render() error {
	// Use regex for below rendering
	x := strings.Contains(gcn.address, "https://") || strings.Contains(gcn.address, "http://")
	if !x {
		gcn.address = "https://" + gcn.address
	}
	parsed, err := url.Parse(gcn.address)
	if err != nil {
		return dmerrors.DMError(ErrParsingGAR, err)
	}
	hostSplit := strings.Split(parsed.Host, ".")

	if len(hostSplit) < 3 {
		return dmerrors.DMError(ErrParsingGARHost, err)
	}

	hostTemp := strings.Split(hostSplit[0], "-")

	if len(hostTemp) < 2 {
		return dmerrors.DMError(ErrParsingGARRegion, err)
	}

	gcn.repositoryType = hostTemp[len(hostTemp)-1]

	gcn.opts.Region = strings.Replace(hostSplit[0], "-"+gcn.repositoryType, "", 1)

	paths := strings.Split(parsed.Path, "/")
	gcn.opts.ProjectID = paths[1]
	gcn.repository = paths[2]
	gcn.parent = "projects/" + gcn.opts.ProjectID + "/locations/" + gcn.opts.Region + "/repositories/" + gcn.repository
	logger.Debug().Msgf("Full GCN --- %v", gcn)
	return nil

}

func (gcn *GcpArManager) Close() error {
	return gcn.Cl.Close()
}

func (gcn *GcpArManager) ListRepoOciArtifacts(ctx context.Context) ([]*arPb.DockerImage, error) {
	var result []*arPb.DockerImage
	req := &arPb.ListDockerImagesRequest{
		Parent:  gcn.parent,
		OrderBy: "NAME",
	}
	res := gcn.Cl.ListDockerImages(context.Background(), req)

	if res == nil {
		return nil, dmerrors.DMError(ErrListingGARPackages, nil)
	}
	defer gcn.Close()
	for {
		resp, err := res.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Err(err).Msgf("Error while listing docker images from gar - %v", gcn.parent)
		}
		result = append(result, resp)

	}
	return result, nil
}
