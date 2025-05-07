package vaws

import (
	"context"
	"net/url"
	"strings"

	ecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtype "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
)

type ECRManager struct {
	Client     *ecr.Client
	Repository []string
	address    string
	RegistryId string
	logger     zerolog.Logger
}

func getRegistryid(address string) string {
	addr := strings.Split(address, ".")
	if len(addr) > 1 {
		return addr[0]
	}
	return ""
}

// NewECRManager creates a new ECRManager instance with the given AWS configuration and address.
// It returns a pointer to the ECRManager and an error, if any.
// The AWS configuration is specified by the opts parameter, which contains the necessary credentials and region information.
// The address parameter is the URL of the ECR registry.
// The function initializes the ECRManager by setting the AWS client configuration and parsing the registry ID and repository from the address.
// If any error occurs during the initialization process, an error is returned.
func NewECRManager(ctx context.Context, opts *AWSConfig, address string, l zerolog.Logger) (*ECRManager, error) {
	configCl, err := opts.getAwsCLientConfig(ctx)
	if err != nil {
		return nil, dmerrors.DMError(ErrAwsConfigLoading, nil)
	}
	obj := &ECRManager{
		Client: ecr.NewFromConfig(configCl),
		logger: l,
	}
	l.Debug().Msg(address)
	parsedAddress, err := url.Parse(address)
	if err != nil {
		l.Err(err).Msg("error while parsing ECR url")
		return nil, dmerrors.DMError(ErrParsingEcrAddress, err)
	}
	l.Debug().Msg(parsedAddress.Host)
	obj.RegistryId = getRegistryid(parsedAddress.Host)
	if obj.RegistryId == "" {
		return nil, dmerrors.DMError(ErrParsingEcrId, nil)
	}
	if parsedAddress.Path != "" {
		obj.Repository = []string{parsedAddress.Path}
	}
	return obj, nil
}

// GetAuthToken retrieves the authorization token for accessing the Amazon Elastic Container Registry (ECR).
// It returns the authorization token as a string and an error if any.
func (ecrM *ECRManager) GetAuthToken(ctx context.Context) (string, error) {
	authz := ecr.GetAuthorizationTokenInput{}
	authzResp, err := ecrM.Client.GetAuthorizationToken(ctx, &authz)
	if err != nil {
		return "", dmerrors.DMError(ErrGettingEcrAuthToken, err)
	}
	if len(authzResp.AuthorizationData) == 0 {
		return "", dmerrors.DMError(ErrGettingEcrAuthToken, ErrECRNoAuthTokenFound)
	}
	return *authzResp.AuthorizationData[0].AuthorizationToken, nil
}

// ListRepoArtifacts retrieves a list of image details for a given repository name.
// It takes a context and the repository name as input parameters.
// It returns a slice of ecrtype.ImageDetail and an error.
// The function makes use of the DescribeImages API to fetch the image details.
// If an error occurs during the API call, it returns a nil slice and an error.
// If there are multiple pages of results, it retrieves all the pages and appends the image details to the result slice.
func (ecrM *ECRManager) ListRepoArtifacts(ctx context.Context, repoName string) ([]ecrtype.ImageDetail, error) {
	req := ecr.DescribeImagesInput{
		RepositoryName: &repoName,
		RegistryId:     &ecrM.RegistryId,
	}
	var result []ecrtype.ImageDetail

	resp, err := ecrM.Client.DescribeImages(ctx, &req)
	if err != nil {
		return nil, dmerrors.DMError(ErrListingECRPackages, err)
	}
	for {
		if err != nil {
			return nil, dmerrors.DMError(ErrListingECRPackages, err)
		}
		result = append(result, resp.ImageDetails...)
		if resp.NextToken == nil {
			break
		}
		req.NextToken = resp.NextToken
		resp, err = ecrM.Client.DescribeImages(ctx, &req)
	}
	return result, nil
}
