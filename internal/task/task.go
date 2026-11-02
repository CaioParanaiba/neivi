// Package task define os tipos do domínio de tarefas
package task

type Task struct {
	ID        int64
	Title     string
	Priority  string
	Due       string
	SubjectId int64
	ProjectId int64
	ParentId  int64
	Extras    string
	doneAt    string
	CreatedAt string
}
