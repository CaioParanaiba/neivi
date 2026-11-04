// Package store guarda e lê os dados do neivi no SQLite.
package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/subject"
	"github.com/CaioParanaiba/neivi/internal/task"
)

//const path_sql = "~/Code/neivi/migrations/001_init.sql"

//go:embed migrations/001_init.sql
var schema string

// scaner de task, para facilitar e refatorar de forma facil
type scanner interface {
	Scan(dest ...any) error
}

const selectTasks = `SELECT tasks.id, tasks.title, COALESCE(tasks.description, ''),
        COALESCE(tasks.priority, ''), COALESCE(tasks.due, ''),
        COALESCE(tasks.estimate_min, 0), COALESCE(tasks.subject_id, 0),
        COALESCE(tasks.project_id, 0), COALESCE(tasks.parent_id, 0),
        COALESCE(tasks.extras, ''), COALESCE(tasks.done_at, ''),
        tasks.created_at, COALESCE(subjects.code, ''),
        COALESCE(projects.code, '')
        FROM tasks
        LEFT JOIN subjects ON subjects.id = tasks.subject_id
        LEFT JOIN projects ON projects.id = tasks.project_id`

func scanTask(s scanner) (task.Task, error) {
	var v task.Task
	err := s.Scan(&v.ID, &v.Title, &v.Description,
		&v.Priority, &v.Due, &v.EstimateMin,
		&v.SubjectID, &v.ProjectID, &v.ParentID,
		&v.Extras, &v.DoneAt, &v.CreatedAt,
		&v.SubjectCode, &v.ProjectCode)
	return v, err
}

func nullInt(n int64) any {
	if n == 0 {
		return nil
	}
	return n
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}
	_, err = db.Exec(schema)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("criar tabelas: %w", err)
	}
	return db, nil
}

func AddTask(db *sql.DB, t task.Task) (int64, error) {
	res, err := db.Exec(`INSERT INTO tasks (title, description, priority, due, estimate_min, subject_id, project_id, extras)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		t.Title,
		nullString(t.Description),
		nullString(t.Priority),
		nullString(t.Due),
		nullInt(t.EstimateMin),
		nullInt(t.SubjectID),
		nullInt(t.ProjectID),
		nullString(t.Extras))
	if err != nil {
		return 0, fmt.Errorf("inserir tarefa: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("retornar id: %w", err)
	}

	return id, nil
}

func ListTasks(db *sql.DB) ([]task.Task, error) {
	rows, err := db.Query(selectTasks + " ORDER BY tasks.id")

	if err != nil {
		return nil, fmt.Errorf("listar tarefas: %w", err)
	}

	defer rows.Close()

	var lista []task.Task
	for rows.Next() {
		v, err := scanTask(rows) // rows é *sql.Rows → serve como scanner
		if err != nil {
			return nil, fmt.Errorf("ler tarefa: %w", err)
		}
		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer tarefas: %w", err)
	}

	return lista, nil
}

func UpdateTask(db *sql.DB, t task.Task) error {
	_, err := db.Exec(`UPDATE tasks SET title = ?, description = ?, priority = ?, due = ?,
		estimate_min = ?, subject_id = ?, project_id = ?, extras = ? WHERE id = ?`,
		t.Title,
		nullString(t.Description),
		nullString(t.Priority),
		nullString(t.Due),
		nullInt(t.EstimateMin),
		nullInt(t.SubjectID),
		nullInt(t.ProjectID),
		nullString(t.Extras),
		t.ID)
	if err != nil {
		return fmt.Errorf("atualizar tarefa %d: %w", t.ID, err)
	}

	return nil
}

func SubjectIDByCode(db *sql.DB, code string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM subjects WHERE code = ?", code).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("disciplina %s nao cadastrada: %w", code, err)
	}
	if err != nil {
		return 0, fmt.Errorf("erro ao consultar disciplinas: %w", err)
	}

	return id, nil
}
func ProjectIDByCode(db *sql.DB, code string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM projects WHERE code = ?", code).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("projeto %s nao cadastrada: %w", code, err)
	}
	if err != nil {
		return 0, fmt.Errorf("erro ao consultar projetos: %w", err)
	}

	return id, nil
}

// AddSubject cadastra uma disciplina e devolve o id criado.
// O code precisa ser único; repetido devolve erro.
func AddSubject(db *sql.DB, code string, name string) (int64, error) {
	res, err := db.Exec("INSERT INTO subjects (code,name) VALUES (?,?)", code, name)

	if err != nil {
		return 0, fmt.Errorf("inserir disciplina: %w", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("retornar id: %w", err)
	}

	return id, nil
}

// ListSubjects devolve todas as disciplinas, na ordem de cadastro.
func ListSubjects(db *sql.DB) ([]subject.Subject, error) {
	rows, err := db.Query("SELECT id, code, name FROM subjects ORDER BY id")

	if err != nil {
		return nil, fmt.Errorf("listar disciplinas: %w", err)
	}

	defer rows.Close()

	var lista []subject.Subject
	for rows.Next() {
		var v subject.Subject

		err := rows.Scan(&v.ID, &v.Code, &v.Name)
		if err != nil {
			return nil, fmt.Errorf("ler disciplina: %w", err)
		}

		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer disciplinas: %w", err)
	}

	return lista, nil
}

// AddProject cadastra um projeto e devolve o id criado.
// O code precisa ser único; repetido devolve erro.
func AddProject(db *sql.DB, code string, name string) (int64, error) {
	res, err := db.Exec("INSERT INTO projects (code,name) VALUES (?,?)", code, name)

	if err != nil {
		return 0, fmt.Errorf("inserir projeto: %w", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("retornar id: %w", err)
	}

	return id, nil
}

// ListProjects devolve todos os projetos, na ordem de cadastro.
func ListProjects(db *sql.DB) ([]project.Project, error) {
	rows, err := db.Query("SELECT id, code, name FROM projects ORDER BY id")

	if err != nil {
		return nil, fmt.Errorf("listar projetos: %w", err)
	}

	defer rows.Close()

	var lista []project.Project
	for rows.Next() {
		var v project.Project

		err := rows.Scan(&v.ID, &v.Code, &v.Name)
		if err != nil {
			return nil, fmt.Errorf("ler projeto: %w", err)
		}

		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer projeto: %w", err)
	}

	return lista, nil
}

// adiquirir a tarefa pelo id
func GetTask(db *sql.DB, id int64) (task.Task, error) {
	v, err := scanTask(db.QueryRow(selectTasks+" WHERE tasks.id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return task.Task{}, fmt.Errorf("A task: %d nao encontrada: %w", id, err)
	} else {
		if err != nil {
			return task.Task{}, fmt.Errorf("erro ao consultar as taks: %w", err)
		}
		return v, nil
	}

}

func DoneTask(db *sql.DB, id int64) error {
	_, err := db.Exec("UPDATE tasks SET done_at = datetime('now', 'localtime') WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("concluir tarefa %d: %w", id, err)
	}

	return nil
}

func ReopenTask(db *sql.DB, id int64) error {
	_, err := db.Exec("UPDATE tasks SET done_at = NULL WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("reabrir tarefa %d: %w", id, err)
	}

	return nil
}

func RmTask(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deletar tarefa %d: %w", id, err)
	}

	return nil
}
