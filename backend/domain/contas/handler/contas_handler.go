//gofi:context contas
package handler

import (
	"net/http"

	"github.com/joaoprofile/gofi-sdk-go/netx/httpx"

	"github.com/Arthur-155/art-books/domain/contas/service"
)

const (
	authPrefix    = "/api/v1/auth"
	accountPrefix = "/api/v1/account"
)

type signUpRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	DisplayName   string `json:"displayName"`
	Username      string `json:"username"`
	TermsAccepted bool   `json:"termsAccepted"`
	AgeDeclared   bool   `json:"ageDeclared"`
}

type tokenRequest struct {
	Token string `json:"token"`
}

type newEmailRequest struct {
	NewEmail string `json:"newEmail"`
}

type passwordRequest struct {
	Password string `json:"password"`
}

// ContasHandler expõe cadastro, e-mail e consulta da conta.
type ContasHandler struct {
	svc service.ContasService
}

func NewContasHandler(svc service.ContasService) *ContasHandler {
	return &ContasHandler{svc: svc}
}

// Handlers implements httpx.RouterHandler.
func (h *ContasHandler) Handlers() []*httpx.Route {
	public := httpx.PublicRoutes(authPrefix,
		httpx.POST("/signup").To(h.signUp),
		httpx.POST("/email-confirmation/confirm").To(h.confirmEmail),
		httpx.POST("/email-change/approve").To(h.approveEmailChange),
		httpx.POST("/email-change/confirm").To(h.confirmEmailChange),
	)
	private := httpx.PrivateRoutes(accountPrefix,
		httpx.GET("/").To(h.getAccount),
		httpx.POST("/email-confirmation/resend").To(h.resendConfirmation),
		httpx.POST("/email-change").To(h.requestEmailChange),
		httpx.DELETE("/").To(h.deleteAccount),
	)
	return append(public, private...)
}

// gofi:implements REQ-01
func (h *ContasHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var req signUpRequest
	if !decodeBody(w, r, &req) {
		return
	}
	result, appErr := h.svc.SignUp(r.Context(), service.SignUpInput{
		Email: req.Email, Password: req.Password, DisplayName: req.DisplayName, Username: req.Username,
		TermsAccepted: req.TermsAccepted, AgeDeclared: req.AgeDeclared,
	})
	if appErr.Exists() {
		httpx.RespondError(w, r, appErr)
		return
	}
	httpx.Response(w, http.StatusCreated, result)
}

// gofi:implements REQ-02
func (h *ContasHandler) confirmEmail(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if !decodeBody(w, r, &req) {
		return
	}
	finish(w, r, h.svc.ConfirmEmail(r.Context(), req.Token), http.StatusNoContent)
}

// gofi:implements REQ-11
func (h *ContasHandler) approveEmailChange(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if !decodeBody(w, r, &req) {
		return
	}
	finish(w, r, h.svc.ApproveEmailChange(r.Context(), req.Token), http.StatusAccepted)
}

// gofi:implements REQ-11
func (h *ContasHandler) confirmEmailChange(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if !decodeBody(w, r, &req) {
		return
	}
	appErr := h.svc.ConfirmEmailChange(r.Context(), req.Token)
	if !appErr.Exists() {
		clearSessionCookies(w)
	}
	finish(w, r, appErr, http.StatusNoContent)
}

// gofi:implements REQ-16
func (h *ContasHandler) getAccount(w http.ResponseWriter, r *http.Request) {
	account, appErr := h.svc.GetAccount(r.Context(), currentUserID(r))
	if appErr.Exists() {
		httpx.RespondError(w, r, appErr)
		return
	}
	httpx.Response(w, http.StatusOK, account)
}

// gofi:implements REQ-03
func (h *ContasHandler) resendConfirmation(w http.ResponseWriter, r *http.Request) {
	appErr := h.svc.ResendConfirmation(r.Context(), currentUserID(r), httpx.ClientIP(r))
	finish(w, r, appErr, http.StatusAccepted)
}

// gofi:implements REQ-11
func (h *ContasHandler) requestEmailChange(w http.ResponseWriter, r *http.Request) {
	var req newEmailRequest
	if !decodeBody(w, r, &req) {
		return
	}
	appErr := h.svc.RequestEmailChange(r.Context(), currentUserID(r), req.NewEmail, httpx.ClientIP(r))
	finish(w, r, appErr, http.StatusAccepted)
}

// gofi:implements REQ-12
func (h *ContasHandler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	var req passwordRequest
	if !decodeBody(w, r, &req) {
		return
	}
	appErr := h.svc.DeleteAccount(r.Context(), currentUserID(r), req.Password)
	if !appErr.Exists() {
		clearSessionCookies(w)
	}
	finish(w, r, appErr, http.StatusNoContent)
}
