package tools

import "errors"

var (
	ErrInvalidTools           = errors.New("invalid tools definition, tools should be a valid json schema")
	ErrInvalidFunctionSchema  = errors.New("invalid function schema, schema should be a valid json schema")
	ErrInvalidJsonSchemaBytes = errors.New("invalid json schema bytes")
	ErrInvalidJSONTools       = errors.New("invalid json tools definition")
)
