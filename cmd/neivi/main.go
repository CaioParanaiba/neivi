package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CaioParanaiba/neivi/internal/store"
)

// dbPath devolve o caminho do banco: NEIVI_DB, se definida (testes), ou
// $XDG_DATA_HOME/neivi/neivi.db (padrao ~/.local/share/neivi/neivi.db).
// Cria a pasta se ela ainda nao existir.
func dbPath() (string, error) {
	if p := os.Getenv("NEIVI_DB"); p != "" {
		return p, nil
	}

	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("achar a pasta pessoal: %w", err)
		}
		dir = filepath.Join(home, ".local", "share")
	}
	dir = filepath.Join(dir, "neivi")

	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return "", fmt.Errorf("criar %s: %w", dir, err)
	}

	return filepath.Join(dir, "neivi.db"), nil
}

func main() {
	path, err := dbPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}

	db, err := store.Open(path)
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
	case "today":
		runToday(db, os.Args[2:])
	case "triage":
		runTriage(db, os.Args[2:])
	case "show":
		runShow(db, os.Args[2:])
	case "subject":
		runSubject(db, os.Args[2:])
	case "project":
		runProject(db, os.Args[2:])
	case "done":
		runDoneTasks(db, os.Args[2:])
	case "reopen":
		runReopenTasks(db, os.Args[2:])
	case "rm":
		runRmTasks(db, os.Args[2:])
	case "edit":
		runEdit(db, os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "Comando", os.Args[1], "eh desconhecido")
		usage()
		os.Exit(1)
	}
}
