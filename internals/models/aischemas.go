package models

import (
	"encoding/json"
	"log"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type StructuredResponseFormat struct {
	Type       string              `json:"type,omitempty" yaml:"type,omitempty"`
	JsonSchema *ResponseJsonSchema `json:"jsonSchema,omitempty" yaml:"jsonSchema,omitempty"`
}

type ResponseJsonSchema struct {
	Name         string             `json:"name,omitempty" yaml:"name,omitempty"`
	Strict       bool               `json:"strict,omitempty" yaml:"strict,omitempty"`
	Schema       *FunctionParameter `json:"schema,omitempty" yaml:"schema,omitempty"`
	AutoGenerate bool               `json:"autoGenerate,omitempty" yaml:"autoGenerate,omitempty"`
	RawJsonParam string             `json:"rawJsonParam,omitempty" yaml:"rawJsonParam,omitempty"`
}

func (s *StructuredResponseFormat) ConvertFromPb(pb *mpb.StructuredResponseFormat) *StructuredResponseFormat {
	if pb == nil || pb.GetJsonSchema() == nil {
		return nil
	}
	jsonSchema := &ResponseJsonSchema{
		Name:   pb.GetJsonSchema().GetName(),
		Strict: pb.GetJsonSchema().GetStrict(),
	}
	schemaParam := &FunctionParameter{}
	err := json.Unmarshal([]byte(pb.GetJsonSchema().GetSchema()), schemaParam)
	if err != nil {
		log.Println("error in marshalling pb.GetJsonSchema().GetSchema(): ", err)
		return nil
	}
	jsonSchema.Schema = schemaParam
	jsonSchema.AutoGenerate = pb.GetJsonSchema().AutoGenerate
	jsonSchema.RawJsonParam = pb.GetJsonSchema().GetRawJsonParams()
	return &StructuredResponseFormat{
		Type:       pb.GetType(),
		JsonSchema: jsonSchema,
	}
}

func (s *StructuredResponseFormat) ConvertToPb() *mpb.StructuredResponseFormat {
	if s != nil {
		schema, err := json.Marshal(s.JsonSchema.Schema)
		if err != nil {
			log.Println("error in marshalling s.JsonSchema.Schema: ", err)
			return nil
		}
		return &mpb.StructuredResponseFormat{
			Type: s.Type,
			JsonSchema: &mpb.ResponseJsonSchema{
				Name:   s.JsonSchema.Name,
				Strict: s.JsonSchema.Strict,
				Schema: string(schema),
			},
		}
	}
	return nil
}

type FunctionCall struct {
	Name        string             `json:"name" yaml:"name"`
	Parameters  *FunctionParameter `json:"parameters" yaml:"parameters"`
	Description string             `json:"description" yaml:"description"`
	Required    []string           `json:"required" yaml:"required"`
	Strict      bool               `json:"strict" yaml:"strict"`
}

type FunctionParameter struct {
	Type        string                          `json:"type" yaml:"type"`
	Properties  map[string]*ParameterProperties `json:"properties" yaml:"properties"`
	Description string                          `json:"description" yaml:"description"`
	Required    []string                        `json:"required" yaml:"required"`
}

type ParameterProperties struct {
	Type                 string                          `json:"type" yaml:"type"`
	Pattern              string                          `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Items                string                          `json:"items,omitempty" yaml:"items,omitempty"`
	Description          string                          `json:"description" yaml:"description"`
	Properties           map[string]*ParameterProperties `json:"properties,omitempty" yaml:"properties,omitempty"`
	Enum                 []string                        `json:"enum,omitempty" yaml:"enum,omitempty"`
	Required             []string                        `json:"required,omitempty" yaml:"required,omitempty"`
	AdditionalProperties bool                            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
}

func (f *FunctionCall) ConvertFromPb(pb *mpb.FunctionCall) *FunctionCall {
	if pb == nil {
		return nil
	}
	param := &FunctionParameter{}
	err := json.Unmarshal([]byte(pb.Parameters), param)
	if err != nil {
		log.Println("error in marshalling pb.Parameters: ", err)
		return nil
	}

	return &FunctionCall{
		Name:        pb.GetName(),
		Parameters:  param,
		Description: pb.GetDescription(),
		Required:    pb.GetRequired(),
	}
}

func (f *FunctionCall) ConvertToPb() *mpb.FunctionCall {
	if f != nil {
		param, err := json.Marshal(f.Parameters)
		if err != nil {
			log.Println("error in marshalling f.Parameters: ", err)
			return nil
		}
		return &mpb.FunctionCall{
			Name:        f.Name,
			Parameters:  string(param),
			Description: f.Description,
			Required:    f.Required,
		}
	}
	return nil
}

func GetFunctionCallFromString(val string) *FunctionCall {
	toolSchema := &FunctionCall{}
	err := json.Unmarshal([]byte(val), toolSchema)
	if err != nil {
		log.Println("error in unmarshalling toolSchema: ", err)
		return nil
	}
	return toolSchema
}

func (f *FunctionCall) GetStringParamSchema() string {
	if f != nil {
		bbytes, err := json.MarshalIndent(f.Parameters, "", "  ")
		if err != nil {
			return ""
		}
		return string(bbytes)
	}
	return ""
}

func (f *FunctionCall) SetFunctionParamsFromString(val string, logger zerolog.Logger) error {
	toolSchema := &FunctionParameter{}
	err := json.Unmarshal([]byte(val), toolSchema)
	if err != nil {
		logger.Error().Err(err).Msg("error in unmarshalling FunctionParameter")
		return err
	}
	f.Parameters = toolSchema
	return nil
}
