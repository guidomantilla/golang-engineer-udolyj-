package facade

import (
	"context"
	"fmt"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/api"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/models"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/services"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/errors"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/security"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/util"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/validation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"strings"
)

var (
	_ api.ApiServer     = (*BankApiFacade)(nil)
	_ api.ApiHTTPServer = (*BankApiFacade)(nil)
)

type BankApiFacade struct {
	api.UnimplementedApiServer
	authenticationService security.AuthenticationService
	authorizationService  security.AuthorizationService
	bankService           services.BankService
}

func NewBankApiFacade(authenticationService security.AuthenticationService, authorizationService security.AuthorizationService, bankService services.BankService) *BankApiFacade {
	return &BankApiFacade{
		authenticationService: authenticationService,
		authorizationService:  authorizationService,
		bankService:           bankService,
	}
}

func (facade *BankApiFacade) Login(ctx context.Context, request *api.LoginRequest) (*api.LoginResponse, error) {

	principal := &security.Principal{
		Username: util.ValueToPtr(request.Username),
		Password: util.ValueToPtr(request.Password),
	}
	if errs := facade.authenticationService.Validate(principal); errs != nil {
		errs = append([]error{fmt.Errorf("error unmarshalling request json to object")}, errs...)
		return nil, errors.ErrJoin(errs...)
	}

	if err := facade.authenticationService.Authenticate(ctx, principal); err != nil {
		return nil, err
	}

	response := &api.LoginResponse{
		Username:  *principal.Username,
		Role:      *principal.Role,
		Resources: principal.Resources,
		Token:     *principal.Token,
	}
	return response, nil
}

func (facade *BankApiFacade) Health(_ context.Context, _ *emptypb.Empty) (*api.HealthResponse, error) {
	response := &api.HealthResponse{
		Status: "alive",
	}
	return response, nil
}

func (facade *BankApiFacade) Info(_ context.Context, _ *emptypb.Empty) (*api.InfoResponse, error) {
	response := &api.InfoResponse{
		AppName: "CHANGE ME", //config.Application,
	}
	return response, nil
}

// CreateAccount Create a new bank account for a customer, with an initial deposit amount.
// A single customer may have multiple bank accounts.
func (facade *BankApiFacade) CreateAccount(ctx context.Context, request *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {

	ctx, err := facade.authorize(ctx)
	if err != nil {
		return nil, err
	}

	accountToSave := &models.Account{
		CustomerId: util.ValueToPtr(request.OwnerId),
		Balance:    util.ValueToPtr(request.Balance),
	}
	if errs := facade.validateCreateAccount(accountToSave); errs != nil {
		errs = append([]error{fmt.Errorf("error validating the object")}, errs...)
		return nil, errors.ErrJoin(errs...)
	}

	if err = facade.bankService.CreateAccount(ctx, accountToSave); err != nil {
		return nil, err
	}

	response := &api.CreateAccountResponse{
		Id:      *accountToSave.Id,
		OwnerId: *accountToSave.CustomerId,
		Balance: *accountToSave.Balance,
	}
	return response, nil
}

func (facade *BankApiFacade) validateCreateAccount(account *models.Account) []error {

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
func (facade *BankApiFacade) Transfer(ctx context.Context, request *api.TransferRequest) (*api.TransferResponse, error) {

	ctx, err := facade.authorize(ctx)
	if err != nil {
		return nil, err
	}

	transferToSave := &models.Transfer{
		FromAccountId: util.ValueToPtr(request.From),
		ToAccountId:   util.ValueToPtr(request.To),
		Amount:        util.ValueToPtr(request.Amount),
	}
	if errs := facade.validateTransfer(transferToSave); errs != nil {
		errs = append([]error{fmt.Errorf("error validating the object")}, errs...)
		return nil, errors.ErrJoin(errs...)
	}

	if err = facade.bankService.Transfer(ctx, transferToSave); err != nil {
		return nil, err
	}

	response := &api.TransferResponse{
		Id:     *transferToSave.Id,
		From:   *transferToSave.FromAccountId,
		To:     *transferToSave.ToAccountId,
		Amount: *transferToSave.Amount,
	}
	return response, nil
}

func (facade *BankApiFacade) validateTransfer(transfer *models.Transfer) []error {

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
func (facade *BankApiFacade) GetAccount(ctx context.Context, request *api.GetAccountRequest) (*api.GetAccountResponse, error) {

	ctx, err := facade.authorize(ctx)
	if err != nil {
		return nil, err
	}

	var account *models.Account
	if account, err = facade.bankService.GetAccount(ctx, &request.Id); err != nil {
		return nil, err
	}

	response := &api.GetAccountResponse{
		Id: *account.Id,
		Owner: &api.Customer{
			Name:  *account.Customer.Name,
			Email: *account.Customer.Email,
			Phone: *account.Customer.Phone,
		},
		Balance: *account.Balance,
	}

	return response, nil
}

func (facade *BankApiFacade) GetAccountWithEntries(ctx context.Context, request *api.GetAccountRequest) (*api.GetAccountWithEntriesResponse, error) {

	ctx, err := facade.authorize(ctx)
	if err != nil {
		return nil, err
	}

	var account *models.Account
	if account, err = facade.bankService.GetAccount(ctx, &request.Id); err != nil {
		return nil, err
	}

	response := &api.GetAccountWithEntriesResponse{
		Id: *account.Id,
		Owner: &api.Customer{
			Name:  *account.Customer.Name,
			Email: *account.Customer.Email,
			Phone: *account.Customer.Phone,
		},
		Entries: func() []*api.Entry {
			var entries []*api.Entry
			for _, entry := range account.Entries {
				entries = append(entries, &api.Entry{
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

func (facade *BankApiFacade) authorize(ctx context.Context) (context.Context, error) {

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
	resource := []string{"CHANGE ME", "N/A", serverStream.Method()}

	var err error
	var principal *security.Principal
	ctxWithResource := context.WithValue(ctx, security.ResourceCtxKey{}, strings.Join(resource, " "))
	if principal, err = facade.authorizationService.Authorize(ctxWithResource, token); err != nil {
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
