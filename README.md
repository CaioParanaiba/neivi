# neivi

Tarefas, hábitos e rotina no terminal. Escrito em Go.

> **Projeto pessoal.** O neivi é feito para a minha rotina e para o meu jeito de
> trabalhar. Não é pensado para outras pessoas usarem: não há instalador, suporte
> nem configurações para outros casos de uso. O código está aberto para quem
> quiser ler, mas pedidos de funcionalidade e pull requests não são o foco.

## O que faz

- Tarefas com uma sintaxe própria, inspirada no
  [todo.txt](https://github.com/todotxt/todo.txt): prioridade, prazo,
  estimativa de tempo, descrição, disciplina e projeto em qualquer posição da
  linha.
- Disciplinas e projetos cadastrados, com as tarefas ligadas a eles.
- Hábitos com meta e periodicidade (todo dia, dias específicos ou N vezes por
  semana).
- Resumo do dia: o que está atrasado, o que vence hoje e o que é importante.

Os dados ficam em SQLite, num único arquivo local.

## Roadmap
(ideia)

| Versão | Conteúdo |
|---|---|
| v0.1 | CLI: tarefas, disciplinas, projetos, hábitos, resumo do dia, backup |
| v0.2 | Avaliações e notas; TUI (Bubble Tea + Lip Gloss) |
| depois | Sequência dos hábitos, cronômetro, fluxo de aula, integrações (Codeforces, GitHub, Google Calendar, Gmail), servidor em casa |

Detalhes em [docs/requisitos.md](docs/requisitos.md) (o que faz) e
[docs/SPEC.md](docs/SPEC.md) (como).

## Nome

Neivi é o nome da minha cachorra :).
