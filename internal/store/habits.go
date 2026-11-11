package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/CaioParanaiba/neivi/internal/habit"
)

const selectHabits = `SELECT id, code, name, goal, freq_kind, COALESCE(freq_value, ''), status, created_at
        FROM habits`

func scanHabit(sc scanner) (habit.Habit, error) {
	var h habit.Habit
	var days string
	err := sc.Scan(&h.ID, &h.Code, &h.Name, &h.Goal, &h.Kind, &days, &h.Status, &h.CreatedAt)
	if days != "" {
		h.Days = habit.DaysFromValue(days)
	}
	return h, err
}

// AddHabit cadastra um habito e devolve o id criado
func AddHabit(db *sql.DB, h habit.Habit) (int64, error) {
	if h.Status == "" {
		h.Status = "active"
	}

	res, err := db.Exec(`INSERT INTO habits (code, name, goal, freq_kind, freq_value, status)
        VALUES (?, ?, ?, ?, ?, ?)`,
		h.Code, h.Name, h.Goal, h.Kind, nullString(habit.DaysValue(h.Days)), h.Status)
	if err != nil {
		return 0, fmt.Errorf("inserir habito: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("retornar id: %w", err)
	}

	return id, nil
}

// GetHabit busca pelo codigo; nao existe = sql.ErrNoRows
func GetHabit(db *sql.DB, code string) (habit.Habit, error) {
	return scanHabit(db.QueryRow(selectHabits+" WHERE code = ?", code))
}

// GetHabitByID busca pelo id; nao existe = sql.ErrNoRows
func GetHabitByID(db *sql.DB, id int64) (habit.Habit, error) {
	return scanHabit(db.QueryRow(selectHabits+" WHERE id = ?", id))
}

// ListHabits devolve os ativos, ou todos (all) com os arquivados depois
func ListHabits(db *sql.DB, all bool) ([]habit.Habit, error) {
	query := selectHabits
	if !all {
		query += " WHERE status = 'active'"
	}
	query += " ORDER BY status = 'archived', id"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar habitos: %w", err)
	}

	defer rows.Close()

	var lista []habit.Habit
	for rows.Next() {
		h, err := scanHabit(rows)
		if err != nil {
			return nil, fmt.Errorf("ler habito: %w", err)
		}
		lista = append(lista, h)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer habitos: %w", err)
	}

	return lista, nil
}

// UpdateHabit grava todos os campos editaveis
func UpdateHabit(db *sql.DB, h habit.Habit) error {
	_, err := db.Exec(`UPDATE habits SET code = ?, name = ?, goal = ?, freq_kind = ?, freq_value = ?, status = ?
        WHERE id = ?`,
		h.Code, h.Name, h.Goal, h.Kind, nullString(habit.DaysValue(h.Days)), h.Status, h.ID)
	if err != nil {
		return fmt.Errorf("atualizar habito %s: %w", h.Code, err)
	}
	return nil
}

// AddLog soma a quantidade no registro do dia (RN-24). Se o dia ainda nao tem
// registro, cria; se ja tem, soma no que existe (o ON CONFLICT da chave primaria).
func AddLog(db *sql.DB, habitID int64, date string, qty int) error {
	_, err := db.Exec(`INSERT INTO habit_logs (habit_id, date, value) VALUES (?, ?, ?)
        ON CONFLICT (habit_id, date) DO UPDATE SET value = value + excluded.value`,
		habitID, date, qty)
	if err != nil {
		return fmt.Errorf("registrar habito: %w", err)
	}
	return nil
}

// UndoLog tira a quantidade do registro do dia, sem deixar negativo (RN-24).
// Devolve quanto tinha antes; chegando a zero, o registro some.
func UndoLog(db *sql.DB, habitID int64, date string, qty int) (int, error) {
	var had int
	err := db.QueryRow("SELECT value FROM habit_logs WHERE habit_id = ? AND date = ?", habitID, date).Scan(&had)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("ler registro: %w", err)
	}

	if had <= qty {
		_, err = db.Exec("DELETE FROM habit_logs WHERE habit_id = ? AND date = ?", habitID, date)
	} else {
		_, err = db.Exec("UPDATE habit_logs SET value = value - ? WHERE habit_id = ? AND date = ?", qty, habitID, date)
	}
	if err != nil {
		return 0, fmt.Errorf("desfazer registro: %w", err)
	}

	return had, nil
}

// SumLogs soma os registros entre duas datas ISO (o dia, ou a semana)
func SumLogs(db *sql.DB, habitID int64, from, to string) (int, error) {
	var n int
	err := db.QueryRow("SELECT COALESCE(SUM(value), 0) FROM habit_logs WHERE habit_id = ? AND date BETWEEN ? AND ?",
		habitID, from, to).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("somar registros: %w", err)
	}
	return n, nil
}

// CountLogs conta os dias com registro (para o rm, RN-22)
func CountLogs(db *sql.DB, habitID int64) (int, error) {
	var n int
	err := db.QueryRow("SELECT COUNT(*) FROM habit_logs WHERE habit_id = ?", habitID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("contar registros: %w", err)
	}
	return n, nil
}

// RmHabit apaga o habito (o cli confere antes que nao tem registros)
func RmHabit(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM habits WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deletar habito %d: %w", id, err)
	}
	return nil
}
