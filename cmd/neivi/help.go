package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// helps eh o texto do neivi <comando> -help (US-22): o formato e os simbolos de cada comando
var helps = map[string]string{
	"add": `uso: neivi add <titulo> [marcas]

  @BD        disciplina (cadastrada)        pri:A    prioridade A, B, C ou D
  +neivi     projeto (cadastrado)           due:sex  prazo: dd/mm, dd/mm/aa, hoje, amanha, seg...dom
  est:1h15   estimativa (30m, 1h, 90)       -- txt   descricao (um segundo -- fecha)

  as marcas vao em qualquer posicao; palavra com : que nao eh marca fica no titulo
  ; & | # ( ) > sao do shell e cortam o comando: use aspas ou troque por , ou -
      neivi add "ler cap 3; fazer resumo"

exemplos:
  neivi add comprar cabo HDMI
  neivi add entregar lista @ED pri:A due:sex est:1h -- exercicios 3 a 7`,

	"edit": `uso: neivi edit <id> [titulo novo] [marcas]

  so muda o que aparece; as marcas sao as do add (neivi add -help)
  -  apaga o campo: pri:-  due:-  est:-  @-  +-  -- -   (o titulo nao pode ser apagado)

exemplos:
  neivi edit 4 due:seg
  neivi edit 4 pri:B @BD
  neivi edit 4 entregar lista 2`,

	"done": `uso: neivi done <id>... [ontem]

  conclui uma ou varias tarefas; ontem vale para todas do comando
  id que nao existe so avisa; outra palavra recusa tudo

exemplos:
  neivi done 3 5 8
  neivi done 3 ontem`,

	"reopen": `uso: neivi reopen <id>...

  desfaz a conclusao (a tarefa volta a pendente, com tudo como estava)

exemplo:
  neivi reopen 5`,

	"rm": `uso: neivi rm <id>...

  apaga de vez, perguntando [s/N] para cada uma (so s confirma)

exemplo:
  neivi rm 3 5`,

	"today": `uso: neivi today

  atrasadas, urgentes (prazo hoje), amanha (na sex, sab e dom: ate segunda)
  e importantes (A sem prazo; sem nenhuma, as sem prazo de maior prioridade)
  sem prazo de projeto (menos as A) so contam na linha "+ N de projetos": ver no neivi project <codigo>
  o tempo estimado do que falta hoje (o amanha nao entra), o contador do dia e os habitos que faltam
  habito "N por semana" sai no dia em que foi feito e volta no dia seguinte`,

	"triage": `uso: neivi triage

  as tarefas sem prioridade (incompletas), das mais antigas para as mais novas
  o mesmo que: neivi ls pri:-`,

	"ls": `uso: neivi ls [filtros] [palavras]

  @BD  @  @-       essa disciplina, qualquer uma, sem disciplina
  +neivi  +  +-    esse projeto, qualquer um, sem projeto
  pri:A  pri:-     essa prioridade, sem prioridade
  new              criadas hoje ou ontem
  late             atrasadas
  all / done       com as concluidas no fim / so as concluidas
  due              ordena por prazo
  desc             mostra a descricao embaixo
  outra palavra    busca no titulo (todas as palavras; sem maiusculas nem acentos)

  mesmo tipo eh "ou" (@BD @ED); tipos diferentes sao "e" (@BD pri:A)

exemplos:
  neivi ls @BD pri:A
  neivi ls +mrt hull`,

	"pri": `prioridade mede a importancia; a urgencia eh o prazo (due:)

     definicao                               pergunta-teste
  A  consequencia seria ou com outra pessoa  se eu atrasar, perco nota, dinheiro ou prejudico alguem?
  B  consequencia real, mas recuperavel      se eu atrasar, eu sofro, mas consigo recuperar?
  C  seria bom, sem consequencia seria       ninguem sofre se atrasar, mas quero fazer?
  D  talvez um dia                           se ficar 3 meses parado, tudo bem?
  -  ainda nao decidida (incompleta)

dar prioridade: neivi edit <id> pri:B        ver as incompletas: neivi triage`,

	"show": `uso: neivi show <id>   (ou so neivi <id>)

  todos os campos de uma tarefa, a situacao e a descricao inteira`,

	"subject": `uso: neivi subject [add|ls|edit|rm] ...

  neivi subject                          o mesmo que subject ls
  neivi subject add <codigo> [nome]
  neivi subject ls [all]                 all: tambem as encerradas
  neivi subject edit <codigo> ...        nome novo, sit:ativa|encerrada, code:novo
  neivi subject rm <codigo>              so sem tarefas nem projetos ligados

exemplos:
  neivi subject add BD Projeto de Banco de Dados
  neivi subject edit BD sit:encerrada`,

	"project": `uso: neivi project [<codigo>|add|ls|show|edit|rm] ...

  neivi project                          o mesmo que project ls
  neivi project <codigo>                 o mesmo que project show (codigo, +codigo ou id)
  neivi project add <codigo> [nome] [campos]
      @BD                 disciplina (projeto academico)
      gh:usuario/repo     GitHub (ou a URL inteira)
      dir:~/Code/x        pasta local
      cont:sim            projeto continuo (sem %)
      sit:pausado         ativo, pausado, concluido ou cancelado
      -- texto            descricao
  neivi project ls [all]                 all: tambem concluidos e cancelados
  neivi project show <codigo>            panorama com as pendentes (prioridade, prazo, id)
  neivi project edit <codigo> ...        os mesmos campos; - apaga (dir:-); code:novo
  neivi project rm <codigo>              so sem tarefas

exemplos:
  neivi project mrt
  neivi project add neivi Meu sucessor do Ganbatte gh:CaioParanaiba/neivi dir:~/Code/neivi
  neivi project edit ganbatte sit:pausado`,

	"habit": `uso: neivi habit [<habito>|add|do|undo|ls|edit|rm] ...

  neivi habit                            o mesmo que habit ls
  neivi habit <habito>... [ontem]        o mesmo que habit do, mas o que ja esta feito no dia so avisa
  neivi habit add <codigo> [nome] [meta:3 | sem:3 | dias:seg,qua,sex]
      sem nada: todo dia, feito ou nao feito
      meta:3   quantidade por dia        sem:3   3 vezes por semana (dom a sab)
      dias:    so nesses dias (dom seg ter qua qui sex sab)
  neivi habit do <habito>... [ontem]     codigo ou id; x2 depois do habito = quantidade (done tambem vale)
  neivi habit undo <habito>... [ontem]   desfaz (nunca fica negativo)
  neivi habit ls [all]                   all: tambem os arquivados
  neivi habit edit <codigo> ...          nome, meta:, sem:, dias:, sit:ativo|arquivado, code:novo
  neivi habit rm <codigo>                so sem registros

exemplos:
  neivi habit digi
  neivi habit add cf Codeforces 3 questoes meta:3
  neivi habit do cf x2 lixo neivi
  neivi habit do lixo ontem`,
}

// runHelp eh o neivi -help (topic vazio) e o neivi <comando> -help (D-39)
func runHelp(topic string) {
	if topic == "" {
		fmt.Print(paintHelp(help))
		return
	}

	text, ok := helps[topic]
	if !ok {
		var names []string
		for k := range helps {
			names = append(names, k)
		}
		sort.Strings(names)
		fmt.Fprintf(os.Stderr, "erro: comando %s desconhecido (ajuda existe para: %s)\n", topic, strings.Join(names, ", "))
		os.Exit(1)
	}

	fmt.Println(paintHelp(text))
}
