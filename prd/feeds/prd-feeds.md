---
tipo: prd
formato: 3
contexto: feeds
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [feed, feed geral, feed de amigos, resumo, republicação, filtro por gênero, ver mais, feed vazio]
atualizado: 2026-11-22
---
# PRD — feeds — Descobrir leituras pelo que outros resumiram

> Sexto contexto. Depende de **leitura-e-resumos**, **amizades**, **interacoes** e do gênero em **catalogo**.

## 1. Problema

O leitor quer descobrir o que ler pelo que outros leitores entenderam dos livros, e não só por notas ou capas. Sem um lugar que reúna esses resumos, a descoberta fica restrita ao que ele já conhece.

## 2. Objetivo e métricas de sucesso

Reunir os resumos públicos em dois feeds, um geral e um só de amigos, com as mesmas funções.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Usuários ativos que abrem o feed em 7 dias | n/a (produto novo) | 60% | Aberturas de feed |
| Itens do feed abertos para ler o resumo inteiro | n/a | 20% | Cliques em "ver mais" ou abertura |
| Livros adicionados à estante a partir do feed | n/a | 5% por usuário ativo | Origem do registro |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora curiosa | Explora o feed geral | Descobrir livros por resumos de desconhecidos | Não há onde ver resumos |
| Leitora com amigos | Acompanha o feed de amigos | Ver só o que quem ela conhece compartilha | Não há conteúdo social |

## 4. Escopo

### Dentro
- Feed geral: todos os resumos públicos, sem republicações.
- Feed de amigos: os resumos públicos dos amigos e as republicações deles; não existe visibilidade "só amigos".
- Ordem cronológica pela data em que o resumo ficou público; editar não reordena.
- Item mostra título, autor, gênero, nome e username de quem escreveu, resumo abreviado (280 caracteres) com "ver mais" no próprio item, e contadores.
- Ao abrir o resumo: resumo do livro, capítulos, curtidas, comentários e republicar, igual nos dois feeds.
- Rolagem infinita, 20 itens por vez, com filtro por gênero.
- Os próprios resumos não aparecem no feed; livros sem gênero identificado aparecem sem a etiqueta.
- Feed vazio convida a buscar amigos por username.

### Fora
- Ordenação por relevância e busca por livro ou autor: fora do primeiro corte.
- Notificações de novos itens: fora do primeiro corte.

## 5. Jornadas

### J1 — Descobrir no feed geral
1. A leitora abre o feed e percorre os resumos, lendo o abreviado.
2. Clica em "ver mais" ou abre o resumo para ver os capítulos.
3. Curte, comenta, republica ou adiciona o livro à estante.

### J2 — Acompanhar amigos
1. A leitora abre o feed de amigos.
2. Vê os resumos públicos e as republicações dos amigos.
3. Se não tem amigos, é convidada a buscá-los por username.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | Privado nunca aparece em nenhum feed | Garante a escolha do autor | Regra explícita |
| Restrição | Só entra no feed o livro público com resumo do livro | Livros sem resumo do livro não circulam | Regra de leitura-e-resumos |
| Premissa | O "ver mais" expande só o resumo do livro; os capítulos ficam no resumo aberto | Mantém o feed leve | Regra explícita |
| Risco | Feed cronológico dominado por quem publica muito | Conteúdo repetitivo | Evolução para relevância |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Mesmo resumo chegando duas vezes ao feed de amigos: aparece uma vez, na posição mais recente.
- Resumo que volta de privado para público: mantém a data de publicação original.
- "Mostrar tudo" no item: só o resumo do livro; os capítulos ficam ao abrir o resumo.
- Filtro por gênero: um por vez.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
