//gofi:context contas
package entity

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type TokenPurpose string

const (
	PurposeEmailConfirmation       TokenPurpose = "email_confirmation"
	PurposePasswordReset           TokenPurpose = "password_reset"
	PurposeEmailChangeApproval     TokenPurpose = "email_change_approval"
	PurposeEmailChangeConfirmation TokenPurpose = "email_change_confirmation"
)

const (
	MaxEmailLength       = 254
	MinPasswordLength    = 8
	MaxPasswordLength    = 128
	MaxDisplayNameLength = 60
)

type RequestType string

const (
	RequestConfirmationResend RequestType = "confirmation_resend"
	RequestPasswordReset      RequestType = "password_reset"
	RequestEmailChange        RequestType = "email_change"
)

type KeyKind string

const (
	KeyKindEmail KeyKind = "email"
	KeyKindIP    KeyKind = "ip"
)

// LimitKey identifica o e-mail ou o IP a que um bloqueio ou limite se aplica.
type LimitKey struct {
	Kind  KeyKind
	Value string
}

type User struct {
	ID               string     `db:"id"`
	Email            string     `db:"email"`
	Username         string     `db:"username"`
	DisplayName      string     `db:"display_name"`
	PasswordHash     string     `db:"password_hash"`
	Role             Role       `db:"role"`
	EmailConfirmedAt *time.Time `db:"email_confirmed_at"`
	TermsAcceptedAt  time.Time  `db:"terms_accepted_at"`
	AgeDeclaredAt    time.Time  `db:"age_declared_at"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

func (u User) IsEmailConfirmed() bool {
	return u.EmailConfirmedAt != nil
}

type Session struct {
	ID        int64     `db:"id"`
	UserID    string    `db:"user_id"`
	TokenHash []byte    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
}

type EmailToken struct {
	UserID    string       `db:"user_id"`
	Purpose   TokenPurpose `db:"purpose"`
	TokenHash []byte       `db:"token_hash"`
	NewEmail  *string      `db:"new_email"`
	ExpiresAt time.Time    `db:"expires_at"`
}

// MailKind é a finalidade de um e-mail transacional.
type MailKind string

const (
	MailEmailConfirmation       MailKind = "email_confirmation"
	MailPasswordReset           MailKind = "password_reset"
	MailEmailChangeApproval     MailKind = "email_change_approval"
	MailEmailChangeConfirmation MailKind = "email_change_confirmation"
)

// Mail descreve um e-mail com link; o token só existe dentro de Link.
type Mail struct {
	Kind   MailKind
	To     string
	Link   string
	UserID string
}

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func IsValidUsername(username string) bool {
	return usernamePattern.MatchString(username)
}

// IsValidEmail aceita só o endereço simples `local@dominio`, sem nome nem comentário.
func IsValidEmail(email string) bool {
	if email == "" || len(email) > MaxEmailLength {
		return false
	}
	parsed, err := mail.ParseAddress(email)
	return err == nil && parsed.Address == email
}

func IsValidPassword(password string) bool {
	length := utf8.RuneCountInString(password)
	return length >= MinPasswordLength && length <= MaxPasswordLength
}

func IsValidDisplayName(displayName string) bool {
	length := utf8.RuneCountInString(strings.TrimSpace(displayName))
	return length >= 1 && length <= MaxDisplayNameLength
}
