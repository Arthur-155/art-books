# art-books

Serviço Go criado com `gofi init` e especificado por contexto (PRD → spec → mudança → entrega).

## Estrutura

```
.
├── go.work
├── backend/
│   ├── go.mod
│   ├── .migrations/        # SQL em pares up/down
│   ├── art-books/          # composition root (main, config, wire)
│   └── domain/<contexto>/  # entity, repository, service, handler, mailer, worker
├── prd/<contexto>/         # PRD de cada contexto
├── specs/<contexto>/       # spec viva, mudanças e design
└── .claude/                # skills, SDK e conhecimento do gofi
```

## Contextos

`contas`, `catalogo`, `leitura-e-resumos`, `recomendacoes`, `feeds`, `interacoes`, `perfil` e `amizades`.
O estado de cada um: `gofi status`.

## Build e execução

```sh
go -C backend run ./art-books
```

A configuração é lida do ambiente em `backend/art-books/config.go`.

## Testes

```sh
gofi test
```

## Fluxo de trabalho

```
/gofi-pd      # descoberta: PRD a partir do problema
/gofi-spec    # contrato: proposta, deltas e design da mudança
gofi run <mudança>   # entrega: o loop implementa, testa, audita e arquiva
```

Mudanças só rodam depois de aprovadas pela pessoa (`gofi change approve`).
