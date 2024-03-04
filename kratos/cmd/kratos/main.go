package main

import (
	"context"
	"fmt"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/environment"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string

	id, _ = os.Hostname()
)

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(gs, hs),
	)
}

func main() {
	ctx, logger := context.Background(), log.Custom()
	environment := environment.NewDefaultEnvironment(environment.WithArraySource("env-vars", os.Environ()))
	var cfg config.Config
	if err := config.Process(ctx, environment, &cfg); err != nil {
		log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
	}

	app, err := wireApp(&cfg, logger)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
