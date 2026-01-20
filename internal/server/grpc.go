package server

import (
	v1 "traceability/api/helloworld/v1"
	trv1 "traceability/api/traceability/v1"
	"traceability/internal/conf"
	"traceability/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	greeter *service.GreeterService,
	traceability *service.TraceabilityService,
	logger log.Logger,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	// Register Greeter
	v1.RegisterGreeterServer(srv, greeter)

	// Register Traceability
	trv1.RegisterTraceabilityServer(srv, traceability)

	return srv
}
