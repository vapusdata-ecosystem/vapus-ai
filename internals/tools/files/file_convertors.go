package filetools

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/phpdave11/gofpdf"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

// ConvertJSONToCSV converts JSON records to a CSV file.
func ConvertJSONToCSV(records []map[string]any, outputPath string) error {
	if len(records) == 0 {
		return errors.New("no records to convert")
	}

	// Create a buffer to write CSV data
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Extract headers from the first record
	var headers []string
	for key := range records[0] {
		headers = append(headers, key)
	}

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write records
	for _, record := range records {
		var row []string
		for _, header := range headers {
			value := record[header]
			// Convert value to string
			strValue := fmt.Sprintf("%v", value)
			row = append(row, strValue)
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	// Flush the writer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	// Write buffer to file
	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// ConvertJSONToYAML converts JSON records to a YAML file.
func ConvertJSONToYAML(records []map[string]any, outputPath string) error {
	if len(records) == 0 {
		return errors.New("no records to convert")
	}

	// Marshal JSON records to YAML
	yamlData, err := yaml.Marshal(&records)
	if err != nil {
		return err
	}

	// Write YAML data to file
	return os.WriteFile(outputPath, yamlData, 0644)
}

// ConvertJSONToText converts JSON records to a plain text file.
func ConvertJSONToText(records []map[string]any, outputPath string) error {
	if len(records) == 0 {
		return errors.New("no records to convert")
	}

	// Marshal JSON records with indentation for readability
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON data as text to file
	return os.WriteFile(outputPath, jsonData, 0644)
}

// ConvertJSONToPDF converts JSON records to a PDF file with a table.
func ConvertJSONToPDF(records []map[string]any, outputPath string) error {
	if len(records) == 0 {
		return errors.New("no records to convert")
	}

	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)

	// Extract headers from the first record
	var headers []string
	for key := range records[0] {
		headers = append(headers, key)
	}

	// Set column widths
	colWidths := make([]float64, len(headers))
	for i := range colWidths {
		colWidths[i] = 40 // Set a default width; adjust as needed
	}

	// Write headers
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Write records
	pdf.SetFont("Arial", "", 12)
	for _, record := range records {
		for i, header := range headers {
			value := record[header]
			strValue := fmt.Sprintf("%v", value)
			// Truncate if too long
			if len(strValue) > 30 {
				strValue = strValue[:27] + "..."
			}
			pdf.CellFormat(colWidths[i], 10, strValue, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Output PDF to file
	return pdf.OutputFileAndClose(outputPath)
}

// ConvertJSONToXLSX converts JSON records to an XLSX (Excel) file.
func ConvertJSONToXLSX(records []map[string]any, outputPath string) error {
	if len(records) == 0 {
		return errors.New("no records to convert")
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Extract headers from the first record
	var headers []string
	for key := range records[0] {
		headers = append(headers, key)
	}

	// Write headers
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return err
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// Write records
	for rowIndex, record := range records {
		for colIndex, header := range headers {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return err
			}
			value := record[header]
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Save the file
	return f.SaveAs(outputPath)
}

// TODO: Add more file type detection functions
func ConvertCsvTo(fileType, expectedformat string, contentBytes []byte) ([]byte, error) {
	var fBytes []byte
	var err error
	if fileType == expectedformat {
		return contentBytes, nil
	}
	result := []map[string]any{}
	result, err = CSVBytesToArrayOfMap(contentBytes)
	if err != nil {
		return nil, err
	}
	log.Println("ConvertFileBytes: result: ", result)
	switch strings.ToLower(expectedformat) {
	case "json":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "yaml":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "csv":
		fBytes, err = MapArrayCSVMarshaler(result)
	default:
		fBytes, err = MapArrayCSVMarshaler(result)
	}
	if err != nil {
		return nil, err
	}
	return fBytes, nil
}

func ConvertFile(fileType, expectedformat string, contentBytes []byte) ([]byte, error) {
	var fBytes []byte
	var err error
	if fileType == expectedformat {
		return contentBytes, nil
	}
	result := []map[string]any{}
	if strings.ToLower(fileType) == "csv" {
		result, err = CSVBytesToArrayOfMap(contentBytes)
		if err != nil {
			return nil, err
		}
	} else {
		err = GenericUnMarshaler(contentBytes, &result, fileType)
		if err != nil {
			return nil, err
		}
	}
	log.Println("ConvertFileBytes: result: ", result)
	switch strings.ToLower(expectedformat) {
	case "json":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "yaml":
		fBytes, err = GenericMarshaler(result, strings.ToUpper(expectedformat))
	case "csv":
		fBytes, err = MapArrayCSVMarshaler(result)
	default:
		fBytes, err = MapArrayCSVMarshaler(result)
	}
	if err != nil {
		return nil, err
	}
	return fBytes, nil
}

func GetCSVColumnsFromStruct(s any) []string {
	var columns []string
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := range t.NumField() {
		field := t.Field(i)
		csvVal := field.Tag.Get("csv")
		if tagName := strings.Split(csvVal, ",")[0]; tagName != "" && tagName != "-" {
			columns = append(columns, tagName)
		} else if csvVal == "-" {
			continue
		} else {
			columns = append(columns, strings.ToLower(field.Name))
		}

	}

	return columns
}
