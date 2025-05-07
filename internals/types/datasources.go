package types

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type StorageEngine string

func (x StorageEngine) String() string {
	return string(x)
}

const (
	StorageEngine_INVALID_DS_SERVICE_TYPE StorageEngine = "INVALID_DS_SERVICE_TYPE"
	StorageEngine_MYSQL                   StorageEngine = "MYSQL"
	StorageEngine_REDIS                   StorageEngine = "REDIS"
	StorageEngine_ELASTICSEARCH           StorageEngine = "ELASTICSEARCH"
	StorageEngine_MONGODB                 StorageEngine = "MONGODB"
	StorageEngine_SECRET_MANAGER          StorageEngine = "SECRET_MANAGER"
	StorageEngine_POSTGRES                StorageEngine = "POSTGRES"
	StorageEngine_OCI                     StorageEngine = "OCI"
	StorageEngine_PYPI                    StorageEngine = "PYPI"
	StorageEngine_BLOB                    StorageEngine = "BLOB"
	StorageEngine_KAFKA_QUEUE             StorageEngine = "KAFKA_QUEUE"
	StorageEngine_RABBITMQ                StorageEngine = "RABBITMQ"
	StorageEngine_RESTAPI                 StorageEngine = "RESTAPI"
	StorageEngine_GRPC                    StorageEngine = "GRPC"
	StorageEngine_BIGQUERY                StorageEngine = "BIGQUERY"
	StorageEngine_SNOWFLAKE               StorageEngine = "SNOWFLAKE"
	StorageEngine_SQL_SERVER              StorageEngine = "SQL_SERVER"
	StorageEngine_OPENSEARCH              StorageEngine = "OPENSEARCH"
	StorageEngine_HIVE                    StorageEngine = "HIVE"
	StorageEngine_ICEBERG                 StorageEngine = "ICEBERG"
	StorageEngine_CLICKHOUSE              StorageEngine = "CLICKHOUSE"
	StorageEngine_MARIADB                 StorageEngine = "MARIADB"
	StorageEngine_DYNAMODB                StorageEngine = "DYNAMODB"
	StorageEngine_ORACLE                  StorageEngine = "ORACLE"
	StorageEngine_SINGLE_STORE            StorageEngine = "SINGLE_STORE"
	StorageEngine_CASSANDRA               StorageEngine = "CASSANDRA"
	StorageEngine_REDSHIFT                StorageEngine = "REDSHIFT"
	StorageEngine_GIT                     StorageEngine = "GIT"
	StorageEngine_ALLOYDB                 StorageEngine = "ALLOYDB"
	StorageEngine_DATABRICKS              StorageEngine = "DATABRICKS"
	StorageEngine_PINECONE                StorageEngine = "PINECONE"
	StorageEngine_QDRANT                  StorageEngine = "QDRANT"
	StorageEngine_KAFKA                   StorageEngine = "KAFKA"
	StorageEngine_DB2                     StorageEngine = "DB2"
	StorageEngine_SAAS                    StorageEngine = "SAAS"
	StorageEngine_FILE_STORE              StorageEngine = "FILE_STORE"
	StorageEngine_ARTIFACT                StorageEngine = "ARTIFACT"
)

var DataSourceTypeMap = map[mpb.DataSourceServices]mpb.DataSourceType{
	mpb.DataSourceServices_RDS_DB2:             mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_IBMDB2:              mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_RDS_MYSQL:           mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_GCP_SQL_MYSQL:       mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_MYSQL:               mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_RDS_POSTGRES:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_GCP_POSTGRES:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_POSTGRES:            mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_RDS_MARIADB:         mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_RDS_SQL_SERVER:      mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_SQLSERVER:           mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_RDS_ORACLEDB:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_AWS_ORACLEDB:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_ORACLE_DB:           mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_AWS_REDSHIFT:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_REDIS_ENTERPRISE:    mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_AWS_ELASTICCACHE:    mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_AWS_OPENSEARCH:      mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_ELASTICSEARCH_CLOUD: mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_MONGODB_ATLAS:       mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_AWS_DOCUMENTDB:      mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_MONGODB:             mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_GCP_BIGQUERY:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_SNOWFLAKE:           mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_HIVE:                mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_ICEBERG:             mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_CLICKHOUSE:          mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_AWS_DYNAMODB:        mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_SINGLESTORE:         mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_CASSANDRA:           mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_ALLOYDB:             mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_DATABRICKS:          mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_PINECONE:            mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_QDRANT:              mpb.DataSourceType_DATABASE,
	mpb.DataSourceServices_GITHUB_SVC:          mpb.DataSourceType_CODE_REPOSITORY,
	mpb.DataSourceServices_GITLAB_SVC:          mpb.DataSourceType_CODE_REPOSITORY,

	// Data stream types
	mpb.DataSourceServices_KAFKA:    mpb.DataSourceType_DATA_STREAM,
	mpb.DataSourceServices_RABBITMQ: mpb.DataSourceType_DATA_STREAM,

	// API stores
	mpb.DataSourceServices_REST_API: mpb.DataSourceType_API_STORE,
	mpb.DataSourceServices_GRPC:     mpb.DataSourceType_API_STORE,

	// Blob stores
	mpb.DataSourceServices_GCP_CLOUD_STORAGE: mpb.DataSourceType_BLOB_STORE,
	mpb.DataSourceServices_AWS_S3:            mpb.DataSourceType_BLOB_STORE,
	mpb.DataSourceServices_GOOGLE_DRIVE:      mpb.DataSourceType_BLOB_STORE,
	mpb.DataSourceServices_DROPBOX:           mpb.DataSourceType_BLOB_STORE,
	mpb.DataSourceServices_BOX:               mpb.DataSourceType_BLOB_STORE,

	// Secret stores
	mpb.DataSourceServices_AWS_SECRET_MANAGER:   mpb.DataSourceType_SECRET_STORE,
	mpb.DataSourceServices_GCP_SECRET_MANAGER:   mpb.DataSourceType_SECRET_STORE,
	mpb.DataSourceServices_AZURE_SECRET_MANAGER: mpb.DataSourceType_SECRET_STORE,
	mpb.DataSourceServices_HASHICORP_VAULT:      mpb.DataSourceType_SECRET_STORE,

	mpb.DataSourceServices_GAR: mpb.DataSourceType_ARTIFACT,
	mpb.DataSourceServices_ECR: mpb.DataSourceType_ARTIFACT,
}

var StorageEngineMap = map[mpb.DataSourceServices]StorageEngine{
	mpb.DataSourceServices_RDS_DB2:              StorageEngine_DB2,
	mpb.DataSourceServices_IBMDB2:               StorageEngine_DB2,
	mpb.DataSourceServices_RDS_MYSQL:            StorageEngine_MYSQL,
	mpb.DataSourceServices_GCP_SQL_MYSQL:        StorageEngine_MYSQL,
	mpb.DataSourceServices_MYSQL:                StorageEngine_MYSQL,
	mpb.DataSourceServices_RDS_POSTGRES:         StorageEngine_POSTGRES,
	mpb.DataSourceServices_GCP_POSTGRES:         StorageEngine_POSTGRES,
	mpb.DataSourceServices_POSTGRES:             StorageEngine_POSTGRES,
	mpb.DataSourceServices_RDS_MARIADB:          StorageEngine_MARIADB,
	mpb.DataSourceServices_RDS_SQL_SERVER:       StorageEngine_SQL_SERVER,
	mpb.DataSourceServices_SQLSERVER:            StorageEngine_SQL_SERVER,
	mpb.DataSourceServices_RDS_ORACLEDB:         StorageEngine_ORACLE,
	mpb.DataSourceServices_AWS_ORACLEDB:         StorageEngine_ORACLE,
	mpb.DataSourceServices_ORACLE_DB:            StorageEngine_ORACLE,
	mpb.DataSourceServices_AWS_REDSHIFT:         StorageEngine_REDSHIFT,
	mpb.DataSourceServices_REDIS_ENTERPRISE:     StorageEngine_REDIS,
	mpb.DataSourceServices_AWS_ELASTICCACHE:     StorageEngine_REDIS,
	mpb.DataSourceServices_AWS_OPENSEARCH:       StorageEngine_OPENSEARCH,
	mpb.DataSourceServices_ELASTICSEARCH_CLOUD:  StorageEngine_ELASTICSEARCH,
	mpb.DataSourceServices_MONGODB_ATLAS:        StorageEngine_MONGODB,
	mpb.DataSourceServices_AWS_DOCUMENTDB:       StorageEngine_MONGODB,
	mpb.DataSourceServices_MONGODB:              StorageEngine_MONGODB,
	mpb.DataSourceServices_GCP_BIGQUERY:         StorageEngine_BIGQUERY,
	mpb.DataSourceServices_SNOWFLAKE:            StorageEngine_SNOWFLAKE,
	mpb.DataSourceServices_HIVE:                 StorageEngine_HIVE,
	mpb.DataSourceServices_ICEBERG:              StorageEngine_ICEBERG,
	mpb.DataSourceServices_CLICKHOUSE:           StorageEngine_CLICKHOUSE,
	mpb.DataSourceServices_AWS_DYNAMODB:         StorageEngine_DYNAMODB,
	mpb.DataSourceServices_SINGLESTORE:          StorageEngine_SINGLE_STORE,
	mpb.DataSourceServices_CASSANDRA:            StorageEngine_CASSANDRA,
	mpb.DataSourceServices_GITHUB_SVC:           StorageEngine_GIT,
	mpb.DataSourceServices_GITLAB_SVC:           StorageEngine_GIT,
	mpb.DataSourceServices_ALLOYDB:              StorageEngine_ALLOYDB,
	mpb.DataSourceServices_DATABRICKS:           StorageEngine_DATABRICKS,
	mpb.DataSourceServices_PINECONE:             StorageEngine_PINECONE,
	mpb.DataSourceServices_QDRANT:               StorageEngine_QDRANT,
	mpb.DataSourceServices_KAFKA:                StorageEngine_KAFKA,
	mpb.DataSourceServices_REST_API:             StorageEngine_RESTAPI,
	mpb.DataSourceServices_GRPC:                 StorageEngine_GRPC,
	mpb.DataSourceServices_RABBITMQ:             StorageEngine_RABBITMQ,
	mpb.DataSourceServices_AWS_SECRET_MANAGER:   StorageEngine_SECRET_MANAGER,
	mpb.DataSourceServices_GCP_SECRET_MANAGER:   StorageEngine_SECRET_MANAGER,
	mpb.DataSourceServices_AZURE_SECRET_MANAGER: StorageEngine_SECRET_MANAGER,
	mpb.DataSourceServices_HASHICORP_VAULT:      StorageEngine_SECRET_MANAGER,
	mpb.DataSourceServices_GCP_CLOUD_STORAGE:    StorageEngine_BLOB,
	mpb.DataSourceServices_AWS_S3:               StorageEngine_BLOB,
	mpb.DataSourceServices_REDIS_STORE:          StorageEngine_REDIS,
	mpb.DataSourceServices_AWS_AURORA_MYSQL:     StorageEngine_MYSQL,
	mpb.DataSourceServices_AWS_AURORA_POSTGRES:  StorageEngine_POSTGRES,
	mpb.DataSourceServices_GOOGLE_DRIVE:         StorageEngine_FILE_STORE,
	mpb.DataSourceServices_DROPBOX:              StorageEngine_FILE_STORE,
	mpb.DataSourceServices_BOX:                  StorageEngine_FILE_STORE,
	mpb.DataSourceServices_ECR:                  StorageEngine_ARTIFACT,
	mpb.DataSourceServices_GAR:                  StorageEngine_ARTIFACT,
}

var StorageEngineLogoMap = map[StorageEngine]string{
	StorageEngine_DB2:            "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c6/IBM_DB2_logo.svg/1200px-IBM_DB2_logo.svg.png",
	StorageEngine_MYSQL:          "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
	StorageEngine_POSTGRES:       "https://wiki.postgresql.org/images/a/a4/PostgreSQL_logo.3colors.svg",
	StorageEngine_MARIADB:        "https://mariadb.org/wp-content/themes/twentynineteen-child/icons/mariadb_org_rgb_h.svg",
	StorageEngine_SQL_SERVER:     "https://www.microsoft.com/en-us/sql-server/img/sql-server-logo.png",
	StorageEngine_ORACLE:         "https://www.oracle.com/a/ocom/img/cb71-database-logo.png",
	StorageEngine_REDSHIFT:       "https://d2908q01vomqb2.cloudfront.net/887309d048beef83ad3eabf2a79a64a389ab1c9f/2023/04/27/redshift-logo.png",
	StorageEngine_REDIS:          "https://redis.io/images/redis-logo.svg",
	StorageEngine_OPENSEARCH:     "https://opensearch.org/assets/brand/SVG/Logo/opensearch_logo_default.svg",
	StorageEngine_ELASTICSEARCH:  "https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt8781708f8f37ed16/5c11ec2edf09df047814db23/logo-elastic-elasticsearch-lt.svg",
	StorageEngine_MONGODB:        "https://www.mongodb.com/assets/images/global/leaf.svg",
	StorageEngine_BIGQUERY:       "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/products/logos/bigquery.svg",
	StorageEngine_SNOWFLAKE:      "https://www.snowflake.com/wp-content/themes/snowflake/assets/img/logo-white.svg",
	StorageEngine_HIVE:           "https://hive.apache.org/images/hive_logo_medium.jpg",
	StorageEngine_ICEBERG:        "https://iceberg.apache.org/assets/images/iceberg-logo.svg",
	StorageEngine_CLICKHOUSE:     "https://avatars.githubusercontent.com/u/54801242",
	StorageEngine_DYNAMODB:       "https://d2908q01vomqb2.cloudfront.net/887309d048beef83ad3eabf2a79a64a389ab1c9f/2020/12/11/DynamoDB-Logo.png",
	StorageEngine_SINGLE_STORE:   "https://www.singlestore.com/images/logo-light.svg",
	StorageEngine_CASSANDRA:      "https://cassandra.apache.org/_static/images/logo.png",
	StorageEngine_GIT:            "https://git-scm.com/images/logos/downloads/Git-Logo-2Color.svg",
	StorageEngine_ALLOYDB:        "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/products/logos/alloydb.svg",
	StorageEngine_DATABRICKS:     "https://www.databricks.com/wp-content/uploads/2022/09/db-nav-logo.svg",
	StorageEngine_PINECONE:       "https://www.pinecone.io/images/pinecone.svg",
	StorageEngine_QDRANT:         "https://github.com/qdrant/qdrant/raw/master/docs/logo.svg",
	StorageEngine_KAFKA:          "https://kafka.apache.org/logos/kafka_logo--simple.png",
	StorageEngine_RESTAPI:        "https://restfulapi.net/wp-content/uploads/rest.png",
	StorageEngine_GRPC:           "https://grpc.io/img/logos/grpc-logo.png",
	StorageEngine_RABBITMQ:       "https://www.rabbitmq.com/img/logo-rabbitmq.svg",
	StorageEngine_SECRET_MANAGER: "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/products/logos/secret-manager.svg",
}

// Map of DataSourceServices to their logo URLs
var DataSourceServicesLogoMap = map[mpb.DataSourceServices]string{
	mpb.DataSourceServices_RDS_DB2:              "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c6/IBM_DB2_logo.svg/1200px-IBM_DB2_logo.svg.png",
	mpb.DataSourceServices_IBMDB2:               "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c6/IBM_DB2_logo.svg/1200px-IBM_DB2_logo.svg.png",
	mpb.DataSourceServices_RDS_MYSQL:            "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
	mpb.DataSourceServices_GCP_SQL_MYSQL:        "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
	mpb.DataSourceServices_MYSQL:                "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
	mpb.DataSourceServices_RDS_POSTGRES:         "https://wiki.postgresql.org/images/a/a4/PostgreSQL_logo.3colors.svg",
	mpb.DataSourceServices_GCP_POSTGRES:         "https://wiki.postgresql.org/images/a/a4/PostgreSQL_logo.3colors.svg",
	mpb.DataSourceServices_POSTGRES:             "https://wiki.postgresql.org/images/a/a4/PostgreSQL_logo.3colors.svg",
	mpb.DataSourceServices_RDS_MARIADB:          "https://mariadb.org/wp-content/themes/twentynineteen-child/icons/mariadb_org_rgb_h.svg",
	mpb.DataSourceServices_RDS_SQL_SERVER:       "https://www.microsoft.com/en-us/sql-server/img/sql-server-logo.png",
	mpb.DataSourceServices_SQLSERVER:            "https://www.microsoft.com/en-us/sql-server/img/sql-server-logo.png",
	mpb.DataSourceServices_RDS_ORACLEDB:         "https://www.oracle.com/a/ocom/img/cb71-database-logo.png",
	mpb.DataSourceServices_AWS_ORACLEDB:         "https://www.oracle.com/a/ocom/img/cb71-database-logo.png",
	mpb.DataSourceServices_ORACLE_DB:            "https://www.oracle.com/a/ocom/img/cb71-database-logo.png",
	mpb.DataSourceServices_AWS_REDSHIFT:         "https://d2908q01vomqb2.cloudfront.net/887309d048beef83ad3eabf2a79a64a389ab1c9f/2023/04/27/redshift-logo.png",
	mpb.DataSourceServices_REDIS_ENTERPRISE:     "https://redis.io/images/redis-logo.svg",
	mpb.DataSourceServices_AWS_ELASTICCACHE:     "https://redis.io/images/redis-logo.svg",
	mpb.DataSourceServices_AWS_OPENSEARCH:       "https://opensearch.org/assets/brand/SVG/Logo/opensearch_logo_default.svg",
	mpb.DataSourceServices_ELASTICSEARCH_CLOUD:  "https://static-www.elastic.co/v3/assets/bltefdd0b53724fa2ce/blt8781708f8f37ed16/5c11ec2edf09df047814db23/logo-elastic-elasticsearch-lt.svg",
	mpb.DataSourceServices_MONGODB_ATLAS:        "https://www.mongodb.com/assets/images/global/leaf.svg",
	mpb.DataSourceServices_AWS_DOCUMENTDB:       "https://www.mongodb.com/assets/images/global/leaf.svg",
	mpb.DataSourceServices_MONGODB:              "https://www.mongodb.com/assets/images/global/leaf.svg",
	mpb.DataSourceServices_GCP_BIGQUERY:         "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/products/logos/bigquery.svg",
	mpb.DataSourceServices_SNOWFLAKE:            "https://www.snowflake.com/wp-content/themes/snowflake/assets/img/logo-white.svg",
	mpb.DataSourceServices_HIVE:                 "https://hive.apache.org/images/hive_logo_medium.jpg",
	mpb.DataSourceServices_ICEBERG:              "https://iceberg.apache.org/assets/images/iceberg-logo.svg",
	mpb.DataSourceServices_CLICKHOUSE:           "https://avatars.githubusercontent.com/u/54801242",
	mpb.DataSourceServices_AWS_DYNAMODB:         "https://d2908q01vomqb2.cloudfront.net/887309d048beef83ad3eabf2a79a64a389ab1c9f/2020/12/11/DynamoDB-Logo.png",
	mpb.DataSourceServices_SINGLESTORE:          "https://www.singlestore.com/images/logo-light.svg",
	mpb.DataSourceServices_CASSANDRA:            "https://cassandra.apache.org/_static/images/logo.png",
	mpb.DataSourceServices_GITHUB_SVC:           "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png",
	mpb.DataSourceServices_GITLAB_SVC:           "https://about.gitlab.com/images/press/logo/svg/gitlab-logo-500.svg",
	mpb.DataSourceServices_ALLOYDB:              "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/products/logos/alloydb.svg",
	mpb.DataSourceServices_DATABRICKS:           "https://www.databricks.com/wp-content/uploads/2022/09/db-nav-logo.svg",
	mpb.DataSourceServices_PINECONE:             "https://www.pinecone.io/images/pinecone.svg",
	mpb.DataSourceServices_QDRANT:               "https://github.com/qdrant/qdrant/raw/master/docs/logo.svg",
	mpb.DataSourceServices_KAFKA:                "https://kafka.apache.org/logos/kafka_logo--simple.png",
	mpb.DataSourceServices_REST_API:             "https://restfulapi.net/wp-content/uploads/rest.png",
	mpb.DataSourceServices_GRPC:                 "https://grpc.io/img/logos/grpc-logo.png",
	mpb.DataSourceServices_RABBITMQ:             "https://www.rabbitmq.com/img/logo-rabbitmq.svg",
	mpb.DataSourceServices_AWS_SECRET_MANAGER:   "https://d1.awsstatic.com/products/SecretsManager/product-page-diagram_AWS-Secrets-Manager%402x.d52255980433ad26dbb0da99595dbc550843a3df.png",
	mpb.DataSourceServices_GCP_SECRET_MANAGER:   "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/products/logos/secret-manager.svg",
	mpb.DataSourceServices_AZURE_SECRET_MANAGER: "https://learn.microsoft.com/azure/media/index/security-center.svg",
	mpb.DataSourceServices_HASHICORP_VAULT:      "https://www.hashicorp.com/img/logo-hashicorp-vault.svg",
}

var ServiceProviderLogoMap = map[mpb.ServiceProvider]string{
	mpb.ServiceProvider_INVALID_PROVIDER: "",
	mpb.ServiceProvider_OPENAI:           "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/OpenAI_Logo.svg/1280px-OpenAI_Logo.svg.png",
	mpb.ServiceProvider_MISTRAL:          "https://avatars.githubusercontent.com/u/99472018",
	mpb.ServiceProvider_HUGGING_FACE:     "https://huggingface.co/front/assets/huggingface_logo.svg",
	mpb.ServiceProvider_VAPUS:            "https://vapus.ai/favicon.ico",
	mpb.ServiceProvider_OLLAMA:           "https://ollama.com/public/ollama.png",
	mpb.ServiceProvider_AZURE_OPENAI:     "https://learn.microsoft.com/azure/media/index/cognitive-services.svg",
	mpb.ServiceProvider_AZURE_PHI:        "https://learn.microsoft.com/azure/media/index/cognitive-services.svg",
	mpb.ServiceProvider_GEMINI:           "https://lh3.googleusercontent.com/UAdRVHEHXTutqH4aDXI3dRETbsGJoYWJ5NULCNXRCNKj8OGYnDMYYoC3Vj4aD26JJhhTs4VFmiCnVJa53ISBGPMiLKgnVOTUvIjpbK0",
	mpb.ServiceProvider_AWS:              "https://a0.awsstatic.com/libra-css/images/logos/aws_logo_smile_1200x630.png",
	mpb.ServiceProvider_META:             "https://about.meta.com/brand/meta/brand-and-resources/images/M-Favicon.jpg",
	mpb.ServiceProvider_ANTHROPIC:        "https://storage.googleapis.com/anthropic-public-website/anthropic-logo-1024.png",
	mpb.ServiceProvider_GENERIC:          "",
	mpb.ServiceProvider_GROQ:             "https://miro.medium.com/v2/resize:fit:1400/1*z9S4GH-fqxGJ9Sm0V6Isfw.png",
	mpb.ServiceProvider_BEDROCK:          "https://a0.awsstatic.com/libra-css/images/logos/aws_logo_smile_1200x630.png",
	mpb.ServiceProvider_TOGETHER:         "https://assets-global.website-files.com/64f6f2c0e3ed0c667c5713c4/650462312b88bbcf356ca397_Together_Logo-01.svg",
	mpb.ServiceProvider_GITHUB:           "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png",
	mpb.ServiceProvider_MONGO_ORG:        "https://www.mongodb.com/assets/images/global/leaf.svg",
	mpb.ServiceProvider_REDIS_ORG:        "https://redis.io/images/redis-logo.svg",
	mpb.ServiceProvider_GITLAB:           "https://about.gitlab.com/images/press/logo/svg/gitlab-logo-500.svg",
	mpb.ServiceProvider_BITBUCKET:        "https://wac-cdn.atlassian.com/assets/img/favicons/bitbucket/favicon-32x32.png",
	mpb.ServiceProvider_MICROSOFT:        "https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31",
	mpb.ServiceProvider_REDHAT:           "https://www.redhat.com/themes/custom/rhdc/img/red-hat-logo.svg",
	mpb.ServiceProvider_GCP:              "https://www.gstatic.com/devrel-devsite/prod/v2f6fb68338062e7c16672db62c4ab042dcb9bfbacf2fa51b6959426b203a4d8a/cloud/images/cloud-logo.svg",
	mpb.ServiceProvider_SELF_HOSTED:      "",
}
