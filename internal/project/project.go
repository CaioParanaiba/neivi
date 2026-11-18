// Package project define o tipo dos projetos
package project

import (
	"fmt"
	"strconv"
	"strings"
)

// Reserved sao os subcomandos do project: o atalho neivi project <codigo> nao
// distinguiria um projeto com esses nomes (RN-26)
var Reserved = []string{"add", "ls", "show", "edit", "rm"}

// ValidCode confere o codigo (RN-26): nunca so numero (se confundiria com o id),
// nunca nome de subcomando
func ValidCode(code string) error {
	if _, err := strconv.Atoi(code); err == nil {
		return fmt.Errorf("o codigo %s nao pode ser so um numero (ele se confundiria com o id)", code)
	}
	for _, r := range Reserved {
		if code == r {
			return fmt.Errorf("o codigo %s eh um comando do project (%s); escolha outro", code, strings.Join(Reserved, ", "))
		}
	}
	return nil
}

type Project struct {
	ID          int64
	Code        string
	Name        string
	Description string
	GitHub      string
	Dir         string
	Status      string // active, paused, done ou canceled (no banco, em ingles)
	Continuous  bool
	SubjectID   int64
	SubjectCode string
	SubjectName string
	CreatedAt   string

	// deduzidos das tarefas (RN-20), nunca guardados
	Done       int
	Pending    int
	LastDoneAt string
}

var statusPT = map[string]string{
	"active":   "ativo",
	"paused":   "pausado",
	"done":     "concluido",
	"canceled": "cancelado",
}

// StatusPT mostra a situacao em portugues
func StatusPT(status string) string {
	return statusPT[status]
}

// StatusFromPT le a situacao digitada (sit:pausado...)
func StatusFromPT(s string) (string, error) {
	for en, pt := range statusPT {
		if s == pt {
			return en, nil
		}
	}
	return "", fmt.Errorf("situacao %s invalida (use ativo, pausado, concluido ou cancelado)", s)
}

// StatusOrder eh a ordem dos grupos no project ls: ativos, pausados, concluidos, cancelados
func StatusOrder(status string) int {
	switch status {
	case "active":
		return 0
	case "paused":
		return 1
	case "done":
		return 2
	}
	return 3
}

// Progress eh o progresso em % (RN-20); ok = false quando nao faz sentido
// (projeto continuo ou sem tarefas)
func (p Project) Progress() (int, bool) {
	total := p.Done + p.Pending
	if p.Continuous || total == 0 {
		return 0, false
	}
	return (p.Done*100 + total/2) / total, true // arredonda para o inteiro mais perto
}
