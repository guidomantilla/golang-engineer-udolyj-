package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/util"
)

// Customer represents a customer of the bank

type Customer struct {
	Id    *int64  `gorm:"primaryKey" json:"id,omitempty"`
	Name  *string `gorm:"column:name" json:"name,omitempty"`
	Email *string `gorm:"column:email" json:"email,omitempty"`
	Phone *string `gorm:"column:phone" json:"phone,omitempty"`
}

func (customer *Customer) TableName() string {
	return "core_customers"
}

// Account represents a bank account

type Account struct {
	Id         *string   `gorm:"primaryKey" json:"id,omitempty"`
	CustomerId *int64    `gorm:"column:owner" json:"ownerId,omitempty"`
	Customer   *Customer `gorm:"foreignKey:CustomerId" json:"owner,omitempty"`
	Balance    *float64  `gorm:"column:balance" json:"balance,omitempty"`
	Entries    []*Entry  `gorm:"foreignKey:AccountId" json:"entries,omitempty"`
}

func (account *Account) TableName() string {
	return "core_accounts"
}

func (account *Account) BeforeCreate(_ *gorm.DB) (err error) {
	account.Id = util.ValueToPtr(uuid.New().String())
	return
}

// Entry represents a bank account entry

type Entry struct {
	Id        *string  `gorm:"primaryKey" json:"id,omitempty"`
	AccountId *string  `gorm:"column:account_id" json:"account,omitempty"`
	Amount    *float64 `gorm:"column:amount" json:"amount,omitempty"`
}

func (entry *Entry) TableName() string {
	return "core_entries"
}

func (entry *Entry) BeforeCreate(_ *gorm.DB) (err error) {
	entry.Id = util.ValueToPtr(uuid.New().String())
	return
}

// Transfer represents a bank account transfer

type Transfer struct {
	Id            *string  `gorm:"primaryKey" json:"id,omitempty"`
	FromAccountId *string  `gorm:"column:from_account_id" json:"from,omitempty"`
	ToAccountId   *string  `gorm:"column:to_account_id" json:"to,omitempty"`
	Amount        *float64 `gorm:"column:amount" json:"amount,omitempty"`
}

func (transfer *Transfer) TableName() string {
	return "core_transfers"
}

func (transfer *Transfer) BeforeCreate(_ *gorm.DB) (err error) {
	transfer.Id = util.ValueToPtr(uuid.New().String())
	return
}
