package snowflake

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
	_ "github.com/snowflakedb/gosnowflake"
	// "gorm.io/driver/sql"
)

var DSN_TEMPLATE = "%s:%s@%s/%s"

// "username:password@AccountName/DataBase"
// SnowflakeOpts holds the Snowflake connection options
type SnowflakeOpts struct {
	URL         string
	Port        int
	Username    string
	Password    string
	AccountName string
	Database    string
	Schema      string
	Warehouse   string
}

// SnowflakeStore represents a Snowflake database connection wrapped with GORM
type SnowflakeStore struct {
	Opts *SnowflakeOpts
	// Conn   *gorm.DB
	DB         *sql.DB
	SchemaName *string
	logger     zerolog.Logger
}

// NewConnectSnowflake initializes a new SnowflakeStore with GORM
func NewConnectSnowflake(opts *SnowflakeOpts, l zerolog.Logger) (*SnowflakeStore, error) {
	// fmt.Println("I am in NewConnectSnowflake: ", &opts)
	dsn := getDsn(opts)
	l.Debug().Msgf("Connecting to Snowflake with DSN: %s", dsn)

	// Open a Snowflake connection using gosnowflake driver
	db, err := sql.Open("snowflake", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Snowflake: %w", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping Snowflake database: %w", err)
	}

	l.Info().Msg("Connected to Snowflake database successfully")
	// fmt.Println(db)

	// Return the SnowflakeStore instance
	return &SnowflakeStore{
		Opts: opts,
		// Conn:   gormDB,
		DB:         db,
		logger:     l,
		SchemaName: &opts.Schema,
	}, nil
}

func (m *SnowflakeStore) Close() {
	m.DB.Close()
}

// getDsn generates the DSN string for Snowflake
func getDsn(opts *SnowflakeOpts) string {
	return fmt.Sprintf(DSN_TEMPLATE, opts.Username, opts.Password, opts.URL, opts.Database)
}
