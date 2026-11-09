
## Ideias de visualização

### Visão do dia (C1)
- CLI: listagem rapida compacta com tarefas do dia por ordem de prioridade (com contexto de disciplina ou projeto)
- TUI: listagem rapida compacta com tarefas do dia por ordem de prioridade (com contexto de disciplina ou projeto)

### Organizar aulas (C7)
- CLI: listagem completa de tarefas vinculadas as aulas apenas
- TUI: listagem completa de tarefas vinculadas as aulas apenas

### Revisao do dia (C8)
- CLI: listagem completa de tarefas, com um campo indicando as novas com as antigas
- TUI: listagem completa de tarefas, com de um lado as novas do outro as antigas, podendo definir novas prioridades

## Secao 8

Movida para o requisitos.md em 2026-11-09. Aqui ficam so as duvidas novas, antes de virarem Q.

## Candidatas a regras de negocio (secao 4)

Para escrever junto com as stories de cada fatia:

- Criterios de cada nivel de prioridade A a D (Q-14, em aberto)
- Tarefa incompleta = sem prioridade definida (D-10)
- Tarefa nova = criada hoje ou ontem (D-11)
- Resumo diario: atrasadas -> urgentes (prazo hoje, qualquer prioridade) -> importantes (A sem prazo); atrasada = prazo antes de hoje (D-27, glossario)
- Habito: ao bater a meta sai da visao do dia; excedente fica no historico (D-12)
- Habito semanal: aparece no dia ate bater a meta da semana; dias especificos: so aparece nos seus dias
- Ordenacoes do resumo diario e do padrao (D-04)
- Filtro por disciplina inclui tarefas de projetos ligados a ela (D-15)
- ~~Nota sem avaliacao: cria a avaliacao na hora (D-14)~~ v0.2 (D-23)
- Habito pode ser registrado em hoje ou ontem (D-21)

## Itens para levar a outras secoes

- Backup automatico -> secao 5 (RNF de confiabilidade)
- Chaves do Ganbatte vao para a descricao na importacao (D-31) -> fatia de importacao
