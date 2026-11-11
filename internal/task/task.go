// Package task define os tipos do domínio de tarefas
package task

import "strings"

type Task struct {
	ID          int64
	Title       string
	Description string
	Priority    string
	Due         string
	EstimateMin int64
	SubjectID   int64
	ProjectID   int64
	DoneAt      string
	CreatedAt   string
	SubjectCode string
	ProjectCode string
	SubjectName string
	ProjectName string

	// SubjectInherited: a disciplina veio do projeto, a tarefa nao tem @ proprio (RN-14)
	SubjectInherited bool
}

// NormalizeTitle deixa o titulo pronto para comparar duplicadas (RN-01):
// minusculo e sem espacos extras
func NormalizeTitle(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}
