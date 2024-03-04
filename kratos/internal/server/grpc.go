package server

import (
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/api"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/config"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"net"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cfg *config.Config, service api.ApiServer) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}

	if *cfg.Host != "" && *cfg.GrpcPort != "" {
		opts = append(opts, grpc.Address(net.JoinHostPort(*cfg.Host, *cfg.GrpcPort)))
	}

	srv := grpc.NewServer(opts...)
	api.RegisterApiServer(srv, service)
	return srv
}
