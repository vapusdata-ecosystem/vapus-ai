package googlegenai

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
	"google.golang.org/api/iterator"
)

func (o *GoogleGenAI) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	modelsIter := o.client.ListModels(ctx)

	for {
		model, err := modelsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		obj := &models.AIModelBase{ // Use the imported type
			ModelId:          model.Name,
			OwnedBy:          "google",
			ModelName:        strings.ReplaceAll(model.Name, "models/", ""),
			Version:          model.Version,
			SupprtedOps:      model.SupportedGenerationMethods,
			InputTokenLimit:  model.InputTokenLimit,
			OutputTokenLimit: model.OutputTokenLimit,
			ModelNature:      []string{},
		}
		if slices.Contains(model.SupportedGenerationMethods, "embedContent") {
			obj.ModelType = mpb.AIModelType_EMBEDDING.String()
			obj.ModelNature = append(obj.ModelNature, mpb.AIModelType_EMBEDDING.String())

		} else {
			obj.ModelType = mpb.AIModelType_LLM.String()
			obj.ModelNature = append(obj.ModelNature, mpb.AIModelType_LLM.String())

		}
		result = append(result, obj)
	}
	return result, nil
}

// getStringFromMap safely gets a string value from the map for a given key.
func getStringFromMap(data map[string]any, key string) (string, bool) {
	val, ok := data[key]
	if !ok {
		return "", false
	}
	strVal, ok := val.(string)
	return strVal, ok
}

// getBoolFromMap safely gets a boolean value from the map for a given key.
func getBoolFromMap(data map[string]any, key string) (bool, bool) {
	val, ok := data[key]
	if !ok {
		return false, false
	}
	boolVal, ok := val.(bool)
	return boolVal, ok
}

// getStringArrayFromMap safely gets a string slice from the map for a given key.
func getStringArrayFromMap(data map[string]any, key string) ([]string, bool, error) {
	val, ok := data[key]
	if !ok {
		return nil, false, nil // Not present is not an error here
	}
	arrVal, ok := val.([]any)
	if !ok {
		return nil, true, fmt.Errorf("invalid type for key '%s', expected array, got %T", key, val)
	}

	result := make([]string, 0, len(arrVal))
	for i, item := range arrVal {
		strItem, ok := item.(string)
		if !ok {
			return nil, true, fmt.Errorf("invalid element type at index %d in array '%s', expected string, got %T", i, key, item)
		}
		result = append(result, strItem)
	}
	return result, true, nil
}

// getMapFromMap safely gets a map[string]any value from the map for a given key.
func getMapFromMap(data map[string]any, key string) (map[string]any, bool, error) {
	val, ok := data[key]
	if !ok {
		return nil, false, nil // Not present is not an error
	}
	mapVal, ok := val.(map[string]any)
	if !ok {
		return nil, true, fmt.Errorf("invalid type for key '%s', expected map[string]any, got %T", key, val)
	}
	return mapVal, true, nil
}

// --- Main Conversion Function ---

// ConvertMapToSchema converts a map[string]any representing a schema definition
// into the target Schema struct. It performs validation and handles nested structures.
func ConvertMapToSchema(data map[string]any, logg zerolog.Logger) (*genai.Schema, error) {
	if data == nil {
		return nil, fmt.Errorf("input schema map cannot be nil")
	}

	schema := &genai.Schema{} // Initialize the result struct

	// 1. Process Required 'type' field
	typeStr, ok := getStringFromMap(data, "type")
	if !ok {
		logg.Error().Msg("required key 'type' is missing or not a string")
		return nil, fmt.Errorf("required key 'type' is missing or not a string")
	}

	// Optional: Validate if typeStr is one of the known Types
	switch typeStr {
	// case genai.TypeString, genai.TypeNumber, genai.TypeInteger, genai.TypeBoolean, genai.TypeArray, genai.TypeObject:
	case "string":
		schema.Type = genai.TypeString // Valid type
	case "number":
		schema.Type = genai.TypeNumber // Valid type
	case "integer":
		schema.Type = genai.TypeInteger // Valid type
	case "boolean":
		schema.Type = genai.TypeBoolean // Valid type
	case "array":
		schema.Type = genai.TypeArray // Valid type
	case "object":
		schema.Type = genai.TypeObject // Valid type
		// Valid type
	default:
		return nil, fmt.Errorf("invalid value for 'type': '%s'", typeStr)
	}

	// 2. Process Optional fields
	schema.Format, _ = getStringFromMap(data, "format")           // Ignore !ok, defaults to ""
	schema.Description, _ = getStringFromMap(data, "description") // Ignore !ok, defaults to ""
	schema.Nullable, _ = getBoolFromMap(data, "nullable")         // Ignore !ok, defaults to false

	// 3. Process 'enum'
	enumVal, present, err := getStringArrayFromMap(data, "enum")
	if err != nil {
		return nil, fmt.Errorf("failed to process 'enum': %w", err) // Wrap error
	}
	if present {
		if schema.Type != genai.TypeString {
			logg.Debug().Msgf("Warning: 'enum' key found but schema type is '%s', ignoring 'enum'", schema.Type)
		}
		schema.Enum = enumVal
	}

	// 4. Process 'items' (Recursive) - Only if type is Array
	if schema.Type == genai.TypeArray {
		itemsMap, present, err := getMapFromMap(data, "items")
		if err != nil {
			return nil, fmt.Errorf("failed to process 'items': %w", err)
		}
		if !present {
			return nil, fmt.Errorf("key 'items' is required when type is '%s'", genai.TypeArray)
		}
		if itemsMap != nil { // Should not be nil if present and no error
			itemsSchema, err := ConvertMapToSchema(itemsMap, logg) // Recursive call
			if err != nil {
				return nil, fmt.Errorf("failed to convert schema for 'items': %w", err)
			}
			schema.Items = itemsSchema
		} else {
			// This case might occur if "items": null was in the input json
			return nil, fmt.Errorf("key 'items' cannot be null when type is '%s'", genai.TypeArray)
		}
	} else if _, itemsPresent := data["items"]; itemsPresent {
		// Items key exists but type is not Array
		logg.Debug().Msgf("Warning: 'items' key found but schema type is '%s', ignoring 'items'", schema.Type)
	}

	// 5. Process 'properties' (Recursive) - Only if type is Object
	if schema.Type == genai.TypeObject {
		propsMap, present, err := getMapFromMap(data, "properties")
		if err != nil {
			return nil, fmt.Errorf("failed to process 'properties': %w", err)
		}
		if present && propsMap != nil { // properties can be optional, allow empty map
			schema.Properties = make(map[string]*genai.Schema)
			for key, propVal := range propsMap {
				propSchemaMap, ok := propVal.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid type for property '%s', expected map[string]any, got %T", key, propVal)
				}
				propSchema, err := ConvertMapToSchema(propSchemaMap, logg) // Recursive call
				if err != nil {
					return nil, fmt.Errorf("failed to convert schema for property '%s': %w", key, err)
				}
				schema.Properties[key] = propSchema
			}
		} else if present && propsMap == nil {
			logg.Debug().Msgf("Warning: 'properties' key found but value is null, initializing empty map")
			schema.Properties = make(map[string]*genai.Schema) // Initialize empty map
		}
	} else if _, propsPresent := data["properties"]; propsPresent {
		logg.Debug().Msgf("Warning: 'properties' key found but schema type is '%s', ignoring 'properties'", schema.Type)
	}

	// 6. Process 'required' - Only if type is Object
	if schema.Type == genai.TypeObject {
		requiredVal, present, err := getStringArrayFromMap(data, "required")
		if err != nil {
			return nil, fmt.Errorf("failed to process 'required': %w", err)
		}
		if present {
			// Validate that required properties actually exist in the properties map
			if schema.Properties != nil {
				for _, reqKey := range requiredVal {
					if _, exists := schema.Properties[reqKey]; !exists {
						logg.Debug().Msgf("Warning: Required property '%s' not found in properties map", reqKey)
					}
				}
			} else if len(requiredVal) > 0 {
				logg.Debug().Msgf("Warning: 'required' key found but properties map is empty, ignoring 'required'")
			}
			schema.Required = requiredVal
		}
	} else if _, requiredPresent := data["required"]; requiredPresent {
		logg.Debug().Msgf("Warning: 'required' key found but schema type is '%s', ignoring 'required'", schema.Type)
	}

	return schema, nil
}

func logToolCallSchema(toolCallSchema *genai.Schema) {
	// Log the schema details
	if toolCallSchema.Items != nil {
		log.Println("Items Schema:")
		logToolCallSchema(toolCallSchema.Items)
	}
	if len(toolCallSchema.Properties) > 0 {
		log.Println("Properties:")
		for key, prop := range toolCallSchema.Properties {
			log.Println("Key:", key)
			logToolCallSchema(prop)
		}
	}
}

// TODO:  https://ai.google.dev/gemini-api/docs/document-processing?lang=go#technical-details (use inline params as well)
func BuildInputContent[T *prompts.SessionMessage | *pb.ChatMessageObject](ctx context.Context, client *genai.Client, obj T) []genai.Part {
	response := make([]genai.Part, 0)
	if obj == nil {
		return response
	}
	content := ""
	structuredContent := make([]*pb.RequestContentPart, 0)
	switch v := any(obj).(type) {
	case *prompts.SessionMessage:
		content = v.Message
		structuredContent = v.StructuredMessage
	case *pb.ChatMessageObject:
		content = v.Content
		structuredContent = v.StructuredContent
	default:
		log.Printf("Unsupported type: %T", v)
		return response
	}
	if content != "" {
		response = append(response, genai.Text(content))
		return response
	}
	if structuredContent != nil {
		for _, value := range structuredContent {
			switch value.Type {
			case aicore.AIResponseFormatText.String():
				response = append(response, genai.Text(value.Text))
			case aicore.AIResponseFormatImageUrl.String():
				if value.ImageUrl != nil {
					name := filepath.Base(value.ImageUrl.GetUrl())
					if value.ImageUrl.GetData() != "" {
						mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
						if !ok {
							log.Printf("Failed to get format for %s: ", name)
							continue
						}
						response = append(response, genai.ImageData(mimeType[0], []byte(value.ImageUrl.GetData())))
					} else if value.ImageUrl.GetUrl() != "" {
						mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
						if !ok {
							log.Printf("Failed to get format for %s: ", name)
							continue
						}
						file, err := client.UploadFile(ctx, name, strings.NewReader(value.ImageUrl.Data), &genai.UploadFileOptions{
							DisplayName: name,
							MIMEType:    mimeType[0],
						})
						if err != nil {
							log.Printf("Failed to Image file %s: %v", name, err)
							continue
						}
						response = append(response, genai.FileData{
							URI:      file.URI,
							MIMEType: file.MIMEType,
						})
					}
				}
			case aicore.AIResponseFormatInputAudio.String():
				if value.InputAudio != nil {
					name := filepath.Base(value.InputAudio.GetUrl())
					if value.InputAudio.GetData() != "" {
						if value.Upload {
							mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
							if !ok {
								log.Printf("Failed to get format for %s: ", name)
								continue
							}
							file, err := client.UploadFile(ctx, name, strings.NewReader(value.InputAudio.Data), &genai.UploadFileOptions{
								DisplayName: name,
								MIMEType:    mimeType[0],
							})
							if err != nil {
								log.Printf("Failed to InputAudio file %s: %v", name, err)
								continue
							}
							response = append(response, genai.FileData{
								URI:      file.URI,
								MIMEType: file.MIMEType,
							})
						} else {
							mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
							if !ok {
								log.Printf("Failed to get format for %s: ", name)
								continue
							}
							response = append(response, genai.Blob{
								MIMEType: mimeType[0],
								Data:     []byte(value.InputAudio.Data),
							})
						}

					}
				}
			case aicore.AIResponseFormatInputFile.String():
				if value.File != nil {
					name := filepath.Base(value.File.GetUrl())
					if value.File.GetFileData() != "" {
						if value.Upload {
							mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
							if !ok {
								log.Printf("Failed to get format for %s: ", name)
								continue
							}
							file, err := client.UploadFile(ctx, name, strings.NewReader(value.File.FileData), &genai.UploadFileOptions{
								DisplayName: name,
								MIMEType:    mimeType[0],
							})
							if err != nil {
								log.Printf("Failed to InputAudio file %s: %v", name, err)
								continue
							}
							response = append(response, genai.FileData{
								URI:      file.URI,
								MIMEType: file.MIMEType,
							})
						} else {
							mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
							if !ok {
								log.Printf("Failed to get format for %s: ", name)
								continue
							}
							response = append(response, genai.Blob{
								MIMEType: mimeType[0],
								Data:     []byte(value.File.FileData),
							})
						}

					}
				}
			default:
				response = append(response, genai.Text(value.Text))
			}
		}
	}
	return response
}
