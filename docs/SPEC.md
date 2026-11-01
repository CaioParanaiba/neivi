# neivi: especificação

Documento vivo. Versão de 31/10/2026, cobrindo o v0.1 em detalhe e as fases
seguintes em linhas gerais.

## 1. Princípios

1. **Projeto unicamente pessoal.** Cada decisão é tomada para um único usuário.
   A pergunta "e se outra pessoa quiser usar?" não entra na conta: nada de
   instalador, temporadas, personalização para casos que não são os meus,
   suporte ou layout pensado para terceiros. Se uma ideia só faz sentido para
   outras pessoas, ela fica fora.
2. **Código comum primeiro, IA por cima.** Toda função precisa funcionar sem IA.
   A IA (via `claude -p` ou similar) entra só onde poupa trabalho manual
   repetitivo e sempre com aprovação antes de gravar.
3. **Utilidade antes de jogo.** Cronômetro, estatísticas e revisão valem mais do
   que mecânicas de jogo. Uma camada de jogo pode vir no futuro, lendo os mesmos
   dados.
4. **Partes que conversam.** A ideia de longo prazo é um hub: concluir algo num
   lugar reflete nos outros (resolver um problema no Codeforces marca o hábito,
   responder um e-mail fecha a tarefa). Isso vem aos poucos, no ritmo das fases.
5. **Ritmo real.** O projeto avança de 2 a 5 horas por semana. Escopo de cada
   versão pequeno o bastante para terminar.

## 2. Contexto

O neivi substitui o Ganbatte, um to-do gamificado anterior (servidor Python +
página web), que passa para outro mantenedor na versão 1.0. O neivi é escrito do
zero, em Go, com CLI e TUI em vez de web.

A substituição é gradual: o Ganbatte continua em uso nas partes que o neivi
ainda não cobre, e os dois convivem pelo mesmo todo.txt durante a transição
(seção 6).

## 3. Decisões fechadas

- **Linguagem:** Go.
- **Interface:** CLI e TUI no mesmo binário. O v0.1 é só CLI; a TUI (Bubble Tea
  + Lip Gloss) entra no v0.2.
- **Arquitetura:** um núcleo (lógica + dados) usado pela CLI e, depois, pela TUI.
  Nenhuma regra de negócio dentro dos comandos.
- **Onde roda:** no PC (Arch Linux). Servidor em casa e celular só no v1.0.
- **Dados:** SQLite é a fonte da verdade. O todo.txt é gerado a partir do banco
  (exportação de mão única) e nunca é lido de volta, exceto na importação.
- **Sem inbox separada:** tudo fica na lista geral; "sem disciplina e sem
  projeto" é um filtro.
- **Prioridade:** A/B/C internamente. O rótulo exibido é só um mapeamento.
- **Subtarefas:** tarefa filha (um nível só) e checklist simples dentro da tarefa.
- **Disciplinas e projetos são entidades** com tabela própria desde o v0.1 (só o
  nome). Os detalhes (professor, avaliações, status, pasta) vêm no v0.3 e no
  v0.4.
- **Sintaxe própria, não a do todo.txt:** símbolo para o que é cadastrado
  (`@` disciplina, `+` projeto); chave para características da tarefa (`pri:`,
  `due:`, `est:`). Tudo em qualquer posição da linha. Sem tags livres e sem
  contextos: só disciplinas e projetos cadastrados.
- **Gamificação leve:** streaks e estatísticas. Nada de XP, cartas ou mundo.
- **Idioma da CLI:** comandos e flags em inglês; tudo o que aparece na tela
  (mensagens, erros, ajuda) em português.
- **Biblioteca de CLI:** pacote `flag` da biblioteca padrão, com os subcomandos
  tratados à mão.
- **Idioma do código:** nomes de pacotes, tipos, funções, tabelas e colunas em
  inglês. Comentários e documentação em português.

## 4. Escopo do v0.1

**Critério de pronto:** o neivi substitui o uso básico do todo.txt no dia a dia,
com hábitos, e tem backup diário funcionando.

### 4.1 Tarefas

- Campos: título; prioridade A/B/C (opcional); prazo (opcional); disciplina
  (opcional); projeto (opcional); estimativa em minutos (opcional);
  datas de criação e de conclusão.
- A data de conclusão é salva desde o v0.1, porque o histórico e as estatísticas
  das fases futuras dependem dela.
- **Tarefa filha:** uma tarefa com `parent_id`. Um nível só.
- **Checklist:** itens de texto dentro da tarefa, cada um marcado ou não, sem
  prazo próprio.
- **Entrada:** o título primeiro, as marcas em qualquer posição:

  ```
  revisar normalização @BD pri:A due:05/11 est:30m
  escrever o parser +neivi est:1h
  ```

  | Marca | Significado |
  |---|---|
  | `@BD` | disciplina (cadastrada) |
  | `+neivi` | projeto (cadastrado) |
  | `pri:A` | prioridade A, B ou C |
  | `due:05/11/26` | prazo |
  | `est:30m` | estimativa |

- Disciplina ou projeto que não existe é erro, com a sugestão de
  cadastrar. Protege contra erro de digitação.
- **Datas de entrada** no formato `dd/mm/aa`, aceitando `/` e `-`. Ano com dois
  dígitos (`26` → `2026`). Sem ano (`05/11`), vale a próxima ocorrência da data.
  No banco e no export, sempre ISO (`2026-11-05`). Na tela, `dd/mm`.
- `est:` aceita minutos e horas: `30m`, `1h`, `1h15`, `90`. Guardado em minutos.
- **Chaves desconhecidas** (`chave:valor` que o neivi não usa, como `rep:`,
  `top:`, `aula:`, `av:`) são guardadas como texto e voltam iguais na
  exportação. Nada se perde na importação.
- **Aspas opcionais:** o texto depois do comando é juntado com espaços. No zsh,
  o alias `alias neivi='noglob neivi'` deixa `?` e `*` passarem sem aspas; aspas
  só para caracteres como `(`, `)`, `'`, `&`, `;`.
- Opcional no v0.1: atalhos de data em `due:` (`hoje`, `amanha`, `seg`…`dom`).
  Se complicar, fica para depois.

### 4.2 Visão do dia

- Mostra tarefas com prazo hoje, atrasadas e hábitos do dia.
- Ordem padrão: prioridade e depois prazo, com opção de ordenar só por prazo. É
  só uma função de ordenação.

### 4.3 Hábitos

- Tipos de marcação: feito/não feito e quantidade (meta, por exemplo 3
  problemas). O tipo tempo entra com o cronômetro, no v0.2.
- Frequências: todo dia, X vezes por semana e dias específicos da semana.
- Dias de folga por hábito (quantos dias pode pular sem zerar a sequência).
- Streak calculado a partir dos registros, nunca guardado como número solto.

### 4.4 todo.txt

- **Exportação:** gera `todo.txt` e `done.txt` a partir do banco, no formato
  padrão. O caminho é configurável (padrão: `~/.todo/`).
- **Importação:** lê um todo.txt (e opcionalmente o done.txt) para o banco. Usada
  na migração inicial e, durante a transição, para trazer tarefas novas criadas
  fora do neivi.
- Não existe sincronização automática nos dois sentidos.
- **Conversões** entre a sintaxe do neivi e a do todo.txt (que o Ganbatte e a
  automação da manhã leem):

  | neivi | todo.txt |
  |---|---|
  | `pri:A` | `(A)` no começo da linha |
  | `@BD` | `+fac.BD` |
  | `+neivi` | `+neivi` |
  | prazo | `due:2026-11-05` (ISO) |
  | (guardado em `extras`) | `@pc` e outros contextos do Ganbatte |

  A exportação converte da esquerda para a direita; a importação, da direita
  para a esquerda.

## 5. Operações do v0.1

Os nomes dos comandos e a forma de usar a CLI ainda não foram definidos (seção
13). O que a CLI precisa permitir:

- Cadastrar e listar disciplinas e projetos (só o nome).
- Criar tarefa a partir de uma linha na sintaxe do neivi.
- Listar tarefas com filtros (disciplina, projeto, sem disciplina e
  sem projeto, atrasadas); por padrão só as pendentes.
- Ver os detalhes de uma tarefa pelo id: todos os campos, checklist (numerado
  na tela, na ordem de criação) e filhas. A listagem mostra só o resumo de cada
  tarefa, com o id para usar nos outros comandos.
- Ver o dia (seção 4.2).
- Concluir e reabrir.
- Editar a partir de um novo texto, mantendo filhas e checklist.
- Remover, pedindo confirmação.
- Criar tarefa filha.
- Adicionar e marcar itens do checklist.
- Criar hábito, listar hábitos com streak e status de hoje, registrar o hábito
  no dia (com valor, para hábitos de quantidade).
- Exportar e importar todo.txt.

## 6. Transição com o Ganbatte

Durante a transição, o Ganbatte e o neivi convivem no mesmo repositório do todo
(`~/.todo`, remoto `CaioParanaiba/todo`). O Ganbatte escreve no `todo.txt`
(sugestões de aula, tarefas criadas por lá); o neivi escreve no banco.

O risco: a exportação reescreve o `todo.txt` inteiro a partir do banco e
apagaria o que o Ganbatte gravou desde o último export. Por isso:

- **Enquanto o Ganbatte estiver ativo, o export é manual** e só roda depois de
  reconciliar: trazer para o banco as tarefas novas e as conclusões feitas no
  Ganbatte. A reconciliação é um procedimento manual, fora do código.
- **Uma tarefa criada no neivi é concluída no neivi.**
- Caçadas, marcos e demais dados do `estado.json` do Ganbatte não são migrados
  no v0.1.

**Ordem de corte:**

1. **Ganbatte:** deixa de ser usado quando o neivi cobrir o fluxo de aula
   (v0.4).
2. **Resumo da manhã:** uma automação externa lê o `todo.txt` do GitHub todo dia
   às 5:50. Ela é desligada no v0.5, quando o neivi tiver o próprio resumo do
   dia com agenda e e-mail.
3. Depois disso, a exportação deixa de rodar automaticamente. Ela continua
   existindo como saída em texto puro dos dados.

## 7. Dados e backup

- **Banco:** `~/.local/share/neivi/neivi.db` (padrão XDG).
- **Configuração:** `~/.config/neivi/`.
- **Driver:** `modernc.org/sqlite` (Go puro, sem cgo) com `database/sql`.
- **Migrações:** arquivos SQL versionados aplicados na inicialização (tabela
  `schema_version`).

O banco é o único lugar onde os dados existem, então o backup faz parte do
critério de pronto do v0.1:

1. **Cópia diária consistente** com `VACUUM INTO` (ou `.backup`), nunca `cp` do
   arquivo aberto.
2. **Retenção:** últimos 7 diários e 4 semanais.
3. **Fora da máquina:** dump em texto (`sqlite3 neivi.db .dump`) commitado num
   repositório **privado**. Com o servidor (v1.0), cópia também entre servidor e
   PC.
4. **Teste de restauração** pelo menos uma vez por mês, num banco separado.

No v0.1, o backup é um script com timer do systemd. No v0.2, passa a ser feito
pelo próprio neivi.

## 8. Modelo de dados (rascunho)

Nomes de tabelas e colunas em inglês, como no código.

```
subjects  (id, name UNIQUE)               -- disciplinas: @BD
projects  (id, name UNIQUE)               -- projetos: +neivi

tasks
  id INTEGER PK AUTOINCREMENT, title TEXT,
  priority TEXT NULL ('A'|'B'|'C'; o banco aceita 'D', reservado para o futuro),
  due TEXT NULL (ISO), estimate_min INTEGER NULL,
  subject_id INTEGER NULL -> subjects.id   (apagar a disciplina é bloqueado)
  project_id INTEGER NULL -> projects.id   (idem)
  parent_id INTEGER NULL -> tasks.id       (um nível só; apagar a mãe é bloqueado)
  extras TEXT NULL                         (chaves desconhecidas, como vieram)
  created_at TEXT (ISO), done_at TEXT NULL (ISO)

checklist_items  (id, task_id, text, done INTEGER 0/1)
  -- task_id com ON DELETE CASCADE: some junto com a tarefa
  -- ordem = ordem de criação (ORDER BY id); sem reordenar no v0.1

habits
  id, name, kind ('bool'|'count'), goal INTEGER NULL,
  freq_kind ('daily'|'times_per_week'|'weekdays'), freq_value TEXT,
  slack INTEGER, created_at, archived BOOL
  -- slack: dias de folga sem zerar a sequência

habit_logs (habit_id, date DATE, value INTEGER)  -- PK(habit_id, date)
```

## 9. Estrutura do repositório

```
neivi/
  cmd/neivi/main.go      ponto de entrada
  internal/task/         struct Task + regras (atrasada, concluir...)
  internal/habit/        struct Habit + cálculo de streak com folga
  internal/todotxt/      parser da sintaxe e exportação
  internal/store/        acesso ao SQLite (única camada que fala SQL)
  internal/cli/          comandos; só chamam o núcleo e imprimem
  migrations/001_init.sql
  docs/SPEC.md
  README.md .gitignore go.mod
```

## 10. Ordem de implementação do v0.1

Marcos pequenos, um ou mais commits cada.

- **M1.** `.gitignore`, `go mod init`, repositório no GitHub, README.
- **M2.** Definir os nomes dos comandos. Store: abrir o banco, aplicar
  migração, inserir e listar tarefas. Primeiros comandos de criar e listar, com
  título só.
- **M3.** Cadastro de disciplinas e projetos. Parser da sintaxe do neivi
  (`@disciplina`, `+projeto`, `pri:`, `due:` com datas `dd/mm/aa`, `est:`,
  chaves extras) com testes (`go test`).
- **M4.** Concluir, reabrir, remover e editar.
- **M5.** Filtros da listagem e visão do dia com ordenação.
- **M6.** Tarefas filhas e checklist.
- **M7.** Hábitos: cadastro, marcação, frequências e streak com folga (com
  testes de datas).
- **M8.** Exportação e importação com as conversões da seção 4.4 (migração
  inicial do `~/.todo/todo.txt`).
- **M9.** Backup diário configurado e testado. Uso real por uma semana em
  paralelo com o Ganbatte, corrigindo o que incomodar. Tag `v0.1.0`.

## 11. Fora do v0.1

TUI, cronômetro, recorrências, detalhes de projetos e disciplinas (além do
nome), fluxo de aula, integrações, assistente, lembretes, servidor e celular.
Pomodoro está fora do plano.

## 12. Roadmap depois do v0.1

| Fase | Conteúdo |
|---|---|
| v0.2 | TUI (Bubble Tea + Lip Gloss); cronômetro por tarefa (intervalos de início e fim; um ativo por vez, pergunta antes de pausar; aviso e parada automática depois de X horas; intervalos editáveis; concluir a tarefa para o cronômetro); hábitos de tempo marcados pelo cronômetro; backup feito pelo próprio neivi |
| v0.3 | Disciplinas: sigla/tag, professor, horários, pasta de material, avaliações com peso e grupos (N1/N2), média por fórmula configurável, segunda chamada; arquivar no fim do semestre. IA monta a lista de avaliações a partir do plano de ensino |
| v0.4 | Projetos (só o nome é obrigatório; modelos que criam tarefas; status; pasta; notas em Markdown; arquivar + resumo ao encerrar) e fluxo de aula (IA por demanda; tarefas sugeridas aprovadas uma a uma; resumo em HTML). **Fim do uso do Ganbatte** |
| v0.5 | Integrações por dificuldade: Codeforces e GitHub, depois Google Calendar, depois Gmail. Resumo do dia no neivi. **Desliga a automação externa da manhã e o export automático** |
| v1.0 | Assistente (propõe e eu aprovo: agenda, e-mails, tarefas, disciplinas); caixa de e-mail com resumo; histórico e estatísticas; recorrências; lembretes no PC; servidor em casa e celular via SSH |

## 13. Decisões pendentes

| Assunto | Situação | Quando decidir |
|---|---|---|
| Atalhos de data em `due:` | `hoje`, `amanha`, `seg`…`dom`: v0.1 ou depois | M3 |
| Importar o done.txt | Traz histórico antigo para as estatísticas; custo é tratar linhas antigas mal formadas | M8 |
| Licença | Sem licença (todos os direitos reservados) ou MIT | M1 |
| Recorrências | Mesmo modelo do Ganbatte (`rep:`); detalhar | Antes da fase em que entrarem |
| Layout e navegação da TUI | Painéis, atalhos, telas | Início do v0.2 |
| Regras de tempo do cronômetro | Valor de X horas para aviso e parada | v0.2 |

## 14. Ideias descartadas

O que já foi considerado e ficou de fora, com o motivo. Serve para não reabrir a
mesma discussão.

- **Prioridade como `(A)` no começo da linha:** o título vem primeiro; virou
  `pri:`. O `(A)` existe só na conversão com o todo.txt.
- **`+fac.BD` para disciplina:** era um projeto fantasiado de disciplina.
  Disciplina virou entidade, marcada com `@`.
- **Contextos (`@pc`, `@caderno`):** quase toda tarefa tinha o mesmo contexto,
  então não filtravam nada. Voltam (como `ctx:`) se fizerem falta. O `@` ficou
  para disciplina.
- **Tags livres:** sem uso claro além de disciplina e projeto. Voltam se fizerem
  falta.
- **todo.txt como fonte da verdade:** não guarda bem filhas, checklist, hábitos
  e histórico; virou só exportação.
- **Inbox separada:** obriga a esvaziar; "sem disciplina e sem projeto" é um
  filtro.
- **Interface web:** o terminal é onde eu trabalho; a TUI resolve o visual.
- **`cobra`:** poucos comandos; o pacote `flag` basta.
- **XP, cartas e mundo:** utilidade antes de jogo.
- **Pomodoro:** o cronômetro por tarefa cobre o que importa.
- **Pensar em outros usuários:** instalador, temporadas e suporte foram o que
  pesou no Ganbatte.
