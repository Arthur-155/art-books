//gofi:context contas
package service_test

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joaoprofile/gofi-sdk-go/base/errs"

	"github.com/Arthur-155/art-books/domain/contas/entity"
	"github.com/Arthur-155/art-books/domain/contas/repository"
	"github.com/Arthur-155/art-books/domain/contas/service"
)

const (
	testBaseURL      = "https://app.test"
	testRefreshTTL   = 720 * time.Hour
	testPassword     = "segredo123"
	testEmail        = "ana@exemplo.com"
	testUsername     = "ana_leitora"
	lockWindow       = 15 * time.Minute
	lockFailureLimit = 5
)

type testClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *testClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *testClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

type fakeMailer struct {
	mu    sync.Mutex
	mails []entity.Mail
	fail  bool
}

func (m *fakeMailer) Send(_ context.Context, mail entity.Mail) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fail {
		return errors.New("resend fora do ar")
	}
	m.mails = append(m.mails, mail)
	return nil
}

func (m *fakeMailer) sent() []entity.Mail {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]entity.Mail(nil), m.mails...)
}

func (m *fakeMailer) last() entity.Mail {
	mails := m.sent()
	return mails[len(mails)-1]
}

func (m *fakeMailer) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mails = nil
}

func (m *fakeMailer) setFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fail = fail
}

type issuedToken struct {
	userID         string
	role           entity.Role
	emailConfirmed bool
}

type fakeTokens struct {
	mu     sync.Mutex
	issued []issuedToken
}

func (f *fakeTokens) Issue(userID string, role entity.Role, emailConfirmed bool) (string, time.Time, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.issued = append(f.issued, issuedToken{userID, role, emailConfirmed})
	return "access-" + userID, time.Time{}, nil
}

func (f *fakeTokens) last() issuedToken {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.issued[len(f.issued)-1]
}

type loginLock struct {
	count       int
	windowStart time.Time
	lockedUntil time.Time
}

type requestRecord struct {
	requestType entity.RequestType
	key         entity.LimitKey
	at          time.Time
}

type fakeRepo struct {
	mu       sync.Mutex
	clock    *testClock
	users    map[string]*entity.User
	sessions []*entity.Session
	tokens   []*entity.EmailToken
	locks    map[entity.LimitKey]*loginLock
	requests []requestRecord
}

func newFakeRepo(clock *testClock) *fakeRepo {
	return &fakeRepo{clock: clock, users: map[string]*entity.User{}, locks: map[entity.LimitKey]*loginLock{}}
}

func (r *fakeRepo) SaveUser(_ context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.users {
		if existing.Email == user.Email {
			return repository.ErrEmailTaken
		}
		if existing.Username == user.Username {
			return repository.ErrUsernameTaken
		}
	}
	stored := *user
	r.users[user.ID] = &stored
	return nil
}

func (r *fakeRepo) FindUserByID(_ context.Context, id string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.copyOf(r.users[id]), nil
}

func (r *fakeRepo) FindUserByEmail(_ context.Context, email string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			return r.copyOf(user), nil
		}
	}
	return nil, nil
}

func (r *fakeRepo) copyOf(user *entity.User) *entity.User {
	if user == nil {
		return nil
	}
	clone := *user
	return &clone
}

func (r *fakeRepo) IsEmailTaken(_ context.Context, email string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeRepo) MarkEmailConfirmed(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if user := r.users[userID]; user != nil && user.EmailConfirmedAt == nil {
		now := r.clock.Now()
		user.EmailConfirmedAt = &now
	}
	return nil
}

func (r *fakeRepo) UpdatePasswordHash(_ context.Context, userID, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if user := r.users[userID]; user != nil {
		user.PasswordHash = hash
	}
	return nil
}

func (r *fakeRepo) UpdateEmail(_ context.Context, userID, email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, other := range r.users {
		if id != userID && other.Email == email {
			return repository.ErrEmailTaken
		}
	}
	if user := r.users[userID]; user != nil {
		now := r.clock.Now()
		user.Email = email
		user.EmailConfirmedAt = &now
	}
	return nil
}

func (r *fakeRepo) SetRoleByEmail(_ context.Context, email string, role entity.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			user.Role = role
			return nil
		}
	}
	return repository.ErrUserNotFound
}

// DeleteUser imita o ON DELETE CASCADE das tabelas de sessão e de token.
func (r *fakeRepo) DeleteUser(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, userID)
	r.removeSessions(func(s *entity.Session) bool { return s.UserID == userID })
	r.removeTokens(func(t *entity.EmailToken) bool { return t.UserID == userID })
	return nil
}

func (r *fakeRepo) removeSessions(drop func(*entity.Session) bool) {
	kept := r.sessions[:0]
	for _, session := range r.sessions {
		if !drop(session) {
			kept = append(kept, session)
		}
	}
	r.sessions = kept
}

func (r *fakeRepo) removeTokens(drop func(*entity.EmailToken) bool) {
	kept := r.tokens[:0]
	for _, token := range r.tokens {
		if !drop(token) {
			kept = append(kept, token)
		}
	}
	r.tokens = kept
}

func (r *fakeRepo) SaveSession(_ context.Context, session *entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	stored := *session
	r.sessions = append(r.sessions, &stored)
	return nil
}

func (r *fakeRepo) RotateSession(_ context.Context, oldHash []byte, session *entity.Session) (*entity.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for index, existing := range r.sessions {
		if string(existing.TokenHash) != string(oldHash) || !existing.ExpiresAt.After(r.clock.Now()) {
			continue
		}
		r.sessions = append(r.sessions[:index], r.sessions[index+1:]...)
		rotated := &entity.Session{UserID: existing.UserID, TokenHash: session.TokenHash, ExpiresAt: session.ExpiresAt}
		stored := *rotated
		r.sessions = append(r.sessions, &stored)
		return rotated, nil
	}
	return nil, repository.ErrInvalidSession
}

func (r *fakeRepo) DeleteSession(_ context.Context, tokenHash []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeSessions(func(s *entity.Session) bool { return string(s.TokenHash) == string(tokenHash) })
	return nil
}

func (r *fakeRepo) DeleteUserSessions(_ context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeSessions(func(s *entity.Session) bool { return s.UserID == userID })
	return nil
}

func (r *fakeRepo) ReplaceEmailToken(_ context.Context, token *entity.EmailToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeTokens(func(t *entity.EmailToken) bool { return t.UserID == token.UserID && t.Purpose == token.Purpose })
	stored := *token
	r.tokens = append(r.tokens, &stored)
	return nil
}

func (r *fakeRepo) ConsumeEmailToken(_ context.Context, purpose entity.TokenPurpose, tokenHash []byte) (*entity.EmailToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for index, token := range r.tokens {
		if token.Purpose != purpose || string(token.TokenHash) != string(tokenHash) || !token.ExpiresAt.After(r.clock.Now()) {
			continue
		}
		r.tokens = append(r.tokens[:index], r.tokens[index+1:]...)
		consumed := *token
		return &consumed, nil
	}
	return nil, nil
}

func (r *fakeRepo) DeleteEmailTokens(_ context.Context, userID string, purposes ...entity.TokenPurpose) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeTokens(func(t *entity.EmailToken) bool {
		if t.UserID != userID {
			return false
		}
		for _, purpose := range purposes {
			if t.Purpose == purpose {
				return true
			}
		}
		return false
	})
	return nil
}

func (r *fakeRepo) FindPendingEmailChange(_ context.Context, userID string) (*string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, token := range r.tokens {
		isChange := token.Purpose == entity.PurposeEmailChangeApproval || token.Purpose == entity.PurposeEmailChangeConfirmation
		if token.UserID == userID && isChange {
			return token.NewEmail, nil
		}
	}
	return nil, nil
}

func (r *fakeRepo) IsLoginLocked(_ context.Context, keys []entity.LimitKey) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, key := range keys {
		if lock := r.locks[key]; lock != nil && lock.lockedUntil.After(r.clock.Now()) {
			return true, nil
		}
	}
	return false, nil
}

// RecordLoginFailure segue o desenho: janela de 15 min, bloqueio de 15 min na quinta falha.
func (r *fakeRepo) RecordLoginFailure(_ context.Context, key entity.LimitKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.clock.Now()
	lock := r.locks[key]
	if lock == nil {
		lock = &loginLock{windowStart: now}
		r.locks[key] = lock
	}
	if now.Sub(lock.windowStart) > lockWindow {
		lock.count, lock.windowStart = 0, now
	}
	lock.count++
	if lock.count >= lockFailureLimit {
		lock.lockedUntil = now.Add(lockWindow)
		lock.count = 0
	}
	return nil
}

func (r *fakeRepo) ClearLoginFailures(_ context.Context, key entity.LimitKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.locks, key)
	return nil
}

// Limites de pedidos de e-mail do REQ-13, espelhados do repositório.
const (
	fakeRequestsPerMinute = 1
	fakeRequestsPerHour   = 5
)

// AcceptEmailRequest imita o limite do repositório: confere e registra o pedido sob o mesmo lock.
func (r *fakeRepo) AcceptEmailRequest(_ context.Context, requestType entity.RequestType, keys []entity.LimitKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.clock.Now()
	for _, key := range keys {
		if r.countRequestsLocked(requestType, key, now.Add(-time.Minute)) >= fakeRequestsPerMinute ||
			r.countRequestsLocked(requestType, key, now.Add(-time.Hour)) >= fakeRequestsPerHour {
			return repository.ErrTooManyRequest
		}
	}
	for _, key := range keys {
		r.requests = append(r.requests, requestRecord{requestType, key, now})
	}
	return nil
}

func (r *fakeRepo) countRequestsLocked(requestType entity.RequestType, key entity.LimitKey, since time.Time) int {
	total := 0
	for _, request := range r.requests {
		if request.requestType == requestType && request.key == key && request.at.After(since) {
			total++
		}
	}
	return total
}

// ResetPassword imita a transação: consome o link, troca a senha e encerra sessões e links de recuperação.
func (r *fakeRepo) ResetPassword(_ context.Context, tokenHash []byte, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	token := r.takeLiveTokenLocked(entity.PurposePasswordReset, tokenHash)
	if token == nil {
		return repository.ErrInvalidLink
	}
	r.replacePasswordLocked(token.UserID, passwordHash)
	return nil
}

func (r *fakeRepo) ChangePassword(_ context.Context, userID, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.replacePasswordLocked(userID, passwordHash)
	return nil
}

func (r *fakeRepo) replacePasswordLocked(userID, passwordHash string) {
	if user := r.users[userID]; user != nil {
		user.PasswordHash = passwordHash
	}
	r.removeTokens(func(t *entity.EmailToken) bool {
		return t.UserID == userID && t.Purpose == entity.PurposePasswordReset
	})
	r.removeSessions(func(s *entity.Session) bool { return s.UserID == userID })
}

// ConfirmEmailChange recusa sem gastar o link quando o novo e-mail já tem dono (REQ-11.7).
func (r *fakeRepo) ConfirmEmailChange(_ context.Context, tokenHash []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	token := r.liveTokenLocked(entity.PurposeEmailChangeConfirmation, tokenHash)
	if token == nil {
		return repository.ErrInvalidLink
	}
	userID, newEmail := token.UserID, *token.NewEmail
	if r.isEmailHeldByOtherLocked(userID, newEmail) {
		return repository.ErrEmailTaken
	}
	r.removeTokens(func(t *entity.EmailToken) bool { return t.UserID == userID })
	r.changeEmailLocked(userID, newEmail)
	return nil
}

func (r *fakeRepo) changeEmailLocked(userID, email string) {
	if user := r.users[userID]; user != nil {
		now := r.clock.Now()
		user.Email, user.EmailConfirmedAt = email, &now
	}
	r.removeSessions(func(s *entity.Session) bool { return s.UserID == userID })
}

// ApproveEmailChange imita a transação: se proceed falha, o link de aprovação volta e nada mais muda.
func (r *fakeRepo) ApproveEmailChange(ctx context.Context, tokenHash []byte, proceed func(context.Context, *entity.EmailToken) error) error {
	r.mu.Lock()
	tokensBefore := append([]*entity.EmailToken(nil), r.tokens...)
	approval := r.takeLiveTokenLocked(entity.PurposeEmailChangeApproval, tokenHash)
	r.mu.Unlock()
	if approval == nil {
		return repository.ErrInvalidLink
	}
	if err := proceed(ctx, approval); err != nil {
		r.mu.Lock()
		r.tokens = tokensBefore
		r.mu.Unlock()
		return err
	}
	return nil
}

func (r *fakeRepo) liveTokenLocked(purpose entity.TokenPurpose, tokenHash []byte) *entity.EmailToken {
	for _, token := range r.tokens {
		if token.Purpose == purpose && string(token.TokenHash) == string(tokenHash) && token.ExpiresAt.After(r.clock.Now()) {
			return token
		}
	}
	return nil
}

func (r *fakeRepo) takeLiveTokenLocked(purpose entity.TokenPurpose, tokenHash []byte) *entity.EmailToken {
	token := r.liveTokenLocked(purpose, tokenHash)
	if token == nil {
		return nil
	}
	r.removeTokens(func(t *entity.EmailToken) bool { return t == token })
	consumed := *token
	return &consumed
}

func (r *fakeRepo) isEmailHeldByOtherLocked(userID, email string) bool {
	for id, other := range r.users {
		if id != userID && other.Email == email {
			return true
		}
	}
	return false
}

func (r *fakeRepo) PurgeExpired(_ context.Context, _ time.Time) error { return nil }

// storedData junta tudo o que o repositório guarda, em texto, para as asserções do REQ-17.
func (r *fakeRepo) storedData() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var data strings.Builder
	for _, user := range r.users {
		fmt.Fprintf(&data, "%+v\n", *user)
	}
	for _, session := range r.sessions {
		fmt.Fprintf(&data, "%+v %x\n", *session, session.TokenHash)
	}
	for _, token := range r.tokens {
		fmt.Fprintf(&data, "%+v %x\n", *token, token.TokenHash)
	}
	return data.String()
}

func (r *fakeRepo) userByEmail(email string) *entity.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			return r.copyOf(user)
		}
	}
	return nil
}

type env struct {
	svc    service.ContasService
	repo   *fakeRepo
	mailer *fakeMailer
	tokens *fakeTokens
	clock  *testClock
}

func newEnv() *env {
	clock := &testClock{now: time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)}
	repo := newFakeRepo(clock)
	mailer := &fakeMailer{}
	tokens := &fakeTokens{}
	svc := service.NewContasService(repo, mailer, tokens, service.Options{
		BaseURL: testBaseURL, RefreshTokenTTL: testRefreshTTL, Now: clock.Now,
	})
	return &env{svc: svc, repo: repo, mailer: mailer, tokens: tokens, clock: clock}
}

func validSignUp(email, username string) service.SignUpInput {
	return service.SignUpInput{
		Email: email, Password: testPassword, DisplayName: "Ana", Username: username,
		TermsAccepted: true, AgeDeclared: true,
	}
}

// signUp cadastra uma conta válida e devolve o id; limpa os e-mails enviados.
func (e *env) signUp(t *testing.T, email, username string) string {
	t.Helper()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(email, username))
	requireNoError(t, appErr)
	e.mailer.clear()
	return e.repo.userByEmail(strings.ToLower(email)).ID
}

func (e *env) confirmEmail(t *testing.T, email, username string) string {
	t.Helper()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(email, username))
	requireNoError(t, appErr)
	requireNoError(t, e.svc.ConfirmEmail(context.Background(), linkToken(t, e.mailer.last())))
	e.mailer.clear()
	return e.repo.userByEmail(strings.ToLower(email)).ID
}

func (e *env) login(t *testing.T, email string) *service.SessionResult {
	t.Helper()
	session, appErr := e.svc.Login(context.Background(), service.LoginInput{Email: email, Password: testPassword})
	requireNoError(t, appErr)
	return session
}

func (e *env) loginWith(email, plain, ip string) (*service.SessionResult, errs.AppError) {
	return e.svc.Login(context.Background(), service.LoginInput{Email: email, Password: plain, IP: ip})
}

func (e *env) refresh(refreshToken string) (*service.SessionResult, errs.AppError) {
	return e.svc.RefreshSession(context.Background(), refreshToken)
}

func linkToken(t *testing.T, mail entity.Mail) string {
	t.Helper()
	parsed, err := url.Parse(mail.Link)
	if err != nil {
		t.Fatalf("link inválido %q: %v", mail.Link, err)
	}
	token := parsed.Query().Get("token")
	if token == "" {
		t.Fatalf("link sem token: %q", mail.Link)
	}
	return token
}

func requireNoError(t *testing.T, appErr errs.AppError) {
	t.Helper()
	if appErr.Exists() {
		t.Fatalf("erro inesperado: %v", appErr)
	}
}

func requireError(t *testing.T, appErr errs.AppError, want errs.AppError) {
	t.Helper()
	if !appErr.Exists() || !errors.Is(appErr, want) {
		t.Fatalf("erro = %v, esperado %s", appErr, want.Code)
	}
}

func errorsIs(appErr errs.AppError, want errs.AppError) bool { return errors.Is(appErr, want) }

func newUnknownID() string { return uuid.NewString() }
