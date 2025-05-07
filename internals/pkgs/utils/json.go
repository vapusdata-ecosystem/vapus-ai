package dmutils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/invopop/jsonschema"
)

func Marshall(val any) ([]byte, error) {
	valType := reflect.TypeOf(val)
	isStruct := valType.Kind() == reflect.Ptr && valType.Elem().Kind() == reflect.Struct
	if isStruct {
		return sonic.Marshal(val)
	} else {
		return json.Marshal(val)
	}
}

func Unmarshall(data []byte, val any) error {
	valType := reflect.TypeOf(val)
	isStruct := valType.Kind() == reflect.Ptr && valType.Elem().Kind() == reflect.Struct
	if isStruct {
		return sonic.Unmarshal(data, val)
	} else {
		return json.Unmarshal(data, val)
	}
}

// func MarshalJsonSchema[T string | map[string]any](schemaMap T) (*jsonschema.Schema, error) {
// 	fMap := make(map[string]any)
// 	log.Println("schemaMap:", schemaMap)
// 	// Use type assertions with comma-ok syntax instead of type switch
// 	if strVal, ok := any(schemaMap).(string); ok {
// 		bbytes, err := sonic.Marshal(strVal)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to marshal schema map: %w", err)
// 		}
// 		log.Println("//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", string(bbytes))
// 		err = sonic.Unmarshal(bbytes, &fMap)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal schema map: %w", err)
// 		}
// 	} else if mapVal, ok := any(schemaMap).(map[string]any); ok {
// 		fMap = mapVal
// 	} else {
// 		return nil, fmt.Errorf("unsupported schema map type: %T", schemaMap)
// 	}

// 	// 1. Marshal the map back into JSON bytes
// 	// This step converts your Go map representation back into the raw JSON format
// 	// that the jsonschema library expects.
// 	jsonData, err := json.Marshal(schemaMap)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal map back to JSON: %w", err)
// 	}

// 	// 2. Create an io.Reader from the JSON bytes
// 	reader := bytes.NewReader(jsonData)

// 	// 3. Use the jsonschema compiler
// 	// The library needs a URL identifier for the schema, even if it's in memory.
// 	// This is used for resolving relative $refs, if any. We'll use a dummy URL.
// 	// You could also use a simple string like "schema.json".
// 	dummyURL := "inmemory://schema.json"

// 	// Create a new compiler instance
// 	compiler := jsonschema.NewCompiler()

// 	// Add the JSON schema data from the reader as a resource to the compiler
// 	// using the dummy URL as its identifier.
// 	if err := compiler.AddResource(dummyURL, reader); err != nil {
// 		return nil, fmt.Errorf("failed to add schema resource to compiler: %w", err)
// 	}

// 	// 4. Compile the schema using its identifier
// 	schema, err := compiler.Compile(dummyURL)
// 	if err != nil {
// 		// Compilation errors often include details about schema validity
// 		return nil, fmt.Errorf("failed to compile JSON schema: %w", err)
// 	}

// 	// Success!
// 	return schema, nil
// }

func ConvertStructToJsonSchema(obj any, description string) string {
	schema := jsonschema.Reflect(obj)
	schema.Description = description

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return ""
	}

	fmt.Println(string(schemaJSON))
	return string(schemaJSON)
}

func ConvertStructToJsonSchemaWithType(obj any, description string) string {
	requestType := reflect.TypeOf(obj).Elem()

	reflector := jsonschema.Reflector{
		DoNotReference:             false, // Allow $ref for nested types (usually good)
		RequiredFromJSONSchemaTags: true,  // Use the 'required' tag explicitly
	}

	schema := reflector.ReflectFromType(requestType)

	schema.Description = description

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return ""
	}

	fmt.Println(string(schemaJSON))
	return string(schemaJSON)
}
