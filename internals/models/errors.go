package models

import "errors"

var (
	ErrDataSourceCredsNotFound = errors.New("no credentials found for the access scope")
)
