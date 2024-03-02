package security

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/config"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/datasource"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/log"
	"git.codesubmit.io/stena-group/golang-engineer-udolyj/vanilla/pkg/util"
)

type GormPrincipalManager struct {
	transactionHandler datasource.TransactionHandler
	passwordManager    PasswordManager
}

func NewGormPrincipalManager(transactionHandler datasource.TransactionHandler, passwordManager PasswordManager) *GormPrincipalManager {

	if transactionHandler == nil {
		log.Fatal("starting up - error setting up principalManager: transactionHandler is nil")
	}

	if passwordManager == nil {
		log.Fatal("starting up - error setting up principalManager: passwordManager is nil")
	}

	return &GormPrincipalManager{
		transactionHandler: transactionHandler,
		passwordManager:    passwordManager,
	}
}

func (manager *GormPrincipalManager) Create(_ context.Context, _ *Principal) error {
	panic("not implemented. no required for this coding challenge")
}

func (manager *GormPrincipalManager) Update(_ context.Context, _ *Principal) error {
	panic("not implemented. no required for this coding challenge")
}

func (manager *GormPrincipalManager) Delete(_ context.Context, _ string) error {
	panic("not implemented. no required for this coding challenge")
}

func (manager *GormPrincipalManager) Find(ctx context.Context, username string) (*Principal, error) {

	var user *Principal
	err := manager.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		var authPrincipals []AuthPrincipal
		if err := tx.Find(&authPrincipals, "username = ? AND application = ?", username, config.Application).Error; err != nil {
			return err
		}
		if len(authPrincipals) == 0 {
			return errors.New("principal does not exists")
		}

		principal := authPrincipals[0]
		if principal.Role == nil || *(principal.Role) == "" {
			return ErrAccountEmptyRole
		}

		if principal.Password == nil || *(principal.Password) == "" {
			return ErrAccountEmptyPassword
		}

		if principal.Enabled != nil && !*(principal.Enabled) {
			return ErrAccountDisabled
		}

		resources := make([]string, 0)
		for _, principal := range authPrincipals {
			resources = append(resources, strings.Join([]string{*principal.Application, *principal.Permission, *principal.Resource}, " "))
		}

		user = &Principal{
			Username:           authPrincipals[0].Username,
			Role:               authPrincipals[0].Role,
			Password:           authPrincipals[0].Password,
			Passphrase:         authPrincipals[0].Passphrase,
			Enabled:            authPrincipals[0].Enabled,
			NonLocked:          util.TruePrt(),
			NonExpired:         util.TruePrt(),
			PasswordNonExpired: util.TruePrt(),
			SignUpDone:         util.TruePrt(),
			Resources:          resources,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (manager *GormPrincipalManager) Exists(_ context.Context, _ string) error {
	panic("not implemented. no required for this coding challenge")
}

func (manager *GormPrincipalManager) ChangePassword(_ context.Context, _ string, _ string) error {
	panic("not implemented. no required for this coding challenge")
}

func (manager *GormPrincipalManager) VerifyResource(ctx context.Context, username string, resource string) error {

	return manager.transactionHandler.HandleTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {

		principals := make([]AuthPrincipal, 0)
		if err := tx.Find(&principals, username, resource).Error; err != nil {
			return err
		}

		return nil
	})
}
