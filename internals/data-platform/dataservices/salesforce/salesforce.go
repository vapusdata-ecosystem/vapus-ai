package salesforce

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
	pkghttp "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/http"
)

const (
	// Salesforce login url
	OAuthLoginUrl = "https://login.salesforce.com/services/oauth2/token"
)

var (
	ListObjectsUrl  = "/services/data/%s/sobjects/"
	DescribeObjects = "/services/data/%s/sobjects/%s/describe"
)

type SalesforceOpts struct {
	ClientId      string
	ClientSecret  string
	Username      string
	Password      string
	JwtToken      string
	LoginUrl      string
	JwtPrivateKey string
	APIVersion    string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
}

type Salesforce struct {
	InstanceUrl string
	AccessToken string
	Opts        *SalesforceOpts
	RestClient  *pkghttp.RestHttp
	logger      zerolog.Logger
}

func NewSalesforceAgent(opts *SalesforceOpts, logger zerolog.Logger) (*Salesforce, error) {
	if opts.LoginUrl == "" {
		opts.LoginUrl = OAuthLoginUrl
	}
	// data := fmt.Sprintf(
	// 	"grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
	// 	opts.ClientId, opts.ClientSecret, opts.Username, opts.Password,
	// )
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("grant_type", "client_credentials")
	writer.WriteField("client_id", opts.ClientId)
	writer.WriteField("client_secret", opts.ClientSecret)
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(opts.LoginUrl, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, pkgs.ErrLoginSalesForceInstance
	}

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return nil, err
	}
	client, err := pkghttp.New(logger,
		pkghttp.WithAddress(authResp.InstanceURL),
		pkghttp.WithBearerAuth(authResp.AccessToken),
	)
	if err != nil {
		return nil, err
	}
	return &Salesforce{
		InstanceUrl: authResp.InstanceURL,
		AccessToken: authResp.AccessToken,
		Opts:        opts,
		logger:      logger,
		RestClient:  client,
	}, err
}
