package services

import (
	"context"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/kratos/internal/models"
)

type BankService interface {

	// CreateAccount Create a new bank account for a customer, with an initial deposit amount.
	// A single customer may have multiple bank accounts.
	CreateAccount(ctx context.Context, account *models.Account) error

	// GetAccount Retrieve balances for a given account.
	GetAccount(ctx context.Context, id *string) (*models.Account, error)

	// GetAccountWithEntries Retrieve transfer history for a given account.
	GetAccountWithEntries(ctx context.Context, id *string) (*models.Account, error)

	// Transfer amounts between any two accounts, including those owned by different customers.
	Transfer(ctx context.Context, transfer *models.Transfer) error
}
