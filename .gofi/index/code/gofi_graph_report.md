# Mapa do codigo — github.com/Arthur-155/art-books

Gerado por gofi index (go, modo `fast`).
Fonte de verdade: `gofi_graph.json`. Este relatorio e o resumo navegavel dele.

> Leia este arquivo **antes** de varrer o repositorio. Ele diz onde olhar.
> Para detalhe de qualquer simbolo use `gofi show <simbolo>` em vez de abrir arquivos.

## Resumo

8 pacotes, 45 tipos, 44 funcoes, 90 metodos em 15 arquivos (2166 linhas).
Grafo com 224 nos e 597 arestas, das quais 178 sao chamadas.
Modo fast: chamadas por heuristica sintatica. 11 ambiguas e 3 nao resolvidas — rode `--deep` para precisao total.
## Pacotes

Instabilidade I = saida/(entrada+saida). Perto de 0 = base estavel, muita gente depende dele. Perto de 1 = camada de borda.

| pacote | simbolos | linhas | dependem dele | ele depende de | I |
|---|---:|---:|---:|---:|---:|
| `domain/contas/entity` | 17 | 139 | 4 | 5 | 0.56 |
| `domain/contas/repository` | 34 | 410 | 2 | 7 | 0.78 |
| `domain/contas/service` | 53 | 753 | 2 | 15 | 0.88 |
| `domain/contas/adapter` | 3 | 42 | 1 | 4 | 0.80 |
| `domain/contas/handler` | 44 | 406 | 1 | 10 | 0.91 |
| `domain/contas/mailer` | 8 | 99 | 1 | 8 | 0.89 |
| `domain/contas/worker` | 9 | 77 | 1 | 6 | 0.86 |
| `art-books` | 11 | 240 | 0 | 24 | 1.00 |

## Pontos centrais

Simbolos com maior acoplamento. Mexer aqui tem alcance largo.

- `domain/contas/entity.User` — struct, `domain/contas/entity/user.go:56`, entrada 12 / saida 1
- `domain/contas/entity.LimitKey` — struct, `domain/contas/entity/user.go:51`, entrada 11 / saida 1
  LimitKey identifica o e-mail ou o IP a que um bloqueio ou limite se aplica.
- `domain/contas/handler.decodeBody` — func, `domain/contas/handler/middleware.go:103`, entrada 10 / saida 0
- `domain/contas/handler.finish` — func, `domain/contas/handler/middleware.go:112`, entrada 9 / saida 0
  finish responde só com o status quando não há erro.
- `domain/contas/service.hashToken` — func, `domain/contas/service/auth_service.go:257`, entrada 8 / saida 0
  hashToken resume com SHA-256 os tokens de link e de renovação.
- `domain/contas/service.contasService.issueAndSend` — method, `domain/contas/service/contas_service.go:424`, entrada 4 / saida 5
  issueAndSend grava o hash do token, que invalida o link anterior, e envia o link; a falha do envio mantém o token.
- `domain/contas/entity.EmailToken` — struct, `domain/contas/entity/user.go:81`, entrada 6 / saida 1
- `domain/contas/service.contasService.findUser` — method, `domain/contas/service/contas_service.go:372`, entrada 5 / saida 2
- `domain/contas/entity.RequestType` — type, `domain/contas/entity/user.go:35`, entrada 6 / saida 0
- `domain/contas/entity.TokenPurpose` — type, `domain/contas/entity/user.go:19`, entrada 6 / saida 0
- `domain/contas/handler.clearSessionCookies` — func, `domain/contas/handler/middleware.go:80`, entrada 5 / saida 1
- `domain/contas/handler.currentUserID` — func, `domain/contas/handler/middleware.go:120`, entrada 5 / saida 1
- `domain/contas/entity.NormalizeEmail` — func, `domain/contas/entity/user.go:109`, entrada 5 / saida 0
- `domain/contas/entity.Role` — type, `domain/contas/entity/user.go:12`, entrada 5 / saida 0
- `art-books.run` — func, `art-books/main.go:31`, entrada 1 / saida 7
- `domain/contas/service.NewContasService` — func, `domain/contas/service/contas_service.go:124`, entrada 2 / saida 5
- `domain/contas/repository.contasRepository.consumeRequired` — method, `domain/contas/repository/contas_repository.go:224`, entrada 3 / saida 3
- `domain/contas/service.contasService.acceptEmailRequest` — method, `domain/contas/service/contas_service.go:448`, entrada 3 / saida 3
  acceptEmailRequest aplica o limite de 1 por minuto e 5 por hora, por e-mail e por IP, e registra o pedido de forma atômica.
- `domain/contas/service.contasService.Login` — method, `domain/contas/service/auth_service.go:33`, entrada 1 / saida 6
  Login abre a sessão; e-mail inexistente e senha errada respondem igual.
- `domain/contas/service.contasService.RequestEmailChange` — method, `domain/contas/service/contas_service.go:225`, entrada 1 / saida 6
  RequestEmailChange envia a aprovação ao e-mail atual; o novo só recebe algo depois da aprovação.

## Comunidades

Grupos de simbolos que se chamam e se referenciam entre si. Nao seguem a arvore de pastas: e o acoplamento real.

**C0 · service** — 156 simbolos em `domain/contas/service`, `domain/contas/handler`, `domain/contas/repository`, `domain/contas/entity`, `art-books`, `domain/contas/mailer`
  nucleo: `domain/contas/entity.User`, `domain/contas/entity.LimitKey`, `domain/contas/handler.decodeBody`, `domain/contas/handler.finish`, `domain/contas/service.hashToken`

**C1 · worker** — 3 simbolos em `domain/contas/worker`
  nucleo: `domain/contas/worker.PurgeWorker.Run`, `domain/contas/worker.PurgeWorker.Stop`, `domain/contas/worker.PurgeWorker.runOnce`

## Como consultar sem abrir arquivos

```sh
gofi show <simbolo>             # tudo sobre um simbolo: origem, vizinhos, doc
gofi find --in code <palavras>  # busca por palavras; substitui o grep -r
gofi path <A> <B>               # como A alcanca B, aresta por aresta
gofi index open                 # abre a visualizacao HTML do grafo
```

Os nomes aceitam forma curta: `NewServer`, `api.NewServer` ou o ID completo funcionam.
