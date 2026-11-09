// Package task define os tipos do domínio de tarefas
package task

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
}
