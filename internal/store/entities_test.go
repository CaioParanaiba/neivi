package store

import (
	"testing"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/subject"
	"github.com/CaioParanaiba/neivi/internal/task"
)

func TestProjectsAndSubjects(t *testing.T) {
	db := openTest(t)

	bd, _ := AddSubject(db, subject.Subject{Code: "BD", Name: "BD"})
	ed, _ := AddSubject(db, subject.Subject{Code: "ED", Name: "ED", Status: "closed"})
	fb, _ := AddProject(db, project.Project{Code: "feedback360", Name: "Trabalho", SubjectID: bd})
	ma, _ := AddProject(db, project.Project{Code: "maratona", Name: "maratona", Continuous: true})
	_, _ = AddProject(db, project.Project{Code: "velho", Name: "velho", Status: "canceled"})

	add := func(tk task.Task) int64 {
		t.Helper()
		id, err := AddTask(db, tk)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	herdada := add(task.Task{Title: "estudar joins", ProjectID: fb})                  // disciplina BD, do projeto
	propria := add(task.Task{Title: "estudar arvores", ProjectID: fb, SubjectID: ed}) // @ED proprio
	feita := add(task.Task{Title: "entregar parte 1", ProjectID: fb})
	_ = add(task.Task{Title: "treino", ProjectID: ma})
	if err := DoneTask(db, feita, 0); err != nil {
		t.Fatal(err)
	}

	t.Run("CA-12.2 progresso contado no banco", func(t *testing.T) {
		p, err := GetProject(db, "feedback360")
		if err != nil {
			t.Fatal(err)
		}
		if p.Done != 1 || p.Pending != 2 || p.LastDoneAt == "" || p.SubjectCode != "BD" {
			t.Errorf("projeto = %+v, queria 1 feita, 2 pendentes, ultima preenchida, disciplina BD", p)
		}
	})

	t.Run("CA-12.1 e CA-12.3 lista sem e com os cancelados", func(t *testing.T) {
		normal, _ := ListProjects(db, false)
		all, _ := ListProjects(db, true)
		if len(normal) != 2 || len(all) != 3 {
			t.Errorf("ativos/pausados = %d, todos = %d; queria 2 e 3", len(normal), len(all))
		}
	})

	t.Run("CA-13.1 ultima tarefa feita", func(t *testing.T) {
		id, title, _, err := LastDone(db, fb)
		if err != nil || id != feita || title != "entregar parte 1" {
			t.Errorf("LastDone = %d %q %v; queria %d", id, title, err, feita)
		}
	})

	t.Run("CA-10.5 RN-14 disciplina herdada no filtro", func(t *testing.T) {
		lista, _ := ListTasksFilter(db, task.Filter{Subjects: []string{"BD"}})
		if len(lista) != 1 || lista[0].ID != herdada || !lista[0].SubjectInherited {
			t.Errorf("@BD = %+v, queria so a herdada", lista)
		}

		lista, _ = ListTasksFilter(db, task.Filter{Subjects: []string{"ED"}})
		if len(lista) != 1 || lista[0].ID != propria || lista[0].SubjectInherited {
			t.Errorf("@ED = %+v, queria so a com @ED proprio", lista)
		}

		lista, _ = ListTasksFilter(db, task.Filter{NoSubject: true})
		if len(lista) != 1 {
			t.Errorf("@- = %d tarefas, queria 1 (a da maratona)", len(lista))
		}
	})

	t.Run("CA-15.2 pendentes pela disciplina efetiva e encerradas so no all", func(t *testing.T) {
		normal, _ := ListSubjects(db, false)
		if len(normal) != 1 || normal[0].Code != "BD" || normal[0].Pending != 1 {
			t.Errorf("ativas = %+v, queria so BD com 1 pendente", normal)
		}
		all, _ := ListSubjects(db, true)
		if len(all) != 2 || all[1].Code != "ED" {
			t.Errorf("todas = %+v, queria BD e depois ED", all)
		}
	})

	t.Run("CA-16.2 e CA-16.3 o que esta ligado", func(t *testing.T) {
		tasks, projects, err := SubjectLinks(db, bd)
		if err != nil || tasks != 0 || len(projects) != 1 || projects[0] != "feedback360" {
			t.Errorf("BD: %d tarefas, projetos %v, %v; queria 0 e [feedback360]", tasks, projects, err)
		}
		tasks, _, _ = SubjectLinks(db, ed)
		if tasks != 1 {
			t.Errorf("ED: %d tarefas, queria 1 (so as com @ED proprio)", tasks)
		}
	})
}
