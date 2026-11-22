---
tipo: prd
formato: 3
contexto: leitura-e-resumos
versao: "1.0"
status: discovery
spec: n/a
operacoes: []
marketplaces: []
keywords: [leitura, estante, quero ler, lendo, lido, capítulo, resumo de capítulo, resumo do livro, visibilidade, público, privado]
atualizado: 2026-11-22
---
# PRD — leitura-e-resumos — Registrar o que se entendeu de cada leitura

> Terceiro contexto. Depende de **contas** e **catalogo**. É a base de **interacoes**, **feeds**, **perfil** e **recomendacoes**.

## 1. Problema

Quem lê não tem onde registrar, com as próprias palavras, o que entendeu de cada capítulo. Sem esse registro, o que se leu se perde e não há o que compartilhar com outros leitores.

## 2. Objetivo e métricas de sucesso

Permitir que o leitor acompanhe seus livros e escreva um resumo por capítulo e, ao fim, um resumo do livro inteiro.

| Métrica | Hoje | Meta (hipótese) | Como medir |
|---|---|---|---|
| Livros "lendo" com ao menos um capítulo resumido | n/a (produto novo) | 60% | Capítulos resumidos por livro |
| Livros "lidos" com resumo do livro | n/a | 50% | Resumos sobre livros lidos |
| Resumos publicados (públicos) sobre o total | n/a | 40% | Visibilidade dos livros |

## 3. Personas

| Persona | Papel | O que precisa | O que a impede hoje |
|---|---|---|---|
| Leitora | Registra e resume | Um lugar para anotar capítulo a capítulo | Anotações soltas, sem estrutura |
| Leitora que compartilha | Publica resumos | Controlar o que é público | Sem rede própria para isso |

## 4. Escopo

### Dentro
- Estante com três estados por livro: "quero ler", "lendo" e "lido"; vários "lendo" ao mesmo tempo.
- Capítulos adicionados livremente, com número e título opcional.
- Resumo de capítulo com até 4.000 caracteres, editável a qualquer momento.
- Resumo do livro, opcional; é ele que aparece no feed.
- Visibilidade pública ou privada, uma por livro, valendo para o resumo e para os capítulos.
- Um registro por usuário e por livro, editável.

### Fora
- Meta de leitura e progresso por página: não faz parte do primeiro corte.
- Notas privadas fora do resumo e importação de capítulos: ficam para depois.
- Segunda leitura do mesmo livro: um registro por livro basta por enquanto.

## 5. Jornadas

### J1 — Resumir um livro em leitura
1. A leitora adiciona o livro à estante como "lendo".
2. Adiciona os capítulos e escreve o que cada um diz.
3. Edita os resumos quando quiser.

### J2 — Concluir e publicar
1. A leitora marca o livro como "lido" e escreve o resumo do livro.
2. Escolhe se o livro é público ou privado.
3. Se for público e tiver resumo do livro, ele passa a aparecer nos feeds.

## 6. Restrições, premissas e riscos

| Tipo | Descrição | Impacto | Mitigação |
|---|---|---|---|
| Restrição | Publicar exige e-mail confirmado | Evita contas descartáveis publicando | Regra de contas |
| Restrição | Só entra no feed o livro público com resumo do livro | Livro sem resumo do livro não circula | Regra explícita |
| Premissa | A visibilidade é única por livro | Simplifica a privacidade | Perfil segue a mesma regra |
| Premissa | O "mais recente" do perfil é o de atividade mais recente | Define o "lendo agora" | Regra explícita |
| Risco | Publicar sem querer | Exposição de conteúdo | Sugestão de visibilidade inicial privada |

## 7. Dúvidas em aberto

Nenhuma. Decisões fechadas em 2026-11-22:

- Resumo do livro: até 4.000 caracteres, igual ao do capítulo.
- Visibilidade inicial de um livro novo: privado.
- Limite de 200 capítulos por livro.
- Troca de e-mail em andamento: os livros públicos não mudam, porque o e-mail antigo segue confirmado até a troca ser concluída.
- Metas da seção 2: mantidas como hipóteses, a revisar depois do lançamento.
