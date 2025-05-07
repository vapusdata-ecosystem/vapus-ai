package charts

// // ConvertToFloatCompute converts an interface{} value to float64.
// func ConvertToFloatCompute(val interface{}) (float64, error) {
// 	switch v := val.(type) {
// 	case int:
// 		return float64(v), nil
// 	case float64:
// 		return v, nil
// 	case string:
// 		return strconv.ParseFloat(v, 64)
// 	default:
// 		return 0, fmt.Errorf("unsupported type: %T", val)
// 	}
// }

// // ProcessDataForChartArrow aggregates data using Apache Arrow's vectorized compute functions.
// // It groups by groupByKey and aggregates the values from aggKey using the specified method ("sum" or "avg").
// // The function returns an array of maps, each representing one group.
// func ProcessDataForChartArrow(data []map[string]any, groupByKey string, aggKey string, aggType string) ([]map[string]any, error) {
// 	ctx := context.Background()
// 	pool := memory.NewGoAllocator()

// 	// Define Arrow schema with two fields: one string and one float64.
// 	fields := []arrow.Field{
// 		{Name: groupByKey, Type: arrow.BinaryTypes.String},
// 		{Name: aggKey, Type: arrow.PrimitiveTypes.Float64},
// 	}
// 	schema := arrow.NewSchema(fields, nil)

// 	// Build Arrow record using a RecordBuilder.
// 	builder := array.NewRecordBuilder(pool, schema)
// 	defer builder.Release()

// 	stringBuilder := builder.Field(0).(*array.StringBuilder)
// 	floatBuilder := builder.Field(1).(*array.Float64Builder)

// 	// Insert data into Arrow arrays.
// 	for _, row := range data {
// 		groupVal, ok := row[groupByKey].(string)
// 		if !ok {
// 			groupVal = fmt.Sprintf("%v", row[groupByKey])
// 		}
// 		aggVal, err := ConvertToFloatCompute(row[aggKey])
// 		if err != nil {
// 			return nil, fmt.Errorf("error converting value for key '%s': %w", aggKey, err)
// 		}
// 		stringBuilder.Append(groupVal)
// 		floatBuilder.Append(aggVal)
// 	}

// 	record := builder.NewRecord()
// 	defer record.Release()

// 	// Extract Arrow arrays.
// 	groupArray := record.Column(0).(*array.String)
// 	aggArray := record.Column(1).(*array.Float64)

// 	// Collect unique groups.
// 	uniqueGroupsMap := make(map[string]struct{})
// 	for i := 0; i < int(record.NumRows()); i++ {
// 		group := groupArray.Value(i)
// 		uniqueGroupsMap[group] = struct{}{}
// 	}
// 	var uniqueGroups []string
// 	for k := range uniqueGroupsMap {
// 		uniqueGroups = append(uniqueGroups, k)
// 	}
// 	sort.Strings(uniqueGroups)

// 	// Map to hold aggregated values per group.
// 	groupedResults := make(map[string]float64)

// 	// For each unique group, use Arrow compute to filter and aggregate.
// 	for _, group := range uniqueGroups {
// 		// Create a scalar from the group string using the scalar package.
// 		groupScalar := scalar.NewStringScalar(group)

// 		// Use vectorized equality to create a boolean mask.
// 		maskDatum, err := compute.CallFunction("equal", []compute.Datum{compute.NewDatum(groupArray), compute.NewDatum(groupScalar)}, nil)
// 		if err != nil {
// 			return nil, fmt.Errorf("error computing equality for group %s: %w", group, err)
// 		}

// 		maskArray, ok := maskDatum.(*array.Boolean)
// 		if !ok {
// 			return nil, fmt.Errorf("expected boolean array, got %T", maskDatum)
// 		}
// 		maskDatum.Release() // Release the datum after extracting the array.

// 		// Filter the aggregation array using the mask.
// 		filteredDatum, err := compute.Filter(ctx, aggArray, maskArray, nil)
// 		if err != nil {
// 			return nil, fmt.Errorf("error filtering aggArray for group %s: %w", group, err)
// 		}
// 		filteredArr, ok := filteredDatum.(*array.Float64)
// 		if !ok {
// 			return nil, fmt.Errorf("expected float64 array for filtered result, got %T", filteredDatum)
// 		}
// 		filteredDatum.Release()

// 		// Use vectorized sum function.
// 		sumDatum, err := compute.Sum(ctx, filteredArr)
// 		if err != nil {
// 			return nil, fmt.Errorf("error summing for group %s: %w", group, err)
// 		}
// 		sumScalar, ok := sumDatum.(compute.Scalar)
// 		if !ok {
// 			return nil, fmt.Errorf("expected a scalar for sum result, got %T", sumDatum)
// 		}
// 		sumVal := sumScalar.Data().(float64)
// 		sumDatum.Release()

// 		// Compute average if requested.
// 		if strings.ToLower(aggType) == "avg" {
// 			count := filteredArr.Len()
// 			if count > 0 {
// 				groupedResults[group] = sumVal / float64(count)
// 			} else {
// 				groupedResults[group] = 0
// 			}
// 		} else { // assume "sum" aggregation
// 			groupedResults[group] = sumVal
// 		}
// 		filteredArr.Release()
// 		maskArray.Release()
// 		// groupScalar is a scalar; call Release() if needed.
// 		groupScalar.Release()
// 	}

// 	// Build the output as an array of maps.
// 	var output []map[string]any
// 	for _, group := range uniqueGroups {
// 		output = append(output, map[string]any{
// 			"group": group,
// 			"value": groupedResults[group],
// 		})
// 	}

// 	return output, nil
// }
