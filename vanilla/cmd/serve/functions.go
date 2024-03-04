package serve

import (
	"context"
	"fmt"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/boot"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/endpoints/rest"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/endpoints/rpc"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/services"
)

func ExecuteCmdFn(_ *cobra.Command, args []string) {

	ctx, logger := context.Background(), log.Custom()
	appName, version := config.Application, config.Version
	enablers := &boot.Enablers{HttpServerEnabled: true, GrpcServerEnabled: true, DatabaseEnabled: true}

	builder := boot.NewBeanBuilder(ctx)
	builder.Config = func(appCtx *boot.ApplicationContext) {
		var cfg config.Config
		if err := config.Process(ctx, appCtx.Environment, &cfg); err != nil {
			log.Fatal(fmt.Sprintf("starting up - error setting up configuration: %s", err.Error()))
		}

		appCtx.HttpConfig = &boot.HttpConfig{Host: cfg.Host, Port: cfg.HttpPort}
		appCtx.GrpcConfig = &boot.GrpcConfig{Host: cfg.Host, Port: cfg.GrpcPort}
		appCtx.SecurityConfig = &boot.SecurityConfig{TokenSignatureKey: cfg.TokenSignatureKey, TokenVerificationKey: cfg.TokenVerificationKey}
		appCtx.DatabaseConfig = &boot.DatabaseConfig{
			DatasourceUrl:      cfg.DatasourceUrl,
			DatasourceServer:   cfg.DatasourceServer,
			DatasourceService:  cfg.DatasourceService,
			DatasourceUsername: cfg.DatasourceUsername,
			DatasourcePassword: cfg.DatasourcePassword,
		}
	}

	var bankService services.BankService
	builder.GrpcServer = func(appCtx *boot.ApplicationContext) (*grpc.ServiceDesc, any) {
		bankService = services.NewDefaultBankService(appCtx.TransactionHandler)
		grpcServer := rpc.NewBankApiGrpcServer(appCtx.AuthenticationService, appCtx.AuthorizationService, bankService)
		return &rpc.Api_ServiceDesc, grpcServer
	}

	bankRestServer := rest.NewDefaultBankApiRestServer(bankService)
	err := boot.Init(appName, version, args, logger, enablers, builder, func(appCtx boot.ApplicationContext) error {

		appCtx.PrivateRouter.POST("/accounts", bankRestServer.CreateAccount)
		appCtx.PrivateRouter.POST("/transfers", bankRestServer.Transfer)
		appCtx.PrivateRouter.GET("/accounts/:number", bankRestServer.GetAccount)
		appCtx.PrivateRouter.GET("/accounts/:number/entries", bankRestServer.GetAccountWithEntries)
		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
