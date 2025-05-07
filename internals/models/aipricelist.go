package models

type AIModelPriceList struct {
	VapusBase              `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	LLMServiceProviderName string          `bun:"llm_service_provider_name,notnull,unique" json:"llm_service_provider_name,omitempty" yaml:"llm_service_provider_name,omitempty"`
	Currency               string          `bun:"currency" json:"currency,omitempty" yaml:"currency,omitempty"`
	ModelPrices            []*ModelPricing `bun:"model_prices,type:jsonb" json:"model_prices,omitempty" yaml:"model_prices,omitempty"`
}

// func (m *AIModelPriceList) SetAccountId(accountId string) {
// 	if m != nil {
// 		m.OwnerAccount = accountId
// 	}
// }

// func (d *AIModelPriceList) SetPromptId() {
// 	if d == nil {
// 		return
// 	}
// 	d.VapusID = fmt.Sprintf(types.PROMPT_ID, guuid.New())
// }

// func (d *AIModelPriceList) PreSaveCreate(authzClaim map[string]string) {
// 	if d == nil {
// 		return
// 	}
// 	d.PreSaveVapusBase(authzClaim)
// }

// func (dn *AIModelPriceList) PreSaveUpdate(userId string) {
// 	if dn == nil {
// 		return
// 	}
// 	dn.UpdatedBy = userId
// 	dn.UpdatedAt = dmutils.GetEpochTime()
// }

//	func (dn *AIModelPriceList) PreSaveDelete(authzClaim map[string]string) {
//		if dn == nil {
//			return
//		}
//		dn.PreDeleteVapusBase(authzClaim)
//	}
//
// Context int32

type ModelPricing struct {
	ModelName                  string         `json:"modelName,omitempty" yaml:"modelName,omitempty"`
	Context                    int32          `json:"context,omitempty" yaml:"context,omitempty"`
	PromptCount                int32          `json:"promptCount,omitempty" yaml:"promptCount,omitempty"`
	PromptUpTo                 int32          `json:"promptUpTo,omitempty" yaml:"promptUpTo,omitempty"`
	PromptLongerThan           int32          `json:"promptLongerThan,omitempty" yaml:"promptLongerThan,omitempty"`
	InputPrice                 *PriceAndToken `json:"inputPrice,omitempty" yaml:"inputPrice,omitempty"`
	CachedInputPrice           *PriceAndToken `json:"cachedInputPrice,omitempty" yaml:"cachedInputPrice,omitempty"`
	CachedOutputPrice          *PriceAndToken `json:"cachedOutputPrice,omitempty" yaml:"cachedOutputPrice,omitempty"`
	ContextCachingPrice        *PriceAndToken `json:"contextCachingPrice,omitempty" yaml:"contextCachingPrice,omitempty"`
	ContextCachingStoragePrice *PriceAndToken `json:"contextCachingStoragePrice,omitempty" yaml:"contextCachingStoragePrice,omitempty"`
	OutputPrice                *PriceAndToken `json:"outputPrice,omitempty" yaml:"outputPrice,omitempty"`
	IsPriceListUpdated         bool           `json:"isPriceListUpdated,omitempty" yaml:"isPriceListUpdated,omitempty"`
	// PricePer1000Request        int32          `json:"pricePer1000Request,omitempty" yaml:"pricePer1000Request,omitempty"`
}

type PriceAndToken struct {
	Amount    float64 `json:"amount,omitempty" yaml:"amount,omitempty"`
	TokenSize int64   `json:"tokenSize,omitempty" yaml:"tokenSize,omitempty"`
}
