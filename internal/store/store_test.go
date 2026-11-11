package store

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

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
