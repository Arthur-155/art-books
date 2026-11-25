---
formato: memoria
contexto: recomendacoes
versao: "1.0"
status: spec
atualizado: 2026-11-22
keywords: [recomendação, IA, resumo breve, motivo, dispensar, quero ler, pedir outras, limite diário, livro confirmado, geração em segundo plano]
---
# recomendacoes

## Estado atual

Gênese escrita em `specs/recomendacoes/` (mudança `recomendacoes`): proposta, `deltas/recomendacoes.md` com REQ-01 a REQ-12 e `design.md` (Go, serviço `art-books`, PostgreSQL, worker de geração no binário, API da Anthropic atrás de `RecommendationGenerator`, estante e catálogo atrás de `ShelfGateway` e `CatalogGateway`). Status `proposto`, sem aprovação, sem marcadores em aberto. `gofi spec check recomendacoes`: 0 erros.

Decisões aplicadas: a escolha da recomendação do perfil fica só no Perfil; nova lista quando um "lido" é posterior à lista, gerada ao abrir a tela; limite de 5 pedidos de outras sugestões por dia em America/Sao_Paulo, contando só o pedido que gerou lista, uma vez; até 2 rodadas extras para completar 10, mostrando o que houver; pedir outras substitui a lista e as não tratadas seguem elegíveis; resumo breve em tabela própria `book_brief`, por título e autor normalizados; falha da IA vai a `failed`, mantém a lista anterior e só o usuário tenta de novo.

A proposta declara `depende_de: [catalogo]`; `contas` e `leitura-e-resumos` ainda não têm mudança e entram no `depende_de` da proposta quando existirem. O delta declara `depende_de: [contas, catalogo, leitura-e-resumos]`. A gênese de `leitura-e-resumos` precisa expor a leitura da estante (título, autor, gênero, estado, data em que foi marcado como "lido") e a operação de colocar como "quero ler" sem duplicar.

Próximo passo: `gofi review recomendacoes` e `gofi change approve recomendacoes` pela pessoa; a entrega espera as gêneses de `contas`, `catalogo` e `leitura-e-resumos`.

## Histórico de versões
