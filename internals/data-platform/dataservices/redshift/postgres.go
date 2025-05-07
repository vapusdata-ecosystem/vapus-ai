package redshift

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Add this line

	"github.com/jackc/pgx/v5/pgxpool" // Add this line
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN_TEMPLATE = "host=%s user=%s password=%s dbname=%s port=%d %s"

type RedshiftOpts struct {
	// PostgresConfig is the configuration for the Postgres
	URL, Database, Schema string
	creds                 *models.GenericCredentialModel
	Port                  int
	WithPool              bool
}

type RedshiftStore struct {
	Opts   *RedshiftOpts
	Conn   *sql.DB
	logger zerolog.Logger
	Pool   *pgxpool.Pool
	Orm    *gorm.DB
}

func NewRedshiftStore(opts *RedshiftOpts, l zerolog.Logger) (*RedshiftStore, error) {
	dsn := getDsn(opts)
	dbs, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrPostgresConnection, err)
	}
	conn, err := dbs.DB()
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrPostgresConnection, err)
	}
	fmt.Println("Succesfully Connected")

	return &RedshiftStore{
		Opts:   opts,
		Conn:   conn,
		logger: l,
		Orm:    dbs,
	}, nil
}

func (m *RedshiftStore) Close() {
	m.Conn.Close()
	m.Orm = nil
}

func getDsn(opts *RedshiftOpts) string {
	// build dsn
	localTime := time.Now()
	localTimeZone := localTime.Location().String()
	log.Println("Local Time Zone ----->>>>>>>>>>>>>> ", localTimeZone)
	return fmt.Sprintf(DSN_TEMPLATE, opts.URL, opts.creds.Username, opts.creds.Password, opts.Database, opts.Port, "")
}
