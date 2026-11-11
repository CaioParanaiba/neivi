package store

import (
	"testing"

	"github.com/CaioParanaiba/neivi/internal/habit"
)

func TestHabitLogs(t *testing.T) {
	db := openTest(t)

	cf, err := AddHabit(db, habit.Habit{Code: "cf", Name: "cf", Goal: 3, Kind: habit.Daily})
	if err != nil {
		t.Fatal(err)
	}

	sum := func(from, to string) int {
		t.Helper()
		n, err := SumLogs(db, cf, from, to)
		if err != nil {
			t.Fatal(err)
		}
		return n
	}

	t.Run("CA-18.1 soma no registro do dia", func(t *testing.T) {
		_ = AddLog(db, cf, "2026-11-11", 2)
		_ = AddLog(db, cf, "2026-11-11", 1)
		if got := sum("2026-11-11", "2026-11-11"); got != 3 {
			t.Errorf("dia = %d, queria 3 (2 + 1 no mesmo registro)", got)
		}
		if n, _ := CountLogs(db, cf); n != 1 {
			t.Errorf("registros = %d, queria 1 por dia", n)
		}
	})

	t.Run("CA-19.3 soma da semana", func(t *testing.T) {
		_ = AddLog(db, cf, "2026-11-09", 1)
		_ = AddLog(db, cf, "2026-11-07", 5) // sabado da semana anterior
		if got := sum("2026-11-08", "2026-11-14"); got != 4 {
			t.Errorf("semana = %d, queria 4", got)
		}
	})

	t.Run("CA-20.1 undo tira uma parte", func(t *testing.T) {
		had, _ := UndoLog(db, cf, "2026-11-11", 2)
		if had != 3 || sum("2026-11-11", "2026-11-11") != 1 {
			t.Errorf("tinha %d, ficou %d; queria 3 e 1", had, sum("2026-11-11", "2026-11-11"))
		}
	})

	t.Run("CA-20.2 undo nunca fica negativo", func(t *testing.T) {
		had, _ := UndoLog(db, cf, "2026-11-11", 9)
		if had != 1 || sum("2026-11-11", "2026-11-11") != 0 {
			t.Errorf("tinha %d, ficou %d; queria 1 e 0", had, sum("2026-11-11", "2026-11-11"))
		}
	})

	t.Run("CA-20.3 undo sem registro", func(t *testing.T) {
		had, err := UndoLog(db, cf, "2026-11-11", 1)
		if err != nil || had != 0 {
			t.Errorf("tinha %d (%v), queria 0", had, err)
		}
	})

	t.Run("CA-19.5 arquivado so no all", func(t *testing.T) {
		_, _ = AddHabit(db, habit.Habit{Code: "digi", Name: "digi", Goal: 1, Kind: habit.Daily, Status: "archived"})
		normal, _ := ListHabits(db, false)
		all, _ := ListHabits(db, true)
		if len(normal) != 1 || len(all) != 2 || all[1].Code != "digi" {
			t.Errorf("ativos = %d, todos = %v", len(normal), all)
		}
	})
}
