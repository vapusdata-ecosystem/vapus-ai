package serp

import (
	"encoding/json"

	serp "github.com/serpapi/google-search-results-golang"
	"github.com/vapusdata-ecosystem/vapusai/core/options"
)

func (s *SerpClient) newSerpCaller(opts *options.SearchInput) *serp.Search {
	if opts == nil {
		opts = &options.SearchInput{}
	}
	if opts.Engine == "" {
		opts.Engine = SerpEngineGoogle.String()
	}
	return &serp.Search{
		HttpSearch: s.HttpClient,
		ApiKey:     s.creds.Credentials.ApiToken,
		Engine:     opts.Engine,
		Parameter:  opts.Params,
	}
}

// func (s *SerpClient) getJson(opts *SerpOpts) (serp.SearchResult, error) {
// 	client := s.NewSerpCaller(opts)
// 	return client.GetJSON()
// }

// func (s *SerpClient) getHTML(opts *SerpOpts) (*string, error) {
// 	client := s.NewSerpCaller(opts)
// 	return client.GetHTML()
// }

// func (s *SerpClient) getSearchArchive(opts *SerpOpts) (serp.SearchResult, error) {
// 	client := s.NewSerpCaller(opts)
// 	return client.GetSearchArchive(opts.SearchId)
// }

// func (s *SerpClient) getLocation(opts *SerpOpts) (serp.SearchResultArray, error) {
// 	client := s.NewSerpCaller(opts)
// 	return client.GetLocation(opts.LocationQuery, int(opts.LocationLimit))
// }

// func (s *SerpClient) getAccount(opts *SerpOpts) (serp.SearchResult, error) {
// 	client := s.NewSerpCaller(opts)
// 	return client.GetAccount()
// }

func (s *SerpClient) SearchFormatted(opts *options.SearchInput) (*options.SearchResult, error) {
	client := s.newSerpCaller(opts)
	resp, err := client.GetJSON()
	if err != nil || resp == nil {
		return nil, err
	}
	return s.formatSearchResult(resp, opts)
}

func (s *SerpClient) SearchRaw(opts *options.SearchInput) (map[string]any, error) {
	client := s.newSerpCaller(opts)
	return client.GetJSON()
}

func (s *SerpClient) formatSearchResult(resp serp.SearchResult, opts *options.SearchInput) (*options.SearchResult, error) {
	bbytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	switch opts.Engine {
	case SerpEngineGoogle.String():
		googleResp := &options.GoogleSearchResult{}
		err = json.Unmarshal(bbytes, googleResp)
		if err != nil {
			return nil, err
		}
		return &options.SearchResult{
			GoogleSearchResult: googleResp,
		}, nil
	default:
		return nil, nil
	}
}
