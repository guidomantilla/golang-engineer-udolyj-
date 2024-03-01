package security

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type ResourceCtxKey struct{}

type Principal struct {
	Username           *string  `json:"username,omitempty" binding:"required"`
	Role               *string  `json:"role,omitempty"`
	Password           *string  `json:"password,omitempty" binding:"required"`
	Passphrase         *string  `json:"passphrase,omitempty" `
	Enabled            *bool    `json:"enabled,omitempty"`
	NonLocked          *bool    `json:"non_locked,omitempty"`
	NonExpired         *bool    `json:"non_expired,omitempty"`
	PasswordNonExpired *bool    `json:"password_non_expired,omitempty"`
	SignUpDone         *bool    `json:"signup_done,omitempty"`
	Resources          []string `json:"resources,omitempty"`
	Token              *string  `json:"token,omitempty"`
}

//

var (
	_ PrincipalManager       = (*InMemoryPrincipalManager)(nil)
	_ AuthenticationEndpoint = (*DefaultAuthenticationEndpoint)(nil)
	_ AuthenticationService  = (*DefaultAuthenticationService)(nil)
	_ AuthorizationFilter    = (*DefaultAuthorizationFilter)(nil)
	_ AuthorizationService   = (*DefaultAuthorizationService)(nil)
	_ TokenManager           = (*JwtTokenManager)(nil)
)

type PrincipalManager interface {
	Create(ctx context.Context, principal *Principal) error
	Update(ctx context.Context, principal *Principal) error
	Delete(ctx context.Context, username string) error
	Find(ctx context.Context, username string) (*Principal, error)
	Exists(ctx context.Context, username string) error

	ChangePassword(ctx context.Context, username string, password string) error
	VerifyResource(ctx context.Context, username string, resource string) error
}

//

type AuthenticationEndpoint interface {
	Authenticate(ctx *gin.Context)
}

type AuthenticationService interface {
	Authenticate(ctx context.Context, principal *Principal) error
	Validate(principal *Principal) []error
}

//

type AuthorizationFilter interface {
	Authorize(ctx *gin.Context)
}

type AuthorizationService interface {
	Authorize(ctx context.Context, tokenString string) (*Principal, error)
}

//

type TokenManager interface {
	Generate(principal *Principal) (*string, error)
	Validate(tokenString string) (*Principal, error)
}

//

const (
	Argon2PrefixKey = "{argon2}"
	BcryptPrefixKey = "{bcrypt}"
	Pbkdf2PrefixKey = "{pbkdf2}"
	ScryptPrefixKey = "{scrypt}"
)

var (
	_ PasswordEncoder   = (*Argon2PasswordEncoder)(nil)
	_ PasswordEncoder   = (*BcryptPasswordEncoder)(nil)
	_ PasswordEncoder   = (*Pbkdf2PasswordEncoder)(nil)
	_ PasswordEncoder   = (*ScryptPasswordEncoder)(nil)
	_ PasswordEncoder   = (*DelegatingPasswordEncoder)(nil)
	_ PasswordEncoder   = (*DefaultPasswordManager)(nil)
	_ PasswordGenerator = (*DefaultPasswordGenerator)(nil)
	_ PasswordGenerator = (*DefaultPasswordManager)(nil)
	_ PasswordManager   = (*DefaultPasswordManager)(nil)
)

type PasswordEncoder interface {
	Encode(rawPassword string) (*string, error)
	Matches(encodedPassword string, rawPassword string) (*bool, error)
	UpgradeEncoding(encodedPassword string) (*bool, error)
}

type PasswordGenerator interface {
	Generate() string
	Validate(rawPassword string) error
}

type PasswordManager interface {
	PasswordEncoder
	PasswordGenerator
}
