---
tipo: mudanca
formato: 3
contextos: [recomendacoes]
versao: "1.0"
status: aprovado
depende_de: [catalogo]
superficies: [backend]
keywords: [recomendação, IA, gênero, autor, motivo, resumo breve, quero ler, dispensar, pedir outras, limite diário, livro confirmado, segundo plano]
atualizado: 2026-11-22
aprovacao_versao: "1.0"
aprovacao_sha: "sha256:3dfb898ccb00b4eaed4a95a486d5d9103b5c80cee2b827ece4004e5dffb33ad2"
---
# Mudança — recomendacoes

## Por quê

Depois de terminar um livro, o leitor quer saber o que ler a seguir e hoje não
tem onde descobrir (`prd/recomendacoes/prd-recomendacoes.md` §1). Esta mudança
cria o contexto: a partir do gênero e do autor dos livros "lidos", uma IA
sugere até 10 livros confirmados, cada um com motivo e resumo breve, que o
leitor adiciona à estante, dispensa ou troca, dentro de um limite diário que
controla o custo.

## O que muda

- recomendacoes: REQ-01 novo (geração da lista a partir dos livros "lidos", sem livro da estante nem dispensado)
- recomendacoes: REQ-02 novo (motivo de uma frase)
- recomendacoes: REQ-03 novo (resumo breve de até 1.000 caracteres, reaproveitado entre usuários)
- recomendacoes: REQ-04 novo (nova lista depois de um livro marcado como "lido")
- recomendacoes: REQ-05 novo (adicionar à estante como "quero ler", sem duplicar no catálogo nem na estante)
- recomendacoes: REQ-06 novo (dispensar, sem voltar)
- recomendacoes: REQ-07 novo (usuário sem livros "lidos")
- recomendacoes: REQ-08 novo (falha da IA e nova tentativa)
- recomendacoes: REQ-09 novo (pedir outras sugestões substitui a lista)
- recomendacoes: REQ-10 novo (limite de 5 pedidos de outras sugestões por dia)
- recomendacoes: REQ-11 novo (só livro confirmado; até 2 rodadas para completar a lista)
- recomendacoes: REQ-12 novo (recomendações pessoais)

## Impacto

O código ainda não tem o contexto. A mudança cria `backend/domain/recomendacoes/`,
cinco tabelas e um worker de geração no binário `art-books`, e chama a API da
Anthropic com as mesmas variáveis do catálogo. Depende de três gêneses,
entregues antes desta: `contas` (tabela `user` e token com o indicador de
e-mail confirmado), `catalogo` (normalização, busca do livro por título e autor
normalizados e registro de livro novo) e `leitura-e-resumos` (estante). A
gênese de `leitura-e-resumos` precisa expor à recomendacoes a leitura da
estante de um usuário (cada livro com título, autor, gênero, estado e a data
em que foi marcado como "lido") e a operação de colocar um livro como "quero
ler" sem duplicar (D8 de Leitura); o design diz a interface. Nenhum outro
contexto muda de comportamento.

## Fora desta mudança

- A recomendação escolhida para o perfil: é do Perfil, que a escolhe livre no
  catálogo, sem depender desta lista.
- Telas web: os textos da mensagem sem livros "lidos", do passo da geração, do
  erro e da explicação do limite ficam para a mudança de interface.
- Recomendação pelo que os amigos leram, por humor ou tema livre, e
  notificação de novas sugestões (PRD §Fora).
