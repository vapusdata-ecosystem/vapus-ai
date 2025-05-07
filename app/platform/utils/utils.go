package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

func EsGenericResponseReader(body io.ReadCloser) (map[string]interface{}, error) {
	var r map[string]interface{}
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}
	val, ok := r["error"]
	if ok {
		errM, ok := val.(map[string]interface{})["type"]
		if ok {
			return nil, errors.New(errM.(string))
		}
	}
	return r, nil
}

func GetFilterParams(query *mpb.SearchParam, key string) string {
	if query != nil {
		for _, param := range query.GetFilters() {
			if param.GetKey() == key {
				if len(param.GetValue()) > 0 {
					return param.GetValue()
				}

			}
		}
	}
	return ""
}

func GetSqlWhereList[T any](vals []T) string {
	result := ""
	for i, val := range vals {
		n := fmt.Sprintf("%v", val)
		if n != "" {
			if i == 0 {
				result = fmt.Sprintf("%v'%v'", result, n)
			} else {
				result = fmt.Sprintf("%v,'%v'", result, n)
			}
		}
	}
	return result
}
