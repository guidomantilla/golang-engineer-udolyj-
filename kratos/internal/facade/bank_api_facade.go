package facade

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"kratos/api"

	"kratos/internal/biz"
)

type BankApiFacade struct {
	api.UnimplementedApiServer

	uc *biz.GreeterUsecase
}

func NewBankService(uc *biz.GreeterUsecase) *BankApiFacade {
	return &BankApiFacade{uc: uc}
}

func (service *BankApiFacade) Login(context.Context, *api.LoginRequest) (*api.LoginResponse, error) {
	return nil, nil
}

func (service *BankApiFacade) Health(context.Context, *emptypb.Empty) (*api.HealthResponse, error) {
	return nil, nil
}

func (service *BankApiFacade) Info(context.Context, *emptypb.Empty) (*api.InfoResponse, error) {
	return nil, nil
}

func (service *BankApiFacade) CreateAccount(context.Context, *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	return nil, nil
}

func (service *BankApiFacade) Transfer(context.Context, *api.TransferRequest) (*api.TransferResponse, error) {
	return nil, nil
}

func (service *BankApiFacade) GetAccount(context.Context, *api.GetAccountRequest) (*api.GetAccountResponse, error) {
	return nil, nil
}

func (service *BankApiFacade) GetAccountWithEntries(context.Context, *api.GetAccountRequest) (*api.GetAccountWithEntriesResponse, error) {
	return nil, nil
}
