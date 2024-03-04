package boot

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/datasource"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/environment"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/security"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/util"
)

type Enablers struct {
	HttpServerEnabled bool
	GrpcServerEnabled bool
	DatabaseEnabled   bool
}

type HttpConfig struct {
	Host            *string
	Port            *string
	SwaggerPort     *string
	CorsAllowOrigin *string
}

type GrpcConfig struct {
	Host *string
	Port *string
}

type SecurityConfig struct {
	TokenSignatureKey       *string
	TokenVerificationKey    *string
	PasswordMinSpecialChars *string
	PasswordMinNumber       *string
	PasswordMinUpperCase    *string
	PasswordLength          *string
}

type DatabaseConfig struct {
	DatasourceUrl      *string
	DatasourceUsername *string
	DatasourcePassword *string
	DatasourceServer   *string
	DatasourceService  *string
}

type ApplicationContext struct {
	AppName                string
	AppVersion             string
	LogLevel               string
	CmdArgs                []string
	Enablers               *Enablers
	HttpConfig             *HttpConfig
	GrpcConfig             *GrpcConfig
	SecurityConfig         *SecurityConfig
	DatabaseConfig         *DatabaseConfig
	Logger                 log.Logger
	Environment            environment.Environment
	DatasourceContext      datasource.DatasourceContext
	Datasource             datasource.Datasource
	TransactionHandler     datasource.TransactionHandler
	PasswordEncoder        security.PasswordEncoder
	PasswordGenerator      security.PasswordGenerator
	PasswordManager        security.PasswordManager
	PrincipalManager       security.PrincipalManager
	TokenManager           security.TokenManager
	AuthenticationService  security.AuthenticationService
	AuthenticationEndpoint security.AuthenticationEndpoint
	AuthorizationService   security.AuthorizationService
	AuthorizationFilter    security.AuthorizationFilter
	PublicRouter           *gin.Engine
	PrivateRouter          *gin.RouterGroup
	GrpcServiceDesc        *grpc.ServiceDesc
	GrpcServiceServer      any
}

func NewApplicationContext(appName string, version string, args []string, logger log.Logger, enablers *Enablers, builder *BeanBuilder) *ApplicationContext {

	if appName == "" {
		log.Fatal("starting up - error setting up the ApplicationContext: appName is empty")
	}

	if version == "" {
		log.Fatal("starting up - error setting up the ApplicationContext: version is empty")
	}

	if args == nil {
		log.Fatal("starting up - error setting up the ApplicationContext: args is nil")
	}

	if logger == nil {
		log.Fatal("starting up - error setting up the application: logger is nil")
	}

	if enablers == nil {
		log.Warn("starting up - warning setting up the application: http server, grpc server and database connectivity are disabled")
		enablers = &Enablers{}
	}

	if builder == nil { //nolint:staticcheck
		log.Fatal("starting up - error setting up the ApplicationContext: builder is nil")
	}

	ctx := &ApplicationContext{
		AppName:    appName,
		AppVersion: version,
		CmdArgs:    args,
		Logger:     logger,
		Enablers:   enablers,
		SecurityConfig: &SecurityConfig{
			TokenSignatureKey:    util.ValueToPtr("SecretYouShouldHide"),
			TokenVerificationKey: util.ValueToPtr("SecretYouShouldHide"),
		},
		HttpConfig: &HttpConfig{
			Host: util.ValueToPtr("localhost"),
			Port: util.ValueToPtr("8080"),
		},
		GrpcConfig: &GrpcConfig{
			Host: util.ValueToPtr("localhost"),
			Port: util.ValueToPtr("50051"),
		},
	}

	log.Debug("starting up - setting up environment variables")
	ctx.Environment = builder.Environment(ctx) //nolint:staticcheck

	log.Debug("starting up - setting up configuration")
	builder.Config(ctx) //nolint:staticcheck

	if ctx.Enablers.DatabaseEnabled {
		log.Debug("starting up - setting up db connectivity")
		ctx.DatasourceContext = builder.DatasourceContext(ctx)   //nolint:staticcheck
		ctx.Datasource = builder.Datasource(ctx)                 //nolint:staticcheck
		ctx.TransactionHandler = builder.TransactionHandler(ctx) //nolint:staticcheck
	} else {
		log.Warn("starting up - warning setting up database configuration. database connectivity is disabled")
	}

	log.Debug("starting up - setting up security")
	ctx.PasswordEncoder = builder.PasswordEncoder(ctx)                                                                          //nolint:staticcheck
	ctx.PasswordGenerator = builder.PasswordGenerator(ctx)                                                                      //nolint:staticcheck
	ctx.PasswordManager = builder.PasswordManager(ctx)                                                                          //nolint:staticcheck
	ctx.PrincipalManager, ctx.TokenManager = builder.PrincipalManager(ctx), builder.TokenManager(ctx)                           //nolint:staticcheck
	ctx.AuthenticationService, ctx.AuthorizationService = builder.AuthenticationService(ctx), builder.AuthorizationService(ctx) //nolint:staticcheck
	ctx.AuthenticationEndpoint, ctx.AuthorizationFilter = builder.AuthenticationEndpoint(ctx), builder.AuthorizationFilter(ctx) //nolint:staticcheck

	if ctx.Enablers.HttpServerEnabled {
		log.Debug("starting up - setting up http server")
		ctx.PublicRouter, ctx.PrivateRouter = builder.HttpServer(ctx) //nolint:staticcheck
	} else {
		log.Warn("starting up - warning setting up http configuration. http server is disabled")
	}

	if ctx.Enablers.GrpcServerEnabled {
		log.Debug("starting up - setting up grpc server")
		ctx.GrpcServiceDesc, ctx.GrpcServiceServer = builder.GrpcServer(ctx) //nolint:staticcheck
	} else {
		log.Warn("starting up - warning setting up grpc configuration. grpc server is disabled")
	}

	return ctx
}

func (ctx *ApplicationContext) Stop() {

	var err error

	if ctx.Datasource != nil && ctx.DatasourceContext != nil {

		var database *gorm.DB
		log.Debug("shutting down - closing up db connection")

		if database, err = ctx.Datasource.GetDatabase(); err != nil {
			log.Error(fmt.Sprintf("shutting down - error db connection: %s", err.Error()))
			return
		}

		var db *sql.DB
		if db, err = database.DB(); err != nil {
			log.Error(fmt.Sprintf("shutting down - error db connection: %s", err.Error()))
			return
		}

		if err = db.Close(); err != nil {
			log.Error(fmt.Sprintf("shutting down - error closing db connection: %s", err.Error()))
			return
		}

		log.Debug("shutting down - db connection closed")
	}

	log.Info(fmt.Sprintf("Application %s stopped", strings.Join([]string{ctx.AppName, ctx.AppVersion}, " - ")))
}
