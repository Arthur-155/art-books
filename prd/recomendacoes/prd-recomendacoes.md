---
tipo: prd
formato: 3
contexto: recomendacoes
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [recomendação, IA, gênero, autor, resumo breve, motivo, dispensar, quero ler, limite diário]
atualizado: 2026-11-22
---
# PRD — recomendacoes — Sugerir o próximo livro a ler

> Sétimo contexto. Depende de **leitura-e-resumos** (livros lidos) e de **catalogo** (gênero e autor). Alimenta a recomendação escolhida em **perfil**.

## 1. Problema

Depois de terminar um livro, o leitor não sabe o que ler a seguir e acaba repetindo o que já conhece ou desistindo de ler. Falta uma sugestão que parta do que ele já leu e explique por que vale a pena.

## 2. Objetivo e métricas de sucesso

Sugerir livros a partir do que o usuário leu, cada um com motivo e resumo breve para ajudar na escolha.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Usuários com livros lidos que veem as recomendações | n/a (produto novo) | 50% | Aberturas da tela |
| Recomendações adicionadas à estante | n/a | 10% das mostradas | Adições sobre exibições |
| Recomendações dispensadas | n/a | menos de 40% | Dispensas sobre exibições |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora que terminou um livro | Busca o próximo | Sugestões com motivo e resumo | Busca por conta própria, sem critério |
| Leitora que está no começo | Tem poucos livros lidos | Entender o que precisa registrar para receber sugestões | Sem lista nem orientação |

## 4. Escopo

### Dentro
- Recomendações pessoais, baseadas só nos livros "lidos", usando gênero e autor.
- A IA pode sugerir livros de fora do catálogo; eles só entram nele se o usuário os registrar.
- 10 recomendações por vez, atualizadas quando um novo livro é marcado como "lido", com botão para pedir outras.
- Cada recomendação traz uma frase de motivo e um resumo breve do livro, de até 1.000 caracteres.
- Adicionar à estante como "quero ler" e dispensar; a dispensada não volta.
- Nunca recomenda um livro que o usuário já tenha em qualquer estado.
- Sem livros lidos, a tela convida a registrar um livro lido, sem lista genérica.
- Limite diário de pedidos por usuário, para controlar o custo, e mensagem de erro com nova tentativa se a IA falhar.

### Fora
- Recomendação baseada no que os amigos leram: fora do primeiro corte.
- Recomendação por humor ou tema livre e notificação de novas sugestões: ficam para depois.

## 5. Jornadas

### J1 — Receber sugestões
1. A leitora marca um livro como "lido".
2. As recomendações são atualizadas, cada uma com motivo e resumo breve.
3. Ela adiciona uma à estante como "quero ler" ou dispensa.

### J2 — Pedir outras sugestões
1. A leitora aciona o botão de pedir outras.
2. Recebe uma nova lista de 10, sem repetir o que já tem ou dispensou.
3. Se passou do limite diário, vê a explicação.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | O resumo breve descreve o livro e não copia o resumo de nenhum usuário | Evita atribuição indevida | Regra explícita |
| Premissa | O livro sugerido que já está no catálogo é reconhecido por título e autor normalizados | Evita duplicados | Regra de catalogo |
| Risco | A IA sugere um livro que não existe | Perda de confiança | Confirmar título e autor antes de mostrar |
| Risco | Custo da IA | Gasto alto por usuário | Limite diário e reaproveitamento do resumo breve |
| Risco | A IA falha | Tela sem sugestões | Mensagem de erro e nova tentativa |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Limite de 5 pedidos de "outras sugestões" por dia, por usuário.
- Livro inexistente: o título e o autor são confirmados antes de mostrar, e o que não for confirmado é descartado.
- Resumo breve: guardado no livro do catálogo e reaproveitado entre usuários.
- Remover um livro "lido" da estante: as recomendações só mudam na próxima atualização.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
