package services

import (
	"context"
	"errors"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/models"

	"gorm.io/gorm"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/datasource"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/util"
)

var (
	_ BankService = (*DefaultBankService)(nil)
)

type DefaultBankService struct {
	transactionHandler datasource.TransactionHandler
}

func NewDefaultBankService(transactionHandler datasource.TransactionHandler) *DefaultBankService {
	return &DefaultBankService{
		transactionHandler: transactionHandler,
	}
}

// CreateAccount Create a new bank account for a customer, with an initial deposit amount.
// A single customer may have multiple bank accounts.
func (service *DefaultBankService) CreateAccount(ctx context.Context, account *models.Account) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		if *(account.Balance) < 0 {
			return errors.New("account cannot be created with negative balance")
		}
		return tx.Create(account).Error
	})
}

// GetAccount Retrieve balances for a given account.
func (service *DefaultBankService) GetAccount(ctx context.Context, id *string) (*models.Account, error) {

	account := &models.Account{
		Id: id,
	}
	err := service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		return tx.Preload("Customer").First(&account).Error
	})
	if err != nil {
		return nil, err
	}

	account.CustomerId = nil
	return account, nil
}

// GetAccountWithEntries Retrieve transfer history for a given account.
func (service *DefaultBankService) GetAccountWithEntries(ctx context.Context, id *string) (*models.Account, error) {

	account := &models.Account{
		Id: id,
	}
	err := service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		return tx.Preload("Customer").Preload("Entries").First(&account).Error
	})
	if err != nil {
		return nil, err
	}

	account.CustomerId = nil
	return account, nil
}

// Transfer amounts between any two accounts, including those owned by different customers.
func (service *DefaultBankService) Transfer(ctx context.Context, transfer *models.Transfer) error {
	return service.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		// Retrieve and validate from account
		fromAccount := &models.Account{
			Id: transfer.FromAccountId,
		}
		if err := tx.First(fromAccount).Error; err != nil {
			return errors.New("from account not found")
		}

		if *(fromAccount.Balance) < 0 {
			return errors.New("from account cannot have negative balance")
		}

		if *(fromAccount.Balance)-*(transfer.Amount) <= 0 {
			return errors.New("from account with insufficient funds")
		}

		// Retrieve and validate to account
		toAccount := &models.Account{
			Id: transfer.ToAccountId,
		}
		if err := tx.First(toAccount).Error; err != nil {
			return errors.New("to account not found")
		}

		if *(toAccount.Balance) < 0 {
			return errors.New("to account cannot have negative balance")
		}

		// Execute transaction
		if err := tx.Create(transfer).Error; err != nil {
			return err
		}

		fromAccount.Balance = util.ValueToPtr(*(toAccount.Balance) - *(transfer.Amount))
		if err := tx.Save(fromAccount).Error; err != nil {
			return err
		}

		fromAccountEntry := &models.Entry{
			AccountId: fromAccount.Id,
			Amount:    util.ValueToPtr(-*(transfer.Amount)),
		}
		if err := tx.Save(fromAccountEntry).Error; err != nil {
			return err
		}

		toAccount.Balance = util.ValueToPtr(*(toAccount.Balance) + *(transfer.Amount))
		if err := tx.Save(toAccount).Error; err != nil {
			return err
		}

		toAccountEntry := &models.Entry{
			AccountId: toAccount.Id,
			Amount:    util.ValueToPtr(*(transfer.Amount)),
		}
		if err := tx.Save(toAccountEntry).Error; err != nil {
			return err
		}

		return nil
	})
}
