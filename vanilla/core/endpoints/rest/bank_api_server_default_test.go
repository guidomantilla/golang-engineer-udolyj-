package rest

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/models"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/services"
)

func TestNewDefaultBankApiRestServer(t *testing.T) {
	type args struct {
		bankService services.BankService
	}
	tests := []struct {
		name string
		args args
		want *DefaultBankApiRestServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultBankApiRestServer(tt.args.bankService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultBankApiRestServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBankApiRestServer_CreateAccount(t *testing.T) {
	type args struct {
		gctx *gin.Context
	}
	tests := []struct {
		name   string
		server *DefaultBankApiRestServer
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.CreateAccount(tt.args.gctx)
		})
	}
}

func TestDefaultBankApiRestServer_validateCreateAccount(t *testing.T) {
	type args struct {
		account *models.Account
	}
	tests := []struct {
		name   string
		server *DefaultBankApiRestServer
		args   args
		want   []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.validateCreateAccount(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultBankApiRestServer.validateCreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBankApiRestServer_Transfer(t *testing.T) {
	type args struct {
		gctx *gin.Context
	}
	tests := []struct {
		name   string
		server *DefaultBankApiRestServer
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.Transfer(tt.args.gctx)
		})
	}
}

func TestDefaultBankApiRestServer_validateTransfer(t *testing.T) {
	type args struct {
		transfer *models.Transfer
	}
	tests := []struct {
		name   string
		server *DefaultBankApiRestServer
		args   args
		want   []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.server.validateTransfer(tt.args.transfer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultBankApiRestServer.validateTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBankApiRestServer_GetAccount(t *testing.T) {
	type args struct {
		gctx *gin.Context
	}
	tests := []struct {
		name   string
		server *DefaultBankApiRestServer
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.GetAccount(tt.args.gctx)
		})
	}
}

func TestDefaultBankApiRestServer_GetAccountWithEntries(t *testing.T) {
	type args struct {
		gctx *gin.Context
	}
	tests := []struct {
		name   string
		server *DefaultBankApiRestServer
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.GetAccountWithEntries(tt.args.gctx)
		})
	}
}
