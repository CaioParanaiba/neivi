# Rascunho

## Fatia atual · ajustes da semana de uso (antes da tag v0.1.0)
Triagem de 18/11 (notas da semana + `neivi-uso`). Codigo primeiro, CAs depois por perguntas guiadas.
- Atalho do grupo: `project <codigo|+codigo|id>` = show; `habit <codigo|id>` = do, avisa `ja feito hoje (para marcar de novo: neivi habit do X)` e nao marca; `habits` = `habit`; grupo sozinho (`project`, `habit`, `subject`) = ls -> US-13, US-17, RN nova
- RN nova: codigo de projeto ou habito nao pode ser nome de subcomando (add, ls, show, edit, rm, do, done, undo); erro de codigo inexistente lista os subcomandos
- `habit done` = sinonimo de `habit do` (errar o verbo)
- `;` `&` `|` `#` `(` na descricao sao do zsh, nao do neivi: dica no help (aspas ou trocar por `,` / ` - `)
- Ficha de projeto e disciplina ordenada por prioridade -> prazo -> id (o trio da maratona fica junto); today e ls sem mudanca
- Habito "N por semana" sai do today no dia em que foi feito e volta no dia seguinte; semana incompleta fecha incompleta -> RN-23
- today: bloco `amanha` entre urgentes e importantes (sex, sab e dom = `ate segunda`), fora do `falta:`, some vazio -> US-07
- today: `sem prazo` so com tarefas sem projeto + linha `+ N de projetos (p1: 4, ...) · neivi project show <codigo>`; A e tarefa com prazo aparecem mesmo de projeto -> RN-16
- [ideia] 18/11 ter uma marcacao para deixar claro que nao farei um hbaito hoje, pois as vezes ja passou ou nao vou fazer, assim mostrando que ele nao sera feito consigo focar nos habitos que ainda estao, como sera a visualizacao disso a decidir -> candidata da migracao 002 (registro novo no banco: coluna em habit_logs ou tabela habit_skips), junto da 191
- [chato] 18/11 percebo que sempre realizo uma acao no neivi eu sempre logo em seguida eu faco neivi today -> descartada em 19/11: as vezes quero so dar done; a TUI resolve (visao sempre na tela)
- [ideia] 19/11 mostrar no today o numero de tarefas incompletas agaurdadno o triage -> feito em 20/11: linha `incompletas: N · neivi triage` antes do `falta:` (US-07, CA-07.13, CA-07.14; CA-07.7 trocado pela linha)
- [duvida] 20/11 o bloco `sem prazo` (B sem prazo avulsas) impede o `nada para hoje`; so atrasadas, urgentes e importantes deveriam contar? -> RN-16 / CA-07.6, decidir depois da v0.1

## Proximas fatias
- Temas dentro do projeto (o trio estudar/implementar/resolver da maratona); migracao 002; perguntas guiadas proprias (checklist fora: outro objetivo)
- Ficha da disciplina: `subject show` + atalho `subject BD` / `subject 2` (projetos ligados, avulsas, ultima atividade, ordem da ficha), junto com as notas da v0.2 (US-15); ate la, `neivi ls @BD` (tentado em 18/11)

## Ideias futuras

- Editar varias tarefas de uma vez (ex.: adiar tudo de BD) (US-03)
- Concluir com datas mais antigas que ontem, se for frequente (US-04)
- Resumo da carga dos proximos dias (`amanha: 3 tarefas (~2h) · qui: 1`), junto com o `neivi week`
- Descricao pelo editor ou stdin, se o `;` do zsh voltar a incomodar (caso raro em 18/11)
- Revisao semanal (C10): `neivi week` com as feitas na semana + o planejado nas proximas 2 semanas -> decidir depois do corte (US-09)
- Revisao semanal da maratona (#106 do Ganbatte) volta quando houver recorrencia
- Ritual mensal de backup (copiar para o HD + testar restauracao) eh um caso real de tarefa recorrente (D-16)

- v0.2 (notas): ao registrar a nota final, oferecer encerrar a disciplina (US-15)

- Compensar o dia de um habito de dias especificos (fiz na segunda, entao a terca nao precisa) (US-18)
- Revisoes semanais (financeiro na sexta, estudos no domingo) como tarefas recorrentes, quando recorrencia entrar (D-16)

- Autocompletar no zsh (US-23, depois da semana de uso): nivel 1 = comandos, palavras do ls e valores das marcas (pri:, due:, sit:); nivel 2 = @disciplinas, +projetos, codigos dos habitos que faltam e ids das pendentes com o titulo, via um comando escondido `neivi __complete ...`

### TUI (futuro)
- Visao do dia (C1): listagem rapida compacta com tarefas do dia por ordem de prioridade (com contexto de disciplina ou projeto)
- Organizar aulas (C7): listagem completa de tarefas vinculadas as aulas apenas
- Revisao do dia (C8): listagem completa de tarefas, com de um lado as novas do outro as antigas, podendo definir novas prioridades

## Duvidas novas (antes de virarem Q)
- (vazio)

---

## Semana de uso real (fechada em 2026-11-18)
Uma semana de uso real (ate ~18/11). Anotar na hora, sem resolver, com `neivi-note <marca> <texto>`:
`[bug]` diferente do requisito · `[chato]` funciona, mas incomoda · `[falta]` precisei e nao tem · `[ideia]` sem precisar agora.
No fim da semana: `neivi-uso` (resumo do log de comandos) e a triagem desta lista.
- ✓ [bug] 11/11 nao sei se eh um bug ou um erro nos requisitos, mas o neivi project show 4 nao funcionou pois ele so funciona o show com o codigo e nao com o id
- ✓ [bug] 11/11 durante a descricao o comando
- ✓ [bug] 11/11 durante a descricao o comando ; acaba nao funcionando encerrando o codigo o mesmo pra r o neivi-note como visto em cima
- ✓ [chato] 12/11 tem algum motivo para o habit usar do e a task done, parece que tenho qeu decorar quando uso cada um
- ✓ [ideia] 12/11 na vizualizacao de tarefas no mrt, como sempre tenho um trio de conteudo, estudar implementar e resolver, devido a est, as tarefas nao ficam sempre juntas, exemplo 31 32 33, ai tenho que ficar proxurando elas, sera que para show de projetos e desciplinas nao fazer sorte pelo est fazer pelo id em si? ou pensamos em outra abordagaem
- ✓ [ideia] 16/11 nao vejo necessidade de mostrar habitos semanais no dia, apos fazer o habito do dia, exemplo: lixo ao fazer uma vez no dia ele tem que sumir do neivi today
- ✓ [chato] 16/11 eh mais importante ver tarefas com proximos prazos, estou dessa forma so vend as proximas tarefas no dia seguinte, as vezes vale mais apena adiantar algo dos proximos dias do que focar nas sem prazo
- Triagem: tudo para a fatia atual, menos temas (proximas fatias) e editor/resumo da carga (ideias futuras)

## Fatia 5 · Importacao (fechada em 2026-11-11)
- ✓ Chaves do Ganbatte vao para a descricao na importacao (D-31) -> UC-01
- ✓ Q-17 e Q-18 -> D-38 (done.txt inteiro entra)
- ✓ Ganbatte desligado (ganbatte.service) em 2026-11-11

## Fatia 3 · Habitos (fechada em 2026-11-11)
- ✓ Habito: ao bater a meta sai da visao do dia; excedente fica no historico (D-12) -> RN-23, RN-24
- ✓ Habito semanal e dias especificos -> RN-23
- ✓ Habito pode ser registrado em hoje ou ontem (D-21) -> RN-24

## Fatia 2 · Disciplinas e projetos (fechada em 2026-11-11)
- ✓ RN-14 (disciplina herdada do projeto) no `ls @`/`@BD` e no `show` (CA-10.5) -> codigo desta fatia
- ✓ Campos da D-06, panorama do projeto (C9) -> US-11 a US-14; Q-15 -> RN-19; Q-16 -> RN-21
- ✓ Esquema mudou direto no 001 (ainda sem dados reais, D-35)

## Fatia 4 · Listagem e resumo diario (fechada em 2026-11-11)
- Tarefa nova = criada hoje ou ontem (D-11)
- ✓ Resumo diario: atrasadas -> urgentes -> importantes (D-27) -> RN-16, US-07
- ✓ Ordenacoes do resumo diario e do padrao (D-04, D-28) -> RN-17
- ✓ Filtro por disciplina usa a disciplina efetiva (RN-14) -> RN-18
- ✓ Ler so as tarefas pendentes no banco, sem carregar o historico -> RNF-01 ampliado
- Ver as tarefas feitas no dia separadas das pendentes (motivacao), sem poluir (US-04)
- ✓ Visao do dia (C1), CLI (-> US-07): listagem rapida compacta com tarefas do dia por ordem de prioridade (com contexto de disciplina ou projeto)
- ✓ Organizar aulas (C7), CLI (-> US-09, filtro `new`): listagem completa de tarefas vinculadas as aulas apenas
- ✓ Revisao do dia (C8), CLI (-> US-09 e US-08): listagem completa de tarefas, com um campo indicando as novas com as antigas

## Fatia 1 · Tarefas (fechada em 2026-11-11)
- ✓ Criterios de cada nivel de prioridade A a D (Q-14) -> RN-15
- ✓ Tarefa incompleta = sem prioridade definida (D-10) -> RN-02
- ✓ Backup automatico -> RNF-02, RNF-03 (script com envio ao GitHub feito)
- ✓ Mensagens de erro dizem o que esta errado e como resolver -> RNF-04
- ✓ Confirmacao curta (`+` `~` `✓` `↺` `✗`) e contador do dia no `done` -> RNF-05
- ✓ Filtro por disciplina inclui as tarefas dos projetos (D-15) -> refinado na RN-14
- ~~Tarefa cancelada (sai das pendentes, fica no historico)~~ descartado: nao quero guardar o que nao fiz (US-06)
- ~~Nota sem avaliacao: cria a avaliacao na hora (D-14)~~ v0.2 (D-23)
