---
tipo: prd
formato: 3
contexto: interacoes
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [curtir, curtida, comentar, comentário, republicar, republicação, feed de amigos, interação]
atualizado: 2026-11-22
---
# PRD — interacoes — Curtir, comentar e republicar resumos

> Quinto contexto. Depende de **leitura-e-resumos** e **amizades**. É consumido por **feeds**.

## 1. Problema

Ler o resumo de alguém só vira conversa quando o leitor pode reagir. Sem curtir, comentar e repassar, o feed é uma vitrine muda e o autor não sabe se o que escreveu ajudou alguém.

## 2. Objetivo e métricas de sucesso

Dar ao leitor formas simples de reagir a um resumo e de indicá-lo aos seus amigos.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Resumos abertos que recebem curtida | n/a (produto novo) | 15% | Curtidas sobre aberturas |
| Resumos públicos com ao menos 1 comentário | n/a | 20% | Comentários por resumo |
| Resumos republicados ao menos uma vez | n/a | 5% | Republicações por resumo |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora | Reage a resumos | Curtir, comentar e indicar a amigos | Não há reação |
| Autora do resumo | Recebe reações | Saber o que gerou interesse e moderar seu espaço | Não há retorno |

## 4. Escopo

### Dentro
- Curtida simples, com contador público; só o autor do resumo vê quem curtiu.
- Comentário no resumo do livro, de até 1.000 caracteres, editável e apagável pelo autor, sem respostas encadeadas.
- Moderação: o autor do resumo também apaga comentários feitos nele.
- Republicação: colocar o resumo público de outra pessoa no feed de amigos dos seus amigos, com crédito ao autor original e texto próprio opcional de até 500 caracteres.
- Quando o resumo vira privado, as curtidas, os comentários e as republicações somem e voltam se ele voltar a ser público.
- Quando um resumo ou uma conta é apagado, a republicação some e os comentários e curtidas de quem apagou a conta são removidos.

### Fora
- Reações diferentes e respostas encadeadas: fora do primeiro corte.
- Notificações: o autor vê só os contadores.
- Denúncia e moderação por administrador: não existem ainda.

## 5. Jornadas

### J1 — Reagir a um resumo
1. A leitora abre um resumo no feed.
2. Curte e escreve um comentário.
3. O contador e a conversa aparecem para quem tem acesso ao resumo.

### J2 — Indicar a amigos
1. A leitora republica um resumo público, com um texto próprio opcional.
2. Os amigos dela veem a republicação no feed de amigos, com o autor original em destaque.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | Só resumos públicos são republicados | Republicar não amplia a visibilidade | Regra explícita |
| Restrição | A republicação aparece só no feed de amigos de quem republicou | Não vai ao feed geral | Regra explícita |
| Premissa | Curtir ou comentar uma republicação age sobre o resumo original | Um só contador e uma só conversa | Regra explícita |
| Premissa | Cada usuário republica o mesmo resumo uma única vez | Evita repetição no feed | Regra explícita |
| Risco | Comentários ofensivos | Ambiente hostil | O autor do resumo apaga; moderação formal fica para depois |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Resumo que volta a ser público: curtidas, comentários e republicações reaparecem.
- Quem republica não vê para quantos amigos a republicação chegou.
- Comentário apagado: some por completo, sem marca.
- Limite de 20 comentários por usuário em um mesmo resumo.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
