
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

- ✓ Criterios de cada nivel de prioridade A a D (Q-14) -> RN-15
- ✓ Tarefa incompleta = sem prioridade definida (D-10) -> RN-02
- Tarefa nova = criada hoje ou ontem (D-11)
- Resumo diario: atrasadas -> urgentes (prazo hoje, qualquer prioridade) -> importantes (A sem prazo); atrasada = prazo antes de hoje (D-27, glossario)
- Habito: ao bater a meta sai da visao do dia; excedente fica no historico (D-12)
- Habito semanal: aparece no dia ate bater a meta da semana; dias especificos: so aparece nos seus dias
- Ordenacoes do resumo diario e do padrao (D-04)
- Filtro por disciplina inclui tarefas de projetos ligados a ela (D-15)
- ~~Nota sem avaliacao: cria a avaliacao na hora (D-14)~~ v0.2 (D-23)
- Habito pode ser registrado em hoje ou ontem (D-21)

## Itens para levar a outras secoes

- ✓ Backup automatico -> secao 5 (RNF de confiabilidade) -> RNF-02, RNF-03
- ✓ RNF de usabilidade: toda mensagem de erro diz o que esta errado e como resolver (US-02, CA-02.6) -> RNF-04
- ✓ RNF de usabilidade: depois de criar ou editar, mostrar a tarefa numa linha curta (`+` ao criar, `~` ao editar, `✓` ao concluir, `↺` ao reabrir, `✗` ao remover: `+ id · pri · dia dd/mm[/aa] · est · titulo @disc +proj …`), sem campos vazios e com o ano so quando nao eh o atual (US-01 a US-06); o `done` mostra no fim o contador do dia (`hoje: N concluidas`, ou `ontem:`) -> RNF-05
- Editar varias tarefas de uma vez (ex.: adiar tudo de BD) -> proximas, so se fizer falta (US-03)
- Concluir com datas mais antigas que ontem -> proximas, so se for frequente (US-04)
- Ver as tarefas feitas no dia separadas das pendentes (motivacao), sem poluir -> fatia 4 (US-04)
- ~~Tarefa cancelada (sai das pendentes, fica no historico)~~ descartado: nao quero guardar o que nao fiz (US-06)
- Chaves do Ganbatte vao para a descricao na importacao (D-31) -> fatia de importacao
- Ler so as tarefas pendentes no banco, sem carregar o historico (RNF-01 a longo prazo) -> fatia 4
- Ritual mensal de backup (copiar para o HD + testar restauracao) eh um caso real de tarefa recorrente (D-16) -> quando recorrencia entrar
- Script de backup: enviar o dump todo dia para `~/.todo/neivi/` + commit e push (RNF-03, D-33); conferir chave SSH/agente no timer -> fim da fatia 1
