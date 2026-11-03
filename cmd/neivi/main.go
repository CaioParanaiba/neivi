package main

import (
	"fmt"
	"os"

	"github.com/CaioParanaiba/neivi/internal/store"
)

const dbPath = "/tmp/neivi-teste.db"

func main() {
	db, err := store.Open(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	defer db.Close()

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		runAdd(db, os.Args[2:])
	case "ls":
		runLs(db, os.Args[2:])
	case "subject":
		runSubject(db, os.Args[2:])
	case "project":
		runProject(db, os.Args[2:])
	case "done":
		runDoneTasks(db, os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "Comando", os.Args[1], "eh desconhecido")
		usage()
		os.Exit(1)
	}
}
