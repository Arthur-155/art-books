---
tipo: insumo-delta
contexto: feeds
versao: "1.0"
status: discovery
prd: prd/feeds/prd-feeds.md
keywords: [feed, feed geral, feed de amigos, resumo, republicação, filtro por gênero, ver mais, feed vazio]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `feeds`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/feeds/prd-feeds.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Feeds

Sexto contexto do discovery. Depende de **Leitura e resumos**, **Amizades**, **Interações** e do gênero do **Catálogo**.

## 1. Problema e objetivo

O usuário quer descobrir leituras pelo que outros leitores resumiram. O feed geral mostra o que é público de todos. O feed de amigos mostra só o que os amigos publicaram ou republicaram. As duas telas têm as mesmas funções.

## 2. Escopo

**Dentro:** feed geral, feed de amigos, item do feed, resumo abreviado com "ver mais", abertura do resumo, filtro por gênero, feed vazio.

**Fora (por enquanto):** ordenação por relevância, busca por livro ou autor, resumos visíveis só para amigos, notificações.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Ordem | Cronológica, do mais recente para o mais antigo |
| D2 | Item do feed | Mostra tudo de forma abreviada: título, autor, gênero, nome e username de quem escreveu, resumo abreviado e contadores. O resumo traz um "ver mais" para ler o texto inteiro no próprio item |
| D3 | Data | Vale a data em que o resumo ficou público. Editar não reordena |
| D4 | Paginação | Rolagem infinita, 20 itens por vez |
| D5 | Filtro | Por gênero. Sem busca por livro ou autor |
| D6 | Feed geral | Todos os resumos públicos, sem republicações |
| D7 | Feed de amigos | Resumos públicos dos amigos e as republicações deles. Não existe visibilidade "só amigos" |
| D8 | Seus resumos | Não aparecem no seu feed. Ficam na estante e no perfil |
| D9 | Abrir o resumo | Mesma tela nos dois feeds, com o resumo do livro, os capítulos, curtidas, comentários e republicar |
| D10 | Feed vazio | Mensagem convidando o usuário a buscar amigos por username |
| D11 | Gênero sem definição | Livros de "gênero não identificado" ou "não verificado" aparecem, sem a etiqueta de gênero |

## 4. Requisitos funcionais

- **REQ-01 Feed geral.** Lista os resumos públicos de todos os usuários, exceto os do próprio usuário, do mais recente para o mais antigo.
- **REQ-02 Feed de amigos.** Lista os resumos públicos dos amigos e as republicações feitas por eles, do mais recente para o mais antigo. A republicação conta pela data em que foi feita.
- **REQ-03 Item do feed.** Cada item mostra título, autor do livro, gênero, quem escreveu e os contadores de curtidas e comentários. Em uma republicação, mostra também quem republicou e o texto próprio dessa pessoa.
- **REQ-04 Ver mais.** O resumo aparece abreviado, com os primeiros 280 caracteres. O botão "ver mais" mostra o texto inteiro no próprio item, sem sair do feed.
- **REQ-05 Abrir o resumo.** Ao abrir um item, o usuário vê o resumo do livro, os capítulos que a pessoa resumiu, as curtidas, os comentários e o botão de republicar.
- **REQ-06 Paginação.** A rolagem carrega mais 20 itens por vez, até acabarem.
- **REQ-07 Filtro por gênero.** O usuário filtra os dois feeds por um gênero da lista.
- **REQ-08 Feed vazio.** Sem itens, o feed mostra uma mensagem. Para quem não tem amigos, ela convida a buscar amigos por username.
- **REQ-09 Atualização.** Quando algo deixa de ser público ou a amizade é desfeita, o item some do feed.

## 5. Regras de negócio

- O feed só mostra o que o usuário tem direito de ver. Privado nunca aparece.
- Só entra no feed o livro público que tem resumo do livro.
- A ordem não muda quando o resumo é editado nem quando recebe curtidas ou comentários.
- O próprio resumo e as próprias republicações do usuário não aparecem no feed dele.
- Os dois feeds usam a mesma tela de item e de resumo aberto, e diferem apenas na origem dos itens.
- O resumo abreviado é cortado no limite de 280 caracteres sem deixar uma palavra pela metade.

## 6. Dependências

- **Leitura e resumos:** resumos, visibilidade e data de publicação.
- **Amizades:** quem são os amigos.
- **Interações:** contadores, comentários e republicações.
- **Catálogo:** título, autor e gênero do livro.

## 7. Critérios de aceite

- Um resumo privado não aparece em nenhum feed, nem para amigos.
- O feed geral não mostra republicações. O feed de amigos mostra.
- Editar um resumo antigo não o leva ao topo do feed.
- Um resumo com mais de 280 caracteres aparece abreviado, e o "ver mais" mostra o texto inteiro sem sair do feed.
- Com 45 itens, a rolagem carrega 20, depois 20 e depois 5.
- Filtrar por um gênero mostra só os livros desse gênero.
- Um livro de gênero não identificado aparece no feed sem etiqueta de gênero.
- Um usuário sem amigos vê, no feed de amigos, a mensagem que o convida a buscar amigos.
- Depois de desfeita a amizade, os itens daquele amigo somem do feed de amigos.
- O usuário não vê os próprios resumos no feed.

## 8. Questões em aberto

- Se o mesmo resumo chegar duas vezes ao feed de amigos (publicado pelo autor e republicado por outro amigo, ou republicado por dois amigos), aparece uma vez ou várias. *(sugestão: uma única vez, na posição mais recente)*
- Se um resumo que volta de privado para público ganha uma nova data de publicação.
- Se "mostrar tudo" no item inclui os capítulos. Assumi que o "ver mais" expande só o resumo do livro e que os capítulos ficam na tela do resumo aberto.
- Se o filtro por gênero aceita mais de um gênero ao mesmo tempo.

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Mesmo resumo chegando duas vezes ao feed de amigos: aparece uma vez, na posição mais recente.
- Resumo que volta de privado para público: mantém a data de publicação original.
- "Mostrar tudo" no item: só o resumo do livro; os capítulos ficam ao abrir o resumo.
- Filtro por gênero: um por vez.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
