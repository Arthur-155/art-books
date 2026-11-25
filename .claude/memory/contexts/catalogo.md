---
formato: memoria
contexto: catalogo
versao: "1.0"
status: spec
atualizado: 2026-11-22
keywords: [catálogo, livro, gênero, normalização, pesquisa em segundo plano, lista de gêneros, livro não verificado, correção do gênero, administrador]
---
# catalogo

## Estado atual

Gênese escrita em `specs/catalogo/` (mudança `catalogo`): proposta, `deltas/catalogo.md` com REQ-01 a REQ-14 e `design.md` (Go, serviço `art-books`, PostgreSQL, worker de pesquisa no binário, API da Anthropic atrás de `BookResearcher`). Status `proposto`, sem aprovação, sem marcadores em aberto.

Decisões aplicadas: troca de gênero recusada enquanto a pesquisa roda (REQ-09.5, `ErrCatalogResearchInProgress` 409); livro não verificado que recebe gênero mantém a marca e mostra o gênero (REQ-05.5, REQ-08, REQ-09.6); a proposta declara `depende_de: [contas]`.

`gofi spec check catalogo` tem um erro: a mudança `contas` ainda não existe. Some quando a gênese de `contas` for criada.

Próximo passo: escrever a gênese de `contas` (`/gofi-spec`, sessão nova), que é entregue antes; depois `gofi review catalogo` e `gofi change approve catalogo`. Em seguida o loop implementa (gofi-eng) e audita (gofi-qa).

## Histórico de versões
