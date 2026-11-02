package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/CaioParanaiba/neivi/internal/store"
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

// runAdd adiciona uma task
func runAdd(db *sql.DB, args []string) {
	title := strings.Join(args, " ")

	if title == "" {
		fmt.Fprintln(os.Stderr, "erro: tarefa sem titulo")
		os.Exit(1)
	}

	id, err := store.AddTask(db, title)
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

	for _, v := range lista {
		fmt.Println(v.ID, v.Title, v.CreatedAt)
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
