package charts

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/apache/arrow/go/v18/arrow"
	"github.com/apache/arrow/go/v18/arrow/array"
	"github.com/apache/arrow/go/v18/arrow/memory"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// **Generic Function: Group, Aggregate & Prepare Data for Chart**
func ProcessDataForChart(data []map[string]any, groupByKey string, aggKey string, aggType string) (map[string]any, error) {
	// Apache Arrow memory allocator
	pool := memory.NewGoAllocator()

	// Define Arrow Schema
	fields := []arrow.Field{
		{Name: groupByKey, Type: arrow.BinaryTypes.String},
		{Name: aggKey, Type: arrow.PrimitiveTypes.Float64},
	}
	schema := arrow.NewSchema(fields, nil)

	// Build Arrow Record
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	// Insert data into Arrow Table
	for _, row := range data {
		groupVal, _ := row[groupByKey].(string)
		aggVal, _ := ConvertToFloat(row[aggKey])

		builder.Field(0).(*array.StringBuilder).Append(groupVal)
		builder.Field(1).(*array.Float64Builder).Append(aggVal)
	}

	// Convert to Record
	record := builder.NewRecord()
	defer record.Release()

	// Perform Grouping & Aggregation
	groupedData := make(map[string]float64)
	counts := make(map[string]int)

	for i := 0; i < int(record.NumRows()); i++ {
		group := record.Column(0).(*array.String).Value(i)
		value := record.Column(1).(*array.Float64).Value(i)

		// Aggregation Logic
		groupedData[group] += value
		counts[group]++
	}

	// Compute Average if needed
	if aggType == "avg" {
		for k := range groupedData {
			groupedData[k] /= float64(counts[k])
		}
	}

	// Sort Results
	sortedKeys := sortedKeys(groupedData)

	// Prepare final data for charting
	result := map[string]any{
		"x_axis": sortedKeys,
		"y_axis": extractSortedValues(groupedData, sortedKeys),
	}
	return result, nil
}

// **Convert Values to Float64**
func ConvertToFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported type")
	}
}

// **Sort Keys for Consistent Chart Order**
func sortedKeys(data map[string]float64) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// **Extract Sorted Values**
func extractSortedValues(data map[string]float64, sortedKeys []string) []float64 {
	values := make([]float64, len(sortedKeys))
	for i, key := range sortedKeys {
		values[i] = data[key]
	}
	return values
}

// // **Generate Line Chart using Go-ECharts**
// func renderChart(w http.ResponseWriter, data map[string]any) {
// 	xLabels := data["x_axis"].([]string)
// 	yValues := data["y_axis"].([]float64)

// 	// Create Line Chart
// 	line := charts.NewLine()
// 	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Aggregated Data"}))
// 	line.SetXAxis(xLabels).AddSeries("Data", generateLineSeries(yValues))

// 	// Render Chart
// 	line.Render(w)
// }

// **Helper: Convert Y-Axis Data to Chart Series**
func generateLineSeries(yValues []float64) []opts.LineData {
	series := []opts.LineData{}
	for _, v := range yValues {
		series = append(series, opts.LineData{Value: v})
	}
	return series
}
