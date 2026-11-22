# Índice do conhecimento portável — carregue por necessidade

> **Não carregue estas pastas por glob.** Ler tudo em toda invocação é o maior
> custo fixo do harness — maior que a spec do contexto em que se está mexendo.
>
> **Núcleo ⬤** é curto e universal: carregue sempre. **O resto é sob demanda**:
> leia a linha *Quando* e carregue só o que a tarefa pede.
>
> Derivado — regenere com `gofi index docs` após `gofi update sdk`.

## Custo

| | Linhas | ~Tokens |
|---|--:|--:|
| Tudo (glob) | 0 | ~0k |
| Núcleo ⬤ | 0 | ~0k |
| Núcleo + 2 módulos | ~600 | ~5k |
