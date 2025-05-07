package charts

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/opts"
)

func GroupAndAggregate(instruction *ChartCallParams, utils *ChartUtils) ([]map[string]any, error) {
	utils.Logger.Info().Msgf("Starting grouping and aggregation with method '%s' on field '%s' and groupby fields - %v",
		instruction.AggregateMethod,
		instruction.ValueField,
		instruction.GroupByFields)

	// Map to hold group keys to slice of float64 values.
	groupedValues := make(map[string][]float64)

	if instruction.ValueField == "" {
		instruction.ValueField = instruction.YAxisField
	}
	// Iterate through each record.
	for i, record := range instruction.Dataset {
		var keyParts []string
		// Build a composite key using all the groupByFields.
		for _, field := range instruction.GroupByFields {
			if val, exists := record[field]; exists {
				keyParts = append(keyParts, fmt.Sprintf("%v", val))
			} else {
				utils.Logger.Error().Msg(fmt.Sprintf("record %d: missing group by field '%s'", i, field))
				return nil, fmt.Errorf("record %d: missing group by field '%s'", i, field)
			}
		}
		groupKey := strings.Join(keyParts, "|")

		// For "count", we simply add a dummy value.
		if instruction.AggregateMethod == "count" {
			groupedValues[groupKey] = append(groupedValues[groupKey], 1.0)
			continue
		}

		// For "sum" or "avg", ensure the valueField exists and is convertible.
		if rawVal, exists := record[instruction.ValueField]; exists {
			if num, ok := toFloat(rawVal); ok {
				groupedValues[groupKey] = append(groupedValues[groupKey], num)
			} else {
				errMsg := fmt.Sprintf("record %d: invalid value for field '%s': %v", i, instruction.ValueField, rawVal)
				utils.Logger.Error().Msg(errMsg)
				return nil, fmt.Errorf(errMsg)
			}
		} else {
			utils.Logger.Error().Msg(fmt.Sprintf("record %d: missing value field '%s'", i, instruction.ValueField))
			return nil, fmt.Errorf("record %d: missing value field '%s'", i, instruction.ValueField)
		}
	}

	// Process grouped data to compute aggregation.
	var results []map[string]any
	for groupKey, values := range groupedValues {
		var aggregated float64
		switch instruction.AggregateMethod {
		case ChartDataSumMethod:
			for _, v := range values {
				aggregated += v
			}
		case ChartDataAvgMethod:
			var sum float64
			for _, v := range values {
				sum += v
			}
			if len(values) > 0 {
				aggregated = sum / float64(len(values))
			}
		case ChartDataCountMethod:
			aggregated = float64(len(values))
		default:
			return nil, fmt.Errorf("unsupported aggregation method: %s", instruction.AggregateMethod)
		}

		results = append(results, map[string]any{
			"group": groupKey,
			"value": aggregated,
		})
	}
	utils.Logger.Info().Msg("Grouping and aggregation completed")
	return results, nil
}

func toFloat(v interface{}) (float64, bool) {
	var n float64
	var valid bool
	switch val := v.(type) {
	case float64:
		n = float64(val)
		valid = true
	case float32:
		n = float64(val)
		valid = true
	case int:
		n = float64(val)
		valid = true
	case int32:
		n = float64(val)
		valid = true
	case int64:
		n = float64(val)
		valid = true
	case uint:
		n = float64(val)
		valid = true
	case uint32:
		n = float64(val)
		valid = true
	case uint64:
		if val <= math.MaxUint64 {
			n = float64(val)
			valid = true
		}
	case string:
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			n = float64(num)
			valid = true
		} else {
			valid = false
		}
	default:
		valid = false
	}
	if valid {
		return math.Round(float64(n)*100) / 100, true
	} else {
		return 0, false
	}

}

func validFloat(v interface{}) (float64, bool) {
	var n float64
	var valid bool
	switch val := v.(type) {
	case float64:
		n = float64(val)
		valid = true
	case float32:
		n = float64(val)
		valid = true
	default:
		valid = false
	}
	if valid {
		return math.Round(float64(n)*100) / 100, true
	} else {
		return 0, false
	}
}

func AggregateScatterData(instruction *ChartCallParams) ([]opts.ScatterData, error) {
	// Use a nested map with keys as strings for consistent grouping.
	aggregation := make(map[string]map[string]float64)

	for _, row := range instruction.Dataset {
		// Get values from the record.
		xRaw, ok1 := row[instruction.XAxisField]
		yRaw, ok2 := row[instruction.YAxisField]
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("fields %s or %s not found in data", instruction.XAxisField, instruction.YAxisField)
		}

		if val, ok := validFloat(xRaw); ok {
			xRaw = val
		}

		if val, ok := validFloat(yRaw); ok {
			yRaw = val
		}

		// Use string representations as keys.
		xKey := fmt.Sprintf("%v", xRaw)
		yKey := fmt.Sprintf("%v", yRaw)

		if aggregation[xKey] == nil {
			aggregation[xKey] = make(map[string]float64)
		}
		aggregation[xKey][yKey]++
	}

	// Collect sorted x keys for consistent output.
	var xKeys []string
	for k := range aggregation {
		xKeys = append(xKeys, k)
	}
	sort.Strings(xKeys)

	var scatterData []opts.ScatterData

	// Iterate over sorted x keys.
	for _, xKey := range xKeys {
		// Attempt to convert xKey to float64.
		xFloat, err := strconv.ParseFloat(xKey, 64)
		if err != nil {
			// If conversion fails, skip this key.
			continue
		}
		// For each xKey, sort the y keys.
		var yKeys []string
		for yKey := range aggregation[xKey] {
			yKeys = append(yKeys, yKey)
		}
		sort.Strings(yKeys)
		for _, yKey := range yKeys {
			// Attempt to convert yKey to float64.
			yFloat, err := strconv.ParseFloat(yKey, 64)
			if err != nil {
				continue
			}
			count := aggregation[xKey][yKey]
			// Append a scatter data point.
			scatterData = append(scatterData, opts.ScatterData{
				Value: [3]interface{}{xFloat, yFloat, count},
			})
		}
	}

	return scatterData, nil
}
