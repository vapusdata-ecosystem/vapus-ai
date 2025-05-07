package apachearrow

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"google.golang.org/protobuf/types/known/structpb"
)

func ConvertStructpbToApacheArrow(mapList []map[string]any) ([]array.Interface, *arrow.Schema, error) {
	var columns []string
	for _, s := range mapList {
		for k := range s {
			columns = append(columns, k)
		}
		break
	}
	if len(columns) == 0 {
		return nil, &arrow.Schema{}, fmt.Errorf("no columns found")
	}
	var arrowFields []arrow.Field
	for _, c := range columns {
		arrowFields = append(arrowFields, arrow.Field{Name: c, Type: arrow.BinaryTypes.String})
	}
	schema := arrow.NewSchema(arrowFields, nil)
	var arrays []array.Interface
	for _, column := range columns {
		var values []string
		for _, s := range mapList {
			if v, ok := s[column]; ok {
				values = append(values, fmt.Sprintf("%v", v))
			} else {
				values = append(values, "")
			}
		}
		arrowArray := array.NewStringData(array.NewData(arrow.BinaryTypes.String, len(values), nil, nil, 0, 0))
		arrays = append(arrays, arrowArray)

	}
	return arrays, schema, nil
}

func ConvertMapToApacheArrow(structs []*structpb.Struct) ([]array.Interface, *arrow.Schema, error) {
	var columns []string
	for _, s := range structs {
		for k := range s.Fields {
			columns = append(columns, k)
		}
		break
	}
	if len(columns) == 0 {
		return nil, &arrow.Schema{}, fmt.Errorf("no columns found")
	}
	var arrowFields []arrow.Field
	for _, c := range columns {
		arrowFields = append(arrowFields, arrow.Field{Name: c, Type: arrow.BinaryTypes.String})
	}
	schema := arrow.NewSchema(arrowFields, nil)
	var arrays []array.Interface
	for _, column := range columns {
		var values []string
		for _, s := range structs {
			if v, ok := s.Fields[column]; ok {
				values = append(values, fmt.Sprintf("%v", v))
			} else {
				values = append(values, "")
			}
		}
		arrowArray := array.NewStringData(array.NewData(arrow.BinaryTypes.String, len(values), nil, nil, 0, 0))
		arrays = append(arrays, arrowArray)
	}
	return arrays, schema, nil
}

func arrowToCSV(columns []array.Interface, schema *arrow.Schema) ([]byte, error) {
	record := array.NewRecord(schema, columns, 3)
	defer record.Release()
	// Create a buffer to hold CSV data
	buffer := &bytes.Buffer{}

	// Create a CSV writer
	writer := csv.NewWriter(buffer)

	// Write header row (column names)
	header := make([]string, len(record.Schema().Fields()))
	for i, field := range record.Schema().Fields() {
		header[i] = field.Name
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write rows
	for row := 0; row < int(record.NumRows()); row++ {
		line := make([]string, record.NumCols())
		for col := 0; col < int(record.NumCols()); col++ {
			column := record.Column(col)
			switch column.DataType().ID() {
			case arrow.INT32:
				line[col] = fmt.Sprintf("%d", column.(*array.Int32).Value(row))
			case arrow.STRING:
				line[col] = column.(*array.String).Value(row)
			default:
				line[col] = fmt.Sprintf("%v", column.Data())
			}
		}
		if err := writer.Write(line); err != nil {
			return nil, err
		}
	}

	// Flush the writer and return bytes
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
