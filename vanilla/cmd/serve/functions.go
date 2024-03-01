package serve

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/endpoints/rpc"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/boot"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/log"
)

func ExecuteCmdFn(_ *cobra.Command, args []string) {

	ctx := context.Background()
	logger := log.Custom()
	appName, version := config.Application, config.Version
	enablers := &boot.Enablers{
		HttpServerEnabled: true,
		GrpcServerEnabled: true,
		DatabaseEnabled:   true,
	}

	builder := boot.NewBeanBuilder(ctx)
	builder.Config = func(appCtx *boot.ApplicationContext) {
		var cfg config.Config
		if err := config.Process(ctx, appCtx.Environment, &cfg); err != nil {
			log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
		}

		appCtx.HttpConfig = &boot.HttpConfig{
			Host: cfg.Host,
			Port: cfg.HttpPort,
		}

		appCtx.GrpcConfig = &boot.GrpcConfig{
			Host: cfg.Host,
			Port: cfg.GrpcPort,
		}

		appCtx.SecurityConfig = &boot.SecurityConfig{
			TokenSignatureKey:    cfg.TokenSignatureKey,
			TokenVerificationKey: cfg.TokenVerificationKey,
		}

		appCtx.DatabaseConfig = &boot.DatabaseConfig{
			DatasourceUrl:      cfg.DatasourceUrl,
			DatasourceServer:   cfg.DatasourceServer,
			DatasourceService:  cfg.DatasourceService,
			DatasourceUsername: cfg.DatasourceUsername,
			DatasourcePassword: cfg.DatasourcePassword,
		}
	}

	builder.GrpcServer = func(appCtx *boot.ApplicationContext) (*grpc.ServiceDesc, any) {
		grpcServer := rpc.NewBankApiGrpcServer(appCtx.AuthenticationService, appCtx.AuthorizationService, appCtx.PrincipalManager)
		return &rpc.BankApi_ServiceDesc, grpcServer
	}

	err := boot.Init(appName, version, args, logger, enablers, builder, func(appCtx boot.ApplicationContext) error {
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}

}
