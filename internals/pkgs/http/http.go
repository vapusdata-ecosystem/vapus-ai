package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type RestHttp struct {
	address       string
	basePath      string
	httpClient    http.Client
	bearerAuth    string
	customHeaders map[string]string
	urlParsed     *url.URL
	logger        zerolog.Logger
}
type RestHttpOpts func(*RestHttp)

func WithAddress(address string) RestHttpOpts {
	return func(r *RestHttp) {
		r.address = address
	}
}

func WithBasePath(basePath string) RestHttpOpts {
	return func(r *RestHttp) {
		r.basePath = basePath
	}
}

func WithBearerAuth(bearerAuth string) RestHttpOpts {
	return func(r *RestHttp) {
		r.bearerAuth = bearerAuth
	}
}

func WithCustomHeaders(h map[string]string) RestHttpOpts {
	return func(r *RestHttp) {
		r.customHeaders = h
	}
}

func New(logger zerolog.Logger, opts ...RestHttpOpts) (*RestHttp, error) {
	r := &RestHttp{}
	for _, opt := range opts {
		opt(r)
	}
	httpClient := http.Client{Transport: &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 600 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Connection timeout
			KeepAlive: 30 * time.Second, // Keep-alive timeout
		}).DialContext,
	},
		Timeout: 60 * 5 * time.Second,
	}
	r.logger = logger
	r.httpClient = httpClient
	r.address = r.validateHttpProtocol(r.address)
	nUrl, err := url.Parse(r.address)
	if err != nil {
		logger.Error().Err(err).Msgf("Error parsing url - %s", r.address)
		return nil, err
	}
	r.urlParsed = nUrl
	return r, nil
}

func (r *RestHttp) validateHttpProtocol(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	} else {
		return "http://" + url
	}
}

func (r *RestHttp) SetBearerToken(token string) {
	r.bearerAuth = token
}

func (r *RestHttp) SetHeader(key, val string) {
	log.Println(key, val)
}

func (r *RestHttp) handleResponse(resp *http.Response, response any) error {
	fmt.Println("I am handling Response =========>>>>>>>>>>>>>>>")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error reading response")
		return err
	}
	log.Println(string(body))
	err = json.Unmarshal(body, response)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error unmarshalling response")
		return err
	}
	return nil
}

func (r *RestHttp) Get(ctx context.Context, subPath string, qParams map[string]string, response any, contentType string) error {
	api := url.URL{
		Scheme: r.urlParsed.Scheme,
		Host:   r.urlParsed.Host,
		Path:   r.urlParsed.Path,
	}
	api.Path = path.Join(r.basePath, subPath)
	fmt.Println("I am in get get request          -------------------->>>>>>>>>>.")
	fmt.Println(api.Path)
	query := api.Query()
	for key, value := range qParams {
		query.Add(key, value)
	}
	api.RawQuery = query.Encode()
	req, err := http.NewRequest(types.GET, api.String(), bytes.NewBuffer(nil))
	if err != nil {
		r.logger.Error().Err(err).Msg("Error creating request")
		return err
	}
	for key, value := range r.customHeaders {
		req.Header.Add(key, value)
	}
	req.Header.Set("Content-Type", contentType)
	if r.bearerAuth != "" {
		req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	}
	log.Println(api.String())
	log.Println(req.Header)
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error sending request")
		return err
	}
	defer resp.Body.Close()
	return r.handleResponse(resp, response)
}

func (r *RestHttp) Post(ctx context.Context, subPath string, inputPayload []byte, response any, contentType string) error {
	log.Println(string(inputPayload), "====================+++++++++++++++++++++++++++++++++POST")
	api := url.URL{
		Scheme: r.urlParsed.Scheme,
		Host:   r.urlParsed.Host,
		Path:   r.urlParsed.Path,
	}
	api.Path = path.Join(r.basePath, subPath)
	fmt.Println("Path: ----->", api.Path)
	req, err := http.NewRequest(types.POST, api.String(), bytes.NewBuffer(inputPayload))
	if err != nil {
		r.logger.Error().Err(err).Msg("Error creating request for POST")
		return err
	}
	for key, value := range r.customHeaders {
		req.Header.Add(key, value)
	}
	if r.bearerAuth != "" {
		req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	}
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error sending request")
		return err
	}
	defer resp.Body.Close()
	return r.handleResponse(resp, response)
}

func (r *RestHttp) StreamPost(ctx context.Context, subPath string, inputPayload []byte, contentType string) (*http.Response, error) {
	api := url.URL{
		Scheme: r.urlParsed.Scheme,
		Host:   r.urlParsed.Host,
		Path:   r.urlParsed.Path,
	}
	api.Path = path.Join(r.basePath, subPath)
	log.Println("||||||||||||||||||||||||||||||||---1", dmutils.GetEpochTime())
	req, err := http.NewRequest(types.POST, api.String(), bytes.NewBuffer(inputPayload))
	if err != nil {
		r.logger.Error().Err(err).Msg("Error creating request for stream POST")
		return nil, err
	}
	for key, value := range r.customHeaders {
		req.Header.Add(key, value)
	}
	if r.bearerAuth != "" {
		req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	}
	log.Println("||||||||||||||||||||||||||||||||---2", dmutils.GetEpochTime())
	req.Header.Set("Content-Type", contentType)
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error sending stream request")
		return nil, err
	}
	log.Println("||||||||||||||||||||||||||||||||---3", dmutils.GetEpochTime())
	return resp, nil
}

func (r *RestHttp) StreamGet(ctx context.Context, subPath string, qParams map[string]string, contentType string) (*http.Response, error) {
	api := url.URL{
		Scheme: r.urlParsed.Scheme,
		Host:   r.urlParsed.Host,
		Path:   r.urlParsed.Path,
	}
	api.Path = path.Join(r.basePath, subPath)
	query := api.Query()
	for key, value := range qParams {
		query.Add(key, value)
	}
	api.RawQuery = query.Encode()
	req, err := http.NewRequest(types.POST, api.String(), bytes.NewBuffer(nil))
	if err != nil {
		r.logger.Error().Err(err).Msg("Error creating request for stream POST")
		return nil, err
	}
	for key, value := range r.customHeaders {
		req.Header.Add(key, value)
	}
	if r.bearerAuth != "" {
		req.Header.Set("Authorization", "Bearer "+r.bearerAuth)
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Error().Err(err).Msg("Error sending stream request")
		return nil, err
	}
	return resp, nil
}
