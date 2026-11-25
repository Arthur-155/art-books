//gofi:context contas
package service_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Arthur-155/art-books/domain/contas/entity"
	"github.com/Arthur-155/art-books/domain/contas/service"
)

func TestSignUp_ValidCreatesUnconfirmedUserWithLowercaseUsername(t *testing.T) { // gofi:tests REQ-01.1
	e := newEnv()
	result, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, "Ana_Leitora"))
	requireNoError(t, appErr)
	if !result.ConfirmationEmailSent {
		t.Fatal("confirmação deveria ter sido enviada")
	}
	user := e.repo.userByEmail(testEmail)
	if user == nil {
		t.Fatal("conta não criada")
	}
	if user.Role != entity.RoleUser || user.EmailConfirmedAt != nil {
		t.Fatalf("papel=%s confirmado=%v", user.Role, user.EmailConfirmedAt)
	}
	if user.Username != "ana_leitora" {
		t.Fatalf("username = %q", user.Username)
	}
}

func TestSignUp_EmailInOtherCaseIsTaken(t *testing.T) { // gofi:tests REQ-01.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	_, appErr := e.svc.SignUp(context.Background(), validSignUp("ANA@Exemplo.com", "outra_ana"))
	requireError(t, appErr, service.ErrContasEmailTaken)
	if len(e.repo.users) != 1 {
		t.Fatalf("contas = %d, esperado 1", len(e.repo.users))
	}
}

func TestSignUp_UsernameInOtherCaseIsTaken(t *testing.T) { // gofi:tests REQ-01.3
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	_, appErr := e.svc.SignUp(context.Background(), validSignUp("beto@exemplo.com", "ANA_LEITORA"))
	requireError(t, appErr, service.ErrContasUsernameTaken)
}

func TestSignUp_WithoutTermsOrAgeIsInvalid(t *testing.T) { // gofi:tests REQ-01.4
	e := newEnv()
	noTerms := validSignUp(testEmail, testUsername)
	noTerms.TermsAccepted = false
	noAge := validSignUp(testEmail, testUsername)
	noAge.AgeDeclared = false
	for _, input := range []service.SignUpInput{noTerms, noAge} {
		_, appErr := e.svc.SignUp(context.Background(), input)
		requireError(t, appErr, service.ErrContasInvalidInput)
	}
	if len(e.repo.users) != 0 {
		t.Fatal("nenhuma conta deveria ter sido criada")
	}
}

func TestSignUp_ShortPasswordIsInvalid(t *testing.T) { // gofi:tests REQ-01.5
	e := newEnv()
	for _, short := range []string{"1234567", "ççççççç"} {
		input := validSignUp(testEmail, testUsername)
		input.Password = short
		_, appErr := e.svc.SignUp(context.Background(), input)
		requireError(t, appErr, service.ErrContasInvalidInput)
	}
}

func TestSignUp_UsernameOutsideRuleIsInvalid(t *testing.T) { // gofi:tests REQ-01.6
	e := newEnv()
	for _, username := range []string{"an", strings.Repeat("a", 31), "ana.leitora"} {
		_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, username))
		requireError(t, appErr, service.ErrContasInvalidInput)
	}
}

func TestSignUp_DisplayNameOutsideLimitIsInvalid(t *testing.T) { // gofi:tests REQ-01.7
	e := newEnv()
	for _, name := range []string{"", "   ", strings.Repeat("n", 61)} {
		input := validSignUp(testEmail, testUsername)
		input.DisplayName = name
		_, appErr := e.svc.SignUp(context.Background(), input)
		requireError(t, appErr, service.ErrContasInvalidInput)
	}
}

func TestSignUp_SimultaneousSameEmailCreatesOneAccount(t *testing.T) { // gofi:tests REQ-01.8
	e := newEnv()
	var wg sync.WaitGroup
	results := make(chan bool, 2)
	for _, username := range []string{"ana_um", "ana_dois"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, username))
			if appErr.Exists() {
				results <- errorsIs(appErr, service.ErrContasEmailTaken)
				return
			}
			results <- true
		}()
	}
	wg.Wait()
	close(results)
	accepted := 0
	for ok := range results {
		if !ok {
			t.Fatal("a recusa deveria ser de e-mail em uso")
		}
	}
	for _, user := range e.repo.users {
		if user.Email == testEmail {
			accepted++
		}
	}
	if accepted != 1 {
		t.Fatalf("contas com o e-mail = %d, esperado 1", accepted)
	}
}

func TestSignUp_MalformedEmailIsInvalid(t *testing.T) { // gofi:tests REQ-01.9
	e := newEnv()
	for _, email := range []string{"ana@", strings.Repeat("a", 250) + "@b.com", "Ana <ana@exemplo.com>"} {
		_, appErr := e.svc.SignUp(context.Background(), validSignUp(email, testUsername))
		requireError(t, appErr, service.ErrContasInvalidInput)
	}
}

func TestSignUp_PasswordAboveMaximumIsInvalid(t *testing.T) { // gofi:tests REQ-01.10
	e := newEnv()
	input := validSignUp(testEmail, testUsername)
	input.Password = strings.Repeat("p", 129)
	_, appErr := e.svc.SignUp(context.Background(), input)
	requireError(t, appErr, service.ErrContasInvalidInput)
}

func TestSignUp_SendsConfirmationToRegisteredAddress(t *testing.T) { // gofi:tests REQ-02.1
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	mails := e.mailer.sent()
	if len(mails) != 1 || mails[0].To != testEmail || mails[0].Kind != entity.MailEmailConfirmation {
		t.Fatalf("e-mails = %+v", mails)
	}
}

func TestConfirmEmail_WithinDeadlineConfirms(t *testing.T) { // gofi:tests REQ-02.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.ResendConfirmation(context.Background(), e.repo.userByEmail(testEmail).ID, ""))
	e.clock.advance(23 * time.Hour)
	requireNoError(t, e.svc.ConfirmEmail(context.Background(), linkToken(t, e.mailer.last())))
	if e.repo.userByEmail(testEmail).EmailConfirmedAt == nil {
		t.Fatal("e-mail deveria constar como confirmado")
	}
}

func TestConfirmEmail_ExpiredLinkIsRejected(t *testing.T) { // gofi:tests REQ-02.3
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	e.clock.advance(25 * time.Hour)
	requireError(t, e.svc.ConfirmEmail(context.Background(), linkToken(t, e.mailer.last())), service.ErrContasInvalidLink)
	if e.repo.userByEmail(testEmail).EmailConfirmedAt != nil {
		t.Fatal("e-mail não deveria ter sido confirmado")
	}
}

func TestConfirmEmail_LinkUsedTwiceIsRejected(t *testing.T) { // gofi:tests REQ-02.4
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	token := linkToken(t, e.mailer.last())
	requireNoError(t, e.svc.ConfirmEmail(context.Background(), token))
	requireError(t, e.svc.ConfirmEmail(context.Background(), token), service.ErrContasInvalidLink)
}

func TestResendConfirmation_InvalidatesPreviousLink(t *testing.T) { // gofi:tests REQ-03.1
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	previous := linkToken(t, e.mailer.last())
	e.clock.advance(time.Minute + time.Second)
	requireNoError(t, e.svc.ResendConfirmation(context.Background(), e.repo.userByEmail(testEmail).ID, ""))
	if len(e.mailer.sent()) != 2 {
		t.Fatalf("e-mails = %d, esperado 2", len(e.mailer.sent()))
	}
	requireError(t, e.svc.ConfirmEmail(context.Background(), previous), service.ErrContasInvalidLink)
	requireNoError(t, e.svc.ConfirmEmail(context.Background(), linkToken(t, e.mailer.last())))
}

func TestResendConfirmation_AlreadyConfirmedIsRejected(t *testing.T) { // gofi:tests REQ-03.2
	e := newEnv()
	userID := e.confirmEmail(t, testEmail, testUsername)
	requireError(t, e.svc.ResendConfirmation(context.Background(), userID, ""), service.ErrContasEmailAlreadyConfirmed)
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado")
	}
}

func TestRequestEmailChange_SendsApprovalToCurrentEmailOnly(t *testing.T) { // gofi:tests REQ-11.1
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	mails := e.mailer.sent()
	if len(mails) != 1 || mails[0].To != testEmail || mails[0].Kind != entity.MailEmailChangeApproval {
		t.Fatalf("e-mails = %+v", mails)
	}
}

func TestApproveEmailChange_SendsConfirmationToNewEmail(t *testing.T) { // gofi:tests REQ-11.2
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	e.mailer.clear()
	e.clock.advance(50 * time.Minute)
	requireNoError(t, e.svc.ApproveEmailChange(context.Background(), approval))
	mails := e.mailer.sent()
	if len(mails) != 1 || mails[0].To != "ana@novo.com" || mails[0].Kind != entity.MailEmailChangeConfirmation {
		t.Fatalf("e-mails = %+v", mails)
	}
	if e.repo.userByEmail(testEmail) == nil {
		t.Fatal("o e-mail da conta deveria continuar o atual")
	}
}

func TestApproveEmailChange_LinkUsedTwiceIsRejected(t *testing.T) { // gofi:tests REQ-11.11
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	requireNoError(t, e.svc.ApproveEmailChange(context.Background(), approval))
	e.mailer.clear()
	requireError(t, e.svc.ApproveEmailChange(context.Background(), approval), service.ErrContasInvalidLink)
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado")
	}
}

func TestApproveEmailChange_NewEmailTakenBeforeApproval(t *testing.T) { // gofi:tests REQ-11.12
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	e.signUp(t, "ana@novo.com", "outra_ana")
	requireError(t, e.svc.ApproveEmailChange(context.Background(), approval), service.ErrContasEmailTaken)
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado")
	}
}

func TestConfirmEmailChange_CompletesChangeAndEndsAccess(t *testing.T) { // gofi:tests REQ-11.3
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	otherDevice := e.login(t, testEmail)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	resetToken := linkToken(t, e.mailer.last())
	requireNoError(t, e.svc.ApproveEmailChange(context.Background(), approval))
	confirmation := linkToken(t, e.mailer.last())
	e.clock.advance(23 * time.Hour)

	requireNoError(t, e.svc.ConfirmEmailChange(context.Background(), confirmation))

	user := e.repo.userByEmail("ana@novo.com")
	if user == nil || user.EmailConfirmedAt == nil {
		t.Fatalf("e-mail novo deveria estar na conta e confirmado: %+v", user)
	}
	if _, appErr := e.loginWith("ana@novo.com", testPassword, ""); appErr.Exists() {
		t.Fatalf("login com o e-mail novo recusado: %v", appErr)
	}
	_, appErr := e.loginWith(testEmail, testPassword, "")
	requireError(t, appErr, service.ErrContasInvalidCredentials)
	_, appErr = e.refresh(otherDevice.RefreshToken)
	requireError(t, appErr, service.ErrContasInvalidSession)
	requireError(t, e.svc.ResetPassword(context.Background(), resetToken, "outrasenha1"), service.ErrContasInvalidLink)
}

func TestApproveEmailChange_ExpiredApprovalIsRejected(t *testing.T) { // gofi:tests REQ-11.4
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	e.mailer.clear()
	e.clock.advance(61 * time.Minute)
	requireError(t, e.svc.ApproveEmailChange(context.Background(), approval), service.ErrContasInvalidLink)
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado ao novo endereço")
	}
}

func TestLogin_CurrentEmailValidWhileChangePending(t *testing.T) { // gofi:tests REQ-11.5
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	if _, appErr := e.loginWith(testEmail, testPassword, ""); appErr.Exists() {
		t.Fatalf("login recusado: %v", appErr)
	}
}

func TestRequestEmailChange_NewEmailInUseIsRejected(t *testing.T) { // gofi:tests REQ-11.6
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	e.signUp(t, "beto@exemplo.com", "beto")
	requireError(t, e.svc.RequestEmailChange(context.Background(), userID, "Beto@Exemplo.com", ""), service.ErrContasEmailTaken)
}

func TestConfirmEmailChange_NewEmailTakenBeforeConfirmation(t *testing.T) { // gofi:tests REQ-11.7
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	requireNoError(t, e.svc.ApproveEmailChange(context.Background(), linkToken(t, e.mailer.last())))
	confirmation := linkToken(t, e.mailer.last())
	e.signUp(t, "ana@novo.com", "outra_ana")

	requireError(t, e.svc.ConfirmEmailChange(context.Background(), confirmation), service.ErrContasEmailTaken)
	if e.repo.userByEmail(testEmail) == nil {
		t.Fatal("o e-mail da conta não deveria mudar")
	}
	requireError(t, e.svc.ConfirmEmailChange(context.Background(), confirmation), service.ErrContasEmailTaken)
}

func TestRequestEmailChange_NewRequestReplacesPending(t *testing.T) { // gofi:tests REQ-11.8
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	previous := linkToken(t, e.mailer.last())
	e.clock.advance(61 * time.Second)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@outro.com", ""))
	requireError(t, e.svc.ApproveEmailChange(context.Background(), previous), service.ErrContasInvalidLink)
}

func TestRequestEmailChange_OwnEmailIsRejected(t *testing.T) { // gofi:tests REQ-11.9
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireError(t, e.svc.RequestEmailChange(context.Background(), userID, "Ana@Exemplo.com", ""), service.ErrContasEmailTaken)
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado")
	}
}

func TestRequestEmailChange_MalformedEmailIsInvalid(t *testing.T) { // gofi:tests REQ-11.10
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@", ""), service.ErrContasInvalidInput)
}

func TestDeleteAccount_RemovesAccountSessionsAndLinks(t *testing.T) { // gofi:tests REQ-12.1
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	confirmation := linkToken(t, e.mailer.last())
	userID := e.repo.userByEmail(testEmail).ID
	first, second := e.login(t, testEmail), e.login(t, testEmail)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	resetToken := linkToken(t, e.mailer.last())

	requireNoError(t, e.svc.DeleteAccount(context.Background(), userID, testPassword))

	_, appErr = e.loginWith(testEmail, testPassword, "")
	requireError(t, appErr, service.ErrContasInvalidCredentials)
	for _, session := range []*service.SessionResult{first, second} {
		_, appErr = e.refresh(session.RefreshToken)
		requireError(t, appErr, service.ErrContasInvalidSession)
	}
	requireError(t, e.svc.ConfirmEmail(context.Background(), confirmation), service.ErrContasInvalidLink)
	requireError(t, e.svc.ResetPassword(context.Background(), resetToken, "outrasenha1"), service.ErrContasInvalidLink)
}

func TestDeleteAccount_WrongPasswordKeepsAccount(t *testing.T) { // gofi:tests REQ-12.2
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireError(t, e.svc.DeleteAccount(context.Background(), userID, "errada123"), service.ErrContasWrongPassword)
	if e.repo.userByEmail(testEmail) == nil {
		t.Fatal("a conta deveria continuar existindo")
	}
}

func TestSignUp_EmailAndUsernameFreeAfterDeletion(t *testing.T) { // gofi:tests REQ-12.3
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.DeleteAccount(context.Background(), userID, testPassword))
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
}

func TestRequestPasswordReset_SecondRequestInSameMinuteFromOtherIPIsRejected(t *testing.T) { // gofi:tests REQ-13.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, "10.0.0.1"))
	e.mailer.clear()
	e.clock.advance(30 * time.Second)
	requireError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, "10.0.0.2"), service.ErrContasTooManyRequests)
	if len(e.mailer.sent()) != 0 {
		t.Fatal("nenhum e-mail deveria ser enviado")
	}
}

func TestResendConfirmation_SixthRequestInHourIsRejected(t *testing.T) { // gofi:tests REQ-13.2
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	for range 5 {
		requireNoError(t, e.svc.ResendConfirmation(context.Background(), userID, ""))
		e.clock.advance(61 * time.Second)
	}
	requireError(t, e.svc.ResendConfirmation(context.Background(), userID, ""), service.ErrContasTooManyRequests)
}

func TestRequestPasswordReset_LimitByIP(t *testing.T) { // gofi:tests REQ-13.3
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.signUp(t, "beto@exemplo.com", "beto")
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, "10.0.0.1"))
	e.clock.advance(30 * time.Second)
	requireError(t, e.svc.RequestPasswordReset(context.Background(), "beto@exemplo.com", "10.0.0.1"), service.ErrContasTooManyRequests)
}

func TestRequestEmailChange_DifferentTypesDoNotAddUp(t *testing.T) { // gofi:tests REQ-13.4
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.ResendConfirmation(context.Background(), userID, ""))
	e.clock.advance(30 * time.Second)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
}

func TestRequestPasswordReset_RejectedRequestDoesNotCount(t *testing.T) { // gofi:tests REQ-13.6
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	e.clock.advance(40 * time.Second)
	requireError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""), service.ErrContasTooManyRequests)
	e.clock.advance(30 * time.Second)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
}

func TestRequestPasswordReset_AfterIntervalIsAccepted(t *testing.T) { // gofi:tests REQ-13.5
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	e.clock.advance(61 * time.Second)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
}

func TestSignUp_EmailServiceDownKeepsAccount(t *testing.T) { // gofi:tests REQ-14.1
	e := newEnv()
	e.mailer.setFail(true)
	result, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	if result.ConfirmationEmailSent {
		t.Fatal("a resposta deveria informar que o e-mail não foi enviado")
	}
	if e.repo.userByEmail(testEmail) == nil {
		t.Fatal("a conta deveria existir")
	}
}

func TestResendConfirmation_NewRequestAfterFailureSendsLink(t *testing.T) { // gofi:tests REQ-14.2
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	e.mailer.setFail(true)
	requireError(t, e.svc.ResendConfirmation(context.Background(), userID, ""), service.ErrContasEmailDeliveryFailed)
	e.mailer.setFail(false)
	e.clock.advance(time.Minute + time.Second)
	requireNoError(t, e.svc.ResendConfirmation(context.Background(), userID, ""))
	if len(e.mailer.sent()) != 1 {
		t.Fatalf("e-mails = %d, esperado 1", len(e.mailer.sent()))
	}
}

func TestApproveEmailChange_DeliveryFailureReportedAndLinkKept(t *testing.T) { // gofi:tests REQ-14.3
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	e.mailer.setFail(true)
	requireError(t, e.svc.ApproveEmailChange(context.Background(), approval), service.ErrContasEmailDeliveryFailed)
	e.mailer.setFail(false)
	e.mailer.clear()
	requireNoError(t, e.svc.ApproveEmailChange(context.Background(), approval))
}

func TestApproveEmailChange_RetryAfterFailureSendsConfirmation(t *testing.T) { // gofi:tests REQ-14.4
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	approval := linkToken(t, e.mailer.last())
	e.mailer.setFail(true)
	requireError(t, e.svc.ApproveEmailChange(context.Background(), approval), service.ErrContasEmailDeliveryFailed)
	e.mailer.setFail(false)
	e.mailer.clear()
	e.clock.advance(10 * time.Minute)
	requireNoError(t, e.svc.ApproveEmailChange(context.Background(), approval))
	mails := e.mailer.sent()
	if len(mails) != 1 || mails[0].To != "ana@novo.com" || mails[0].Kind != entity.MailEmailChangeConfirmation {
		t.Fatalf("e-mails = %+v", mails)
	}
}

func TestRequestPasswordReset_DeliveryFailureKeepsResponseAndCounts(t *testing.T) { // gofi:tests REQ-14.5
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	e.mailer.setFail(true)
	requireNoError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""))
	requireError(t, e.svc.RequestPasswordReset(context.Background(), testEmail, ""), service.ErrContasTooManyRequests)
}

func TestGrantAdmin_AppliesOnNextRefresh(t *testing.T) { // gofi:tests REQ-15.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	requireNoError(t, e.svc.GrantAdmin(context.Background(), testEmail))
	renewed, appErr := e.refresh(session.RefreshToken)
	requireNoError(t, appErr)
	if renewed.User.Role != entity.RoleAdmin || e.tokens.last().role != entity.RoleAdmin {
		t.Fatalf("papel = %s", renewed.User.Role)
	}
}

func TestRevokeAdmin_AppliesOnNextRefresh(t *testing.T) { // gofi:tests REQ-15.2
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.GrantAdmin(context.Background(), testEmail))
	session := e.login(t, testEmail)
	requireNoError(t, e.svc.RevokeAdmin(context.Background(), testEmail))
	renewed, appErr := e.refresh(session.RefreshToken)
	requireNoError(t, appErr)
	if renewed.User.Role != entity.RoleUser || e.tokens.last().role != entity.RoleUser {
		t.Fatalf("papel = %s", renewed.User.Role)
	}
}

func TestGrantAdmin_UnknownEmailFails(t *testing.T) { // gofi:tests REQ-15.3
	e := newEnv()
	requireError(t, e.svc.GrantAdmin(context.Background(), "ninguem@exemplo.com"), service.ErrContasUserNotFound)
}

func TestSignUp_AlwaysCreatesRegularUser(t *testing.T) { // gofi:tests REQ-15.4
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	if session.User.Role != entity.RoleUser {
		t.Fatalf("papel = %s", session.User.Role)
	}
}

func TestGetAccount_ShowsOwnDataAndPendingEmail(t *testing.T) { // gofi:tests REQ-16.1
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.RequestEmailChange(context.Background(), userID, "ana@novo.com", ""))
	account, appErr := e.svc.GetAccount(context.Background(), userID)
	requireNoError(t, appErr)
	if account.ID != userID || account.Email != testEmail || account.Username != testUsername ||
		account.DisplayName != "Ana" || account.Role != entity.RoleUser || account.EmailConfirmed {
		t.Fatalf("conta = %+v", account)
	}
	if account.PendingEmail == nil || *account.PendingEmail != "ana@novo.com" {
		t.Fatalf("troca pendente = %v", account.PendingEmail)
	}
}

func TestGetAccount_ReturnsOnlyOwnEmail(t *testing.T) { // gofi:tests REQ-16.3
	e := newEnv()
	anaID := e.signUp(t, testEmail, testUsername)
	e.signUp(t, "beto@exemplo.com", "beto")
	account, appErr := e.svc.GetAccount(context.Background(), anaID)
	requireNoError(t, appErr)
	if account.Email != testEmail || account.PendingEmail != nil {
		t.Fatalf("conta = %+v", account)
	}
}

func TestGetAccount_DeletedAccountIsUnauthenticated(t *testing.T) { // gofi:tests REQ-16.4
	e := newEnv()
	userID := e.signUp(t, testEmail, testUsername)
	requireNoError(t, e.svc.DeleteAccount(context.Background(), userID, testPassword))
	_, appErr := e.svc.GetAccount(context.Background(), userID)
	requireError(t, appErr, service.ErrContasUnauthenticated)
}

func TestGetAccount_UnknownUserIsUnauthenticated(t *testing.T) { // gofi:tests REQ-16.2
	e := newEnv()
	_, appErr := e.svc.GetAccount(context.Background(), newUnknownID())
	requireError(t, appErr, service.ErrContasUnauthenticated)
}

func TestSignUp_PasswordStoredOnlyAsHash(t *testing.T) { // gofi:tests REQ-17.1
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	if strings.Contains(e.repo.storedData(), testPassword) {
		t.Fatal("a senha apareceu nos dados guardados")
	}
	e.login(t, testEmail)
}

func TestSignUp_LinkTokenStoredOnlyAsHash(t *testing.T) { // gofi:tests REQ-17.2
	e := newEnv()
	_, appErr := e.svc.SignUp(context.Background(), validSignUp(testEmail, testUsername))
	requireNoError(t, appErr)
	token := linkToken(t, e.mailer.last())
	if strings.Contains(e.repo.storedData(), token) {
		t.Fatal("o token do link apareceu nos dados guardados")
	}
}

func TestLogin_RefreshCredentialStoredOnlyAsHash(t *testing.T) { // gofi:tests REQ-17.3
	e := newEnv()
	e.signUp(t, testEmail, testUsername)
	session := e.login(t, testEmail)
	if session.RefreshToken == "" || strings.Contains(e.repo.storedData(), session.RefreshToken) {
		t.Fatal("a credencial de renovação apareceu nos dados guardados")
	}
}
