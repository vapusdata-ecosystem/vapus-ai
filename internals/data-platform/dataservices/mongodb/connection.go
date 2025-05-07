package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MONGO_TEMPLATE = "mongodb+srv://%s:%s@%s/?%s"

type MongoDBOpts struct {
	// PostgresConfig is the configuration for the Postgres
	URL, Username, Password, Database, Schema, ParamString string
	Port                                                   int
	WithPool                                               bool
}

type MongoDBStore struct {
	Opts   *MongoDBOpts
	Client *mongo.Client
	logger zerolog.Logger
}

func NewMongoDBStore(opts *MongoDBOpts, l zerolog.Logger) (*MongoDBStore, error) {
	uri := getDsn(opts)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	optss := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), optss)
	if err != nil {
		return nil, dmerrors.DMError(pkgs.ErrMongoDBConnection, err)
	}
	return &MongoDBStore{
		Client: client,
	}, nil
}

func (m *MongoDBStore) Close(ctx context.Context) {
	m.Client.Disconnect(ctx)
}

func getDsn(opts *MongoDBOpts) string {
	// build dsn
	localTime := time.Now()
	localTimeZone := localTime.Location().String()
	log.Println("Local Time Zone ----->>>>>>>>>>>>>> ", localTimeZone)
	return fmt.Sprintf(MONGO_TEMPLATE, opts.Username, opts.Password, opts.URL, opts.ParamString)
}
