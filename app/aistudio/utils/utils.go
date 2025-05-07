package utils

import (
	"encoding/json"
	"errors"
	"io"

	guuid "github.com/google/uuid"
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

func Int32Ptr(x int32) *int32 { return &x }

func Int64Ptr(x int64) *int64 { return &x }

func IntPtr(x int) *int { return &x }

func Float32Ptr(x float32) *float32 { return &x }

func Float64Ptr(x float64) *float64 { return &x }

func Bool2Ptr(x bool) *bool { return &x }

func GetSecretName(resource, resourceId, attribute string) string {
	if resource == "" {
		return guuid.NewString() + "_" + attribute
	}
	return resourceId + "_" + attribute
}
