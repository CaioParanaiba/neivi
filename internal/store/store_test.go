package store

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/subject"
	"github.com/CaioParanaiba/neivi/internal/task"
)

// openTest abre um banco novo numa pasta temporaria (apagada no fim do teste)
func openTest(t *testing.T) *sql.DB {
	t.Helper()

	db, err := Open(filepath.Join(t.TempDir(), "teste.db"))
	if err != nil {
		t.Fatalf("abrir banco: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

// addTest cria uma tarefa so com o titulo e devolve o id
func addTest(t *testing.T, db *sql.DB, title string) int64 {
	t.Helper()

	id, err := AddTask(db, task.Task{Title: title})
	if err != nil {
		t.Fatalf("AddTask(%q): %v", title, err)
	}

	return id
}

func TestFindDuplicate(t *testing.T) {
	db := openTest(t)

	pendente := addTest(t, db, "comprar cabo HDMI") // 1
	feita := addTest(t, db, "ligar")                // 2
	err := DoneTask(db, feita, 0)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		title    string
		exceptID int64
		want     int64
	}{
		{"CA-01.4 caixa e espaco diferentes", "Comprar  cabo hdmi", 0, pendente},
		{"CA-01.5 concluida nao conta", "ligar", 0, 0},
		{"titulo novo", "revisar slide", 0, 0},
		{"CA-03.7 edit acha outra tarefa", "comprar cabo hdmi", 99, pendente},
		{"edit nao conta a propria tarefa", "comprar cabo hdmi", pendente, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FindDuplicate(db, tc.title, tc.exceptID)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("FindDuplicate(%q, %d) = %d, queria %d", tc.title, tc.exceptID, got, tc.want)
			}
		})
	}
}

func TestDoneTaskECountDone(t *testing.T) {
	db := openTest(t)

	a := addTest(t, db, "a")
	b := addTest(t, db, "b")
	c := addTest(t, db, "c")

	// CA-04.1: duas hoje
	for _, id := range []int64{a, b} {
		err := DoneTask(db, id, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	// CA-04.2 e CA-04.3: uma ontem (criada hoje, concluida antes da criacao)
	err := DoneTask(db, c, 1)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		daysAgo int
		want    int
	}{
		{"CA-04.1 contador de hoje", 0, 2},
		{"CA-04.2 contador de ontem", 1, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CountDone(db, tc.daysAgo)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("CountDone(%d) = %d, queria %d", tc.daysAgo, got, tc.want)
			}
		})
	}

	t.Run("CA-04.3 data de ontem", func(t *testing.T) {
		v, err := GetTask(db, c)
		if err != nil {
			t.Fatal(err)
		}

		ontem := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if v.DoneAt[:10] != ontem {
			t.Errorf("done_at = %s, queria o dia %s", v.DoneAt, ontem)
		}
	})
}

func TestReopenTask(t *testing.T) {
	db := openTest(t)

	id := addTest(t, db, "comprar cabo HDMI")
	err := DoneTask(db, id, 0)
	if err != nil {
		t.Fatal(err)
	}

	// CA-05.1: volta a pendente, sem data de conclusao
	err = ReopenTask(db, id)
	if err != nil {
		t.Fatal(err)
	}

	v, err := GetTask(db, id)
	if err != nil {
		t.Fatal(err)
	}
	if v.DoneAt != "" {
		t.Errorf("CA-05.1 done_at = %q, queria vazio", v.DoneAt)
	}
	if v.Title != "comprar cabo HDMI" {
		t.Errorf("CA-05.1 titulo = %q, queria igual ao de antes", v.Title)
	}
}

func TestRmTask(t *testing.T) {
	db := openTest(t)

	id := addTest(t, db, "entregar lista")

	// CA-06.1: a tarefa some
	err := RmTask(db, id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = GetTask(db, id)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("CA-06.1 GetTask depois do rm: erro %v, queria sql.ErrNoRows", err)
	}
}

func TestSubjectIDByCode(t *testing.T) {
	db := openTest(t)

	_, err := SubjectIDByCode(db, "XX")
	want := "disciplina XX nao cadastrada (cadastrar: neivi subject add XX)"
	if err == nil || err.Error() != want {
		t.Errorf("CA-02.6 erro = %v, queria %q", err, want)
	}
}

func TestListTasksFilter(t *testing.T) {
	db := openTest(t)

	bd, _ := AddSubject(db, subject.Subject{Code: "BD", Name: "BD"})
	ed, _ := AddSubject(db, subject.Subject{Code: "ED", Name: "ED"})
	neivi, _ := AddProject(db, project.Project{Code: "neivi", Name: "neivi"})

	add := func(tk task.Task) int64 {
		t.Helper()
		id, err := AddTask(db, tk)
		if err != nil {
			t.Fatal(err)
		}
		return id
	}

	// as tarefas do CA-09.1 (os ids sao outros, mas o papel eh o mesmo)
	t3 := add(task.Task{Title: "entregar parte do feedback360", Priority: "A", SubjectID: bd})
	t2 := add(task.Task{Title: "revisar normalizacao", Priority: "A", SubjectID: bd})
	t11 := add(task.Task{Title: "parser", Priority: "B", ProjectID: neivi})
	t9 := add(task.Task{Title: "estudar sorting", Priority: "B", SubjectID: ed})
	t4 := add(task.Task{Title: "comprar cabo HDMI"})
	t1 := add(task.Task{Title: "ligar"})
	if err := DoneTask(db, t1, 0); err != nil {
		t.Fatal(err)
	}
	// a 9 foi criada ha uma semana (nao eh nova)
	if _, err := db.Exec("UPDATE tasks SET created_at = datetime('now', 'localtime', '-7 days') WHERE id = ?", t9); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		f    task.Filter
		want []int64
	}{
		{"CA-09.1 so pendentes", task.Filter{}, []int64{t3, t2, t11, t9, t4}},
		{"CA-09.2 @BD", task.Filter{Subjects: []string{"BD"}}, []int64{t3, t2}},
		{"CA-09.2 @", task.Filter{AnySubject: true}, []int64{t3, t2, t9}},
		{"CA-09.2 @-", task.Filter{NoSubject: true}, []int64{t11, t4}},
		{"CA-09.2 +neivi", task.Filter{Projects: []string{"neivi"}}, []int64{t11}},
		{"CA-09.2 +-", task.Filter{NoProject: true}, []int64{t3, t2, t9, t4}},
		{"CA-09.2 pri:A", task.Filter{Priorities: []string{"A"}}, []int64{t3, t2}},
		{"CA-09.2 pri:- (triage)", task.Filter{Priorities: []string{"-"}}, []int64{t4}},
		{"CA-09.2 new", task.Filter{New: true}, []int64{t3, t2, t11, t4}},
		{"CA-09.3 @BD @ED eh ou", task.Filter{Subjects: []string{"BD", "ED"}}, []int64{t3, t2, t9}},
		{"CA-09.3 @BD pri:B eh e", task.Filter{Subjects: []string{"BD"}, Priorities: []string{"B"}}, nil},
		{"CA-09.4 busca sem acento", task.Filter{Words: []string{"Normalização"}}, []int64{t2}},
		{"CA-09.5 all inclui concluidas", task.Filter{All: true}, []int64{t3, t2, t11, t9, t4, t1}},
		{"CA-09.5 done so concluidas", task.Filter{Done: true}, []int64{t1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lista, err := ListTasksFilter(db, tc.f)
			if err != nil {
				t.Fatal(err)
			}
			var got []int64
			for _, v := range lista {
				got = append(got, v.ID)
			}
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("ids = %v, queria %v", got, tc.want)
			}
		})
	}
}
