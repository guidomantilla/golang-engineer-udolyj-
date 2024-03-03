package rest

import (
	"github.com/gin-gonic/gin"
)

type BankApiRestServer interface {
	CreateAccount(gctx *gin.Context)
	Transfer(gctx *gin.Context)
	GetAccount(gctx *gin.Context)
	GetAccountWithEntries(gctx *gin.Context)
}
