package serp

type SerpEngines string

func (s SerpEngines) String() string {
	return string(s)
}

const (
	SerpEngineGoogle        SerpEngines = "google"
	SerpEngineBing          SerpEngines = "bing"
	SerpEngineBaidu         SerpEngines = "baidu"
	SerpEngineYahoo         SerpEngines = "yahoo"
	SerpEngineYandex        SerpEngines = "yandex"
	SerpEngineEbay          SerpEngines = "ebay"
	SerpEngineYoutube       SerpEngines = "youtube"
	SerpEngineWalmart       SerpEngines = "walmart"
	SerpEngineHomeDepot     SerpEngines = "home_depot"
	SerpEngineGoogleMaps    SerpEngines = "google_maps"
	SerpEngineGoogleProduct SerpEngines = "google_product"
	SerpEngineGoogleScholar SerpEngines = "google_scholar"
	SerpEngineGoogleNews    SerpEngines = "google_news"
	SerpEngineGoogleFinance SerpEngines = "google_finance"
	SerpEngineAppleAppStore SerpEngines = "apple_app_store"
)
