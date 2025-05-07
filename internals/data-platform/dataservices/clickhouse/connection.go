package clickhouse

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	// "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	"crypto/tls"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

var DSN_TEMPLATE = "host=%s user=%s password=%s dbname=%s port=%d %s"

type clickhouseOpts struct {
	// PostgresConfig is the configuration for the Postgres
	URL, Database, Schema string
	creds                 *models.GenericCredentialModel
	Port                  int
	WithPool              bool
}

type ClickhouseStore struct {
	Opts   *clickhouseOpts
	Conn   *sql.DB
	logger zerolog.Logger
	Pool   *pgxpool.Pool
	// Orm    *gorm.DB
}

func NewClickhouseStore(opts *clickhouseOpts, l zerolog.Logger) (*ClickhouseStore, error) {

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{opts.URL},
		Protocol: clickhouse.Native,
		TLS:      &tls.Config{}, // enable secure TLS
		Auth: clickhouse.Auth{
			Username: opts.creds.Username,
			Password: opts.creds.Password,
		},
	})

	return &ClickhouseStore{
		Opts:   opts,
		Conn:   conn,
		logger: l,
	}, nil
}

func (m *ClickhouseStore) Close() {
	m.Conn.Close()
	// m.Orm = nil
}

func getDsn(opts *clickhouseOpts) string {
	// build dsn
	localTime := time.Now()
	localTimeZone := localTime.Location().String()
	log.Println("Local Time Zone ----->>>>>>>>>>>>>> ", localTimeZone)
	return fmt.Sprintf(DSN_TEMPLATE, opts.URL, opts.creds.Username, opts.creds.Password, opts.Database, opts.Port, "")
}
