package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qmdx00/lifecycle"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/log"
)

var _ lifecycle.Server = (*HttpServer)(nil)

type HttpServer struct {
	internal *http.Server
}

func BuildHttpServer(server *http.Server) lifecycle.Server {
	return &HttpServer{
		internal: server,
	}
}

func (server *HttpServer) Run(ctx context.Context) error {

	log.Info(fmt.Sprintf("starting up - starting http server: %s", server.internal.Addr))

	if err := server.internal.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(fmt.Sprintf("starting up - starting http server error: %s", err.Error()))
		return err
	}

	return nil
}

func (server *HttpServer) Stop(ctx context.Context) error {

	log.Info("shutting down - stopping http server")

	if err := server.internal.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("shutting down - forced shutdown: %s", err.Error()))
		return err
	}

	log.Info("shutting down - http server stopped")
	return nil
}
