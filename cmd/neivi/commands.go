package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CaioParanaiba/neivi/internal/store"
	"github.com/CaioParanaiba/neivi/internal/task"
	"github.com/CaioParanaiba/neivi/internal/todotxt"
)

const help = `uso: neivi <comando> [argumentos]

tarefas:
  add <texto>                     adicionar tarefa
  today                           o que tem para hoje (tarefas e habitos)
  triage                          tarefas incompletas (sem prioridade)
  ls [filtros] [palavras]         listar tarefas
  show <id>                       ver todos os campos de uma tarefa
  edit <id> <texto>               mudar so os campos do texto
  done <id>... [ontem]            concluir tarefas
  reopen <id>...                  reabrir tarefas concluidas
  rm <id>...                      remover tarefas (pede confirmacao)

cadastros:
  subject <add|ls|edit|rm>        disciplinas
  project <add|ls|show|edit|rm>   projetos
  habit <add|do|undo|ls|edit|rm>  habitos

detalhes de um comando: neivi help <comando>   (ex.: neivi help add)
`

func usage() {
	fmt.Fprint(os.Stderr, help)
}

// dada a string dos codigos retorna os ids correspondentes
func resolveCode(db *sql.DB, subjectCode string, projectCode string) (int64, int64, error) {

	var idSubject int64
	var err error
	if subjectCode != "" {
		idSubject, err = store.SubjectIDByCode(db, subjectCode)
		if err != nil {
			return 0, 0, err
		}
	}

	var idProject int64
	if projectCode != "" {
		idProject, err = store.ProjectIDByCode(db, projectCode)
		if err != nil {
			return idSubject, 0, err
		}
	}

	return idSubject, idProject, nil
}

// runAdd adiciona uma task
func runAdd(db *sql.DB, args []string) {
	entry, err := todotxt.ParseTask(args, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	idSubject, idProject, err := resolveCode(db, entry.SubjectCode, entry.ProjectCode)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	t := task.Task{
		Title:       entry.Title,
		Description: entry.Description,
		Priority:    entry.Priority,
		Due:         entry.Due,
		EstimateMin: entry.EstimateMin,
		SubjectID:   idSubject,
		ProjectID:   idProject,
	}

	// procura antes de inserir, para a propria tarefa nao contar (RN-01)
	dup, err := store.FindDuplicate(db, t.Title, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	id, err := store.AddTask(db, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	t.ID = id
	t.SubjectCode = entry.SubjectCode
	t.ProjectCode = entry.ProjectCode
	fmt.Println(formatLine("+", t, time.Now()))

	if dup != 0 {
		fmt.Fprintf(os.Stderr, "aviso: ja existe a tarefa %d com esse titulo\n", dup)
	}
	warnInactive(db, entry.SubjectCode, entry.ProjectCode)
}

// formatEst mostra a estimativa como 30m, 1h ou 1h15; vazia quando não tem
func formatEst(min int64) string {
	if min <= 0 {
		return ""
	}
	h, m := min/60, min%60
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh%02d", h, m)
}

// dias da semana na ordem do time.Weekday (domingo = 0)
var weekdays = [...]string{"dom", "seg", "ter", "qua", "qui", "sex", "sab"}

// formatDue mostra o prazo como "sex 13/11"; com o ano so quando nao eh o atual
func formatDue(due string, today time.Time) string {
	d, err := time.Parse("2006-01-02", due)
	if err != nil {
		return due
	}

	if d.Year() != today.Year() {
		return weekdays[d.Weekday()] + " " + d.Format("02/01/06")
	}
	return weekdays[d.Weekday()] + " " + d.Format("02/01")
}

// linePart monta "id · pri · prazo · est · titulo @disc +proj …", sem os campos vazios
func linePart(t task.Task, today time.Time) string {
	parts := []string{strconv.FormatInt(t.ID, 10)}

	if t.Priority != "" {
		parts = append(parts, t.Priority)
	}
	if t.Due != "" {
		parts = append(parts, formatDue(t.Due, today))
	}
	if t.EstimateMin > 0 {
		parts = append(parts, formatEst(t.EstimateMin))
	}

	title := t.Title
	if t.SubjectCode != "" {
		title += " @" + t.SubjectCode
	}
	if t.ProjectCode != "" {
		title += " +" + t.ProjectCode
	}
	if t.Description != "" {
		title += " …"
	}
	parts = append(parts, title)

	return strings.Join(parts, " · ")
}

// formatLine eh a confirmacao curta (RNF-05): + criar, ~ editar, ✓ concluir, ↺ reabrir
func formatLine(symbol string, t task.Task, today time.Time) string {
	return symbol + " " + linePart(t, today)
}

// parseIDs le os ids de done, reopen e rm antes de agir em qualquer um (RN-12).
// Com allowOntem, aceita a palavra "ontem". Qualquer outra palavra recusa tudo.
func parseIDs(args []string, allowOntem bool, example string) ([]int64, bool, error) {
	var ids []int64
	var ontem bool

	for _, v := range args {
		if allowOntem && v == "ontem" {
			ontem = true
			continue
		}

		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			if allowOntem {
				return nil, false, fmt.Errorf("%q nao eh um id nem \"ontem\" (exemplo: %s)", v, example)
			}
			return nil, false, fmt.Errorf("%q nao eh um id (exemplo: %s)", v, example)
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, false, fmt.Errorf("falta o id da tarefa (exemplo: %s)", example)
	}

	return ids, ontem, nil
}

// testar fazer uma funcao bonita ksksks
func applyText(field *string, v string) {
	if v == "-" {
		*field = ""
	} else if v != "" {
		*field = v
	}
}

func applyEdit(db *sql.DB, t *task.Task, e todotxt.Entry) error {
	applyText(&t.Title, e.Title)
	applyText(&t.Description, e.Description)
	applyText(&t.Priority, e.Priority)
	applyText(&t.Due, e.Due)

	if e.EstimateMin == -1 {
		t.EstimateMin = 0
	} else if e.EstimateMin != 0 {
		t.EstimateMin = e.EstimateMin
	}

	sCode := e.SubjectCode
	pCode := e.ProjectCode
	if sCode == "-" {
		sCode = ""
	}
	if pCode == "-" {
		pCode = ""
	}

	idS, idP, err := resolveCode(db, sCode, pCode)
	if err != nil {
		return err
	}

	if e.SubjectCode == "-" {
		t.SubjectID = 0
	} else if e.SubjectCode != "" {
		t.SubjectID = idS
	}

	if e.ProjectCode == "-" {
		t.ProjectID = 0
	} else if e.ProjectCode != "" {
		t.ProjectID = idP
	}

	return nil
}

func runEdit(db *sql.DB, args []string) {

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "erro: falta o id da tarefa (exemplo: neivi edit 1 pri:B)")
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

	args = args[1:]

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "erro: nada para mudar na tarefa %d (exemplo: neivi edit %d pri:B due:sex)\n", id, id)
		os.Exit(1)
	}

	entry, err := todotxt.ParseEdit(args, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	err = applyEdit(db, &t, entry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	// titulo novo: procura duplicada sem contar a propria tarefa (RN-01)
	var dup int64
	if entry.Title != "" {
		dup, err = store.FindDuplicate(db, t.Title, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
	}

	err = store.UpdateTask(db, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	// busca de novo para os codigos de @ e + virem atualizados
	t, err = store.GetTask(db, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	fmt.Println(formatLine("~", t, time.Now()))

	if t.DoneAt != "" {
		fmt.Fprintf(os.Stderr, "aviso: a tarefa %d esta concluida\n", id)
	}
	if dup != 0 {
		fmt.Fprintf(os.Stderr, "aviso: ja existe a tarefa %d com esse titulo\n", dup)
	}
	warnInactive(db, entry.SubjectCode, entry.ProjectCode)
}

func runDoneTasks(db *sql.DB, args []string) {
	ids, ontem, err := parseIDs(args, true, "neivi done 3 5 ontem")
	if err != nil {
		if len(args) == 0 {
			err = fmt.Errorf("falta o id da tarefa (exemplo: neivi done 3 5)")
		}
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	daysAgo := 0
	suffix := ""
	label := "hoje"
	if ontem {
		daysAgo = 1
		suffix = " (ontem)"
		label = "ontem"
	}

	today := time.Now()
	var changed int
	for _, id := range ids {
		v, err := store.GetTask(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: tarefa %d nao existe (veja os ids no neivi ls)\n", id)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}

		// concluir de novo nao muda a data (RN-11)
		if v.DoneAt != "" {
			fmt.Fprintf(os.Stderr, "aviso: a tarefa %d ja estava concluida\n", id)
			continue
		}

		err = store.DoneTask(db, id, daysAgo)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
		fmt.Println(formatLine("✓", v, today) + suffix)
		changed++
	}

	if changed > 0 {
		n, err := store.CountDone(db, daysAgo)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
		fmt.Printf("%s: %d concluidas\n", label, n)
	}
}

func runReopenTasks(db *sql.DB, args []string) {
	ids, _, err := parseIDs(args, false, "neivi reopen 5")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	today := time.Now()
	for _, id := range ids {
		v, err := store.GetTask(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: tarefa %d nao existe (veja os ids no neivi ls)\n", id)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}

		if v.DoneAt == "" {
			fmt.Fprintf(os.Stderr, "aviso: a tarefa %d nao estava concluida\n", id)
			continue
		}

		err = store.ReopenTask(db, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
		fmt.Println(formatLine("↺", v, today))
	}
}

func runRmTasks(db *sql.DB, args []string) {
	ids, _, err := parseIDs(args, false, "neivi rm 3")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	// respostas linha a linha
	// um reader só para o laço
	reader := bufio.NewReader(os.Stdin)

	today := time.Now()
	for _, id := range ids {
		v, err := store.GetTask(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: tarefa %d nao existe (veja os ids no neivi ls)\n", id)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}

		// concluida sai do historico: avisa na pergunta (RN-13)
		warn := ""
		if v.DoneAt != "" {
			warn = " (concluida, sai do historico)"
		}

		fmt.Fprintf(os.Stdout, "remover %s%s? [s/N] ", linePart(v, today), warn)
		line, _ := reader.ReadString('\n')
		ans := strings.TrimSpace(line)

		ans = strings.ToUpper(ans)

		if ans != "S" {
			fmt.Fprintf(os.Stdout, "%d mantida\n", id)
			continue
		}

		err = store.RmTask(db, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "✗ %d removida\n", id)
	}
}
