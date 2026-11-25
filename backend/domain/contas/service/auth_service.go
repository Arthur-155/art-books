//gofi:context contas
package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/joaoprofile/gofi-sdk-go/base/errs"
	"github.com/joaoprofile/gofi-sdk-go/iam/provider/password"
	"github.com/joaoprofile/gofi-sdk-go/obs/logging"

	"github.com/Arthur-155/art-books/domain/contas/entity"
	"github.com/Arthur-155/art-books/domain/contas/repository"
)

// dummyPasswordHash iguala o custo do login com e-mail inexistente ao de senha errada.
var dummyPasswordHash = sync.OnceValue(func() string {
	hash, err := password.Hash("dummy-password-for-timing")
	if err != nil {
		return ""
	}
	return hash
})

// Login abre a sessão; e-mail inexistente e senha errada respondem igual.
//
// gofi:implements REQ-05
func (s *contasService) Login(ctx context.Context, in LoginInput) (*SessionResult, errs.AppError) {
	if in.Email == "" || in.Password == "" {
		return nil, ErrContasInvalidInput.New()
	}
	email := entity.NormalizeEmail(in.Email)
	keys := loginKeys(email, in.IP)
	if appErr := s.ensureNotLocked(ctx, keys); appErr.Exists() {
		return nil, appErr
	}
	user, appErr := s.authenticate(ctx, email, in.Password, keys)
	if appErr.Exists() {
		return nil, appErr
	}
	if err := s.repo.ClearLoginFailures(ctx, keys[0]); err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	return s.issueSession(ctx, user)
}

func loginKeys(email, ip string) []entity.LimitKey {
	keys := []entity.LimitKey{{Kind: entity.KeyKindEmail, Value: email}}
	if ip != "" {
		keys = append(keys, entity.LimitKey{Kind: entity.KeyKindIP, Value: ip})
	}
	return keys
}

// ensureNotLocked recusa antes de conferir a senha; a recusa não conta como erro.
//
// gofi:implements REQ-06
func (s *contasService) ensureNotLocked(ctx context.Context, keys []entity.LimitKey) errs.AppError {
	locked, err := s.repo.IsLoginLocked(ctx, keys)
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	if locked {
		return ErrContasLoginLocked.New()
	}
	return errs.AppError{}
}

func (s *contasService) authenticate(ctx context.Context, email, plain string, keys []entity.LimitKey) (*entity.User, errs.AppError) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	if user != nil && verifyPassword(user.PasswordHash, plain) {
		return user, errs.AppError{}
	}
	if user == nil {
		verifyPassword(dummyPasswordHash(), plain)
	}
	for _, key := range keys {
		if err := s.repo.RecordLoginFailure(ctx, key); err != nil {
			return nil, ErrContasOperationFailed.Wrap(err)
		}
	}
	return nil, ErrContasInvalidCredentials.New()
}

// issueSession abre uma sessão nova com credencial de renovação opaca, guardada só como hash.
func (s *contasService) issueSession(ctx context.Context, user *entity.User) (*SessionResult, errs.AppError) {
	refresh, err := newRandomToken()
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	session := &entity.Session{UserID: user.ID, TokenHash: hashToken(refresh), ExpiresAt: s.now().Add(s.refreshTTL)}
	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	return s.buildSession(user, refresh, session.ExpiresAt)
}

// buildSession traz o papel e o e-mail confirmado atuais no token de acesso.
//
// gofi:implements REQ-04
func (s *contasService) buildSession(user *entity.User, refresh string, refreshExpiresAt time.Time) (*SessionResult, errs.AppError) {
	access, accessExpiresAt, err := s.tokens.Issue(user.ID, user.Role, user.IsEmailConfirmed())
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	return &SessionResult{
		AccessToken: access, AccessExpiresAt: accessExpiresAt,
		RefreshToken: refresh, RefreshExpiresAt: refreshExpiresAt,
		User: SessionUser{
			ID: user.ID, Username: user.Username, DisplayName: user.DisplayName,
			Role: user.Role, EmailConfirmed: user.IsEmailConfirmed(),
		},
	}, errs.AppError{}
}

// RefreshSession troca a credencial de renovação por outra e relê o usuário.
//
// gofi:implements REQ-07
func (s *contasService) RefreshSession(ctx context.Context, refreshToken string) (*SessionResult, errs.AppError) {
	if refreshToken == "" {
		return nil, ErrContasInvalidSession.New()
	}
	next, err := newRandomToken()
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	fresh := &entity.Session{TokenHash: hashToken(next), ExpiresAt: s.now().Add(s.refreshTTL)}
	rotated, err := s.repo.RotateSession(ctx, hashToken(refreshToken), fresh)
	if errors.Is(err, repository.ErrInvalidSession) {
		return nil, ErrContasInvalidSession.New()
	}
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	user, err := s.repo.FindUserByID(ctx, rotated.UserID)
	if err != nil {
		return nil, ErrContasOperationFailed.Wrap(err)
	}
	if user == nil {
		return nil, ErrContasInvalidSession.New()
	}
	return s.buildSession(user, next, rotated.ExpiresAt)
}

// Logout encerra só a sessão do dispositivo; credencial inexistente não é erro.
//
// gofi:implements REQ-08
func (s *contasService) Logout(ctx context.Context, refreshToken string) errs.AppError {
	if refreshToken == "" {
		return errs.AppError{}
	}
	if err := s.repo.DeleteSession(ctx, hashToken(refreshToken)); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// RequestPasswordReset responde igual exista ou não a conta e mesmo se o envio falhar.
//
// gofi:implements REQ-09
func (s *contasService) RequestPasswordReset(ctx context.Context, email, ip string) errs.AppError {
	if !entity.IsValidEmail(email) {
		return ErrContasInvalidInput.New()
	}
	normalized := entity.NormalizeEmail(email)
	if appErr := s.acceptEmailRequest(ctx, entity.RequestPasswordReset, normalized, ip); appErr.Exists() {
		return appErr
	}
	user, err := s.repo.FindUserByEmail(ctx, normalized)
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	if user != nil {
		s.sendPasswordReset(ctx, user)
	}
	return errs.AppError{}
}

// sendPasswordReset envia na própria requisição: o envio só roda para conta existente, e a espera é aceita.
func (s *contasService) sendPasswordReset(ctx context.Context, user *entity.User) {
	delivery := s.issueAndSend(ctx, linkSpec{
		purpose: entity.PurposePasswordReset, kind: entity.MailPasswordReset, path: passwordResetPath,
		ttl: passwordResetTTL, userID: user.ID, to: user.Email,
	})
	if delivery.Exists() {
		logging.FromContext(ctx).Warn("contas: password reset email not sent", slog.String("userId", user.ID))
	}
}

// ResetPassword troca a senha pelo link e encerra todas as sessões; senha inválida não gasta o link.
//
// gofi:implements REQ-09
func (s *contasService) ResetPassword(ctx context.Context, token, newPassword string) errs.AppError {
	if !entity.IsValidPassword(newPassword) {
		return ErrContasInvalidInput.New()
	}
	hash, err := hashPassword(newPassword)
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	err = s.repo.ResetPassword(ctx, hashToken(token), hash)
	if errors.Is(err, repository.ErrInvalidLink) {
		return ErrContasInvalidLink.New()
	}
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// ChangePassword exige a senha atual e encerra todas as sessões, a atual incluída.
//
// gofi:implements REQ-10
func (s *contasService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) errs.AppError {
	if !entity.IsValidPassword(newPassword) {
		return ErrContasInvalidInput.New()
	}
	user, appErr := s.findUser(ctx, userID)
	if appErr.Exists() {
		return appErr
	}
	if !verifyPassword(user.PasswordHash, currentPassword) {
		return ErrContasWrongPassword.New()
	}
	hash, err := hashPassword(newPassword)
	if err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	if err := s.repo.ChangePassword(ctx, user.ID, hash); err != nil {
		return ErrContasOperationFailed.Wrap(err)
	}
	return errs.AppError{}
}

// hashPassword guarda a senha inteira em argon2id, formato PHC.
//
// gofi:implements REQ-17
func hashPassword(plain string) (string, error) {
	return password.Hash(plain)
}

func verifyPassword(hash, plain string) bool {
	return password.Verify(hash, plain) == nil
}

// hashToken resume com SHA-256 os tokens de link e de renovação.
//
// gofi:implements REQ-17
func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

func trimmed(value string) string {
	return strings.TrimSpace(value)
}
