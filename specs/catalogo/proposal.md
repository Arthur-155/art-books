---
tipo: mudanca
formato: 3
contextos: [catalogo]
versao: "1.0"
status: proposto
depende_de: [contas]
superficies: [backend]
keywords: [catálogo, livro, título, autor, normalização, gênero, IA, pesquisa em segundo plano, lista de gêneros, livro não verificado, correção do gênero, administrador]
atualizado: 2026-11-22
---
# Mudança — catalogo

## Por quê

Sem um catálogo único, cada leitor informaria o livro de novo, o mesmo livro
apareceria duplicado e com gêneros diferentes, e não haveria como agrupar
leituras, filtrar por gênero nem recomendar (`prd/catalogo/prd-catalogo.md` §1).
Esta mudança cria o contexto: o livro entra uma vez, a IA define o gênero uma
vez, dentro de uma lista curada, e todos o reaproveitam.

## O que muda

- catalogo: REQ-01 novo (registro de livro novo)
- catalogo: REQ-02 novo (reconhecimento do livro existente por título e autor normalizados)
- catalogo: REQ-03 novo (livro e pesquisa únicos, mesmo em registro simultâneo)
- catalogo: REQ-04 novo (registro exige e-mail confirmado)
- catalogo: REQ-05 novo (estado visível do livro; livro não verificado com gênero mantém a marca)
- catalogo: REQ-06 novo (gênero escolhido pela IA só da lista ativa)
- catalogo: REQ-07 novo (3 novas tentativas, em 1, 5 e 15 min)
- catalogo: REQ-08 novo (livro não verificado)
- catalogo: REQ-09 novo (correção do gênero, valendo para todos; recusada enquanto a pesquisa roda)
- catalogo: REQ-10 novo (correção só do registrante e do administrador; sem registrante, só o administrador)
- catalogo: REQ-11 novo (lista inicial de gêneros inativa até a aprovação)
- catalogo: REQ-12 novo (inclusão, renomeação e desativação de gêneros pelo administrador)
- catalogo: REQ-13 novo (desativar gênero preserva os livros)
- catalogo: REQ-14 novo (lista de gêneros ativos)

## Impacto

O código ainda não tem nenhum contexto: só `backend/art-books/main.go`. A
mudança cria `backend/domain/catalogo/`, as tabelas `genre` e `book` e um
worker de pesquisa no binário `art-books`, e passa a chamar a API da Anthropic.
Depende da gênese de `contas`, entregue antes desta: a tabela `user` (FK do
registrante) e o token com o papel e o indicador de e-mail confirmado.
`leitura-e-resumos`, `recomendacoes` e
`feeds` vão ler livro e gênero daqui, sem mudar o catálogo.

## Fora desta mudança

- Telas web: o aviso de digitar título e autor exatamente (PRD REQ-01) e a
  exibição do passo da pesquisa ficam para a mudança de interface do catálogo.
- Capa, sinopse, ISBN, edição e mais de um gênero por livro (PRD §4 Fora).
- Busca por similaridade, listagem de livros e fusão de duplicados por erro de
  digitação (PRD §4 Fora).
- Tela de revisão de livros não verificados: o administrador age direto na base.
- Troca do provedor de IA por um roteador de LLMs.
