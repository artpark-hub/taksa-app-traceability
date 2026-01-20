package server

import (
	v1 "traceability/api/helloworld/v1"
	trv1 "traceability/api/traceability/v1"
	"traceability/internal/conf"
	"traceability/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, traceability *service.TraceabilityService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	
	// Register Services
	v1.RegisterGreeterHTTPServer(srv, greeter)
	trv1.RegisterTraceabilityHTTPServer(srv, traceability)

	return srv
}
