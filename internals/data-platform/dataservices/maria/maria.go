package maria

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rs/zerolog"
)

var defaultDsnParams = "charset=utf8mb4&parseTime=True&loc=Local"
var DSN_TEMPLATE = "%s:%s@tcp(%s:%d)/%s?%s"

// username, password, host, port, dbname, url
type MariaOpts struct {
	URL      string //host
	Port     int
	Username string
	Password string
	Database string
}

type MariaStore struct {
	Opts   *MariaOpts
	DB     *sql.DB
	logger zerolog.Logger
}

func NewMariaStore(opts *MariaOpts, l zerolog.Logger) (*MariaStore, error) {

	dsn := getDsn(opts)
	// data Source name
	l.Debug().Msgf("Connecting to Sql Server with dsn: %s", dsn)

	// driverSourceName := "mssql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// defer db.Close()
	// Ping the database to check connectivity
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	l.Info().Msg("Connected to Maria DB database successfully")
	// err = SetupDatabase(db)
	// if err != nil {
	// 	log.Fatalf("Database setup failed: %v", err)
	// }

	// // Insert dummy users
	// err = InsertDummyUsers(db)
	// if err != nil {
	// 	log.Fatalf("Data insertion failed: %v", err)
	// }
	return &MariaStore{
		Opts: opts,
		// Conn:   conn,
		logger: l,
		DB:     db,
	}, nil
}

func (m *MariaStore) Close() {
	m.DB.Close()
	// m.DB = nil
}

func getDsn(opts *MariaOpts) string {
	// build dsn
	// rds_endpoint, Port, your_database_name, your_username, your_password
	return fmt.Sprintf(DSN_TEMPLATE, opts.Username, opts.Password, opts.URL, opts.Port, opts.Database, defaultDsnParams)
}
