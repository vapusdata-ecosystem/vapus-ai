package charts

import (
	"fmt"
	"log"
	"reflect"
)

func aggregateHeatmapData(data []map[string]any, xField string, yField string) (map[any]map[any]float64, error) {
	aggregation := make(map[interface{}]map[interface{}]float64)
	for _, item := range data {
		xValue, xExists := item[xField]
		yValue, yExists := item[yField]

		if !xExists || !yExists {
			return nil, fmt.Errorf("fields %s or %s not found in data", xField, yField)
		}

		xVal := reflect.ValueOf(xValue).Interface()
		yVal := reflect.ValueOf(yValue).Interface()

		if _, exists := aggregation[xVal]; !exists {
			aggregation[xVal] = make(map[any]float64)
		}

		aggregation[xVal][yVal]++
	}
	log.Println("aggregation >>>>>>>: ", aggregation)
	return aggregation, nil
}

func groupAndAggregate(data []map[string]any, groupByFields []string, valueField, method string, utils *ChartUtils) ([]map[string]any, error) {
	grouped := make(map[string][]float64)
	utils.Logger.Info().Msg("Grouping and aggregating data started")
	for _, record := range data {
		groupKey := ""
		for _, field := range groupByFields {
			if val, ok := record[field]; ok {
				groupKey += fmt.Sprintf("%v|", val)
			} else {
				return nil, fmt.Errorf("missing group by field: %s", field)
			}
		}
		if method == "count" {
			grouped[groupKey] = append(grouped[groupKey], 1.0)
		} else if val, ok := record[valueField]; ok {
			if num, valid := toFloat(val); valid {
				grouped[groupKey] = append(grouped[groupKey], num)
			} else {
				utils.Logger.Error().Msg("Invalid value for field " + valueField + ": " + fmt.Sprintf("%v", val))
				return nil, fmt.Errorf("invalid value for field %s: %v", valueField, val)
			}
		} else {
			return nil, fmt.Errorf("missing value field: %s", valueField)
		}
	}
	var result []map[string]any
	for groupKey, values := range grouped {
		var aggregatedValue float64
		switch method {
		case "sum":
			for _, v := range values {
				aggregatedValue += v
			}
		case "avg":
			var sum float64
			for _, v := range values {
				sum += v
			}
			aggregatedValue = sum / float64(len(values))
		case "count":
			aggregatedValue = float64(len(values))
		default:
			return nil, fmt.Errorf("unsupported aggregation method: %s", method)
		}
		result = append(result, map[string]any{
			"group": groupKey,
			"value": aggregatedValue,
		})
	}
	log.Println("Grouping and aggregating data completed")
	log.Println("Result: ", result)
	return result, nil
}
