package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CaioParanaiba/neivi/internal/habit"
	"github.com/CaioParanaiba/neivi/internal/store"
	"github.com/CaioParanaiba/neivi/internal/todotxt"
)

func runHabit(db *sql.DB, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "erro: falta o comando (exemplo: neivi habit do cf; ou add, undo, ls, edit, rm)")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		runHabitAdd(db, args[1:])
	case "do":
		runHabitDo(db, args[1:], false)
	case "undo":
		runHabitDo(db, args[1:], true)
	case "ls":
		runHabitLs(db, args[1:])
	case "edit":
		runHabitEdit(db, args[1:])
	case "rm":
		runHabitRm(db, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "erro: comando habit %s desconhecido (use add, do, undo, ls, edit ou rm)\n", args[0])
		os.Exit(1)
	}
}

// doneIn soma o que foi feito no periodo da meta do habito (o dia, ou a semana)
func doneIn(db *sql.DB, h habit.Habit, day time.Time) int {
	from, to := h.Period(day)
	n, err := store.SumLogs(db, h.ID, from, to)
	if err != nil {
		fail(err)
	}
	return n
}

// applyPeriod aplica meta:, sem: e dias: no habito (add e edit)
func applyPeriod(h *habit.Habit, e todotxt.HabitEntry) error {
	switch {
	case e.Weekly > 0:
		h.Kind = habit.Weekly
		h.Goal = e.Weekly
		h.Days = nil
	case len(e.Days) > 0:
		h.Kind = habit.Weekdays
		h.Days = e.Days
		if h.Goal == 0 {
			h.Goal = 1
		}
	}
	if e.Goal > 0 {
		if h.Kind == habit.Weekly {
			return fmt.Errorf("no sem: o numero ja eh a meta da semana (use so sem:%d)", h.Goal)
		}
		h.Goal = e.Goal
	}
	return nil
}

func runHabitAdd(db *sql.DB, args []string) {
	e, err := todotxt.ParseHabit(args, false)
	if err != nil {
		fail(err)
	}

	_, err = store.GetHabit(db, e.Code)
	if err == nil {
		fail(fmt.Errorf("ja existe um habito %s (para mudar: neivi habit edit %s)", e.Code, e.Code))
	}

	h := habit.Habit{Code: e.Code, Name: e.Name, Kind: habit.Daily, Goal: 1, Status: "active"}
	if h.Name == "" {
		h.Name = h.Code
	}
	if e.Status != "" {
		h.Status = e.Status
	}
	err = applyPeriod(&h, e)
	if err != nil {
		fail(err)
	}

	id, err := store.AddHabit(db, h)
	if err != nil {
		fail(err)
	}

	fmt.Println(confirmEntity("+", id, h.Code, h.Name, h.MetaLabel()))
}

// missingLine eh o fim do habit do: "faltam hoje: digi, cf (1)" ou "habitos: tudo feito hoje ✓"
func missingLine(db *sql.DB, today time.Time) string {
	lista, err := store.ListHabits(db, false)
	if err != nil {
		fail(err)
	}

	var missing []string
	for _, h := range lista {
		done := doneIn(db, h, today)
		if !h.Missing(done, today) {
			continue
		}
		name := h.Code
		if h.Goal > 1 {
			name += fmt.Sprintf(" (%d)", h.Goal-done)
		}
		missing = append(missing, name)
	}

	if len(missing) == 0 {
		return "habitos: tudo feito hoje ✓"
	}
	return "faltam hoje: " + strings.Join(missing, ", ")
}

// runHabitDo eh o habit do e, com undo = true, o habit undo (RN-24)
func runHabitDo(db *sql.DB, args []string, undo bool) {
	cmd := "do"
	if undo {
		cmd = "undo"
	}

	refs, ontem, err := todotxt.ParseHabitRefs(args, cmd)
	if err != nil {
		fail(err)
	}

	now := time.Now()
	day := now
	suffix := ""
	when := "hoje"
	if ontem {
		day = now.AddDate(0, 0, -1)
		suffix = " (ontem)"
		when = "ontem"
	}
	date := day.Format("2006-01-02")

	// le todos antes de agir (RN-12): codigo que nao existe recusa tudo
	type target struct {
		h   habit.Habit
		qty int
	}
	var targets []target
	for _, r := range refs {
		if r.Code != "" {
			h, err := store.GetHabit(db, r.Code)
			if errors.Is(err, sql.ErrNoRows) {
				fail(fmt.Errorf("%q nao eh um habito, quantidade (x2) nem \"ontem\" (veja os habitos no neivi habit ls)", r.Code))
			}
			if err != nil {
				fail(err)
			}
			targets = append(targets, target{h, r.Qty})
			continue
		}

		h, err := store.GetHabitByID(db, r.ID)
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Fprintf(os.Stderr, "aviso: habito %d nao existe (veja os habitos no neivi habit ls)\n", r.ID)
			continue
		}
		if err != nil {
			fail(err)
		}
		targets = append(targets, target{h, r.Qty})
	}

	for _, t := range targets {
		h := t.h
		if undo {
			had, err := store.UndoLog(db, h.ID, date, t.qty)
			if err != nil {
				fail(err)
			}
			if had == 0 {
				fmt.Fprintf(os.Stderr, "aviso: %s nao tem registro %s\n", h.Code, when)
				continue
			}
			fmt.Printf("↺ %s · %s%s\n", h.Code, h.Confirm(doneIn(db, h, day), ontem), suffix)
			if had < t.qty {
				fmt.Fprintf(os.Stderr, "aviso: %s so tinha %d %s\n", h.Code, had, when)
			}
			continue
		}

		err := store.AddLog(db, h.ID, date, t.qty)
		if err != nil {
			fail(err)
		}
		fmt.Printf("✓ %s · %s%s\n", h.Code, h.Confirm(doneIn(db, h, day), ontem), suffix)

		if !h.IsDay(day) {
			fmt.Fprintf(os.Stderr, "aviso: %s nao eh de %s\n", h.Code, habit.DayName(day))
		}
		if h.Status == "archived" {
			fmt.Fprintf(os.Stderr, "aviso: o habito %s esta arquivado (reativar: neivi habit edit %s sit:ativo)\n", h.Code, h.Code)
		}
	}

	if len(targets) > 0 {
		fmt.Println(missingLine(db, now))
	}
}

func runHabitLs(db *sql.DB, args []string) {
	all := len(args) == 1 && args[0] == "all"
	if len(args) > 0 && !all {
		fail(fmt.Errorf("habit ls so aceita all (exemplo: neivi habit ls all)"))
	}

	lista, err := store.ListHabits(db, all)
	if err != nil {
		fail(err)
	}

	if len(lista) == 0 {
		fmt.Println("nenhum habito (cadastrar: neivi habit add codigo nome)")
		return
	}

	now := time.Now()
	head := []string{"ID", "HABITO", "META", "FEITO", "NOME"}
	if all {
		head = []string{"ID", "HABITO", "META", "FEITO", "SITUACAO", "NOME"}
	}

	var tb table
	tb.fit(head)
	var rows [][]string
	for _, h := range lista {
		r := []string{strconv.FormatInt(h.ID, 10), h.Code, h.MetaLabel(), h.DoneLabel(doneIn(db, h, now), now)}
		if all {
			r = append(r, habit.StatusPT(h.Status))
		}
		r = append(r, h.Name)
		tb.fit(r)
		rows = append(rows, r)
	}

	fmt.Println(tb.line(head))
	for _, r := range rows {
		fmt.Println(tb.line(r))
	}
}

// printHabitsToday eh o bloco de habitos no fim do today (US-19, CA-19.1)
func printHabitsToday(db *sql.DB, today time.Time) {
	lista, err := store.ListHabits(db, false)
	if err != nil {
		fail(err)
	}
	if len(lista) == 0 {
		return // sem habitos cadastrados, nada aparece
	}

	head := []string{"ID", "HABITO", "FEITO"}
	var tb table
	tb.fit(head)
	var rows [][]string
	for _, h := range lista {
		done := doneIn(db, h, today)
		if !h.Missing(done, today) {
			continue
		}
		r := []string{strconv.FormatInt(h.ID, 10), h.Code, h.DoneLabel(done, today)}
		tb.fit(r)
		rows = append(rows, r)
	}

	fmt.Println()
	if len(rows) == 0 {
		fmt.Println("habitos: tudo feito hoje ✓")
		return
	}
	fmt.Println(paint(stSection, fmt.Sprintf("habitos (%d)", len(rows))))
	fmt.Println(paint(stDim, tb.line(head)))
	for _, r := range rows {
		fmt.Println(tb.paintedLine(r, []string{paint(stDim, r[0]), r[1], r[2]}))
	}
}

// getHabit busca pelo codigo ou pelo id e sai com a mensagem certa se nao existir
func getHabit(db *sql.DB, ref string) habit.Habit {
	var h habit.Habit
	var err error
	if id, errID := strconv.ParseInt(ref, 10, 64); errID == nil {
		h, err = store.GetHabitByID(db, id)
	} else {
		h, err = store.GetHabit(db, ref)
	}
	if errors.Is(err, sql.ErrNoRows) {
		fail(fmt.Errorf("habito %s nao existe (veja os habitos no neivi habit ls all)", ref))
	}
	if err != nil {
		fail(err)
	}
	return h
}

func runHabitEdit(db *sql.DB, args []string) {
	e, err := todotxt.ParseHabit(args, true)
	if err != nil {
		fail(fmt.Errorf("%s", strings.Replace(err.Error(), "habit add digi Treino de digitacao", "habit edit cf meta:5", 1)))
	}

	h := getHabit(db, e.Code)

	if e.IsEmpty() {
		fail(fmt.Errorf("nada para mudar no habito %s (exemplo: neivi habit edit %s meta:5)", h.Code, h.Code))
	}

	if e.NewCode != "" && e.NewCode != h.Code {
		_, err := store.GetHabit(db, e.NewCode)
		if err == nil {
			fail(fmt.Errorf("ja existe um habito %s", e.NewCode))
		}
		h.Code = e.NewCode
	}
	if e.Name != "" {
		h.Name = e.Name
	}
	if e.Status != "" {
		h.Status = e.Status
	}
	err = applyPeriod(&h, e)
	if err != nil {
		fail(err)
	}

	err = store.UpdateHabit(db, h)
	if err != nil {
		fail(err)
	}

	status := h.MetaLabel()
	if h.Status == "archived" {
		status = "arquivado"
	}
	fmt.Println(confirmEntity("~", h.ID, h.Code, h.Name, status))
}

func runHabitRm(db *sql.DB, args []string) {
	if len(args) != 1 {
		fail(fmt.Errorf("falta o habito (exemplo: neivi habit rm x)"))
	}

	h := getHabit(db, args[0])

	// so sem registros (RN-22)
	n, err := store.CountLogs(db, h.ID)
	if err != nil {
		fail(err)
	}
	if n > 0 {
		fail(fmt.Errorf("o habito %s tem %d registros (para parar sem perder o historico: neivi habit edit %s sit:arquivado)", h.Code, n, h.Code))
	}

	if !askYes(fmt.Sprintf("remover o habito %s (sem registros)?", h.Code)) {
		fmt.Printf("%s mantido\n", h.Code)
		return
	}

	err = store.RmHabit(db, h.ID)
	if err != nil {
		fail(err)
	}
	fmt.Printf("✗ %s removido\n", h.Code)
}
