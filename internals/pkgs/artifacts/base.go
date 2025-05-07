package artifacts

import (
	"context"

	"github.com/rs/zerolog"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type NabhikArtifactAgent struct {
	logger    zerolog.Logger
	Error     error
	mountPath string
}

type ArtifactOpts struct {
	ArtifactURL  string
	MetaData     map[string]any
	Username     string
	Password     string
	AccesToken   string
	RefreshToken string
}

type opts func(*NabhikArtifactAgent)

func WithRegistryURL(l zerolog.Logger) opts {
	return func(n *NabhikArtifactAgent) {
		n.logger = l
	}
}

func New(options ...opts) *NabhikArtifactAgent {
	n := &NabhikArtifactAgent{}
	for _, option := range options {
		option(n)
	}
	return n
}

func (x *ArtifactOpts) GetOrasCred(ctx context.Context) func(ctx context.Context, hostport string) (auth.Credential, error) {
	if x.AccesToken != "" {
		return func(ctx context.Context, hostport string) (auth.Credential, error) {
			return auth.Credential{
				AccessToken: x.AccesToken,
			}, nil
		}
	}
	if x.RefreshToken != "" {
		return func(ctx context.Context, hostport string) (auth.Credential, error) {
			return auth.Credential{
				RefreshToken: x.RefreshToken,
			}, nil
		}
	}
	return func(ctx context.Context, hostport string) (auth.Credential, error) {
		return auth.Credential{
			Username: x.Username,
			Password: x.Password,
		}, nil
	}
}
