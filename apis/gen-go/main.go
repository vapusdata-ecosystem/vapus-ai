package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func main() {
	// Load OpenAPI v2 JSON
	v2File := "apidocs.swagger.json" // Path to your OpenAPI v2 JSON file
	data, err := os.ReadFile(v2File)
	if err != nil {
		log.Fatalf("Failed to read OpenAPI v2 file: %v", err)
	}

	// Parse the OpenAPI v2 spec
	v2Loader := &openapi2.T{}
	if err := json.Unmarshal(data, &v2Loader); err != nil {
		log.Fatalf("Failed to parse OpenAPI v2 JSON: %v", err)
	}

	// Convert to OpenAPI v3
	v3Spec, err := openapi2conv.ToV3(v2Loader)
	if err != nil {
		log.Fatalf("Failed to convert to OpenAPI v3: %v", err)
	}

	// Marshal the OpenAPI v3 spec to JSON
	v3Data, err := json.MarshalIndent(v3Spec, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI v3 spec: %v", err)
	}

	err = os.Remove(v2File)
	if err != nil {
		log.Fatalf("Failed to delete OpenAPI v2 file: %v", err)
	}
	// Write the converted spec to a new file
	v3File := "../api-spec/apidocs.openapiv3.json"
	err = os.WriteFile(v3File, v3Data, 0644)
	if err != nil {
		log.Fatalf("Failed to write OpenAPI v3 file: %v", err)
	}

	fmt.Printf("Successfully converted OpenAPI v2 to v3. Output written to %s\n", v3File)
}
