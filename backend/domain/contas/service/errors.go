//gofi:context contas
package service

import "github.com/joaoprofile/gofi-sdk-go/base/errs"

var (
	ErrContasInvalidInput          = errs.RegisterValidation("CONTAS_INVALID_INPUT", "invalid input")
	ErrContasEmailTaken            = errs.RegisterConflict("CONTAS_EMAIL_TAKEN", "email already in use")
	ErrContasUsernameTaken         = errs.RegisterConflict("CONTAS_USERNAME_TAKEN", "username already in use")
	ErrContasInvalidLink           = errs.RegisterValidation("CONTAS_INVALID_LINK", "invalid or expired link")
	ErrContasInvalidCredentials    = errs.RegisterUnauthorized("CONTAS_INVALID_CREDENTIALS", "invalid email or password")
	ErrContasLoginLocked           = errs.RegisterTooManyRequests("CONTAS_LOGIN_LOCKED", "too many attempts, try again later")
	ErrContasInvalidSession        = errs.RegisterUnauthorized("CONTAS_INVALID_SESSION", "invalid session")
	ErrContasEmailAlreadyConfirmed = errs.RegisterConflict("CONTAS_EMAIL_ALREADY_CONFIRMED", "email already confirmed")
	ErrContasWrongPassword         = errs.RegisterUnprocessable("CONTAS_WRONG_PASSWORD", "wrong password")
	ErrContasTooManyRequests       = errs.RegisterTooManyRequests("CONTAS_TOO_MANY_REQUESTS", "too many requests, try again later")
	ErrContasEmailDeliveryFailed   = errs.RegisterExternalError("CONTAS_EMAIL_DELIVERY_FAILED", "email could not be sent")
	ErrContasUserNotFound          = errs.RegisterNotFound("CONTAS_USER_NOT_FOUND", "account not found")
	ErrContasUnauthenticated       = errs.RegisterUnauthorized("CONTAS_UNAUTHENTICATED", "not authenticated")
	ErrContasOperationFailed       = errs.RegisterOperation("CONTAS_OPERATION_FAILED", "operation failed")
)
