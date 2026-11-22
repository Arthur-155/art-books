---
tipo: prd
formato: 3
contexto: contas
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [cadastro, login, senha, confirmação de e-mail, recuperação de senha, troca de e-mail, exclusão de conta, administrador, resend]
atualizado: 2026-11-22
---
# PRD — contas — Cadastro e login de leitores

> Primeiro contexto da rede social de leitura. Os demais dependem dele pelo identificador do usuário, pelo username, pelo estado "e-mail confirmado" e pelo papel de administrador.

## 1. Problema

Quem lê livros não tem hoje um lugar para registrar o que entendeu de cada capítulo. Sem uma identidade própria, ninguém publica resumos, tem perfil ou se conecta a amigos. A conta precisa ser segura o bastante para que o dono dela seja de fato quem age nela, porque o que se publica carrega o nome da pessoa.

## 2. Objetivo e métricas de sucesso

Permitir que um leitor crie uma conta, entre e recupere o acesso sozinho, com e-mail e senha, de forma segura.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Cadastros concluídos sobre cadastros iniciados | n/a (produto novo) | 70% | Eventos de cadastro |
| Contas com e-mail confirmado em 24 h | n/a | 80% | Registro de confirmações |
| Pedidos de recuperação de senha concluídos | n/a | 85% | Registro de tokens usados |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora | Usuária comum | Entrar rápido e ter a conta protegida | Não há produto para registrar leituras |
| Administrador | Mantém o catálogo | Papel com poderes sobre gêneros | Não existe esse papel |

## 4. Escopo

### Dentro
- Cadastro, login e logout com e-mail e senha, na web responsiva.
- Nome de exibição e `@username` único; o e-mail nunca aparece para outros usuários.
- Confirmação de e-mail por link, enviado pelo Resend; sem confirmar, o usuário não publica resumos.
- Recuperação de senha por link de uso único e prazo curto, enviado pelo Resend.
- Senha com mínimo de 8 caracteres e limite de tentativas por e-mail e por IP.
- Sessão longa com renovação; trocar a senha encerra todos os dispositivos.
- Troca de e-mail aprovada a partir do e-mail antigo e confirmada no novo.
- Exclusão de conta, removendo tudo do usuário.
- Aceite de termos e privacidade, com idade mínima de 13 anos.
- Dois papéis: usuário e administrador, o segundo definido direto na base.

### Fora
- Login social: não é prioridade agora.
- Autenticação em dois fatores: fica para uma evolução.
- App mobile: a plataforma inicial é a web responsiva.

## 5. Jornadas

### J1 — Criar conta e confirmar e-mail
1. A pessoa informa e-mail, senha, nome, username e aceita os termos.
2. Recebe o link de confirmação por e-mail.
3. Confirma e passa a poder publicar resumos.

### J2 — Recuperar a senha
1. A pessoa pede a recuperação com o e-mail.
2. Recebe um link de uso único.
3. Define uma nova senha e todos os dispositivos são desconectados.

### J3 — Trocar o e-mail
1. A pessoa informa o novo e-mail.
2. Aprova o pedido a partir do e-mail antigo.
3. Confirma o novo endereço, e a troca passa a valer.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | Envio de e-mail pelo Resend | Falha no envio impede confirmação e recuperação | Reenvio e erro claro ao usuário |
| Restrição | LGPD: exclusão remove os dados do usuário | Afeta resumos, comentários e republicações | Regra de cascata definida nos outros contextos |
| Premissa | O primeiro administrador é criado direto na base | Sem tela de cadastro de administrador | Procedimento operacional |
| Risco | Abuso de login por tentativa e erro | Contas invadidas | Limite de tentativas por e-mail e IP |
| Risco | Troca de e-mail por quem invadiu a conta | Perda da conta | Aprovação a partir do e-mail antigo |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Prazos: confirmação de e-mail em 24 h, recuperação de senha em 30 min, aprovação da troca de e-mail em 1 h; sessão de 30 dias com renovação.
- Bloqueio de login: 5 tentativas erradas bloqueiam por 15 min, por e-mail e por IP.
- Exclusão de conta: imediata, com confirmação por senha.
- Primeiro administrador: criado por um comando operacional na base, sem tela no app; o papel é concedido e retirado do mesmo modo.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
