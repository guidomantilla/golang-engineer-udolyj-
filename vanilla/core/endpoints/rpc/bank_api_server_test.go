package rpc

import (
	"context"
	"reflect"
	"testing"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/security"
	"google.golang.org/protobuf/types/known/emptypb"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/models"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/services"
)

func TestBankApiGrpcServer_mustEmbedUnimplementedApiServer(t *testing.T) {
	tests := []struct {
		name   string
		server *BankApiGrpcServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.mustEmbedUnimplementedApiServer()
		})
	}
}

func TestNewBankApiGrpcServer(t *testing.T) {
	type args struct {
		authenticationService security.AuthenticationService
		authorizationService  security.AuthorizationService
		bankService           services.BankService
	}
	tests := []struct {
		name string
		args args
		want *BankApiGrpcServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBankApiGrpcServer(tt.args.authenticationService, tt.args.authorizationService, tt.args.bankService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBankApiGrpcServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_Login(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *LoginRequest
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *LoginResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.Login(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_Health(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *HealthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.Health(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.Health() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_Info(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *emptypb.Empty
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *InfoResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.Info(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_CreateAccount(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *CreateAccountRequest
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *CreateAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.CreateAccount(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.CreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_validateCreateAccount(t *testing.T) {
	type args struct {
		account *models.Account
	}
	tests := []struct {
		name   string
		server *BankApiGrpcServer
		args   args
		want   []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.validateCreateAccount(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.validateCreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_Transfer(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *TransferRequest
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *TransferResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.Transfer(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.Transfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_validateTransfer(t *testing.T) {
	type args struct {
		transfer *models.Transfer
	}
	tests := []struct {
		name   string
		server *BankApiGrpcServer
		args   args
		want   []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.validateTransfer(tt.args.transfer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.validateTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_GetAccount(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *GetAccountRequest
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *GetAccountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.GetAccount(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_GetAccountWithEntries(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *GetAccountRequest
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    *GetAccountWithEntriesResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.GetAccountWithEntries(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.GetAccountWithEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.GetAccountWithEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBankApiGrpcServer_authorize(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		server  *BankApiGrpcServer
		args    args
		want    context.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.server.authorize(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BankApiGrpcServer.authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BankApiGrpcServer.authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
