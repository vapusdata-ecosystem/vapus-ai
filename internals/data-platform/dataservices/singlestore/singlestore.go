package singlestore

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	// "gorm.io/driver/sql"
)

// var defaultDsnParams = "charset=utf8mb4&parseTime=True&loc=Local"
var DSN_TEMPLATE = "%s:%s@tcp(%s:%v)/%s?tls=skip-verify"

// Username, Password Host Port  DB
// "username:password@AccountName/DataBase"

type SingleStoreOpts struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type SingleStoreStore struct {
	Opts *SingleStoreOpts
	// Conn   *gorm.DB
	DB         *sql.DB
	SchemaName *string
	logger     zerolog.Logger
}

func NewConnectSingleStore(opts *SingleStoreOpts, l zerolog.Logger) (*SingleStoreStore, error) {
	dsn := getDsn(opts)
	l.Debug().Msgf("Connecting to Single Store with DSN: %s", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Single Store: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping Single Store database: %w", err)
	}

	l.Info().Msg("Connected to Single Store database successfully")

	return &SingleStoreStore{
		Opts:   opts,
		DB:     db,
		logger: l,
	}, nil
}

func (m *SingleStoreStore) Close() {
	m.DB.Close()
	m.DB = nil
}

// getDsn generates the DSN string for SingleStore
func getDsn(opts *SingleStoreOpts) string {
	// Username, Password Host Port DB
	return fmt.Sprintf(DSN_TEMPLATE, opts.Username, opts.Password, opts.Host, opts.Port, opts.Database)
}

// Elastic Search,  Maria,  sqlserver
