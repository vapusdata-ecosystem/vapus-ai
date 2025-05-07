package sqlops

import (
	"github.com/rs/zerolog"
	"vitess.io/vitess/go/vt/sqlparser"
)

type SQLOperator struct {
	logger zerolog.Logger
	parser *sqlparser.Parser
}

func New(logger zerolog.Logger) (*SQLOperator, error) {
	p, err := sqlparser.New(sqlparser.Options{
		TruncateErrLen: 50,
	})
	return &SQLOperator{
		logger: logger,
		parser: p,
	}, err
}
