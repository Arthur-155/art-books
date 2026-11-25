---
formato: memoria
contexto: contas
versao: "1.0"
status: spec
atualizado: 2026-11-24
keywords: [cadastro, login, sessão, confirmação de e-mail, recuperação de senha, troca de e-mail, exclusão de conta, limite de tentativas, administrador, resend]
---
# contas

## Estado atual

Gênese escrita em `specs/contas/` (mudança `contas`): proposta, `deltas/contas.md` com REQ-01 a REQ-17 e `design.md` (Go, serviço `art-books`, PostgreSQL, tabelas `"user"`, `user_session`, `email_token`, `login_lock` e `email_request`, Resend atrás de `Mailer`, expurgo no binário, subcomando `art-books admin grant|revoke`). Status `proposto`, sem aprovação, sem marcadores em aberto; `gofi spec check contas` com 0 erros e duas rodadas de `gofi review` já corrigidas.

Decisões aplicadas: sessão em cookies HTTP-only (ADR-06); encerrar sessões impede só a renovação e o token de acesso vale até vencer (ADR-05); recuperação de senha dá a mesma resposta mesmo se o envio falhar; falha no envio da aprovação da troca preserva o link; cadastro e aprovação ficam fora do limite de pedidos de e-mail; senha errada na troca de senha e na exclusão fica fora do bloqueio de login; limites em tabela do PostgreSQL; argon2id; username de 3 a 30 caracteres em `a-z`, `0-9` e `_`, nome de exibição de 1 a 60; papel de administrador só pelo subcomando; variáveis `APP_BASE_URL`, `RESEND_API_KEY` e `RESEND_FROM_EMAIL`.

Próximo passo: a pessoa decide se aprova (`gofi change approve contas`) ou roda a revisão de novo. Em seguida o loop implementa (gofi-eng) e audita (gofi-qa). As migrações de contas são 0001 a 0004; as de catalogo precisam começar depois.

## Histórico de versões
