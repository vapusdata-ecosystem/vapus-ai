package mysql

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var defaultDsnParams = "charset=utf8mb4&parseTime=True&loc=Local"
var DSN_TEMPLATE = "%s:%s@tcp(%s:%d)/%s?%s"

type MysqlOpts struct {
	URL      string
	Port     int
	Username string
	Password string
	Database string
}

type MysqlStore struct {
	Opts   *MysqlOpts
	Conn   *sql.DB
	DB     *gorm.DB
	logger zerolog.Logger
}

func NewMysqlStore(opts *MysqlOpts, l zerolog.Logger) (*MysqlStore, error) {
	dsn := getDsn(opts)
	l.Debug().Msgf("Connecting to mysql with dsn: %s", dsn)
	dbs, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, dmerrors.DMError(ErrMySQLConnection, err)
	}
	conn, err := dbs.DB()
	if err != nil {
		return nil, dmerrors.DMError(ErrMySQLConnection, err)
	}
	return &MysqlStore{
		Opts:   opts,
		Conn:   conn,
		logger: l,
		DB:     dbs,
	}, nil
}

func (m *MysqlStore) Close() {
	m.Conn.Close()
	m.DB = nil
}

func getDsn(opts *MysqlOpts) string {
	// build dsn
	return fmt.Sprintf(DSN_TEMPLATE, opts.Username, opts.Password, opts.URL, opts.Port, opts.Database, defaultDsnParams)
}
