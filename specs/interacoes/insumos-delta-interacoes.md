---
tipo: insumo-delta
contexto: interacoes
versao: "1.0"
status: discovery
prd: prd/interacoes/prd-interacoes.md
keywords: [curtir, curtida, comentar, comentário, republicar, republicação, feed de amigos, interação]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `interacoes`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/interacoes/prd-interacoes.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Interações

Quinto contexto do discovery. Depende de **Leitura e resumos** (resumo e visibilidade) e de **Amizades**. É consumido pelos **Feeds**.

## 1. Problema e objetivo

Ler o resumo de alguém só vira conversa se o leitor puder reagir. Este contexto define curtir, comentar e republicar um resumo, e o que acontece com essas ações quando o resumo, a visibilidade ou a conta mudam.

## 2. Escopo

**Dentro:** curtida, comentário, republicação, remoção em cascata.

**Fora (por enquanto):** reações diferentes, respostas encadeadas a comentários, notificações, denúncia e moderação.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Curtir | Curtida simples, com contador público. O autor do resumo vê quem curtiu |
| D2 | Onde comentar | Só no resumo do livro, não nos capítulos |
| D3 | Comentário | Máximo de 1.000 caracteres. O autor do comentário edita e apaga. Sem respostas encadeadas |
| D4 | Quem apaga comentário | O autor do comentário e o autor do resumo |
| D5 | Republicar | Colocar o resumo de outra pessoa no feed dos **seus amigos**, com crédito ao autor original. Texto próprio opcional de até 500 caracteres |
| D6 | O que se republica | Só resumos públicos |
| D7 | Onde aparece a republicação | Só no feed de amigos de quem republicou. Não vai ao feed geral |
| D8 | Resumo vira privado | Curtidas, comentários e republicações somem de todos os lugares e voltam se o resumo voltar a ser público |
| D9 | Resumo ou conta apagados | A republicação some. Os comentários e as curtidas de quem apagou a conta são removidos |
| D10 | Notificações | Fora do primeiro corte. O autor vê só os contadores |
| D11 | Denúncia | Fora do primeiro corte |

## 4. Requisitos funcionais

- **REQ-01 Curtir.** O usuário curte e descurte um resumo visível para ele. Cada usuário curte uma vez por resumo.
- **REQ-02 Contador e lista.** Todos veem o total de curtidas. Só o autor do resumo vê quem curtiu.
- **REQ-03 Comentar.** O usuário comenta o resumo do livro, com até 1.000 caracteres, e edita ou apaga o próprio comentário.
- **REQ-04 Moderar no próprio resumo.** O autor do resumo apaga qualquer comentário feito nele.
- **REQ-05 Republicar.** O usuário republica o resumo público de outra pessoa. O resumo aparece no feed de amigos de cada amigo dele, com o autor original em destaque e o texto próprio, se houver.
- **REQ-06 Desfazer republicação.** Quem republicou remove a republicação a qualquer momento.
- **REQ-07 Cascata.** Mudanças no resumo ou nas contas aplicam as regras D8 e D9 de forma automática.

## 5. Regras de negócio

- Só interage quem tem acesso ao resumo: público para todos, privado só para o autor.
- Curtir e comentar uma republicação age sobre o resumo original, com um único contador e uma única conversa.
- O usuário republica um mesmo resumo uma única vez. Republicar uma republicação republica o original, nunca forma corrente.
- O usuário não republica o próprio resumo.
- Republicar não torna o resumo mais visível do que ele já era: os amigos de quem republicou veem apenas o que já é público.
- Se o autor desfaz a amizade com quem republicou, a republicação continua valendo, porque o resumo original é público.
- Comentar, curtir e republicar exigem e-mail confirmado. *(assumido a partir de D3 de Contas, a confirmar)*

## 6. Dependências

- **Leitura e resumos:** resumo, visibilidade e autor.
- **Amizades:** para saber quem são os amigos de quem republica.
- **Feeds** exibem republicações, contadores e comentários daqui.

## 7. Critérios de aceite

- Um usuário não curte duas vezes o mesmo resumo.
- O contador de curtidas é igual para todos, mas só o autor vê a lista de quem curtiu.
- Um comentário com 1.001 caracteres é recusado.
- O autor do resumo apaga o comentário de outra pessoa nele. Um terceiro não apaga.
- Um resumo privado não pode ser republicado.
- Uma republicação aparece no feed de amigos dos amigos de quem republicou e não aparece no feed geral.
- Ao tornar o resumo privado, as curtidas, os comentários e as republicações somem. Ao voltar a público, reaparecem.
- Ao apagar a conta, os comentários e as curtidas dela somem dos resumos dos outros.
- Curtir uma republicação aumenta o contador do resumo original.

## 8. Questões em aberto

- Se quem republica vê quantos amigos receberam a republicação.
- Se um comentário apagado deixa marca ("comentário removido") ou some por completo.
- Se existe limite de comentários por usuário em um mesmo resumo, para evitar abuso.
- O comportamento de reaparecer após voltar a público (D8) pode surpreender o autor: confirmar se isso é desejado ou se é preferível apagar de vez.

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Resumo que volta a ser público: curtidas, comentários e republicações reaparecem.
- Quem republica não vê para quantos amigos a republicação chegou.
- Comentário apagado: some por completo, sem marca.
- Limite de 20 comentários por usuário em um mesmo resumo.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
