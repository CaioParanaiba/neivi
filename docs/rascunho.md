# Rascunho

## Fatia atual · 4 Listagem e resumo diario
- Tarefa nova = criada hoje ou ontem (D-11)
- ✓ Resumo diario: atrasadas -> urgentes -> importantes (D-27) -> RN-16, US-07
- ✓ Ordenacoes do resumo diario e do padrao (D-04, D-28) -> RN-17
- ✓ Filtro por disciplina usa a disciplina efetiva (RN-14) -> RN-18
- ✓ Ler so as tarefas pendentes no banco, sem carregar o historico -> RNF-01 ampliado
- Ver as tarefas feitas no dia separadas das pendentes (motivacao), sem poluir (US-04)
- ✓ Visao do dia (C1), CLI (-> US-07): listagem rapida compacta com tarefas do dia por ordem de prioridade (com contexto de disciplina ou projeto)
- ✓ Organizar aulas (C7), CLI (-> US-09, filtro `new`): listagem completa de tarefas vinculadas as aulas apenas
- ✓ Revisao do dia (C8), CLI (-> US-09 e US-08): listagem completa de tarefas, com um campo indicando as novas com as antigas

## Proximas fatias (ordem da D-34)

### Fatia 5 · Importacao e backup
- Chaves do Ganbatte vao para a descricao na importacao (D-31)
- Q-17 (quais arquivos) e Q-18 (done.txt)

- Depois da importacao: corte das tarefas (tarefas so no neivi; prompt das 5:50 neutro)

### Fatia 2 · Disciplinas e projetos
- RN-14 (disciplina herdada do projeto) no `ls @`/`@BD` e no `show` (CA-10.5): falta o projeto ter disciplina
- Campos da D-06, panorama do projeto (C9), Q-15, Q-16
- Controle de versoes do esquema + migracao 002 (o 001 nao muda mais depois do corte)

### Fatia 3 · Habitos
- Habito: ao bater a meta sai da visao do dia; excedente fica no historico (D-12)
- Habito semanal: aparece no dia ate bater a meta da semana; dias especificos: so aparece nos seus dias
- Habito pode ser registrado em hoje ou ontem (D-21)

## Ideias futuras

- Editar varias tarefas de uma vez (ex.: adiar tudo de BD) (US-03)
- Concluir com datas mais antigas que ontem, se for frequente (US-04)
- Revisao semanal (C10): `neivi week` com as feitas na semana + o planejado nas proximas 2 semanas -> decidir depois do corte das tarefas (US-09)
- Ritual mensal de backup (copiar para o HD + testar restauracao) eh um caso real de tarefa recorrente (D-16)

### TUI (futuro)
- Visao do dia (C1): listagem rapida compacta com tarefas do dia por ordem de prioridade (com contexto de disciplina ou projeto)
- Organizar aulas (C7): listagem completa de tarefas vinculadas as aulas apenas
- Revisao do dia (C8): listagem completa de tarefas, com de um lado as novas do outro as antigas, podendo definir novas prioridades

## Duvidas novas (antes de virarem Q)
- (vazio)

---

## Fatia 1 · Tarefas (fechada em 2026-11-11)
- ✓ Criterios de cada nivel de prioridade A a D (Q-14) -> RN-15
- ✓ Tarefa incompleta = sem prioridade definida (D-10) -> RN-02
- ✓ Backup automatico -> RNF-02, RNF-03 (script com envio ao GitHub feito)
- ✓ Mensagens de erro dizem o que esta errado e como resolver -> RNF-04
- ✓ Confirmacao curta (`+` `~` `✓` `↺` `✗`) e contador do dia no `done` -> RNF-05
- ✓ Filtro por disciplina inclui as tarefas dos projetos (D-15) -> refinado na RN-14
- ~~Tarefa cancelada (sai das pendentes, fica no historico)~~ descartado: nao quero guardar o que nao fiz (US-06)
- ~~Nota sem avaliacao: cria a avaliacao na hora (D-14)~~ v0.2 (D-23)
