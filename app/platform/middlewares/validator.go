package middlewares

import (
	"context"

	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	grpc "google.golang.org/grpc"
)

func UnaryRequestValidator() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger = pkgs.GetSubDMLogger("Middleware", "Request Validator")
		// Validate the request
		methodName, _ := grpc.Method(ctx)
		logger.Info().Msgf("Validating Request for method - %v ", methodName)
		if err := pkgs.SvcPackageManager.GrpcRequestValidator.Validate(request); err != nil {
			// return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %v", err)
			return handler(ctx, request)
		}

		// Continue to handler
		return handler(ctx, request)
	}
}
