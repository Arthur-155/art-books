---
tipo: insumo-delta
contexto: perfil
versao: "1.0"
status: discovery
prd: prd/perfil/prd-perfil.md
keywords: [perfil, lendo agora, últimos livros lidos, livros lidos, recomendação escolhida, perfil de amigo, visibilidade]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `perfil`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/perfil/prd-perfil.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Perfil

Oitavo e último contexto do discovery. Depende de **Contas**, **Leitura e resumos**, **Amizades** e **Recomendações**. Só lê o que os outros contextos produzem.

## 1. Problema e objetivo

O perfil é a vitrine do leitor: o que ele está lendo, o que leu e o que indica. Qualquer pessoa pode abrir o perfil de outra e ver só o que é público.

## 2. Escopo

**Dentro:** cabeçalho, livro "lendo agora", últimos 5 livros lidos, todos os livros lidos, recomendação escolhida, ações de amizade no perfil, visão do próprio perfil.

**Fora (por enquanto):** foto, biografia, data de entrada, seguidores, estatísticas de leitura.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Lendo agora | O livro "lendo" de atividade mais recente. Sem nenhum, o bloco some |
| D2 | Últimos 5 lidos | Ordenados pela data em que foram marcados como "lido", do mais recente para o mais antigo |
| D3 | Todos os lidos | Tela à parte, com rolagem infinita e filtro por gênero |
| D4 | Visibilidade dos livros | Segue a visibilidade do livro, definida em Leitura. O livro público aparece com título, autor e resumo. O livro privado só aparece para o dono |
| D5 | Amigo e não amigo | Veem o mesmo. A amizade muda o feed, não o perfil |
| D6 | Recomendação escolhida | Um livro do catálogo, com frase opcional de até 280 caracteres. Pode ficar vazia e é editável |
| D7 | Cabeçalho | Nome de exibição, username e total de livros lidos |
| D8 | Ações em perfil alheio | Enviar pedido de amizade e desfazer a amizade |
| D9 | Perfil próprio | O usuário o vê como os outros o veem, com um aviso do que é privado e não aparece para eles |

## 4. Requisitos funcionais

- **REQ-01 Abrir perfil.** Qualquer usuário abre o perfil de outro pelo username. O e-mail nunca aparece.
- **REQ-02 Cabeçalho.** Mostra o nome de exibição, o username e o total de livros lidos.
- **REQ-03 Lendo agora.** Mostra o livro "lendo" com a atividade mais recente. Sem livro "lendo", o bloco não aparece.
- **REQ-04 Últimos 5 lidos.** Mostra os 5 livros marcados como "lido" mais recentemente.
- **REQ-05 Todos os lidos.** Leva a uma tela com todos os livros lidos, com rolagem infinita e filtro por gênero.
- **REQ-06 Recomendação escolhida.** Mostra um único livro que o dono indica, com a frase dele, se houver. O dono escolhe, edita e remove.
- **REQ-07 Ações de amizade.** No perfil de outro usuário, o botão muda conforme a relação: enviar pedido, cancelar pedido enviado, responder pedido recebido ou desfazer amizade.
- **REQ-08 Perfil próprio.** O dono vê o perfil como os outros veem e um aviso dos livros privados que não aparecem para eles.

## 5. Regras de negócio

- Para quem não é o dono, só aparecem livros públicos. Os blocos "lendo agora", "últimos 5" e "todos os lidos" e o total de livros lidos contam apenas livros públicos.
- Para o dono, o total e as listas incluem também os livros privados.
- A recomendação escolhida aparece sempre, porque é uma indicação deliberada do usuário e não um livro da estante dele.
- Se o livro da recomendação escolhida sair do catálogo, o espaço fica vazio.
- Quando a conta é apagada, o perfil deixa de existir (D7 de Contas).
- A amizade nunca dá acesso a livros privados (D5).

## 6. Dependências

- **Contas:** nome de exibição e username.
- **Leitura e resumos:** estados, datas e visibilidade dos livros.
- **Amizades:** estado da relação entre dois usuários.
- **Recomendações:** escolha do livro indicado (REQ-09 de Recomendações).
- **Catálogo:** título, autor e gênero.

## 7. Critérios de aceite

- O perfil de outra pessoa nunca mostra e-mail nem livros privados.
- O total de livros lidos que outra pessoa vê é menor que o do dono, se o dono tem livros lidos privados.
- Sem livro "lendo", o bloco "lendo agora" não aparece.
- Os últimos 5 lidos vêm ordenados pela data de conclusão, do mais recente para o mais antigo.
- A tela de todos os lidos pagina por rolagem infinita e filtra por gênero.
- Uma frase de recomendação com 281 caracteres é recusada.
- Amigo e não amigo veem exatamente o mesmo perfil.
- No perfil de um usuário com pedido pendente enviado por mim, o botão é "cancelar pedido".
- O dono vê o aviso dos livros privados que não aparecem para os outros.

## 8. Questões em aberto

- Se o usuário pode esconder o perfil inteiro de quem não é amigo.
- Se o total de livros lidos aparece para os outros como o número real ou só dos livros públicos (assumi só os públicos).
- Se o livro de "lendo agora" mostra também o progresso (capítulos já resumidos).

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Perfil de quem não é amigo: abre para todos, sem opção de esconder no primeiro corte.
- Total de livros lidos visto por outros: conta só os públicos.
- "Lendo agora": sem progresso.
- Livro privado: privado por inteiro, sem aparecer só com título e autor.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
