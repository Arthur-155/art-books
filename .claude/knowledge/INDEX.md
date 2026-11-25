# Índice do conhecimento portável — carregue por necessidade

> **Não carregue estas pastas por glob.** Ler tudo em toda invocação é o maior
> custo fixo do harness — maior que a spec do contexto em que se está mexendo.
>
> **Núcleo ⬤** é curto e universal: carregue sempre. **O resto é sob demanda**:
> leia a linha *Quando* e carregue só o que a tarefa pede.
>
> Derivado — regenere com `gofi index docs` após `gofi update sdk`.

## `.claude/sdk/go/knowledge/`

_go — padrões e armadilhas do SDK_

| | Arquivo | Linhas | Quando |
|---|---|--:|---|
| ⬤ | `absolute-rules.md` | 32 | Regras invioláveis. Cada uma é Go-specific e nasceu de erro repetido em |
|  | `api-endpoint-rules.md` | 48 | Índice; detalhe no arquivo apontado. API: .claude/sdk/go/api/netx-httpx.md, |
|  | `application-bridge-rules.md` | 95 | Índice; detalhe no arquivo apontado. Cliente HTTP do SDK: |
|  | `bridge-factory-adapter-pattern.md` | 663 | Use quando o contexto integra com N implementações intercambiáveis de |
|  | `bucket-storage.md` | 235 | Padrão para subir/baixar arquivos em object storage (S3/MinIO/R2, OCI Object |
|  | `cache-layer.md` | 88 | Cache é detalhe de acesso a dados: vive só no repository. O service não |
|  | `cloud-identity.md` | 79 | Referência gerada: .claude/sdk/go/api/base-cloud-aws.md, base-cloud-oci.md. |
|  | `configuration.md` | 175 | API: .claude/sdk/go/api/base-environment.md, gofi-config.md, |
|  | `consumer-bootstrap.md` | 205 | mq := messaging.New() // MESSAGING_ + import _ ".../msq/provider/<nome>" |
|  | `cronjob.md` | 96 | API: .claude/sdk/go/api/base-cronjob.md. Como o job vive dentro do serviço |
|  | `database-connection.md` | 151 | API: .claude/sdk/go/api/gofi-component-database.md, sqln-connection.md, |
|  | `dynamic-filter.md` | 216 | Quando o cliente escolhe os filtros (em vez de query params fixos). API: |
|  | `env-vars-standard.md` | 94 | (struct environment.Environment — .claude/sdk/go/api/base-environment.md). |
|  | `error-handling.md` | 140 | Import: github.com/joaoprofile/gofi-sdk-go/base/errs. Referência completa: |
|  | `events-scheduling-rules.md` | 102 | Índice; detalhe no arquivo apontado. Mensageria do SDK (producer, consumer, |
|  | `full-text-search.md` | 191 | API: .claude/sdk/go/api/sqln-fts.md, sqln-criteria.md, sqln-connection.md, |
|  | `gofi-orchestrator.md` | 164 | API: .claude/sdk/go/api/gofi.md e gofi-component-.md. Programas completos |
|  | `grpc.md` | 106 | API: .claude/sdk/go/api/netx-grpcx.md e gofi-component-grpcserver.md. Exemplo executável com |
|  | `handler-test-pattern.md` | 86 | Todo handler tem {contexto}_handler_test.go cobrindo o óbvio de cada endpoint. |
|  | `http-auth-middleware.md` | 238 | Setup do serviço iam, login, refresh e logout: iam.md. Gate de permissão: |
|  | `http-client-signing.md` | 84 | Referência: .claude/sdk/go/api/netx-httpx.md (§HttpClient, Request, Signature), |
|  | `http-server.md` | 95 | Referência: .claude/sdk/go/api/netx-httpx.md, gofi-component-httpserver.md. |
|  | `iam-adapter-pattern.md` | 126 | Caminho padrão: componente iamc.New com UserPort/TenantPort do projeto |
|  | `iam-social-login.md` | 109 | Pré-requisito: iam.md (ports, sessão, entrega token/cookie). Referência: |
|  | `iam.md` | 206 | Referência: .claude/sdk/go/api/iam.md e gofi-component-iam.md. Programa |
|  | `kafka-consumer-naming.md` | 131 | Convenção para nomes de consumer groups quando o projeto usa |
|  | `layer-adapter.md` | 128 | Leia este primeiro; o documento completo só na situação da última seção. |
|  | `layer-handler.md` | 140 | Leia este primeiro; o documento completo só na situação da última seção. |
|  | `layer-model.md` | 110 | Leia este primeiro; o documento completo só na situação da última seção. |
|  | `layer-repository.md` | 132 | Leia este primeiro; o documento completo só na situação da última seção. |
|  | `layer-service.md` | 142 | Leia este primeiro; o documento completo só na situação da última seção. |
| ⬤ | `layers.md` | 67 | Nenhuma camada lê ambiente, abre conexão ou encerra o processo: recebe tudo |
|  | `loader-pattern.md` | 399 | Pattern Go para carregar snapshot consistente de múltiplas tabelas em |
|  | `logging.md` | 109 | API: .claude/sdk/go/api/obs-logging.md. Import: |
|  | `lookup-endpoints.md` | 127 | Quando um campo do {Ctx}FilterMapping (ver dynamic-filter.md) é um |
|  | `mail.md` | 101 | Referência gerada: .claude/sdk/go/api/base-mail.md; env → config em |
|  | `messaging-msq.md` | 237 | Referência gerada: .claude/sdk/go/api/msq.md (fachada), msq-types.md |
|  | `migrations.md` | 75 | API: .claude/sdk/go/api/sqln-migrate.md e connection.WithMigrations em |
| ⬤ | `naming.md` | 57 | nunca currency, currencyId, coin ou moeda. Carrega o código ISO 4217 |
|  | `observability-otel.md` | 424 | API: .claude/sdk/go/api/gofi-component-observability.md, obs.md, |
|  | `pagination.md` | 76 | API: .claude/sdk/go/api/sqln-pagination.md e sqln.md. Exemplo que percorre |
|  | `persistence-rules.md` | 83 | Índice das regras; o detalhe mora no arquivo apontado. API do SDK em |
|  | `postgres-index-strategy.md` | 198 | Aplicado quando o contexto materializa entidades em PostgreSQL via |
|  | `qa-checklist-postgres.md` | 136 | Itens que se somam a qa-checklist.md quando o contexto toca banco |
|  | `qa-checklist.md` | 361 | Aplicado por gofi-qa em todo contexto Go implementado. Itens de PostgreSQL |
|  | `rbac.md` | 121 | Referência: .claude/sdk/go/api/iam-provider-rbac-roles.md, iam-port.md. |
|  | `read-endpoints.md` | 193 | Padrão para qualquer rota GET de listagem/detalhe servida ao front. Mantém o |
|  | `report-export.md` | 83 | Export de listagem = o mesmo filtro dinâmico do contexto |
|  | `repository-aggregate-pattern.md` | 213 | Mecânica de transação (propagação por ctx, isolamento, retry, aninhamento): |
|  | `repository-close.md` | 58 | O repository não é dono de recurso de banco: não guarda sql.Stmt nem |
|  | `repository-primitive-return.md` | 100 | Método que consulta um único valor primitivo (existência, flag, escalar |
|  | `repository-update-simple.md` | 71 | estrutural do banco. |
|  | `service-bootstrap.md` | 134 | Mecânica do orquestrador: gofi-orchestrator.md. Ambiente e segredos: |
|  | `session-store.md` | 74 | Referência: .claude/sdk/go/api/base-session.md, gofi-component-session.md. |
|  | `shared-packages-rules.md` | 77 | Índice; detalhe no arquivo apontado. helpers e money são pacotes do |
| ⬤ | `structure.md` | 117 | Layout obrigatório para um contexto de domínio em projetos gofi/Go. |
|  | `transactions.md` | 73 | API: .claude/sdk/go/api/sqln-transaction.md e sqln-connection.md |
|  | `validation.md` | 159 | Import: github.com/joaoprofile/gofi-sdk-go/base/validator (tags do |
|  | `value-objects.md` | 109 | API: .claude/sdk/go/api/sqln-mapping.md. Todo sqln.Find escaneia em T |
|  | `worker-bootstrap.md` | 275 | broker), roda Run no ListenAndServe e chama Stop no shutdown — |

## `.claude/sdk/go/api/`

_go — API do SDK por pacote (gerada do código)_

| | Arquivo | Linhas | Quando |
|---|---|--:|---|
|  | `base-bucket-buckettest.md` | 29 | Package buckettest is the contract every bucket.Store must satisfy. Run it |
|  | `base-bucket-file.md` | 163 | Package file implements bucket.Store on a local directory, for development |
|  | `base-bucket-mem.md` | 127 | Package mem implements bucket.Store in memory, for tests and local runs. |
|  | `base-bucket-oci.md` | 102 | Package oci implements bucket.Store on top of OCI Object Storage. |
|  | `base-bucket-s3.md` | 168 | Package s3 implements bucket.Store for Amazon S3 and S3-compatible services |
|  | `base-bucket.md` | 604 | Package bucket defines a provider-agnostic object-storage abstraction. |
|  | `base-cloud-aws.md` | 71 | Package aws loads AWS credentials for every gofi integration (S3, SQS, RDS, |
|  | `base-cloud-oci.md` | 100 | Package oci resolves OCI credentials for every gofi integration (object |
|  | `base-common.md` | 114 | const ( |
|  | `base-cronjob.md` | 279 | const ( |
|  | `base-debug.md` | 117 | const DefaultAddr = "localhost:6060" |
|  | `base-environment.md` | 685 | const APP_MAX_PARALLEL_WORKERS = 1 |
|  | `base-errs.md` | 269 | func JoinErrors(errs error) error |
|  | `base-mail.md` | 321 | Package mail sends transactional and bulk e-mail over any SMTP provider, with |
|  | `base-observer.md` | 102 | var ( |
|  | `base-redact.md` | 90 | Package redact renders configuration structs without their secrets, for fmt, |
|  | `base-secrets-awssm.md` | 53 | Package awssm resolves secret://awssm/<name-or-arn>#key from AWS Secrets |
|  | `base-secrets-ocivault.md` | 72 | Package ocivault resolves secret://ocivault/<secret-ocid>#key and |
|  | `base-secrets.md` | 143 | Package secrets resolves secret references of the form |
|  | `base-session.md` | 369 | const ( |
|  | `base-timezone.md` | 40 | const BrazilName = "America/Sao_Paulo" |
|  | `base-validator.md` | 64 | func IsStruct(s any) error |
|  | `examples.md` | 115 | A people.v1.PeopleService that moves a Person entity and its photo through every gRPC call shape… |
|  | `gofi-component-cache.md` | 94 | Package cache is the gofi component for the shared Redis client used by the |
|  | `gofi-component-database.md` | 128 | Package database is the gofi component for the SQL database (sqln). |
|  | `gofi-component-grpcserver.md` | 129 | Package grpcserver is the gofi component for the gRPC server (netx/grpcx). |
|  | `gofi-component-httpserver.md` | 143 | Package httpserver is the gofi component for the HTTP server (netx/httpx). It is a |
|  | `gofi-component-iam.md` | 97 | Package iam is the gofi component for the identity service (iam): JWT |
|  | `gofi-component-messaging.md` | 113 | Package messaging is the gofi component for the message broker (msq). |
|  | `gofi-component-observability.md` | 98 | Package observability is the gofi component that exports traces, metrics |
|  | `gofi-component-session.md` | 66 | Package session is the gofi component that installs the global session |
|  | `gofi-config-core.md` | 163 | Package core holds the process-wide settings gofi's Build applies before any |
|  | `gofi-config.md` | 151 | Package config is gofi's composition layer between the environment loader |
|  | `gofi.md` | 239 | type Builder interface { |
|  | `iam-config.md` | 170 | const DefaultSessionMaxLifetime = 30  24  time.Hour |
|  | `iam-core.md` | 330 | const MaxClockSkew = 5  time.Minute |
|  | `iam-middleware.md` | 144 | func AuthInterceptor(svc core.IAMService, opts ...GRPCOption) grpc.UnaryServerInterceptor |
|  | `iam-port.md` | 390 | type AuthInput struct { |
|  | `iam-provider-bcrypt.md` | 50 | Package bcrypt provides password hashing utilities for use in UserPort implementations. |
|  | `iam-provider-google.md` | 73 | Package google implements IDPAuthPort for Google OAuth2/OIDC. |
|  | `iam-provider-jwks.md` | 187 | Package jwks verifies JWTs signed by a third party that publishes its public |
|  | `iam-provider-jwt.md` | 103 | var ErrInvalidKey = errors.New("iam/jwt: invalid key") |
|  | `iam-provider-memory.md` | 188 | Package memory implements port.SessionPort in memory. |
|  | `iam-provider-microsoft.md` | 88 | Package microsoft implements IDPAuthPort for Microsoft Identity Platform (Entra ID). |
|  | `iam-provider-oidc.md` | 122 | Package oidc implements a generic IDPAuthPort for any OIDC-compliant provider. |
|  | `iam-provider-password.md` | 98 | Package password hashes passwords with Argon2id (RFC 9106, OWASP's first |
|  | `iam-provider-rbac-roles.md` | 63 | Package roles implements port.RBACPort with a simple role-to-permissions model. |
|  | `iam-provider-redis.md` | 220 | Package redis implements port.SessionPort using Redis. |
|  | `iam-types.md` | 227 | type Claims struct { |
|  | `iam.md` | 65 | Package iam is an identity and access facade for the gofi SDK. |
|  | `msq-core.md` | 277 | Package core implements the orchestration layer of the msq messaging system. |
|  | `msq-msqtest.md` | 51 | Package msqtest is the contract every msq provider must satisfy. Provider |
|  | `msq-port.md` | 119 | type Broker interface { |
|  | `msq-provider-kafka.md` | 138 | Package kafka implements port.Broker for Apache Kafka using the Sarama library. |
|  | `msq-provider-nats.md` | 123 | Package nats implements port.Broker for NATS JetStream: persistent |
|  | `msq-provider-oci.md` | 52 | type Broker struct { |
|  | `msq-provider-rabbitmq.md` | 158 | Package rabbitmq implements port.Broker for RabbitMQ using AMQP 0-9-1. |
|  | `msq-provider-redis.md` | 135 | type Broker struct { |
|  | `msq-provider-sqs.md` | 58 | Package sqs implements the msq broker on Amazon SQS (aws-sdk-go-v2). |
|  | `msq-types.md` | 421 | const ( |
|  | `msq-worker.md` | 208 | const ( |
|  | `msq.md` | 563 | const ( |
|  | `netx-awssign.md` | 68 | Package awssign signs httpx client requests with AWS Signature V4 (API Gateway IAM |
|  | `netx-grpcx.md` | 403 | Package grpcx serves and calls gRPC services with the same defaults as the |
|  | `netx-httpx.md` | 1241 | Package httpx serves and calls HTTP APIs: a chi-based server with hardened |
|  | `netx.md` | 94 | Package netx holds what the HTTP and gRPC transports share: the hardened |
|  | `obs-logging.md` | 263 | const ( |
|  | `obs-metrics.md` | 80 | Package metrics is the OpenTelemetry metric API used by gofi and by |
|  | `obs.md` | 178 | var ErrInvalidTLS = errors.New("obs: invalid OTLP TLS configuration") |
|  | `sqln-cache.md` | 155 | func Close() error |
|  | `sqln-connection.md` | 729 | const ( |
|  | `sqln-criteria.md` | 601 | const ( |
|  | `sqln-driver-mysql.md` | 220 | var ErrConcurrentIndex = errors.New("mysql: a FULLTEXT index cannot be added without blocking wr… |
|  | `sqln-driver-oracle.md` | 100 | type Driver struct{} |
|  | `sqln-driver-postgres.md` | 218 | Package postgres registers the PostgreSQL driver, built on pgx/v5. |
|  | `sqln-driver-sqlserver.md` | 111 | type Driver struct{} |
|  | `sqln-driver.md` | 452 | const ( |
|  | `sqln-filter.md` | 323 | const ( |
|  | `sqln-fts.md` | 582 | Package fts defines full-text indexes and the searches that read them. |
|  | `sqln-mapping.md` | 59 | func EachT any(rows sql.Rows) iter.Seq2T, error |
|  | `sqln-migrate.md` | 56 | func RegisterDriver(d Driver) |
|  | `sqln-pagination.md` | 134 | const ( |
|  | `sqln-query.md` | 32 | type Query interface { |
|  | `sqln-rdsauth.md` | 29 | Package rdsauth signs RDS / Aurora IAM authentication tokens, so databases |
|  | `sqln-statement.md` | 40 | type Statement interface { |
|  | `sqln-transaction.md` | 48 | type Options struct { |
|  | `sqln.md` | 577 | const ( |

## `.claude/sdk/go/boilerplates/`

_go — esqueletos por camada_

| | Arquivo | Linhas | Quando |
|---|---|--:|---|
|  | `handler-test.md` | 187 | Regras: .claude/sdk/go/knowledge/handler-test-pattern.md. Casa com |
|  | `handler.md` | 222 | Servidor, ordem Use → UseAuth → Handlers e o que já vem ligado: |
|  | `main.md` | 159 | Regras e ordem: .claude/sdk/go/knowledge/service-bootstrap.md. Orquestrador: |
|  | `model.md` | 137 | //gofi:context {contexto} |
|  | `repository.md` | 188 | Regras: .claude/sdk/go/knowledge/persistence-rules.md. API: |
|  | `service.md` | 283 | Regras de erro: .claude/sdk/go/knowledge/error-handling.md. Validação: |

## Custo

| | Linhas | ~Tokens |
|---|--:|--:|
| Tudo (glob) | 27322 | ~243k |
| Núcleo ⬤ | 273 | ~4k |
| Núcleo + 2 módulos | ~873 | ~10k |
