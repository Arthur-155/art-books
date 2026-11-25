---
tipo: delta
formato: 3
contexto: contas
versao: "1.0"
status: aprovado
keywords: [cadastro, login, sessão, confirmação de e-mail, recuperação de senha, troca de senha, troca de e-mail, exclusão de conta, limite de tentativas, administrador, username, privacidade]
atualizado: 2026-11-24
aprovacao_versao: "1.0"
aprovacao_sha: "sha256:d7d8c912c5d1cc1020ce66e1abcbbd918d52cd3ebdd424a0da55233ec86448ab"
---
# Delta — contas — contas

## Propósito

Contas é a identidade do leitor na rede social de leitura: cadastro, login,
sessão e recuperação de acesso com e-mail e senha, confirmação e troca de
e-mail, exclusão da conta e o papel de administrador. Os outros contextos
dependem dele pelo id do usuário, pelo username, pelo estado "e-mail
confirmado" e pelo papel. Termina na identidade: perfil público, amizades e
publicações são de outros contextos.

## ADICIONADO

### Requisito: REQ-01 — Cadastro
O sistema DEVE criar uma conta com e-mail, senha, nome de exibição e username QUANDO a pessoa aceita os termos e a política de privacidade e declara ter 13 anos ou mais, com e-mail em endereço simples válido (`local@dominio`, sem nome nem comentário) de até 254 caracteres e username únicos sem diferenciar maiúsculas de minúsculas, senha de 8 a 128 pontos de código Unicode, considerada inteira, username de 3 a 30 caracteres só com letras de a a z, algarismos e sublinhado, guardado em minúsculas, e nome de exibição de 1 a 60 pontos de código depois de aparados os espaços das pontas.

#### Cenário: REQ-01.1 — Cadastro válido
- DADO que não existe conta com o e-mail `ana@exemplo.com` nem com o username `ana_leitora`
- QUANDO a pessoa se cadastra com esse e-mail, a senha `segredo123`, o nome "Ana", o username `Ana_Leitora`, o aceite dos termos e a declaração de 13 anos ou mais
- ENTÃO a conta é criada com o papel usuário e o e-mail não confirmado
- E o username fica guardado como `ana_leitora`

#### Cenário: REQ-01.2 — E-mail já cadastrado em outra caixa
- DADO uma conta com o e-mail `ana@exemplo.com`
- QUANDO outra pessoa se cadastra com `ANA@Exemplo.com`
- ENTÃO o cadastro é recusado com a mensagem de que o e-mail já está em uso
- E nenhuma conta nova é criada

#### Cenário: REQ-01.3 — Username já em uso em outra caixa
- DADO uma conta com o username `ana_leitora`
- QUANDO outra pessoa se cadastra com o username `ANA_LEITORA`
- ENTÃO o cadastro é recusado com a mensagem de que o username já está em uso

#### Cenário: REQ-01.4 — Sem aceite dos termos ou sem declaração de idade
- DADO um cadastro com os demais dados válidos
- QUANDO a pessoa não aceita os termos ou não declara ter 13 anos ou mais
- ENTÃO o cadastro é recusado como dado inválido

#### Cenário: REQ-01.5 — Senha curta
- DADO um cadastro com os demais dados válidos
- QUANDO a senha tem 7 caracteres
- ENTÃO o cadastro é recusado como dado inválido

#### Cenário: REQ-01.6 — Username fora da regra
- DADO um cadastro com os demais dados válidos
- QUANDO o username é `an`, tem 31 caracteres ou contém `ana.leitora`
- ENTÃO o cadastro é recusado como dado inválido

#### Cenário: REQ-01.7 — Nome de exibição fora do limite
- DADO um cadastro com os demais dados válidos
- QUANDO o nome de exibição é vazio, só espaços ou tem 61 caracteres
- ENTÃO o cadastro é recusado como dado inválido

#### Cenário: REQ-01.8 — Cadastro simultâneo com o mesmo e-mail
- DADO que não existe conta com o e-mail `ana@exemplo.com`
- QUANDO dois cadastros com esse e-mail chegam ao mesmo tempo
- ENTÃO só um cria a conta
- E o outro é recusado com a mensagem de que o e-mail já está em uso

#### Cenário: REQ-01.9 — E-mail malformado
- DADO um cadastro com os demais dados válidos
- QUANDO o e-mail é `ana@` ou tem mais de 254 caracteres
- ENTÃO o cadastro é recusado como dado inválido

#### Cenário: REQ-01.10 — Senha acima do máximo
- DADO um cadastro com os demais dados válidos
- QUANDO a senha tem 129 caracteres
- ENTÃO o cadastro é recusado como dado inválido

### Requisito: REQ-02 — Confirmação do e-mail
QUANDO uma conta é criada, o sistema DEVE enviar ao e-mail dela um link de confirmação de uso único que vale por 24 horas, e DEVE marcar o e-mail como confirmado QUANDO o link é usado dentro do prazo.

#### Cenário: REQ-02.1 — Envio no cadastro
- DADO um cadastro válido
- QUANDO a conta é criada
- ENTÃO um e-mail com o link de confirmação é enviado ao endereço cadastrado

#### Cenário: REQ-02.2 — Link usado no prazo
- DADO um link de confirmação emitido há 23 horas e nunca usado
- QUANDO a pessoa o usa
- ENTÃO o e-mail da conta passa a constar como confirmado

#### Cenário: REQ-02.3 — Link vencido
- DADO um link de confirmação emitido há 25 horas
- QUANDO a pessoa o usa
- ENTÃO a confirmação é recusada como link inválido ou vencido
- E o e-mail continua não confirmado

#### Cenário: REQ-02.4 — Link usado duas vezes
- DADO um link de confirmação já usado
- QUANDO a pessoa o usa de novo
- ENTÃO a confirmação é recusada como link inválido ou vencido

### Requisito: REQ-03 — Reenvio da confirmação
QUANDO o usuário logado com e-mail não confirmado pede o reenvio, o sistema DEVE enviar um link de confirmação novo e DEVE invalidar o link anterior.

#### Cenário: REQ-03.1 — Reenvio
- DADO um usuário logado com e-mail não confirmado e um link de confirmação ainda no prazo
- QUANDO ele pede o reenvio
- ENTÃO um link novo é enviado
- E o link anterior passa a ser recusado como inválido

#### Cenário: REQ-03.2 — E-mail já confirmado
- DADO um usuário logado com e-mail confirmado
- QUANDO ele pede o reenvio da confirmação
- ENTÃO o pedido é recusado informando que o e-mail já está confirmado
- E nenhum e-mail é enviado

### Requisito: REQ-04 — Sessão informa o e-mail confirmado
O sistema DEVE informar, em toda sessão que emite ou renova, se o e-mail do usuário está confirmado, para que os outros contextos impeçam de publicar quem não confirmou.

#### Cenário: REQ-04.1 — Usuário sem e-mail confirmado
- DADO um usuário com e-mail não confirmado
- QUANDO ele entra
- ENTÃO a sessão informa o e-mail como não confirmado
- E o usuário navega normalmente

#### Cenário: REQ-04.2 — Confirmação vale na renovação
- DADO um usuário logado cuja sessão informa o e-mail como não confirmado
- QUANDO ele confirma o e-mail e renova a sessão
- ENTÃO a sessão renovada informa o e-mail como confirmado

### Requisito: REQ-05 — Login
O sistema DEVE abrir uma sessão QUANDO a pessoa informa um e-mail cadastrado, sem diferenciar maiúsculas de minúsculas, e a senha correta, e NÃO DEVE revelar, numa recusa, se o e-mail existe.

#### Cenário: REQ-05.1 — Login válido
- DADO uma conta com o e-mail `ana@exemplo.com` e a senha `segredo123`
- QUANDO a pessoa entra com `Ana@Exemplo.com` e `segredo123`
- ENTÃO uma sessão é aberta com o id do usuário e o papel dele

#### Cenário: REQ-05.2 — Senha errada e e-mail inexistente respondem igual
- DADO uma conta com o e-mail `ana@exemplo.com`
- QUANDO a pessoa entra com `ana@exemplo.com` e a senha errada, e depois com `ninguem@exemplo.com` e qualquer senha
- ENTÃO as duas tentativas são recusadas com a mesma resposta e a mesma mensagem

### Requisito: REQ-06 — Bloqueio por tentativas erradas
SE houver 5 tentativas de login erradas numa janela de 15 minutos contada a partir do primeiro erro, para o mesmo e-mail ou para o mesmo IP, ENTÃO o sistema DEVE recusar todo login desse e-mail ou desse IP por 15 minutos contados da quinta tentativa, mesmo com a senha correta, sem revelar se o e-mail existe. Login correto zera a contagem do e-mail, e tentativa recusada por bloqueio não conta como erro nem prorroga o bloqueio. As senhas erradas na troca de senha e na exclusão da conta ficam fora desse bloqueio.

#### Cenário: REQ-06.1 — Bloqueio por e-mail
- DADO 5 tentativas erradas para `ana@exemplo.com` nos últimos 15 minutos, vindas de IPs diferentes
- QUANDO a pessoa entra com `ana@exemplo.com` e a senha correta
- ENTÃO o login é recusado como bloqueado temporariamente

#### Cenário: REQ-06.2 — Bloqueio por IP
- DADO 5 tentativas erradas de um mesmo IP nos últimos 15 minutos, para e-mails diferentes
- QUANDO desse IP alguém entra com outro e-mail e a senha correta
- ENTÃO o login é recusado como bloqueado temporariamente

#### Cenário: REQ-06.3 — Fim do bloqueio
- DADO um e-mail bloqueado pela quinta tentativa errada há 16 minutos
- QUANDO a pessoa entra com esse e-mail e a senha correta, de um IP sem bloqueio
- ENTÃO a sessão é aberta

#### Cenário: REQ-06.4 — E-mail inexistente também bloqueia
- DADO 5 tentativas erradas para `ninguem@exemplo.com` nos últimos 15 minutos
- QUANDO mais uma tentativa chega para esse e-mail
- ENTÃO a recusa é a mesma de um e-mail cadastrado bloqueado

#### Cenário: REQ-06.5 — Erros espaçados não bloqueiam
- DADO 4 tentativas erradas para `ana@exemplo.com`, a primeira há mais de 15 minutos
- QUANDO a pessoa erra mais uma vez e depois entra com a senha correta
- ENTÃO a sessão é aberta

#### Cenário: REQ-06.6 — Tentativa durante o bloqueio
- DADO um e-mail bloqueado há 10 minutos
- QUANDO a pessoa tenta entrar com ele mais 3 vezes
- ENTÃO as tentativas são recusadas como bloqueio temporário
- E o bloqueio termina 15 minutos depois da quinta tentativa errada, sem prorrogação

#### Cenário: REQ-06.7 — Login correto zera a contagem
- DADO 4 tentativas erradas para `ana@exemplo.com` há 2 minutos
- QUANDO a pessoa entra com a senha correta e, depois, erra 4 vezes
- ENTÃO nenhum bloqueio é aplicado

### Requisito: REQ-07 — Sessão longa com renovação
O sistema DEVE manter a sessão válida por 30 dias e DEVE renová-la por mais 30 dias a cada renovação, trazendo na sessão renovada o papel e o estado do e-mail atuais; NÃO DEVE aceitar a mesma credencial de renovação duas vezes. Encerrar uma sessão impede só a renovação: o token de acesso já emitido vale até vencer, no prazo padrão do SDK.

#### Cenário: REQ-07.1 — Renovação no prazo
- DADO uma sessão aberta há 29 dias
- QUANDO o usuário a renova
- ENTÃO recebe uma sessão nova, válida por mais 30 dias

#### Cenário: REQ-07.2 — Sessão vencida
- DADO uma sessão aberta há 31 dias e nunca renovada
- QUANDO o usuário tenta renová-la
- ENTÃO a renovação é recusada e ele precisa entrar de novo

#### Cenário: REQ-07.3 — Credencial de renovação já usada
- DADO uma credencial de renovação já usada numa renovação
- QUANDO alguém a usa de novo
- ENTÃO a renovação é recusada
- E a sessão que a renovação anterior gerou continua valendo

#### Cenário: REQ-07.4 — Sessão encerrada não renova
- DADO uma sessão encerrada pela troca de senha
- QUANDO o usuário tenta renová-la
- ENTÃO a renovação é recusada

### Requisito: REQ-08 — Logout
QUANDO o usuário sai, o sistema DEVE encerrar só a sessão do dispositivo atual.

#### Cenário: REQ-08.1 — Sair num dispositivo
- DADO um usuário com sessões abertas em dois dispositivos
- QUANDO ele sai no primeiro
- ENTÃO a sessão do primeiro não renova mais
- E a sessão do segundo continua renovando

#### Cenário: REQ-08.2 — Sair sem sessão válida
- DADO uma credencial de renovação inexistente ou já encerrada
- QUANDO o usuário pede para sair
- ENTÃO o pedido é aceito sem erro, como se tivesse saído

### Requisito: REQ-09 — Recuperação de senha
QUANDO a pessoa pede a recuperação com um e-mail, o sistema DEVE dar a mesma resposta exista ou não a conta, inclusive na recusa por excesso de pedidos, DEVE enviar à conta existente um link de uso único que vale por 30 minutos e invalida o link de recuperação anterior, e QUANDO o link é usado no prazo com uma senha nova de 8 a 128 caracteres, DEVE trocar a senha e encerrar todas as sessões do usuário.

#### Cenário: REQ-09.1 — Pedido com e-mail cadastrado
- DADO uma conta com o e-mail `ana@exemplo.com`
- QUANDO a pessoa pede a recuperação com `ana@exemplo.com`
- ENTÃO um link de recuperação é enviado a esse e-mail
- E a resposta diz que, se o e-mail estiver cadastrado, um link foi enviado

#### Cenário: REQ-09.2 — Pedido com e-mail inexistente
- DADO que não existe conta com `ninguem@exemplo.com`
- QUANDO a pessoa pede a recuperação com esse e-mail
- ENTÃO a resposta é a mesma do REQ-09.1
- E nenhum e-mail é enviado

#### Cenário: REQ-09.3 — Nova senha pelo link no prazo
- DADO um link de recuperação emitido há 29 minutos e uma sessão aberta em outro dispositivo
- QUANDO a pessoa usa o link com a senha nova `outrasenha1`
- ENTÃO a senha passa a ser `outrasenha1`
- E a sessão do outro dispositivo não renova mais

#### Cenário: REQ-09.4 — Link vencido
- DADO um link de recuperação emitido há 31 minutos
- QUANDO a pessoa o usa com uma senha nova válida
- ENTÃO o pedido é recusado como link inválido ou vencido
- E a senha não muda

#### Cenário: REQ-09.5 — Link usado duas vezes
- DADO um link de recuperação já usado
- QUANDO a pessoa o usa de novo
- ENTÃO o pedido é recusado como link inválido ou vencido

#### Cenário: REQ-09.6 — Senha nova curta
- DADO um link de recuperação no prazo
- QUANDO a pessoa o usa com uma senha nova de 7 caracteres
- ENTÃO o pedido é recusado como dado inválido
- E o link continua valendo

#### Cenário: REQ-09.7 — Recusa por excesso igual para e-mail inexistente
- DADO um pedido de recuperação há 30 segundos para `ana@exemplo.com` e outro para `ninguem@exemplo.com`
- QUANDO chega outro pedido para cada um dos dois e-mails
- ENTÃO os dois são recusados como excesso de pedidos, com a mesma resposta

#### Cenário: REQ-09.8 — Novo pedido invalida o link anterior
- DADO um link de recuperação ainda no prazo
- QUANDO a pessoa faz um novo pedido de recuperação, passado 1 minuto
- ENTÃO o link anterior é recusado como inválido ou vencido

### Requisito: REQ-10 — Troca de senha
QUANDO o usuário logado informa a senha atual correta e uma senha nova de 8 a 128 caracteres, o sistema DEVE trocar a senha e encerrar todas as sessões dele.

#### Cenário: REQ-10.1 — Troca com a senha atual
- DADO um usuário com sessões abertas em dois dispositivos
- QUANDO ele troca a senha no primeiro informando a senha atual correta
- ENTÃO a senha nova passa a valer no login
- E nenhuma das duas sessões renova mais

#### Cenário: REQ-10.2 — Senha atual errada
- DADO um usuário logado
- QUANDO ele tenta trocar a senha informando a senha atual errada
- ENTÃO a troca é recusada
- E a senha e as sessões não mudam

### Requisito: REQ-11 — Troca de e-mail
QUANDO o usuário logado pede a troca para um e-mail em formato válido que nenhuma conta usa, o dele incluído, o sistema DEVE enviar ao e-mail atual um link de aprovação que vale por 1 hora; aprovado no prazo, DEVE enviar ao novo e-mail um link de confirmação que vale por 24 horas; confirmado no prazo, DEVE trocar o e-mail, mantê-lo como confirmado e encerrar todas as sessões do usuário. Enquanto a troca não é confirmada, o e-mail atual DEVE continuar valendo para o login.

#### Cenário: REQ-11.1 — Pedido envia a aprovação ao e-mail atual
- DADO um usuário logado com o e-mail `ana@exemplo.com`
- QUANDO ele pede a troca para `ana@novo.com`
- ENTÃO um link de aprovação é enviado a `ana@exemplo.com`
- E nada é enviado a `ana@novo.com`

#### Cenário: REQ-11.2 — Aprovação envia a confirmação ao novo e-mail
- DADO um link de aprovação da troca emitido há 50 minutos
- QUANDO a pessoa o usa
- ENTÃO um link de confirmação é enviado a `ana@novo.com`
- E o e-mail da conta ainda é `ana@exemplo.com`

#### Cenário: REQ-11.11 — Link de aprovação já usado
- DADO um link de aprovação da troca já usado
- QUANDO a pessoa o usa de novo
- ENTÃO a aprovação é recusada como link inválido ou vencido
- E nenhum e-mail é enviado

#### Cenário: REQ-11.12 — Novo e-mail tomado antes da aprovação
- DADO um pedido de troca para `ana@novo.com` e, depois, outra conta passa a usar `ana@novo.com`
- QUANDO a pessoa usa o link de aprovação
- ENTÃO a aprovação é recusada com a mensagem de que o e-mail já está em uso
- E nenhum e-mail é enviado

#### Cenário: REQ-11.3 — Confirmação conclui a troca
- DADO uma troca aprovada, um link de confirmação emitido há 23 horas e uma sessão aberta em outro dispositivo
- QUANDO a pessoa usa o link de confirmação
- ENTÃO o e-mail da conta passa a ser `ana@novo.com`, confirmado
- E o login passa a aceitar `ana@novo.com` e recusar `ana@exemplo.com`
- E a sessão do outro dispositivo não renova mais
- E um link de recuperação pendente é recusado como inválido ou vencido

#### Cenário: REQ-11.4 — Sem aprovação a troca não avança
- DADO um pedido de troca cujo link de aprovação venceu, emitido há 61 minutos
- QUANDO a pessoa o usa
- ENTÃO a aprovação é recusada como link inválido ou vencido
- E nenhum e-mail é enviado ao novo endereço

#### Cenário: REQ-11.5 — E-mail atual vale durante a troca
- DADO uma troca pedida e ainda não confirmada
- QUANDO o usuário entra com o e-mail atual e a senha
- ENTÃO a sessão é aberta

#### Cenário: REQ-11.6 — Novo e-mail já em uso
- DADO uma conta com o e-mail `beto@exemplo.com`
- QUANDO outro usuário pede a troca para `Beto@Exemplo.com`
- ENTÃO o pedido é recusado com a mensagem de que o e-mail já está em uso

#### Cenário: REQ-11.7 — Novo e-mail tomado antes da confirmação
- DADO uma troca aprovada para `ana@novo.com`
- QUANDO outra conta passa a usar `ana@novo.com` e a pessoa usa o link de confirmação
- ENTÃO a confirmação é recusada com a mensagem de que o e-mail já está em uso
- E o e-mail da conta não muda
- E o link de confirmação continua valendo até o prazo

#### Cenário: REQ-11.8 — Pedido novo substitui o pendente
- DADO uma troca pendente para `ana@novo.com`
- QUANDO o usuário pede a troca para `ana@outro.com`
- ENTÃO os links da troca para `ana@novo.com` passam a ser recusados como inválidos

#### Cenário: REQ-11.9 — Troca para o próprio e-mail
- DADO um usuário com o e-mail `ana@exemplo.com`
- QUANDO ele pede a troca para `Ana@Exemplo.com`
- ENTÃO o pedido é recusado com a mensagem de que o e-mail já está em uso
- E nenhum e-mail é enviado

#### Cenário: REQ-11.10 — Novo e-mail malformado
- DADO um usuário logado
- QUANDO ele pede a troca para `ana@`
- ENTÃO o pedido é recusado como dado inválido

### Requisito: REQ-12 — Exclusão da conta
QUANDO o usuário logado confirma com a senha correta, o sistema DEVE excluir a conta na hora, com os dados dela e todas as sessões, deixando o e-mail e o username livres para um novo cadastro.

#### Cenário: REQ-12.1 — Exclusão com a senha
- DADO um usuário `ana_leitora` com o e-mail `ana@exemplo.com` e sessões abertas em dois dispositivos
- QUANDO ele exclui a conta informando a senha correta
- ENTÃO o login com `ana@exemplo.com` é recusado como credencial inválida
- E nenhuma das sessões renova mais
- E um link de confirmação ou de recuperação pendente dessa conta é recusado como inválido

#### Cenário: REQ-12.2 — Senha errada
- DADO um usuário logado
- QUANDO ele pede a exclusão com a senha errada
- ENTÃO a exclusão é recusada
- E a conta continua existindo

#### Cenário: REQ-12.3 — E-mail e username livres depois da exclusão
- DADO que a conta `ana_leitora` com o e-mail `ana@exemplo.com` foi excluída
- QUANDO alguém se cadastra com esse e-mail e esse username
- ENTÃO a conta nova é criada

### Requisito: REQ-13 — Limite de pedidos de e-mail
O sistema NÃO DEVE aceitar, para o mesmo e-mail ou o mesmo IP, mais de 1 pedido por minuto nem mais de 5 por hora de cada tipo, contando só os pedidos aceitos, isto é, os que passaram pelo limite e tiveram o envio tentado, mesmo que o envio falhe, e inclusive o de recuperação para e-mail sem conta. A chave por e-mail é o informado na recuperação e o e-mail atual da conta no reenvio e na troca de e-mail: reenvio da confirmação, recuperação de senha e troca de e-mail. O cadastro e o uso do link de aprovação da troca ficam fora desse limite.

#### Cenário: REQ-13.1 — Segundo pedido no mesmo minuto
- DADO um pedido de recuperação para `ana@exemplo.com` há 30 segundos
- QUANDO chega outro pedido de recuperação para esse e-mail, de outro IP
- ENTÃO o pedido é recusado como excesso de pedidos
- E nenhum e-mail é enviado

#### Cenário: REQ-13.2 — Sexto pedido na hora
- DADO 5 pedidos de reenvio da confirmação do mesmo usuário na última hora, com mais de 1 minuto entre eles
- QUANDO ele pede o reenvio de novo
- ENTÃO o pedido é recusado como excesso de pedidos

#### Cenário: REQ-13.3 — Limite por IP
- DADO um pedido de recuperação de um IP há 30 segundos
- QUANDO desse IP chega um pedido de recuperação para outro e-mail
- ENTÃO o pedido é recusado como excesso de pedidos

#### Cenário: REQ-13.4 — Tipos diferentes não se somam
- DADO um pedido de reenvio da confirmação do usuário há 30 segundos
- QUANDO ele pede a troca de e-mail
- ENTÃO o pedido de troca é aceito

#### Cenário: REQ-13.6 — Pedido recusado não conta
- DADO um pedido de recuperação aceito há 70 segundos e outro recusado há 30 segundos
- QUANDO chega outro pedido de recuperação para o mesmo e-mail
- ENTÃO o pedido é aceito

#### Cenário: REQ-13.5 — Depois do intervalo
- DADO um pedido de recuperação para `ana@exemplo.com` há 61 segundos e nenhum outro na hora
- QUANDO chega outro pedido de recuperação para esse e-mail
- ENTÃO o pedido é aceito

### Requisito: REQ-14 — Falha no envio do e-mail
SE o envio de um e-mail de confirmação, de reenvio ou da troca de e-mail falhar, ENTÃO o sistema DEVE informar ao usuário que o e-mail não foi enviado e DEVE permitir pedir o envio de novo, respeitado o REQ-13; na recuperação de senha, a resposta continua a mesma e a falha só é registrada.

#### Cenário: REQ-14.1 — Falha no cadastro
- DADO um cadastro válido e o serviço de e-mail fora do ar
- QUANDO a conta é criada
- ENTÃO a conta existe
- E a resposta informa que o e-mail de confirmação não foi enviado

#### Cenário: REQ-14.2 — Novo pedido depois da falha
- DADO um reenvio da confirmação que falhou no envio
- QUANDO o usuário pede o reenvio de novo, passado 1 minuto, com o serviço de e-mail de volta
- ENTÃO o link novo é enviado

#### Cenário: REQ-14.3 — Falha na aprovação da troca de e-mail
- DADO um link de aprovação da troca no prazo e o serviço de e-mail fora do ar
- QUANDO a pessoa usa o link
- ENTÃO a resposta informa que o e-mail de confirmação não foi enviado
- E o link de aprovação continua valendo até o prazo

#### Cenário: REQ-14.4 — Nova tentativa da aprovação depois da falha
- DADO um link de aprovação cujo uso falhou no envio, com o serviço de e-mail de volta
- QUANDO a pessoa usa o mesmo link de novo, dentro do prazo
- ENTÃO o link de confirmação é enviado ao novo e-mail

#### Cenário: REQ-14.5 — Falha no envio da recuperação
- DADO uma conta com o e-mail `ana@exemplo.com` e o serviço de e-mail fora do ar
- QUANDO a pessoa pede a recuperação com esse e-mail
- ENTÃO a resposta é a mesma do REQ-09.1
- E o pedido conta no limite do REQ-13

### Requisito: REQ-15 — Papéis
Toda conta DEVE nascer com o papel usuário, e o papel de administrador NÃO DEVE ser concedido nem retirado por nenhuma operação da aplicação, só pelo comando operacional, valendo a partir da próxima sessão aberta ou renovada.

#### Cenário: REQ-15.1 — Concessão pelo comando operacional
- DADO uma conta com o papel usuário e o e-mail `ana@exemplo.com`
- QUANDO o operador concede o papel de administrador a `ana@exemplo.com` pelo comando operacional e o usuário renova a sessão
- ENTÃO a sessão renovada traz o papel administrador

#### Cenário: REQ-15.2 — Retirada pelo comando operacional
- DADO uma conta administradora
- QUANDO o operador retira o papel pelo comando operacional e o usuário renova a sessão
- ENTÃO a sessão renovada traz o papel usuário

#### Cenário: REQ-15.3 — E-mail inexistente no comando
- DADO que não existe conta com `ninguem@exemplo.com`
- QUANDO o operador concede o papel a esse e-mail
- ENTÃO o comando falha informando que a conta não existe

#### Cenário: REQ-15.4 — Cadastro não escolhe o papel
- DADO um cadastro que informa o papel administrador entre os dados
- QUANDO a conta é criada
- ENTÃO a conta tem o papel usuário

### Requisito: REQ-16 — Consulta da própria conta
O sistema DEVE mostrar ao usuário logado o id, o e-mail, o username, o nome de exibição, o papel, se o e-mail está confirmado e o e-mail de uma troca pendente, e NÃO DEVE expor o e-mail de uma conta a outro usuário.

#### Cenário: REQ-16.1 — Dados da própria conta
- DADO um usuário logado com o e-mail `ana@exemplo.com` e uma troca pendente para `ana@novo.com`
- QUANDO ele consulta a própria conta
- ENTÃO vê `ana@exemplo.com`, o username, o nome de exibição, o papel, o estado do e-mail e `ana@novo.com` como troca pendente

#### Cenário: REQ-16.3 — Só a própria conta
- DADO as contas `ana@exemplo.com` e `beto@exemplo.com`
- QUANDO a conta de Ana consulta a própria conta
- ENTÃO a resposta traz `ana@exemplo.com`
- E nenhuma resposta do contexto traz o e-mail de outra conta

#### Cenário: REQ-16.4 — Conta excluída com token de acesso ainda válido
- DADO um usuário que excluiu a conta e ainda tem um token de acesso no prazo
- QUANDO ele consulta a conta
- ENTÃO é recusado como não autenticado

#### Cenário: REQ-16.2 — Sem sessão
- DADO uma requisição sem sessão válida
- QUANDO ela consulta a conta
- ENTÃO é recusada como não autenticada

### Requisito: REQ-17 — Senha e tokens nunca em texto
O sistema NÃO DEVE guardar a senha nem os tokens dos links de e-mail e de renovação da sessão em texto, só um resumo criptográfico que não permita recuperá-los.

#### Cenário: REQ-17.1 — Senha guardada só como hash
- DADO uma conta criada com a senha `segredo123`
- QUANDO os dados guardados da conta são lidos
- ENTÃO nenhum deles contém `segredo123`
- E o login com `segredo123` é aceito

#### Cenário: REQ-17.2 — Token de link guardado só como hash
- DADO um link de confirmação emitido
- QUANDO os dados guardados dos tokens são lidos
- ENTÃO nenhum deles contém o token do link

#### Cenário: REQ-17.3 — Credencial de renovação guardada só como hash
- DADO uma sessão aberta
- QUANDO os dados guardados das sessões são lidos
- ENTÃO nenhum deles contém a credencial de renovação
