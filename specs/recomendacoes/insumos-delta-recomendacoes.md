---
tipo: insumo-delta
contexto: recomendacoes
versao: "1.0"
status: discovery
prd: prd/recomendacoes/prd-recomendacoes.md
keywords: [recomendação, IA, gênero, autor, resumo breve, motivo, dispensar, quero ler, limite diário]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `recomendacoes`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/recomendacoes/prd-recomendacoes.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Recomendações

Sétimo contexto do discovery. Depende de **Leitura e resumos** (livros lidos e estantes) e de **Catálogo** (gênero e autor). Alimenta a tela de recomendações e a escolha da recomendação do **Perfil**.

## 1. Problema e objetivo

Depois de terminar um livro, o leitor quer saber o que ler a seguir. O sistema usa os livros que o usuário leu para sugerir novos livros, cada um com um motivo e um resumo breve para ajudar na escolha.

## 2. Escopo

**Dentro:** lista de recomendações por usuário, motivo, resumo breve de cada livro, adicionar à estante, dispensar, pedir outras, limite de uso.

**Fora (por enquanto):** recomendação baseada no que os amigos leram, recomendação por humor ou tema livre, notificação de novas recomendações.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Fonte | A IA pode sugerir livros de fora do catálogo. Eles só entram no catálogo se o usuário os registrar |
| D2 | Base | Só os livros marcados como "lido" |
| D3 | Sinais | Gênero e autor dos livros lidos. Não usa o que os amigos leram |
| D4 | Quantidade | 10 por vez |
| D5 | Atualização | Quando o usuário marca um novo livro como "lido". Há um botão para pedir outras |
| D6 | Repetição | Nunca recomenda um livro que o usuário já tenha em qualquer estado, inclusive "quero ler" |
| D7 | Sem livros lidos | Mensagem convidando a registrar um livro lido, sem recomendação genérica |
| D8 | Motivo | Cada recomendação traz uma frase curta, por exemplo "porque você leu um livro de ficção científica" |
| D9 | Resumo breve | Cada recomendação traz um resumo breve do livro, de no máximo 1.000 caracteres |
| D10 | Ações | Adicionar à estante como "quero ler" e dispensar. A dispensada não volta |
| D11 | Recomendação do perfil | O usuário escolhe qualquer livro do catálogo, mesmo que não esteja na lista de recomendações |
| D12 | Falha da IA | Mensagem de erro com nova tentativa. Há um limite diário de pedidos por usuário, a definir |

## 4. Requisitos funcionais

- **REQ-01 Gerar recomendações.** Com base no gênero e no autor dos livros "lidos", a IA sugere 10 livros, sem repetir nenhum que o usuário já tenha ou tenha dispensado.
- **REQ-02 Motivo.** Cada livro recomendado traz uma frase curta explicando por que foi sugerido.
- **REQ-03 Resumo breve.** Cada livro recomendado traz um resumo breve gerado pela IA, com no máximo 1.000 caracteres, sobre o livro em si. Não é o resumo de nenhum usuário.
- **REQ-04 Atualizar.** A lista é refeita quando o usuário marca um livro como "lido". O usuário também pede outras sugestões por um botão.
- **REQ-05 Adicionar à estante.** O usuário adiciona a recomendação à estante como "quero ler". Se o livro ainda não existe no catálogo, ele é registrado com o título e o autor da sugestão, e a pesquisa de gênero segue o fluxo do Catálogo.
- **REQ-06 Dispensar.** O usuário dispensa uma recomendação, que não volta a aparecer para ele.
- **REQ-07 Sem livros lidos.** Sem nenhum livro "lido", a tela mostra a mensagem do D7.
- **REQ-08 Falha e limite.** Se a IA falhar, o usuário vê um erro e pode tentar de novo. Ao passar do limite diário, o botão de pedir outras fica indisponível, com a explicação.
- **REQ-09 Escolha do perfil.** O usuário escolhe uma única recomendação para o perfil, entre os livros do catálogo (detalhada em Perfil).

## 5. Regras de negócio

- As recomendações são pessoais: cada usuário tem a sua lista e ninguém vê a de outro.
- A lista nunca inclui livro que o usuário já tenha na estante, em qualquer estado.
- Um livro dispensado não é sugerido de novo para o mesmo usuário.
- O resumo breve descreve o livro e não substitui nem copia o resumo de nenhum usuário.
- A recomendação de um livro fora do catálogo não cria o livro no catálogo até o usuário registrá-lo.
- O livro recomendado que já existe no catálogo é reconhecido pelo título e autor normalizados (D1 e D2 de Catálogo), para não gerar duplicado.
- A falha da IA nunca bloqueia o resto do produto.

## 6. Dependências

- **Leitura e resumos:** livros "lidos" e estado de cada livro do usuário.
- **Catálogo:** gênero, autor e reconhecimento do mesmo livro.
- **IA:** geração da lista, dos motivos e dos resumos breves, em segundo plano e com o passo visível ao usuário, como no Catálogo.
- **Perfil** usa a escolha do REQ-09.

## 7. Critérios de aceite

- Um usuário sem livros "lidos" vê a mensagem convidando a registrar um livro, sem lista.
- A lista tem 10 livros, e nenhum deles está na estante do usuário em qualquer estado.
- Cada livro mostra uma frase de motivo e um resumo breve de no máximo 1.000 caracteres.
- Marcar um novo livro como "lido" atualiza as recomendações.
- Um livro dispensado não reaparece nas listas seguintes.
- Adicionar uma recomendação que já está no catálogo não cria um segundo livro.
- Adicionar uma recomendação de fora do catálogo a registra e inicia a pesquisa de gênero.
- Se a IA falhar, o usuário vê o erro e a opção de tentar de novo.
- Depois do limite diário, o botão de pedir outras fica indisponível.

## 8. Questões em aberto

- O valor do limite diário de pedidos por usuário.
- Como evitar que a IA sugira um livro que não existe. *(sugestão: confirmar título e autor antes de mostrar, e descartar o que não for confirmado)*
- Se o resumo breve de um livro do catálogo é guardado no livro e reaproveitado entre usuários, para reduzir custo. *(sugestão: sim)*
- Se o sistema pode sugerir livros de um autor que o usuário já leu, ou deve variar. *(sugestão: pode, e um dos sinais é o autor)*
- O que acontece com as recomendações quando o usuário remove da estante um livro "lido".

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Limite de 5 pedidos de "outras sugestões" por dia, por usuário.
- Livro inexistente: o título e o autor são confirmados antes de mostrar, e o que não for confirmado é descartado.
- Resumo breve: guardado no livro do catálogo e reaproveitado entre usuários.
- Remover um livro "lido" da estante: as recomendações só mudam na próxima atualização.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
