package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

// isHelp diz se a palavra eh um pedido de ajuda (D-39)
func isHelp(s string) bool {
	return s == "-h" || s == "-help" || s == "--help"
}

// helpTopic acha o pedido de ajuda: so como primeira palavra depois do comando
// ou, nos grupos, depois do subcomando (CA-22.1, CA-22.5)
func helpTopic(args []string) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	if isHelp(args[0]) {
		return "", true // neivi -help: a lista de comandos
	}
	topic := args[0]
	if topic == "habits" {
		topic = "habit"
	}
	if len(args) > 1 && isHelp(args[1]) {
		return topic, true
	}
	group := topic == "subject" || topic == "project" || topic == "habit"
	if group && len(args) > 2 && isHelp(args[2]) {
		return topic, true
	}
	return "", false
}

func main() {
	if topic, ok := helpTopic(os.Args[1:]); ok {
		runHelp(topic)
		return
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

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
	case "habit", "habits":
		runHabit(db, os.Args[2:])
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
		// atalho: neivi 211 = neivi show 211 (RN-26)
		if _, err := strconv.ParseInt(os.Args[1], 10, 64); err == nil {
			runShow(db, os.Args[1:])
			return
		}
		fmt.Fprintf(os.Stderr, "erro: comando %s desconhecido (ajuda: neivi -help)\n", os.Args[1])
		os.Exit(1)
	}
}
