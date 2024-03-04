package providers

import (
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/api"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/facade"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/server"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/services"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/datasource"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/security"
	"github.com/google/wire"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

func ProvideGormLogger(logger log.Logger) logger.Interface {
	return slogGorm.New(slogGorm.WithLogger(logger.RetrieveLogger().(*slog.Logger)), slogGorm.WithTraceAll(), slogGorm.WithRecordNotFoundError())
}

func ProvideGormConfig(logger logger.Interface) *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger,
	}
}

var GormProviderSet = wire.NewSet(ProvideGormLogger, ProvideGormConfig)

//

func ProvideDatasourceContext(cfg *config.Config) datasource.DatasourceContext {
	return datasource.NewDefaultDatasourceContext(*cfg.DatasourceUrl, *cfg.DatasourceUsername, *cfg.DatasourcePassword, *cfg.DatasourceServer, *cfg.DatasourceService)
}

func ProvideDatasource(datasourceCtx datasource.DatasourceContext, config *gorm.Config) datasource.Datasource {
	return datasource.NewDefaultDatasource(datasourceCtx, mysql.Open(datasourceCtx.GetUrl()), config)
}

func ProvideTransactionHandler(datasrc datasource.Datasource) datasource.TransactionHandler {
	return datasource.NewTransactionHandler(datasrc)
}

var DatasourceProviderSet = wire.NewSet(ProvideDatasourceContext, ProvideDatasource, ProvideTransactionHandler)

//

func ProvidePasswordEncoder() security.PasswordEncoder {
	return security.NewBcryptPasswordEncoder()
}

func ProvidePasswordGenerator() security.PasswordGenerator {
	return security.NewDefaultPasswordGenerator()
}

func ProvidePasswordManager(encoder security.PasswordEncoder, generator security.PasswordGenerator) security.PasswordManager {
	return security.NewDefaultPasswordManager(encoder, generator)
}

func ProvidePrincipalManager(transactionHandler datasource.TransactionHandler, passwordManager security.PasswordManager) security.PrincipalManager {
	return security.NewGormPrincipalManager(transactionHandler, passwordManager)
}

func ProvideTokenManager(cfg *config.Config) security.TokenManager {
	return security.NewJwtTokenManager(security.WithIssuer(config.Application),
		security.WithSigningKey([]byte(*cfg.TokenSignatureKey)), security.WithVerifyingKey([]byte(*cfg.TokenVerificationKey)))
}

var SecurityProviderSet = wire.NewSet(ProvidePasswordEncoder, ProvidePasswordGenerator, ProvidePasswordManager, ProvidePrincipalManager, ProvideTokenManager)

//

func ProvideAuthenticationService(passwordManager security.PasswordManager, principalManager security.PrincipalManager, tokenManager security.TokenManager) security.AuthenticationService {
	return security.NewDefaultAuthenticationService(passwordManager, principalManager, tokenManager)
}

func ProvideAuthorizationService(principalManager security.PrincipalManager, tokenManager security.TokenManager) security.AuthorizationService {
	return security.NewDefaultAuthorizationService(tokenManager, principalManager)
}

var AuthProviderSet = wire.NewSet(ProvideAuthenticationService, ProvideAuthorizationService)

//

func ProvideBankService(transactionHandler datasource.TransactionHandler) services.BankService {
	return services.NewDefaultBankService(transactionHandler)
}

func ProvideBankApiFacade(authenticationService security.AuthenticationService, authorizationService security.AuthorizationService, bankService services.BankService) api.ApiServer {
	return facade.NewBankApiFacade(authenticationService, authorizationService, bankService)
}
func ProvideBankApiHttpFacade(authenticationService security.AuthenticationService, authorizationService security.AuthorizationService, bankService services.BankService) api.ApiHTTPServer {
	return facade.NewBankApiFacade(authenticationService, authorizationService, bankService)
}

//

var ProviderSet = wire.NewSet(GormProviderSet, DatasourceProviderSet, SecurityProviderSet, AuthProviderSet,
	ProvideBankService, ProvideBankApiFacade, ProvideBankApiHttpFacade,
	server.NewGRPCServer, server.NewHTTPServer)

/*
// ProviderSet is service providers.
var ProviderSet = wire.NewSet(datasource.NewTransactionHandler, security.NewBcryptPasswordEncoder,
	security.NewDefaultPasswordGenerator, security.NewDefaultPasswordManager, security.NewGormPrincipalManager,
	security.NewDefaultAuthenticationService, security.NewDefaultAuthorizationService, services.NewDefaultBankService,
	facade.NewBankApiFacade, server.NewGRPCServer, server.NewHTTPServer)
*/
