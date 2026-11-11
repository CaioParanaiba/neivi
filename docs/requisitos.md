# Requisitos do Neivi

Cobre: v0.1 · Versao do doc: 0.5 · 2026-11-11

| Versao    | data          | Mudanca                   |
|-----------|---------------|---------------------------|
| 0.1       | 2026-11-07    | Introducao                |
| 0.2       | 2026-11-09    | Contexto 2.1 e decisoes (secao 8) |
| 0.3       | 2026-11-09    | Contexto completo (2.2 e 2.3) |
| 0.4       | 2026-11-09    | Revisao contra o SPEC: escopo, glossario e decisoes D-22 a D-32; 1 e 2 fechados |
| 0.5       | 2026-11-11    | Fatia 1 (tarefas) fechada: US-01 a US-06, RN-01 a RN-15, RNF-01 a RNF-05, modelo conceitual, D-33 e rastreabilidade |

# 1. Introducao

## 1.1 Proposito do documento

Especifica os requisitos do projeto Neivi da v0.1, para guiar a implementação e os testes.

## 1.2 Problema e visao

**Problema**: Eu usava o TODO.txt para me organizar em todas as minhas tarefas e habitos diarios, com o tempo vi a necessidade de criar um sistema maior, desenvolvendo o Ganbatte (por nao ser um projeto pessoal, foi muito utilizada a AI), que acabou sendo usado por mais de uma pessoa e um grupo de amigos onde foi criado um sistema ate de multiplayer envolvendo a plataforma web.

- **P1:** Ao anotar tarefas no TODO.txt faltava a capacidade de diferenciar dados como habitos e tarefas;
- **P2:** TODO.txt limitava o uso de mais funcionalidades, gerenciamento de habitos, acompanhamento de disciplinas e projetos... Nao era robusto o suficiente;
- **P3:** Ao desenvolver o Ganbatte para mais pessoas o sistema perdeu a minha identidade e as funcionalidades viraram genericas e muitos desejos pessoais nao podiam entrar no projeto;
- **P4:** Ganbatte se tornou muito complexo, com codigos onde eu nao poderia contribuir;
- **P5:** Meu sistema era completamente dependente do uso de AI; como programador e estudante, estava deixando de lado um possivel treinamento e desenvolvimento de software por comodismo;
- **P6:** Possuia muita gamificacao, perdendo o foco do que realmente importa, o proprio sistema de tarefas, gastando mais tempo desenvolvendo gamificacoes do que de fato usando o sistema;


**Para quem**: eu como: pessoa, estudante, programador e maratonista.

**Visao**: com a neivi, eu consigo ter um sistema robusto, com funcionalidades diversas alem de apenas um simples TODO formato txt, mas nao complexo, me permitindo gerir mudancas e entender todos os processos, onde posso registrar, visualizar e analisar minhas tarefas e habitos do meu dia a dia. Pensado para mim e desenvolvido por mim, um hub dedicado para minhas necessidades, para gerenciamento de vida academica + organizacao de projetos pessoais e desenvolvimento das habilidades de programacao competitiva.

## 1.3 Escopo

**Pronto (v0.1)**: uso o neivi todo dia no lugar do TODO.txt para tarefas e habitos, com o resumo diario, e com backup funcionando. (D-22)

**Dentro**:
- CRUD de tarefas, disciplinas e projetos;                                       (C2, C8)
- CRUD de habitos com meta e periodicidade (diario, dias especificos, N por semana);  (C5, C11)
- Registro rapido de tarefas com o minimo de campos, completando depois;                 (C2)
- Descricao da tarefa (detalhe que fica fora da listagem);                                (C4)
- Vinculo de tarefas a disciplinas e projetos;                                            (C3, C7)
- Concluir e reabrir tarefas;                                                             (C4)
- Registro de habitos realizados (hoje ou ontem), por quantidade, podendo passar da meta;  (C5, C11)
- Listagem de tarefas, habitos, disciplinas e projetos com filtros e ordenacao;           (C4, C5, C7, C10, C11)
- Historico de tarefas concluidas e habitos realizados;                                   (C9, C10)
- Resumo diario (leitura rapida do dia: atrasadas, urgentes e importantes);              (C1)
- Panorama do projeto, com o link do github (so hyperlink);                              (C9)
- Importacao unica das tarefas e disciplinas do TODO.txt e do Ganbatte (arquivos locais), antes de desligar o Ganbatte; (D-24)
- Backup diario automatico (RNF, secao 5).                                               (D-29)

**Fora (v0.1)**:
- Avaliacoes e notas (vao para o v0.2, antes do semestre 2027-1);                         (C3, C6, D-23)
- Tarefas filhas e checklist;                                                             (D-25)
- Sequencia (streak) e dias de folga dos habitos (entram com a TUI);                      (D-26)
- Exportacao para o TODO.txt;                                                             (D-24)
- Visao propria do C11 juntando habito e projeto (entra com a TUI);                       (D-30)
- Interface em TUI;
- Claude incluido como assistente pessoal;
- Transcricao de audios de aula e criacao de resumos;
- Cadastro de resumos de aula;
- Recomendacao de tarefas de estudo a partir de aulas;
- Vinculacao ao Google Calendario;
- Calendario de avaliacoes (visao por data, alertas);
- Plano de avaliacao da disciplina com pesos, N1/N2 e calculo de media;                   (C6)
- Painel de planejamento do dia (panorama para reorganizar o dia);                        (C8)
- Painel de planejamento semanal;                                                         (C10)
- Inbox (lugar separado para as tarefas incompletas);                                     (C2)
- Contagem automatica das questoes do Codeforces nos habitos;                            (C11)
- Ligar um habito a um projeto;
- Automatizar prioridade com base em prazo e EST;
- Tarefas recorrentes;
- Servidor proprio acessado pelos meus computadores (sincronizacao e backup no servidor);

**Fora Sempre**:
- Sistema de gamificacao avancado e altamente personalizavel
- Sistema multiplayer e pensado para outras pessoas
- Sistema em formato html/web


## 1.4 Glossário

### Sistema

| Termo    | Definicao |
|----------|-----------|
| TODO     | Sistema de gerenciamento de tarefas |
| Hub      | Centro de conexao de varias frentes, como TODO, maratona, habitos e assistente de AI |
| Robusto  | Existe mais de um tipo de dado, os dados se relacionam entre si e eh possivel extrair informacoes importantes de um conjunto de dados |
| Complexo | Alteracoes simples sao dificeis, o sistema eh muito dificil de entender e tem mais funcionalidades do que de fato precisa |

### Entidades

| Termo      | Definicao |
|------------|-----------|
| Tarefa     | Atividade que precisa ser feita uma unica vez e, depois de feita, eh concluida |
| Habito     | Atividade recorrente, feita pelo menos 1 vez por semana; diferente da tarefa, nao eh concluida, e sim marcada com metas em um intervalo de tempo, por isso fica separada das tarefas comuns. Nao atrasa e nao acumula: se nao for feito, so nao bate a meta naquele periodo (diferente de uma tarefa recorrente, que continua devendo) |
| Disciplina | Materia academica da faculdade, a qual se vinculam tarefas, avaliacoes e notas |
| Projeto    | Projeto individual focado em programacao, podendo ser pessoal ou academico, ao qual se vinculam tarefas, podendo ser continuo ou possuir um fim determinado/indeterminado. Pode ter uma disciplina (opcional); com disciplina eh academico, sem disciplina eh pessoal. Situacao: ativo, pausado, concluido ou cancelado |
| Avaliacao  | (v0.2) Atividade avaliativa de uma disciplina (prova, trabalho, lista...), com data e conteudos (texto livre), que recebe uma nota |
| Nota       | (v0.2) Valor obtido em uma avaliacao |
| Resumo     | Resumo ou anotacao referente a uma aula (futuro, fora do v0.1) |

### Atributos do habito

| Termo          | Definicao |
|----------------|-----------|
| Meta           | Objetivo de curto prazo responsavel por metrificar um habito, podendo ser diaria, em dias especificos ou uma quantidade determinada durante a semana |
| Periodicidade  | Quando a meta do habito vale: diario, dias especificos (ex.: seg/qua/sex) ou N vezes por semana (dias livres) |

### Atributos e estados da tarefa

| Termo             | Definicao |
|-------------------|-----------|
| Prazo             | Data limite para realizar uma tarefa |
| EST               | Estimativa de tempo para realizar uma tarefa |
| Descricao         | Detalhe da tarefa, visto so nos detalhes, fora da listagem |
| Prioridade        | Nivel de importancia de uma tarefa (A a D), usado para decidir o que fazer primeiro; D eh uma prioridade definida (baixa), diferente de vazio (ainda nao decidida); criterio de cada nivel na RN-15 |
| Atrasada          | Tarefa nao concluida cuja a data ja passou, sem considerar o dia atual |
| Urgente           | Tarefa nao concluida com prazo no dia atual, de qualquer prioridade (D-27) |
| Importante        | Tarefa nao concluida de prioridade A e sem prazo (D-27) |
| Tarefa incompleta | Tarefa sem prioridade definida, registrada de forma rapida para completar na revisao do dia |
| Tarefa pronta     | Tarefa com prioridade definida (A a D) |
| Tarefa nova       | Tarefa criada hoje ou ontem (deduzido pela data de criacao, sem marcacao) |

### Tempo

| Termo  | Definicao |
|--------|-----------|
| Dia    | Periodo de 24 horas, das 00:00 as 23:59 |
| Semana | Periodo de domingo a sabado |
| Mes    | Do dia 01 ao ultimo dia do respectivo mes, conforme o calendario |

### Visoes

| Termo                          | Definicao |
|--------------------------------|-----------|
| Resumo diario                  | Visao rapida, so de leitura, do que tem para hoje: tarefas atrasadas, urgentes e importantes, nessa ordem (C1) |
| Painel de planejamento do dia  | Panorama do dia para reorganizar: tarefas novas e antigas de todas as disciplinas, ajustando prioridades (C8) |
| Painel de planejamento semanal | Panorama da semana: o que foi feito e o que esta por vir na semana/mes (C10) |
| Inbox                          | Local onde ficam as tarefas incompletas, que ainda nao foram transformadas em tarefas prontas para serem realizadas. Nao usado no neivi: a tarefa incompleta eh deduzida (D-10) |

# 2. Contexto
## 2.1 Usuário e contextos de uso

Um estudante de ciencia da computacao, que faz maratona de programacao juntamente com projetos da faculdade e pessoais, utiliza no dia a dia seu notebook com Arch Linux em todo lugar que vai sempre ao lado do terminal e neovim. Seu objetivo eh organizar o dia a dia.

- **C1 - Inicio do dia:** ao ligar o PC, de forma rapida (~1 min). Leio as tarefas atrasadas e as urgentes para ja me preparar para o dia.
  -> Exige ver atrasadas e urgentes juntas, de uma vez, sem montar filtro.
- **C2 - Registro de tarefas:** apos lembrar/receber uma nova tarefa, de forma instantanea. Se a tarefa nao apresentar uma urgencia grande, eh registrada de forma simples (deixando campos pendentes para a proxima revisao do dia); do contrario, de forma completa (com todos os campos).
  -> Exige um preenchimento que difere o simples do completo: registrar com o minimo de campos e reencontrar depois as tarefas incompletas, na revisao do dia.
- **C3 - Professor marcou a prova:** durante a aula, com pressa. Registro a data e os conteudos da prova.
  -> Exige registrar em poucos segundos, ligado a disciplina.
  -> Atendido no v0.2 (D-23).
- **C4 - Consultas e conclusoes:** apos terminar uma tarefa (~2 min). Marco a tarefa concluida como feita e verifico quais tarefas ainda faltam ser feitas no dia; com a EST e a prioridade, defino qual a proxima tarefa que vou fazer.
  -> Exige visualizacao facil das tarefas e de seus principais atributos: EST, prioridade, atrasadas.
- **C5 - Acompanhamento de habitos:** ao longo do dia, com pressa. Vejo quais habitos ainda faltam ser feitos e registro algum ja realizado.
  -> Exige visualizar rapidamente os habitos pendentes e marcar os ja feitos de forma simples.
- **C6 - Registro de notas:** apos receber a nota de uma avaliacao, com pressa (~5 min). Registro a nota na respectiva avaliacao.
  -> Exige o registro facil das notas.
  -> Atendido no v0.2 (D-23).
- **C7 - Organizar aulas:** apos a ultima aula do dia, com calma (~5 min). Revejo as tarefas novas associadas as ultimas aulas, para definir novas prioridades.
  -> Exige a visualizacao das tarefas novas relacionadas as ultimas aulas.
- **C8 - Revisao do dia:** no primeiro bloco livre do dia, com calma (~15 min). Comparo as novas prioridades com as antigas e faco os ultimos ajustes que preciso para comecar o dia.
  -> Exige uma visualizacao das tarefas novas do dia junto com as antigas e uma edicao rapida das prioridades.
- **C9 - Acompanhamento de Projetos:** durante um bloco ou tarefa que envolva um projeto, com calma (~15 min). Me contextualizo sobre o projeto, revisando do que se trata, juntamente com a ultima tarefa feita associada e uma visualizacao das pendentes.
  -> Exige um panorama de todo o projeto, com suas informacoes e caracteristicas (hyperlink para o github se houver), e vislumbre das suas tarefas ja feitas e as pendentes.
- **C10 - Revisao semanal:** no domingo em casa, com calma (~30 min). Analiso as tarefas atrasadas e pendentes junto com as proximas avaliacoes para redefinir prioridades.
  -> Exige um panorama de tudo que foi feito na semana, com tudo que esta por vir na semana/mes.
  -> As proximas avaliacoes entram no v0.2 (D-23).
- **C11 - Foco na Maratona:** no bloco de maratona, com rapidez (~2 min). Verifico quantas questoes do Codeforces foram realizadas e qual o conteudo que irei focar hoje.
  -> Exige consultar rapidamente o habito de questoes do Codeforces juntamente com o proximo conteudo de estudo.

## 2.2 Restrições

- **R-01 Linguagem:** Go.
  Porque: para eu ter que aprender uma linguagem nova, sem cair na tentacao de copiar codigos ja escritos antes.
- **R-02 Banco de Dados:** SQLite.
  Porque: para eu aplicar os conhecimentos adquiridos na disciplina *Projeto de Banco de Dados*.
- **R-03 Maquina e Dados:** Maquinas proprias.
  Porque: ao nao utilizar nenhum tipo de terceiro nao ha problemas de privacidade, e permite a expansao futura para um servidor proprio (futuro, fora do escopo).
- **R-04 Interface:** Terminal: CLI no v0.1, TUI no futuro.
  Porque: para nao ter que aprender mais uma linguagem nova para interface web/software, alem do terminal fazer parte do cotidiano.
- **R-05 Uso:** Unicamente pessoal, um usuario so.
  Porque: fazer um sistema com as minhas demandas apenas (P3).
- **R-06 Desenvolvimento:** Codigo escrito por mim, com a IA apenas como mentora.
  Porque: permite eu me autoaperfeicoar e nao depender de IA (P4, P5).
- **R-07 Internet:** Nenhuma funcionalidade do v0.1 precisa de internet.
  Porque: o neivi funciona em qualquer lugar, inclusive em aula sem rede (S-01).


## 2.3 Suposições e dependências

### Suposicoes

- **S-01:** O notebook esta comigo em todos os momentos de uso.
  Se mudar: C2, C3 e C5 ficam sem registro na hora; anoto em outro lugar e passo depois.
- **S-02:** O relogio local do computador registra a data e hora atual corretamente.
  Se mudar: dia, tarefa nova, atrasada, urgente e habitos ficam errados.
- **S-03:** Uso de notebook em aula eh sempre permitido.
  Se mudar: o C3 fica sem registro na hora; anoto em outro lugar e passo depois.
- **S-04:** Uso o neivi quase todos os dias.
  Se mudar: tarefas deixam de ser novas (D-11) sem eu ter revisado. Risco aceito: o sistema mantem tudo o que ja foi feito e segue como se nada tivesse acontecido.
- **S-05:** Os dados ficam em um computador so.
  Se mudar: preciso de um servidor ou sincronizacao (Fora v0.1).
- **S-06:** A avaliacao eh registrada antes da nota chegar.
  Se mudar: ja coberto pela D-14 (a avaliacao eh criada na hora, com a nota).
- **S-07:** Sempre sei a disciplina de uma tarefa academica no momento do registro.
  Se mudar: a tarefa sem disciplina nao aparece como vinda das aulas no C7 (D-20).
- **S-08:** Uso o neivi sempre no mesmo fuso horario.
  Se mudar: dia, tarefa nova, atrasada, urgente e habitos ficam errados. Ajuste futuro.
- **S-09:** Copio o backup para um HD externo uma vez por mes (RNF-03).
  Se mudar: se o computador der pau, perco os dados desde a ultima copia.

### Dependencias

- **DEP-01:** A importacao depende do formato atual dos arquivos do Ganbatte e do TODO.txt (todo.txt, done.txt, disciplinas.txt, hunter/avaliacoes.txt e arquivos .json do hunter).
  Se mudar: a importacao precisa ser refeita. A importacao eh feita uma vez so, antes de desligar o Ganbatte; depois dela essa dependencia acaba (D-24, Q-17).
- **DEP-02:** Para rodar, o neivi precisa apenas do binario e do arquivo do banco (SQLite em Go puro, sem SQLite instalado nem bibliotecas em C). Para compilar, precisa do Go e dos modulos do go.mod.
  Se mudar: o binario se recompila; o arquivo do banco eh a unica coisa insubstituivel.
- **DEP-03:** HD externo para guardar as copias do backup.
  Se mudar: ver S-09.


# 3. User stories

## Indice

| ID    | Titulo | Prioridade | Fatia | Cenarios | Status |
|-------|--------|------------|-------|----------|--------|
| US-01 | Registro rapido de tarefa | Must | 1 | C2 | testada |
| US-02 | Registro completo de tarefa | Must | 1 | C2, C1 | testada |
| US-03 | Completar ou mudar uma tarefa | Must | 1 | C8, C2 | testada |
| US-04 | Concluir tarefas | Must | 1 | C4, C10 | testada |
| US-05 | Reabrir tarefas | Must | 1 | C4 | testada |
| US-06 | Remover tarefas | Must | 1 | C8, C2 | testada |
| US-07 | Ver o que tem para hoje | Must | 4 | C1, C4 | testada |
| US-08 | Ver as tarefas incompletas | Must | 4 | C8 | testada |
| US-09 | Ver o panorama das tarefas, com foco quando quiser | Must | 4 | C8, C7, C9, C4 | testada |
| US-10 | Consultar uma tarefa | Must | 4 | C4, C9 | testada |

Prioridade: Must (sem ela o v0.1 nao fica pronto, D-22), Should (importante, mas da para viver sem), Could (se sobrar tempo).
Status: - (nao escrita), rascunho, pronta, testada.

## Modelo

Copiar para cada story nova. Cada CA vira um teste.

```
### US-XX · Titulo curto (o que precisa existir)

Prioridade: Must · Fatia: 1 · Cenarios: C2 · Regras: RN-XX, D-XX

**Como** eu [em qual situacao: com pressa na aula, revisando o dia...],
**quero** [acao, na minha linguagem, nao na do sistema],
**para** [o porque; se nao tiver, talvez a story nao precise existir].

**Criterios de aceitacao**
- **CA-XX.1** *(caminho feliz)* Dado [estado inicial com valores concretos],
  quando [acao],
  entao [resultado que eu vejo].
- **CA-XX.2** *(limite)* Dado ..., quando ..., entao ...
- **CA-XX.3** *(entrada invalida)* Dado ..., quando ..., entao ...
- **CA-XX.4** *(estado incompativel)* Dado ..., quando ..., entao ...

**Notas** (opcional): duvidas que viraram Q-XX, o que ficou para outra fatia.
```


Checklist antes de marcar como pronta (INVEST): independente, negociavel, valiosa (o "para" eh real), estimavel, pequena, testavel (consigo imaginar o teste de cada CA).

## Fatia 1 · Tarefas


---

### US-01 · Registro rapido de tarefa

Prioridade: Must · Fatia: 1 · Cenarios: C2 · Regras: RN-01, RN-02, RN-03, RN-04

**Como** eu, no meio de outra tarefa, quando surge uma tarefa nova (uma ideia ou algo que a tarefa atual gerou),
**quero** anotar rapido so o titulo,
**para** voltar pro que estava fazendo sem perder o foco, e completar o resto (prioridade, EST) na revisao do dia (C8).

**Criterios de aceitacao**
- **CA-01.1** *(caminho feliz)* Dado que nao tenho tarefas,
  quando registro `neivi add comprar cabo HDMI`,
  entao a tarefa 1 eh criada com esse titulo, sem prioridade (incompleta, D-10), e aparece `+ 1 · comprar cabo HDMI`.
- **CA-01.2** *(entrada invalida)* Dado que nao tenho tarefas,
  quando registro `neivi add` (sem texto),
  entao aparece `erro: tarefa sem titulo` e nenhuma tarefa eh criada.
- **CA-01.3** *(entrada invalida)* Dado que nao tenho tarefas,
  quando registro `neivi add pri:A` (so marcas, sem titulo),
  entao aparece `erro: tarefa sem titulo` e nenhuma tarefa eh criada.
- **CA-01.4** *(estado incompativel)* Dado que tenho `comprar cabo HDMI` (id 1) nao concluida,
  quando registro `neivi add Comprar  cabo hdmi pri:B` (outra caixa, espaco a mais, com prioridade),
  entao a tarefa 2 eh criada e aparece `+ 2 · B · Comprar cabo hdmi` e `aviso: ja existe a tarefa 1 com esse titulo` (RN-01).
- **CA-01.5** *(limite)* Dado que tenho `comprar cabo HDMI` (id 1) concluida,
  quando registro `neivi add comprar cabo HDMI`,
  entao a tarefa 2 eh criada e aparece `+ 2 · comprar cabo HDMI`, sem aviso (RN-01: so conta tarefa nao concluida).
- **CA-01.6** *(limite)* Dado que nao tenho tarefas,
  quando registro `neivi add reuniao 10:30`,
  entao a tarefa 1 eh criada com o titulo "reuniao 10:30" e aparece `+ 1 · reuniao 10:30` (D-31).

**Notas:** CA-01.4 e CA-01.5 nao existem no codigo hoje (o add cria sem avisar), e a confirmacao hoje eh so `Tarefa: 1 adicionada`. Os erros das marcas (`pri:E`, `@XX`...) estao na US-02.

---

### US-02 · Registro completo de tarefa

Prioridade: Must · Fatia: 1 · Cenarios: C2, C1 · Regras: RN-03, RN-05, RN-06, RN-07, RN-08, RN-09

**Como** eu, organizando o dia com tempo livre, ou no meio de outra tarefa quando surge uma tarefa urgente,
**quero** registrar a tarefa numa linha so, ja com os campos que eu sei (prioridade, prazo, disciplina ou projeto; EST e descricao quando der),
**para** a tarefa urgente nao ser esquecida nem perder o prazo (ela ja aparece no resumo diario, D-27), e eu entender depois do que ela se trata.

**Criterios de aceitacao** (hoje = ter 10/11/2026; @ED e @BD cadastradas)
- **CA-02.1** *(caminho feliz)* Quando registro `neivi add entregar lista @ED pri:A due:sex est:1h -- exercicios 3 a 7`,
  entao a tarefa 1 eh criada e aparece `+ 1 · A · sex 13/11 · 1h · entregar lista @ED …` (o `…` indica descricao).
- **CA-02.2** *(caminho feliz, so alguns campos)* Quando registro `neivi add revisar slide @BD pri:B`,
  entao aparece `+ 1 · B · revisar slide @BD` (campos vazios nao aparecem).
- **CA-02.3** *(limite)* Quando registro `neivi add prova due:ter`,
  entao o prazo eh a proxima terca: `+ 1 · ter 17/11 · prova` (RN-05).
- **CA-02.4** *(limite)* Quando registro `neivi add c due:09/11` (sem ano, data que ja passou),
  entao o prazo vai para o ano seguinte e o ano aparece: `+ 1 · ter 09/11/27 · c` (RN-05).
- **CA-02.5** *(limite)* Quando registro `neivi add c due:09/11/26` (com ano, no passado),
  entao a tarefa eh criada ja atrasada: `+ 1 · seg 09/11 · c` (RN-05).
- **CA-02.6** *(entrada invalida)* Quando registro uma tarefa com uma marca invalida,
  entao nada eh criado e aparece a mensagem com o que esta errado e como resolver (RN-07):

  | Entrada                      | Mensagem |
  |------------------------------|----------|
  | `neivi add lista @XX pri:A`  | `erro: disciplina XX nao cadastrada (cadastrar: neivi subject add XX)` |
  | `neivi add lista pri:E`      | `erro: prioridade E invalida (use A, B, C ou D)` |
  | `neivi add lista due:32/11`  | `erro: data 32/11 invalida (use dd/mm, dd/mm/aa, hoje, amanha ou seg...dom)` |
  | `neivi add lista due:29/02`  | `erro: 29/02 nao existe em 2027` |
  | `neivi add lista est:abc`    | `erro: estimativa abc invalida (use 30m, 1h, 1h15 ou 90)` |
  | `neivi add lista pri:A pri:B`| `erro: mais de uma prioridade (pri:A e pri:B)` |

**Notas:** estado incompativel (tarefa duplicada) ja coberto pela RN-01 e pelo CA-01.4. Diferencas com o codigo hoje: a confirmacao (hoje so `Tarefa: 1 adicionada`) e as mensagens de erro.

---

### US-03 · Completar ou mudar uma tarefa

Prioridade: Must · Fatia: 1 · Cenarios: C8, C2 · Regras: RN-01, RN-02, RN-03, RN-07, RN-10

**Como** eu, na revisao do dia (ou, raramente, no meio do dia quando chega uma demanda urgente),
**quero** mudar so os campos que preciso de uma tarefa por vez: completar as incompletas, redefinir a prioridade ou adiar o prazo,
**para** reorganizar o dia rapido, sem redigitar a tarefa e sem perder informacao (o id e a data de criacao continuam os mesmos, entao uma tarefa antiga nao volta a parecer nova, D-11).

**Criterios de aceitacao** (hoje = ter 10/11/2026; tarefa 4 = `entregar lista @ED pri:A due:sex est:1h -- exercicios 3 a 7`)
- **CA-03.1** *(caminho feliz, adiar)* Quando registro `neivi edit 4 due:seg`,
  entao so o prazo muda e aparece `~ 4 · A · seg 16/11 · 1h · entregar lista @ED …`; o id e a data de criacao nao mudam.
- **CA-03.2** *(caminho feliz, completar)* Dado a tarefa 5 `comprar cabo HDMI` sem prioridade (incompleta),
  quando registro `neivi edit 5 pri:B due:amanha`,
  entao aparece `~ 5 · B · qua 11/11 · comprar cabo HDMI` e ela deixa de ser incompleta (RN-02).
- **CA-03.3** *(caminho feliz, titulo)* Quando registro `neivi edit 4 entregar lista 2`,
  entao so o titulo muda: `~ 4 · A · sex 13/11 · 1h · entregar lista 2 @ED …`.
- **CA-03.4** *(limite, apagar campos)* O valor `-` apaga o campo (RN-10):

  | Entrada               | Resultado |
  |-----------------------|-----------|
  | `neivi edit 4 due:-`  | `~ 4 · A · 1h · entregar lista @ED …` |
  | `neivi edit 4 pri:-`  | `~ 4 · sex 13/11 · 1h · entregar lista @ED …` (volta a ser incompleta, RN-02) |
  | `neivi edit 4 @-`     | `~ 4 · A · sex 13/11 · 1h · entregar lista …` |
  | `neivi edit 4 -- -`   | `~ 4 · A · sex 13/11 · 1h · entregar lista @ED` (sem o `…`) |

- **CA-03.5** *(entrada invalida)* Nada muda e aparece a mensagem com o que esta errado e como resolver (RN-07):

  | Entrada                  | Mensagem |
  |--------------------------|----------|
  | `neivi edit 99 pri:B`    | `erro: tarefa 99 nao existe (veja os ids no neivi ls)` |
  | `neivi edit 4`           | `erro: nada para mudar na tarefa 4 (exemplo: neivi edit 4 pri:B due:sex)` |
  | `neivi edit`             | `erro: falta o id da tarefa (exemplo: neivi edit 1 pri:B)` |
  | `neivi edit abc pri:B`   | `erro: id abc invalido (use o numero que aparece no ls)` |
  | `neivi edit 4 due:32/11` | `erro: data 32/11 invalida (use dd/mm, dd/mm/aa, hoje, amanha ou seg...dom)` |

- **CA-03.6** *(estado incompativel)* Dado a tarefa 2 `comprar cabo` concluida,
  quando registro `neivi edit 2 pri:C`,
  entao a tarefa muda e aparece `~ 2 · C · comprar cabo` e `aviso: a tarefa 2 esta concluida`.
- **CA-03.7** *(estado incompativel)* Dado a tarefa 1 `comprar cabo HDMI` nao concluida,
  quando registro `neivi edit 4 comprar cabo hdmi`,
  entao o titulo muda e aparece `aviso: ja existe a tarefa 1 com esse titulo` (RN-01).

**Notas:** editar varias tarefas de uma vez fica para depois (rascunho). O titulo nunca fica vazio (RN-03): nao existe jeito de apaga-lo, so de troca-lo. Diferencas com o codigo hoje: a confirmacao (hoje `tarefa 4 editada`), as mensagens de `edit` sem texto/sem id/id invalido (hoje mostra a ajuda), o aviso de concluida e o aviso de duplicada.

---

### US-04 · Concluir tarefas

Prioridade: Must · Fatia: 1 · Cenarios: C4, C10 · Regras: RN-11, RN-12

**Como** eu, logo depois de terminar uma tarefa, ou no fim do dia / comeco do dia seguinte marcando o que esqueci,
**quero** marcar uma ou varias tarefas como feitas, dizendo se foi ontem quando for o caso,
**para** guardar o historico do que fiz (estatistica e revisao da semana, C10) e ter a motivacao de ver quanto ja fiz no dia.

**Criterios de aceitacao** (hoje = ter 10/11/2026; tarefas 3 `entregar lista @ED pri:A`, 5 `comprar cabo HDMI pri:B`, 8 `revisar slide @BD` pendentes; a tarefa 1 ja foi concluida hoje)
- **CA-04.1** *(caminho feliz)* Quando registro `neivi done 3 5 8`,
  entao as tres ficam concluidas em 10/11 e aparece:
  `✓ 3 · A · entregar lista @ED`
  `✓ 5 · B · comprar cabo HDMI`
  `✓ 8 · revisar slide @BD`
  `hoje: 4 concluidas`
- **CA-04.2** *(caminho feliz, ontem)* Dado que nada foi concluido ontem,
  quando registro `neivi done 3 5 ontem`,
  entao as duas ficam concluidas em 09/11 e aparece:
  `✓ 3 · A · entregar lista @ED (ontem)`
  `✓ 5 · B · comprar cabo HDMI (ontem)`
  `ontem: 2 concluidas`
- **CA-04.3** *(limite)* Dado a tarefa 9 criada hoje,
  quando registro `neivi done 9 ontem`,
  entao ela fica concluida em 09/11, mesmo sendo antes da criacao (RN-11).
- **CA-04.4** *(entrada invalida, id que nao existe)* Quando registro `neivi done 3 99`,
  entao a 3 eh concluida e aparece `✓ 3 · A · entregar lista @ED`, `aviso: tarefa 99 nao existe (veja os ids no neivi ls)` e o contador (RN-12).
- **CA-04.5** *(entrada invalida, recusa tudo)* Nada eh concluido e aparece (RN-12):

  | Entrada              | Mensagem |
  |----------------------|----------|
  | `neivi done 3 otem`  | `erro: "otem" nao eh um id nem "ontem" (exemplo: neivi done 3 5 ontem)` |
  | `neivi done abc`     | `erro: "abc" nao eh um id nem "ontem" (exemplo: neivi done 3 5 ontem)` |
  | `neivi done`         | `erro: falta o id da tarefa (exemplo: neivi done 3 5)` |

- **CA-04.6** *(estado incompativel)* Quando registro `neivi done 1 3` (a 1 ja estava concluida),
  entao a 3 eh concluida, aparece `aviso: a tarefa 1 ja estava concluida` e a data de conclusao da 1 nao muda (RN-11).

**Notas:** concluir com datas mais antigas que ontem fica para depois, se for frequente (rascunho). Ver as feitas do dia separadas fica para a fatia 4. Diferencas com o codigo hoje: confirmacao com `✓` e contador, o `ontem`, as mensagens e a recusa de palavra invalida.

---

### US-05 · Reabrir tarefas

Prioridade: Must · Fatia: 1 · Cenarios: C4 · Regras: RN-12

**Como** eu, quando marquei uma tarefa como feita por engano ou percebi que faltou algo nela,
**quero** desfazer a conclusao de uma ou mais tarefas,
**para** ela voltar para as pendentes com tudo como estava, sem recriar e sem contar no meu historico algo que nao terminei.

**Criterios de aceitacao** (tarefas 5 `comprar cabo HDMI pri:B` e 8 `revisar slide @BD` concluidas; 2 `ligar` pendente)
- **CA-05.1** *(caminho feliz)* Quando registro `neivi reopen 5`,
  entao ela volta a ser pendente, sem data de conclusao, com os outros campos iguais, e aparece `↺ 5 · B · comprar cabo HDMI` (sem contador).
- **CA-05.2** *(caminho feliz, varias)* Quando registro `neivi reopen 5 8`,
  entao aparecem `↺ 5 · B · comprar cabo HDMI` e `↺ 8 · revisar slide @BD`.
- **CA-05.3** *(estado incompativel)* Quando registro `neivi reopen 2`,
  entao nada muda e aparece `aviso: a tarefa 2 nao estava concluida`.
- **CA-05.4** *(entrada invalida)* (RN-12):

  | Entrada               | Resultado |
  |-----------------------|-----------|
  | `neivi reopen 5 99`   | reabre a 5 e aparece `aviso: tarefa 99 nao existe (veja os ids no neivi ls)` |
  | `neivi reopen abc`    | nada muda; `erro: "abc" nao eh um id (exemplo: neivi reopen 5)` |
  | `neivi reopen`        | nada muda; `erro: falta o id da tarefa (exemplo: neivi reopen 5)` |

**Notas:** diferencas com o codigo hoje: confirmacao com `↺` e as mensagens.

---

### US-06 · Remover tarefas

Prioridade: Must · Fatia: 1 · Cenarios: C8, C2 · Regras: RN-01, RN-12, RN-13

**Como** eu, quando registrei algo muito errado (editar sairia mais caro que refazer), quando desisti de uma tarefa que perdeu o sentido, ou quando o aviso de duplicada mostra que eh mesmo a mesma tarefa,
**quero** apagar a tarefa de vez, confirmando uma por uma,
**para** a lista ficar so com o que eu ainda vou fazer, sem guardar o que nao fiz.

**Criterios de aceitacao** (tarefas 3 `entregar lista @ED pri:A due:sex` e 5 `comprar cabo HDMI pri:B` pendentes; 1 `ligar` concluida ontem)
- **CA-06.1** *(caminho feliz)* Quando registro `neivi rm 3`,
  entao aparece a pergunta `remover 3 · A · sex 13/11 · entregar lista @ED? [s/N]`;
  respondendo `s`, aparece `✗ 3 removida` e a tarefa some (das pendentes e do historico).
- **CA-06.2** *(limite, resposta)* So `s` remove (RN-13):

  | Resposta                          | Resultado |
  |-----------------------------------|-----------|
  | `s`                               | `✗ 3 removida` |
  | `N`, `n`, Enter ou qualquer outra | `3 mantida`, nada muda |

- **CA-06.3** *(caminho feliz, varias)* Quando registro `neivi rm 3 5`, respondo `s` para a 3 e Enter para a 5,
  entao a 3 eh removida e a 5 continua; uma pergunta por tarefa.
- **CA-06.4** *(estado incompativel)* Quando registro `neivi rm 1` (concluida),
  entao a pergunta avisa: `remover 1 · ligar (concluida, sai do historico)? [s/N]`.
- **CA-06.5** *(entrada invalida)* (RN-12):

  | Entrada          | Resultado |
  |------------------|-----------|
  | `neivi rm 3 99`  | pergunta da 3 normalmente e aparece `aviso: tarefa 99 nao existe (veja os ids no neivi ls)` |
  | `neivi rm abc`   | nada muda; `erro: "abc" nao eh um id (exemplo: neivi rm 3)` |
  | `neivi rm`       | nada muda; `erro: falta o id da tarefa (exemplo: neivi rm 3)` |

**Notas:** sem estado "cancelada" (descartado: nao quero guardar o que nao fiz). Diferencas com o codigo hoje: a pergunta com a linha curta e o aviso de concluida, `✗ removida`/`mantida` e as mensagens.

---

## Fatia 4 · Listagem e resumo diario

---

### US-07 · Ver o que tem para hoje

Prioridade: Must · Fatia: 4 · Cenarios: C1, C4 · Regras: RN-14, RN-16, RN-17

**Como** eu, ao ligar o PC (C1) e ao longo do dia, depois de concluir uma tarefa (C4),
**quero** ver de uma vez o que tem para hoje (atrasadas, urgentes e importantes), com o tempo estimado do que falta,
**para** me preparar para o dia em ~1 min e escolher a proxima tarefa sem montar filtro.

**Criterios de aceitacao** (hoje = qua 11/11/2026; 3 tarefas ja concluidas hoje)
- **CA-07.1** *(caminho feliz)* Dado as pendentes:
  5 `lista 3 @ED pri:B due:09/11/26 est:30m`, 3 `entregar parte do feedback360 @BD pri:A due:hoje est:1h`,
  7 `pagar boleto pri:C due:hoje est:30m`, 2 `revisar normalizacao @BD pri:A est:2h`, 8 `ligar para o dentista pri:A`,
  9 `estudar sorting pri:B` e 10 `comprar cabo pri:A due:sex`,
  quando registro `neivi today`, entao aparece:

      atrasadas (1)
      ID  P  PRAZO      EST  TAREFA
      5   B  seg 09/11  30m  lista 3 @ED

      urgentes (2)
      ID  P  PRAZO      EST  TAREFA
      3   A  qua 11/11  1h   entregar parte do feedback360 @BD
      7   C  qua 11/11  30m  pagar boleto

      importantes (2)
      ID  P  PRAZO      EST  TAREFA
      2   A             2h   revisar normalizacao @BD
      8   A                  ligar para o dentista

      falta: ~4h (1 sem EST)
      hoje: 3 concluidas

  (cada grupo eh uma tabela com o proprio cabecalho; as colunas tem a mesma largura em todos os grupos)

  (a 9 nao aparece porque ha A sem prazo; a 10 so aparece no dia dela ou quando atrasar; concluidas nunca aparecem)
- **CA-07.2** *(limite, sem A sem prazo)* Dado que nao ha nenhuma A sem prazo e ha as B sem prazo 9 `estudar sorting` e 12 `levantamento de notebook`,
  entao o bloco 3 vira `sem prazo (2)` com as duas (RN-16).
- **CA-07.3** *(limite, incompleta com prazo hoje)* Dado 13 `pagar aluguel due:hoje` sem prioridade,
  entao ela aparece como uma linha em urgentes, sem prioridade (RN-16).
- **CA-07.4** *(limite, incompleta sem prazo)* Dado 14 `ideia qualquer`, sem prioridade e sem prazo,
  entao ela nao aparece em nenhum bloco (RN-16).
- **CA-07.5** *(limite, bloco vazio)* Dado que nao ha atrasadas,
  entao o bloco `atrasadas` nao aparece.
- **CA-07.6** *(estado, dia livre)* Dado que nada se encaixa nos tres blocos e nao ha incompletas,
  entao aparece `nada para hoje, bom trabalho!` e `hoje: 3 concluidas`.
- **CA-07.7** *(estado, dia livre com incompletas)* Dado que nada se encaixa nos tres blocos e ha 4 incompletas,
  entao aparece `nada para hoje, bom trabalho!`, `voce tem 4 tarefas incompletas: que tal completar? (neivi triage)` e `hoje: 3 concluidas`.
- **CA-07.8** *(limite, EST)* Dado que nenhuma tarefa mostrada tem EST,
  entao aparece `falta: sem EST (3 tarefas)` no lugar da soma.
- **CA-07.9** *(entrada invalida)* Quando registro `neivi today xyz`,
  entao aparece `erro: today nao recebe argumentos (exemplo: neivi today)`.

**Notas:** o contador aparece sempre, mesmo `hoje: 0 concluidas`. Os habitos do dia entram abaixo na fatia 3.

---

### US-08 · Ver as tarefas incompletas

Prioridade: Must · Fatia: 4 · Cenarios: C8 · Regras: RN-02, RN-17

**Como** eu, na revisao do dia (C8), ou quando o today sugere num dia livre,
**quero** ver so as tarefas incompletas, das mais antigas para as mais novas,
**para** decidir a prioridade delas e parar de adiar essa escolha.

**Criterios de aceitacao** (hoje = qua 11/11/2026)
- **CA-08.1** *(caminho feliz)* Dado as incompletas 4 `comprar cabo HDMI` (criada em 05/11), 9 `ligar` (10/11) e 12 `pagar aluguel due:sex` (hoje), e outras tarefas com prioridade,
  quando registro `neivi triage`, entao aparece:

      incompletas (3)
      ID  P  PRAZO      EST  TAREFA
      4                      comprar cabo HDMI
      9                      ligar
      12     sex 13/11       pagar aluguel

- **CA-08.2** *(caminho feliz, alternativa)* `neivi ls pri:-` mostra as mesmas tarefas (story do ls).
- **CA-08.3** *(limite)* Dado a incompleta 4 concluida, entao ela nao aparece.
- **CA-08.4** *(estado)* Dado que nao ha incompletas, entao aparece `nenhuma tarefa incompleta`.
- **CA-08.5** *(entrada invalida)* Quando registro `neivi triage xyz`,
  entao aparece `erro: triage nao recebe argumentos (exemplo: neivi triage)`.

**Notas:** sem mostrar ha quantos dias cada uma esta esperando; se fizer falta, entra depois.

---

### US-09 · Ver o panorama das tarefas, com foco quando quiser

Prioridade: Must · Fatia: 4 · Cenarios: C8, C7, C9, C4 · Regras: RN-14, RN-17, RN-18

**Como** eu, quando chega uma demanda nova e preciso reavaliar tudo, ou quando quero focar so nas aulas (C7) ou num projeto (C9),
**quero** ver todas as tarefas pendentes numa tabela, e filtrar por disciplina, projeto, prioridade, tarefas novas ou palavras do titulo,
**para** redefinir prioridades olhando o todo, ou organizar uma parte sem ver o resto.

**Criterios de aceitacao** (hoje = qua 11/11/2026; @BD e @ED cadastradas; +neivi cadastrado)
Pendentes: 3 `entregar parte do feedback360 @BD pri:A due:hoje est:1h`, 2 `revisar normalizacao @BD pri:A est:2h -- slides 4 a 9`,
11 `parser +neivi pri:B due:sex`, 9 `estudar sorting @ED pri:B` (criada em 05/11) e 4 `comprar cabo HDMI` (incompleta).
Concluida hoje: 1 `ligar`. Todas, menos a 9, criadas hoje.

- **CA-09.1** *(caminho feliz)* Quando registro `neivi ls`, entao aparecem as pendentes na ordem padrao (RN-17):

      ID  P  PRAZO      EST  TAREFA
      3   A  qua 11/11  1h   entregar parte do feedback360 @BD
      2   A             2h   revisar normalizacao @BD …
      11  B  sex 13/11       parser +neivi
      9   B                  estudar sorting @ED
      4                      comprar cabo HDMI

- **CA-09.2** *(caminho feliz, foco)* Os filtros mostram so uma parte (RN-18):

  | Entrada          | Tarefas        |
  |------------------|----------------|
  | `neivi ls @BD`   | 3, 2           |
  | `neivi ls @`     | 3, 2, 9        |
  | `neivi ls @-`    | 11, 4          |
  | `neivi ls +neivi`| 11             |
  | `neivi ls +`     | 11             |
  | `neivi ls +-`    | 3, 2, 9, 4     |
  | `neivi ls pri:A` | 3, 2           |
  | `neivi ls pri:-` | 4 (igual ao `neivi triage`) |
  | `neivi ls new`   | 3, 2, 11, 4 (criadas hoje ou ontem, D-11) |

- **CA-09.3** *(caminho feliz, combinar)* Mesmo tipo eh "ou", tipos diferentes sao "e" (RN-18):
  `neivi ls @BD @ED` mostra 3, 2, 9; `neivi ls @BD pri:A` mostra 3, 2; `neivi ls new @BD` mostra 3, 2.
- **CA-09.4** *(caminho feliz, busca)* Palavras buscam no titulo; todas precisam aparecer; maiusculas e acentos nao importam (RN-18):
  `neivi ls Normalização` mostra 2; `neivi ls cabo hdmi` mostra 4.
- **CA-09.5** *(caminho feliz, concluidas)* `neivi ls all` mostra as pendentes na ordem padrao e, no fim, as concluidas das mais recentes para as mais antigas, com a coluna `X`:

      ID  X  P  PRAZO      EST  TAREFA
      3      A  qua 11/11  1h   entregar parte do feedback360 @BD
      ...
      4                         comprar cabo HDMI
      1   x                     ligar

  `neivi ls done` mostra so as concluidas, das mais recentes para as mais antigas, sem a coluna `X`.
- **CA-09.6** *(caminho feliz, descricao)* `neivi ls desc` mostra a descricao embaixo da tarefa, recuada, no lugar do `…`.
- **CA-09.7** *(caminho feliz, prazo)* `neivi ls due` ordena por prazo -> prioridade -> EST -> id (a ordem do resumo diario, RN-17): 3, 11, 2, 9, 4.
- **CA-09.8** *(estado, nada encontrado)* Sem pendentes, `neivi ls` mostra `nenhuma tarefa pendente`.
  Com filtro ou busca sem resultado, `neivi ls cabo azul` mostra `nenhuma tarefa com esses filtros (tente menos palavras ou filtros)`.
- **CA-09.9** *(entrada invalida)* `neivi ls @XX` mostra `erro: disciplina XX nao cadastrada (cadastrar: neivi subject add XX)`; o mesmo para `+XX`.

**Notas:** a revisao semanal (C10) fica fora (decidir depois do corte das tarefas). `all`, `done`, `late`, `new`, `due` e `desc` sao reservadas e nao servem para busca.

---

### US-10 · Consultar uma tarefa

Prioridade: Must · Fatia: 4 · Cenarios: C4, C9 · Regras: RN-02, RN-14

**Como** eu, antes de comecar uma tarefa ou quando preciso lembrar do que ela se trata,
**quero** ver todos os campos de uma tarefa, inclusive a descricao inteira,
**para** entender a tarefa sem depender so do titulo.

**Criterios de aceitacao** (hoje = qua 11/11/2026; @BD = "Projeto de Banco de Dados")
- **CA-10.1** *(caminho feliz)* Dado a tarefa 2 `revisar normalizacao @BD pri:A due:sex est:2h -- slides 4 a 9 da aula de segunda`, criada ontem,
  quando registro `neivi show 2`, entao aparece:

      tarefa 2 · revisar normalizacao
        prioridade:  A
        prazo:       sex 13/11 (em 2 dias)
        EST:         2h
        disciplina:  BD (Projeto de Banco de Dados)
        projeto:     -
        criada:      ter 10/11 (nova)
        situacao:    pendente

        slides 4 a 9 da aula de segunda

- **CA-10.2** *(limite, campos vazios)* Campos sem valor aparecem como `-`; sem descricao, nao aparece a linha em branco nem o texto.
- **CA-10.3** *(limite, prazo deduzido)* O parenteses do prazo diz `hoje`, `amanha`, `em N dias` ou `atrasada ha N dias`.
- **CA-10.4** *(limite, situacao)* `situacao` eh `pendente`, `incompleta` (sem prioridade, RN-02) ou `concluida em qua 11/11`; tarefa concluida nao tem o parenteses do prazo.
- **CA-10.5** *(limite, disciplina herdada)* Dado uma tarefa sem `@` de um projeto de BD, `disciplina` mostra `BD (do projeto)` (RN-14).
- **CA-10.6** *(entrada invalida)* (mesmas mensagens do edit):

  | Entrada            | Mensagem |
  |--------------------|----------|
  | `neivi show 99`    | `erro: tarefa 99 nao existe (veja os ids no neivi ls)` |
  | `neivi show abc`   | `erro: id abc invalido (use o numero que aparece no ls)` |
  | `neivi show`       | `erro: falta o id da tarefa (exemplo: neivi show 2)` |
  | `neivi show 2 3`   | `erro: show mostra uma tarefa por vez (exemplo: neivi show 2)` |

**Notas:** o formato pode mudar com o uso.

---

# 4. Regras de negocio

Regras que valem para qualquer comando ou tela. Cada uma diz de onde veio (Fonte) e onde eh usada.

- **RN-01** Tarefa duplicada: tarefa nao concluida com o mesmo titulo, ignorando maiusculas/minusculas e espacos extras. Marcas (prioridade, prazo, disciplina, projeto...) nao contam. Duplicada gera aviso, nunca erro; tarefa concluida com o mesmo titulo nao gera aviso.
  Fonte: CA-01.4, CA-01.5 · Usada em: US-01, US-02, US-03, US-06
- **RN-02** Tarefa incompleta: tarefa sem prioridade definida. Eh deduzida, sem marcacao.
  Fonte: D-10, D-17 · Usada em: US-01, US-03
- **RN-03** Titulo obrigatorio: toda tarefa tem titulo. O titulo eh o texto que sobra depois de tirar as marcas; so marcas (ex.: `pri:A`) nao formam titulo.
  Fonte: CA-01.2, CA-01.3 · Usada em: US-01, US-02, US-03
- **RN-04** Palavra com `:` que nao eh uma marca conhecida faz parte do titulo (ex.: `reuniao 10:30`, `aula:3`).
  Fonte: D-31, CA-01.6 · Usada em: US-01
- **RN-05** Prazo: aceita `dd/mm`, `dd/mm/aa` ou `dd/mm/aaaa` (com `/` ou `-`), `hoje`, `amanha` e `seg`...`dom`. Sem ano vale a proxima ocorrencia (hoje conta como hoje; data que ja passou vai para o ano seguinte). Dia da semana eh sempre o proximo (`ter` numa terca = terca seguinte). `29/02` sem ano so vale em ano bissexto. Prazo no passado so com o ano escrito.
  Fonte: SPEC 4, CA-02.3 a CA-02.5 · Usada em: US-02
- **RN-06** Disciplina e projeto da tarefa precisam estar cadastrados.
  Fonte: SPEC 4 · Usada em: US-02
- **RN-07** Uma marca invalida recusa a operacao inteira (registro ou edicao); nada eh criado ou mudado pela metade.
  Fonte: CA-02.6, CA-03.5 · Usada em: US-02, US-03
- **RN-08** Prioridade: A, B, C ou D, aceita minuscula; no maximo uma por tarefa.
  Fonte: D-10, glossario · Usada em: US-02
- **RN-09** Estimativa: minutos ou horas (`30m`, `1h`, `1h15`, `90` = 90 minutos), sempre maior que zero.
  Fonte: SPEC 4 · Usada em: US-02
- **RN-10** Edicao parcial: so muda o que aparece no texto; o resto fica como esta. O valor `-` apaga o campo (`pri:-`, `due:-`, `est:-`, `@-`, `+-`, `-- -`). Texto solto vira o titulo novo. O id e a data de criacao nunca mudam.
  Fonte: SPEC 4, D-11, CA-03.1 a CA-03.4 · Usada em: US-03
- **RN-11** Data de conclusao: hoje, ou ontem quando escrevo `ontem`; nenhuma outra data no v0.1. Pode ser anterior a data de criacao. Concluir de novo uma tarefa concluida nao muda a data.
  Fonte: D-21 (mesma ideia dos habitos), CA-04.2, CA-04.3, CA-04.6 · Usada em: US-04
- **RN-12** Comandos com varios ids: id que nao existe gera aviso e os outros seguem; palavra que nao eh id (nem `ontem`, no done) recusa tudo, porque nao da para saber o que eu quis dizer.
  Fonte: CA-04.4, CA-04.5, CA-05.4, CA-06.5 · Usada em: US-04, US-05, US-06
- **RN-13** Remover apaga de vez: a tarefa sai das pendentes e do historico, sem estado de cancelada. Sempre pergunta, uma por tarefa, e so `s` confirma. Se a tarefa esta concluida, a pergunta avisa que ela sai do historico.
  Fonte: CA-06.1 a CA-06.4 · Usada em: US-06
- **RN-14** Disciplina efetiva: a disciplina de uma tarefa eh o `@` dela, se tiver; senao, a do projeto. Eh deduzida, nunca copiada (D-17): trocar a disciplina do projeto muda a de todas as tarefas dele sem `@`. Filtro por disciplina usa a disciplina efetiva, entao tarefa com `@` diferente do projeto nunca aparece na disciplina do projeto.
  Fonte: D-15, D-17, secao 7 · Usada em: US-02 (e nos filtros da fatia 4)
- **RN-15** Criterio de prioridade. A prioridade mede a importancia; a urgencia eh o prazo (`due:`).

  | Nivel | Definicao | Pergunta-teste | Exemplos |
  |-------|-----------|----------------|----------|
  | **A** | Consequencia seria ou com outra pessoa | Se eu atrasar, perco nota ou dinheiro, ou alguem alem de mim eh prejudicado? | entregar parte do feedback360, pagar contas, preparar aula para aluno, trabalho em grupo |
  | **B** | Consequencia real, mas recuperavel, que so me afeta | Se eu atrasar, eu sofro, mas consigo recuperar? | estudar o conteudo da aula que faltei, neivi, levantamento de notebook para o meu pai |
  | **C** | Seria bom, mas sem consequencia seria se atrasar | Ninguem sofre se atrasar, mas quero fazer? | projetos pessoais (costumam ser C), configurar o Hyprland |
  | **D** | Talvez um dia | Se ficar 3 meses parado, tudo bem? | testar um emulador, curiosidades |
  | vazia | Ainda nao decidida (incompleta, RN-02) | Ainda nao pensei nela? | o que entrou correndo pela US-01 |

  Fonte: Q-14 · Usada em: US-02, US-03

- **RN-16** Resumo diario: tres blocos de tarefas pendentes, nesta ordem: atrasadas (prazo antes de hoje), urgentes (prazo hoje, qualquer prioridade, inclusive sem prioridade) e importantes (prioridade A sem prazo). Se nao ha nenhuma A sem prazo, o terceiro bloco vira "sem prazo" e mostra todas do nivel mais alto que existir entre as sem prazo (B, senao C, senao D). Incompleta sem prazo nunca aparece. Bloco vazio nao aparece; se os tres estao vazios, "nada para hoje" (com a dica do triage se houver incompletas).
  Fonte: D-27, CA-07.1 a CA-07.7 · Usada em: US-07
- **RN-17** Ordenacao: no resumo diario, prazo (mais antigo primeiro) -> prioridade (A a D, sem prioridade por ultimo) -> EST (mais rapidas primeiro, sem EST por ultimo) -> id. Na listagem padrao, prioridade -> prazo -> EST -> id. No triage, data de criacao (mais antiga primeiro) -> id.
  Fonte: D-04, D-28, US-08 · Usada em: US-07, US-08, US-09
- **RN-18** Filtros da listagem: `@codigo`, `+codigo`, `pri:X` escolhem; `@` e `+` sozinhos = qualquer disciplina / qualquer projeto; `@-`, `+-`, `pri:-` = sem. Disciplina usa a disciplina efetiva (RN-14). `new` = criadas hoje ou ontem (D-11); `late` = atrasadas. Filtros do mesmo tipo sao "ou"; de tipos diferentes, "e". Outras palavras buscam no titulo: todas precisam aparecer, ignorando maiusculas e acentos. Palavras reservadas (`all`, `done`, `late`, `new`, `due`, `desc`) nao entram na busca.
  Fonte: decisoes de 04/11 (SPEC 5), D-11, US-09 · Usada em: US-08, US-09

# 5. Requisitos nao funcionais

- **RNF-01** Desempenho: todo comando de tarefa (`add`, `edit`, `done`, `reopen`, `rm`) e de consulta (`today`, `ls`, `triage`, `show`) responde em menos de 1 segundo com 10 000 tarefas no banco (5x o uso esperado de um ano, ~2000). As consultas leem so as pendentes quando der.
  Como verificar: `time` num banco de teste gerado com 10 000 tarefas. · Fonte: C2
- **RNF-02** Confiabilidade (copia local): uma copia diaria consistente do banco (nunca `cp` do arquivo aberto), conferida com `integrity_check`, guardando as 7 ultimas diarias e 4 semanais. Perda maxima aceita: 1 dia.
  Como verificar: o timer roda todo dia; uma vez por mes, restaurar uma copia num banco separado e conferir que o neivi abre e lista as tarefas. · Fonte: DEP-02, D-29
- **RNF-03** Confiabilidade (fora da maquina): ate existir o servidor, o dump vai todo dia para o repositorio privado `todo` no GitHub (excecao a R-03, D-33) e, uma vez por mes, para o HD externo (S-09). Com o servidor, ele vira a copia diaria e o GitHub fica espacado.
  Como verificar: o commit do dia aparece no repositorio; a copia do mes esta no HD. · Fonte: D-33, S-09
- **RNF-04** Usabilidade (mensagens): toda mensagem de erro ou aviso diz o que esta errado e como resolver, sem detalhes internos do programa (ex.: nada de `strconv.ParseInt`).
  Como verificar: os CAs de entrada invalida das stories. · Fonte: US-02 a US-06
- **RNF-05** Usabilidade (confirmacao): depois de criar, editar, concluir, reabrir ou remover, o neivi mostra a tarefa numa linha curta: simbolo (`+` criar, `~` editar, `✓` concluir, `↺` reabrir, `✗` remover), id, prioridade, prazo (`dia dd/mm`, com o ano so quando nao eh o atual), EST, titulo, `@disciplina`, `+projeto` e `…` se tiver descricao. Campos vazios nao aparecem. O `done` termina com o contador do dia.
  A linha curta eh so para confirmacoes (uma tarefa). Listas (`today`, `ls`, `triage`) usam tabela: `ID P PRAZO EST TAREFA` (o `X` so no `ls all`), prazo como `qua 11/11`, colunas alinhadas.
  Como verificar: os CAs de caminho feliz das stories. · Fonte: US-01 a US-06

# 6. Casos de uso

Nenhum na fatia 1 (as stories cobrem os fluxos). O primeiro sera a importacao do Ganbatte (fatia 5).

# 7. Modelo conceitual

Entidades e relacionamentos da fatia 1 (sem tabelas nem colunas; isso eh o SPEC 8).

```
    Disciplina 1 ─────── 0..N Tarefa
        │ 1                   │ 0..N
        │                     │
        │ 0..N                │
        └──────── Projeto 1 ──┘
```

- **Tarefa:** titulo, descricao, prioridade, prazo, EST, data de criacao, data de conclusao.
- **Disciplina:** codigo (o que vai depois do `@`), nome.
- **Projeto:** codigo (o que vai depois do `+`), nome. Os outros atributos (D-06) entram na fatia 2.

Relacionamentos:
- Tarefa **eh de** no maximo uma Disciplina (opcional); uma Disciplina tem varias Tarefas.
- Tarefa **pertence a** no maximo um Projeto (opcional); um Projeto tem varias Tarefas.
- Projeto **eh de** no maximo uma Disciplina (opcional, D-06); uma Disciplina tem varios Projetos.
- Tarefa pode ter Disciplina e Projeto ao mesmo tempo, inclusive com uma Disciplina diferente da do Projeto (RN-14).

Deduzido, nunca guardado (D-17):
- tarefa incompleta (RN-02) e tarefa nova (D-11);
- disciplina efetiva da tarefa (RN-14);
- tipo do projeto, academico ou pessoal (D-06);
- atrasada, urgente e importante (D-27, fatia 4).

# 8. Decisoes e questoes em aberto

### Questoes

- **Q-01** Urgente: o que torna uma tarefa urgente? -> resolvida em D-09 (glossario)
- **Q-02** Atrasada: prazo hoje ja conta? -> resolvida no glossario (nao conta)
- **Q-03** Tarefa pronta: quais campos? -> resolvida em D-10
- **Q-04** Tarefa nova: desde quando? -> resolvida em D-11
- **Q-05** Periodicidade dos habitos no v0.1 -> resolvida em D-08
- **Q-06** Habito: passar da meta conta? -> resolvida em D-12
- **Q-07** Desempate por EST -> resolvida em D-04
- **Q-08** Projeto: situacao e tipo -> resolvida em D-06
- **Q-09** Conteudo de avaliacao nao definido pelo professor -> resolvida em D-13
- **Q-10** Ordem do C7 e C8 -> confirmada: C7 (organizar aulas) antes do C8 (revisao do dia)
- **Q-11** Nota sem avaliacao registrada -> resolvida em D-14
- **Q-12** Habito x tarefa recorrente -> resolvida no glossario (Habito)
- **Q-13** Tarefa recorrente no v0.1? -> resolvida em D-16
- **Q-14** Criterios de cada nivel de prioridade (A a D), para nem tudo ficar no A nem tudo vazio. -> resolvida na RN-15
- **Q-15** Adicionar ou concluir tarefa em disciplina ou projeto inativo: o sistema bloqueia, avisa, ou eu so suponho que nao vou fazer? (ex-S-08) *(em aberto)*
- **Q-16** Disciplina precisa de situacao ou semestre, para saber quais sao as atuais? Hoje so o projeto tem situacao. (ex-S-08) *(em aberto)*
- **Q-17** Quais arquivos do Ganbatte importar: todo.txt, done.txt, disciplinas.txt, hunter/avaliacoes.txt, os .json do hunter? (DEP-01) -> parcial: no v0.1 so tarefas e disciplinas; o hunter fica para o v0.2 (D-23); o done.txt eh a Q-18
- **Q-18** Importar o done.txt (historico antigo de tarefas concluidas)? *(em aberto, fatia de importacao)*
- **Q-19** Para onde vai a copia do backup: repositorio privado, HD externo (S-09) ou os dois? -> resolvida em D-33

### Decisoes

- **D-01** Maratona nao tem estrutura propria: eh um projeto continuo (tarefas = conteudos de estudo) + um habito de questoes no Codeforces.
  Porque: evita criar estrutura para um caso so; mantem o sistema nao complexo (glossario: Complexo; P6).
- **D-02** Todo habito tem uma meta numerica; feito/nao feito eh o caso de meta 1.
  Porque: uma regra so para os dois tipos de habito.
- **D-03** Questoes do Codeforces sao registradas manualmente no v0.1; contagem automatica (API do Codeforces) fica para o futuro.
  Porque: simples, sem dependencia externa; o cenario C11 funciona igual.
- **D-04** Duas ordenacoes. Resumo diario (C1): prazo -> prioridade -> EST. Padrao e planejamento (C8): prioridade -> prazo -> EST. No EST, as mais rapidas primeiro. Projetos nao tem ordem sequencial propria.
  Porque: o resumo eh sobre urgencia, o planejamento eh sobre importancia; a prioridade ja define o proximo conteudo de um projeto.
- **D-05** Habitos nao se ligam a projetos. O C11 eh atendido por uma visao propria que junta o habito e as tarefas do projeto.
  Porque: a maratona eh o unico caso; um caso unico nao justifica mudar a estrutura.
- **D-06** Projeto guarda: nome, descricao, link do github, data de criacao (registrada pelo sistema), situacao (ativo, pausado, concluido, cancelado), disciplina (opcional) e caminho da pasta local (opcional, ex.: ~/Code/...). O tipo pessoal/academico eh deduzido: com disciplina eh academico. Disciplina nao aponta para projeto.
  Porque: minimo para o C9; tipo deduzido evita um campo a mais; o caminho local me leva direto para a pasta do projeto.
- **D-07** Link do github no projeto sobe para o v0.1 (so hyperlink).
  Porque: o C9 precisa dele e o custo eh baixo.
- **D-08** Periodicidade dos habitos: diario, N vezes por semana (dias livres) e dias especificos (seg/qua/sex). Os tres entram direto no v0.1.
  Porque: quero os habitos por tipo e meta desde o inicio.
- **D-09** Urgente = tarefa de maior prioridade com prazo no dia atual. *(revista na D-27)*
  Porque: o resumo diario (C1) mostra so o que nao pode esperar.
- **D-10** Tarefa pronta = prioridade definida (A a D). Prioridade vazia = tarefa incompleta.
  Porque: separa "ainda nao decidi" de "decidi que eh baixa"; disciplina e projeto nao servem, pois tarefas pessoais nao tem nenhum dos dois.
- **D-11** Tarefa nova = criada hoje ou ontem, deduzido pela data de criacao, sem marcacao.
  Porque: atende o C7 e o C8 e deixa de ser nova sozinha, sem passo a mais.
- **D-12** Habito pode ser registrado alem da meta; o valor real fica no historico. Ao bater a meta, o habito sai da visao do dia; registros extras sao feitos por outro caminho.
  Porque: quero saber quanto fiz de verdade, sem poluir a visao do dia.
- **D-13** Conteudo de avaliacao nao definido vira uma tarefa ("perguntar o conteudo da avaliacao de X").
  Porque: eh raro e uma tarefa ja resolve, sem funcionalidade nova.
- **D-14** No v0.1, notas guardam so o valor obtido em cada avaliacao (sem pesos e media); o plano de avaliacao (pesos, N1/N2, media) fica para o futuro. Se a nota chegar e a avaliacao nao existir, a avaliacao eh criada na hora, ja com a nota.
  Porque: o esquema de pesos da PUC varia por professor e muda no semestre; merece uma fatia propria. Registrar a nota ja atende o C6.
- **D-15** O filtro por disciplina inclui as tarefas dos projetos ligados a ela.
  Porque: as tarefas de um projeto de uma disciplina sao dela; evita repetir a disciplina em cada tarefa. *(refinada na RN-14: so as tarefas sem @ proprio)*
- **D-16** Sem tarefas recorrentes no v0.1.
  Porque: nao preciso agora.
- **D-17** Principio: deduzir em vez de guardar (tarefa nova, tarefa incompleta, tipo do projeto).
  Porque: menos passos manuais e menos estado para manter; segue o "nao complexo".
- **D-18** Dia = das 00:00 as 23:59 (glossario). (ex-S-04)
  Porque: *(preencher)*
- **D-19** Semana = de domingo a sabado (glossario). (ex-S-05)
  Porque: *(preencher)*
- **D-20** Tarefa vinda das aulas = tarefa nova (D-11) com disciplina. (ex-S-07)
  Porque: deduzido, sem marcacao a mais (D-17); atende o C7.
- **D-21** Habitos podem ser registrados em hoje ou em ontem.
  Porque: as vezes esqueco de marcar ou nao tenho acesso ao computador quando faco.
- **D-22** Pronto do v0.1: uso o neivi todo dia no lugar do TODO.txt para tarefas e habitos, com o resumo diario, e com backup funcionando.
  Porque: sem uma frase de pronto, toda funcionalidade boa entra no v0.1 e ele nunca termina.
- **D-23** Avaliacoes e notas saem do v0.1 e vao para o v0.2, prontas antes do semestre 2027-1. A D-14 continua valendo.
  Porque: o semestre 2026-2 acaba em dezembro; quando o v0.1 ficar pronto quase nao vai ter prova para registrar. Eh a maior entidade nova do escopo.
- **D-24** Sem exportacao para o TODO.txt. O Ganbatte tem corte seco: importo uma vez quando o v0.1 estiver pronto e desligo. A automacao das 5:50 continua so como gatilho do dia ate a TUI.
  Porque: o export servia ao Ganbatte e a automacao, nao ao backup; o backup eh o dump do proprio banco. Estou no fim do semestre e nao vou usar as sugestoes de aula do Ganbatte.
- **D-25** Tarefas filhas e checklist fora do v0.1, e fora do banco.
  Porque: nenhum cenario (C1 a C11) pede; se fizerem falta no uso, voltam com uma migracao.
- **D-26** Sem sequencia (streak) e sem dias de folga no v0.1. A folga, quando vier, eh uma regra de calculo sobre os registros: os dias que faltaram sao perdoados na conta, sem mudar os dados, inclusive no passado.
  Porque: a sequencia so faz sentido com a visualizacao da TUI; guardando os registros de cada dia, ela sai deles depois (D-17).
- **D-27** Resumo diario em tres blocos, nesta ordem: atrasadas, urgentes (prazo hoje, qualquer prioridade) e importantes (prioridade A sem prazo). Tarefa A com prazo futuro so aparece no dia dela ou quando atrasar. Substitui a D-09.
  Porque: uma tarefa B com prazo hoje tambem nao pode esperar; e o resumo continua curto, de ~1 min (C1). O resumo eh o C1 (primeira vez do dia), diferente da revisao do C8 (primeiro bloco livre). *(completada na RN-16)*
- **D-28** Ordenacao: a da D-04, com o id como ultimo desempate.
  Porque: a ordem nunca muda entre duas listagens iguais.
- **D-29** Backup diario automatico entra no v0.1, como RNF (secao 5). A S-09 complementa, nao substitui.
  Porque: o banco eh a unica coisa insubstituivel (DEP-02).
- **D-30** C9 atendido por um comando de panorama do projeto; C11 atendido pela listagem do projeto da maratona junto com a lista de habitos, sem visao propria no v0.1.
  Porque: na CLI dois comandos ja resolvem; a visao que junta tudo fica para a TUI.
- **D-31** Sem campo para as chaves do Ganbatte (`aula:`, `av:`, `top:`, `gh:`...). Na importacao elas vao para a descricao da tarefa ("Ganbatte: aula:2026-09-30 av:A12"). Fora da importacao, palavra com `:` desconhecida faz parte do titulo.
  Porque: o campo so existia para elas voltarem iguais no export (D-24); na descricao a informacao fica a vista sem mudar o banco, e o TODO.txt original continua no historico do git.
- **D-32** O requisitos.md diz o que o neivi faz (escopo, cenarios, regras, stories); o SPEC.md diz como (sintaxe, arquitetura, modelo de dados, comandos).
  Porque: os dois falando de escopo iam se contradizer.
- **D-33** Excecao a R-03: o backup vai para o repositorio privado `todo` no GitHub, sem criptografia, ate existir o servidor em casa; depois, so de forma espacada. Q-19 resolvida.
  Porque: ate la eh o unico jeito automatico de tirar os dados de casa; o HD externo eh manual (uma vez por mes). Ressalva: no v0.2 as notas tambem iriam para o GitHub; se incomodar, criptografar com `age`.

# 9. Rastreabilidade

Cada criterio aponta para o teste que o confere (`go test ./...`). Os nomes dos casos comecam pelo id do criterio.

| Story | Criterios | Testes |
|-------|-----------|--------|
| US-01 | CA-01.1 a CA-01.6 | `cmd/neivi`: TestAdd (CA-01.1 a 01.6); `internal/task`: TestNormalizeTitle (CA-01.4); `internal/store`: TestFindDuplicate (CA-01.4, 01.5) |
| US-02 | CA-02.1, 02.2 | `cmd/neivi`: TestFormatLine (02.1, 02.2), TestAdd (02.2) |
| US-02 | CA-02.3 a 02.5 | `cmd/neivi`: TestFormatDue; `internal/todotxt`: TestDueDate (datas relativas a um "hoje" fixo) |
| US-02 | CA-02.6 | `cmd/neivi`: TestAdd (cada linha da tabela, menos o 29/02); `internal/todotxt`: TestDueDate (29/02); `internal/store`: TestSubjectIDByCode |
| US-03 | CA-03.1, 03.2 | `cmd/neivi`: TestFormatLine (03.1), TestFormatDue (03.1, 03.2) |
| US-03 | CA-03.3 a 03.7 | `cmd/neivi`: TestEdit; `internal/store`: TestFindDuplicate (03.7) |
| US-04 | CA-04.1 a 04.6 | `cmd/neivi`: TestDone, TestParseIDs; `internal/store`: TestDoneTaskECountDone (04.1 a 04.3) |
| US-05 | CA-05.1 a 05.4 | `cmd/neivi`: TestReopen, TestParseIDs; `internal/store`: TestReopenTask |
| US-06 | CA-06.1 a 06.5 | `cmd/neivi`: TestRm; `internal/store`: TestRmTask |

| US-07 | CA-07.1 a CA-07.9 | `cmd/neivi`: TestToday (07.1, 07.2, 07.3, 07.6 a 07.9); `internal/task`: TestToday (07.1 a 07.6), TestIsNewIsLate |
| US-08 | CA-08.1 a CA-08.5 | `cmd/neivi`: TestTriage; `internal/task`: TestSorts (08.1); `internal/store`: TestListTasksFilter (08.2) |
| US-09 | CA-09.1 a CA-09.9 | `cmd/neivi`: TestLs; `internal/store`: TestListTasksFilter (09.1 a 09.5); `internal/task`: TestSorts (09.1, 09.5, 09.7), TestMatchWords (09.4), TestFilterIsEmpty (09.8); `internal/todotxt`: TestParseFilter |
| US-10 | CA-10.1 a CA-10.6 | `cmd/neivi`: TestShow, TestDaysLabel (10.1, 10.3), TestCodeName (10.1) |

**Pendente da fatia 2:** CA-10.5 e a parte "herdada" do `ls @` (RN-14) dependem de o projeto ter disciplina (D-06). Ate la, a disciplina eh so o `@` da propria tarefa.

Os testes do `cmd/neivi` rodam o binario de verdade, com um banco temporario por caso. Os criterios que dependem do dia de hoje (prazos com dia da semana) sao conferidos nas funcoes de formatacao e de data, com o "hoje" fixo em ter 10/11/2026, porque os comandos usam o relogio real.
- **D-34** Nova ordem das fatias: 4 (listagem e resumo) -> 5 (importacao) -> corte das tarefas -> 2 (disciplinas e projetos) -> 3 (habitos) -> corte dos habitos e v0.1. O corte seco (D-24) vira duas etapas: depois da importacao, tarefas so no neivi e o Ganbatte fica ligado so para os habitos; depois da fatia 3, o Ganbatte eh desligado. Os numeros das fatias continuam os mesmos.
  Porque: preciso comecar a usar o basico (anotar e ver minhas tarefas) o quanto antes, e o uso real melhora os requisitos das fatias seguintes. Custo aceito: por um tempo, tarefas no neivi e habitos no Ganbatte; mudancas no banco depois do corte viram migracoes (sem perder dados).
