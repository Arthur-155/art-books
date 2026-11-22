---
tipo: prd
formato: 3
contexto: perfil
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [perfil, lendo agora, últimos livros lidos, livros lidos, recomendação escolhida, perfil de amigo, visibilidade]
atualizado: 2026-11-22
---
# PRD — perfil — A vitrine do leitor

> Oitavo e último contexto. Depende de **contas**, **leitura-e-resumos**, **amizades** e **recomendacoes**. Só lê o que os outros contextos produzem.

## 1. Problema

Um leitor não tem como mostrar a outras pessoas o que está lendo, o que já leu e o que indica. Sem essa vitrine, conhecer o gosto de um amigo exige perguntar a ele.

## 2. Objetivo e métricas de sucesso

Dar a cada leitor uma página pública com a sua identidade de leitor e permitir que outras pessoas a visitem.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Usuários que escolhem uma recomendação no perfil | n/a (produto novo) | 30% | Perfis com recomendação |
| Visitas a perfis de outros usuários por usuário ativo | n/a | 2 por semana | Aberturas de perfil |
| Pedidos de amizade enviados a partir do perfil | n/a | 20% dos pedidos | Origem do pedido |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora dona do perfil | Mostra suas leituras | Uma página com o que leu e indica | Não há perfil |
| Visitante | Quer conhecer o gosto de alguém | Ver o que é público | Não há como ver |

## 4. Escopo

### Dentro
- Cabeçalho com nome de exibição, username e total de livros lidos.
- Livro "lendo agora": o de atividade mais recente; sem nenhum, o bloco some.
- Os últimos 5 livros lidos, pela data de conclusão, do mais recente para o mais antigo.
- Tela com todos os livros lidos, com rolagem infinita e filtro por gênero.
- Uma única recomendação escolhida pelo dono: um livro do catálogo, com frase opcional de até 280 caracteres, editável.
- Perfil aberto por qualquer usuário, mostrando só livros públicos; amigo e não amigo veem o mesmo.
- Botões de amizade no perfil de outra pessoa; o dono vê o próprio perfil como os outros, com aviso do que é privado.

### Fora
- Foto, biografia, data de entrada e estatísticas de leitura: fora do primeiro corte.
- Perfil visível só para amigos: não existe por enquanto.

## 5. Jornadas

### J1 — Montar o próprio perfil
1. A leitora marca livros como lendo e lidos.
2. Escolhe uma recomendação e escreve uma frase.
3. Vê o perfil como os outros o veem.

### J2 — Visitar o perfil de alguém
1. A leitora abre o perfil de um username.
2. Vê o livro que a pessoa lê agora, os últimos 5 lidos e a recomendação.
3. Pode abrir todos os lidos e enviar um pedido de amizade.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | Livro privado só aparece para o dono | A amizade nunca dá acesso a privados | Regra de leitura-e-resumos |
| Restrição | O e-mail nunca aparece no perfil | Privacidade | Regra de contas |
| Premissa | O total de livros lidos visto por outros conta só os públicos | Total menor que o do dono | Regra explícita |
| Premissa | A recomendação escolhida aparece sempre | É uma indicação deliberada | Regra explícita |
| Risco | Perfil de conta apagada | Link quebrado | Perfil deixa de existir com a conta |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Perfil de quem não é amigo: abre para todos, sem opção de esconder no primeiro corte.
- Total de livros lidos visto por outros: conta só os públicos.
- "Lendo agora": sem progresso.
- Livro privado: privado por inteiro, sem aparecer só com título e autor.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
