---
tipo: delta
formato: 3
contexto: recomendacoes
versao: "1.0"
status: aprovado
depende_de: [contas, catalogo, leitura-e-resumos]
keywords: [recomendação, IA, gênero, autor, motivo, resumo breve, quero ler, dispensar, pedir outras, limite diário, livro confirmado, segundo plano]
atualizado: 2026-11-22
aprovacao_versao: "1.0"
aprovacao_sha: "sha256:292b00b46ab8f847f288c0e895cf8cffedaa7c4aa1cc5e50daa9f2e3eab067c5"
---
# Delta — recomendacoes — recomendacoes

## Propósito

Recomendações sugere ao leitor o próximo livro a ler. Com base no gênero e no
autor dos livros que ele marcou como "lido", uma IA gera, em segundo plano, uma
lista pessoal de até 10 livros, cada um com um motivo e um resumo breve. O
leitor adiciona uma sugestão à estante como "quero ler", dispensa ou pede
outras, dentro de um limite diário. O contexto termina na lista: o livro e o
gênero são do Catálogo, a estante é de Leitura e resumos, e a recomendação
escolhida para o perfil é do Perfil.

## ADICIONADO

### Requisito: REQ-01 — Geração da lista de recomendações
QUANDO a lista de um usuário com ao menos um livro "lido" é gerada, o sistema
DEVE pedir à IA livros com base no gênero e no autor dos livros "lidos" dele e
montar uma lista de até 10 livros, e NÃO DEVE incluir livro que o usuário tem
na estante em qualquer estado nem livro que ele dispensou. A geração DEVE
correr em segundo plano, e o usuário DEVE ver que a lista está sendo gerada e
em que passo ela está. A primeira geração acontece quando o usuário abre as
recomendações e ainda não tem lista.

#### Cenário: REQ-01.1 — Primeira lista
- DADO um usuário com os livros "Dom Casmurro" e "Quincas Borba" marcados como "lido" e sem lista de recomendações
- QUANDO ele abre as recomendações
- ENTÃO vê que a lista está sendo gerada
- E quando a geração termina, vê até 10 livros recomendados

#### Cenário: REQ-01.2 — Livro da estante em qualquer estado fica de fora
- DADO um usuário com "Dom Casmurro" como "lido", "Helena" como "lendo" e "O Alienista" como "quero ler"
- QUANDO a IA sugere "Helena", "O Alienista" e outros livros
- ENTÃO a lista não traz "Helena" nem "O Alienista"

#### Cenário: REQ-01.3 — Livro dispensado fica de fora
- DADO um usuário que dispensou "Iaiá Garcia", de "Machado de Assis"
- QUANDO a IA sugere "iaia garcia", de "Machado de Assis", numa nova geração
- ENTÃO a lista não traz esse livro

#### Cenário: REQ-01.4 — Livro repetido na mesma resposta
- DADO uma geração em andamento
- QUANDO a IA sugere duas vezes o mesmo título e autor, com diferença só de maiúsculas ou acentos
- ENTÃO o livro aparece uma única vez na lista

### Requisito: REQ-02 — Motivo da recomendação
O sistema DEVE mostrar, em cada livro recomendado, uma frase curta, escrita
pela IA, que explica por que o livro foi sugerido.

#### Cenário: REQ-02.1 — Recomendação com motivo
- DADO uma lista de recomendações gerada
- QUANDO o usuário a consulta
- ENTÃO cada livro traz uma frase de motivo, como "porque você leu um livro de ficção científica"

### Requisito: REQ-03 — Resumo breve do livro
O sistema DEVE mostrar, em cada livro recomendado, um resumo breve do livro,
gerado pela IA, com no máximo 1.000 caracteres, que descreve o livro e NÃO DEVE
ser o resumo de nenhum usuário. O resumo breve de um livro, identificado pelo
título e pelo autor normalizados, DEVE ser gerado uma vez e reaproveitado em
toda lista, de qualquer usuário, que recomende o mesmo livro.

#### Cenário: REQ-03.1 — Recomendação com resumo breve
- DADO uma lista de recomendações gerada
- QUANDO o usuário a consulta
- ENTÃO cada livro traz um resumo breve de no máximo 1.000 caracteres

#### Cenário: REQ-03.2 — Resumo acima do limite
- DADO uma geração em andamento
- QUANDO a IA devolve, para um livro, um resumo breve com 1.001 caracteres
- ENTÃO esse livro não entra na lista

#### Cenário: REQ-03.3 — Resumo reaproveitado entre usuários
- DADO o livro "Dom Casmurro", de "Machado de Assis", já recomendado ao usuário A, com um resumo breve
- QUANDO "dom casmurro", de "Machado de Assis", é recomendado ao usuário B
- ENTÃO B vê o mesmo resumo breve que A
- E nenhum resumo novo é pedido à IA para esse livro

### Requisito: REQ-04 — Atualização quando um livro é marcado como "lido"
QUANDO o usuário marca um livro como "lido" depois da geração da lista atual,
o sistema DEVE gerar uma nova lista na próxima vez que ele abrir as
recomendações. Tirar um livro do estado "lido" ou da estante NÃO DEVE mudar a
lista atual; a mudança só aparece na próxima geração.

#### Cenário: REQ-04.1 — Novo livro "lido"
- DADO um usuário com uma lista de recomendações gerada
- QUANDO ele marca "Helena" como "lido" e depois abre as recomendações
- ENTÃO vê que uma nova lista está sendo gerada
- E quando a geração termina, vê a nova lista

#### Cenário: REQ-04.2 — Nenhum livro novo "lido"
- DADO um usuário com uma lista de recomendações gerada e nenhum livro marcado como "lido" depois dela
- QUANDO ele abre as recomendações
- ENTÃO vê a mesma lista, sem nova geração

#### Cenário: REQ-04.3 — Livro "lido" removido da estante
- DADO um usuário com uma lista de recomendações gerada
- QUANDO ele remove da estante um livro "lido" e depois abre as recomendações
- ENTÃO vê a mesma lista, sem nova geração

### Requisito: REQ-05 — Adicionar a recomendação à estante
QUANDO o usuário adiciona um livro recomendado à estante, o sistema DEVE
colocá-lo na estante dele como "quero ler" e tirá-lo da lista. O livro que já
existe no catálogo DEVE ser reconhecido pelo título e pelo autor normalizados,
sem criar outro; o livro que não existe DEVE ser registrado no catálogo com o
título e o autor da recomendação, pela regra do Catálogo, que inicia a
pesquisa de gênero e exige e-mail confirmado. Livro que já está na estante do
usuário NÃO DEVE ser duplicado nem mudar de estado.

#### Cenário: REQ-05.1 — Livro que já está no catálogo
- DADO o livro "Esaú e Jacó", de "Machado de Assis", no catálogo
- E a recomendação "esau e jaco", de "Machado de Assis", na lista do usuário
- QUANDO ele adiciona a recomendação à estante
- ENTÃO "Esaú e Jacó" entra na estante dele como "quero ler"
- E nenhum livro novo é criado no catálogo
- E a recomendação sai da lista

#### Cenário: REQ-05.2 — Livro fora do catálogo
- DADO um usuário com e-mail confirmado
- E a recomendação "Memorial de Aires", de "Machado de Assis", que não está no catálogo
- QUANDO ele adiciona a recomendação à estante
- ENTÃO o livro é registrado no catálogo com esse título e esse autor, em "Pesquisando informações sobre o livro"
- E entra na estante dele como "quero ler"
- E a recomendação sai da lista

#### Cenário: REQ-05.3 — Livro fora do catálogo sem e-mail confirmado
- DADO um usuário com o e-mail ainda não confirmado
- E uma recomendação de um livro que não está no catálogo
- QUANDO ele adiciona a recomendação à estante
- ENTÃO a operação é recusada por e-mail não confirmado
- E nenhum livro é criado no catálogo
- E a recomendação continua na lista

#### Cenário: REQ-05.4 — Livro que já está na estante
- DADO a recomendação "Helena", de "Machado de Assis", na lista do usuário
- E "Helena" colocado na estante dele como "lendo" depois da geração da lista
- QUANDO ele adiciona a recomendação à estante
- ENTÃO a estante continua com um único "Helena", como "lendo"
- E a recomendação sai da lista

### Requisito: REQ-06 — Dispensar a recomendação
QUANDO o usuário dispensa um livro recomendado, o sistema DEVE tirá-lo da lista
e NÃO DEVE recomendá-lo de novo a esse usuário.

#### Cenário: REQ-06.1 — Dispensar
- DADO a recomendação "Iaiá Garcia", de "Machado de Assis", na lista do usuário
- QUANDO ele a dispensa
- ENTÃO a recomendação sai da lista
- E nas gerações seguintes "Iaiá Garcia", de "Machado de Assis", não aparece para ele

#### Cenário: REQ-06.2 — Dispensado por um usuário, recomendado a outro
- DADO o usuário A que dispensou "Iaiá Garcia", de "Machado de Assis"
- QUANDO a IA sugere esse livro ao usuário B
- ENTÃO o livro pode entrar na lista de B

### Requisito: REQ-07 — Usuário sem livros "lidos"
SE o usuário não tem nenhum livro "lido" e não tem lista, ENTÃO o sistema DEVE
responder que não há livros "lidos", convidando-o a registrar um, sem lista e
sem chamar a IA.

#### Cenário: REQ-07.1 — Nenhum livro "lido"
- DADO um usuário com "Helena" como "lendo" e nenhum livro "lido"
- QUANDO ele abre as recomendações
- ENTÃO recebe a indicação de que não há livros "lidos", sem lista
- E nenhuma geração é iniciada

#### Cenário: REQ-07.2 — Pedir outras sem livros "lidos"
- DADO um usuário com uma lista gerada que depois removeu da estante todos os livros "lidos"
- QUANDO ele pede outras sugestões
- ENTÃO o pedido é recusado por falta de livros "lidos"
- E o pedido não conta no limite diário

### Requisito: REQ-08 — Falha da IA e nova tentativa
SE a IA falha durante a geração, ENTÃO o sistema DEVE mostrar ao usuário o
erro e permitir que ele tente de novo. A falha NÃO DEVE afetar o resto do
produto e NÃO DEVE iniciar uma nova geração sozinha.

#### Cenário: REQ-08.1 — Falha na geração
- DADO uma geração em andamento
- QUANDO a IA falha por erro do serviço
- ENTÃO o usuário, ao abrir as recomendações, vê o erro e a opção de tentar de novo

#### Cenário: REQ-08.2 — Nova tentativa com sucesso
- DADO uma geração que falhou
- QUANDO o usuário tenta de novo e a IA responde
- ENTÃO ele vê a nova lista

#### Cenário: REQ-08.3 — IA fora do ar e o resto do produto
- DADO o serviço de IA fora do ar
- QUANDO um usuário consulta um livro do catálogo ou a estante dele
- ENTÃO a consulta funciona normalmente

#### Cenário: REQ-08.4 — Abrir as recomendações depois da falha
- DADO uma geração que falhou e o usuário não pediu nova tentativa
- QUANDO ele abre as recomendações de novo
- ENTÃO vê o erro e a opção de tentar de novo
- E nenhuma geração nova é iniciada

### Requisito: REQ-09 — Pedir outras sugestões
QUANDO o usuário pede outras sugestões, o sistema DEVE gerar uma nova lista,
pela regra do REQ-01, que substitui a lista inteira, e NÃO DEVE tratar como
dispensadas as recomendações da lista anterior que o usuário não dispensou nem
adicionou à estante: elas continuam elegíveis nas gerações futuras. O sistema NÃO DEVE aceitar o pedido enquanto uma geração do
mesmo usuário está em andamento.

#### Cenário: REQ-09.1 — Nova lista substitui a anterior
- DADO um usuário com uma lista de 10 recomendações, sem nenhuma dispensada nem adicionada
- QUANDO ele pede outras sugestões e a geração termina
- ENTÃO vê só os livros da nova lista

#### Cenário: REQ-09.2 — Pedido durante uma geração
- DADO um usuário com uma geração em andamento
- QUANDO ele pede outras sugestões
- ENTÃO o pedido é recusado por geração em andamento
- E o pedido não conta no limite diário

### Requisito: REQ-10 — Limite diário de pedidos de outras sugestões
O sistema NÃO DEVE aceitar mais de 5 pedidos de outras sugestões por usuário no
mesmo dia, contado no horário de Brasília, e DEVE indicar ao usuário que o
botão de pedir outras está indisponível, com a explicação, quando o limite é
atingido. Conta no limite só o pedido de outras sugestões que gerou uma lista,
uma vez por pedido; a geração automática, a geração que falhou e a nova
tentativa depois de uma falha NÃO DEVEM contar.

#### Cenário: REQ-10.1 — Sexto pedido no dia
- DADO um usuário com 5 pedidos de outras sugestões que geraram lista hoje
- QUANDO ele pede outras sugestões
- ENTÃO o pedido é recusado por limite diário atingido
- E ele vê que o botão de pedir outras está indisponível, com a explicação

#### Cenário: REQ-10.2 — Virada do dia
- DADO um usuário que atingiu o limite às 23h50 no horário de Brasília
- QUANDO ele pede outras sugestões às 00h10 do dia seguinte, no horário de Brasília
- ENTÃO o pedido é aceito

#### Cenário: REQ-10.3 — Geração automática não conta
- DADO um usuário com 4 pedidos de outras sugestões que geraram lista hoje
- QUANDO ele marca um livro como "lido" e a nova lista é gerada ao abrir as recomendações
- ENTÃO ele ainda pode pedir outras sugestões uma vez hoje

#### Cenário: REQ-10.4 — Falha e nova tentativa contam uma vez
- DADO um usuário com 4 pedidos de outras sugestões que geraram lista hoje
- QUANDO ele pede outras sugestões, a geração falha, e ele tenta de novo com sucesso
- ENTÃO o dia fica com 5 pedidos contados
- E o próximo pedido de outras sugestões hoje é recusado por limite diário atingido

### Requisito: REQ-11 — Só livro confirmado entra na lista
O sistema DEVE confirmar que existe um livro com o título e o autor sugeridos
pela IA antes de mostrá-lo, e NÃO DEVE mostrar o livro que não for confirmado.
SE, depois de tirar os livros não confirmados, os que estão na estante e os
dispensados, a lista tem menos de 10 livros, ENTÃO o sistema DEVE pedir à IA
mais sugestões, no máximo 2 vezes por geração, e mostrar os livros que houver,
mesmo menos de 10.

#### Cenário: REQ-11.1 — Livro não confirmado
- DADO uma geração em andamento
- QUANDO a IA sugere um título e um autor que a confirmação não encontra
- ENTÃO esse livro não entra na lista

#### Cenário: REQ-11.2 — Lista completada com novas sugestões
- DADO uma geração em que 3 das 10 sugestões foram descartadas
- QUANDO a IA sugere outros livros e 3 deles são confirmados
- ENTÃO a lista tem 10 livros

#### Cenário: REQ-11.3 — Lista com menos de 10
- DADO uma geração em que, depois das 2 novas rodadas de sugestões, só 7 livros foram confirmados
- QUANDO a geração termina
- ENTÃO a lista mostra os 7 livros
- E nenhuma outra sugestão é pedida à IA nessa geração

### Requisito: REQ-12 — Recomendações pessoais
O sistema DEVE mostrar as recomendações só ao usuário dono delas, autenticado,
e NÃO DEVE permitir que outro usuário, nem o administrador, consulte, adicione
ou dispense uma recomendação que não é dele.

#### Cenário: REQ-12.1 — Recomendação de outro usuário
- DADO uma recomendação na lista da usuária A
- QUANDO o usuário B tenta adicioná-la à estante ou dispensá-la
- ENTÃO recebe a resposta de recomendação não encontrada
- E a lista de A não muda

#### Cenário: REQ-12.2 — Sem autenticação
- DADO uma requisição sem usuário autenticado
- QUANDO ela tenta abrir as recomendações
- ENTÃO é recusada por falta de autenticação
