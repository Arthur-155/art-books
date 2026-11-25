//gofi:context contas
package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/Arthur-155/art-books/domain/contas/entity"
)

const (
	DefaultEndpoint = "https://api.resend.com/emails"
	sendTimeout     = 10 * time.Second
)

type Config struct {
	APIKey   string
	From     string
	Endpoint string
}

// ResendMailer envia os e-mails pela API HTTP do Resend, de forma síncrona.
type ResendMailer struct {
	cfg    Config
	client *http.Client
}

func NewResendMailer(cfg Config) *ResendMailer {
	if cfg.Endpoint == "" {
		cfg.Endpoint = DefaultEndpoint
	}
	return &ResendMailer{cfg: cfg, client: &http.Client{Timeout: sendTimeout}}
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type mailTemplate struct {
	subject string
	intro   string
	action  string
}

var templates = map[entity.MailKind]mailTemplate{
	entity.MailEmailConfirmation:       {"Confirme seu e-mail", "Confirme seu e-mail para publicar resumos.", "Confirmar e-mail"},
	entity.MailPasswordReset:           {"Redefinição de senha", "Use o link para criar uma nova senha. Ele vale por 30 minutos.", "Redefinir senha"},
	entity.MailEmailChangeApproval:     {"Aprove a troca de e-mail", "Recebemos um pedido para trocar o e-mail da sua conta. O link vale por 1 hora.", "Aprovar troca"},
	entity.MailEmailChangeConfirmation: {"Confirme o novo e-mail", "Confirme este endereço para concluir a troca de e-mail.", "Confirmar novo e-mail"},
}

// Send falha se o Resend recusar ou não responder em 10 s; o erro nunca leva o endereço nem o link.
//
// gofi:implements REQ-14
func (m *ResendMailer) Send(ctx context.Context, mail entity.Mail) error {
	template, ok := templates[mail.Kind]
	if !ok {
		return fmt.Errorf("resend: unknown mail kind %q", mail.Kind)
	}
	body, err := json.Marshal(resendPayload{
		From: m.cfg.From, To: []string{mail.To}, Subject: template.subject, HTML: renderHTML(template, mail.Link),
	})
	if err != nil {
		return err
	}
	return m.post(ctx, body)
}

func (m *ResendMailer) post(ctx context.Context, body []byte) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, m.cfg.Endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("resend: build request failed")
	}
	request.Header.Set("Authorization", "Bearer "+m.cfg.APIKey)
	request.Header.Set("Content-Type", "application/json")
	response, err := m.client.Do(request)
	if err != nil {
		return fmt.Errorf("resend: request failed")
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("resend: unexpected status %d", response.StatusCode)
	}
	return nil
}

func renderHTML(template mailTemplate, link string) string {
	escaped := html.EscapeString(link)
	return fmt.Sprintf(`<p>%s</p><p><a href="%s">%s</a></p><p>%s</p>`, template.intro, escaped, template.action, escaped)
}
