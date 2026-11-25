//gofi:context contas
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/joaoprofile/gofi-sdk-go/base/environment"
	"github.com/joaoprofile/gofi-sdk-go/base/errs"
	"github.com/joaoprofile/gofi-sdk-go/iam/types"
	"github.com/joaoprofile/gofi-sdk-go/netx/httpx"

	"github.com/Arthur-155/art-books/domain/contas/service"
)

const (
	accessCookie      = "access_token"
	refreshCookie     = "refresh_token"
	refreshCookiePath = "/api/v1/auth"
	bearerPrefix      = "Bearer "
)

var errUnauthenticated = errors.New("not authenticated")

type claimsKey struct{}

// TokenParser valida a assinatura e o prazo do token de acesso, sem consultar o banco.
type TokenParser interface {
	ParseToken(token string) (*types.Claims, error)
}

// Middleware guarda as rotas privadas; registre-o com UseAuth antes de Handlers.
//
// gofi:implements REQ-16
func Middleware(parser TokenParser) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := parser.ParseToken(accessToken(r))
			if err != nil || claims == nil {
				httpx.Error(w, http.StatusUnauthorized, errUnauthenticated)
				return
			}
			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

func accessToken(r *http.Request) string {
	if token, ok := strings.CutPrefix(r.Header.Get("Authorization"), bearerPrefix); ok {
		return token
	}
	if cookie, err := r.Cookie(accessCookie); err == nil {
		return cookie.Value
	}
	return ""
}

func WithClaims(ctx context.Context, claims *types.Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// ClaimsFromContext devolve nil fora de uma rota privada.
func ClaimsFromContext(ctx context.Context) *types.Claims {
	claims, _ := ctx.Value(claimsKey{}).(*types.Claims)
	return claims
}

func cookieSecure() bool {
	return !environment.IsLocalEnvironment()
}

func setSessionCookies(w http.ResponseWriter, session *service.SessionResult) {
	http.SetCookie(w, newCookie(accessCookie, session.AccessToken, "/", session.AccessExpiresAt))
	http.SetCookie(w, newCookie(refreshCookie, session.RefreshToken, refreshCookiePath, session.RefreshExpiresAt))
}

func clearSessionCookies(w http.ResponseWriter) {
	for name, path := range map[string]string{accessCookie: "/", refreshCookie: refreshCookiePath} {
		cookie := newCookie(name, "", path, time.Time{})
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
}

func newCookie(name, value, path string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name: name, Value: value, Path: path, Expires: expiresAt,
		HttpOnly: true, Secure: cookieSecure(), SameSite: http.SameSiteLaxMode,
	}
}

func refreshCredential(r *http.Request) string {
	cookie, err := r.Cookie(refreshCookie)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func decodeBody(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := httpx.ParseRequestBody(w, r, target, httpx.AllowUnknownFields()); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

// finish responde só com o status quando não há erro.
func finish(w http.ResponseWriter, r *http.Request, appErr errs.AppError, status int) {
	if appErr.Exists() {
		httpx.RespondError(w, r, appErr)
		return
	}
	httpx.JSON(w, status, nil)
}

func currentUserID(r *http.Request) string {
	return ClaimsFromContext(r.Context()).UserID
}
