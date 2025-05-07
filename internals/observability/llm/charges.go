package llm_observability

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	metrics "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio/metrics"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

type PricingStore struct {
	DmStore           *apppkgs.VapusStore
	PriceObj          []*models.AIModelPriceList
	ModelPriceDetails *models.ModelPricing
}

func NewLLMPricingStore(ctx context.Context, dmstore *apppkgs.VapusStore, obj *mpb.ModelObservability, logger zerolog.Logger) (*PricingStore, error) {
	priceList, err := metrics.LLMPriceList(ctx, dmstore, logger)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching models pricing")
		return nil, err
	}

	PricingData := &models.ModelPricing{}
	for _, vals := range priceList {
		// Service Provider Name Check
		// fmt.Println("Provider vals: ", vals.LLMServiceProviderName)
		if vals.LLMServiceProviderName == obj.ModelProvider {
			for _, val := range vals.ModelPrices {
				// Model Name Check
				// fmt.Println("Model Name val: ", val.ModelName)
				if isValidSubstring(strings.ToLower(obj.ModelName), strings.ToLower(val.ModelName)) {
					fmt.Println("Okay I have reached here...")
					fmt.Println("Found the LLM service Provider:", obj.ModelProvider, "Model Name", val.ModelName)
					PricingData = val
					break
				}
			}
			break
		}
	}

	return &PricingStore{
		DmStore:           dmstore,
		PriceObj:          priceList,
		ModelPriceDetails: PricingData,
	}, nil
}

func (priceStore *PricingStore) Calculate(ctx context.Context, usageobj *mpb.ModelObservability, logger zerolog.Logger) *models.AIStudioUsages {
	var inputAmountPerToken float64
	var outputAmountPerToken float64
	var cachedInputAmountPerToken float64
	var cachedOutputAmountPerToken float64
	// log.Println("Reached in Calculated Function")
	if priceStore.ModelPriceDetails == nil {
		return &models.AIStudioUsages{
			Charges: 0,
		}
	}
	if priceStore.ModelPriceDetails.InputPrice != nil {
		inputAmountPerToken = priceStore.ModelPriceDetails.InputPrice.Amount / float64(priceStore.ModelPriceDetails.InputPrice.TokenSize)
		fmt.Println("inputAmountPerToken: ", inputAmountPerToken)
	}
	if priceStore.ModelPriceDetails.OutputPrice != nil {
		outputAmountPerToken = priceStore.ModelPriceDetails.OutputPrice.Amount / float64(priceStore.ModelPriceDetails.OutputPrice.TokenSize)
		fmt.Println("outputAmountPerToken: ", outputAmountPerToken)
	}
	if priceStore.ModelPriceDetails.CachedInputPrice != nil {
		cachedInputAmountPerToken = priceStore.ModelPriceDetails.CachedInputPrice.Amount / float64(priceStore.ModelPriceDetails.CachedInputPrice.TokenSize)
		fmt.Println("cachedInputAmountPerToken: ", cachedInputAmountPerToken)

	}
	if priceStore.ModelPriceDetails.CachedOutputPrice != nil {
		cachedOutputAmountPerToken = priceStore.ModelPriceDetails.CachedOutputPrice.Amount / float64(priceStore.ModelPriceDetails.CachedOutputPrice.TokenSize)
		fmt.Println("cachedOutputAmountPerToken: ", cachedOutputAmountPerToken)

	}
	// if priceStore.ModelPriceDetails.ContextCachingPrice != nil {
	// 	contextCachingAmountPerToken = priceStore.ModelPriceDetails.ContextCachingPrice.Amount / float64(priceStore.ModelPriceDetails.ContextCachingPrice.TokenSize)
	// }
	amount := inputAmountPerToken*float64(usageobj.InputTokens) + outputAmountPerToken*float64(usageobj.OutputTokens) + cachedInputAmountPerToken*float64(usageobj.InputCachedTokens) + cachedOutputAmountPerToken*float64(usageobj.OutputCachedTokens)
	logger.Info().Msgf("Total Charges of using AI '%v'", amount)
	// log.Println("Reached in Calculated Function")
	return &models.AIStudioUsages{
		Charges: amount,
	}
}

func isValidSubstring(mainStr, subStr string) bool {
	if mainStr == subStr {
		return true
	}
	reSuffix := regexp.MustCompile(`-\d+$`)
	mainStrCleaned := reSuffix.ReplaceAllString(mainStr, "")

	mainStrPattern := regexp.QuoteMeta(mainStrCleaned)
	mainStrPattern = strings.ReplaceAll(mainStrPattern, `-`, `[-.]`)

	pattern := `(?i)\b` + mainStrPattern + `\b`
	re := regexp.MustCompile(pattern)

	loc := re.FindStringIndex(subStr)
	if loc == nil {
		return false
	}

	endIndex := loc[1]
	if endIndex < len(subStr) && subStr[endIndex] == '-' {
		return false
	}

	return true
}

func ToUpdate(ctx context.Context, dmstore *apppkgs.VapusStore, serviceProvider string, logger zerolog.Logger) error {
	fileDataInBytes, err := os.ReadFile("/home/ashutosh/workspaces/vapusdata/vapusdata/internals/aistudio/usages/llmModels.json")
	if err != nil {
		logger.Err(err).Msg("Error while reading the JSON file")
		return err
	}

	fileData := map[string]models.AIModelPriceList{}
	err = json.Unmarshal(fileDataInBytes, &fileData)
	if err != nil {
		logger.Err(err).Msg("Error while Unmarshal the JSON file")
		return err
	}
	// fmt.Println("UnMarshal: ", fileData)
	// llmProvider := &serviceProvider
	modelUpdateInfo := fileData[serviceProvider]
	fmt.Println("")
	fmt.Println("LLMServiceProviderName: ", modelUpdateInfo.LLMServiceProviderName)
	fmt.Println("Currency: ", modelUpdateInfo.Currency)
	// fmt.Println("modelUpdateInfo: ", modelUpdateInfo.LLMServiceProviderName)
	for _, val := range modelUpdateInfo.ModelPrices {
		fmt.Println("Model Name", val.ModelName)
	}
	err = metrics.UpdatePriceList(ctx, dmstore, &modelUpdateInfo, logger)
	if err != nil {
		logger.Err(err).Msg("error while getting the updating the model price list")
		return err
	}
	logger.Info().Msgf("Sucessfully Updated the AI Price List")
	return nil
}

func UpdateAllAIModelPrices(ctx context.Context, dmstore *apppkgs.VapusStore, logger zerolog.Logger) error {
	query := `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'OPENAI',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"OPENAI-GPT-o1","inputPrice":{"amount":15,"tokenSize":1000000},"cachedInputPrice":{"amount":7.5,"tokenSize":1000000},"outputPrice":{"amount":60,"tokenSize":1000000}},{"modelName":"OPENAI-GPT-o3-mini","inputPrice":{"amount":1.1,"tokenSize":1000000},"cachedInputPrice":{"amount":0.55,"tokenSize":1000000},"outputPrice":{"amount":4.4,"tokenSize":1000000}},{"modelName":"OPENAI-GPT-4o","inputPrice":{"amount":2.5,"tokenSize":1000000},"cachedInputPrice":{"amount":1.25,"tokenSize":1000000},"outputPrice":{"amount":10,"tokenSize":1000000}},{"modelName":"OPENAI-GPT-4o-mini","inputPrice":{"amount":0.15,"tokenSize":1000000},"cachedInputPrice":{"amount":0.075,"tokenSize":1000000},"outputPrice":{"amount":0.6,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'OPENAI')`

	_, err := dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the OPENAI pricing")
	}

	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'ANTHROPIC',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"Claude-3.5-Sonnet","inputPrice":{"amount":3,"tokenSize":1000000},"cachedInputPrice":{"amount":0.3,"tokenSize":1000000},"cachedOutputPrice":{"amount":3.75,"tokenSize":1000000},"outputPrice":{"amount":15,"tokenSize":1000000}},{"modelName":"Claude-3.5-Haiku","inputPrice":{"amount":0.8,"tokenSize":1000000},"cachedInputPrice":{"amount":0.08,"tokenSize":1000000},"cachedOutputPrice":{"amount":1,"tokenSize":1000000},"outputPrice":{"amount":4,"tokenSize":1000000}},{"modelName":"Claude-3-Opus","inputPrice":{"amount":2.5,"tokenSize":1000000},"cachedInputPrice":{"amount":0.08,"tokenSize":1000000},"cachedOutputPrice":{"amount":1,"tokenSize":1000000},"outputPrice":{"amount":10,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'ANTHROPIC')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the ANTHROPIC pricing")
	}

	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'DEEPSEEK',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"deepseek-reasoner","inputPrice":{"amount":0.14,"tokenSize":1000000},"cachedInputPrice":{"amount":0.55,"tokenSize":1000000},"outputPrice":{"amount":2.19,"tokenSize":1000000}},{"modelName":"deepseek-chat","inputPrice":{"amount":0.27,"tokenSize":1000000},"cachedInputPrice":{"amount":0.07,"tokenSize":1000000},"outputPrice":{"amount":1.1,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'DEEPSEEK')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the DEEPSEEK pricing")
	}

	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'GEMINI',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"Gemini-2.0-Flash","inputPrice":{"amount":0.1,"tokenSize":1000000},"contextCachingPrice":{"amount":0.025,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":1,"tokenSize":1000000},"outputPrice":{"amount":0.4,"tokenSize":1000000}},{"modelName":"Gemini-2.0-Flash-Lite","inputPrice":{"amount":0.075,"tokenSize":1000000},"contextCachingPrice":{"amount":0.01875,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":1,"tokenSize":1000000},"outputPrice":{"amount":0.3,"tokenSize":1000000}},{"modelName":"Gemini-1.5-Flash","promptUpTo":128000,"inputPrice":{"amount":0.075,"tokenSize":1000000},"contextCachingPrice":{"amount":0.01875,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":1,"tokenSize":1000000},"outputPrice":{"amount":0.3,"tokenSize":1000000}},{"modelName":"Gemini-1.5-Flash","promptLongerThan":128000,"inputPrice":{"amount":0.15,"tokenSize":1000000},"contextCachingPrice":{"amount":0.0375,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":1,"tokenSize":1000000},"outputPrice":{"amount":0.6,"tokenSize":1000000}},{"modelName":"Gemini-1.5-Flash-8B","promptUpTo":128000,"inputPrice":{"amount":0.0375,"tokenSize":1000000},"contextCachingPrice":{"amount":0.01,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":0.25,"tokenSize":1000000},"outputPrice":{"amount":0.15,"tokenSize":1000000}},{"modelName":"Gemini-1.5-Flash-8B","promptLongerThan":128000,"inputPrice":{"amount":0.075,"tokenSize":1000000},"contextCachingPrice":{"amount":0.02,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":0.25,"tokenSize":1000000},"outputPrice":{"amount":0.3,"tokenSize":1000000}},{"modelName":"Gemini-1.5-Flash-Pro","promptUpTo":128000,"inputPrice":{"amount":1.25,"tokenSize":1000000},"contextCachingPrice":{"amount":0.3125,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":4.5,"tokenSize":1000000},"outputPrice":{"amount":5,"tokenSize":1000000}},{"modelName":"Gemini-1.5-Flash-Pro","promptLongerThan":128000,"inputPrice":{"amount":2.5,"tokenSize":1000000},"contextCachingPrice":{"amount":0.625,"tokenSize":1000000},"contextCachingStoragePrice":{"amount":4.5,"tokenSize":1000000},"outputPrice":{"amount":10,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'GEMINI')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the GEMINI pricing")
	}

	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'AZURE_OPENAI',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"o1-Global","inputPrice":{"amount":15,"tokenSize":1000000},"cachedInputPrice":{"amount":7.5,"tokenSize":1000000},"outputPrice":{"amount":60,"tokenSize":1000000}},{"modelName":"o1-Regional","inputPrice":{"amount":16.5,"tokenSize":1000000},"cachedInputPrice":{"amount":8.25,"tokenSize":1000000},"outputPrice":{"amount":66,"tokenSize":1000000}},{"modelName":"o3-mini-Global","inputPrice":{"amount":1.1,"tokenSize":1000000},"cachedInputPrice":{"amount":0.55,"tokenSize":1000000},"outputPrice":{"amount":4.4,"tokenSize":1000000}},{"modelName":"o3-mini-Regional","inputPrice":{"amount":1.21,"tokenSize":1000000},"cachedInputPrice":{"amount":0.605,"tokenSize":1000000},"outputPrice":{"amount":4.48,"tokenSize":1000000}},{"modelName":"4o-Global","inputPrice":{"amount":5,"tokenSize":1000000},"cachedInputPrice":{"amount":2.5,"tokenSize":1000000},"outputPrice":{"amount":20,"tokenSize":1000000}},{"modelName":"4o-Global","inputPrice":{"amount":5.5,"tokenSize":1000000},"cachedInputPrice":{"amount":2.75,"tokenSize":1000000},"outputPrice":{"amount":22,"tokenSize":1000000}},{"modelName":"4o-mini-Global","inputPrice":{"amount":0.6,"tokenSize":1000000},"cachedInputPrice":{"amount":0.3,"tokenSize":1000000},"outputPrice":{"amount":2.64,"tokenSize":1000000}},{"modelName":"4o-mini-Regional","inputPrice":{"amount":0.66,"tokenSize":1000000},"cachedInputPrice":{"amount":0.33,"tokenSize":1000000},"outputPrice":{"amount":2.64,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'AZURE_OPENAI')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the AZURE_OPENAI pricing")
	}

	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'AZURE_PHI',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"Phi-3-mini-4k-instruct","context":4000,"inputPrice":{"amount":0.00013,"tokenSize":1000000},"outputPrice":{"amount":0.00052,"tokenSize":1000000}},{"modelName":"Phi-3-mini-128k-instruct","context":128000,"inputPrice":{"amount":0.00013,"tokenSize":1000000},"outputPrice":{"amount":0.00052,"tokenSize":1000000}},{"modelName":"Phi-3.5-mini-instruct","context":128000,"inputPrice":{"amount":0.00013,"tokenSize":1000000},"outputPrice":{"amount":0.00052,"tokenSize":1000000}},{"modelName":"Phi-3-small-8k-instruct","context":8000,"inputPrice":{"amount":0.00015,"tokenSize":1000000},"outputPrice":{"amount":0.0006,"tokenSize":1000000}},{"modelName":"Phi-3-small-128k-instruct","context":128000,"inputPrice":{"amount":0.00015,"tokenSize":1000000},"outputPrice":{"amount":0.0006,"tokenSize":1000000}},{"modelName":"Phi-3-medium-4k-instruct","context":4000,"inputPrice":{"amount":0.00017,"tokenSize":1000000},"outputPrice":{"amount":0.00068,"tokenSize":1000000}},{"modelName":"Phi-3-medium-128k-instruct","context":128000,"inputPrice":{"amount":0.00017,"tokenSize":1000000},"outputPrice":{"amount":0.00068,"tokenSize":1000000}},{"modelName":"Phi-3.5-MoE-instruct","context":128000,"inputPrice":{"amount":0.00016,"tokenSize":1000000},"outputPrice":{"amount":0.00064,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'AZURE_PHI')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the AZURE_PHI pricing")
	}

	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'MISTRAL',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"Mistral-Large-24.11","inputPrice":{"amount":2,"tokenSize":1000000},"outputPrice":{"amount":6,"tokenSize":1000000}},{"modelName":"Pixtral-Large","inputPrice":{"amount":2,"tokenSize":1000000},"outputPrice":{"amount":6,"tokenSize":1000000}},{"modelName":"Mistral-Small-3","inputPrice":{"amount":0.1,"tokenSize":1000000},"outputPrice":{"amount":0.3,"tokenSize":1000000}},{"modelName":"Mistral-Saba","inputPrice":{"amount":0.2,"tokenSize":1000000},"outputPrice":{"amount":0.6,"tokenSize":1000000}},{"modelName":"Phi-3-small-128k-instruct","context":128000,"inputPrice":{"amount":0.00015,"tokenSize":1000000},"outputPrice":{"amount":0.0006,"tokenSize":1000000}},{"modelName":"Codestral","inputPrice":{"amount":0.3,"tokenSize":1000000},"outputPrice":{"amount":0.9,"tokenSize":1000000}},{"modelName":"Ministral-8B-24.10","inputPrice":{"amount":0.1,"tokenSize":1000000},"outputPrice":{"amount":0.1,"tokenSize":1000000}},{"modelName":"Ministral-3B-24.10","inputPrice":{"amount":0.04,"tokenSize":1000000},"outputPrice":{"amount":0.04,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'MISTRAL')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the MISTRAL pricing")
	}
	query = `UPDATE
			ai_model_price_list
		SET
			"llm_service_provider_name" = 'GROQ',
			"currency" = 'Dollars',
			"model_prices" = '[{"modelName":"DeepSeek-R1-Distill-Llama-70B","inputPrice":{"amount":0.75,"tokenSize":1000000},"outputPrice":{"amount":0.99,"tokenSize":1000000}},{"modelName":"DeepSeek-R1-Distill-Qwen-32B-128k","promptCount":128000,"inputPrice":{"amount":0.69,"tokenSize":1000000},"outputPrice":{"amount":0.69,"tokenSize":1000000}},{"modelName":"Qwen-2.5-32B-Instruct-128k","promptCount":128000,"inputPrice":{"amount":0.79,"tokenSize":1000000},"outputPrice":{"amount":0.79,"tokenSize":1000000}},{"modelName":"Qwen-2.5-Coder-32B-Instruct-128k","promptCount":128000,"inputPrice":{"amount":0.79,"tokenSize":1000000},"outputPrice":{"amount":0.79,"tokenSize":1000000}},{"modelName":"Llama-3.2-1B-(Preview)-8k","promptCount":8000,"inputPrice":{"amount":0.04,"tokenSize":1000000},"outputPrice":{"amount":0.04,"tokenSize":1000000}},{"modelName":"Llama-3.2-3B-(Preview)-8k","promptCount":8000,"inputPrice":{"amount":0.06,"tokenSize":1000000},"outputPrice":{"amount":0.06,"tokenSize":1000000}},{"modelName":"Llama-3.3-70B-Versatile-128k","promptCount":128000,"inputPrice":{"amount":0.59,"tokenSize":1000000},"outputPrice":{"amount":0.79,"tokenSize":1000000}},{"modelName":"Llama-3.1-8B-Instant-128k","promptCount":128000,"inputPrice":{"amount":0.05,"tokenSize":1000000},"outputPrice":{"amount":0.08,"tokenSize":1000000}},{"modelName":"Llama-3-70B-8k","promptCount":8000,"inputPrice":{"amount":0.59,"tokenSize":1000000},"outputPrice":{"amount":0.79,"tokenSize":1000000}},{"modelName":"Llama-3-8B-8k","promptCount":8000,"inputPrice":{"amount":0.05,"tokenSize":1000000},"outputPrice":{"amount":0.08,"tokenSize":1000000}},{"modelName":"Mixtral-8x7B-Instruct-32k","promptCount":32000,"inputPrice":{"amount":0.24,"tokenSize":1000000},"outputPrice":{"amount":0.24,"tokenSize":1000000}},{"modelName":"Gemma-2-9B-8k","promptCount":8000,"inputPrice":{"amount":0.2,"tokenSize":1000000},"outputPrice":{"amount":0.2,"tokenSize":1000000}},{"modelName":"Llama-Guard-3-8B-8k","promptCount":8000,"inputPrice":{"amount":0.2,"tokenSize":1000000},"outputPrice":{"amount":0.2,"tokenSize":1000000}},{"modelName":"Llama-3.3-70B-SpecDec-8k","promptCount":8000,"inputPrice":{"amount":0.59,"tokenSize":1000000},"outputPrice":{"amount":0.99,"tokenSize":1000000}}]',
			"created_at" = 0,
			"created_by" = '',
			"deleted_at" = NULL,
			"deleted_by" = '',
			"updated_at" = NULL,
			"updated_by" = '',
			"owner_account" = '',
			"vapus_id" = '',
			"last_audit_id" = '',
			"error_logs" = 'null',
			"ORGANIZATION" = '',
			"status" = '',
			"editors" = NULL,
			"scope" = ''
		WHERE
			("llm_service_provider_name" = 'GROQ')`

	_, err = dmstore.Db.PostgresClient.Conn.Exec(query)
	if err != nil {
		logger.Err(err).Msg("Error While updating the GROQ pricing")
	}
	return nil
}
