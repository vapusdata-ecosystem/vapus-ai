package alloydb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN_TEMPLATE = "host=%s port=%d user=%s password=%s dbname=%s sslmode=require"

type AlloyDbOpts struct {
	// PostgresConfig is the configuration for the Postgres
	URL, Username, Password, Database, Schema string
	Port                                      int
	WithPool                                  bool
}

type AlloyDbStore struct {
	Opts   *AlloyDbOpts
	Conn   *sql.DB
	DB     *bun.DB
	logger zerolog.Logger
	Pool   *pgxpool.Pool
	Orm    *gorm.DB
}

func NewAlloyDbStore(opts *AlloyDbOpts, l zerolog.Logger) (*AlloyDbStore, error) {
	dsn := getDsn(opts)
	dbs, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrPostgresConnection, err)
	}
	conn, err := dbs.DB()
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrPostgresConnection, err)
	}
	return &AlloyDbStore{
		Opts:   opts,
		Conn:   conn,
		logger: l,
		Orm:    dbs,
	}, nil
}

func (m *AlloyDbStore) Close() {
	m.Conn.Close()
	m.DB = nil
	m.Orm = nil
}

func getDsn(opts *AlloyDbOpts) string {
	// build dsn
	localTime := time.Now()
	localTimeZone := localTime.Location().String()
	log.Println("Local Time Zone ----->>>>>>>>>>>>>> ", localTimeZone)
	return fmt.Sprintf(DSN_TEMPLATE, opts.URL, opts.Port, opts.Username, opts.Password, opts.Database)
}
