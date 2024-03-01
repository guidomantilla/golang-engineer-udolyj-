package rpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/security"
)

type BankApiGrpcServer struct {
	authenticationService security.AuthenticationService
	authorizationService  security.AuthorizationService
	principalManager      security.PrincipalManager
}

func (server *BankApiGrpcServer) mustEmbedUnimplementedBankApiServer() {
}

func NewBankApiGrpcServer(authenticationService security.AuthenticationService, authorizationService security.AuthorizationService, principalManager security.PrincipalManager) *BankApiGrpcServer {
	return &BankApiGrpcServer{
		authenticationService: authenticationService,
		authorizationService:  authorizationService,
		principalManager:      principalManager,
	}
}

func (server *BankApiGrpcServer) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}

func (server *BankApiGrpcServer) GetPrincipal(ctx context.Context, _ *emptypb.Empty) (*Principal, error) {
	return nil, nil
}
