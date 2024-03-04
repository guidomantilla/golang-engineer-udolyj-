package datasource

import (
	"context"

	"gorm.io/gorm"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
)

type DefaultTransactionHandler struct {
	datasource Datasource
}

func NewTransactionHandler(datasource Datasource) *DefaultTransactionHandler {
	return &DefaultTransactionHandler{datasource: datasource}
}

func (handler *DefaultTransactionHandler) HandleTransaction(ctx context.Context, fn TransactionHandlerFunction) error {
	dbx, err := handler.datasource.GetDatabase()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return dbx.Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}
