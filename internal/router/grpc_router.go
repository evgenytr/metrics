package router

import (
	pb "github.com/evgenytr/metrics.git/gen/go/metrics/v1"
	"github.com/evgenytr/metrics.git/internal/handlers"
	"google.golang.org/grpc"
)

func GrpcRouter(h *handlers.StorageHandler) (g *grpc.Server, err error) {
	g = grpc.NewServer()

	pb.RegisterMetricsServiceV1Server(g, h)
	return
}
