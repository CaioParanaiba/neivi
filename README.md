# neivi

Tarefas, hábitos e rotina no terminal. Escrito em Go.

> **Projeto pessoal.** O neivi é feito para a minha rotina e para o meu jeito de
> trabalhar. Não é pensado para outras pessoas usarem: não há instalador, suporte
> nem configurações para outros casos de uso. O código está aberto para quem
> quiser ler, mas pedidos de funcionalidade e pull requests não são o foco.

## O que faz

- Tarefas na sintaxe do [todo.txt](https://github.com/todotxt/todo.txt), com
  prioridade, prazo, tags, contextos e estimativa de tempo.
- Tarefas filhas e checklist dentro da tarefa.
- Hábitos com meta, frequência e dias de folga, com sequência (streak) calculada
  a partir do histórico.
- Visão do dia: o que vence hoje, o que está atrasado e os hábitos do dia.
- Exportação para todo.txt.

Os dados ficam em SQLite, num único arquivo local.

## Roadmap
(ideia)

| Versão | Conteúdo |
|---|---|
| v0.1 | CLI: tarefas, filhas, checklist, hábitos, visão do dia, exportação para todo.txt |
| v0.2 | TUI (Bubble Tea + Lip Gloss), cronômetro por tarefa, backup integrado |
| v0.3 | Disciplinas e avaliações |
| v0.4 | Projetos e fluxo de aula |
| v0.5 | Integrações (Codeforces, GitHub, Google Calendar, Gmail) e resumo do dia |
| v1.0 | Assistente, estatísticas, recorrências, servidor em casa e acesso pelo celular |

Detalhes em [docs/SPEC.md](docs/SPEC.md).

## Nome

Neivi é o nome da minha cachorra :).
