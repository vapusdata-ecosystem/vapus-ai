package dmcontrollers

import (
	"context"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthCheck struct{}

func NewHealthCheckController() pb.HealthServer {
	return new(HealthCheck)
}

func (h *HealthCheck) Check(c context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthCheck) Watch(hc *pb.HealthCheckRequest, hs pb.Health_WatchServer) error {
	return nil
}

func (h *HealthCheck) List(context.Context, *pb.HealthListRequest) (*pb.HealthListResponse, error) {
	return &pb.HealthListResponse{
		Statuses: map[string]*pb.HealthCheckResponse{
			"vapusdata": {
				Status: pb.HealthCheckResponse_SERVING,
			},
		},
	}, nil
}
