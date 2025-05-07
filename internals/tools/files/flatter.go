package filetools

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Global map to collect all unique headers
var allHeaders map[string]struct{}

func main() {
	// --- Command Line Flags ---
	inputFile := flag.String("input", "", "Path to the input JSON or YAML file (required)")
	outputFile := flag.String("output", "", "Path to the output CSV file (required)")
	sep := flag.String("sep", ".", "Separator for nested keys")
	listSep := flag.String("listsep", "|", "Separator for joining items in arrays of primitives")
	forceFlatten := flag.Bool("force-flatten", false, "Force flattening even if nesting isn't detected")

	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Error: --input and --output flags are required.")
		flag.Usage()
		os.Exit(1)
	}

	// --- Read Input File ---
	log.Printf("Reading input file: %s\n", *inputFile)
	inputBytes, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file %s: %v", *inputFile, err)
	}

	// --- Unmarshal Data ---
	var data interface{}
	fileExt := strings.ToLower(filepath.Ext(*inputFile))

	log.Println("Detecting format and unmarshaling...")
	switch fileExt {
	case ".json":
		err = json.Unmarshal(inputBytes, &data)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal(inputBytes, &data)
		if err != nil {
			log.Fatalf("Error unmarshaling YAML: %v", err)
		}
	default:
		log.Fatalf("Error: Unsupported file extension '%s'. Only .json, .yaml, .yml are supported.", fileExt)
	}
	log.Println("Unmarshaling successful.")

	// --- Prepare for CSV ---
	allHeaders = make(map[string]struct{}) // Initialize header collector
	var finalRows []map[string]string      // To store the rows for CSV writing

	// --- Check for Nesting and Flatten/Convert ---
	shouldFlatten := *forceFlatten || isNested(data, 0)

	if shouldFlatten {
		log.Println("Nesting detected (or forced). Flattening data...")
		flattenedRows := make([]map[string]string, 0)
		// Handle top-level array or single object
		switch v := data.(type) {
		case []interface{}:
			for _, item := range v {
				flattenRecursive(item, "", *sep, *listSep, make(map[string]string), &flattenedRows)
			}
		case map[string]interface{}:
			flattenRecursive(v, "", *sep, *listSep, make(map[string]string), &flattenedRows)
		default:
			log.Println("Warning: Input data is not an array or object. Cannot flatten effectively.")
			// Handle primitive top-level if necessary, e.g., create a single cell CSV
		}
		finalRows = flattenedRows
		log.Printf("Flattening complete. Generated %d rows (including expansions).\n", len(finalRows))

	} else {
		log.Println("No significant nesting detected. Performing direct conversion...")
		var headersList []string
		finalRows, headersList = convertToCSVDirect(data)
		// Add headers collected from direct conversion
		for _, h := range headersList {
			allHeaders[h] = struct{}{}
		}
		log.Printf("Direct conversion complete. Generated %d rows.\n", len(finalRows))
	}

	if len(finalRows) == 0 {
		log.Println("Warning: No data rows were generated. Output CSV will be empty or headers only.")
	}

	// --- Sort Headers ---
	sortedHeaders := make([]string, 0, len(allHeaders))
	for header := range allHeaders {
		sortedHeaders = append(sortedHeaders, header)
	}
	sort.Strings(sortedHeaders)

	// --- Write CSV Output ---
	log.Printf("Writing CSV data to: %s\n", *outputFile)
	outFile, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v", *outputFile, err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush() // Ensure all buffered data is written

	// Write Header
	if err := writer.Write(sortedHeaders); err != nil {
		log.Fatalf("Error writing CSV header: %v", err)
	}

	// Write Rows
	record := make([]string, len(sortedHeaders)) // Pre-allocate record slice
	for _, rowMap := range finalRows {
		for i, header := range sortedHeaders {
			record[i] = rowMap[header] // Map lookup defaults to zero value ("") if key missing
		}
		if err := writer.Write(record); err != nil {
			// Log error but continue trying to write other rows if possible
			log.Printf("Error writing CSV record: %v (row data: %v)", err, rowMap)
		}
	}

	// Check for final errors from Flush()
	if err := writer.Error(); err != nil {
		log.Fatalf("Error flushing CSV writer: %v", err)
	}

	log.Println("CSV file successfully created.")
}

// --- Nesting Detection ---

// isNested checks if the data structure contains maps or slices within maps/slices.
func isNested(data interface{}, depth int) bool {
	switch v := data.(type) {
	case map[string]interface{}:
		if depth > 0 { // Found a map inside something else
			return true
		}
		for _, val := range v {
			if isNested(val, depth+1) {
				return true
			}
		}
	case []interface{}:
		// Treat top-level array differently - check its elements for nesting
		// An array itself isn't nesting unless it's inside another array/map
		// OR its elements are nested.
		if depth > 0 && len(v) > 0 { // Found a non-empty array inside something else
			// Check if it's just primitives
			isPrim := true
			for _, item := range v {
				switch item.(type) {
				case map[string]interface{}, []interface{}:
					isPrim = false
					break
				}
			}
			if !isPrim { // It's an array of complex types nested
				return true
			}
			// If it's primitives, it doesn't count as *structural* nesting for this check
		}
		// Always check elements regardless of depth
		for _, item := range v {
			if isNested(item, depth+1) {
				return true
			}
		}
	}
	return false // Primitive or nil
}

// --- Flattening Logic ---

func flattenRecursive(data interface{}, parentKey, sep, listSep string, baseRow map[string]string, allRows *[]map[string]string) {
	switch value := data.(type) {
	case map[string]interface{}:
		// Create a snapshot of the baseRow *before* processing this map's items
		// This is crucial for list expansion later.
		currentLevelBase := cloneMap(baseRow)
		expansionTargets := make([]struct {
			key  string
			list []interface{}
		}, 0, 1) // Usually expand only first list

		// 1. Process non-list-of-map items first to build the base row additions
		tempItems := make(map[string]string) // Items flattened *at this map level*

		for k, v := range value {
			newKey := k
			if parentKey != "" {
				newKey = parentKey + sep + k
			}

			// Check if it's a list potentially needing expansion
			if listVal, ok := v.([]interface{}); ok && len(listVal) > 0 {
				isListOfMaps := true
				for _, item := range listVal {
					if _, isMap := item.(map[string]interface{}); !isMap {
						isListOfMaps = false
						break
					}
				}
				if isListOfMaps {
					// Mark for expansion later, process other items first
					expansionTargets = append(expansionTargets, struct {
						key  string
						list []interface{}
					}{newKey, listVal})
					continue // Skip normal processing for this key for now
				}
			}

			// Process normally (recurse for maps, join lists, handle primitives)
			// This directly modifies tempItems or calls recursively
			processSingleValue(v, newKey, sep, listSep, tempItems, allRows)
		}

		// Merge items flattened at this level into the current base
		for k, v := range tempItems {
			currentLevelBase[k] = v
		}

		// 2. Handle Expansion (if any target was found)
		if len(expansionTargets) > 0 {
			target := expansionTargets[0] // Expand only the first list encountered at this level
			log.Printf("Expanding list at key: %s (%d items)\n", target.key, len(target.list))
			for _, listItem := range target.list {
				// For each item in the list, start recursion with a *copy*
				// of the base built *up to this point* (currentLevelBase).
				flattenRecursive(listItem, target.key, sep, listSep, cloneMap(currentLevelBase), allRows)
			}
		} else {
			// No expansion happened at this level. This map represents a complete
			// (potentially nested) part of a final row.
			// We only add to allRows if this call completes a branch originating
			// from a top-level item or a list expansion.
			// Avoid adding intermediate baseRows during deep recursion.
			// Heuristic: If parentKey is empty (top level) or contains listSep (unlikely, check sep usage)
			// A better way might be to track if an expansion occured *below* this point.
			// Simpler: Only the caller that processes a list item adds the final row.
			// If this function call completes without triggering expansion below it,
			// it means this 'baseRow' represents a final state for its branch.
			// Let's add it here, the caller (list expansion loop) will handle its items.
			*allRows = append(*allRows, currentLevelBase) // Add the fully processed state for this branch
			// Register headers found in this completed row
			for k := range currentLevelBase {
				allHeaders[k] = struct{}{}
			}
		}

	case []interface{}:
		// This case is primarily for the top-level call if the input is an array.
		// If encountered nestedly, it's handled by processSingleValue.
		// Flatten each item in the top-level array independently.
		log.Println("Processing top-level array...")
		for _, item := range value {
			flattenRecursive(item, parentKey, sep, listSep, make(map[string]string), allRows)
		}

	default: // Primitive or nil at top level (less common)
		log.Printf("Warning: Top-level data is a primitive or nil: %v. Handling as single value.", value)
		newRow := make(map[string]string)
		key := "value" // Default key if no parent
		if parentKey != "" {
			key = parentKey
		}
		newRow[key] = fmt.Sprint(value) // Convert primitive to string
		allHeaders[key] = struct{}{}
		*allRows = append(*allRows, newRow)
	}
}

// processSingleValue handles flattening non-map values encountered within flattenRecursive
func processSingleValue(value interface{}, currentKey, sep, listSep string, currentRow map[string]string, allRows *[]map[string]string) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Nested map encountered - continue recursion, passing current state
		// This map doesn't represent a final row itself yet.
		flattenRecursive(v, currentKey, sep, listSep, currentRow, allRows)
	case []interface{}:
		// Handle array of primitives vs array of complex types
		if len(v) > 0 {
			primitives := make([]string, 0, len(v))
			isPrimitiveList := true
			for _, item := range v {
				switch item.(type) {
				case map[string]interface{}, []interface{}:
					isPrimitiveList = false
					break // Found nested structure
				default:
					primitives = append(primitives, fmt.Sprint(item))
				}
			}

			if isPrimitiveList {
				// Join primitives and add to current row
				currentRow[currentKey] = strings.Join(primitives, listSep)
				allHeaders[currentKey] = struct{}{}
			} else {
				// It's a list of complex types but not maps (e.g., list of lists)
				// Represent it as string or handle more specifically if needed
				log.Printf("Warning: Nested list at key '%s' contains non-primitive, non-map items. Converting to string.", currentKey)
				currentRow[currentKey] = fmt.Sprintf("%v", v) // Basic string representation
				allHeaders[currentKey] = struct{}{}
			}

		} else {
			// Empty list
			currentRow[currentKey] = "" // Represent empty list as empty string
			allHeaders[currentKey] = struct{}{}
		}
	case nil:
		currentRow[currentKey] = "" // Represent nil as empty string
		allHeaders[currentKey] = struct{}{}
	default: // Primitive type (string, number, bool)
		currentRow[currentKey] = fmt.Sprintf("%v", v)
		allHeaders[currentKey] = struct{}{}
	}
}

// cloneMap creates a shallow copy of a map. Needed for recursive calls.
func cloneMap(original map[string]string) map[string]string {
	newMap := make(map[string]string, len(original))
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

// --- Direct Conversion (for Non-Nested Data) ---

func convertToCSVDirect(data interface{}) ([]map[string]string, []string) {
	rows := make([]map[string]string, 0)
	headers := make(map[string]struct{})
	tempRow := make(map[string]string) // Reusable map for efficiency

	processItem := func(item interface{}) {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("Warning: Item in direct conversion is not an object/map: %T. Skipping.", item)
			return
		}
		// Clear tempRow before reuse IF maps are large often
		// Otherwise, overwriting is fine. Let's clear for safety.
		for k := range tempRow {
			delete(tempRow, k)
		}

		for key, val := range itemMap {
			headers[key] = struct{}{} // Collect header
			// Convert value to string - ignore nested structures in direct mode
			switch v := val.(type) {
			case map[string]interface{}, []interface{}:
				tempRow[key] = "" // Or "[nested]" or skip
				log.Printf("Warning: Skipping nested structure at key '%s' in direct conversion mode.", key)
			case nil:
				tempRow[key] = ""
			default:
				tempRow[key] = fmt.Sprintf("%v", v)
			}
		}
		rows = append(rows, cloneMap(tempRow)) // Add a copy
	}

	switch v := data.(type) {
	case []interface{}: // Array of objects
		for _, item := range v {
			processItem(item)
		}
	case map[string]interface{}: // Single top-level object
		processItem(v)
	default:
		log.Printf("Warning: Direct conversion input is not an array or object: %T", data)
		// Handle primitive - create single row/cell?
		tempRow["value"] = fmt.Sprintf("%v", v)
		headers["value"] = struct{}{}
		rows = append(rows, cloneMap(tempRow))

	}

	// Convert header set to list
	headerList := make([]string, 0, len(headers))
	for h := range headers {
		headerList = append(headerList, h)
	}
	// Note: Headers are not sorted here, main function sorts the combined list

	return rows, headerList
}
