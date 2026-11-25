# Mapa do codigo — github.com/joaoprofile/gofi-sdk-go/base

Gerado por gofi index (go, modo `fast`).
Fonte de verdade: `gofi_graph.json`. Este relatorio e o resumo navegavel dele.

> Leia este arquivo **antes** de varrer o repositorio. Ele diz onde olhar.
> Para detalhe de qualquer simbolo use `gofi show <simbolo>` em vez de abrir arquivos.

## Resumo

117 pacotes, 595 tipos, 999 funcoes, 1177 metodos em 322 arquivos (40653 linhas).
Grafo com 3039 nos e 8846 arestas, das quais 2400 sao chamadas.
Modo fast: chamadas por heuristica sintatica. 1138 ambiguas e 246 nao resolvidas — rode `--deep` para precisao total.
## Pacotes

Instabilidade I = saida/(entrada+saida). Perto de 0 = base estavel, muita gente depende dele. Perto de 1 = camada de borda.

| pacote | simbolos | linhas | dependem dele | ele depende de | I |
|---|---:|---:|---:|---:|---:|
| `github.com/joaoprofile/gofi-sdk-go/obs/logging` | 43 | 441 | 30 | 11 | 0.27 |
| `github.com/joaoprofile/gofi-sdk-go/gofi` | 29 | 455 | 22 | 15 | 0.41 |
| `github.com/joaoprofile/gofi-sdk-go/iam/port` | 21 | 334 | 13 | 3 | 0.19 |
| `github.com/joaoprofile/gofi-sdk-go/msq` | 54 | 490 | 13 | 13 | 0.50 |
| `github.com/joaoprofile/gofi-sdk-go/iam/types` | 9 | 193 | 12 | 1 | 0.08 |
| `environment` | 66 | 716 | 11 | 15 | 0.58 |
| `github.com/joaoprofile/gofi-sdk-go/netx/httpx` | 223 | 3668 | 11 | 47 | 0.81 |
| `github.com/joaoprofile/gofi-sdk-go/sqln/connection` | 106 | 1530 | 11 | 24 | 0.69 |
| `github.com/joaoprofile/gofi-sdk-go/gofi/config/core` | 16 | 268 | 10 | 13 | 0.57 |
| `github.com/joaoprofile/gofi-sdk-go/msq/port` | 7 | 110 | 10 | 2 | 0.17 |
| `github.com/joaoprofile/gofi-sdk-go/msq/types` | 30 | 433 | 10 | 9 | 0.47 |
| `github.com/joaoprofile/gofi-sdk-go/sqln/driver` | 24 | 434 | 10 | 3 | 0.23 |
| `github.com/joaoprofile/gofi-sdk-go/iam/core` | 74 | 1242 | 9 | 18 | 0.67 |
| `github.com/joaoprofile/gofi-sdk-go/msq/worker` | 23 | 258 | 7 | 6 | 0.46 |
| `github.com/joaoprofile/gofi-sdk-go/sqln/fts` | 64 | 1355 | 7 | 18 | 0.72 |
| `bucket` | 48 | 665 | 6 | 13 | 0.68 |
| `github.com/joaoprofile/gofi-sdk-go/gofi/component/messaging` | 13 | 237 | 6 | 14 | 0.70 |
| `github.com/joaoprofile/gofi-sdk-go/iam/config` | 5 | 175 | 6 | 5 | 0.45 |
| `errs` | 39 | 240 | 5 | 5 | 0.50 |
| `github.com/joaoprofile/gofi-sdk-go/gofi/component/httpserver` | 20 | 286 | 5 | 11 | 0.69 |
| `github.com/joaoprofile/gofi-sdk-go/netx` | 14 | 247 | 5 | 10 | 0.67 |
| `github.com/joaoprofile/gofi-sdk-go/netx/grpcx` | 76 | 1340 | 5 | 33 | 0.87 |
| `cloud/aws` | 10 | 103 | 4 | 7 | 0.64 |
| `redact` | 12 | 226 | 4 | 7 | 0.64 |
| `github.com/joaoprofile/gofi-sdk-go/examples/grpc/people` | 11 | 175 | 4 | 11 | 0.73 |
| `github.com/joaoprofile/gofi-sdk-go/examples/obs/store` | 8 | 85 | 4 | 4 | 0.50 |
| `github.com/joaoprofile/gofi-sdk-go/examples/obs/telemetry` | 5 | 84 | 4 | 8 | 0.67 |
| `github.com/joaoprofile/gofi-sdk-go/obs/metrics` | 9 | 118 | 4 | 5 | 0.56 |
| `github.com/joaoprofile/gofi-sdk-go/sqln` | 98 | 992 | 4 | 33 | 0.89 |
| `github.com/joaoprofile/gofi-sdk-go/sqln/criteria` | 85 | 1055 | 4 | 9 | 0.69 |
| `github.com/joaoprofile/gofi-sdk-go/sqln/migrate` | 10 | 168 | 4 | 11 | 0.73 |
| `cloud/oci` | 12 | 133 | 3 | 6 | 0.67 |
| `common` | 14 | 253 | 3 | 13 | 0.81 |
| `observer` | 15 | 144 | 3 | 7 | 0.70 |
| `secrets` | 15 | 191 | 3 | 7 | 0.70 |
| `github.com/joaoprofile/gofi-sdk-go/examples/grpc/grpcapi` | 12 | 211 | 3 | 10 | 0.77 |
| `github.com/joaoprofile/gofi-sdk-go/gofi/component/database` | 18 | 309 | 3 | 15 | 0.83 |
| `github.com/joaoprofile/gofi-sdk-go/iam/provider/oidc` | 22 | 484 | 3 | 21 | 0.88 |
| `github.com/joaoprofile/gofi-sdk-go/iam/provider/password` | 9 | 141 | 3 | 8 | 0.73 |
| `github.com/joaoprofile/gofi-sdk-go/sqln/driver/postgres` | 38 | 490 | 3 | 14 | 0.82 |

_(+77 pacotes omitidos)_

## Pontos centrais

Simbolos com maior acoplamento. Mexer aqui tem alcance largo.

- `session.Session.Unlock` — method, `base/session/session.go:182`, entrada 104 / saida 2
  Unlock releases the lock for key only while token (from TryLock) still owns
- `environment.Environment` — struct, `base/environment/environment.go:61`, entrada 57 / saida 0
  Environment holds all configuration loaded from environment variables.
- `github.com/joaoprofile/gofi-sdk-go/sqln/criteria.Predicate` — struct, `sqln/criteria/predicates.go:17`, entrada 43 / saida 0
  Predicate represents a single WHERE/HAVING condition or a logical connector (AND/OR).
- `github.com/joaoprofile/gofi-sdk-go/msq/provider/kafka.scramClient.Done` — method, `msq/provider/kafka/scram.go:39`, entrada 37 / saida 0
- `github.com/joaoprofile/gofi-sdk-go/sqln/driver.TextIndex` — struct, `sqln/driver/fulltext.go:62`, entrada 34 / saida 2
  TextIndex is the validated definition of a full-text index. The same
- `github.com/joaoprofile/gofi-sdk-go/obs/logging.Error` — func, `obs/logging/logging.go:90`, entrada 32 / saida 1
- `github.com/joaoprofile/gofi-sdk-go/sqln/filter.Filters.Add` — method, `sqln/filter/dynamic_filter.go:154`, entrada 29 / saida 2
- `github.com/joaoprofile/gofi-sdk-go/obs/logging.Warn` — func, `obs/logging/logging.go:92`, entrada 26 / saida 1
- `github.com/joaoprofile/gofi-sdk-go/netx/httpx.Response` — func, `netx/httpx/response_parse.go:129`, entrada 22 / saida 4
- `github.com/joaoprofile/gofi-sdk-go/iam/internal/flight.Group.Do` — method, `iam/internal/flight/flight.go:31`, entrada 22 / saida 2
  Do runs fn unless a call is already in flight, then waits for it.
- `bucket.Object` — struct, `base/bucket/bucket.go:77`, entrada 23 / saida 0
  Object is the provider-agnostic metadata of a stored object.
- `errs.AppError` — struct, `base/errs/errs.go:32`, entrada 20 / saida 1
- `github.com/joaoprofile/gofi-sdk-go/sqln/driver.Dialect` — interface, `sqln/driver/dialect.go:9`, entrada 19 / saida 1
- `github.com/joaoprofile/gofi-sdk-go/obs/logging.FromContext` — func, `obs/logging/logging.go:99`, entrada 18 / saida 2
- `github.com/joaoprofile/gofi-sdk-go/netx/httpx.Error` — func, `netx/httpx/response_parse.go:145`, entrada 18 / saida 1
  Error writes err as an ErrorResponse. A *RequestError (ParseRequestBody,
- `github.com/joaoprofile/gofi-sdk-go/msq/port.MessageHandler` — interface, `msq/port/handler.go:15`, entrada 18 / saida 0
  MessageHandler is the interface for processing broker messages.
- `github.com/joaoprofile/gofi-sdk-go/obs/logging.Info` — func, `obs/logging/logging.go:89`, entrada 17 / saida 1
- `github.com/joaoprofile/gofi-sdk-go/sqln/fts.Index` — struct, `sqln/fts/index.go:114`, entrada 17 / saida 1
  Index is a validated full-text index definition. Declare each one once
- `github.com/joaoprofile/gofi-sdk-go/sqln/filter.Filter` — struct, `sqln/filter/dynamic_filter.go:114`, entrada 17 / saida 0
  Domain Types
- `github.com/joaoprofile/gofi-sdk-go/gofi.gofiInstance.With` — method, `gofi/builder.go:44`, entrada 16 / saida 0

## Comunidades

Grupos de simbolos que se chamam e se referenciam entre si. Nao seguem a arvore de pastas: e o acoplamento real.

**C0 · httpx** — 2080 simbolos em `github.com/joaoprofile/gofi-sdk-go/netx/httpx`, `github.com/joaoprofile/gofi-sdk-go/sqln/connection`, `github.com/joaoprofile/gofi-sdk-go/sqln`, `github.com/joaoprofile/gofi-sdk-go/sqln/criteria`, `github.com/joaoprofile/gofi-sdk-go/iam/core`, `github.com/joaoprofile/gofi-sdk-go/netx/grpcx`
  nucleo: `session.Session.Unlock`, `environment.Environment`, `github.com/joaoprofile/gofi-sdk-go/sqln/criteria.Predicate`, `github.com/joaoprofile/gofi-sdk-go/msq/provider/kafka.scramClient.Done`, `github.com/joaoprofile/gofi-sdk-go/sqln/driver.TextIndex`

**C1 · environment** — 51 simbolos em `environment`, `redact`, `bucket`, `debug`, `mail`
  nucleo: `redact.GoSprint`, `redact.Format`, `redact.JSON`, `redact.LogValue`, `redact.Sprint`

**C2 · redis** — 14 simbolos em `github.com/joaoprofile/gofi-sdk-go/iam/provider/redis`
  nucleo: `github.com/joaoprofile/gofi-sdk-go/iam/provider/redis.Provider.sessionKey`, `github.com/joaoprofile/gofi-sdk-go/iam/provider/redis.Provider.userKey`, `github.com/joaoprofile/gofi-sdk-go/iam/provider/redis.markRevoked`, `github.com/joaoprofile/gofi-sdk-go/iam/provider/redis.Provider.prune`, `github.com/joaoprofile/gofi-sdk-go/iam/provider/redis.revokedTTL`

**C3 · gofi** — 9 simbolos em `github.com/joaoprofile/gofi-sdk-go/gofi`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/cache`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/database`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/grpcserver`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/httpserver`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/iam`
  nucleo: `github.com/joaoprofile/gofi-sdk-go/gofi.Stage`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/cache.Component.Stage`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/database.Component.Stage`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/grpcserver.Component.Stage`, `github.com/joaoprofile/gofi-sdk-go/gofi/component/httpserver.Component.Stage`

**C4 · msq** — 7 simbolos em `github.com/joaoprofile/gofi-sdk-go/msq`
  nucleo: `github.com/joaoprofile/gofi-sdk-go/msq.redact`, `github.com/joaoprofile/gofi-sdk-go/msq.OCICredentials.redacted`, `github.com/joaoprofile/gofi-sdk-go/msq.ProviderConfig.redacted`, `github.com/joaoprofile/gofi-sdk-go/msq.RedisConnection.redacted`, `github.com/joaoprofile/gofi-sdk-go/msq.plainOCICredentials`

**C5 · httpx** — 6 simbolos em `github.com/joaoprofile/gofi-sdk-go/netx/httpx`
  nucleo: `github.com/joaoprofile/gofi-sdk-go/netx/httpx.NewSafeClient`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx.destinationGuard`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx.safeRedirect`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx.isBlocked`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx.sameOriginRedirect`

**C6 · common** — 4 simbolos em `common`, `validator`
  nucleo: `validator.IsStruct`, `common.CamelToSnake`, `common.ParseStructName`, `common.ParseStructColumns`

**C7 · observer** — 4 simbolos em `observer`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx`
  nucleo: `observer.Registry`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx.HandlerRegistry`, `github.com/joaoprofile/gofi-sdk-go/netx/httpx.NewHandlerRegistry`, `observer.New`

**C8 · jwt** — 4 simbolos em `github.com/joaoprofile/gofi-sdk-go/iam/provider/jwt`, `github.com/joaoprofile/gofi-sdk-go/iam/port`
  nucleo: `github.com/joaoprofile/gofi-sdk-go/iam/provider/jwt.buildRefreshToken`, `github.com/joaoprofile/gofi-sdk-go/iam/provider/jwt.buildRefreshTokenInternal`, `github.com/joaoprofile/gofi-sdk-go/iam/provider/jwt.Provider.IssueRefreshToken`, `github.com/joaoprofile/gofi-sdk-go/iam/port.TokenPort`

**C9 · connection** — 4 simbolos em `github.com/joaoprofile/gofi-sdk-go/sqln/connection`
  nucleo: `github.com/joaoprofile/gofi-sdk-go/sqln/connection.hasControl`, `github.com/joaoprofile/gofi-sdk-go/sqln/connection.hasSpaceOrControl`, `github.com/joaoprofile/gofi-sdk-go/sqln/connection.validateHost`, `github.com/joaoprofile/gofi-sdk-go/sqln/connection.Settings.Validate`

**C10 · aws** — 3 simbolos em `cloud/aws`
  nucleo: `cloud/aws.mask`, `cloud/aws.Config.redacted`, `cloud/aws.plain`

**C11 · oci** — 3 simbolos em `cloud/oci`
  nucleo: `cloud/oci.mask`, `cloud/oci.Config.redacted`, `cloud/oci.plain`

_(+28 comunidades menores omitidas)_
## Como consultar sem abrir arquivos

```sh
gofi show <simbolo>             # tudo sobre um simbolo: origem, vizinhos, doc
gofi find --in code <palavras>  # busca por palavras; substitui o grep -r
gofi path <A> <B>               # como A alcanca B, aresta por aresta
gofi index open                 # abre a visualizacao HTML do grafo
```

Os nomes aceitam forma curta: `NewServer`, `api.NewServer` ou o ID completo funcionam.
