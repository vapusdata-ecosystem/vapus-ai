package mysql

import (
	"errors"
)

var (
	// Error constants for MySQL operations
	ErrMySQLConnection  = errors.New("error while connecting to MySQL")
	ErrMySQLDBList      = errors.New("error while listing databases in MySQL")
	ErrMySQLTableList   = errors.New("error while listing tables in MySQL")
	ErrMySQLColumnList  = errors.New("error while listing columns in MySQL")
	ErrMySQLInsertQuery = errors.New("error while preparing insert query for MySQL")
)
