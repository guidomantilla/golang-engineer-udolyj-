package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/models"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/datasource"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/mocks"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/util"
)

func TestNewDefaultBankService(t *testing.T) {

	transactionHandler, _ := mocks.BuildMockGormTransactionHandler()

	type args struct {
		transactionHandler datasource.TransactionHandler
	}
	tests := []struct {
		name string
		args args
		want *DefaultBankService
	}{
		{
			name: "Test Default Bank Service Created Successfully",
			args: args{
				transactionHandler: transactionHandler,
			},
			want: NewDefaultBankService(transactionHandler),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultBankService(tt.args.transactionHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultBankService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBankService_CreateAccount(t *testing.T) {

	type args struct {
		ctx     context.Context
		account *models.Account
	}
	tests := []struct {
		name    string
		service *DefaultBankService
		args    args
		wantErr bool
	}{
		{
			name: "Test Create Account Successfully",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `core_accounts` (`id`,`owner`,`balance`) VALUES (?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				account: &models.Account{
					CustomerId: util.ValueToPtr(int64(1)),
					Balance:    util.ValueToPtr(float64(9)),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Create Account with negative balance",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `core_accounts` (`id`,`owner`,`balance`) VALUES (?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				account: &models.Account{
					CustomerId: util.ValueToPtr(int64(1)),
					Balance:    util.ValueToPtr(float64(-9)),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service.CreateAccount(tt.args.ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("DefaultBankService.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultBankService_GetAccount(t *testing.T) {
	type args struct {
		ctx context.Context
		id  *string
	}
	tests := []struct {
		name    string
		service *DefaultBankService
		args    args
		want    *models.Account
		wantErr bool
	}{
		{
			name: "Test Get Account Successfully",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 1))
				mock.ExpectQuery("SELECT * FROM `core_customers` WHERE `core_customers`.`id` = ?").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
						AddRow(1, "some_name", "some_email", "some_phone"))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				id:  util.ValueToPtr("some_id"),
			},
			want: &models.Account{
				Id: util.ValueToPtr("some_id"),
				Customer: &models.Customer{
					Id:    util.ValueToPtr(int64(1)),
					Name:  util.ValueToPtr("some_name"),
					Email: util.ValueToPtr("some_email"),
					Phone: util.ValueToPtr("some_phone"),
				},
				Balance: util.ValueToPtr(float64(1)),
			},
			wantErr: false,
		},
		{
			name: "Test Get Account Error",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(errors.New("some error"))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				id:  util.ValueToPtr("some_id"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.service.GetAccount(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultBankService.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultBankService.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBankService_GetAccountWithEntries(t *testing.T) {
	type args struct {
		ctx context.Context
		id  *string
	}
	tests := []struct {
		name    string
		service *DefaultBankService
		args    args
		want    *models.Account
		wantErr bool
	}{
		{
			name: "Test Get Account With Entries Successfully",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 1))
				mock.ExpectQuery("SELECT * FROM `core_customers` WHERE `core_customers`.`id` = ?").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
						AddRow(1, "some_name", "some_email", "some_phone"))
				mock.ExpectQuery("SELECT * FROM `core_entries` WHERE `core_entries`.`account_id` = ?").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "amount"}).
						AddRow("some_id", "some_id", 1))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				id:  util.ValueToPtr("some_id"),
			},
			want: &models.Account{
				Id: util.ValueToPtr("some_id"),
				Customer: &models.Customer{
					Id:    util.ValueToPtr(int64(1)),
					Name:  util.ValueToPtr("some_name"),
					Email: util.ValueToPtr("some_email"),
					Phone: util.ValueToPtr("some_phone"),
				},
				Balance: util.ValueToPtr(float64(1)),
				Entries: []*models.Entry{
					{
						Id:        util.ValueToPtr("some_id"),
						AccountId: util.ValueToPtr("some_id"),
						Amount:    util.ValueToPtr(float64(1)),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test Get Account With Entries Error",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(errors.New("some error"))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				id:  util.ValueToPtr("some_id"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.service.GetAccountWithEntries(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultBankService.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultBankService.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBankService_Transfer(t *testing.T) {
	type args struct {
		ctx      context.Context
		transfer *models.Transfer
	}
	tests := []struct {
		name    string
		service *DefaultBankService
		args    args
		wantErr bool
	}{
		{
			name: "Test Transfer Error: from account not found",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnError(errors.New("some error"))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				transfer: &models.Transfer{
					FromAccountId: util.ValueToPtr("some_id"),
					ToAccountId:   util.ValueToPtr("some_id"),
					Amount:        util.ValueToPtr(float64(1)),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Transfer Error: from account cannot have negative balance",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, -1))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				transfer: &models.Transfer{
					FromAccountId: util.ValueToPtr("some_id"),
					ToAccountId:   util.ValueToPtr("some_id"),
					Amount:        util.ValueToPtr(float64(1)),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Transfer Error: from account with insufficient funds",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 1))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				transfer: &models.Transfer{
					FromAccountId: util.ValueToPtr("some_id"),
					ToAccountId:   util.ValueToPtr("some_id"),
					Amount:        util.ValueToPtr(float64(1)),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Transfer Error: to account not found",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 10))
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnError(errors.New("some error"))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				transfer: &models.Transfer{
					FromAccountId: util.ValueToPtr("some_id"),
					ToAccountId:   util.ValueToPtr("some_id"),
					Amount:        util.ValueToPtr(float64(1)),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Transfer Error: to account cannot have negative balance",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 10))
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, -10))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				transfer: &models.Transfer{
					FromAccountId: util.ValueToPtr("some_id"),
					ToAccountId:   util.ValueToPtr("some_id"),
					Amount:        util.ValueToPtr(float64(1)),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Transfer Error: execute transaction",
			service: func() *DefaultBankService {
				transactionHandler, mock := mocks.BuildMockGormTransactionHandler()
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 10))
				mock.ExpectQuery("SELECT * FROM `core_accounts` WHERE `core_accounts`.`id` = ? ORDER BY `core_accounts`.`id` LIMIT 1").
					WithArgs("some_id").
					WithArgs("some_id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "owner", "balance"}).
						AddRow("some_id", 1, 10))
				mock.ExpectExec("INSERT INTO `core_transfers` (`id`,`from_account_id`,`to_account_id`,`amount`) VALUES (?,?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `core_accounts` SET `owner`=?,`balance`=? WHERE `id` = ?").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `core_entries` (`id`,`account_id`,`amount`) VALUES (?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `core_accounts` SET `owner`=?,`balance`=? WHERE `id` = ?").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `core_entries` (`id`,`account_id`,`amount`) VALUES (?,?,?)").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				return NewDefaultBankService(transactionHandler)
			}(),
			args: args{
				ctx: context.TODO(),
				transfer: &models.Transfer{
					FromAccountId: util.ValueToPtr("some_id"),
					ToAccountId:   util.ValueToPtr("some_id"),
					Amount:        util.ValueToPtr(float64(1)),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service.Transfer(tt.args.ctx, tt.args.transfer); (err != nil) != tt.wantErr {
				t.Errorf("DefaultBankService.Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
