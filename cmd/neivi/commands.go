package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/CaioParanaiba/neivi/internal/store"
	"github.com/CaioParanaiba/neivi/internal/task"
	"github.com/CaioParanaiba/neivi/internal/todotxt"
)

const help = `uso: neivi <comando> [argumentos]

comandos:
  add <texto>                adicionar tarefa
  ls                         listar tarefas
  done <id>...               concluir tarefas
  reopen <id>...             reabrir tarefas concluidas
  edit <id> <texto>          mudar so os campos do texto
  rm <id>...                 remover tarefas (pede confirmacao)
  subject add <code> [nome]  adicionar disciplina
  subject ls                 listar disciplinas
  project add <code> [nome]  adicionar projeto
  project ls                 listar projetos
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

	id, err := store.AddTask(db, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	fmt.Println("Tarefa:", id, "adicionada")
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

func runLs(db *sql.DB, args []string) {
	lista, err := store.ListTasks(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tX\tP\tPRAZO\tEST\tTAREFA")

	for _, v := range lista {
		done := ""
		if v.DoneAt != "" {
			done = "x"
		}

		if v.Due != "" {
			d, err := time.Parse("2006-01-02", v.Due)
			if err != nil {
				fmt.Fprintln(os.Stderr, "erro:", err)
				os.Exit(1)
			}
			v.Due = d.Format("02/01")
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s", v.ID, done, v.Priority, v.Due, formatEst(v.EstimateMin), v.Title)
		if v.SubjectCode != "" {
			fmt.Fprintf(w, " @%s", v.SubjectCode)
		}
		if v.ProjectCode != "" {
			fmt.Fprintf(w, " +%s", v.ProjectCode)
		}
		fmt.Fprintf(w, "\n")
	}

	w.Flush()
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

	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: id %s eh invalido: %s\n", args[0], err)
		os.Exit(1)
	}

	t, err := store.GetTask(db, id)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintf(os.Stderr, "erro: tarefa %d nao existe\n", id)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	args = args[1:]

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

	err = store.UpdateTask(db, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "tarefa %d editada\n", id)
}

func runDoneTasks(db *sql.DB, args []string) {

	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	var idChanged, idSkipped []string
	for _, sId := range args {
		id, err := strconv.ParseInt(sId, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: id %s eh invalido: %s\n", sId, err)
			continue
		}

		v, err := store.GetTask(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: tarefa %d nao existe\n", id)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}

		if v.DoneAt != "" {
			idSkipped = append(idSkipped, sId)
			continue
		}

		err = store.DoneTask(db, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
		idChanged = append(idChanged, sId)
	}

	if len(idChanged) > 0 {
		s := strings.Join(idChanged, " ")
		fmt.Fprintf(os.Stdout, "tarefa(s): %s feitas\n", s)
	}
	if len(idSkipped) > 0 {
		s := strings.Join(idSkipped, " ")
		fmt.Fprintf(os.Stdout, "tarefa(s): %s ja feitas\n", s)
	}
}

func runReopenTasks(db *sql.DB, args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	var idChanged, idSkipped []string
	for _, sId := range args {
		id, err := strconv.ParseInt(sId, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: id %s eh invalido: %s\n", sId, err)
			continue
		}

		v, err := store.GetTask(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: tarefa %d nao existe\n", id)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}

		if v.DoneAt == "" {
			idSkipped = append(idSkipped, sId)
			continue
		}

		err = store.ReopenTask(db, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
		idChanged = append(idChanged, sId)
	}

	if len(idChanged) > 0 {
		s := strings.Join(idChanged, " ")
		fmt.Fprintf(os.Stdout, "tarefa(s): %s reabertas\n", s)
	}
	if len(idSkipped) > 0 {
		s := strings.Join(idSkipped, " ")
		fmt.Fprintf(os.Stdout, "tarefa(s): %s nao estavam concluidas\n", s)
	}
}

func runRmTasks(db *sql.DB, args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	// respostas linha a linha
	// um reader só para o laço
	reader := bufio.NewReader(os.Stdin)

	for _, sId := range args {
		id, err := strconv.ParseInt(sId, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: id %s eh invalido: %s\n", sId, err)
			continue
		}

		v, err := store.GetTask(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: tarefa %d nao existe\n", id)
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "remover a tarefa %d \"%s\"? [s/N] ", id, v.Title)
		line, _ := reader.ReadString('\n')
		ans := strings.TrimSpace(line)

		ans = strings.ToUpper(ans)

		if ans != "S" {
			continue
		}

		err = store.RmTask(db, id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "tarefa %d removida\n", id)
	}
}

func runSubject(db *sql.DB, args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		runSubjectAdd(db, args[1:])
	case "ls":
		runSubjectLs(db, args[1:])
	default:
		fmt.Fprintln(os.Stderr, "Comando", args[0], "eh desconhecido")
		usage()
		os.Exit(1)
	}
}

func runSubjectAdd(db *sql.DB, args []string) {
	if len(args) < 1 || args[0] == "" {
		fmt.Fprintln(os.Stderr, "erro: disciplina sem codigo")
		os.Exit(1)
	}

	code := args[0]

	name := code
	if len(args) > 1 {
		name = strings.Join(args[1:], " ")
	}

	id, err := store.AddSubject(db, code, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	fmt.Println("Disciplina:", id, "adicionada")
}

func runSubjectLs(db *sql.DB, args []string) {
	lista, err := store.ListSubjects(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	for _, v := range lista {
		fmt.Println(v.ID, v.Code, v.Name)
	}
}

func runProject(db *sql.DB, args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		runProjectAdd(db, args[1:])
	case "ls":
		runProjectLs(db, args[1:])
	default:
		fmt.Fprintln(os.Stderr, "Comando", args[0], "eh desconhecido")
		usage()
		os.Exit(1)
	}
}

func runProjectAdd(db *sql.DB, args []string) {
	if len(args) < 1 || args[0] == "" {
		fmt.Fprintln(os.Stderr, "erro: projeto sem codigo")
		os.Exit(1)
	}

	code := args[0]

	name := code
	if len(args) > 1 {
		name = strings.Join(args[1:], " ")
	}

	id, err := store.AddProject(db, code, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	fmt.Println("Projeto:", id, "adicionado")
}

func runProjectLs(db *sql.DB, args []string) {
	lista, err := store.ListProjects(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	for _, v := range lista {
		fmt.Println(v.ID, v.Code, v.Name)
	}
}
