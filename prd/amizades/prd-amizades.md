---
tipo: prd
formato: 3
contexto: amizades
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [amizade, pedido de amizade, aceitar, recusar, cancelar pedido, desfazer amizade, amigos, username]
atualizado: 2026-11-22
---
# PRD — amizades — Conexão mútua entre leitores

> Quarto contexto. Depende de **contas**. É a base do feed de amigos em **feeds**, de **interacoes** e do perfil de amigo em **perfil**.

## 1. Problema

O leitor quer acompanhar o que pessoas conhecidas estão lendo e resumindo, e não só descobrir desconhecidos. Sem uma relação entre usuários, não existe feed de amigos nem uma forma de ver o que quem se conhece compartilhou.

## 2. Objetivo e métricas de sucesso

Permitir que dois leitores se conectem de forma mútua, e que cada um controle quando essa conexão termina.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Usuários com ao menos 1 amigo em 7 dias | n/a (produto novo) | 50% | Amizades por usuário |
| Pedidos aceitos sobre pedidos enviados | n/a | 60% | Registro de pedidos |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora | Quer seguir o que amigos leem | Encontrar e adicionar pessoas | Não há como se conectar |
| Leitora que recebe pedido | Decide quem entra | Aceitar ou recusar sem constrangimento | Não há pedido |

## 4. Escopo

### Dentro
- Busca de usuário por username.
- Pedido de amizade, aceite, recusa e cancelamento pelo remetente; a amizade é mútua.
- Pedido pendente que não expira.
- Recusa sem aviso ao remetente; o reenvio só vale depois de 7 dias.
- Desfazer a amizade por qualquer lado, sem aviso ao outro.
- Lista de amigos privada: cada usuário vê só a própria.
- Sem limite de amigos.

### Fora
- Bloqueio de usuário: fica como evolução.
- Convite por link, sugestão de amigos e grupos: fora do primeiro corte.
- Denúncia e moderação: não existem ainda.

## 5. Jornadas

### J1 — Fazer um amigo
1. A leitora busca um username e envia o pedido.
2. A outra pessoa vê o pedido e aceita.
3. As duas passam a se ver na lista de amigos e no feed de amigos.

### J2 — Recusar ou desfazer
1. A pessoa recusa um pedido, ou desfaz uma amizade existente.
2. O outro lado não é avisado.
3. O conteúdo deixa de aparecer no feed de amigos de cada um.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Premissa | A amizade é simétrica: vale para os dois ou para nenhum | Simplifica as regras do feed | Regra explícita |
| Premissa | Pedidos cruzados entre A e B viram amizade sem aceite | Evita impasse | Regra explícita |
| Restrição | Só quem tem e-mail confirmado envia pedido | Reduz abuso | Regra de contas |
| Risco | Assédio por pedidos repetidos | Incômodo ao destinatário | Espera de 7 dias após recusa |
| Risco | Sem bloqueio nem denúncia | Pouca defesa contra abuso | Evolução prevista |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Pedido recusado: o remetente continua vendo "pendente", para não revelar a recusa.
- Reenvio: até 3 pedidos para a mesma pessoa; depois disso, não reenvia.
- Administrador: não mexe em amizades; não há denúncia nem moderação no primeiro corte.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
