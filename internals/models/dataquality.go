package models

type DataQualityMetric struct {
	Completeness  float64  `json:"completeness"`   // Completeness rate as a percentage
	Validity      float64  `json:"validity"`       // Validity rate as a percentage
	Accuracy      float64  `json:"accuracy"`       // Accuracy rate as a percentage
	Consistency   float64  `json:"consistency"`    // Consistency rate as a percentage
	Uniqueness    float64  `json:"uniqueness"`     // Uniqueness rate as a percentage
	Timeliness    string   `json:"timeliness"`     // Timeliness status ("Up-to-date", "Outdated")
	LastEvaluated int64    `json:"last_evaluated"` // When the quality was last evaluated
	Alerts        []string `json:"alerts"`         // List of alerts triggered for this data product
}

// func (x *DataQualityMetric) ConvertToPb() *pb.DataQualityMetric {
// 	return &pb.DataQualityMetric{
// 		Completeness:  x.Completeness,
// 		Validity:      x.Validity,
// 		Accuracy:      x.Accuracy,
// 		Consistency:   x.Consistency,
// 		Uniqueness:    x.Uniqueness,
// 		Timeliness:    x.Timeliness,
// 		LastEvaluated: x.LastEvaluated,
// 		Alerts:        x.Alerts,
// 	}
// }

func calculateDataQualityMetrics(nonMissingValues, totalValues, validEntries, totalEntries, accurateEntries, consistentRecords, uniqueRecords, totalRecords int, timelinessThreshold int64, lastUpdated int64) DataQualityMetric {
	metric := DataQualityMetric{}

	// Completeness
	metric.Completeness = (float64(nonMissingValues) / float64(totalValues)) * 100

	// Validity
	metric.Validity = (float64(validEntries) / float64(totalEntries)) * 100

	// Accuracy
	metric.Accuracy = (float64(accurateEntries) / float64(totalEntries)) * 100

	// Consistency
	metric.Consistency = (float64(consistentRecords) / float64(totalRecords)) * 100

	// Uniqueness
	metric.Uniqueness = (float64(uniqueRecords) / float64(totalRecords)) * 100

	// Timeliness
	if lastUpdated >= timelinessThreshold {
		metric.Timeliness = "Up-to-date"
	} else {
		metric.Timeliness = "Outdated"
	}

	return metric
}
