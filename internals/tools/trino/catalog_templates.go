package trinocl

import "github.com/vapusdata-ecosystem/vapusai/core/types"

var TrinoCatalogContentMap = map[string]string{
	types.StorageEngine_POSTGRES.String(): `connector.name=postgresql
	connection-url=jdbc:postgresql://%s:%d/%s
	credential-provider.type=%s
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_MYSQL.String(): `connector.name=mysql
	connection-url=jdbc:mysql://%s:%d
	credential-provider.type=%s
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_ORACLE.String(): `connector.name=oracle
	connection-url=jdbc:oracle://%s:%s/%s
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_ELASTICSEARCH.String(): `connector.name=elasticsearch
	elasticsearch.host=%s
	elasticsearch.port=%s
	elasticsearch.default-schema-name=%s
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_MONGODB.String(): `connector.name=mongodb
	mongodb.connection-url=%s
	mongodb.schema-collection=%s
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_OPENSEARCH.String(): `connector.name=elasticsearch
	opensearch.host=%s
	opensearch.port=%s
	opensearch.default-schema-name=%s
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}
	opensearch.security=%s`,
	types.StorageEngine_SQL_SERVER.String(): `connector.name=sqlserver
	connection-url=jdbc:sqlserver://%s:%s/%s
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_REDSHIFT.String(): `connector.name=redshift
	connection-url=jdbc:redshift://%s:%d/%s
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_BIGQUERY.String(): `connector.name=bigquery
	bigquery.project-id=%s
	bigquery.credentials-file={secretsFileFullPath}`,
	types.StorageEngine_CLICKHOUSE.String(): `connector.name=clickhouse
	connection-url=jdbc:clickhouse://%s:%d/
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_MARIADB.String(): `connector.name=mariadb
	connection-url=jdbc:mariadb://%s:%d
	credential-provider.type=%s
	connection-credential-file={secretsFileFullPath}`,
	types.StorageEngine_SINGLE_STORE.String(): `connector.name=singlestore
	connection-url=jdbc:singlestore://%s:%d
	credential-provider.type=FILE
	connection-credential-file={secretsFileFullPath}`,
}

var TrinoCatalogSecretsMap = map[string]string{
	types.StorageEngine_POSTGRES.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_MYSQL.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_ELASTICSEARCH.String(): `elasticsearch.auth.user=%s
		elasticsearch.auth.password=%s`,
	types.StorageEngine_OPENSEARCH.String(): `opensearch.auth.user=%s
		opensearch.auth.password=%s`,
	types.StorageEngine_ORACLE.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_MYSQL.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_REDSHIFT.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_BIGQUERY.String(): `%s`,
	types.StorageEngine_CLICKHOUSE.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_MARIADB.String(): `connection-user=%s
		connection-password=%s`,
	types.StorageEngine_SINGLE_STORE.String(): `aws.access-key=%s
		aws.secret-key=%s`,
}
