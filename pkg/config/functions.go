package config

import (
	"context"

	envconfig "github.com/sethvargo/go-envconfig"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/environment"
)

func Process(ctx context.Context, environment environment.Environment, config interface{}) error {
	return envconfig.ProcessWith(ctx, config, &EnvironmentLookup{
		environment: environment,
	})
}
