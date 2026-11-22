# Mapa do codigo — https://github.com/Arthur-155/art-books.git

Gerado por gofi index (go, modo `fast`).
Fonte de verdade: `gofi_graph.json`. Este relatorio e o resumo navegavel dele.

> Leia este arquivo **antes** de varrer o repositorio. Ele diz onde olhar.
> Para detalhe de qualquer simbolo use `gofi show <simbolo>` em vez de abrir arquivos.

## Resumo

1 pacotes, 0 tipos, 1 funcoes, 0 metodos em 1 arquivos (8 linhas).
Grafo com 3 nos e 2 arestas, das quais 0 sao chamadas.
Modo fast: chamadas por heuristica sintatica. 0 ambiguas e 0 nao resolvidas — rode `--deep` para precisao total.
## Pacotes

Instabilidade I = saida/(entrada+saida). Perto de 0 = base estavel, muita gente depende dele. Perto de 1 = camada de borda.

| pacote | simbolos | linhas | dependem dele | ele depende de | I |
|---|---:|---:|---:|---:|---:|
| `art-books` | 1 | 8 | 0 | 1 | 1.00 |

## Pontos centrais

Simbolos com maior acoplamento. Mexer aqui tem alcance largo.


## Como consultar sem abrir arquivos

```sh
gofi show <simbolo>             # tudo sobre um simbolo: origem, vizinhos, doc
gofi find --in code <palavras>  # busca por palavras; substitui o grep -r
gofi path <A> <B>               # como A alcanca B, aresta por aresta
gofi index open                 # abre a visualizacao HTML do grafo
```

Os nomes aceitam forma curta: `NewServer`, `api.NewServer` ou o ID completo funcionam.
