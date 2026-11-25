# Guia do gofi: como funciona e onde paramos

## 1. A ideia em uma frase

Você descreve **o quê** quer (PRD e spec), aprova, e o `gofi run` põe agentes de IA para
implementar, testar, auditar e arquivar sozinhos. Você decide; o gofi executa.

## 2. O pipeline, na ordem

| Passo | Comando | O que faz | Quem decide |
|---|---|---|---|
| 1. Discovery | `/gofi-pd` | Gera o PRD (problema, escopo) em `prd/<contexto>/` | você revisa |
| 2. Contrato | `/gofi-spec` | Escreve uma **mudança**: proposta, deltas e design, em `specs/<contexto>/` | você revisa |
| 3. Revisão | `gofi review <mudança>` | Procura ambiguidade, contradição e cenário não testável | você corrige |
| 4. Aprovação | `gofi change approve <mudança>` | Sela o contrato. Só a pessoa aprova | **você** |
| 5. Entrega | `gofi run <mudança>` | O loop implementa, testa, audita e arquiva | o loop |
| 6. Estado | `gofi status` | Onde cada contexto está e o próximo passo | informativo |

Termos:
- **Contexto:** uma área do produto (`contas`, `catalogo`, `recomendacoes`, `feeds`, `interacoes`,
  `leitura-e-resumos`, `perfil`, `amizades`).
- **Spec viva** (`specs/<contexto>/spec.md`): a verdade de um contexto. Não se edita à mão.
  Ela só aparece quando uma mudança é aceita e arquivada.
- **Mudança:** todo trabalho é uma mudança, criada por `gofi change new`, nunca à mão.
- **Gênese:** a primeira mudança de um contexto novo.

## 3. O que o `gofi run` faz por dentro

Cada execução (`r-AAAAMMDD-HHMMSS`) percorre um grafo de etapas:

```
plan → implement → test → sense → (repair ↺) → audit → policy → arquivar
```

- **implement:** um agente escreve o código seguindo o design e o SDK Go.
- **test:** outro agente escreve os testes.
- **sense:** sensores automáticos (build, `go vet`, testes, rastreio, segurança).
- **repair:** corrige o que o sensor reprovou. Não implementa o que falta.
- **audit:** um QA revisa e grava achados com severidade (`gofi verdict`).
- **gate (portão):** se falta uma decisão que a mudança não cobre, o agente grava uma
  pergunta e para. Quem responde é você (`gofi answer <run>`).

Limites: o loop tem `loop.max_iterations` voltas de correção (hoje 3). Estourou, termina
`exhausted`: não é erro técnico, é "gastei as voltas e ainda há achados".

Regras que o gofi impõe:
- Só roda mudança **aprovada**.
- A árvore git precisa estar **limpa** (ou use `--worktree`).
- Execução não aceita não deixa commit: o trabalho vai para a branch `gofi/<mudança>`
  e a árvore volta ao início. Cada nova tentativa cria outra branch (`-2`, `-3`...).
- Execução aceita vira um commit no branch atual.

## 4. As 5 decisões que são só suas

1. **Revisar** os achados de `gofi review`.
2. **Aprovar** o contrato (`gofi change approve`).
3. **Responder portões** (`gofi answer <run>`).
4. **Decidir uma emenda** quando termina em `impasse`.
5. **Aceitar ou recusar lições** (`gofi learn accept` / `reject`).

## 5. O que fizemos hoje (24/11/2026)

1. **Erro de versão.** Binário `1917106` (build local) e skills `30af369` não combinavam.
   Compilei o commit `30af369` do repo público do gofi e instalei em `~/.local/bin/gofi`.
   O binário antigo está em `~/.local/bin/gofi.bak-1917106`.
2. **Remoto git quebrado.** A URL do `origin` tinha o comando `git remote add origin ...`
   dentro dela, então nada subia. Corrigido; `master` está no GitHub.
3. **Árvore suja.** O loop recusou rodar com arquivos soltos; fizemos commit das specs e memória.
4. **7 execuções falharam** (`exhausted`) porque o **SDK Go não estava instalado**: o agente
   parou numa pergunta, o fluxo seguiu sem implementação, os testes quebraram o build.
5. **Instalamos o SDK** (`gofi update sdk`), o índice (`gofi index`) e a expertise
   (`gofi update expertise`).
6. **`gofi run contas` funcionou:** escreveu o código inteiro (entity, repository, service,
   handler, mailer, worker, composition root). Parou em `exhausted` porque a auditoria achou
   26 problemas e só havia 3 voltas.
7. **Merge** da branch `gofi/contas-4` no `master`. Build, vet e testes do service passam.
8. **Limpeza:** branches `gofi/*` apagadas, menos `gofi/contas` no GitHub (veja abaixo).
9. **README** atualizado (`src/` virou `backend/`).

Custo total das execuções: cerca de US$ 7.

## 6. Onde estamos

| Contexto | PRD | Mudança | Situação |
|---|---|---|---|
| contas | ✓ | aprovada | código no `master`, **26 achados de auditoria abertos** |
| recomendacoes | ✓ | aprovada | depende de contas, catalogo e leitura-e-resumos |
| catalogo | ✓ | aberta, **não aprovada** | falta `gofi review` e `gofi change approve` |
| leitura-e-resumos, feeds, interacoes, perfil, amizades | ✓ | sem spec | falta `/gofi-spec <contexto>` |

Achados de `contas`: 1 maior (falta transação em `ResetPassword`, `ChangePassword`,
`ConfirmEmailChange`) e cerca de 25 menores. Detalhe em
`.gofi/runs/r-20261124-232625/findings.ndjson`.

## 7. Roteiro para amanhã

1. **GitHub:** Settings → Branches → Default branch → `master`.
   Depois `git push origin --delete gofi/contas`.
2. **Decidir duas perguntas** que o agente de reparo deixou:
   - E-mail de reset de senha: enviar **fora da resposta** (recomendado) ou aceitar o
     risco de enumeração por tempo?
   - `HTTP_TRUSTED_PROXIES` vazio: **falhar no boot** (recomendado) ou só documentar?
3. **Subir as voltas:** em `.gofi.yaml`, `loop.max_iterations` de 3 para 6.
4. **Rodar de novo:** `gofi run contas` (árvore limpa). Se parar num portão,
   `gofi answer <run>`. Rodar sempre em terminal sem limite de tempo; um corte no meio
   deixa a execução `interrompida` e se retoma com `gofi run --resume <run>`.
5. Quando `contas` for aceita: aprovar `catalogo` (`gofi review catalogo`,
   `gofi change approve catalogo`) e entregar.
6. Escrever a spec de `leitura-e-resumos` (`/gofi-spec leitura-e-resumos`), aprovar e entregar.
7. Só então rodar `recomendacoes`, que depende das três.

## 8. Comandos de diagnóstico

| Comando | Para quê |
|---|---|
| `gofi status` | estado e próximo passo de cada contexto |
| `gofi ps` | lista as execuções e se acabaram |
| `gofi watch` | acompanha execuções ao vivo |
| `gofi logs <run>` | cada nó, a aresta tomada e por quê |
| `gofi report <run>` | o que a execução entregou, custo e achados |
| `gofi spec check` | valida as specs |
| `gofi update audit` | o que ficou desatualizado no projeto |
| `gofi run --resume <run>` | retoma uma execução interrompida |
| `gofi cancel <run>` | para uma execução em curso |

Arquivos de uma execução: `.gofi/runs/<run>/` (`events.ndjson`, `report.md`,
`findings.ndjson`, `agents/`).

## 9. Armadilhas que já pegamos

- **Versão do binário diferente da das skills:** o gofi recusa. `GOFI_ALLOW_SKEW=1` contorna,
  mas o certo é alinhar os dois (já feito).
- **SDK ausente:** nunca rode `gofi run` antes de `gofi update sdk`. Cada tentativa custa e
  cria branch inútil.
- **Árvore suja:** faça commit antes do `run`.
- **Responder portão:** `gofi answer <run>`. Digitar um comando na pergunta não vale; só o número
  da opção ou o texto dela.
- **Busca no projeto:** a guarda do gofi bloqueia `grep`, `find` e `Read` varrendo a árvore.
  Use `gofi find "<pergunta>"`, `gofi find --text "<literal>"` e `gofi show <símbolo>`.
- **Branch padrão do GitHub:** foi parar em `gofi/contas` porque foi a primeira enviada.
