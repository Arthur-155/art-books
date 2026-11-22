---
tipo: insumo-delta
contexto: contas
versao: "1.0"
status: discovery
prd: prd/contas/prd-contas.md
keywords: [cadastro, login, senha, confirmação de e-mail, recuperação de senha, troca de e-mail, exclusão de conta, administrador, resend]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `contas`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/contas/prd-contas.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Contas

Primeiro contexto do discovery da rede social de leitura. Todos os outros dependem dele.

## 1. Problema e objetivo

Leitores precisam de uma identidade própria para registrar resumos, ter perfil e se conectar a amigos. Este contexto entrega cadastro e login com e-mail e senha, com segurança suficiente para que o dono da conta seja de fato quem age nela.

## 2. Escopo

**Dentro:** cadastro, confirmação de e-mail, login, logout, recuperação de senha, troca de senha, troca de e-mail, exclusão de conta.

**Fora (por enquanto):** login social, autenticação em dois fatores, app mobile (a plataforma é web responsiva).

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Plataforma | Web responsiva |
| D2 | Identidade pública | Nome de exibição e `@username` único. Amigos são encontrados por username. O e-mail nunca aparece para outros usuários |
| D3 | Confirmação de e-mail | Link enviado pelo Resend. Sem confirmar, o usuário não publica resumos |
| D4 | Recuperação de senha | Link enviado pelo Resend, de uso único e prazo curto, que leva à definição de uma nova senha |
| D5 | Senha | Mínimo de 8 caracteres. Limite de tentativas por e-mail e por IP |
| D6 | Sessão | Sessão longa com renovação. Ao trocar a senha, todos os dispositivos são deslogados |
| D7 | Exclusão de conta | Permitida. Remove tudo do usuário: resumos, comentários e republicações |
| D8 | Termos e idade | Aceite obrigatório de termos e privacidade no cadastro. Idade mínima de 13 anos |
| D9 | Troca de e-mail | O usuário pede a troca e aprova a partir do e-mail **antigo**, para provar que é o dono da conta. O novo e-mail também precisa ser confirmado |
| D10 | Papéis | Existem dois: usuário e administrador. O administrador é definido fora do app, direto na base, sem tela de cadastro. O Catálogo usa o papel para corrigir gêneros e revisar a lista |

## 4. Requisitos funcionais

- **REQ-01 Cadastro.** O usuário informa e-mail, senha, nome de exibição e username, e aceita os termos. O e-mail e o username devem ser únicos. O usuário declara ter 13 anos ou mais.
- **REQ-02 Confirmação de e-mail.** Após o cadastro, o sistema envia um link de confirmação pelo Resend. O usuário sem e-mail confirmado entra e navega, mas não publica resumos. O link pode ser reenviado.
- **REQ-03 Login.** E-mail e senha. A mensagem de erro não revela se o e-mail existe.
- **REQ-04 Limite de tentativas.** Tentativas erradas de login são limitadas por e-mail e por IP, com bloqueio temporário.
- **REQ-05 Sessão.** A sessão dura muito tempo e se renova com o uso. O usuário pode sair no dispositivo atual.
- **REQ-06 Recuperação de senha.** O usuário informa o e-mail e recebe um link pelo Resend. O link é de uso único e expira em pouco tempo. A resposta da tela é a mesma exista ou não a conta. Ao concluir, todas as sessões são encerradas.
- **REQ-07 Troca de senha.** O usuário logado troca a senha informando a atual. Ao concluir, todas as sessões são encerradas.
- **REQ-08 Troca de e-mail.** O usuário informa o novo e-mail. O sistema envia um pedido de aprovação ao e-mail **antigo**. Só depois da aprovação ele envia a confirmação ao novo e-mail. A troca vale quando o novo e-mail é confirmado. Os links são de uso único e expiram.
- **REQ-09 Exclusão de conta.** O usuário confirma com a senha. O sistema remove todos os dados dele (D7) e encerra as sessões.

## 5. Regras de negócio

- Username e e-mail são únicos sem diferenciar maiúsculas de minúsculas.
- A senha é guardada apenas com hash forte, nunca em texto.
- Os tokens de e-mail (confirmação, recuperação, troca) são de uso único, com prazo, e guardados com hash.
- Enquanto a troca de e-mail estiver pendente, o e-mail antigo continua valendo para login.

## 6. Fora do escopo e dependências

- **Resend** é a dependência externa de envio. Se o envio falhar, o usuário vê um erro claro e pode pedir o reenvio.
- Os contextos seguintes dependem de Contas apenas pelo identificador do usuário, pelo username e pelo estado "e-mail confirmado".

## 7. Critérios de aceite

- Cadastro com e-mail ou username repetido é recusado com mensagem clara.
- Usuário sem e-mail confirmado não consegue publicar um resumo.
- O link de recuperação funciona uma única vez e deixa de valer após o prazo.
- Depois da troca de senha, uma sessão aberta em outro dispositivo deixa de funcionar.
- A troca de e-mail não avança sem a aprovação feita pelo e-mail antigo.
- Depois da exclusão, a conta e os dados dela não aparecem em nenhum lugar do sistema.

## 8. Questões em aberto

- Prazo exato dos links de e-mail e da sessão.
- Quantas tentativas erradas disparam o bloqueio, e por quanto tempo.
- Se a exclusão é imediata ou tem um período de arrependimento.
- Como o primeiro administrador é criado e como se concede ou retira o papel depois.

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Prazos: confirmação de e-mail em 24 h, recuperação de senha em 30 min, aprovação da troca de e-mail em 1 h; sessão de 30 dias com renovação.
- Bloqueio de login: 5 tentativas erradas bloqueiam por 15 min, por e-mail e por IP.
- Exclusão de conta: imediata, com confirmação por senha.
- Primeiro administrador: criado por um comando operacional na base, sem tela no app; o papel é concedido e retirado do mesmo modo.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
