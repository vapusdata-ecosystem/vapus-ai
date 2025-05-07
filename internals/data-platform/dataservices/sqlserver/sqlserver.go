package sqlserver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/rs/zerolog"
)

// "Server=your-rds-endpoint.region.rds.amazonaws.com,1433;Database=your_database_name;User Id=your_username;Password=your_password;Encrypt=True;TrustServerCertificate=False;"

// var rds_endpoint = "microsoft-mysql-database.chaoqk4cybzf.ap-south-1.rds.amazonaws.com"
var DSN_TEMPLATE = "Server=%s,%d;Database=%s;User Id=%s;Password=%s;"

// rds_endpoint, Port, your_database_name, your_username, your_password
type SqlServerOpts struct {
	URL      string
	Port     int
	Username string
	Password string
	Database string
}

type SqlServerStore struct {
	Opts *SqlServerOpts
	// Conn
	DB     *sql.DB
	logger zerolog.Logger
}

func NewSqlServerStore(opts *SqlServerOpts, l zerolog.Logger) (*SqlServerStore, error) {

	dsn := getDsn(opts)
	// data Source name
	l.Debug().Msgf("Connecting to Sql Server with dsn: %s", dsn)

	// driverSourceName := "mssql"
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// defer db.Close()

	// Ping the database to check connectivity
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	l.Info().Msg("Connected to SqlServer database successfully")
	// err = SetupDatabase(db)
	// if err != nil {
	// 	log.Fatalf("Database setup failed: %v", err)
	// }

	// // Insert dummy users
	// err = InsertDummyUsers(db)
	// if err != nil {
	// 	log.Fatalf("Data insertion failed: %v", err)
	// }
	return &SqlServerStore{
		Opts: opts,
		// Conn:   conn,
		logger: l,
		DB:     db,
	}, nil
}

func (m *SqlServerStore) Close() {
	m.DB.Close()
	// m.DB = nil
}

func getDsn(opts *SqlServerOpts) string {
	// build dsn
	// rds_endpoint, Port, your_database_name, your_username, your_password
	return fmt.Sprintf(DSN_TEMPLATE, opts.URL, opts.Port, opts.Database, opts.Username, opts.Password)
}
