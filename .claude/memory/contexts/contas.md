---
formato: memoria
contexto: contas
versao: "1.0"
status: auditado
atualizado: 2026-11-24
keywords: [cadastro, login, sessão, confirmação de e-mail, recuperação de senha, troca de e-mail, exclusão de conta, limite de tentativas, administrador, resend]
ressalvas: 2
---
# contas

## Estado atual

Mudança `contas` implementada em `backend/domain/contas/` (REQ-01 a REQ-17), `build` e `static` verdes; testes ainda por escrever (gofi-test).

Camadas: `entity` (User, Session, EmailToken, normalização e validações puras), `repository` (arquivo único, SQL por `sqln`, `RotateSession` em transação, restrições `uq_user_email` e `uq_user_username` traduzidas em erros), `service` (`contas_service.go` e `auth_service.go` sob a interface `ContasService`, `errors.go`, `Mailer` e `AccessTokenIssuer` como interfaces), `adapter` (JWT do SDK com `role` e `emailConfirmed`), `mailer` (Resend síncrono, 10 s), `handler` (`contas_handler.go`, `auth_handler.go`, `middleware.go` com cookies `access_token` e `refresh_token`), `worker` (`PurgeWorker`, `gofi.Runner`, a cada 1 h).

Binário `backend/art-books/`: `main.go`, `config.go` (`APP_BASE_URL`, `RESEND_API_KEY`, `RESEND_FROM_EMAIL` obrigatórias; `JWT_SECRET` ≥ 32 bytes), `wire.go`, `admin.go` (`art-books admin grant|revoke <e-mail>`). Migrações `backend/.migrations/0001` a `0004` com `up`/`down`. `.env` na raiz com placeholders.

Decisões de implementação: senha pelo `iam/provider/password` do SDK (argon2id, mesmos parâmetros do design); `uuid.NewV7`; falha no envio da aprovação ou conflito na confirmação da troca de e-mail devolve o token consumido ao banco (`ReplaceEmailToken`); a troca e a recuperação de senha só encerram sessões, e a confirmação da troca de e-mail apaga também os links pendentes (REQ-11.3); o `go.mod` passou a declarar o módulo `github.com/Arthur-155/art-books`.

As migrações de contas são 0001 a 0004; as de catalogo começam depois.

## Histórico de versões
