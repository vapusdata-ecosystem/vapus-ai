package bedrock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	bedrockService "github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/rs/zerolog"
	awsConfig "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/aws"
)

type BedrockOpts struct {
	SecretAccessKey, AccessKeyId, Region string
}

type BedrockGuardrailsStore struct {
	*awsConfig.BedrockConfig
}

type BedrockGuardrailOpts struct {
	Id      string
	Version string
	Text    string
}

func NewBedrockGuardrails(ctx context.Context, aws *BedrockOpts, logger zerolog.Logger) (*BedrockGuardrailsStore, error) {
	brClient, err := awsConfig.GetBedrockConfig(ctx, &awsConfig.AWSConfig{
		Region:          aws.Region,
		AccessKeyId:     aws.AccessKeyId,
		SecretAccessKey: aws.SecretAccessKey,
	})
	if err != nil {
		return nil, err
	}
	return &BedrockGuardrailsStore{
		BedrockConfig: brClient,
	}, nil
}

func (b BedrockGuardrailsStore) ListGuardrails(ctx context.Context, identifier *string, logger zerolog.Logger) ([]*BedrockGuardrailList, error) {
	fmt.Println("I am in List Guradrails")
	listGuargrails, err := b.BedrockService.ListGuardrails(ctx, &bedrockService.ListGuardrailsInput{})
	if err != nil {
		logger.Err(err).Msg("Unable to list the guardrail from bedrock")
		return nil, err
	}
	result := []*BedrockGuardrailList{}
	for _, val := range listGuargrails.Guardrails {
		temp := &BedrockGuardrailList{
			ARN:  *val.Arn,
			Id:   *val.Id,
			Name: *val.Name,
		}
		result = append(result, temp)
	}
	fmt.Println("List result: ", reflect.ValueOf(result))
	// b.BedrockRuntime.ApplyGuardrail()
	return result, nil
}

func (b BedrockGuardrailsStore) ApplyGuardrail(ctx context.Context, data *BedrockGuardrailOpts, logger zerolog.Logger) (*bedrockruntime.ApplyGuardrailOutput, error) {
	fmt.Println("I am in Apply Guardrails")
	if data == nil || data.Text == "" {
		return nil, errors.New("invalid or nil response")
	}
	applyInput := &bedrockruntime.ApplyGuardrailInput{
		GuardrailIdentifier: aws.String(data.Id),
		GuardrailVersion:    aws.String(data.Version),
		Source:              types.GuardrailContentSource(*aws.String("INPUT")),
		Content: []types.GuardrailContentBlock{
			&types.GuardrailContentBlockMemberText{
				Value: types.GuardrailTextBlock{
					Text: &data.Text,
				},
			},
		},
	}
	fmt.Println("Input: ", reflect.ValueOf(applyInput))
	// payload need to build
	resp, err := b.BedrockRuntime.ApplyGuardrail(ctx, applyInput)
	if err != nil {
		logger.Err(err).Msg("failed to guard messages from bedrock")
		return nil, err
	}
	converted := ConvertFromBedrockOutput(resp)
	jsonBytes, err := json.MarshalIndent(converted, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal struct:", err)
	} else {
		fmt.Println("Guardrail Response: ", string(jsonBytes))
	}

	fmt.Println(*resp.Outputs[0].Text)
	return resp, nil
}
