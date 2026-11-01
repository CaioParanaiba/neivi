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
	add	adicionar tarefa
	ls	listar tarefas
`

func usage() {
	fmt.Fprint(os.Stderr, help)
}

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
	tl, err := store.ListTasks(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	for _, t := range tl {
		fmt.Println(t.ID, t.Title, t.CreatedAt)
	}
}
