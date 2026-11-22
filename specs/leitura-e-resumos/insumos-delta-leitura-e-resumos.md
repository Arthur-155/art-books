---
tipo: insumo-delta
contexto: leitura-e-resumos
versao: "1.0"
status: discovery
prd: prd/leitura-e-resumos/prd-leitura-e-resumos.md
keywords: [leitura, estante, quero ler, lendo, lido, capítulo, resumo de capítulo, resumo do livro, visibilidade, público, privado]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `leitura-e-resumos`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/leitura-e-resumos/prd-leitura-e-resumos.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Leitura e resumos

Terceiro contexto do discovery. Depende de **Contas** (usuário e e-mail confirmado) e de **Catálogo** (livro e gênero). É a base de **Interações**, **Feeds**, **Perfil** e **Recomendações**.

## 1. Problema e objetivo

Quem lê não tem onde registrar o que entendeu de cada capítulo. Este contexto permite acompanhar a leitura de um livro, escrever com as próprias palavras o que cada capítulo diz e, ao terminar, escrever um resumo do livro inteiro.

## 2. Escopo

**Dentro:** estado da leitura de cada livro, capítulos, resumo de capítulo, resumo do livro, visibilidade pública ou privada, edição.

**Fora (por enquanto):** meta de leitura, progresso por página, notas privadas fora do resumo, importação de capítulos.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Estados da leitura | Três: "quero ler", "lendo" e "lido" |
| D2 | Leituras simultâneas | O usuário pode ter vários livros "lendo". O perfil mostra o mais recente |
| D3 | Capítulos | O usuário adiciona capítulos livremente, com número e título opcional |
| D4 | Resumo de capítulo | Máximo de 4.000 caracteres. Editável a qualquer momento |
| D5 | Resumo do livro | Opcional. É ele que aparece no feed |
| D6 | Visibilidade | Uma única, pública ou privada, por livro. Vale para o resumo do livro e para os capítulos |
| D7 | Mudança de visibilidade | Ao virar privado, o resumo some do feed. O destino de curtidas, comentários e republicações é decidido em Interações |
| D8 | Mesmo livro de novo | Um registro por usuário e por livro, editável. Não existe segunda leitura |

## 4. Requisitos funcionais

- **REQ-01 Adicionar livro.** O usuário escolhe ou registra um livro pelo Catálogo e o adiciona à sua estante, em um dos estados da D1.
- **REQ-02 Mudar estado.** O usuário muda o estado de "quero ler" para "lendo" e depois "lido". Ele pode voltar atrás.
- **REQ-03 Capítulos.** O usuário adiciona, renomeia, reordena e remove capítulos do livro dele.
- **REQ-04 Resumo de capítulo.** Para cada capítulo, o usuário escreve com as próprias palavras o que ele diz, até 4.000 caracteres, e pode editar depois.
- **REQ-05 Resumo do livro.** O usuário pode escrever o resumo do livro inteiro, sem limite de estado. Livros "lidos" sem resumo seguem válidos.
- **REQ-06 Visibilidade.** O usuário escolhe se o livro dele é público ou privado e pode mudar depois (D6).
- **REQ-07 Publicação.** Só aparece para outros usuários o que é público. Publicar exige e-mail confirmado (D3 de Contas).
- **REQ-08 Estante.** O usuário vê seus livros por estado, com o resumo e os capítulos de cada um.
- **REQ-09 Remover livro.** O usuário remove o livro da estante, apagando os capítulos e os resumos dele. O livro continua no catálogo.

## 5. Regras de negócio

- Existe um único registro por usuário e livro (D8).
- Só entra no feed o livro público que tem resumo do livro (D5). Capítulos de livro público ficam visíveis ao abrir o resumo.
- O resumo é do usuário. Ninguém além dele edita o texto, nem o administrador.
- Dois capítulos do mesmo livro não repetem o número.
- O livro "mais recente" do perfil (D2) é o que teve a última atividade: mudança para "lendo" ou edição de um capítulo.
- A mudança de visibilidade é imediata e vale para tudo do livro, inclusive os capítulos.

## 6. Dependências

- **Contas:** usuário e estado "e-mail confirmado".
- **Catálogo:** livro e gênero.
- **Interações**, **Feeds** e **Perfil** leem os resumos e a visibilidade definidos aqui.
- **Recomendações** usa os livros "lidos".

## 7. Critérios de aceite

- O usuário tem dois livros "lendo" ao mesmo tempo e o perfil mostra o de atividade mais recente.
- Um resumo de capítulo com 4.001 caracteres é recusado, com mensagem clara. Com 4.000 é aceito.
- O usuário não consegue registrar o mesmo livro duas vezes na estante.
- Um livro privado nunca aparece para outro usuário, nem por link direto.
- Um livro público sem resumo do livro não aparece no feed.
- Ao mudar um livro de público para privado, ele some do feed na hora.
- O usuário com e-mail não confirmado não consegue tornar um livro público.
- Remover o livro da estante não remove o livro do catálogo.

## 8. Questões em aberto

- Limite de caracteres do resumo do livro. Só o do capítulo foi definido (4.000).
- Visibilidade inicial de um livro novo: privado ou público. *(sugestão: privado, para ninguém publicar sem querer)*
- Se existe um limite de capítulos por livro.
- O que acontece com os livros públicos do usuário se o e-mail dele deixa de estar confirmado, como numa troca de e-mail em andamento.

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Resumo do livro: até 4.000 caracteres, igual ao do capítulo.
- Visibilidade inicial de um livro novo: privado.
- Limite de 200 capítulos por livro.
- Troca de e-mail em andamento: os livros públicos não mudam, porque o e-mail antigo segue confirmado até a troca ser concluída.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
