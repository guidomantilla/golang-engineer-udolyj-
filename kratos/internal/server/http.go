package server

import (
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/api"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/config"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"net"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cfg *config.Config, service api.ApiHTTPServer) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}

	if *cfg.Host != "" && *cfg.HttpPort != "" {
		opts = append(opts, http.Address(net.JoinHostPort(*cfg.Host, *cfg.HttpPort)))
	}

	srv := http.NewServer(opts...)
	api.RegisterApiHTTPServer(srv, service)
	return srv
}
