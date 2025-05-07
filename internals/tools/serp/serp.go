package serp

import (
	"net"
	"net/http"
	"time"

	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

type Options func(*SerpClient)

type SerpClient struct {
	Params     []*models.Mapper
	creds      *models.PluginNetworkParams
	HttpClient *http.Client
}

func WithParams(params []*models.Mapper) Options {
	return func(c *SerpClient) {
		c.Params = params
	}
}

func WithCreds(creds *models.PluginNetworkParams) Options {
	return func(c *SerpClient) {
		c.creds = creds
	}
}

func New(opts ...Options) *SerpClient {
	serpCl := &SerpClient{}
	for _, opt := range opts {
		opt(serpCl)
	}
	serpCl.HttpClient = &http.Client{Transport: &http.Transport{
		MaxIdleConns:    50,
		IdleConnTimeout: 3600 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second, // Connection timeout
			KeepAlive: 60 * time.Second, // Keep-alive timeout
		}).DialContext,
	},
		Timeout: 60 * 5 * time.Second,
	}
	return serpCl
}
