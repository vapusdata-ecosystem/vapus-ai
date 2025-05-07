package pkgs

import (
	"errors"
)

var (
	// Error constants for Postgres operations
	ErrPostgresConnection      = errors.New("error while connecting to Postgres")
	ErrInvalidModelStruct      = errors.New("invalid model struct for creating table, must be a pointer")
	ErrInvalidDataTablesOpts   = errors.New("invalid data tables options")
	ErrLoginSalesForceInstance = errors.New("error while logging into salesforce instance")
	ErrMongoDBConnection  = errors.New("error initializing mongoDB client")
	ErrGithubConnection      = errors.New("error while connecting to Github")
)
