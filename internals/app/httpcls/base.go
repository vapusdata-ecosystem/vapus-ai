package httpcls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
)

type VapusHttpClient struct {
	Address string
	logger  zerolog.Logger
	Version string
}

type httpRequestGeneric struct {
	uri    string
	method string
	body   []byte
	token  string
}

type HttpRequestGeneric struct {
	Uri    string
	Method string
	Body   []byte
	Token  string
}

func New(address string, logger zerolog.Logger) (*VapusHttpClient, error) {
	return &VapusHttpClient{
		Address: address,
		logger:  logger,
	}, nil
}

func (x *VapusHttpClient) httpClient(httpOpts *httpRequestGeneric) ([]byte, error) {
	url, err := url.Parse(x.Address)
	url.Path = httpOpts.uri

	req, err := http.NewRequest(httpOpts.method, url.String(), bytes.NewBuffer(httpOpts.body))
	if err != nil {
		x.logger.Error().Err(err).Msg("Error creating request")
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+httpOpts.token)
	// req.Header.Add("Authorization", httpOpts.token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		x.logger.Error().Err(err).Msg("Error sending request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		x.logger.Error().Err(err).Msg("Error reading response")
		return nil, err
	}

	return body, nil
}

func (x *VapusHttpClient) VapusHttpClient(httpOpts *HttpRequestGeneric) ([]byte, error) {
	urls, err := url.Parse(x.Address)
	urls.Path = httpOpts.Uri
	var api string
	// payload := strings.NewReader(``)
	if urls.Scheme != "https" {
		urls.Scheme = "http"
		api = "http://" + x.Address + httpOpts.Uri
	} else {
		api = x.Address + httpOpts.Uri
	}
	client := &http.Client{}
	var bd io.Reader
	if httpOpts.Body == nil {
		bd = strings.NewReader(``)
	} else {
		bd = strings.NewReader(fmt.Sprintf(`%s`, string(httpOpts.Body)))
	}
	req, err := http.NewRequest(httpOpts.Method, api, bd)

	if err != nil {
		x.logger.Error().Err(err).Msg("Error creating request")
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+httpOpts.Token)

	res, err := client.Do(req)
	if err != nil {
		x.logger.Error().Err(err).Msg("Error sending request")
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		x.logger.Error().Err(err).Msg("Error reading response")
		return nil, err
	}
	if !slices.Contains([]int{200, 201, 202}, res.StatusCode) {
		errBody := map[string]any{}
		err := json.Unmarshal(body, &errBody)
		if err != nil {
			return nil, fmt.Errorf("Error: %v", res.Status)
		}
		errBody["statusCode"] = res.StatusCode
		bbytes, err := json.Marshal(errBody)
		if err != nil {
			return body, fmt.Errorf(" %v", errBody)
		}
		return body, fmt.Errorf(" %v", string(bbytes))
	}
	return body, nil
}

func ConvertSearchParamToUrlValues(searchParam *mpb.SearchParam, urlParams url.Values) url.Values {
	if searchParam == nil {
		return urlParams
	}
	for i, value := range searchParam.Filters {
		if value == nil {
			continue
		}
		urlParams.Add("search_param.filters["+strconv.Itoa(i)+"].key", value.Key)
		urlParams.Add("search_param.filters["+strconv.Itoa(i)+"].value", value.GetValue())
	}
	if searchParam.Q != "" {
		urlParams.Add("search_param.q", searchParam.Q)
	}
	return urlParams
}

func appendIdToUrl(url string, id string) string {
	return fmt.Sprintf("%s/%s", url, id)
}

type GenericParams interface {
	*pb.AuthzGetterRequest | *pb.UserGetterRequest | *pb.OrganizationGetterRequest |
		*pb.DataSourceGetterRequest |
		*aipb.AIModelNodeGetterRequest | *aipb.PromptGetterRequest |
		*aipb.GuardrailsGetterRequest | *aipb.AgentGetterRequest
}
