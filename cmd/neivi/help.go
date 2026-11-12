package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// helps eh o texto do neivi help <comando> (US-22): o formato e os simbolos de cada comando
var helps = map[string]string{
	"add": `uso: neivi add <titulo> [marcas]

  @BD        disciplina (cadastrada)        pri:A    prioridade A, B, C ou D
  +neivi     projeto (cadastrado)           due:sex  prazo: dd/mm, dd/mm/aa, hoje, amanha, seg...dom
  est:1h15   estimativa (30m, 1h, 90)       -- txt   descricao (um segundo -- fecha)

  as marcas vao em qualquer posicao; palavra com : que nao eh marca fica no titulo

exemplos:
  neivi add comprar cabo HDMI
  neivi add entregar lista @ED pri:A due:sex est:1h -- exercicios 3 a 7`,

	"edit": `uso: neivi edit <id> [titulo novo] [marcas]

  so muda o que aparece; as marcas sao as do add (neivi help add)
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

  atrasadas, urgentes (prazo hoje) e importantes (A sem prazo; sem nenhuma, as sem prazo de maior prioridade),
  o tempo estimado do que falta, o contador do dia e os habitos que faltam`,

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

	"show": `uso: neivi show <id>

  todos os campos de uma tarefa, a situacao e a descricao inteira`,

	"subject": `uso: neivi subject <add|ls|edit|rm> ...

  neivi subject add <codigo> [nome]
  neivi subject ls [all]                 all: tambem as encerradas
  neivi subject edit <codigo> ...        nome novo, sit:ativa|encerrada, code:novo
  neivi subject rm <codigo>              so sem tarefas nem projetos ligados

exemplos:
  neivi subject add BD Projeto de Banco de Dados
  neivi subject edit BD sit:encerrada`,

	"project": `uso: neivi project <add|ls|show|edit|rm> ...

  neivi project add <codigo> [nome] [campos]
      @BD                 disciplina (projeto academico)
      gh:usuario/repo     GitHub (ou a URL inteira)
      dir:~/Code/x        pasta local
      cont:sim            projeto continuo (sem %)
      sit:pausado         ativo, pausado, concluido ou cancelado
      -- texto            descricao
  neivi project ls [all]                 all: tambem concluidos e cancelados
  neivi project show <codigo>            panorama com as pendentes
  neivi project edit <codigo> ...        os mesmos campos; - apaga (dir:-); code:novo
  neivi project rm <codigo>              so sem tarefas

exemplos:
  neivi project add neivi Meu sucessor do Ganbatte gh:CaioParanaiba/neivi dir:~/Code/neivi
  neivi project edit ganbatte sit:pausado`,

	"habit": `uso: neivi habit <add|do|undo|ls|edit|rm> ...

  neivi habit add <codigo> [nome] [meta:3 | sem:3 | dias:seg,qua,sex]
      sem nada: todo dia, feito ou nao feito
      meta:3   quantidade por dia        sem:3   3 vezes por semana (dom a sab)
      dias:    so nesses dias (dom seg ter qua qui sex sab)
  neivi habit do <habito>... [ontem]     codigo ou id; x2 depois do habito = quantidade
  neivi habit undo <habito>... [ontem]   desfaz (nunca fica negativo)
  neivi habit ls [all]                   all: tambem os arquivados
  neivi habit edit <codigo> ...          nome, meta:, sem:, dias:, sit:ativo|arquivado, code:novo
  neivi habit rm <codigo>                so sem registros

exemplos:
  neivi habit add cf Codeforces 3 questoes meta:3
  neivi habit do cf x2 lixo neivi
  neivi habit do lixo ontem`,
}

// runHelp eh o neivi help [comando]
func runHelp(args []string) {
	if len(args) == 0 {
		fmt.Print(paintHelp(help))
		return
	}

	text, ok := helps[args[0]]
	if !ok {
		var names []string
		for k := range helps {
			names = append(names, k)
		}
		sort.Strings(names)
		fmt.Fprintf(os.Stderr, "erro: comando %s desconhecido (ajuda existe para: %s)\n", args[0], strings.Join(names, ", "))
		os.Exit(1)
	}

	fmt.Println(paintHelp(text))
}
