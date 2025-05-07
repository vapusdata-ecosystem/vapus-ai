package vaws

import (
	"context"
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
)

type RedshiftStore struct {
	host     string `json:"host,omitempty" yaml:"host"`
	port     int32  `json:"port,omitempty" yaml:"port"`
	database string `json:"database,omitempty" yaml:"database"`
	username string `json:"username,omitempty" yaml:"username"`
	opts     *AWSConfig
	client   *redshiftdata.Client
}

func NewRedshiftStore(ctx context.Context, host string, port int32, database, username string, opts *AWSConfig) (*RedshiftStore, error) {
	configCl, err := opts.getAwsCLientConfig(ctx)
	if err != nil {
		return nil, dmerrors.DMError(ErrAwsConfigLoading, nil)
	}
	return &RedshiftStore{
		opts:     opts,
		host:     host,
		port:     port,
		database: database,
		username: username,
		client:   redshiftdata.NewFromConfig(configCl),
	}, nil
}

func (rs *RedshiftStore) ExecuteQuery(ctx context.Context, query string, logger zerolog.Logger) ([]map[string]any, error) {
	input := &redshiftdata.ExecuteStatementInput{
		ClusterIdentifier: aws.String("your-cluster-identifier"),
		Database:          aws.String("your-database"),
		DbUser:            aws.String("your-db-user"),
		Sql:               aws.String("SELECT current_timestamp"),
	}

	executeResult, err := rs.client.ExecuteStatement(ctx, input)
	if err != nil {
		logger.Error().Err(err).Msg("Error executing query on redshift")
		return nil, err
	}
	var status types.StatusString
	for {
		describeInput := &redshiftdata.DescribeStatementInput{
			Id: executeResult.Id,
		}
		describeResult, err := rs.client.DescribeStatement(ctx, describeInput)
		if err != nil {
			logger.Error().Err(err).Msg("Error describing query on redshift")
			return nil, err
		}
		status = describeResult.Status
		if status == types.StatusStringFinished || status == types.StatusStringFailed || status == types.StatusStringAborted {
			break
		}
		time.Sleep(1 * time.Second) // Wait before polling again
	}
	if status != types.StatusStringFinished {
		logger.Error().Msg("Query execution failed")
		return nil, err
	}
	getResultInput := &redshiftdata.GetStatementResultInput{
		Id: executeResult.Id,
	}
	getResultOutput, err := rs.client.GetStatementResult(context.TODO(), getResultInput)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting result of query on redshift")
		return nil, err
	}
	// Map the results to a slice of maps
	results := make([]map[string]any, len(getResultOutput.Records))
	for i, record := range getResultOutput.Records {
		row := make(map[string]any)
		for j, field := range record {
			columnName := *getResultOutput.ColumnMetadata[j].Name
			switch v := field.(type) {
			case *types.FieldMemberStringValue:
				row[columnName] = v.Value
			case *types.FieldMemberLongValue:
				row[columnName] = v.Value
			case *types.FieldMemberDoubleValue:
				row[columnName] = v.Value
			case *types.FieldMemberBooleanValue:
				row[columnName] = v.Value
			case *types.FieldMemberIsNull:
				row[columnName] = nil
			default:
				row[columnName] = fmt.Sprintf("%v", field)
			}
		}
		results[i] = row
	}
	return results, nil
}
