package databricks

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/databricks/databricks-sdk-go"
	_ "github.com/databricks/databricks-sql-go"
)

var DSN_TEMPLATE = "token:%s@%s:%v%s"

type DatabricksOpts struct {
	Host  string
	Token string
	Path  string
	Port  int
}

type DatabricksStore struct {
	Opts   *DatabricksOpts
	Client *databricks.WorkspaceClient
	logger zerolog.Logger
	DB     *sql.DB
}

// Ex:- token:dapi12345678901234567890123456789012@dbc-a1b2345c-d6e7.cloud.databricks.com:443/sql/1.0/endpoints/a1b234c5678901d2

func NewConnectDatabricks(opts *DatabricksOpts, logger zerolog.Logger) (*DatabricksStore, error) {
	dsn := getDsn(opts)
	logger.Debug().Msgf("Connecting to Databricks with DSN: %s", dsn)

	db, err := sql.Open("databricks", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Databricks: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping Databricks database: %w", err)
	}

	logger.Info().Msg("Connected to Databricks Datalake successfully")

	return &DatabricksStore{
		Opts:   opts,
		DB:     db,
		logger: logger,
	}, nil
}

func (m *DatabricksStore) Close() {
	m.DB.Close()
}

func getDsn(opts *DatabricksOpts) string {
	// Token, Host, Http Path
	return fmt.Sprintf(DSN_TEMPLATE, opts.Token, opts.Host, opts.Port, opts.Path)
}
