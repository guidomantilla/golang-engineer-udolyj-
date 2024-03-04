package rpc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/models"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/services"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/errors"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/security"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/util"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/validation"
)

type BankApiGrpcServer struct {
	authenticationService security.AuthenticationService
	authorizationService  security.AuthorizationService
	bankService           services.BankService
}

func (server *BankApiGrpcServer) mustEmbedUnimplementedApiServer() {
}

func NewBankApiGrpcServer(authenticationService security.AuthenticationService, authorizationService security.AuthorizationService, bankService services.BankService) *BankApiGrpcServer {
	return &BankApiGrpcServer{
		authenticationService: authenticationService,
		authorizationService:  authorizationService,
		bankService:           bankService,
	}
}

func (server *BankApiGrpcServer) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {

	principal := &security.Principal{
		Username: util.ValueToPtr(request.Username),
		Password: util.ValueToPtr(request.Password),
	}
	if errs := server.authenticationService.Validate(principal); errs != nil {
		errs = append([]error{fmt.Errorf("error unmarshalling request json to object")}, errs...)
		return nil, errors.ErrJoin(errs...)
	}

	if err := server.authenticationService.Authenticate(ctx, principal); err != nil {
		return nil, err
	}

	response := &LoginResponse{
		Username:  *principal.Username,
		Role:      *principal.Role,
		Resources: principal.Resources,
		Token:     *principal.Token,
	}
	return response, nil
}

func (server *BankApiGrpcServer) Health(_ context.Context, _ *emptypb.Empty) (*HealthResponse, error) {
	response := &HealthResponse{
		Status: "alive",
	}
	return response, nil
}

func (server *BankApiGrpcServer) Info(_ context.Context, _ *emptypb.Empty) (*InfoResponse, error) {
	response := &InfoResponse{
		AppName: config.Application,
	}
	return response, nil
}

// CreateAccount Create a new bank account for a customer, with an initial deposit amount.
// A single customer may have multiple bank accounts.
func (server *BankApiGrpcServer) CreateAccount(ctx context.Context, request *CreateAccountRequest) (*CreateAccountResponse, error) {

	ctx, err := server.authorize(ctx)
	if err != nil {
		return nil, err
	}

	accountToSave := &models.Account{
		CustomerId: util.ValueToPtr(request.OwnerId),
		Balance:    util.ValueToPtr(request.Balance),
	}
	if errs := server.validateCreateAccount(accountToSave); errs != nil {
		errs = append([]error{fmt.Errorf("error validating the object")}, errs...)
		return nil, errors.ErrJoin(errs...)
	}

	if err = server.bankService.CreateAccount(ctx, accountToSave); err != nil {
		return nil, err
	}

	response := &CreateAccountResponse{
		Id:      *accountToSave.Id,
		OwnerId: *accountToSave.CustomerId,
		Balance: *accountToSave.Balance,
	}
	return response, nil
}

func (server *BankApiGrpcServer) validateCreateAccount(account *models.Account) []error {

	var errors []error

	if err := validation.ValidateFieldMustBeUndefined("this", "id", account.Id); err != nil {
		errors = append(errors, err)
	}
	if err := validation.ValidateFieldIsRequired("this", "owner", account.CustomerId); err != nil {
		errors = append(errors, err)
	}

	if err := validation.ValidateFieldIsRequired("this", "balance", account.Balance); err != nil {
		errors = append(errors, err)
	}

	return errors
}

// Transfer amounts between any two accounts, including those owned by different customers.
func (server *BankApiGrpcServer) Transfer(ctx context.Context, request *TransferRequest) (*TransferResponse, error) {

	ctx, err := server.authorize(ctx)
	if err != nil {
		return nil, err
	}

	transferToSave := &models.Transfer{
		FromAccountId: util.ValueToPtr(request.From),
		ToAccountId:   util.ValueToPtr(request.To),
		Amount:        util.ValueToPtr(request.Amount),
	}
	if errs := server.validateTransfer(transferToSave); errs != nil {
		errs = append([]error{fmt.Errorf("error validating the object")}, errs...)
		return nil, errors.ErrJoin(errs...)
	}

	if err = server.bankService.Transfer(ctx, transferToSave); err != nil {
		return nil, err
	}

	response := &TransferResponse{
		Id:     *transferToSave.Id,
		From:   *transferToSave.FromAccountId,
		To:     *transferToSave.ToAccountId,
		Amount: *transferToSave.Amount,
	}
	return response, nil
}

func (server *BankApiGrpcServer) validateTransfer(transfer *models.Transfer) []error {

	var errors []error

	if err := validation.ValidateFieldMustBeUndefined("this", "id", transfer.Id); err != nil {
		errors = append(errors, err)
	}
	if err := validation.ValidateFieldIsRequired("this", "from", transfer.FromAccountId); err != nil {
		errors = append(errors, err)
	}

	if err := validation.ValidateFieldIsRequired("this", "to", transfer.ToAccountId); err != nil {
		errors = append(errors, err)
	}

	if err := validation.ValidateFieldIsRequired("this", "amount", transfer.Amount); err != nil {
		errors = append(errors, err)
	}

	return errors
}

// GetAccount Retrieve balances for a given account.
func (server *BankApiGrpcServer) GetAccount(ctx context.Context, request *GetAccountRequest) (*GetAccountResponse, error) {

	ctx, err := server.authorize(ctx)
	if err != nil {
		return nil, err
	}

	var account *models.Account
	if account, err = server.bankService.GetAccount(ctx, &request.Id); err != nil {
		return nil, err
	}

	response := &GetAccountResponse{
		Id: *account.Id,
		Owner: &Customer{
			Name:  *account.Customer.Name,
			Email: *account.Customer.Email,
			Phone: *account.Customer.Phone,
		},
		Balance: *account.Balance,
	}

	return response, nil
}

// GetAccountWithEntries Retrieve transfer history for a given account
func (server *BankApiGrpcServer) GetAccountWithEntries(ctx context.Context, request *GetAccountRequest) (*GetAccountWithEntriesResponse, error) {

	ctx, err := server.authorize(ctx)
	if err != nil {
		return nil, err
	}

	var account *models.Account
	if account, err = server.bankService.GetAccount(ctx, &request.Id); err != nil {
		return nil, err
	}

	response := &GetAccountWithEntriesResponse{
		Id: *account.Id,
		Owner: &Customer{
			Name:  *account.Customer.Name,
			Email: *account.Customer.Email,
			Phone: *account.Customer.Phone,
		},
		Entries: func() []*Entry {
			var entries []*Entry
			for _, entry := range account.Entries {
				entries = append(entries, &Entry{
					Id:      *entry.Id,
					Account: *entry.AccountId,
					Amount:  *entry.Amount,
				})
			}
			return entries
		}(),
		Balance: *account.Balance,
	}

	return response, nil
}

//

func (server *BankApiGrpcServer) authorize(ctx context.Context) (context.Context, error) {

	var ok bool
	var md metadata.MD
	if md, ok = metadata.FromIncomingContext(ctx); !ok {
		return ctx, fmt.Errorf("invalid authorization header")
	}
	if len(md.Get("Authorization")) == 0 {
		return ctx, fmt.Errorf("invalid authorization header")
	}

	header := strings.Fields(md.Get("Authorization")[0])
	if len(header) < 2 {
		return ctx, fmt.Errorf("invalid authorization header")
	}

	if !strings.HasPrefix(header[0], "Bearer") {
		return ctx, fmt.Errorf("invalid authorization header")
	}

	token := header[1]
	serverStream := grpc.ServerTransportStreamFromContext(ctx)
	resource := []string{config.Application, "N/A", serverStream.Method()}

	var err error
	var principal *security.Principal
	ctxWithResource := context.WithValue(ctx, security.ResourceCtxKey{}, strings.Join(resource, " "))
	if principal, err = server.authorizationService.Authorize(ctxWithResource, token); err != nil {
		return ctx, err
	}

	PrincipalCtxKey := struct {
		PrincipalCtxKey string
	}{
		PrincipalCtxKey: security.PrincipalCtxKey,
	}
	ctx = context.WithValue(ctx, PrincipalCtxKey, principal)
	return ctx, nil
}
