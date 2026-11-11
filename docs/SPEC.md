# neivi: especificação

Documento vivo. Versão de 09/11/2026.

**O que o neivi faz e por quê** (escopo, cenários de uso, regras, decisões e o
que fica fora de cada versão) está em [requisitos.md](requisitos.md). Este
documento descreve **como**: sintaxe, arquitetura, dados e comandos.

**Sobre este documento:** a especificação é escrita e mantida principalmente
com uma IA (Claude). As ideias e as decisões são minhas; a IA ajuda a
organizá-las, aponta lacunas e redige o texto. O código segue outro modelo: é
escrito por mim, e a IA atua como mentora: explica conceitos, revisa e recomenda
materiais de estudo. No máximo, quando eu peço, faz alguma parte mecânica e
repetitiva para poupar tempo, nunca funcionalidade nova.

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
   responder um e-mail fecha a tarefa). Isso vem aos poucos.
5. **Ritmo real.** O projeto avança de 2 a 5 horas por semana. Escopo de cada
   versão pequeno o bastante para terminar.

## 2. Contexto

O neivi substitui o Ganbatte, um to-do gamificado anterior (servidor Python +
página web), que passa para outro mantenedor na versão 1.0. O neivi é escrito do
zero, em Go, com CLI e, depois, TUI em vez de web.

A troca é um corte seco: quando o v0.1 estiver pronto, os dados do Ganbatte são
importados uma vez e ele é desligado (seção 6).

## 3. Decisões técnicas

- **Linguagem:** Go.
- **Interface:** CLI e TUI no mesmo binário. O v0.1 é só CLI; a TUI (Bubble Tea
  + Lip Gloss) entra depois.
- **Arquitetura:** um núcleo (lógica + dados) usado pela CLI e, depois, pela TUI.
  Nenhuma regra de negócio dentro dos comandos.
- **Onde roda:** no PC (Arch Linux).
- **Dados:** SQLite é a fonte da verdade e o único lugar onde os dados existem.
  Não há exportação para todo.txt; o todo.txt só é lido na importação.
- **Sem inbox separada:** a tarefa incompleta (sem prioridade) é deduzida, e
  "sem disciplina e sem projeto" é um filtro.
- **Prioridade:** A, B, C e D. D é uma prioridade definida (baixa), diferente de
  vazio (ainda não decidida).
- **Disciplinas e projetos são entidades** com tabela própria (código curto,
  como `BD`, e nome opcional). O código é o que vai depois do `@` ou do `+`.
- **Sintaxe própria, não a do todo.txt:** símbolo para o que é cadastrado
  (`@` disciplina, `+` projeto); chave para características da tarefa (`pri:`,
  `due:`, `est:`). Tudo em qualquer posição da linha. Sem tags livres e sem
  contextos: só disciplinas e projetos cadastrados.
- **Deduzir em vez de guardar:** tarefa nova, tarefa incompleta, tipo do
  projeto e, no futuro, a sequência dos hábitos saem dos dados, sem campo
  próprio.
- **Idioma da CLI:** comandos e flags em inglês; tudo o que aparece na tela
  (mensagens, erros, ajuda) em português.
- **Biblioteca de CLI:** pacote `flag` da biblioteca padrão, com os subcomandos
  tratados à mão.
- **Formato dos comandos:** tarefas no nível de cima (`neivi add`, `neivi ls`),
  por serem o uso diário. O resto em grupos, substantivo e depois verbo
  (`neivi subject add`, `neivi project ls`).
- **Idioma do código:** nomes de pacotes, tipos, funções, tabelas e colunas em
  inglês. Comentários e documentação em português.

## 4. Sintaxe da tarefa

- **Campos:** título; descrição (opcional); prioridade A–D (opcional); prazo
  (opcional); disciplina (opcional); projeto (opcional); estimativa em minutos
  (opcional); datas de criação e de conclusão (gravadas pelo sistema).
- **Entrada:** o título primeiro, as marcas em qualquer posição:

  ```
  revisar normalização @BD pri:A due:05/11 est:30m
  escrever o parser +neivi est:1h
  ```

  | Marca | Significado |
  |---|---|
  | `@BD` | disciplina (cadastrada) |
  | `+neivi` | projeto (cadastrado) |
  | `pri:A` | prioridade A, B, C ou D |
  | `due:05/11/26` | prazo |
  | `est:30m` | estimativa |
  | `-- texto` | descrição |

- Disciplina ou projeto que não existe é erro, com a sugestão de
  cadastrar. Protege contra erro de digitação.
- **Datas de entrada** no formato `dd/mm/aa`, aceitando `/` e `-`. Ano com dois
  ou quatro dígitos (`26` ou `2026`). Sem ano (`05/11`), vale a próxima
  ocorrência da data: hoje conta como hoje; data que já passou vai para o ano
  seguinte. `29/02` sem ano só vale se esse ano for bissexto; senão é erro.
  No banco, sempre ISO (`2026-11-05`). Na tela, `dd/mm`.
- `est:` aceita minutos e horas: `30m`, `1h`, `1h15`, `90`. Guardado em minutos.
- `pri:` aceita minúscula (`pri:a` vira `A`).
- **Outras palavras com `:`** (`aula:3`, `verificar:`, `10:30`, links) fazem
  parte do título.
- **Aspas opcionais:** o texto depois do comando é juntado com espaços. No zsh,
  o alias `alias neivi='noglob neivi'` deixa `?` e `*` passarem sem aspas; aspas
  só para caracteres como `(`, `)`, `'`, `&`, `;`.
- **Atalhos de data** em `due:`: `hoje`, `amanha`, `seg`…`dom` (também
  `segunda`…`domingo`, sem acento). Dia da semana é sempre o próximo: `seg`
  numa segunda é a segunda seguinte.
- **Descrição:** o título fica curto e o detalhe vai depois de `--`. O primeiro
  `--` abre a descrição; um segundo `--` fecha, e o que vem depois volta a ser
  lido como marcas. Sem o segundo, a descrição vai até o fim da linha. Dentro
  dela nada é interpretado (`@`, `+`, `:` são texto). A descrição não pode
  conter `--`. Aparece nos detalhes da tarefa, não na listagem.

  ```
  neivi add fazer trabalho @BD -- conforme o slide 2
  neivi add fazer trabalho -- conforme o slide 2 -- @BD pri:A
  ```

- **Edição parcial:** a edição usa a mesma sintaxe da entrada, e só muda o que
  aparece no texto; o resto fica como está. Texto solto vira o título novo;
  sem texto solto, o título não muda. O valor `-` apaga o campo (`pri:-`,
  `due:-`, `est:-`, `@-`, `+-`, `-- -`). O título não pode ser apagado.

  ```
  neivi edit 1 pri:B
  neivi edit 1 @BD due:sex
  neivi edit 1 due:-
  ```

## 5. Comandos

Comandos já definidos:

| Comando | Faz |
|---|---|
| `neivi add <texto>` | cria uma tarefa |
| `neivi ls [filtros] [palavras]` | lista as tarefas (ver abaixo) |
| `neivi today` | resumo diário (ver abaixo) |
| `neivi triage` | tarefas incompletas (sem prioridade), as mais antigas primeiro |
| `neivi show <id>` | mostra todos os campos de uma tarefa |
| `neivi done <id>... [ontem]` | conclui uma ou mais tarefas, hoje ou ontem |
| `neivi reopen <id>...` | reabre uma ou mais tarefas concluídas |
| `neivi edit <id> <texto>` | muda só os campos que aparecem no texto (seção 4) |
| `neivi rm <id>...` | remove uma ou mais tarefas, pedindo confirmação de cada uma |
| `neivi subject add <código> [nome]` | cadastra uma disciplina (ativa) |
| `neivi subject ls [all]` | lista as ativas (com `all`, também as encerradas) |
| `neivi subject edit <código> ...` | muda nome, `sit:ativa`/`sit:encerrada`, `code:novo` |
| `neivi subject rm <código>` | remove, só sem tarefas nem projetos ligados |
| `neivi project add <código> [nome] ...` | cadastra um projeto: `@BD`, `gh:usuario/repo`, `dir:pasta`, `cont:sim`, `sit:...`, `-- descrição` |
| `neivi project ls [all]` | lista ativos e pausados (com `all`, todos) com progresso e última atividade |
| `neivi project show <código>` | panorama do projeto (C9) com as pendentes |
| `neivi project edit <código> ...` | muda só o que aparece; `-` apaga; `code:novo` troca o código |
| `neivi project rm <código>` | remove, só sem tarefas |

**Saída dos comandos de tarefa** (requisitos RNF-04 e RNF-05):

- Confirmação numa linha curta: `+` criar, `~` editar, `✓` concluir, `↺`
  reabrir, `✗` remover; depois `id · prioridade · prazo · EST · título
  @disciplina +projeto …`, sem os campos vazios. Prazo como `sex 13/11`, com o
  ano só quando não é o atual. O `done` termina com `hoje: N concluidas` (ou
  `ontem:`).
- Avisos e erros dizem o que está errado e como resolver. Avisos (duplicada,
  concluída, id que não existe) não impedem a ação; erros não mudam nada.
- `done`, `reopen` e `rm` leem todos os ids antes de agir: id que não existe
  só avisa; palavra que não é id (nem `ontem`, no `done`) recusa tudo.

**Listagem.** Sem nada, mostra as pendentes na ordem padrão. Filtros e
palavras usam a mesma sintaxe da entrada e se combinam:

| Parte | Efeito |
|---|---|
| `@BD`, `+neivi`, `pri:A` | só essa disciplina, projeto ou prioridade |
| `@`, `+` | qualquer disciplina, qualquer projeto |
| `@-`, `+-`, `pri:-` | sem disciplina, sem projeto, sem prioridade (`pri:-` = `triage`) |
| `new` | só as novas (criadas hoje ou ontem) |
| `all` | inclui as concluídas, depois das pendentes |
| `done` | só as concluídas |
| `late` | só as atrasadas |
| `due` | ordena só por prazo |
| `desc` | mostra a descrição embaixo de cada tarefa, recuada |
| outra palavra | busca no título (todas as palavras, sem maiúsculas nem acentos) |

- Mesmo tipo é "ou" (`@BD @ED` = BD ou ED); tipos diferentes são "e"
  (`@BD pri:A`). Disciplina ou projeto que não existe é erro.
- `all`, `done`, `late`, `new`, `due` e `desc` são reservadas: não servem para
  buscar no título. `all` mostra as concluídas no fim, das mais recentes para
  as mais antigas, com a coluna `X`; `done` mostra só elas, na mesma ordem.
- Sem `desc`, a tarefa com descrição ganha um `…` depois do título.
- Ordens (requisitos D-04 e D-28): padrão = prioridade (A–D, sem prioridade por
  último) → prazo (mais perto primeiro, sem prazo por último) → estimativa
  (mais rápidas primeiro) → id. Com `due` e no resumo diário: prazo →
  prioridade → estimativa → id.
- Filtros e ordenação ficam fora da CLI (store e `internal/task`), para a TUI
  usar os mesmos.

**Resumo diário (`today`).** Três blocos de tarefas pendentes, nesta ordem
(requisitos RN-16): atrasadas, urgentes (prazo hoje, qualquer prioridade) e
importantes (prioridade A sem prazo); sem nenhuma A sem prazo, o terceiro vira
"sem prazo" com o nível mais alto que existir. Cada bloco é uma tabela com o
próprio cabeçalho, todas com a mesma largura de coluna. No fim, `falta: ~4h`
(soma das EST) e o contador do dia. Dia vazio: `nada para hoje`, com a dica do
`triage` se houver incompletas. Os hábitos do dia entram abaixo.

**Listas em tabela** (`today`, `ls`, `triage`): `ID P PRAZO EST TAREFA`, o `X`
só no `ls all`, prazo como `qua 11/11`. A tabela é montada à mão
(`cmd/neivi/views.go`) e não com o `tabwriter`, porque o `tabwriter` não mantém
a largura entre os grupos do `today`.

**Detalhes (`show`).** Todos os campos, um por linha (vazios como `-`), com o
prazo deduzido (`em 2 dias`, `atrasada ha 3 dias`), `(nova)` na criação, a
situação (pendente, incompleta, concluída em…) e a descrição inteira no fim.

Os demais comandos seguem o formato da seção 3 e entram aqui quando forem
definidos, fatia a fatia. O que ainda falta definir:

- Campos de projeto da D-06 (descrição, GitHub, situação, disciplina, pasta) no
  `project add` e o panorama do projeto (C9).
- Criar hábito, listar hábitos com o status de hoje e registrar o hábito em
  hoje ou ontem, com valor.
- Importação única do todo.txt e das disciplinas do Ganbatte.

## 6. Corte com o Ganbatte

Sem convivência: o neivi não escreve no `todo.txt`.

Um corte só, com o v0.1 completo (requisitos D-35):

1. Importar uma vez (tarefas e disciplinas), conferir e desligar o
   `ganbatte.service`. Tarefas e hábitos passam a ser só no neivi. As chaves
   do Ganbatte que o neivi não usa (`aula:`, `av:`, `top:`, `gh:`, `esf:`,
   `rec:`, `rep:`) vão para a descrição da tarefa importada.
2. Nesse dia, a automação externa das 5:50 troca o prompt por algo neutro (o
   `todo.txt` para de ser atualizado) e fica só como gatilho do dia, até a TUI.
4. Dados que ficam para o v0.2: avaliações e notas (`hunter/avaliacoes.txt`).
   Caçadas, marcos e demais dados do `estado.json` não são migrados.

## 7. Dados e backup

- **Banco:** `$XDG_DATA_HOME/neivi/neivi.db` (padrão
  `~/.local/share/neivi/neivi.db`); a pasta é criada na primeira execução. A
  variável `NEIVI_DB` troca o caminho, para testar sem tocar no banco real
  (`NEIVI_DB=/tmp/neivi-teste.db neivi ls`).
- **Configuração:** `~/.config/neivi/`.
- **Driver:** `modernc.org/sqlite` (Go puro, sem cgo) com `database/sql`.
- **Esquema:** `001_init.sql` embutido no binário (`embed`), com
  `CREATE TABLE IF NOT EXISTS`, executado sempre que o banco é aberto. Enquanto
  não houver dados reais, o esquema muda direto no `001`. Depois do corte, o
  `001` não muda mais (o `IF NOT EXISTS` não altera tabela que já existe):
  toda mudança vira uma migração nova. O controle de versões
  (`PRAGMA user_version`, só aplicar o que falta) entra junto com a primeira
  migração `002`, quando ela for necessária.
- **Chaves estrangeiras** ligadas em toda conexão (`_pragma=foreign_keys(1)` no
  caminho do banco).

O banco é o único lugar onde os dados existem, então o backup faz parte do
pronto do v0.1 (requisitos D-29):

1. **Cópia diária consistente** com `VACUUM INTO` (ou `.backup`), nunca `cp` do
   arquivo aberto.
2. **Retenção:** últimos 7 diários e 4 semanais.
3. **Fora da máquina:** dump em texto (`sqlite3 neivi.db .dump`), que volta
   para um banco novo sem conversão. Até existir o servidor: todo dia no
   repositório privado `todo` (em `neivi/`), sem criptografia (requisitos
   D-33), e uma vez por mês no HD externo. Com o servidor, ele faz a cópia
   diária e o GitHub fica espaçado.
4. **Teste de restauração** pelo menos uma vez por mês, num banco separado.

No v0.1, o backup é um script (`~/.local/bin/neivi-backup`) com timer do
systemd do usuário (`neivi-backup.timer`, todo dia às 12:00, `Persistent=true`
para rodar depois se o PC estava desligado). As cópias ficam em
`~/.local/share/neivi/backup/` (`diario/`, `semanal/` e o dump `neivi.sql`).

## 8. Modelo de dados (rascunho)

Nomes de tabelas e colunas em inglês, como no código.

```
subjects  (id, code UNIQUE, name, status 'active'|'closed')   -- disciplinas: @BD
projects  (id, code UNIQUE, name, description NULL, github NULL, dir NULL,
           status 'active'|'paused'|'done'|'canceled', continuous 0/1,
           subject_id NULL -> subjects.id, created_at)          -- projetos: +neivi
  -- o código é o que vai depois do @ ou do +; sem nome, a CLI grava o código
  -- valores de status em inglês no banco; na tela e na entrada, em português
  -- progresso e última atividade são contados nas consultas, nunca guardados
  -- disciplina efetiva da tarefa = COALESCE(tasks.subject_id, projects.subject_id)

tasks
  id INTEGER PK AUTOINCREMENT, title TEXT, description TEXT NULL,
  priority TEXT NULL ('A'|'B'|'C'|'D'),
  due TEXT NULL (ISO), estimate_min INTEGER NULL,
  subject_id INTEGER NULL -> subjects.id   (apagar a disciplina é bloqueado)
  project_id INTEGER NULL -> projects.id   (idem)
  created_at TEXT (ISO), done_at TEXT NULL (ISO)

habits
  id, name, goal INTEGER NOT NULL DEFAULT 1,
  freq_kind ('daily'|'times_per_week'|'weekdays'), freq_value TEXT,
  created_at, archived BOOL
  -- todo hábito tem meta numérica; feito/não feito é meta 1 (requisitos D-02)

habit_logs (habit_id, date DATE, value INTEGER)  -- PK(habit_id, date)
  -- o valor real, mesmo acima da meta (D-12); a sequência sai daqui no futuro
```

## 9. Estrutura do repositório

```
neivi/
  cmd/neivi/             ponto de entrada e comandos
  internal/task/         struct Task + regras (atrasada, ordenação...)
  internal/subject/      struct Subject (disciplina)
  internal/project/      struct Project (projeto)
  internal/habit/        struct Habit + regras dos hábitos
  internal/todotxt/      parser da sintaxe do neivi (e a leitura do todo.txt
                         na importação)
  internal/store/        acesso ao SQLite (única camada que fala SQL)
    migrations/001_init.sql   (dentro do store por causa do embed)
  docs/SPEC.md, docs/requisitos.md
  README.md .gitignore go.mod
```

## 10. Ordem de implementação do v0.1

Feito:

- **M1.** `.gitignore`, `go mod init`, repositório no GitHub, README.
- **M2.** Nomes dos comandos; store; criar e listar tarefas.
- **M3.** Cadastro de disciplinas e projetos; parser da sintaxe com testes.
- **M4.** Concluir, reabrir, remover e editar; descrição.

Daqui em diante, por fatias (cada uma: requisitos → código → testes):

1. Tarefas (retroativa). Fechada em 11/11/2026.
2. Listagem, filtros e resumo diário (fatia 4 dos requisitos). Fechada em
   11/11/2026.
3. Disciplinas e projetos (fatia 2: D-06, panorama do projeto).
4. Hábitos (fatia 3).
5. Importação (fatia 5; o backup já está pronto).
6. **Corte** (seção 6) e uma semana de uso real. Tag `v0.1.0`.

As fatias que mexem no banco (2 e 3) vêm antes da importação, para o esquema
ficar certo antes de existirem dados reais, sem migração (requisitos D-35); os
números das fatias nos requisitos continuam os mesmos.

## 11. Decisões pendentes

| Assunto | Situação | Quando decidir |
|---|---|---|
| Licença | Sem licença (todos os direitos reservados) ou MIT | Quando quiser |
| Nome do pacote `internal/todotxt` | Hoje é o parser da sintaxe do neivi; trocar o nome é opcional | Fatia 1 |
| Layout e navegação da TUI | Painéis, atalhos, telas | Início da TUI |
| Edição pelo `$EDITOR` | Edição sem texto abre a linha completa da tarefa no editor (como o `git commit`); depende de converter a tarefa de volta para texto | Depois do v0.1 |
| Edição por campos numerados | Mostra a tarefa com os campos numerados e o usuário escolhe quais editar; é interface interativa | Layout da TUI |

## 12. Ideias descartadas

O que já foi considerado e ficou de fora, com o motivo. Serve para não reabrir a
mesma discussão.

- **Prioridade como `(A)` no começo da linha:** o título vem primeiro; virou
  `pri:`. O `(A)` só é lido na importação.
- **`+fac.BD` para disciplina:** era um projeto fantasiado de disciplina.
  Disciplina virou entidade, marcada com `@`.
- **Contextos (`@pc`, `@caderno`):** quase toda tarefa tinha o mesmo contexto,
  então não filtravam nada. Voltam (como `ctx:`) se fizerem falta. O `@` ficou
  para disciplina.
- **Tags livres:** sem uso claro além de disciplina e projeto. Voltam se fizerem
  falta.
- **todo.txt como fonte da verdade:** não guarda bem relacionamentos, hábitos e
  histórico.
- **Exportação para todo.txt:** servia ao Ganbatte e à automação da manhã, não
  ao backup; com o corte seco, não tem mais uso (requisitos D-24).
- **Campo `extras` para as chaves do Ganbatte:** só existia para elas voltarem
  iguais no export; na importação vão para a descrição (requisitos D-31).
- **Inbox separada:** obriga a esvaziar; a tarefa incompleta é deduzida.
- **Interface web:** o terminal é onde eu trabalho; a TUI resolve o visual.
- **`cobra`:** poucos comandos; o pacote `flag` basta.
- **XP, cartas e mundo:** utilidade antes de jogo.
- **Pomodoro:** o cronômetro por tarefa (futuro) cobre o que importa.
- **Pensar em outros usuários:** instalador, temporadas e suporte foram o que
  pesou no Ganbatte.
