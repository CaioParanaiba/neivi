package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/CaioParanaiba/neivi/internal/store"
	"github.com/CaioParanaiba/neivi/internal/task"
	"github.com/CaioParanaiba/neivi/internal/todotxt"
)

// table guarda as linhas de uma ou mais tabelas para alinhar todas com a mesma
// largura de coluna (o today tem um grupo por tabela, RNF-05)
type table struct {
	widths []int
}

// fit aumenta as larguras para caber a linha (conta letras, nao bytes: "…" e acentos)
func (tb *table) fit(cells []string) {
	for i, c := range cells {
		n := utf8.RuneCountInString(c)
		if i >= len(tb.widths) {
			tb.widths = append(tb.widths, n)
		} else if n > tb.widths[i] {
			tb.widths[i] = n
		}
	}
}

// line monta a linha com as colunas alinhadas; a ultima coluna nao tem espaco no fim
func (tb *table) line(cells []string) string {
	return tb.paintedLine(cells, cells)
}

// paintedLine eh o line com cor: escreve as celulas pintadas, mas o espaco vem
// do texto puro (os codigos de cor nao ocupam lugar na tela)
func (tb *table) paintedLine(cells, painted []string) string {
	var b strings.Builder
	for i, c := range cells {
		b.WriteString(painted[i])
		if i == len(cells)-1 {
			break
		}
		pad := tb.widths[i] - utf8.RuneCountInString(c) + 2
		b.WriteString(strings.Repeat(" ", pad))
	}
	return strings.TrimRight(b.String(), " ")
}

// header eh o cabecalho das listas; o X so aparece no ls all
func header(withX bool) []string {
	if withX {
		return []string{"ID", "X", "P", "PRAZO", "EST", "TAREFA"}
	}
	return []string{"ID", "P", "PRAZO", "EST", "TAREFA"}
}

// row monta as celulas de uma tarefa na tabela
func row(t task.Task, withX bool, withDesc bool, today time.Time) []string {
	due := ""
	if t.Due != "" {
		due = formatDue(t.Due, today)
	}

	title := t.Title
	if t.SubjectCode != "" {
		title += " @" + t.SubjectCode
	}
	if t.ProjectCode != "" {
		title += " +" + t.ProjectCode
	}
	if t.Description != "" && !withDesc {
		title += " …"
	}

	cells := []string{strconv.FormatInt(t.ID, 10)}
	if withX {
		x := ""
		if t.DoneAt != "" {
			x = "x"
		}
		cells = append(cells, x)
	}
	return append(cells, t.Priority, due, formatEst(t.EstimateMin), title)
}

// descRow eh a linha da descricao embaixo da tarefa (ls desc), na coluna do titulo
func descRow(t task.Task, withX bool) []string {
	n := len(header(withX))
	cells := make([]string, n)
	cells[n-1] = "  " + t.Description
	return cells
}

// printTasks imprime uma tabela simples (ls e triage)
func printTasks(ts []task.Task, withX bool, withDesc bool, today time.Time) {
	var tb table
	tb.fit(header(withX))

	var rows, painted [][]string
	for _, t := range ts {
		r := row(t, withX, withDesc, today)
		tb.fit(r)
		rows = append(rows, r)
		painted = append(painted, paintRow(r, withX))
		if withDesc && t.Description != "" {
			d := descRow(t, withX)
			rows = append(rows, d)
			painted = append(painted, d)
		}
	}

	fmt.Println(paint(stDim, tb.line(header(withX))))
	for i, r := range rows {
		fmt.Println(tb.paintedLine(r, painted[i]))
	}
}

// noArgs recusa argumentos nos comandos que nao recebem nenhum
func noArgs(cmd string, args []string) {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "erro: %s nao recebe argumentos (exemplo: neivi %s)\n", cmd, cmd)
		os.Exit(1)
	}
}

func runLs(db *sql.DB, args []string) {
	f, err := todotxt.ParseFilter(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	// disciplina ou projeto que nao existe eh erro (CA-09.9)
	for _, code := range f.Subjects {
		_, err := store.SubjectIDByCode(db, code)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
	}
	for _, code := range f.Projects {
		_, err := store.ProjectIDByCode(db, code)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
	}

	lista, err := store.ListTasksFilter(db, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	if len(lista) == 0 {
		if f.IsEmpty() {
			fmt.Println("nenhuma tarefa pendente")
		} else {
			fmt.Println("nenhuma tarefa com esses filtros (tente menos palavras ou filtros)")
		}
		return
	}

	// ordem (RN-17): pendentes pela ordem padrao ou por prazo; concluidas, mais recentes primeiro
	var pending, done []task.Task
	for _, t := range lista {
		if t.DoneAt == "" {
			pending = append(pending, t)
		} else {
			done = append(done, t)
		}
	}
	if f.Due {
		task.SortByDue(pending)
	} else {
		task.SortDefault(pending)
	}
	task.SortDoneRecent(done)

	printTasks(append(pending, done...), f.All, f.Desc, time.Now())
}

func runToday(db *sql.DB, args []string) {
	noArgs("today", args)

	pending, err := store.ListTasksFilter(db, task.Filter{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	now := time.Now()
	s := task.Today(pending, now)

	if s.IsEmpty() {
		fmt.Println("nada para hoje, bom trabalho!")

		incomplete := 0
		for _, t := range pending {
			if t.Priority == "" {
				incomplete++
			}
		}
		if incomplete > 0 {
			fmt.Printf("voce tem %d tarefas incompletas: que tal completar? (neivi triage)\n", incomplete)
		}
	} else {
		third := "importantes"
		if !s.Important {
			third = "sem prazo"
		}

		groups := []struct {
			title string
			ts    []task.Task
		}{
			{"atrasadas", s.Late},
			{"urgentes", s.Urgent},
			{third, s.Third},
		}

		// todas as linhas primeiro, para as larguras valerem para todos os grupos
		var tb table
		tb.fit(header(false))
		var est int64
		var noEst, shown int
		for _, g := range groups {
			for _, t := range g.ts {
				tb.fit(row(t, false, false, now))
				shown++
				if t.EstimateMin > 0 {
					est += t.EstimateMin
				} else {
					noEst++
				}
			}
		}

		first := true
		for _, g := range groups {
			if len(g.ts) == 0 {
				continue // bloco vazio nao aparece (RN-16)
			}
			if !first {
				fmt.Println()
			}
			first = false

			fmt.Println(paint(stSection, fmt.Sprintf("%s (%d)", g.title, len(g.ts))))
			fmt.Println(paint(stDim, tb.line(header(false))))
			for _, t := range g.ts {
				r := row(t, false, false, now)
				fmt.Println(tb.paintedLine(r, paintRow(r, false)))
			}
		}

		fmt.Println()
		var falta string
		switch {
		case est == 0:
			falta = fmt.Sprintf("falta: sem EST (%d tarefas)", shown)
		case noEst > 0:
			falta = fmt.Sprintf("falta: ~%s (%d sem EST)", formatEst(est), noEst)
		default:
			falta = fmt.Sprintf("falta: ~%s", formatEst(est))
		}
		fmt.Println(paint(stFooter, falta))
	}

	n, err := store.CountDone(db, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	fmt.Println(paint(stFooter, fmt.Sprintf("hoje: %d concluidas", n)))

	printHabitsToday(db, now)
}

func runTriage(db *sql.DB, args []string) {
	noArgs("triage", args)

	lista, err := store.ListTasksFilter(db, task.Filter{Priorities: []string{"-"}})
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	if len(lista) == 0 {
		fmt.Println("nenhuma tarefa incompleta")
		return
	}

	task.SortByCreated(lista)

	fmt.Println(paint(stSection, fmt.Sprintf("incompletas (%d)", len(lista))))
	printTasks(lista, false, false, time.Now())
	fmt.Println()
	fmt.Println(paint(stDim, "criterio de prioridade: neivi help pri"))
}

// daysLabel diz quanto falta para o prazo: hoje, amanha, em N dias, atrasada ha N dias
func daysLabel(due string, today time.Time) string {
	d, err := time.ParseInLocation("2006-01-02", due, today.Location())
	if err != nil {
		return ""
	}
	t := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	n := int(d.Sub(t).Hours() / 24)

	switch {
	case n == 0:
		return "hoje"
	case n == 1:
		return "amanha"
	case n > 1:
		return fmt.Sprintf("em %d dias", n)
	case n == -1:
		return "atrasada ha 1 dia"
	default:
		return fmt.Sprintf("atrasada ha %d dias", -n)
	}
}

// codeName mostra "BD (Projeto de Banco de Dados)"; sem nome proprio, so o codigo
func codeName(code, name string) string {
	if code == "" {
		return "-"
	}
	if name == "" || name == code {
		return code
	}
	return code + " (" + name + ")"
}

func runShow(db *sql.DB, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "erro: falta o id da tarefa (exemplo: neivi show 2)")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "erro: show mostra uma tarefa por vez (exemplo: neivi show 2)")
		os.Exit(1)
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: id %s invalido (use o numero que aparece no ls)\n", args[0])
		os.Exit(1)
	}

	t, err := store.GetTask(db, id)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintf(os.Stderr, "erro: tarefa %d nao existe (veja os ids no neivi ls)\n", id)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	now := time.Now()
	dash := func(s string) string {
		if s == "" {
			return "-"
		}
		return s
	}

	due := "-"
	if t.Due != "" {
		due = formatDue(t.Due, now)
		if t.DoneAt == "" {
			due += " (" + daysLabel(t.Due, now) + ")"
		}
	}

	created := formatDue(t.CreatedAt[:10], now)
	if task.IsNew(t, now) {
		created += " (nova)"
	}

	situation := "pendente"
	switch {
	case t.DoneAt != "":
		situation = "concluida em " + formatDue(t.DoneAt[:10], now)
	case t.Priority == "":
		situation = "incompleta"
	}

	fmt.Println(paint(stSection, fmt.Sprintf("tarefa %d", t.ID)) + " · " + paint(stBold, t.Title))
	// campo vazio ("-") fica apagado; os outros usam a cor do papel, como no today
	field := func(name, value string, st style) {
		if value == "-" {
			st = stDim
		}
		fmt.Printf("  %s %s\n", paint(stDim, fmt.Sprintf("%-12s", name+":")), paint(st, value))
	}
	field("prioridade", dash(t.Priority), priStyle(t.Priority))
	field("prazo", due, stDue)
	field("EST", dash(formatEst(t.EstimateMin)), stEst)
	disc := codeName(t.SubjectCode, t.SubjectName)
	if t.SubjectInherited {
		disc = t.SubjectCode + " (do projeto)" // RN-14
	}
	field("disciplina", disc, stSubject)
	field("projeto", codeName(t.ProjectCode, t.ProjectName), stProject)
	field("criada", created, "")
	field("situacao", situation, "")

	if t.Description != "" {
		fmt.Println()
		fmt.Println("  " + t.Description)
	}
}
