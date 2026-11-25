---
tipo: delta
formato: 3
contexto: catalogo
versao: "1.0"
status: proposto
depende_de: [contas]
keywords: [catálogo, livro, título, autor, normalização, gênero, pesquisa em segundo plano, nova tentativa, lista de gêneros, livro não verificado, correção do gênero, administrador]
atualizado: 2026-11-22
---
# Delta — catalogo — catalogo

## Propósito

O catálogo é a lista única e compartilhada de livros do produto. Cada livro
entra uma vez, identificado pelo título e pelo autor normalizados, e recebe um
único gênero, escolhido por uma IA em segundo plano dentro de uma lista fechada
que o administrador cura. Quem registrou o livro e o administrador corrigem o
gênero, e a correção vale para todos. O catálogo termina no livro e no gênero:
estante, resumos e recomendações são de outros contextos, que só leem daqui.

## ADICIONADO

### Requisito: REQ-01 — Registro de livro novo
QUANDO um usuário registra um título e um autor que não correspondem a nenhum
livro do catálogo, o sistema DEVE criar o livro, tornar esse usuário o
registrante dele e responder sem esperar a pesquisa do gênero, com o livro em
"Pesquisando informações sobre o livro".

#### Cenário: REQ-01.1 — Livro novo entra no catálogo
- DADO um usuário com e-mail confirmado
- E nenhum livro "Dom Casmurro", de "Machado de Assis", no catálogo
- QUANDO ele registra o título "Dom Casmurro" e o autor "Machado de Assis"
- ENTÃO o livro é criado com esse título e esse autor
- E o usuário é o registrante do livro
- E o livro aparece em "Pesquisando informações sobre o livro", sem gênero

#### Cenário: REQ-01.2 — Título ou autor em branco
- DADO um usuário com e-mail confirmado
- QUANDO ele registra um livro com o título ou o autor vazio, ou só com espaços
- ENTÃO o registro é recusado por dado inválido
- E nenhum livro é criado

### Requisito: REQ-02 — Reconhecimento do livro existente
QUANDO um usuário registra um título e um autor que, normalizados, são iguais
aos de um livro do catálogo, o sistema DEVE devolver esse livro, com o estado e
o gênero que ele já tem, sem iniciar nova pesquisa e sem mudar o registrante.
Normalizar é ignorar maiúsculas, acentos, espaços repetidos e pontuação no
início e no fim do texto; todo o resto precisa ser igual.

#### Cenário: REQ-02.1 — Diferença só de maiúsculas e espaços
- DADO o livro "Dom Casmurro", de "Machado de Assis", registrado por outro usuário
- QUANDO um usuário registra o título "dom  casmurro " e o autor "machado de assis"
- ENTÃO recebe o livro já existente, com o gênero que ele tem
- E nenhuma pesquisa nova é iniciada
- E o registrante do livro continua o mesmo

#### Cenário: REQ-02.2 — Diferença só de acentos e pontuação nas pontas
- DADO o livro "Memórias Póstumas de Brás Cubas", de "Machado de Assis"
- QUANDO um usuário registra o título "Memorias Postumas de Bras Cubas." e o autor "Machado de Assis"
- ENTÃO recebe o livro já existente

#### Cenário: REQ-02.3 — Diferença de grafia é outro livro
- DADO o livro "Dom Casmurro", de "Machado de Assis"
- QUANDO um usuário registra o título "Dom Casmuro" e o autor "Machado de Assis"
- ENTÃO um livro novo é criado, com esse usuário como registrante

#### Cenário: REQ-02.4 — Pontuação no meio do texto conta
- DADO o livro "O Alienista", de "Machado de Assis"
- QUANDO um usuário registra o título "O, Alienista" e o autor "Machado de Assis"
- ENTÃO um livro novo é criado

### Requisito: REQ-03 — Livro único no catálogo
O sistema NÃO DEVE manter dois livros com o mesmo título e o mesmo autor
normalizados, nem abrir mais de uma pesquisa de gênero para o mesmo livro,
mesmo quando os registros chegam ao mesmo tempo.

#### Cenário: REQ-03.1 — Registro simultâneo do mesmo livro novo
- DADO nenhum livro "Quincas Borba", de "Machado de Assis", no catálogo
- QUANDO dois usuários registram "Quincas Borba" e "Machado de Assis" ao mesmo tempo
- ENTÃO o catálogo fica com um único livro "Quincas Borba"
- E os dois recebem esse mesmo livro
- E existe uma única pesquisa de gênero para ele

#### Cenário: REQ-03.2 — Segundo registro durante a pesquisa
- DADO o livro "Helena", de "Machado de Assis", em "Pesquisando informações sobre o livro"
- QUANDO outro usuário registra "Helena" e "Machado de Assis"
- ENTÃO recebe o mesmo livro, no passo em que a pesquisa está
- E quando a pesquisa termina, os dois veem o mesmo resultado

### Requisito: REQ-04 — Registro exige e-mail confirmado
O sistema NÃO DEVE aceitar o registro de livro de um usuário cujo e-mail não
está confirmado.

#### Cenário: REQ-04.1 — Usuário sem e-mail confirmado
- DADO um usuário autenticado com o e-mail ainda não confirmado
- QUANDO ele registra um livro
- ENTÃO o registro é recusado por e-mail não confirmado
- E nenhum livro é criado

#### Cenário: REQ-04.2 — Sem autenticação
- DADO uma requisição sem usuário autenticado
- QUANDO ela tenta registrar um livro
- ENTÃO é recusada por falta de autenticação

### Requisito: REQ-05 — Estado visível do livro
O sistema DEVE mostrar a qualquer usuário autenticado que consulta um livro o
passo em que ele está: "Pesquisando informações sobre o livro", "Definindo o
gênero", "Gênero definido", "Gênero não identificado" ou "Livro não
verificado", com o gênero quando houver. Um livro "Livro não verificado"
continua com essa marca depois de receber um gênero.

#### Cenário: REQ-05.1 — Livro em pesquisa
- DADO um livro cuja pesquisa ainda não confirmou as informações dele
- QUANDO um usuário consulta o livro
- ENTÃO vê "Pesquisando informações sobre o livro", sem gênero

#### Cenário: REQ-05.2 — Livro em classificação
- DADO um livro cujas informações já foram confirmadas e cujo gênero está sendo escolhido
- QUANDO um usuário consulta o livro
- ENTÃO vê "Definindo o gênero", sem gênero

#### Cenário: REQ-05.3 — Livro com gênero
- DADO um livro com o gênero "Romance" definido
- QUANDO um usuário que não é o registrante consulta o livro
- ENTÃO vê "Gênero definido" e o gênero "Romance"

#### Cenário: REQ-05.4 — Livro inexistente
- DADO um identificador que não corresponde a nenhum livro
- QUANDO um usuário consulta esse livro
- ENTÃO recebe a resposta de livro não encontrado

#### Cenário: REQ-05.5 — Livro não verificado com gênero
- DADO um livro em "Livro não verificado" ao qual o registrante deu o gênero "Poesia"
- QUANDO um usuário consulta o livro
- ENTÃO vê "Livro não verificado" e o gênero "Poesia"

### Requisito: REQ-06 — Gênero escolhido pela IA dentro da lista ativa
QUANDO a pesquisa confirma que o livro existe, o sistema DEVE fazer a IA
escolher um único gênero da lista ativa e gravá-lo, e NÃO DEVE gravar um
gênero que não esteja ativo na lista no momento da escolha; se nenhum se
encaixa, o livro fica em "Gênero não identificado".

#### Cenário: REQ-06.1 — IA escolhe um gênero ativo
- DADO um livro em pesquisa e o gênero "Romance" ativo na lista
- QUANDO a IA confirma o livro e escolhe "Romance"
- ENTÃO o livro fica em "Gênero definido", com o gênero "Romance"

#### Cenário: REQ-06.2 — IA devolve um gênero fora da lista ativa
- DADO um livro em pesquisa e o gênero "Ficção científica" inativo na lista
- QUANDO a IA confirma o livro e devolve "Ficção científica", ou um nome que não está na lista
- ENTÃO nenhum gênero é gravado
- E o livro fica em "Gênero não identificado"

#### Cenário: REQ-06.3 — Nenhum gênero se encaixa
- DADO um livro em pesquisa
- QUANDO a IA confirma o livro e responde que nenhum gênero da lista se encaixa
- ENTÃO o livro fica em "Gênero não identificado", sem gênero

### Requisito: REQ-07 — Nova tentativa quando a pesquisa falha
SE a pesquisa do gênero falha por erro, ENTÃO o sistema DEVE tentar de novo até
3 vezes, 1, 5 e 15 minutos depois de cada falha, mantendo o livro no passo
"Pesquisando informações sobre o livro"; esgotadas as tentativas, o livro DEVE
ficar em "Gênero não identificado". A falha da pesquisa NÃO DEVE desfazer nem
impedir o registro do livro.

#### Cenário: REQ-07.1 — Falha seguida de sucesso
- DADO um livro cuja primeira pesquisa falhou por erro do serviço de IA
- QUANDO a nova tentativa, 1 minuto depois, confirma o livro e escolhe "Romance", ativo
- ENTÃO o livro fica em "Gênero definido", com "Romance"

#### Cenário: REQ-07.2 — Tentativas esgotadas
- DADO um livro cuja pesquisa falhou por erro
- QUANDO as 3 novas tentativas, 1, 5 e 15 minutos depois de cada falha, também falham
- ENTÃO o livro fica em "Gênero não identificado"
- E nenhuma outra tentativa é feita

#### Cenário: REQ-07.3 — Serviço de IA indisponível no registro
- DADO o serviço de IA fora do ar
- QUANDO um usuário com e-mail confirmado registra um livro novo
- ENTÃO o livro é criado e devolvido em "Pesquisando informações sobre o livro"

#### Cenário: REQ-07.4 — Entre uma tentativa e outra
- DADO um livro cuja primeira pesquisa falhou há 30 segundos
- QUANDO um usuário consulta o livro
- ENTÃO vê "Pesquisando informações sobre o livro"

### Requisito: REQ-08 — Livro não verificado
QUANDO a IA conclui que o livro não existe, o sistema DEVE manter o livro no
catálogo, marcado como "Livro não verificado" e sem gênero até que o
registrante ou o administrador escolha um.

#### Cenário: REQ-08.1 — IA não encontra o livro
- DADO um livro em pesquisa
- QUANDO a IA conclui que não existe livro com esse título e esse autor
- ENTÃO o livro continua no catálogo, em "Livro não verificado"
- E nenhum gênero é gravado

### Requisito: REQ-09 — Correção do gênero pelo registrante ou pelo administrador
QUANDO o registrante do livro ou um administrador escolhe um gênero ativo da
lista para o livro, o sistema DEVE gravar esse gênero no livro, que passa a ser
o gênero que todos os usuários veem. O sistema NÃO DEVE aceitar a troca
enquanto o livro está em "Pesquisando informações sobre o livro" ou em
"Definindo o gênero".

#### Cenário: REQ-09.1 — Registrante corrige o gênero
- DADO um livro em "Gênero definido", com "Romance", registrado pela usuária A
- E o gênero "Conto" ativo na lista
- QUANDO A troca o gênero do livro para "Conto"
- ENTÃO o usuário B, ao consultar o livro, vê o gênero "Conto"

#### Cenário: REQ-09.2 — Administrador corrige o gênero
- DADO um livro com o gênero "Romance", registrado pela usuária A
- QUANDO um administrador troca o gênero para "Conto"
- ENTÃO o livro passa a ter o gênero "Conto"

#### Cenário: REQ-09.3 — Registrante escolhe o gênero de um livro não identificado
- DADO um livro em "Gênero não identificado", registrado pela usuária A
- QUANDO A escolhe o gênero "Poesia", ativo na lista
- ENTÃO o livro fica em "Gênero definido", com "Poesia"

#### Cenário: REQ-09.4 — Gênero inativo ou inexistente
- DADO um livro registrado pela usuária A e o gênero "Epopeia" inativo
- QUANDO A escolhe "Epopeia", ou um gênero que não existe
- ENTÃO a correção é recusada
- E o gênero do livro não muda

#### Cenário: REQ-09.5 — Troca durante a pesquisa
- DADO um livro registrado pela usuária A, em "Pesquisando informações sobre o livro" ou em "Definindo o gênero"
- E o gênero "Conto" ativo na lista
- QUANDO A, ou um administrador, escolhe "Conto"
- ENTÃO a troca é recusada por pesquisa em andamento
- E o livro continua no mesmo passo, sem gênero

#### Cenário: REQ-09.6 — Gênero para um livro não verificado
- DADO um livro em "Livro não verificado", registrado pela usuária A
- QUANDO A escolhe o gênero "Poesia", ativo na lista
- ENTÃO o livro continua em "Livro não verificado", com o gênero "Poesia"

### Requisito: REQ-10 — Correção restrita ao registrante e ao administrador
O sistema NÃO DEVE aceitar a troca do gênero de um livro por um usuário que não
é o registrante dele nem administrador; um livro sem registrante só tem o
gênero trocado pelo administrador.

#### Cenário: REQ-10.1 — Outro usuário tenta corrigir
- DADO um livro registrado pela usuária A
- QUANDO o usuário B, que não é administrador, tenta trocar o gênero
- ENTÃO a troca é recusada por falta de permissão
- E o gênero do livro não muda

#### Cenário: REQ-10.2 — Registrante excluiu a conta
- DADO um livro registrado pela usuária A
- QUANDO a conta de A é excluída
- ENTÃO o livro continua no catálogo, com o mesmo gênero e sem registrante

#### Cenário: REQ-10.3 — Livro sem registrante
- DADO um livro sem registrante
- QUANDO um usuário que não é administrador tenta trocar o gênero
- ENTÃO a troca é recusada por falta de permissão
- E quando um administrador troca o gênero, a troca é aceita

### Requisito: REQ-11 — Lista inicial de gêneros inativa até a aprovação
O sistema DEVE trazer a lista inicial de 20 a 25 gêneros, tirada da Wikipédia,
com todos os gêneros inativos, e só um administrador DEVE poder ativá-los, um a
um, ao aprová-los.

#### Cenário: REQ-11.1 — Sistema recém-instalado
- DADO o sistema recém-instalado, sem nenhuma aprovação do administrador
- QUANDO um usuário pede a lista de gêneros ativos
- ENTÃO recebe a lista vazia
- E o administrador, ao consultar a lista completa, vê de 20 a 25 gêneros inativos

#### Cenário: REQ-11.2 — Administrador aprova um gênero
- DADO o gênero "Romance" inativo
- QUANDO um administrador o ativa
- ENTÃO "Romance" passa a aparecer na lista de gêneros ativos

### Requisito: REQ-12 — Manutenção da lista de gêneros pelo administrador
O sistema DEVE permitir ao administrador incluir um gênero, que já nasce ativo,
renomear um gênero e desativá-lo; e NÃO DEVE aceitar dois gêneros com o mesmo
nome, sem diferenciar maiúsculas e acentos, nem essas operações vindas de quem
não é administrador.

#### Cenário: REQ-12.1 — Incluir gênero
- DADO um administrador e nenhum gênero "Distopia" na lista
- QUANDO ele inclui "Distopia"
- ENTÃO "Distopia" aparece na lista de gêneros ativos

#### Cenário: REQ-12.2 — Renomear gênero
- DADO o gênero "Sci-fi", usado por um livro
- QUANDO um administrador o renomeia para "Ficção científica"
- ENTÃO o livro passa a mostrar o gênero "Ficção científica"

#### Cenário: REQ-12.3 — Nome repetido
- DADO o gênero "Romance" na lista
- QUANDO um administrador inclui "romance", ou renomeia outro gênero para "Romance"
- ENTÃO a operação é recusada por nome já existente

#### Cenário: REQ-12.4 — Quem não é administrador
- DADO um usuário que não é administrador
- QUANDO ele tenta incluir, renomear, ativar ou desativar um gênero, ou consultar a lista completa
- ENTÃO a operação é recusada por falta de permissão

### Requisito: REQ-13 — Desativar gênero preserva os livros
QUANDO um administrador desativa um gênero, o sistema DEVE tirá-lo da lista
ativa e NÃO DEVE mudar o gênero dos livros que já o usam.

#### Cenário: REQ-13.1 — Gênero desativado em uso
- DADO o gênero "Romance" ativo e um livro com o gênero "Romance"
- QUANDO um administrador desativa "Romance"
- ENTÃO "Romance" sai da lista de gêneros ativos
- E o livro continua com o gênero "Romance"

### Requisito: REQ-14 — Lista de gêneros ativos
O sistema DEVE mostrar a qualquer usuário autenticado a lista de gêneros
ativos, em ordem alfabética, sem os inativos.

#### Cenário: REQ-14.1 — Consulta da lista ativa
- DADO os gêneros "Romance" e "Conto" ativos e "Epopeia" inativo
- QUANDO um usuário autenticado pede a lista de gêneros
- ENTÃO recebe "Conto" e "Romance", nessa ordem
