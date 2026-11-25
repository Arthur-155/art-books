//gofi:context contas
package handler

import (
	"net/http"

	"github.com/joaoprofile/gofi-sdk-go/base/errs"
	"github.com/joaoprofile/gofi-sdk-go/netx/httpx"

	"github.com/Arthur-155/art-books/domain/contas/service"
)

const passwordResetMessage = "Se o e-mail estiver cadastrado, enviamos um link de recuperação. Se ele não chegar, peça de novo."

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type passwordResetRequest struct {
	Email string `json:"email"`
}

type passwordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type sessionResponse struct {
	User service.SessionUser `json:"user"`
}

// AuthHandler expõe login, sessão e senha.
type AuthHandler struct {
	svc service.ContasService
}

func NewAuthHandler(svc service.ContasService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Handlers implements httpx.RouterHandler.
func (h *AuthHandler) Handlers() []*httpx.Route {
	public := httpx.PublicRoutes(authPrefix,
		httpx.POST("/login").To(h.login),
		httpx.POST("/refresh").To(h.refresh),
		httpx.POST("/logout").To(h.logout),
		httpx.POST("/password-reset/request").To(h.requestPasswordReset),
		httpx.POST("/password-reset/confirm").To(h.resetPassword),
	)
	private := httpx.PrivateRoutes(accountPrefix, httpx.PUT("/password").To(h.changePassword))
	return append(public, private...)
}

// gofi:implements REQ-05
func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeBody(w, r, &req) {
		return
	}
	session, appErr := h.svc.Login(r.Context(), service.LoginInput{
		Email: req.Email, Password: req.Password, IP: httpx.ClientIP(r),
	})
	respondSession(w, r, session, appErr)
}

// gofi:implements REQ-07
func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	session, appErr := h.svc.RefreshSession(r.Context(), refreshCredential(r))
	respondSession(w, r, session, appErr)
}

func respondSession(w http.ResponseWriter, r *http.Request, session *service.SessionResult, appErr errs.AppError) {
	if appErr.Exists() {
		httpx.RespondError(w, r, appErr)
		return
	}
	setSessionCookies(w, session)
	httpx.Response(w, http.StatusOK, sessionResponse{User: session.User})
}

// gofi:implements REQ-08
func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	appErr := h.svc.Logout(r.Context(), refreshCredential(r))
	clearSessionCookies(w)
	finish(w, r, appErr, http.StatusNoContent)
}

// gofi:implements REQ-09
func (h *AuthHandler) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequest
	if !decodeBody(w, r, &req) {
		return
	}
	appErr := h.svc.RequestPasswordReset(r.Context(), req.Email, httpx.ClientIP(r))
	if appErr.Exists() {
		httpx.RespondError(w, r, appErr)
		return
	}
	httpx.Response(w, http.StatusAccepted, map[string]string{"message": passwordResetMessage})
}

// gofi:implements REQ-09
func (h *AuthHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req passwordResetConfirmRequest
	if !decodeBody(w, r, &req) {
		return
	}
	appErr := h.svc.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if !appErr.Exists() {
		clearSessionCookies(w)
	}
	finish(w, r, appErr, http.StatusNoContent)
}

// gofi:implements REQ-10
func (h *AuthHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	var req changePasswordRequest
	if !decodeBody(w, r, &req) {
		return
	}
	appErr := h.svc.ChangePassword(r.Context(), currentUserID(r), req.CurrentPassword, req.NewPassword)
	if !appErr.Exists() {
		clearSessionCookies(w)
	}
	finish(w, r, appErr, http.StatusNoContent)
}
