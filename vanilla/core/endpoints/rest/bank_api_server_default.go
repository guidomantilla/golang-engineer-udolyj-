package rest

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/models"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/core/services"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/rest"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/validation"
)

type DefaultBankApiRestServer struct {
	bankService services.BankService
}

func NewDefaultBankApiRestServer(bankService services.BankService) *DefaultBankApiRestServer {
	return &DefaultBankApiRestServer{
		bankService: bankService,
	}
}

// CreateAccount Create a new bank account for a customer, with an initial deposit amount.
// A single customer may have multiple bank accounts.
func (server *DefaultBankApiRestServer) CreateAccount(gctx *gin.Context) {

	var err error
	var accountToSave *models.Account
	if err = gctx.ShouldBindJSON(&accountToSave); err != nil {
		ex := rest.BadRequestException("error unmarshalling request json to object")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := server.validateCreateAccount(accountToSave); errs != nil {
		ex := rest.BadRequestException("error validating the object", errs...)
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = server.bankService.CreateAccount(gctx.Request.Context(), accountToSave); err != nil {
		ex := rest.UnauthorizedException(err.Error())
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	gctx.JSON(http.StatusCreated, accountToSave)
}

func (server *DefaultBankApiRestServer) validateCreateAccount(account *models.Account) []error {

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
func (server *DefaultBankApiRestServer) Transfer(gctx *gin.Context) {

	var err error
	var transferToSave *models.Transfer
	if err = gctx.ShouldBindJSON(&transferToSave); err != nil {
		ex := rest.BadRequestException("error unmarshalling request json to object")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := server.validateTransfer(transferToSave); errs != nil {
		ex := rest.BadRequestException("error validating the object", errs...)
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = server.bankService.Transfer(gctx.Request.Context(), transferToSave); err != nil {
		ex := rest.UnauthorizedException(err.Error())
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	gctx.JSON(http.StatusOK, transferToSave)
}

func (server *DefaultBankApiRestServer) validateTransfer(transfer *models.Transfer) []error {

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
func (server *DefaultBankApiRestServer) GetAccount(gctx *gin.Context) {

	var err error
	var body []byte
	if body, err = io.ReadAll(gctx.Request.Body); err != nil {
		ex := rest.BadRequestException("error reading body")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if len(body) != 0 {
		ex := rest.BadRequestException("body is not empty")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	id := gctx.Params.ByName("number")
	if id == "" {
		ex := rest.BadRequestException("object id not defined in url path")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	var account *models.Account
	if account, err = server.bankService.GetAccount(gctx.Request.Context(), &id); err != nil {
		ex := rest.UnauthorizedException(err.Error())
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	gctx.JSON(http.StatusOK, account)
}

// GetAccountWithEntries Retrieve transfer history for a given account.
func (server *DefaultBankApiRestServer) GetAccountWithEntries(gctx *gin.Context) {

	var err error
	var body []byte
	if body, err = io.ReadAll(gctx.Request.Body); err != nil {
		ex := rest.BadRequestException("error reading body")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if len(body) != 0 {
		ex := rest.BadRequestException("body is not empty")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	id := gctx.Params.ByName("number")
	if id == "" {
		ex := rest.BadRequestException("object id not defined in url path")
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	var account *models.Account
	if account, err = server.bankService.GetAccountWithEntries(gctx.Request.Context(), &id); err != nil {
		ex := rest.UnauthorizedException(err.Error())
		gctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	gctx.JSON(http.StatusOK, account)
}
