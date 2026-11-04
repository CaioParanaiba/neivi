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
		Priority:    entry.Priority,
		Due:         entry.Due,
		EstimateMin: entry.EstimateMin,
		SubjectID:   idSubject,
		ProjectID:   idProject,
		//parenteID
		Extras: entry.Extras,
	}

	id, err := store.AddTask(db, t)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	fmt.Println("Tarefa:", id, "adicionada")
}

func runLs(db *sql.DB, args []string) {
	lista, err := store.ListTasks(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tX\tP\tPRAZO\tTAREFA")

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

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s", v.ID, done, v.Priority, v.Due, v.Title)
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
