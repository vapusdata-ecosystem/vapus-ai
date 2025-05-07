package appcl

import (
	"context"
	"log"

	"github.com/rs/zerolog"
	aiutilitypb "github.com/vapusdata-ecosystem/apis/protos/vapus-aiutilities/v1alpha1"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (x *VapusSvcInternalClients) GenerateSummary(ctx context.Context, req *aiutilitypb.SummarizerRequest, logger zerolog.Logger, retryCount int) (*aiutilitypb.SummarizerResponse, error) {
	if x == nil {
		return nil, ErrAIStudioConnNotInitialized
	}
	if x.AIUtilityServerClient == nil {
		err := x.PingTestAndReconnect(ctx, x.AiUtilityServerDns, logger)
		if err != nil || x.AIUtilityServerClient == nil {
			return nil, ErrAIStudioConnNotInitialized
		}
	}
	resp, err := x.AIUtilityServerClient.Summarizer(pbtools.SwapNewContextWithAuthToken(ctx), req)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while Generating content.")
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if status.Code(err).String() == codes.Unavailable.String() {
				log.Println("Retry count", retryCount)
				if retryCount > 3 {
					return nil, ErrGeneratingContent
				}
				retryCount++
				logger.Err(err).Ctx(ctx).Msgf("error while calling AIUtilityServerClient for summarizer, retrying %d", retryCount)
				return x.GenerateSummary(ctx, req, logger, retryCount)
			}
		}
	}
	if err != nil {
		return nil, ErrGeneratingContent
	}
	log.Println("Summary Generator Called ====================>>>>>>>>>>>>>>>>>>>>>>>>3333", resp)
	return resp, nil
}
