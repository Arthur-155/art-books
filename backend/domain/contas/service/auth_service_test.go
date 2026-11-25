//gofi:context contas
package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Arthur-155/art-books/domain/contas/entity"
	"github.com/Arthur-155/art-books/domain/contas/service"
)

const (
	wrongPassword = "senhaerrada1"
	newPassword   = "outrasenha1"
	unknownEmail  = "ninguem@exemplo.com"
	otherEmail    = "beto@exemplo.com"
	day           = 24 * time.Hour
)

func (e *env) failLogins(email, ip string, times int) {
	for range times {
		e.loginWith(email, wrongPassword, ip)
	}
}

func (e *env) requireLoginOK(t *testing.T, email, ip string) {
	t.Helper()
	if _, appErr := e.loginWith(email, testPassword, ip); appErr.Exists() {
		t.Fatalf("login deveria abrir a sessão: %v", appErr)
	}
}

func TestLogin_UnconfirmedUserSessionReportsEmailNotConfirmed(t *testing.T) { // gofi:tests REQ-04.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	if session.User.EmailConfirmed || e.tokens.last().emailConfirmed {
		t.Fatal("a sessão deveria informar o e-mail como não confirmado")
	}
}

func TestRefreshSession_ReportsEmailConfirmedAfterConfirmation(t *testing.T) { // gofi:tests REQ-04.2
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	confirmation := linkToken(t, e.mailer.last())
	session := e.login(t, testEmail)
	if session.User.EmailConfirmed {
		t.Fatal("antes da confirmação o e-mail deveria constar como não confirmado")
	}
	requireNoError(t, e.svc.ConfirmEmail(context.Background(), confirmation))
	renewed, appErr := e.refresh(session.RefreshToken)
	requireNoError(t, appErr)
	if !renewed.User.EmailConfirmed || !e.tokens.last().emailConfirmed {
		t.Fatal("a sessão renovada deveria informar o e-mail como confirmado")
	}
}

func TestLogin_ValidCredentialsOpenSessionWithIDAndRole(t *testing.T) { // gofi:tests REQ-05.1
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	session, appErr := e.loginWith("Ana@Exemplo.com", testPassword, "")
	requireNoError(t, appErr)
	if session.User.ID != userID || session.User.Role != entity.RoleUser {
		t.Fatalf("usuário da sessão = %+v", session.User)
	}
	if issued := e.tokens.last(); issued.userID != userID || issued.role != entity.RoleUser {
		t.Fatalf("token emitido = %+v", issued)
	}
}

func TestLogin_WrongPasswordAndUnknownEmailRespondTheSame(t *testing.T) { // gofi:tests REQ-05.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	_, wrong := e.loginWith(testEmail, wrongPassword, "")
	_, unknown := e.loginWith(unknownEmail, wrongPassword, "")
	requireError(t, wrong, service.ErrContasInvalidCredentials)
	requireError(t, unknown, service.ErrContasInvalidCredentials)
	if wrong.Error() != unknown.Error() {
		t.Fatalf("mensagens diferentes: %q e %q", wrong.Error(), unknown.Error())
	}
}

func TestLogin_FiveFailuresLockEmailAcrossIPs(t *testing.T) { // gofi:tests REQ-06.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	for index := range lockFailureLimit {
		e.failLogins(testEmail, fmt.Sprintf("10.0.0.%d", index+1), 1)
	}
	_, appErr := e.loginWith(testEmail, testPassword, "10.0.0.99")
	requireError(t, appErr, service.ErrContasLoginLocked)
}

func TestLogin_FiveFailuresLockIPAcrossEmails(t *testing.T) { // gofi:tests REQ-06.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	for index := range lockFailureLimit {
		e.failLogins(fmt.Sprintf("leitor%d@exemplo.com", index), "10.0.0.1", 1)
	}
	_, appErr := e.loginWith(testEmail, testPassword, "10.0.0.1")
	requireError(t, appErr, service.ErrContasLoginLocked)
}

func TestLogin_LockEndsAfterFifteenMinutes(t *testing.T) { // gofi:tests REQ-06.3
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.failLogins(testEmail, "10.0.0.1", lockFailureLimit)
	e.clock.advance(16 * time.Minute)
	e.requireLoginOK(t, testEmail, "10.0.0.2")
}

func TestLogin_UnknownEmailLocksLikeRegisteredEmail(t *testing.T) { // gofi:tests REQ-06.4
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.failLogins(unknownEmail, "", lockFailureLimit)
	e.failLogins(testEmail, "", lockFailureLimit)
	_, unknown := e.loginWith(unknownEmail, wrongPassword, "")
	_, registered := e.loginWith(testEmail, wrongPassword, "")
	requireError(t, unknown, service.ErrContasLoginLocked)
	requireError(t, registered, service.ErrContasLoginLocked)
	if unknown.Error() != registered.Error() {
		t.Fatalf("mensagens diferentes: %q e %q", unknown.Error(), registered.Error())
	}
}

func TestLogin_SpacedFailuresDoNotLock(t *testing.T) { // gofi:tests REQ-06.5
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.failLogins(testEmail, "", 1)
	e.clock.advance(10 * time.Minute)
	e.failLogins(testEmail, "", 3)
	e.clock.advance(6 * time.Minute)
	e.failLogins(testEmail, "", 1)
	e.requireLoginOK(t, testEmail, "")
}

func TestLogin_AttemptsDuringLockDoNotExtendIt(t *testing.T) { // gofi:tests REQ-06.6
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.failLogins(testEmail, "", lockFailureLimit)
	e.clock.advance(10 * time.Minute)
	for range 3 {
		_, appErr := e.loginWith(testEmail, testPassword, "")
		requireError(t, appErr, service.ErrContasLoginLocked)
	}
	e.clock.advance(5*time.Minute + time.Second)
	e.requireLoginOK(t, testEmail, "")
}

func TestLogin_CorrectLoginResetsFailureCount(t *testing.T) { // gofi:tests REQ-06.7
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.failLogins(testEmail, "", lockFailureLimit-1)
	e.clock.advance(2 * time.Minute)
	e.requireLoginOK(t, testEmail, "")
	e.failLogins(testEmail, "", lockFailureLimit-1)
	e.requireLoginOK(t, testEmail, "")
}

func TestRefreshSession_WithinDeadlineGivesNewThirtyDaySession(t *testing.T) { // gofi:tests REQ-07.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	e.clock.advance(29 * day)
	renewed, appErr := e.refresh(session.RefreshToken)
	requireNoError(t, appErr)
	if renewed.RefreshToken == "" || renewed.RefreshToken == session.RefreshToken {
		t.Fatal("a renovação deveria trazer uma credencial nova")
	}
	if want := e.clock.Now().Add(30 * day); !renewed.RefreshExpiresAt.Equal(want) {
		t.Fatalf("validade = %v, esperado %v", renewed.RefreshExpiresAt, want)
	}
}

func TestRefreshSession_ExpiredSessionIsRejected(t *testing.T) { // gofi:tests REQ-07.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	e.clock.advance(31 * day)
	_, appErr := e.refresh(session.RefreshToken)
	requireError(t, appErr, service.ErrContasInvalidSession)
}

func TestRefreshSession_UsedCredentialIsRejectedAndNewOneKeepsWorking(t *testing.T) { // gofi:tests REQ-07.3
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	renewed, appErr := e.refresh(session.RefreshToken)
	requireNoError(t, appErr)
	_, appErr = e.refresh(session.RefreshToken)
	requireError(t, appErr, service.ErrContasInvalidSession)
	_, appErr = e.refresh(renewed.RefreshToken)
	requireNoError(t, appErr)
}

func TestRefreshSession_EndedSessionIsRejected(t *testing.T) { // gofi:tests REQ-07.4
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	requireNoError(t, e.svc.ChangePassword(context.Background(), userID, testPassword, newPassword))
	_, appErr := e.refresh(session.RefreshToken)
	requireError(t, appErr, service.ErrContasInvalidSession)
}

func TestLogout_EndsOnlyCurrentDeviceSession(t *testing.T) { // gofi:tests REQ-08.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	first := e.login(t, testEmail)
	second := e.login(t, testEmail)
	requireNoError(t, e.svc.Logout(context.Background(), first.RefreshToken))
	_, appErr := e.refresh(first.RefreshToken)
	requireError(t, appErr, service.ErrContasInvalidSession)
	_, appErr = e.refresh(second.RefreshToken)
	requireNoError(t, appErr)
}

func TestLogout_UnknownOrEndedCredentialIsAccepted(t *testing.T) { // gofi:tests REQ-08.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	requireNoError(t, e.svc.Logout(context.Background(), "credencial-inexistente"))
	requireNoError(t, e.svc.Logout(context.Background(), session.RefreshToken))
	requireNoError(t, e.svc.Logout(context.Background(), session.RefreshToken))
}

func TestRequestPasswordReset_KnownEmailReceivesLink(t *testing.T) { // gofi:tests REQ-09.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	mails := e.mailer.sent()
	if len(mails) != 1 || mails[0].To != testEmail || mails[0].Kind != entity.MailPasswordReset {
		t.Fatalf("e-mails = %+v", mails)
	}
}

func TestRequestPasswordReset_UnknownEmailGetsSameResponseAndNoMail(t *testing.T) { // gofi:tests REQ-09.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), unknownEmail, ""))
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado")
	}
}

func TestResetPassword_WithinDeadlineChangesPasswordAndEndsSessions(t *testing.T) { // gofi:tests REQ-09.3
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	otherDevice := e.login(t, testEmail)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	token := linkToken(t, e.mailer.last())
	e.clock.advance(29 * time.Minute)
	requireNoError(t, e.svc.ResetPassword(context.Background(), token, newPassword))
	if _, appErr := e.loginWith(testEmail, newPassword, ""); appErr.Exists() {
		t.Fatalf("login com a senha nova recusado: %v", appErr)
	}
	_, appErr := e.loginWith(testEmail, testPassword, "")
	requireError(t, appErr, service.ErrContasInvalidCredentials)
	_, appErr = e.refresh(otherDevice.RefreshToken)
	requireError(t, appErr, service.ErrContasInvalidSession)
}

func TestResetPassword_ExpiredLinkIsRejectedAndPasswordKept(t *testing.T) { // gofi:tests REQ-09.4
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	token := linkToken(t, e.mailer.last())
	e.clock.advance(31 * time.Minute)
	requireError(t, e.svc.ResetPassword(context.Background(), token, newPassword), service.ErrContasInvalidLink)
	e.requireLoginOK(t, testEmail, "")
}

func TestResetPassword_LinkUsedTwiceIsRejected(t *testing.T) { // gofi:tests REQ-09.5
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	token := linkToken(t, e.mailer.last())
	requireNoError(t, e.svc.ResetPassword(context.Background(), token, newPassword))
	requireError(t, e.svc.ResetPassword(context.Background(), token, "terceirasenha1"), service.ErrContasInvalidLink)
}

func TestResetPassword_ShortPasswordIsInvalidAndLinkKeepsWorking(t *testing.T) { // gofi:tests REQ-09.6
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	token := linkToken(t, e.mailer.last())
	requireError(t, e.svc.ResetPassword(context.Background(), token, "1234567"), service.ErrContasInvalidInput)
	requireNoError(t, e.svc.ResetPassword(context.Background(), token, newPassword))
}

func TestRequestPasswordReset_ExcessRefusalIsSameForUnknownEmail(t *testing.T) { // gofi:tests REQ-09.7
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), unknownEmail, ""))
	e.clock.advance(30 * time.Second)
	known := e.svc.RequestPasswordReset(context.Background(), testEmail, "")
	unknown := e.svc.RequestPasswordReset(context.Background(), unknownEmail, "")
	requireError(t, known, service.ErrContasTooManyRequests)
	requireError(t, unknown, service.ErrContasTooManyRequests)
	if known.Error() != unknown.Error() {
		t.Fatalf("respostas diferentes: %q e %q", known.Error(), unknown.Error())
	}
}

func TestRequestPasswordReset_NewRequestInvalidatesPreviousLink(t *testing.T) { // gofi:tests REQ-09.8
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	previous := linkToken(t, e.mailer.last())
	e.clock.advance(time.Minute + time.Second)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	requireError(t, e.svc.ResetPassword(context.Background(), previous, newPassword), service.ErrContasInvalidLink)
	requireNoError(t, e.svc.ResetPassword(context.Background(), linkToken(t, e.mailer.last()), newPassword))
}

func TestChangePassword_CorrectCurrentPasswordChangesItAndEndsAllSessions(t *testing.T) { // gofi:tests REQ-10.1
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	first := e.login(t, testEmail)
	second := e.login(t, testEmail)
	requireNoError(t, e.svc.ChangePassword(context.Background(), userID, testPassword, newPassword))
	if _, appErr := e.loginWith(testEmail, newPassword, ""); appErr.Exists() {
		t.Fatalf("login com a senha nova recusado: %v", appErr)
	}
	for _, session := range []*service.SessionResult{first, second} {
		_, appErr := e.refresh(session.RefreshToken)
		requireError(t, appErr, service.ErrContasInvalidSession)
	}
}

func TestChangePassword_WrongCurrentPasswordIsRefusedAndNothingChanges(t *testing.T) { // gofi:tests REQ-10.2
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	if appErr := e.svc.ChangePassword(context.Background(), userID, wrongPassword, newPassword); !appErr.Exists() {
		t.Fatal("a troca com a senha atual errada deveria ser recusada")
	}
	e.requireLoginOK(t, testEmail, "")
	_, appErr := e.refresh(session.RefreshToken)
	requireNoError(t, appErr)
}
