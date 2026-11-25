//gofi:context contas
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/joaoprofile/gofi-sdk-go/sqln"
	"github.com/joaoprofile/gofi-sdk-go/sqln/connection"
	"github.com/joaoprofile/gofi-sdk-go/sqln/transaction"

	"github.com/Arthur-155/art-books/domain/contas/entity"
)

var (
	ErrEmailTaken     = errors.New("email already in use")
	ErrUsernameTaken  = errors.New("username already in use")
	ErrInvalidSession = errors.New("session not found or expired")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidLink    = errors.New("link not found, used or expired")
	ErrTooManyRequest = errors.New("too many email requests")
)

const (
	uqUserEmail    = "uq_user_email"
	uqUserUsername = "uq_user_username"

	userSelectFields = `id, email, username, display_name, password_hash, role, email_confirmed_at,
		terms_accepted_at, age_declared_at, created_at, updated_at`

	userInsertQuery = `INSERT INTO "user" (id, email, username, display_name, password_hash, role,
		terms_accepted_at, age_declared_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	userFindByIDQuery    = `SELECT ` + userSelectFields + ` FROM "user" WHERE id = $1`
	userFindByEmailQuery = `SELECT ` + userSelectFields + ` FROM "user" WHERE email = $1`
	userEmailTakenQuery  = `SELECT TRUE FROM "user" WHERE email = $1 LIMIT 1`
	userConfirmQuery     = `UPDATE "user" SET email_confirmed_at = now(), updated_at = now() WHERE id = $1 AND email_confirmed_at IS NULL`
	userPasswordQuery    = `UPDATE "user" SET password_hash = $2, updated_at = now() WHERE id = $1`
	userEmailQuery       = `UPDATE "user" SET email = $2, email_confirmed_at = now(), updated_at = now() WHERE id = $1`
	userRoleQuery        = `UPDATE "user" SET role = $2, updated_at = now() WHERE email = lower($1)`
	userDeleteQuery      = `DELETE FROM "user" WHERE id = $1`

	sessionInsertQuery  = `INSERT INTO user_session (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	sessionConsumeQuery = `DELETE FROM user_session WHERE token_hash = $1 AND expires_at > now() RETURNING user_id`
	sessionDeleteQuery  = `DELETE FROM user_session WHERE token_hash = $1`
	sessionPurgeUser    = `DELETE FROM user_session WHERE user_id = $1`

	tokenUpsertQuery = `INSERT INTO email_token (user_id, purpose, token_hash, new_email, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, purpose) DO UPDATE
		SET token_hash = EXCLUDED.token_hash, new_email = EXCLUDED.new_email,
			expires_at = EXCLUDED.expires_at, created_at = now()`
	tokenConsumeQuery = `DELETE FROM email_token WHERE purpose = $1 AND token_hash = $2 AND expires_at > now()
		RETURNING user_id, purpose, token_hash, new_email, expires_at`
	tokenDeleteQuery  = `DELETE FROM email_token WHERE user_id = $1 AND purpose = $2`
	tokenPendingQuery = `SELECT new_email FROM email_token
		WHERE user_id = $1 AND purpose IN ('email_change_approval', 'email_change_confirmation')
			AND expires_at > now() LIMIT 1`
	tokenDeleteAllQuery = `DELETE FROM email_token WHERE user_id = $1`

	loginLockedQuery = `SELECT TRUE FROM login_lock WHERE key_kind = $1 AND key_value = $2 AND locked_until > now() LIMIT 1`
	loginFailQuery   = `INSERT INTO login_lock (key_kind, key_value, failure_count, window_started_at)
		VALUES ($1, $2, 1, now())
		ON CONFLICT (key_kind, key_value) DO UPDATE SET
			locked_until = CASE WHEN (CASE WHEN login_lock.window_started_at < now() - interval '15 minutes'
					THEN 1 ELSE login_lock.failure_count + 1 END) >= 5
				THEN now() + interval '15 minutes' ELSE login_lock.locked_until END,
			failure_count = CASE WHEN (CASE WHEN login_lock.window_started_at < now() - interval '15 minutes'
					THEN 1 ELSE login_lock.failure_count + 1 END) >= 5
				THEN 0
				WHEN login_lock.window_started_at < now() - interval '15 minutes' THEN 1
				ELSE login_lock.failure_count + 1 END,
			window_started_at = CASE WHEN login_lock.window_started_at < now() - interval '15 minutes'
				THEN now() ELSE login_lock.window_started_at END`
	loginClearQuery = `DELETE FROM login_lock WHERE key_kind = $1 AND key_value = $2`

	requestLockQuery   = `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`
	requestCountQuery  = `SELECT count(*) FROM email_request WHERE request_type = $1 AND key_kind = $2 AND key_value = $3 AND requested_at > now() - $4::interval`
	requestInsertQuery = `INSERT INTO email_request (request_type, key_kind, key_value) VALUES ($1, $2, $3)`

	requestsPerMinute = 1
	requestsPerHour   = 5

	purgeSessionsQuery = `DELETE FROM user_session WHERE expires_at <= $1`
	purgeTokensQuery   = `DELETE FROM email_token WHERE expires_at <= $1`
	purgeRequestsQuery = `DELETE FROM email_request WHERE requested_at < $1::timestamptz - interval '1 hour'`
	purgeLocksQuery    = `DELETE FROM login_lock WHERE (locked_until IS NULL OR locked_until <= $1::timestamptz)
		AND window_started_at < $1::timestamptz - interval '15 minutes'`
)

type ContasRepository interface {
	SaveUser(ctx context.Context, user *entity.User) error
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	MarkEmailConfirmed(ctx context.Context, userID string) error
	SetRoleByEmail(ctx context.Context, email string, role entity.Role) error
	DeleteUser(ctx context.Context, userID string) error

	// ResetPassword consome o link, troca a senha e encerra sessões e links pendentes numa transação.
	ResetPassword(ctx context.Context, tokenHash []byte, passwordHash string) error
	// ChangePassword troca a senha e encerra sessões e links de recuperação numa transação.
	ChangePassword(ctx context.Context, userID, passwordHash string) error
	// ConfirmEmailChange consome o link e troca o e-mail; qualquer falha devolve o link.
	ConfirmEmailChange(ctx context.Context, tokenHash []byte) error
	// ApproveEmailChange consome a aprovação e roda proceed; se proceed falha, o link é devolvido.
	ApproveEmailChange(ctx context.Context, tokenHash []byte, proceed func(ctx context.Context, approval *entity.EmailToken) error) error

	SaveSession(ctx context.Context, session *entity.Session) error
	RotateSession(ctx context.Context, oldHash []byte, session *entity.Session) (*entity.Session, error)
	DeleteSession(ctx context.Context, tokenHash []byte) error

	ReplaceEmailToken(ctx context.Context, token *entity.EmailToken) error
	ConsumeEmailToken(ctx context.Context, purpose entity.TokenPurpose, tokenHash []byte) (*entity.EmailToken, error)
	DeleteEmailTokens(ctx context.Context, userID string, purposes ...entity.TokenPurpose) error
	FindPendingEmailChange(ctx context.Context, userID string) (*string, error)

	IsLoginLocked(ctx context.Context, keys []entity.LimitKey) (bool, error)
	RecordLoginFailure(ctx context.Context, key entity.LimitKey) error
	ClearLoginFailures(ctx context.Context, key entity.LimitKey) error
	// AcceptEmailRequest confere o limite e registra o pedido de forma atômica; excesso devolve ErrTooManyRequest.
	AcceptEmailRequest(ctx context.Context, requestType entity.RequestType, keys []entity.LimitKey) error
	PurgeExpired(ctx context.Context, now time.Time) error
}

type contasRepository struct {
	stm sqln.Statement
	tx  transaction.Transaction
}

// NewContasRepository não toca o banco: a conexão global só existe depois do Build.
func NewContasRepository() ContasRepository {
	return &contasRepository{stm: sqln.NewStatement(), tx: transaction.NewTransaction()}
}

// gofi:implements REQ-01
func (r *contasRepository) SaveUser(ctx context.Context, user *entity.User) error {
	_, err := r.stm.Execute(ctx, userInsertQuery, user.ID, user.Email, user.Username, user.DisplayName,
		user.PasswordHash, string(user.Role), user.TermsAcceptedAt, user.AgeDeclaredAt)
	if connection.IsUniqueViolation(err, uqUserEmail) {
		return ErrEmailTaken
	}
	if connection.IsUniqueViolation(err, uqUserUsername) {
		return ErrUsernameTaken
	}
	return err
}

func (r *contasRepository) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	return sqln.Find[entity.User](ctx, userFindByIDQuery, id).UniqueResult()
}

func (r *contasRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return sqln.Find[entity.User](ctx, userFindByEmailQuery, email).UniqueResult()
}

func (r *contasRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	found, err := sqln.Find[bool](ctx, userEmailTakenQuery, email).UniqueResult()
	return found != nil, err
}

func (r *contasRepository) MarkEmailConfirmed(ctx context.Context, userID string) error {
	_, err := r.stm.Execute(ctx, userConfirmQuery, userID)
	return err
}

// gofi:implements REQ-09
func (r *contasRepository) ResetPassword(ctx context.Context, tokenHash []byte, passwordHash string) error {
	return r.tx.Execute(ctx, func(ctx context.Context) error {
		token, err := r.consumeRequired(ctx, entity.PurposePasswordReset, tokenHash)
		if err != nil {
			return err
		}
		return r.replacePassword(ctx, token.UserID, passwordHash)
	})
}

// gofi:implements REQ-10
func (r *contasRepository) ChangePassword(ctx context.Context, userID, passwordHash string) error {
	return r.tx.Execute(ctx, func(ctx context.Context) error {
		return r.replacePassword(ctx, userID, passwordHash)
	})
}

// replacePassword roda dentro da transação do chamador.
func (r *contasRepository) replacePassword(ctx context.Context, userID, passwordHash string) error {
	if _, err := r.stm.Execute(ctx, userPasswordQuery, userID, passwordHash); err != nil {
		return err
	}
	if _, err := r.stm.Execute(ctx, tokenDeleteQuery, userID, string(entity.PurposePasswordReset)); err != nil {
		return err
	}
	return r.deleteUserSessions(ctx, userID)
}

// gofi:implements REQ-11
func (r *contasRepository) ConfirmEmailChange(ctx context.Context, tokenHash []byte) error {
	return r.tx.Execute(ctx, func(ctx context.Context) error {
		token, err := r.consumeRequired(ctx, entity.PurposeEmailChangeConfirmation, tokenHash)
		if err != nil {
			return err
		}
		if err := r.updateEmail(ctx, token.UserID, *token.NewEmail); err != nil {
			return err
		}
		if _, err := r.stm.Execute(ctx, tokenDeleteAllQuery, token.UserID); err != nil {
			return err
		}
		return r.deleteUserSessions(ctx, token.UserID)
	})
}

// gofi:implements REQ-14
func (r *contasRepository) ApproveEmailChange(ctx context.Context, tokenHash []byte, proceed func(ctx context.Context, approval *entity.EmailToken) error) error {
	return r.tx.Execute(ctx, func(ctx context.Context) error {
		approval, err := r.consumeRequired(ctx, entity.PurposeEmailChangeApproval, tokenHash)
		if err != nil {
			return err
		}
		return proceed(ctx, approval)
	})
}

func (r *contasRepository) consumeRequired(ctx context.Context, purpose entity.TokenPurpose, tokenHash []byte) (*entity.EmailToken, error) {
	token, err := r.ConsumeEmailToken(ctx, purpose, tokenHash)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, ErrInvalidLink
	}
	return token, nil
}

// updateEmail é decidido pela restrição uq_user_email quando outro usuário tomou o endereço.
func (r *contasRepository) updateEmail(ctx context.Context, userID, email string) error {
	_, err := r.stm.Execute(ctx, userEmailQuery, userID, email)
	if connection.IsUniqueViolation(err, uqUserEmail) {
		return ErrEmailTaken
	}
	return err
}

// gofi:implements REQ-15
func (r *contasRepository) SetRoleByEmail(ctx context.Context, email string, role entity.Role) error {
	result, err := r.stm.Execute(ctx, userRoleQuery, email, string(role))
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// gofi:implements REQ-12
func (r *contasRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.stm.Execute(ctx, userDeleteQuery, userID)
	return err
}

func (r *contasRepository) SaveSession(ctx context.Context, session *entity.Session) error {
	_, err := r.stm.Execute(ctx, sessionInsertQuery, session.UserID, session.TokenHash, session.ExpiresAt)
	return err
}

// gofi:implements REQ-07
func (r *contasRepository) RotateSession(ctx context.Context, oldHash []byte, session *entity.Session) (*entity.Session, error) {
	rotated := &entity.Session{TokenHash: session.TokenHash, ExpiresAt: session.ExpiresAt}
	err := r.tx.Execute(ctx, func(ctx context.Context) error {
		userID, err := sqln.Find[string](ctx, sessionConsumeQuery, oldHash).UniqueResult()
		if err != nil {
			return err
		}
		if userID == nil {
			return ErrInvalidSession
		}
		rotated.UserID = *userID
		return r.SaveSession(ctx, rotated)
	})
	if err != nil {
		return nil, err
	}
	return rotated, nil
}

// gofi:implements REQ-08
func (r *contasRepository) DeleteSession(ctx context.Context, tokenHash []byte) error {
	_, err := r.stm.Execute(ctx, sessionDeleteQuery, tokenHash)
	return err
}

func (r *contasRepository) deleteUserSessions(ctx context.Context, userID string) error {
	_, err := r.stm.Execute(ctx, sessionPurgeUser, userID)
	return err
}

// gofi:implements REQ-17
func (r *contasRepository) ReplaceEmailToken(ctx context.Context, token *entity.EmailToken) error {
	_, err := r.stm.Execute(ctx, tokenUpsertQuery, token.UserID, string(token.Purpose), token.TokenHash,
		token.NewEmail, token.ExpiresAt)
	return err
}

func (r *contasRepository) ConsumeEmailToken(ctx context.Context, purpose entity.TokenPurpose, tokenHash []byte) (*entity.EmailToken, error) {
	return sqln.Find[entity.EmailToken](ctx, tokenConsumeQuery, string(purpose), tokenHash).UniqueResult()
}

func (r *contasRepository) DeleteEmailTokens(ctx context.Context, userID string, purposes ...entity.TokenPurpose) error {
	for _, purpose := range purposes {
		if _, err := r.stm.Execute(ctx, tokenDeleteQuery, userID, string(purpose)); err != nil {
			return err
		}
	}
	return nil
}

func (r *contasRepository) FindPendingEmailChange(ctx context.Context, userID string) (*string, error) {
	return sqln.Find[string](ctx, tokenPendingQuery, userID).UniqueResult()
}

func (r *contasRepository) IsLoginLocked(ctx context.Context, keys []entity.LimitKey) (bool, error) {
	for _, key := range keys {
		found, err := sqln.Find[bool](ctx, loginLockedQuery, string(key.Kind), key.Value).UniqueResult()
		if err != nil || found != nil {
			return found != nil, err
		}
	}
	return false, nil
}

func (r *contasRepository) RecordLoginFailure(ctx context.Context, key entity.LimitKey) error {
	_, err := r.stm.Execute(ctx, loginFailQuery, string(key.Kind), key.Value)
	return err
}

func (r *contasRepository) ClearLoginFailures(ctx context.Context, key entity.LimitKey) error {
	_, err := r.stm.Execute(ctx, loginClearQuery, string(key.Kind), key.Value)
	return err
}

// gofi:implements REQ-13
func (r *contasRepository) AcceptEmailRequest(ctx context.Context, requestType entity.RequestType, keys []entity.LimitKey) error {
	return r.tx.Execute(ctx, func(ctx context.Context) error {
		for _, key := range keys {
			if err := r.lockRequestKey(ctx, requestType, key); err != nil {
				return err
			}
		}
		for _, key := range keys {
			if err := r.ensureWithinLimit(ctx, requestType, key); err != nil {
				return err
			}
		}
		return r.insertRequests(ctx, requestType, keys)
	})
}

// lockRequestKey serializa os pedidos da mesma chave até o fim da transação.
func (r *contasRepository) lockRequestKey(ctx context.Context, requestType entity.RequestType, key entity.LimitKey) error {
	lockName := string(requestType) + "|" + string(key.Kind) + "|" + key.Value
	_, err := r.stm.Execute(ctx, requestLockQuery, lockName)
	return err
}

func (r *contasRepository) ensureWithinLimit(ctx context.Context, requestType entity.RequestType, key entity.LimitKey) error {
	lastMinute, err := r.countRequests(ctx, requestType, key, "1 minute")
	if err != nil {
		return err
	}
	lastHour, err := r.countRequests(ctx, requestType, key, "1 hour")
	if err != nil {
		return err
	}
	if lastMinute >= requestsPerMinute || lastHour >= requestsPerHour {
		return ErrTooManyRequest
	}
	return nil
}

func (r *contasRepository) countRequests(ctx context.Context, requestType entity.RequestType, key entity.LimitKey, window string) (int, error) {
	total, err := sqln.Find[int](ctx, requestCountQuery, string(requestType), string(key.Kind), key.Value, window).UniqueResult()
	if err != nil || total == nil {
		return 0, err
	}
	return *total, nil
}

func (r *contasRepository) insertRequests(ctx context.Context, requestType entity.RequestType, keys []entity.LimitKey) error {
	for _, key := range keys {
		if _, err := r.stm.Execute(ctx, requestInsertQuery, string(requestType), string(key.Kind), key.Value); err != nil {
			return err
		}
	}
	return nil
}

func (r *contasRepository) PurgeExpired(ctx context.Context, now time.Time) error {
	for _, query := range []string{purgeSessionsQuery, purgeTokensQuery, purgeRequestsQuery, purgeLocksQuery} {
		if _, err := r.stm.Execute(ctx, query, now); err != nil {
			return err
		}
	}
	return nil
}
