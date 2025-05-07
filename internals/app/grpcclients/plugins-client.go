package appcl

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rs/zerolog"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (x *VapusSvcInternalClients) PluginAction(ctx context.Context, spec any, request *pb.PluginActionRequest, logger zerolog.Logger, retryCount int) error {
	if x == nil {
		return ErrUserConnNotInitialized
	}
	if x.PlConn == nil {
		err := x.PingTestAndReconnect(ctx, x.PlDns, logger)
		if err != nil {
			return ErrUserConnNotInitialized
		}
	}
	if request != nil && request.Spec == nil {
		reqbytes, err := json.Marshal(spec)
		if err != nil {
			logger.Err(err).Ctx(ctx).Msg("error while marshalling request spec for plugin action")
			return ErrInvalidPluginActionRequest
		}
		request.Spec = reqbytes
	} else if request == nil {
		return ErrInvalidPluginActionRequest
	}
	_, err := x.PluginServiceClient.Action(pbtools.SwapNewContextWithAuthToken(ctx), request)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while calling plugin action for plugin type %v", request.PluginType)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if status.Code(err).String() == codes.Unavailable.String() {
				log.Println("Retry count", retryCount)
				if retryCount > 3 {
					return ErrPluginActionFailed
				}
				retryCount++
				logger.Err(err).Ctx(ctx).Msgf("error while calling platform server for sending email, retrying with new connection. Count = %v", retryCount)
				return x.PluginAction(ctx, spec, request, logger, retryCount)
			}
		}
	}
	if err != nil {
		return ErrPluginActionFailed
	}
	return nil

}
