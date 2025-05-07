package githubconnector

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"golang.org/x/oauth2"
)

type GithubOpts struct {
	// PostgresConfig is the configuration for the Postgres
	URL, Username, Password, Database, Schema, Pat string
	Port                                           int
	WithPool                                       bool
}

type GithubStore struct {
	Opts *GithubOpts
	Conn *github.Client
	User string
}

func New(opts *GithubOpts, l zerolog.Logger) (*GithubStore, error) {
	token := opts.Pat
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrGithubConnection, err)
	}
	l.Info().Msgf("Authenticated user %v", user.GetLogin())
	//l.Println("Authenticated user %v",user.GetName())
	return &GithubStore{
		Opts: opts,
		Conn: client,
		User: user.GetLogin(),
	}, nil
}

func (m *GithubStore) Close() {
	m.Conn = nil
}
