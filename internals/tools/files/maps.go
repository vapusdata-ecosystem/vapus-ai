package filetools

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

var FileMimeMap = map[mpb.ContentFormats][]string{
	mpb.ContentFormats_YAML:       {"application/x-yaml"},
	mpb.ContentFormats_JSON:       {"application/json"},
	mpb.ContentFormats_TOML:       {"application/toml"},
	mpb.ContentFormats_CSV:        {"text/csv"},
	mpb.ContentFormats_HTML:       {"text/html"},
	mpb.ContentFormats_PDF:        {"application/pdf"},
	mpb.ContentFormats_PYTHON:     {"text/x-python", "application/x-python"},
	mpb.ContentFormats_GOLANG:     {"text/x-go", "application/x-go"},
	mpb.ContentFormats_JAVASCRIPT: {"application/javascript", "text/javascript"},
	mpb.ContentFormats_JPG:        {"image/jpeg"},
	mpb.ContentFormats_JPEG:       {"image/jpeg"},
	mpb.ContentFormats_PNG:        {"image/png"},
}
