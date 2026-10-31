package main

import (
	"fmt"
	"os"
)

const help = `uso: neivi <comando> [argumentos]

comandos:
	add	adicionar tarefa
	ls	listar tarefas
`

func usage() {
	fmt.Fprint(os.Stderr, help)
}

func runAdd(args []string) {
	fmt.Println("Foi adicionado a tarefa:", args)
}

func runLs(args []string) {
	fmt.Println("Rodando o comando associado a ls:", args)
}
