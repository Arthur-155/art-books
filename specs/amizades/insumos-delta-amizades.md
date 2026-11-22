---
tipo: insumo-delta
contexto: amizades
versao: "1.0"
status: discovery
prd: prd/amizades/prd-amizades.md
keywords: [amizade, pedido de amizade, aceitar, recusar, cancelar pedido, desfazer amizade, amigos, username]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `amizades`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/amizades/prd-amizades.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Amizades

Quarto contexto do discovery. Depende de **Contas** (usuário e username). É a base do **Feed de amigos**, de **Interações** e do **Perfil** de amigo.

## 1. Problema e objetivo

O usuário quer acompanhar o que as pessoas que ele conhece estão lendo. A amizade é mútua: um envia o pedido e o outro aceita. Este contexto define como ela nasce, como termina e quem enxerga o quê por causa dela.

## 2. Escopo

**Dentro:** busca por username, pedido, aceite, recusa, cancelamento, desfazer amizade, lista de amigos.

**Fora (por enquanto):** bloqueio de usuário, convite por link, sugestão de amigos, grupos, limite de amigos.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Como pedir | Só por busca de username, como definido em Contas |
| D2 | Pedido pendente | Não expira. Quem enviou pode cancelar |
| D3 | Recusa | O remetente não é avisado. Ele pode reenviar depois de 7 dias |
| D4 | Desfazer | Qualquer lado pode desfazer, sem aviso ao outro |
| D5 | Bloqueio | Fora do primeiro corte, anotado como evolução |
| D6 | Lista de amigos | Privada: cada usuário vê só a própria |
| D7 | Limite | Sem limite de amigos |
| D8 | Perfil de quem não é amigo | Abre normalmente e mostra só o que é público |

## 4. Requisitos funcionais

- **REQ-01 Buscar usuário.** O usuário busca outro pelo username. A busca não revela o e-mail de ninguém.
- **REQ-02 Enviar pedido.** O usuário envia um pedido de amizade a outro. Não é possível enviar a si mesmo nem a quem já é amigo ou já tem pedido pendente.
- **REQ-03 Responder pedido.** Quem recebe vê os pedidos pendentes e aceita ou recusa. Ao aceitar, a amizade passa a valer para os dois lados.
- **REQ-04 Cancelar pedido.** Quem enviou cancela o pedido enquanto ele estiver pendente.
- **REQ-05 Recusa sem aviso.** Recusar não gera notificação ao remetente. O pedido some da lista de quem recusou.
- **REQ-06 Reenvio.** Depois de recusado, o remetente só consegue enviar um novo pedido após 7 dias.
- **REQ-07 Desfazer amizade.** Qualquer um dos dois desfaz a amizade, sem aviso ao outro.
- **REQ-08 Lista de amigos.** O usuário vê a própria lista de amigos e abre o perfil de cada um. Outros usuários não veem essa lista.

## 5. Regras de negócio

- A amizade é simétrica: existe para os dois lados ou para nenhum.
- Entre dois usuários existe no máximo um pedido pendente ou uma amizade.
- Se A e B enviarem pedido um ao outro ao mesmo tempo, os dois viram amigos, sem precisar de aceite.
- Pedido recusado não aparece para o remetente como recusado: ele continua vendo "pendente" até cancelar, e o reenvio só vale depois de 7 dias.
- Ao desfazer a amizade, o conteúdo que o outro publicou como "só para amigos" deixa de ser visto na hora, e o feed de amigos deixa de mostrá-lo.
- Ao apagar a conta, os pedidos e as amizades do usuário são removidos (D7 de Contas).
- Só usuário com e-mail confirmado envia pedido. *(assumido a partir de D3 de Contas, a confirmar)*

## 6. Dependências

- **Contas:** identificador, username e e-mail confirmado.
- **Feeds**, **Interações** e **Perfil** consultam este contexto para saber se dois usuários são amigos.

## 7. Critérios de aceite

- Enviar pedido a si mesmo é recusado.
- Enviar pedido a quem já é amigo ou já tem pedido pendente é recusado, com mensagem clara.
- Depois do aceite, os dois veem um ao outro na lista de amigos.
- Quem recusou um pedido não gera nenhuma notificação ao remetente.
- Um novo pedido ao mesmo usuário, feito 3 dias depois da recusa, é recusado. Feito 8 dias depois, é aceito.
- Dois pedidos cruzados entre A e B resultam em uma amizade, sem pedido pendente.
- Depois de desfeita a amizade, o feed de amigos de cada um deixa de mostrar o conteúdo do outro.
- A lista de amigos de A não é visível para B.

## 8. Questões em aberto

- Como o remetente vê um pedido recusado nos 7 dias de espera: como "pendente" (hoje) ou sem pedido algum.
- Se o reenvio de pedido tem um limite total de tentativas.
- Se o administrador pode ver ou desfazer amizades por denúncia. Hoje não existe denúncia nem moderação.

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Pedido recusado: o remetente continua vendo "pendente", para não revelar a recusa.
- Reenvio: até 3 pedidos para a mesma pessoa; depois disso, não reenvia.
- Administrador: não mexe em amizades; não há denúncia nem moderação no primeiro corte.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
