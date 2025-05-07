package filetools

import (
	"encoding/json"
	"log"
	os "os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"gopkg.in/yaml.v3"
)

func ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func FileDatasetLoader(bytes []byte, format string, isColumnerData bool) (*options.DataSetSummary, error) {
	var result = []map[string]any{}
	var columnnerData = make([][]string, 0)
	var headers = make([]string, 0)
	dataLength := 0
	var err error
	log.Println("IsColumnerData =========================================================", isColumnerData)
	switch strings.ToLower(format) {
	case "yaml":
		err = yaml.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	case "json":
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	case "toml":
		err = toml.Unmarshal(bytes, &result)
		if err != nil {
			return nil, err
		}
	case "csv":
		if isColumnerData {
			columnnerData, headers, err = CSVBytesToArrayColumn(bytes)
			if err != nil {
				return nil, err
			}
		} else {
			result, err = CSVBytesToArrayOfMap(bytes)
			if err != nil {
				return nil, err
			}
		}

	default:
		return nil, dmerrors.ErrInvalidArgs
	}
	if isColumnerData {
		dataLength = len(columnnerData)
	} else {
		if len(result) > 0 {
			for key := range result[0] {
				headers = append(headers, key)
			}
		}
		dataLength = len(result)
	}
	log.Println("Data Length =========================================================", dataLength)
	log.Println("Headers =========================================================", headers)
	return &options.DataSetSummary{
		DataFields:     headers,
		ResultMap:      result,
		ColumnsArray:   columnnerData,
		IsColumnerData: isColumnerData,
		ResultLength:   int64(dataLength),
	}, err
}
