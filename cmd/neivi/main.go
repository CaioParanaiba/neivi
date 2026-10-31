package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		runAdd(os.Args[2:])
	case "ls":
		runLs(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, "Comando", os.Args[1], "eh desconhecido")
		usage()
		os.Exit(1)
	}
}
