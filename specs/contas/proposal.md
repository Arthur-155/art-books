---
tipo: mudanca
formato: 3
contextos: [contas]
versao: "1.0"
status: aprovado
depende_de: []
superficies: [backend]
keywords: [cadastro, login, sessão, confirmação de e-mail, recuperação de senha, troca de senha, troca de e-mail, exclusão de conta, limite de tentativas, administrador, resend, argon2id]
atualizado: 2026-11-24
aprovacao_versao: "1.0"
aprovacao_sha: "sha256:80af203f60056de6e6fbb17a3fba9dd258ae8bd58f0e6c8ec15af70e340cf2d5"
---
# Mudança — contas

## Por quê

Sem uma identidade própria, nenhum leitor publica resumos, tem perfil ou se
conecta a amigos, e todo outro contexto depende do id do usuário, do username,
do e-mail confirmado e do papel (`prd/contas/prd-contas.md` §1). Esta mudança
cria o contexto: cadastro, login e recuperação de acesso com e-mail e senha,
seguros o bastante para que o dono da conta seja quem age nela.

## O que muda

- contas: REQ-01 novo (cadastro com e-mail, senha, nome de exibição, username, termos e idade mínima)
- contas: REQ-02 novo (confirmação do e-mail por link de 24 h e uso único)
- contas: REQ-03 novo (reenvio da confirmação, invalidando o link anterior)
- contas: REQ-04 novo (a sessão informa se o e-mail está confirmado)
- contas: REQ-05 novo (login com e-mail e senha, sem revelar se o e-mail existe)
- contas: REQ-06 novo (bloqueio de 15 min após 5 tentativas erradas, por e-mail e por IP)
- contas: REQ-07 novo (sessão de 30 dias renovada com o uso)
- contas: REQ-08 novo (logout encerra só o dispositivo atual)
- contas: REQ-09 novo (recuperação de senha por link de 30 min e uso único, encerrando todas as sessões)
- contas: REQ-10 novo (troca de senha com a senha atual, encerrando todas as sessões)
- contas: REQ-11 novo (troca de e-mail aprovada no e-mail antigo e confirmada no novo)
- contas: REQ-12 novo (exclusão imediata da conta com a senha)
- contas: REQ-13 novo (limite de 1 pedido por minuto e 5 por hora para os envios de e-mail)
- contas: REQ-14 novo (falha no envio do e-mail vira erro claro e permite novo pedido)
- contas: REQ-15 novo (papéis usuário e administrador; o administrador só por comando operacional)
- contas: REQ-16 novo (consulta da própria conta; o e-mail nunca é exposto a outro usuário)
- contas: REQ-17 novo (senha e tokens de e-mail nunca guardados em texto)

## Impacto

O código ainda não tem nenhum contexto: só `backend/art-books/main.go`. A
mudança cria `backend/domain/contas/`, as tabelas `"user"`, `user_session`,
`email_token`, `login_lock` e `email_request`, o subcomando
`art-books admin grant|revoke <e-mail>`, um expurgo periódico no binário e a
chamada à API do Resend, com três variáveis novas (`APP_BASE_URL`,
`RESEND_API_KEY`, `RESEND_FROM_EMAIL`).

`catalogo` e `recomendacoes` declaram `depende_de: [contas]`: usam a tabela
`"user"` como FK (o catálogo com `ON DELETE SET NULL`) e leem do token o id, o
papel e o indicador de e-mail confirmado. `perfil`, `amizades` e
`leitura-e-resumos` consomem só o id, o username e o e-mail confirmado; o que
cada um apaga na exclusão da conta é da FK de cada um. As migrações de
`catalogo` foram desenhadas a partir de `0001`: se o serviço tem uma sequência
única de migrações, as dele passam a começar depois das de `contas`.

## Fora desta mudança

- Telas web: as páginas que os links de e-mail abrem (caminhos fixados no
  design) e os formulários ficam para a mudança de interface de contas.
- Login social, autenticação em dois fatores e app mobile (PRD §4 Fora).
- Perfil público, foto e busca de usuário por username: são de `perfil` e
  `amizades`.
- Período de arrependimento na exclusão: a exclusão é imediata (PRD §7).

## Dúvidas em aberto

Nenhuma.
