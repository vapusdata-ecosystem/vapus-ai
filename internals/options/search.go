package options

import "github.com/rs/zerolog"

type SearchInput struct {
	Engine        string            `json:"engine"`
	Params        map[string]string `json:"params"`
	Logger        zerolog.Logger    `json:"logger"`
	Q             string            `json:"q"`
	IsSafe        bool              `json:"isSafe"`
	LocationQuery string            `json:"locationQuery"`
	LocationLimit int               `json:"locationLimit"`
	SearchId      string            `json:"searchId"`
	SearchRaw     bool              `json:"searchRaw"`
}

type SearchResult struct {
	*GoogleSearchResult
}

type GoogleSearchResult struct {
	AnswerBox      map[string]any `json:"answerBox"`
	OrganicResults []struct {
		Position                int      `json:"position"`
		Title                   string   `json:"title"`
		Link                    string   `json:"link"`
		RedirectLink            string   `json:"redirectLink"`
		DisplayedLink           string   `json:"displayedLink"`
		Favicon                 string   `json:"favicon"`
		Snippet                 string   `json:"snippet"`
		SnippetHighlightedWords []string `json:"snippetHighlightedWords"`
		Source                  string   `json:"source"`
	}
}
