package crane

import (
	"context"
	"os"
	"strings"

	authn "github.com/google/go-containerregistry/pkg/authn"
	crane "github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/rs/zerolog"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
)

type CraneClient struct {
	authOpts  *CraneAuth
	logger    zerolog.Logger
	mountPath string
}

// CraneAuth represents the authentication details for the Crane service.
type CraneAuth struct {
	ImageURL  string // ImageURL is the URL of the image.
	Username  string // Username is the username for authentication.
	Password  string // Password is the password for authentication.
	AuthToken string // AuthToken is the authentication token for authentication.
	Address   string // Address is the address of the Crane service.
	Tag       string // Tag is the tag of the image.
}

type CraneOpts func(*CraneClient)

func WithLogger(l zerolog.Logger) CraneOpts {
	return func(c *CraneClient) {
		c.logger = l
	}
}

func WithMountPath(m string) CraneOpts {
	return func(c *CraneClient) {
		c.mountPath = m
	}
}

func WithAuthOpts(a *CraneAuth) CraneOpts {
	return func(c *CraneClient) {
		c.authOpts = a
	}
}

// authOptions returns a crane.Option that sets the authentication token for the crane client.
func authOptions(opts CraneAuth) crane.Option {
	if opts.Username != "" && opts.Password != "" {
		return crane.WithAuth(
			authn.FromConfig(authn.AuthConfig{
				Username: opts.Username,
				Password: opts.Password,
			},
			),
		)
	}
	return crane.WithAuth(
		authn.FromConfig(authn.AuthConfig{
			Auth: opts.AuthToken,
		},
		),
	)
}

func NewCraneClient(ctx context.Context, opts ...CraneOpts) *CraneClient {
	cl := &CraneClient{}
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

// GetCatalog retrieves the list of available images in the catalog.
// It takes CraneAuth options as input and returns a slice of strings representing the image names and an error, if any.
func GetCatalog(opts CraneAuth) ([]string, error) {
	address := strings.Replace(opts.Address, "https://", "", 1)
	res, err := crane.Catalog(address, authOptions(opts))
	if err != nil {
		return []string{}, err
	}
	return res, nil
}

func (c *CraneClient) GetCatalog() ([]string, error) {
	return GetCatalog(*c.authOpts)
}
func (c *CraneClient) PullImage() error {
	_, err := crane.Pull(c.authOpts.ImageURL, authOptions(*c.authOpts))
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error pulling image %s", c.authOpts.ImageURL)
		return err
	}
	return nil
}

func (c *CraneClient) PushImage(img v1.Image) (string, error) {
	err := crane.Push(img, c.authOpts.ImageURL, authOptions(*c.authOpts))
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error pushing image %s", c.authOpts.ImageURL)
		return "", err
	}
	dg, err := img.Digest()
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error getting digest of image %s", c.authOpts.ImageURL)
		return "", err
	}
	return dg.String(), nil
}

func (c *CraneClient) GetFullOciURL() string {
	return c.authOpts.ImageURL
}

func (c *CraneClient) AppendFiles2OCIAndPush(destination *CraneAuth, files []string, identifier string) (string, string, error) {
	var err error
	img, err := crane.Pull(c.authOpts.ImageURL, authOptions(*c.authOpts))
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error pulling image %s", c.authOpts.ImageURL)
		return "", "", err
	}
	oldManifest, err := img.Manifest()
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error getting new manifest of image %s", c.authOpts.ImageURL)
	}
	c.logger.Info().Msgf("New manifest of image %s - %v", c.authOpts.ImageURL, oldManifest)
	oldLayers, err := img.Layers()
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error getting layers of image %s", c.authOpts.ImageURL)
	}
	c.logger.Info().Msgf("New Layers of image %s - %v", c.authOpts.ImageURL, oldLayers)
	tarFile := identifier + ".tar"
	err = filetools.CreateTarFile(tarFile, files, c.mountPath)
	// err = dmutils.CreateTarFile(tarFile, files, "/mnt")
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error creating tar file %s", tarFile)
		return "", "", err
	}
	defer os.Remove(tarFile)
	layer, err := tarball.LayerFromFile(tarFile)
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error creating layer from file %s", tarFile)
		return "", "", err
	}
	c.logger.Info().Msgf("Successfully created layer from file %s", tarFile)
	newImg, err := mutate.AppendLayers(img, layer)
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error appending files to image %s", c.authOpts.ImageURL)
	}
	newManifest, err := newImg.Manifest()
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error getting new manifest of image %s", c.authOpts.ImageURL)
	}
	c.logger.Info().Msgf("New manifest of image %s - %v", c.authOpts.ImageURL, newManifest)
	layers, err := newImg.Layers()
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error getting layers of image %s", c.authOpts.ImageURL)
	}
	c.logger.Info().Msgf("New Layers of image %s - %v", c.authOpts.ImageURL, layers)
	digest, err := newImg.Digest()
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error getting digest of image %s", destination.ImageURL)
		return "", "", err
	}
	err = crane.Push(newImg, destination.ImageURL, authOptions(*destination))
	if err != nil {
		c.logger.Error().Err(err).Msgf("Error pushing image %s", destination.ImageURL)
		return "", "", err
	}
	c.logger.Info().Msgf("Successfully pushed image %s to %s with digest - %v", c.authOpts.ImageURL, destination.ImageURL, digest.String())
	return destination.ImageURL, digest.String(), nil
}

func (c *CraneClient) CopyOCI(destination *CraneAuth) error {
	return crane.Copy(c.authOpts.ImageURL, destination.ImageURL, authOptions(*c.authOpts), authOptions(*destination))
}
