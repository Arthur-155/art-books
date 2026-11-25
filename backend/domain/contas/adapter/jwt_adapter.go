//gofi:context contas
package adapter

import (
	"time"

	"github.com/joaoprofile/gofi-sdk-go/iam/provider/jwt"
	"github.com/joaoprofile/gofi-sdk-go/iam/types"

	"github.com/Arthur-155/art-books/domain/contas/entity"
)

// ClaimEmailConfirmed é a chave do indicador de e-mail confirmado dentro do token de acesso.
const ClaimEmailConfirmed = "emailConfirmed"

// JWTAccessTokenIssuer emite o token de acesso com o provider JWT do SDK.
type JWTAccessTokenIssuer struct {
	provider *jwt.Provider
	ttl      time.Duration
	now      func() time.Time
}

func NewJWTAccessTokenIssuer(provider *jwt.Provider, ttl time.Duration) *JWTAccessTokenIssuer {
	return &JWTAccessTokenIssuer{provider: provider, ttl: ttl, now: time.Now}
}

// Issue grava o id, o papel e o e-mail confirmado no token (REQ-04, REQ-15).
//
// gofi:implements REQ-04
func (i *JWTAccessTokenIssuer) Issue(userID string, role entity.Role, emailConfirmed bool) (string, time.Time, error) {
	issuedAt := i.now()
	expiresAt := issuedAt.Add(i.ttl)
	token, err := i.provider.IssueAccessToken(types.Claims{
		UserID:    userID,
		Roles:     []string{string(role)},
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Extra:     map[string]any{ClaimEmailConfirmed: emailConfirmed},
	})
	return token, expiresAt, err
}
