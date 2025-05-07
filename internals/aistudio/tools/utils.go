package tools

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func GenerateAIToolSchemaFromJson[T any](data T) map[string]any {
	d := any(data)
	if d == nil {
		return nil
	}
	switch v := d.(type) {
	case []any:
		if len(v) == 0 {
			return map[string]any{
				"type":  "array",
				"items": map[string]any{},
			}
		}
		return map[string]any{
			"type":  "array",
			"items": GenerateAIToolSchemaFromJson(v[0]),
		}
	case map[string]any:
		schema := map[string]any{
			"type":       "object",
			"properties": map[string]any{},
			"required":   []string{},
		}
		for key, val := range v {
			schema["properties"].(map[string]any)[key] = GenerateAIToolSchemaFromJson(val)
			schema["required"] = append(schema["required"].([]string), key)
		}
		return schema
	case string:
		return map[string]any{
			"type": "string",
		}
	case float64, float32:
		return map[string]any{
			"type": "number",
		}
	case int32, int64, int16, int8, int:
		return map[string]any{
			"type": "number",
		}
	case bool:
		return map[string]any{
			"type": "boolean",
		}
	case []byte:
		var res map[string]any
		if err := json.Unmarshal(v, &res); err != nil {
			return nil
		}
		return GenerateAIToolSchemaFromJson(res)
	default:
		return map[string]any{}
	}
}

func GenerateAIToolSchema[T any](data T, req *models.FunctionCall, logger zerolog.Logger) (*models.FunctionCall, error) {
	d := any(data)
	if d == nil {
		return nil, ErrInvalidTools
	}
	if req == nil {
		return nil, ErrInvalidFunctionSchema
	}
	schemaMap := GenerateAIToolSchemaFromJson(d)
	if schemaMap == nil {
		return nil, ErrInvalidJSONTools
	}
	schemaBytes, err := json.Marshal(schemaMap)
	if err != nil {
		return nil, err
	}
	err = req.SetFunctionParamsFromString(string(schemaBytes), logger)

	if err != nil {
		logger.Err(err).Msg("error while setting function params")
		return nil, err
	}
	return req, nil
}
