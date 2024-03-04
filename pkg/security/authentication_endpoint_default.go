package security

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/log"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/pkg/rest"
)

type DefaultAuthenticationEndpoint struct {
	authenticationService AuthenticationService
}

func NewDefaultAuthenticationEndpoint(authenticationService AuthenticationService) *DefaultAuthenticationEndpoint {

	if authenticationService == nil {
		log.Fatal("starting up - error setting up authenticationEndpoint: authenticationService is nil")
	}

	return &DefaultAuthenticationEndpoint{
		authenticationService: authenticationService,
	}
}

func (endpoint *DefaultAuthenticationEndpoint) Authenticate(ctx *gin.Context) {

	var err error
	var principal *Principal
	if err = ctx.ShouldBindJSON(&principal); err != nil {
		ex := rest.BadRequestException("error unmarshalling request json to object")
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if errs := endpoint.authenticationService.Validate(principal); errs != nil {
		ex := rest.BadRequestException("error validating the principal", errs...)
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	if err = endpoint.authenticationService.Authenticate(ctx.Request.Context(), principal); err != nil {
		ex := rest.UnauthorizedException(err.Error())
		ctx.AbortWithStatusJSON(ex.Code, ex)
		return
	}

	ctx.JSON(http.StatusOK, principal)
}
