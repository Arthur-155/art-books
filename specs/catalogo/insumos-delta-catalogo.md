---
tipo: insumo-delta
contexto: catalogo
versao: "1.0"
status: discovery
prd: prd/catalogo/prd-catalogo.md
keywords: [catálogo, livro, título, autor, normalização, gênero, IA, pesquisa em segundo plano, lista de gêneros, livro não verificado]
atualizado: 2026-11-22
---
> **Insumo do delta da primeira mudança de `catalogo`** — requisitos, regras e critérios levantados no discovery.
> Não é o PRD (ver `prd/catalogo/prd-catalogo.md`): entra no delta quando a mudança for aberta com `gofi change new`.


# PRD — Catálogo

Segundo contexto do discovery. Depende de **Contas** (identidade do usuário e papel de administrador). É a base de **Leitura e resumos** e de **Recomendações**.

## 1. Problema e objetivo

Cada usuário que registra um livro hoje teria de informar tudo de novo. O catálogo é único e compartilhado: o livro entra uma vez, a IA define o gênero uma vez, e os demais usuários já encontram o livro pronto.

## 2. Escopo

**Dentro:** registro do livro por título e autor, identificação do mesmo livro, pesquisa de gênero pela IA em segundo plano, lista fechada de gêneros, correção do gênero.

**Fora (por enquanto):** capa, sinopse, ISBN, edição, busca por similaridade de títulos, vários gêneros por livro.

## 3. Decisões

| # | Tema | Decisão |
|---|------|---------|
| D1 | Identidade do livro | Título e autor normalizados. A tela avisa que o nome deve ser digitado exatamente igual |
| D2 | Normalização | Ignora maiúsculas, acentos, espaços duplicados e pontuação nas pontas. Todo o resto precisa ser igual |
| D3 | Pesquisa da IA | Em segundo plano, sem travar o usuário. A tela mostra o passo atual |
| D4 | Gêneros | Lista fechada. A IA só escolhe um gênero dessa lista e nunca cria um novo |
| D5 | Origem da lista | Um agente lê uma fonte confiável uma única vez e copia os gêneros. A lista é salva no sistema e revisada por um administrador antes de valer |
| D6 | Quantidade | Um único gênero por livro |
| D7 | Quem corrige o gênero | Quem registrou o livro e o administrador. Se o registrante deixar de existir, só o administrador |
| D8 | Alcance da correção | A correção vale para o livro inteiro no catálogo, para todos os usuários |
| D9 | Falha da IA | O livro fica como "gênero não identificado". O sistema tenta de novo algumas vezes e o registrante pode escolher o gênero na lista |
| D10 | Livro desconhecido | Se a IA concluir que o livro não existe, ele entra no catálogo marcado como "não verificado" |

### Fonte proposta para a lista de gêneros (D5)

Proponho como fonte as **categorias de nível mais alto da BISAC**, o padrão da indústria editorial para classificar livros. O agente reduz essa lista a cerca de 20 a 25 gêneros com nomes em português.

> **Precisa de validação antes de usar:** não confirmei os termos de uso da BISAC para copiar a lista. Se não permitirem, a alternativa é uma lista curada pelo administrador a partir de uma fonte aberta, como a Wikipédia (licença CC BY-SA). Evitar sites que proíbem coleta automática, como o Goodreads.

## 4. Requisitos funcionais

- **REQ-01 Registro.** O usuário informa título e autor. O formulário avisa que ambos devem ser digitados exatamente como o livro é conhecido, porque é isso que identifica o livro no catálogo.
- **REQ-02 Livro já existente.** O sistema normaliza título e autor (D2) e procura no catálogo. Se achar, usa o livro existente com o gênero que já tem, sem nova pesquisa. O usuário não vira dono dele.
- **REQ-03 Livro novo.** Se não achar, cria o livro no catálogo e o usuário vira o registrante dele. A pesquisa do gênero começa em segundo plano e o usuário segue usando o produto.
- **REQ-04 Estados visíveis.** O livro mostra o passo da pesquisa: "Pesquisando informações sobre o livro", "Definindo o gênero" e, ao fim, o gênero, "gênero não identificado" ou "livro não verificado". O usuário nunca vê uma tela parada sem explicação.
- **REQ-05 Lista de gêneros.** A lista inicial é montada pelo agente (D5) e fica inativa até o administrador revisar e aprovar. O administrador pode incluir, renomear e desativar gêneros.
- **REQ-06 Definição do gênero.** A IA pesquisa o livro e escolhe um único gênero da lista ativa. Se nenhum se encaixar, o resultado é "gênero não identificado".
- **REQ-07 Falha e nova tentativa.** Em caso de erro da pesquisa, o sistema tenta de novo algumas vezes. Esgotadas as tentativas, o livro fica "gênero não identificado" e o registrante pode escolher o gênero na lista.
- **REQ-08 Correção do gênero.** O registrante ou o administrador troca o gênero por outro da lista ativa. A mudança vale para o livro no catálogo todo.
- **REQ-09 Livro não verificado.** Se a IA concluir que o livro não existe, ele entra marcado como "não verificado", sem gênero. O registrante ou o administrador ainda pode definir o gênero.

## 5. Regras de negócio

- Dois registros com o mesmo título e autor normalizados são o mesmo livro. Nunca existem dois no catálogo.
- Só existe uma pesquisa por livro. Se um segundo usuário registrar o mesmo livro enquanto a pesquisa roda, ele aguarda o mesmo resultado e não abre outra.
- O gênero é do livro, não do usuário. Corrigir um gênero muda o que todos veem.
- Desativar um gênero da lista não apaga o gênero dos livros que já o usam.
- Quando o registrante apaga a conta, o livro permanece no catálogo, sem dono, e só o administrador corrige o gênero. O que o usuário escreveu (resumos e outros) segue a regra de exclusão de Contas.
- O gênero é definido pelo livro, não pelo texto dos resumos dos usuários.

## 6. Dependências

- **Contas:** identificador do usuário e o papel de administrador (D10 de Contas).
- **IA:** serviço de pesquisa e classificação, chamado em segundo plano. Falhas dele não podem derrubar o registro do livro.
- Os contextos **Leitura e resumos** e **Recomendações** leem o livro e o gênero daqui.

## 7. Critérios de aceite

- Registrar um livro que já existe, com título e autor iguais após a normalização, não cria outro livro nem dispara nova pesquisa.
- "Dom Casmurro" e "dom casmurro " são o mesmo livro. "Dom Casmurro" e "Dom Casmuro" não são.
- Durante a pesquisa o usuário vê o passo em andamento e pode continuar usando o produto.
- A IA nunca grava um gênero que não esteja na lista ativa.
- Se a IA falhar, o livro aparece como "gênero não identificado" e o registrante consegue escolher o gênero.
- Um usuário que não registrou o livro e não é administrador não consegue corrigir o gênero.
- Depois da correção feita pelo registrante ou pelo administrador, outro usuário vê o novo gênero.
- Dois usuários registrando o mesmo livro novo ao mesmo tempo resultam em um único livro e uma única pesquisa.

## 8. Questões em aberto

- Validar os termos de uso da fonte dos gêneros (BISAC) ou trocar por uma alternativa aberta.
- Número de tentativas da pesquisa e intervalo entre elas.
- Quantos gêneros terá a lista final e quem os escolhe na primeira revisão.
- Se o administrador precisa de uma tela para revisar livros "não verificados" ou isso fica fora do primeiro corte.
- Como o produto trata títulos digitados com erro: hoje o erro de digitação cria um livro novo, e a fusão de duplicados fica para depois.

## Decisões fechadas (2026-11-22)

> Superam as "Questões em aberto" acima e qualquer menção anterior em contrário.

- Fonte dos gêneros: Wikipédia (CC BY-SA), curada pelo administrador na primeira revisão; a BISAC fica descartada por não ter os termos validados.
- Pesquisa de gênero: 3 tentativas, com intervalos de 1, 5 e 15 min.
- Lista de gêneros: de 20 a 25, revisados pelo administrador.
- Livros "não verificados": sem tela de revisão no primeiro corte; o administrador age direto na base.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
