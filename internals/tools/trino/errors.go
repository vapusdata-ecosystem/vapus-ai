package trinocl

import (
	"errors"
)

var (
	ErrTrinoIsReadOnly = errors.New("trino is in read only mode")
)
