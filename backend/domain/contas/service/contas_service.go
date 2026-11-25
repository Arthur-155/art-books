//gofi:context contas
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/joaoprofile/gofi-sdk-go/base/errs"
	"github.com/joaoprofile/gofi-sdk-go/obs/logging"

	"github.com/Arthur-155/art-books/domain/contas/entity"
	"github.com/Arthur-155/art-books/domain/contas/repository"
)

const (
	linkTokenBytes = 32

	confirmationTTL   = 24 * time.Hour
	passwordResetTTL  = 30 * time.Minute
	emailApprovalTTL  = time.Hour

	confirmationPath      = "/confirmar-email"
	passwordResetPath     = "/redefinir-senha"
	emailApprovalPath     = "/aprovar-troca-de-email"
	emailConfirmationPath = "/confirmar-novo-email"
)

// Mailer envia um e-mail transacional; a única implementação é o cliente do Resend.
type Mailer interface {
	Send(ctx context.Context, mail entity.Mail) error
}

// AccessTokenIssuer emite o token de acesso (JWT) da sessão.
type AccessTokenIssuer interface {
	Issue(userID string, role entity.Role, emailConfirmed bool) (token string, expiresAt time.Time, err error)
}

type Options struct {
	BaseURL         string
	RefreshTokenTTL time.Duration
	Now             func() time.Time
}

type SignUpInput struct {
	Email         string
	Password      string
	DisplayName   string
	Username      string
	TermsAccepted bool
	AgeDeclared   bool
}

type SignUpResult struct {
	ConfirmationEmailSent bool `json:"confirmationEmailSent"`
}

type LoginInput struct {
	Email    string
	Password string
	IP       string
}

type SessionUser struct {
	ID             string      `json:"id"`
	Username       string      `json:"username"`
	DisplayName    string      `json:"displayName"`
	Role           entity.Role `json:"role"`
	EmailConfirmed bool        `json:"emailConfirmed"`
}

// SessionResult leva os tokens para o handler pôr em cookie; eles nunca vão no corpo.
type SessionResult struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
	User             SessionUser
}

type Account struct {
	ID             string      `json:"id"`
	Email          string      `json:"email"`
	Username       string      `json:"username"`
	DisplayName    string      `json:"displayName"`
	Role           entity.Role `json:"role"`
	EmailConfirmed bool        `json:"emailConfirmed"`
	PendingEmail   *string     `json:"pendingEmail"`
}

type ContasService interface {
	SignUp(ctx context.Context, in SignUpInput) (SignUpResult, errs.AppError)
	ConfirmEmail(ctx context.Context, token string) errs.AppError
	ResendConfirmation(ctx context.Context, userID, ip string) errs.AppError
	Login(ctx context.Context, in LoginInput) (*SessionResult, errs.AppError)
	RefreshSession(ctx context.Context, refreshToken string) (*SessionResult, errs.AppError)
	Logout(ctx context.Context, refreshToken string) errs.AppError
	RequestPasswordReset(ctx context.Context, email, ip string) errs.AppError
	ResetPassword(ctx context.Context, token, newPassword string) errs.AppError
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) errs.AppError
	RequestEmailChange(ctx context.Context, userID, newEmail, ip string) errs.AppError
	ApproveEmailChange(ctx context.Context, token string) errs.AppError
	ConfirmEmailChange(ctx context.Context, token string) errs.AppError
	DeleteAccount(ctx context.Context, userID, password string) errs.AppError
	GetAccount(ctx context.Context, userID string) (*Account, errs.AppError)
	GrantAdmin(ctx context.Context, email string) errs.AppError
	RevokeAdmin(ctx context.Context, email string) errs.AppError
	PurgeExpired(ctx context.Context) errs.AppError
}

type contasService struct {
	repo       repository.ContasRepository
	mailer     Mailer
	tokens     AccessTokenIssuer
	baseURL    string
	refreshTTL time.Duration
	now        func() time.Time
}

func NewContasService(repo repository.ContasRepository, mailer Mailer, tokens AccessTokenIssuer, opts Options) ContasService {
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	return &contasService{
		repo: repo, mailer: mailer, tokens: tokens,
		baseURL: opts.BaseURL, refreshTTL: opts.RefreshTokenTTL, now: now,
	}
}

// SignUp cria a conta e tenta enviar a confirmação; a falha no envio não desfaz a conta (REQ-14.1).
//
// gofi:implements REQ-01
func (s *contasService) SignUp(ctx context.Context, in SignUpInput) (SignUpResult, errs.AppError) {
	if !isValidSignUp(in) {
		return SignUpResult{}, ErrContasInvalidInput.New()
	}
	user, appErr := s.newUser(in)
	if appErr.Exists() {
		return SignUpResult{}, appErr
	}
	if appErr := s.saveUser(ctx, user); appErr.Exists() {
		return SignUpResult{}, appErr
	}
	delivery := s.sendConfirmation(ctx, user)
	return SignUpResult{ConfirmationEmailSent: !delivery.Exists()}, errs.AppError{}
}

func isValidSignUp(in SignUpInput) bool {
	return in.TermsAccepted && in.AgeDeclared &&
		entity.IsValidEmail(in.Email) &&
		entity.IsValidPassword(in.Password) &&
		entity.IsValidUsername(in.Username) &&
		entity.IsValidDisplayName(in.DisplayName)
}

func (s *contasService) newUser(in SignUpInput) (*entity.User, errs.AppError) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	hash, err := hashPassword(in.Password)
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	now := s.now()
	return &entity.User{
		ID: id.String(), Email: entity.NormalizeEmail(in.Email),
		Username: entity.NormalizeUsername(in.Username), DisplayName: trimmed(in.DisplayName),
		PasswordHash: hash, Role: entity.RoleUser, TermsAcceptedAt: now, AgeDeclaredAt: now,
	}, errs.AppError{}
}

func (s *contasService) saveUser(ctx context.Context, user *entity.User) errs.AppError {
	err := s.repo.SaveUser(ctx, user)
	switch {
	case errors.Is(err, repository.ErrEmailTaken):
		return ErrContasEmailTaken.New()
	case errors.Is(err, repository.ErrUsernameTaken):
		return ErrContasUsernameTaken.New()
	case err != nil:
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// ConfirmEmail marca o e-mail como confirmado; o link é de uso único e vence em 24 h.
//
// gofi:implements REQ-02
func (s *contasService) ConfirmEmail(ctx context.Context, token string) errs.AppError {
	consumed, appErr := s.consumeLink(ctx, entity.PurposeEmailConfirmation, token)
	if appErr.Exists() {
		return appErr
	}
	if err := s.repo.MarkEmailConfirmed(ctx, consumed.UserID); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// ResendConfirmation emite um link novo, o que invalida o anterior.
//
// gofi:implements REQ-03
func (s *contasService) ResendConfirmation(ctx context.Context, userID, ip string) errs.AppError {
	user, appErr := s.findUser(ctx, userID)
	if appErr.Exists() {
		return appErr
	}
	if user.IsEmailConfirmed() {
		return ErrContasEmailAlreadyConfirmed.New()
	}
	if appErr := s.acceptEmailRequest(ctx, entity.RequestConfirmationResend, user.Email, ip); appErr.Exists() {
		return appErr
	}
	return s.sendConfirmation(ctx, user)
}

// RequestEmailChange envia a aprovação ao e-mail atual; o novo só recebe algo depois da aprovação.
//
// gofi:implements REQ-11
func (s *contasService) RequestEmailChange(ctx context.Context, userID, newEmail, ip string) errs.AppError {
	if !entity.IsValidEmail(newEmail) {
		return ErrContasInvalidInput.New()
	}
	user, appErr := s.findUser(ctx, userID)
	if appErr.Exists() {
		return appErr
	}
	normalized := entity.NormalizeEmail(newEmail)
	if appErr := s.ensureEmailFree(ctx, normalized); appErr.Exists() {
		return appErr
	}
	if appErr := s.acceptEmailRequest(ctx, entity.RequestEmailChange, user.Email, ip); appErr.Exists() {
		return appErr
	}
	return s.sendEmailChangeApproval(ctx, user, normalized)
}

func (s *contasService) sendEmailChangeApproval(ctx context.Context, user *entity.User, newEmail string) errs.AppError {
	purposes := []entity.TokenPurpose{entity.PurposeEmailChangeApproval, entity.PurposeEmailChangeConfirmation}
	if err := s.repo.DeleteEmailTokens(ctx, user.ID, purposes...); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return s.issueAndSend(ctx, linkSpec{
		purpose: entity.PurposeEmailChangeApproval, kind: entity.MailEmailChangeApproval, path: emailApprovalPath,
		ttl: emailApprovalTTL, userID: user.ID, to: user.Email, newEmail: &newEmail,
	})
}

// ApproveEmailChange troca a aprovação pela confirmação enviada ao novo e-mail.
//
// gofi:implements REQ-11
func (s *contasService) ApproveEmailChange(ctx context.Context, token string) errs.AppError {
	var refusal errs.AppError
	err := s.repo.ApproveEmailChange(ctx, hashToken(token), func(ctx context.Context, approval *entity.EmailToken) error {
		refusal = s.sendEmailChangeConfirmation(ctx, approval)
		if refusal.Exists() {
			return errApprovalRefused
		}
		return nil
	})
	switch {
	case errors.Is(err, errApprovalRefused):
		return refusal
	case errors.Is(err, repository.ErrInvalidLink):
		return ErrContasInvalidLink.New()
	case err != nil:
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// errApprovalRefused desfaz a transação da aprovação, devolvendo o link, quando a troca não segue (REQ-14.4).
var errApprovalRefused = errors.New("email change approval refused")

func (s *contasService) sendEmailChangeConfirmation(ctx context.Context, approval *entity.EmailToken) errs.AppError {
	if appErr := s.ensureEmailFree(ctx, *approval.NewEmail); appErr.Exists() {
		return appErr
	}
	return s.issueAndSend(ctx, linkSpec{
		purpose: entity.PurposeEmailChangeConfirmation, kind: entity.MailEmailChangeConfirmation,
		path: emailConfirmationPath, ttl: confirmationTTL, userID: approval.UserID,
		to: *approval.NewEmail, newEmail: approval.NewEmail,
	})
}

// ConfirmEmailChange conclui a troca de e-mail e encerra todas as sessões.
//
// gofi:implements REQ-11
func (s *contasService) ConfirmEmailChange(ctx context.Context, token string) errs.AppError {
	err := s.repo.ConfirmEmailChange(ctx, hashToken(token))
	switch {
	case errors.Is(err, repository.ErrInvalidLink):
		return ErrContasInvalidLink.New()
	case errors.Is(err, repository.ErrEmailTaken):
		return ErrContasEmailTaken.New()
	case err != nil:
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// DeleteAccount exclui a conta na hora; as FKs ON DELETE CASCADE levam sessões e tokens.
//
// gofi:implements REQ-12
func (s *contasService) DeleteAccount(ctx context.Context, userID, password string) errs.AppError {
	user, appErr := s.findUser(ctx, userID)
	if appErr.Exists() {
		return appErr
	}
	if !verifyPassword(user.PasswordHash, password) {
		return ErrContasWrongPassword.New()
	}
	if err := s.repo.DeleteUser(ctx, user.ID); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// GetAccount devolve só a conta do próprio usuário.
//
// gofi:implements REQ-16
func (s *contasService) GetAccount(ctx context.Context, userID string) (*Account, errs.AppError) {
	user, appErr := s.findUser(ctx, userID)
	if appErr.Exists() {
		return nil, appErr
	}
	pending, err := s.repo.FindPendingEmailChange(ctx, user.ID)
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	return &Account{
		ID: user.ID, Email: user.Email, Username: user.Username, DisplayName: user.DisplayName,
		Role: user.Role, EmailConfirmed: user.IsEmailConfirmed(), PendingEmail: pending,
	}, errs.AppError{}
}

// GrantAdmin é chamado só pelo comando operacional.
//
// gofi:implements REQ-15
func (s *contasService) GrantAdmin(ctx context.Context, email string) errs.AppError {
	return s.setRole(ctx, email, entity.RoleAdmin)
}

// gofi:implements REQ-15
func (s *contasService) RevokeAdmin(ctx context.Context, email string) errs.AppError {
	return s.setRole(ctx, email, entity.RoleUser)
}

func (s *contasService) setRole(ctx context.Context, email string, role entity.Role) errs.AppError {
	err := s.repo.SetRoleByEmail(ctx, entity.NormalizeEmail(email), role)
	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrContasUserNotFound.New()
	}
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

func (s *contasService) PurgeExpired(ctx context.Context) errs.AppError {
	if err := s.repo.PurgeExpired(ctx, s.now()); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

func (s *contasService) findUser(ctx context.Context, userID string) (*entity.User, errs.AppError) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	if user == nil {
		return nil, ErrContasUnauthenticated.New()
	}
	return user, errs.AppError{}
}

func (s *contasService) ensureEmailFree(ctx context.Context, email string) errs.AppError {
	taken, err := s.repo.IsEmailTaken(ctx, email)
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	if taken {
		return ErrContasEmailTaken.New()
	}
	return errs.AppError{}
}

// consumeLink usa o link uma única vez; inexistente, usado ou vencido são indistinguíveis.
func (s *contasService) consumeLink(ctx context.Context, purpose entity.TokenPurpose, token string) (*entity.EmailToken, errs.AppError) {
	consumed, err := s.repo.ConsumeEmailToken(ctx, purpose, hashToken(token))
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	if consumed == nil {
		return nil, ErrContasInvalidLink.New()
	}
	return consumed, errs.AppError{}
}

func (s *contasService) sendConfirmation(ctx context.Context, user *entity.User) errs.AppError {
	return s.issueAndSend(ctx, linkSpec{
		purpose: entity.PurposeEmailConfirmation, kind: entity.MailEmailConfirmation, path: confirmationPath,
		ttl: confirmationTTL, userID: user.ID, to: user.Email,
	})
}

type linkSpec struct {
	purpose  entity.TokenPurpose
	kind     entity.MailKind
	path     string
	ttl      time.Duration
	userID   string
	to       string
	newEmail *string
}

// issueAndSend grava o hash do token, que invalida o link anterior, e envia o link; a falha do envio mantém o token.
func (s *contasService) issueAndSend(ctx context.Context, spec linkSpec) errs.AppError {
	raw, err := newRandomToken()
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	token := &entity.EmailToken{
		UserID: spec.userID, Purpose: spec.purpose, TokenHash: hashToken(raw),
		NewEmail: spec.newEmail, ExpiresAt: s.now().Add(spec.ttl),
	}
	if err := s.repo.ReplaceEmailToken(ctx, token); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	mail := entity.Mail{Kind: spec.kind, To: spec.to, Link: s.baseURL + spec.path + "?token=" + raw, UserID: spec.userID}
	if err := s.mailer.Send(ctx, mail); err != nil {
		logging.FromContext(ctx).Error("contas: email delivery failed",
			slog.String("kind", string(spec.kind)), slog.String("userId", spec.userID))
		return ErrContasEmailDeliveryFailed.Wrap(err)
	}
	return errs.AppError{}
}

// acceptEmailRequest aplica o limite de 1 por minuto e 5 por hora, por e-mail e por IP, e registra o pedido de forma atômica.
//
// gofi:implements REQ-13
func (s *contasService) acceptEmailRequest(ctx context.Context, requestType entity.RequestType, email, ip string) errs.AppError {
	err := s.repo.AcceptEmailRequest(ctx, requestType, loginKeys(email, ip))
	if errors.Is(err, repository.ErrTooManyRequest) {
		return ErrContasTooManyRequests.New()
	}
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

func newRandomToken() (string, error) {
	buffer := make([]byte, linkTokenBytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
