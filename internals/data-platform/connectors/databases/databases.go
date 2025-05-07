package databases

import (
	"context"

	"github.com/rs/zerolog"
	// velasticsearch "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/elasticsearch"
	vdatabricks "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/databricks"
	vdynamodb "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/dynamodb"
	ves "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/elasticsearch"
	vgithubconnector "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/githubconnector"
	vmariadb "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/maria"
	vmongodb "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/mongodb"
	vmysql "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/mysql"
	vopensearch "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/opensearch"
	vpinecone "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pinecone"
	vpostgres "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/postgres"
	vqdrant "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/qdrant"
	vredis "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/redis"
	vsinglestore "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/singlestore"
	vsnowflake "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/snowflake"
	vsqlserver "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/sqlserver"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/types"

	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	logger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
)

type DataStoreClient struct {
	RedisClient         *vredis.RedisStore
	MysqlClient         *vmysql.MysqlStore
	PostgresClient      *vpostgres.PostgresStore
	PineconeClient      *vpinecone.PineconeStore
	QdrantClient        *vqdrant.QdrantStore
	MongoDBClient       *vmongodb.MongoDBStore
	ElasticSearchClient *ves.ElasticSearchStore
	DynamoDBClient      *vdynamodb.DynamoDBStore
	DatabricksClient    *vdatabricks.DatabricksStore
	SnowflakeClient     *vsnowflake.SnowflakeStore
	OpenSearchClient    *vopensearch.OpenSearchStore
	SingleStoreClient   *vsinglestore.SingleStoreStore
	SqlServerClient     *vsqlserver.SqlServerStore
	MariaDbClient       *vmariadb.MariaStore
	GithubClient        *vgithubconnector.GithubStore
	DataStoreParams     *models.DataSourceCredsParams
	Logger              zerolog.Logger
	Debug               bool
	InApp               bool
}

var (
	log = logger.CoreLogger
)

func (d *DataStoreClient) Close() {
	if d.RedisClient != nil {
		d.RedisClient.Close()
	}
	if d.MysqlClient != nil {
		d.MysqlClient.Close()
	}
	if d.PostgresClient != nil {
		d.PostgresClient.Close()
	}
	if d.OpenSearchClient != nil {
		d.OpenSearchClient.Close()
	}
	if d.SingleStoreClient != nil {
		d.SingleStoreClient.Close()
	}
	if d.GithubClient != nil {
		d.GithubClient.Close()
	}
	if d.DynamoDBClient != nil {
		d.DynamoDBClient.Close()
	}
	if d.DatabricksClient != nil {
		d.DatabricksClient.Close()
	}
	if d.PineconeClient != nil {
		d.PineconeClient.Close()
	}
	if d.QdrantClient != nil {
		d.QdrantClient.Close()
	}
}

type DataStoreOpts func(*DataStoreClient)

func WithDebug(debug bool) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.Debug = debug
	}
}

func WithInApp(inApp bool) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.InApp = inApp
	}
}

func WithLogger(log zerolog.Logger) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.Logger = log
	}
}

func WithDataSourceCredsParams(params *models.DataSourceCredsParams) DataStoreOpts {
	return func(d *DataStoreClient) {
		d.DataStoreParams = params
	}
}

// NewDataConnClient creates a new DataConnClient
func New(ctx context.Context, opts ...DataStoreOpts) (*DataStoreClient, error) {
	var err error
	dsc := &DataStoreClient{}
	for _, opt := range opts {
		opt(dsc)
	}
	if dsc.DataStoreParams == nil {
		return nil, dmerrors.DMError(dputils.ErrDataStoreParams404, dputils.ErrDataStoreParams404)
	}
	log.Debug().Msgf("Creating new data source connection for %s", dsc.DataStoreParams.DataSourceEngine)
	switch dsc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_REDIS.String():
		client, err := connectRedis(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to redis")
			return nil, err
		}
		return &DataStoreClient{RedisClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_MYSQL.String():
		client, err := connectMysql(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to mysql")
			return nil, err
		}
		return &DataStoreClient{MysqlClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_ELASTICSEARCH.String():
		client, err := connectElaticSearch(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to elasticsearch")
			return nil, err
		}
		return &DataStoreClient{ElasticSearchClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil

	case types.StorageEngine_POSTGRES.String():
		var client *vpostgres.PostgresStore
		crds := &vpostgres.PostgresOpts{
			URL:      dsc.DataStoreParams.DataSourceCreds.URL,
			Username: dsc.DataStoreParams.DataSourceCreds.GenericCredentialModel.Username,
			Password: dsc.DataStoreParams.DataSourceCreds.GenericCredentialModel.Password,
			Database: dsc.DataStoreParams.DataSourceCreds.DB,
			Port:     int(dsc.DataStoreParams.DataSourceCreds.Port),
		}
		if dsc.InApp {
			client, err = vpostgres.NewPostgresStoreLocal(crds, log)
		} else {
			client, err = vpostgres.NewPostgresStore(crds, log)
		}

		if err != nil {
			log.Err(err).Msg("Error connecting to postgres")
			return nil, err
		}
		return &DataStoreClient{PostgresClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_SNOWFLAKE.String():
		client, err := connectSnowflake(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Snowflake")
			return nil, err
		}
		return &DataStoreClient{SnowflakeClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_OPENSEARCH.String():
		client, err := connectOpenSearch(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Opensearch")
			return nil, err
		}
		return &DataStoreClient{OpenSearchClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_SINGLE_STORE.String():
		client, err := connectSingleStore(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Single Store")
			return nil, err
		}
		return &DataStoreClient{SingleStoreClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_SQL_SERVER.String():
		client, err := connectSqlServer(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Sql Server")
			return nil, err
		}
		return &DataStoreClient{SqlServerClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_MARIADB.String():
		client, err := connectMariaDb(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Maria DB")
			return nil, err
		}
		return &DataStoreClient{MariaDbClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_DYNAMODB.String():
		client, err := connectDynamoDB(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Dynamo DB")
			return nil, err
		}
		return &DataStoreClient{DynamoDBClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_DATABRICKS.String():
		client, err := connectDatabricks(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Databricks")
			return nil, err
		}
		return &DataStoreClient{DatabricksClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_PINECONE.String():
		client, err := connectPinecone(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Databricks")
			return nil, err
		}
		return &DataStoreClient{PineconeClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	case types.StorageEngine_QDRANT.String():
		client, err := connectQdrant(ctx, dsc.DataStoreParams.DataSourceCreds)
		if err != nil {
			log.Err(err).Msg("Error connecting to Databricks")
			return nil, err
		}
		return &DataStoreClient{QdrantClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil

	// case types.StorageEngine_AURORA_MYSQL.String():
	// 	client, err := connectMysql(ctx, dsc.DataStoreParams.DataSourceCreds)
	// 	if err != nil {
	// 		log.Err(err).Msg("Error connecting to mysql")
	// 		return nil, err
	// 	}
	// 	return &DataStoreClient{MysqlClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	// case types.StorageEngine_AURORA_POSTGRESS.String():
	// 	var client *vpostgres.PostgresStore
	// 	crds := &vpostgres.PostgresOpts{
	// 		URL:      dsc.DataStoreParams.DataSourceCreds.URL,
	// 		Username: dsc.DataStoreParams.DataSourceCreds.GenericCredentialModel.Username,
	// 		Password: dsc.DataStoreParams.DataSourceCreds.GenericCredentialModel.Password,
	// 		Database: dsc.DataStoreParams.DataSourceCreds.DB,
	// 		Port:     int(dsc.DataStoreParams.DataSourceCreds.Port),
	// 	}
	// 	if dsc.InApp {
	// 		client, err = vpostgres.NewPostgresStoreLocal(crds, log)
	// 	} else {
	// 		client, err = vpostgres.NewPostgresStore(crds, log)
	// 	}

	// 	if err != nil {
	// 		log.Err(err).Msg("Error connecting to postgres")
	// 		return nil, err
	// 	}
	// 	return &DataStoreClient{PostgresClient: client, DataStoreParams: dsc.DataStoreParams, Logger: log}, nil
	default:
		return nil, dmerrors.DMError(dputils.ErrInvalidDataStoreEngine, dputils.ErrDataStoreConn)
	}
}

func connectRedis(ctx context.Context, opts *models.DataSourceSecrets) (*vredis.RedisStore, error) {
	return vredis.NewRedisStore(ctx, &vredis.RedisOpts{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Password: opts.GenericCredentialModel.Password,
		Username: opts.GenericCredentialModel.Username,
	}, log)
}

func connectMysql(ctx context.Context, opts *models.DataSourceSecrets) (*vmysql.MysqlStore, error) {
	return vmysql.NewMysqlStore(&vmysql.MysqlOpts{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
		Database: opts.DB,
	}, log)
}

func connectElaticSearch(ctx context.Context, opts *models.DataSourceSecrets) (*ves.ElasticSearchStore, error) {
	return ves.NewElasticSearchStore(&ves.ElasticSearchOpts{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
		ApiKey:   opts.GenericCredentialModel.ApiToken,
	}, log)
}

func connectSnowflake(ctx context.Context, opts *models.DataSourceSecrets) (*vsnowflake.SnowflakeStore, error) {
	return vsnowflake.NewConnectSnowflake(&vsnowflake.SnowflakeOpts{
		URL: opts.URL,
		// Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
		Database: opts.DB,
	}, log)
}

func connectOpenSearch(ctx context.Context, opts *models.DataSourceSecrets) (*vopensearch.OpenSearchStore, error) {
	return vopensearch.NewOpenSearchStore(&vopensearch.OpenSearchOpts{
		Endpoint: opts.URL,
		// Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
	}, log)
}

func connectSingleStore(ctx context.Context, opts *models.DataSourceSecrets) (*vsinglestore.SingleStoreStore, error) {
	return vsinglestore.NewConnectSingleStore(&vsinglestore.SingleStoreOpts{
		Host:     opts.URL,
		Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
		Database: opts.DB,
	}, log)
}

func connectSqlServer(ctx context.Context, opts *models.DataSourceSecrets) (*vsqlserver.SqlServerStore, error) {
	return vsqlserver.NewSqlServerStore(&vsqlserver.SqlServerOpts{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
		Database: opts.DB,
	}, log)
}

func connectMariaDb(ctx context.Context, opts *models.DataSourceSecrets) (*vmariadb.MariaStore, error) {
	return vmariadb.NewMariaStore(&vmariadb.MariaOpts{
		URL:      opts.URL,
		Port:     int(opts.Port),
		Username: opts.GenericCredentialModel.Username,
		Password: opts.GenericCredentialModel.Password,
		Database: opts.DB,
	}, log)
}

func connectDynamoDB(ctx context.Context, opts *models.DataSourceSecrets) (*vdynamodb.DynamoDBStore, error) {
	return vdynamodb.NewDynamoDBStore(&vdynamodb.DynamoDBOpts{
		Region:    opts.GenericCredentialModel.AwsCreds.Region,
		AccessKey: opts.GenericCredentialModel.AwsCreds.AccessKeyId,
		SecretKey: opts.GenericCredentialModel.AwsCreds.SecretAccessKey,
	}, log)
}

func connectDatabricks(ctx context.Context, opts *models.DataSourceSecrets) (*vdatabricks.DatabricksStore, error) {

	return vdatabricks.NewConnectDatabricks(&vdatabricks.DatabricksOpts{
		Host:  opts.URL,
		Token: opts.GenericCredentialModel.ApiToken,
		Path:  opts.Params["Path"].(string),
		Port:  int(opts.Port),
	}, log)
}

func connectPinecone(ctx context.Context, opts *models.DataSourceSecrets) (*vpinecone.PineconeStore, error) {
	return vpinecone.NewConnectPinecone(&vpinecone.PineconeOpts{
		ApiKey: opts.GenericCredentialModel.ApiToken,
	}, log)
}

func connectQdrant(ctx context.Context, opts *models.DataSourceSecrets) (*vqdrant.QdrantStore, error) {

	return vqdrant.NewConnectQdrant(&vqdrant.QdrantOpts{
		Host:   opts.URL,
		Port:   int(opts.Port),
		ApiKey: opts.GenericCredentialModel.ApiToken,
	}, log)
}
