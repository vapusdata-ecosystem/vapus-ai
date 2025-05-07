package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	// Add this line
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool" // Add this line
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN_TEMPLATE = "host=%s user=%s password=%s dbname=%s port=%d %s"

type PostgresOpts struct {
	// PostgresConfig is the configuration for the Postgres
	URL, Username, Password, Database, Schema string
	Port                                      int
	WithPool                                  bool
}

type PostgresStore struct {
	Opts   *PostgresOpts
	Conn   *sql.DB
	DB     *bun.DB
	logger zerolog.Logger
	Pool   *pgxpool.Pool
	Orm    *gorm.DB
}

func NewPostgresStoreLocal(opts *PostgresOpts, l zerolog.Logger) (*PostgresStore, error) {
	dsn := getDsn(opts)
	l.Debug().Msgf("Connecting to postgres with dsn: %s", dsn)
	if opts.WithPool {
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}
		pool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			return nil, err
		}
		sqldb := stdlib.OpenDBFromPool(pool)
		db := bun.NewDB(sqldb, pgdialect.New())
		// db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		db.AddQueryHook(bundebug.NewQueryHook())
		return &PostgresStore{
			Opts:   opts,
			logger: l,
			Conn:   sqldb,
			DB:     db,
			Pool:   pool,
		}, nil
	} else {
		config, err := pgx.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}
		// bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
		sqldb := stdlib.OpenDB(*config)
		db := bun.NewDB(sqldb, pgdialect.New())
		db.AddQueryHook(bundebug.NewQueryHook())
		return &PostgresStore{
			Opts:   opts,
			logger: l,
			Conn:   sqldb,
			DB:     db,
		}, nil
	}
}

func NewPostgresStore(opts *PostgresOpts, l zerolog.Logger) (*PostgresStore, error) {
	dsn := getDsn(opts)
	dbs, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrPostgresConnection, err)
	}
	conn, err := dbs.DB()
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrPostgresConnection, err)
	}
	return &PostgresStore{
		Opts:   opts,
		Conn:   conn,
		logger: l,
		Orm:    dbs,
	}, nil
}

func (m *PostgresStore) Close() {
	m.Conn.Close()
	m.DB = nil
	m.Orm = nil
}

func getDsn(opts *PostgresOpts) string {
	// build dsn
	localTime := time.Now()
	localTimeZone := localTime.Location().String()
	log.Println("Local Time Zone ----->>>>>>>>>>>>>> ", localTimeZone)
	return fmt.Sprintf(DSN_TEMPLATE, opts.URL, opts.Username, opts.Password, opts.Database, opts.Port, "")
}
