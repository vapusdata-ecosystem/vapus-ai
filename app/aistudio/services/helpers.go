package services

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	aimodels "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/providers"
	aitool "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/tools"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

func crawlAIModels(ctx context.Context, modelNode *models.AIModelNode, logger zerolog.Logger) error {
	aiConfig, err := aimodels.NewAIModelNode(aimodels.WithAIModelNode(modelNode), aimodels.WithLogger(logger))
	if err != nil {
		logger.Err(err).Msgf("error while creating AI model node for model - %v", modelNode.Name)
		return err
	}
	// TO merge the upcoming model and previous model
	resultGenerativeModels := make([]*models.AIModelBase, 0)
	resultEmbeddingModels := make([]*models.AIModelBase, 0)
	tempGenerative := make(map[string]bool)
	tempEmbedding := make(map[string]bool)

	// storing the unique value of generative Models in resultGenerativeModels
	for _, val := range modelNode.GenerativeModels {
		key := val.ModelName
		if !tempGenerative[key] {
			tempGenerative[key] = true
			resultGenerativeModels = append(resultGenerativeModels, val)
		}
	}
	// storing the unique value of Embedding Models in resultEmbeddingModels
	for _, val := range modelNode.EmbeddingModels {
		key := val.ModelName
		if !tempEmbedding[key] {
			tempEmbedding[key] = true
			resultEmbeddingModels = append(resultEmbeddingModels, val)
		}
	}

	result, err := aiConfig.CrawlModels(ctx)
	if err != nil {
		logger.Err(err).Msgf("error while crawling models for model - %v", modelNode.Name)
		return err
	}

	for _, model := range result {
		if model.ModelType == mpb.AIModelType_EMBEDDING.String() {
			key := model.ModelName
			if !tempEmbedding[key] {
				tempEmbedding[key] = true
				resultEmbeddingModels = append(resultEmbeddingModels, model)
			}
		} else {
			key := model.ModelName
			if !tempGenerative[key] {
				tempGenerative[key] = true
				resultGenerativeModels = append(resultGenerativeModels, model)
			}
		}
	}
	modelNode.GenerativeModels = make([]*models.AIModelBase, 0)
	modelNode.EmbeddingModels = make([]*models.AIModelBase, 0)
	modelNode.EmbeddingModels = resultEmbeddingModels
	modelNode.GenerativeModels = resultGenerativeModels
	return nil
}

func BuildAIPromptTemplate(obj *models.AIPrompt) {
	if obj == nil {
		return
	}
	templaterMap := []map[string]string{}
	template := obj.Spec.UserMessage + "\n"
	if obj.Spec.InputTag != "" {
		template = template +
			strings.Replace(aicore.StartTagTemplate, "TAG", obj.Spec.InputTag, 1) +
			"[" + obj.Spec.InputTag + "]" +
			strings.Replace(aicore.EndTagTemplate, "TAG", obj.Spec.InputTag, 1) +
			"\n"
	} else {
		template = template + "[Input] - \n"
	}
	if obj.Spec.ContextTag != "" {
		template = template +
			strings.Replace(aicore.StartTagTemplate, "TAG", obj.Spec.ContextTag, 1) +
			"[" + obj.Spec.ContextTag + "]" +
			strings.Replace(aicore.EndTagTemplate, "TAG", obj.Spec.ContextTag, 1) +
			"\n"
	}
	if obj.Spec.OutputTag != "" {
		template = template +
			strings.Replace(aicore.StartTagTemplate, "TAG", obj.Spec.OutputTag, 1) +
			strings.Replace(aicore.EndTagTemplate, "TAG", obj.Spec.OutputTag, 1) +
			"\n"
	}
	if obj.Spec.Sample != nil {
		if obj.Spec.Sample.InputText != "" {
			template = template + "[Sample Input] " + obj.Spec.Sample.InputText + "\n"
		}
		if obj.Spec.Sample.Response != "" {
			template = template + "[Sample Output] " + obj.Spec.Sample.Response + "\n"
		}
		// template = template + "[Sample Input] " + obj.Spec.Sample.InputText + "\n"
		// template = template + "[Sample Output] " + obj.Spec.Sample.Response + "\n"
	}
	sysTemplate := obj.Spec.SystemMessage
	obj.UserTemplate = template
	obj.Spec.SystemMessage = sysTemplate
	templaterMap = append(templaterMap, map[string]string{
		"content": obj.UserTemplate,
		"role":    "user",
	})
	templaterMap = append(templaterMap, map[string]string{
		"content": obj.Spec.SystemMessage,
		"role":    "system",
	})
	bbytes, err := json.MarshalIndent(templaterMap, "", "  ")
	if err != nil {
		return
	}
	obj.Template = string(bbytes)
	// template = template + "[Instruction] - Please follow the above instructions to get proper content, input and expected output in desired format.\n Do not change the format of the content, input and expected output."
	return
}

func BuildPromptSchema(obj *models.AIPrompt, logger zerolog.Logger) error {
	var err error
	for _, tool := range obj.Spec.Tools {
		if tool.AutoGenerate {
			if tool.RawJsonParam != "" {
				if tool.Schema.Name != "" || tool.Schema.Description != "" {
					tool.Schema, err = aitool.GenerateAIToolSchema(tool.RawJsonParam, tool.Schema, logger)
					if err != nil {
						logger.Err(err).Msg("error while setting function params from raw json param")
						return err
					} else {
						return nil
					}
				}
			}
		}
	}
	return nil
}

func BuildPromptResponseFormat(obj *models.AIPrompt, logger zerolog.Logger) error {
	if obj.Spec.ResponseFormat == nil {
		return nil
	}
	switch obj.Spec.ResponseFormat.Type {
	case aicore.AIResponseFormatJSONObject.String():
		return nil
	case aicore.AIResponseFormatJSONSchema.String():
		if obj.Spec.ResponseFormat.JsonSchema.AutoGenerate {
			if obj.Spec.ResponseFormat.JsonSchema.RawJsonParam != "" && obj.Spec.ResponseFormat.JsonSchema.Schema != nil {
				schema, err := aitool.GenerateAIToolSchema(obj.Spec.ResponseFormat.JsonSchema.RawJsonParam, &models.FunctionCall{
					Name:   obj.Spec.ResponseFormat.JsonSchema.Name,
					Strict: obj.Spec.ResponseFormat.JsonSchema.Strict,
				}, logger)
				if err != nil {
					logger.Err(err).Msg("error while setting function params from raw json param")
					return err
				}
				if schema != nil {
					obj.Spec.ResponseFormat.JsonSchema.Schema = schema.Parameters
					return nil
				}
			} else {
				logger.Error().Msg("error: invalid ai response format, missing raw json param")
				return apperr.ErrInvalidAIResponseFormat
			}
		}
	default:
		return apperr.ErrInvalidAIResponseFormat
	}

	return aitool.ErrInvalidFunctionSchema
}

func BuildPromptVariables(obj *models.PromptSpec, logger zerolog.Logger) error {
	if obj == nil || obj.UserMessage == "" {
		logger.Error().Msg("error: invalid ai prompt spec, missing spec")
		return apperr.ErrInvalidAIPromptRequestSpec
	}
	matches := dmutils.TemplateVarRegex.FindAllString(obj.UserMessage, -1)
	matches = lo.Map(matches, func(s string, index int) string {
		val := strings.Replace(s, "{{", "", -1)
		val = strings.Replace(val, "}}", "", -1)
		return val
	})
	obj.Variables = lo.Union(obj.Variables, matches)
	return nil
}

// func BuildAgentFuntioncallRenderer(obj *models.VapusAgents) string {
// 	if obj == nil {
// 		return ""
// 	}
// 	required := []string{}
// 	tool := &models.FunctionCall{
// 		Name: dmutils.SlugifyBase(obj.Name),
// 		Parameters: &models.FunctionParameter{
// 			Type:        aitool.FuncParamType,
// 			Properties:  map[string]*models.ParameterProperties{},
// 			Description: obj.Description,
// 		},
// 	}
// 	for _, step := range obj.Steps {
// 		tool.Parameters.Properties[step.Id] = &models.ParameterProperties{
// 			Type:        strings.ToLower(step.ValueType),
// 			Description: step.Prompt,
// 		}
// 		if step.Required {
// 			required = append(required, step.Id)
// 		}
// 	}
// 	// schema, err := dmutils.GenericMarshaler(tool, mpb.ContentFormats_JSON.String())
// 	// if err != nil {
// 	// 	return ""
// 	// }
// 	schema, err := json.MarshalIndent(tool, "", "  ")
// 	if err != nil {
// 		return ""
// 	}
// 	return string(schema)
// }

func BuildGuardrailSchema(obj *models.AIGuardrails) string {
	if obj == nil {
		return ""
	}
	tool := &models.FunctionCall{
		Name:        dmutils.SlugifyBase(obj.Name),
		Description: obj.Description,
		Parameters: &models.FunctionParameter{
			Type:        aitool.FuncParamType,
			Properties:  map[string]*models.ParameterProperties{},
			Description: "Guardrails for AI models, based on the different topics and its conent, the content will be filtered.",
		},
	}

	topicParams := &models.ParameterProperties{
		Type: strings.ToLower(mpb.AgentStepValueType_OBJECT.String()),
		Description: `Guardrails for different topics,based on the different topics, the content will be filtered.
		 If the content is not related to the topic, return false for below properties. If they are related to the topic, return true.`,
		Properties: map[string]*models.ParameterProperties{},
	}
	contentParams := &models.ParameterProperties{
		Type: strings.ToLower(mpb.AgentStepValueType_OBJECT.String()),
		Description: `Guardrails for different contents,based on the different contents, the content will be filtered. Related levels & types of content is defined in enums. 
		If user input is not related to the content, return NONE for below properties. If they are related to the content, return the severity.
		Valid values for severity are: NONE, LOW, MEDIUM, HIGH that should be returned based on the content.`,
		Properties: map[string]*models.ParameterProperties{},
	}
	for _, topic := range obj.Topics {
		topicParams.Properties[dmutils.SlugifyBase(topic.Topic)] = &models.ParameterProperties{
			Type:        strings.ToLower(mpb.AgentStepValueType_BOOLEAN.String()),
			Description: topic.Description,
		}
	}
	tool.Parameters.Properties["topic_guardrails"] = topicParams
	contentParams.Properties[dmutils.SlugifyBase("Hate Speech")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: "Hate speech is a communication that carries no meaning other than the expression of hatred for some group, especially in circumstances in which the communication is likely to provoke violence.",
	}
	contentParams.Properties[dmutils.SlugifyBase("Insults")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Insults are words or actions that are intended to be rude or offensive, often because they are directed at a particular person or group.`,
	}
	contentParams.Properties[dmutils.SlugifyBase("Threats")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Threats are statements of an intention to inflict pain, injury, damage, or other hostile action on someone in retribution for something done or not done.`,
	}
	contentParams.Properties[dmutils.SlugifyBase("Sexual")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Sexual content is any material depicting, describing, or alluding to sexual behavior or anatomy that is intended to arouse sexual feelings.`,
	}
	contentParams.Properties[dmutils.SlugifyBase("Misconduct")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Misconduct is behavior that is illegal or dishonest, or that is considered morally wrong by most people.`,
	}
	tool.Parameters.Properties["content_guardrails"] = contentParams
	schema, err := json.MarshalIndent(tool, "", "  ")
	if err != nil {
		return ""
	}
	return string(schema)
}
