//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/providers"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(cfg *config.Config, logger log.Logger) (*kratos.App, error) {
	panic(wire.Build(providers.ProviderSet, newApp))
}
