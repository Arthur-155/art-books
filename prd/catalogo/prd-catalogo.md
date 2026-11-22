---
tipo: prd
formato: 3
contexto: catalogo
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [catálogo, livro, título, autor, normalização, gênero, IA, pesquisa em segundo plano, lista de gêneros, livro não verificado]
atualizado: 2026-11-22
---
# PRD — catalogo — Catálogo compartilhado de livros

> Segundo contexto. Depende de **contas** (usuário e administrador) e é a base de **leitura-e-resumos** e **recomendacoes**.

## 1. Problema

Cada leitor que registra um livro teria de informar tudo de novo, e o mesmo livro apareceria duplicado, com gêneros diferentes. Sem um catálogo único, não há como agrupar leituras, filtrar por gênero nem recomendar.

## 2. Objetivo e métricas de sucesso

Ter um catálogo único em que cada livro entra uma vez, ganha um gênero definido pela IA e é reaproveitado por todos os usuários.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Registros que reaproveitam um livro já existente | n/a (produto novo) | 60% após 3 meses | Registros com livro existente sobre o total |
| Livros com gênero definido | n/a | 90% | Livros com gênero sobre o total |
| Livros com gênero corrigido depois da IA | n/a | menos de 10% | Correções sobre livros com gênero |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora | Registra livros | Achar o livro já pronto, com gênero | Não há catálogo |
| Administrador | Cuida da lista de gêneros | Revisar gêneros e corrigir livros | Não há lista nem ferramenta |

## 4. Escopo

### Dentro
- Registro de livro por título e autor; o mesmo livro é reconhecido por título e autor normalizados (ignora maiúsculas, acentos, espaços duplicados e pontuação nas pontas), e a tela avisa que o nome deve ser digitado exatamente.
- Pesquisa do gênero pela IA em segundo plano, com o passo visível ("Pesquisando informações sobre o livro").
- Lista fechada de gêneros, de um único gênero por livro; a IA só escolhe dela. A lista inicial vem de um agente que lê uma fonte aberta (Wikipédia) uma única vez e é revisada pelo administrador.
- Correção do gênero por quem registrou o livro e pelo administrador, valendo para todos.
- Livro com "gênero não identificado" quando a IA falha, com nova tentativa e escolha manual; livro "não verificado" quando a IA conclui que ele não existe.

### Fora
- Capa, sinopse, ISBN e edição: fora do primeiro corte.
- Busca por similaridade e fusão de duplicados por erro de digitação: ficam para depois.
- Mais de um gênero por livro: um só simplifica as recomendações.

## 5. Jornadas

### J1 — Registrar um livro novo
1. A leitora informa título e autor, digitados exatamente.
2. O livro entra no catálogo e ela vê "Pesquisando informações sobre o livro".
3. O gênero aparece quando a pesquisa termina; ela segue usando o produto enquanto isso.

### J2 — Registrar um livro que já existe
1. A leitora informa título e autor.
2. O livro é reconhecido e já traz o gênero, sem nova pesquisa.

### J3 — Corrigir o gênero
1. Quem registrou o livro, ou o administrador, escolhe outro gênero da lista.
2. A mudança vale para todos os usuários.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | A IA só escolhe gêneros da lista ativa | Evita variações como "Sci-fi" e "Ficção científica" | Lista revisada pelo administrador |
| Restrição | Se quem registrou apaga a conta, o livro fica sem dono | Só o administrador corrige o gênero | Regra explícita |
| Premissa | Dois usuários registrando o mesmo livro ao mesmo tempo geram uma única pesquisa | Evita custo duplicado | A pesquisa é única por livro |
| Risco | Erro de digitação cria um livro duplicado | Catálogo com repetições | Aviso na tela; fusão fica para depois |
| Premissa | A fonte dos gêneros é aberta (Wikipédia, CC BY-SA) e curada pelo administrador | Evita risco de licença | Curadoria na primeira revisão |
| Risco | A IA falha ou inventa | Livro sem gênero ou com gênero errado | Estados "não identificado" e "não verificado" e correção |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Fonte dos gêneros: Wikipédia (CC BY-SA), curada pelo administrador na primeira revisão; a BISAC fica descartada por não ter os termos validados.
- Pesquisa de gênero: 3 tentativas, com intervalos de 1, 5 e 15 min.
- Lista de gêneros: de 20 a 25, revisados pelo administrador.
- Livros "não verificados": sem tela de revisão no primeiro corte; o administrador age direto na base.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
