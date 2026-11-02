// Package store guarda e lê os dados do neivi no SQLite.
package store

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/subject"
	"github.com/CaioParanaiba/neivi/internal/task"
)

//const path_sql = "~/Code/neivi/migrations/001_init.sql"

//go:embed migrations/001_init.sql
var schema string

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

func AddTask(db *sql.DB, title string) (int64, error) {
	res, err := db.Exec("INSERT INTO tasks (title) VALUES (?)", title)
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
	rows, err := db.Query("SELECT id, title, created_at FROM tasks ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("listar tarefas: %w", err)
	}

	defer rows.Close()

	var lista []task.Task
	for rows.Next() {
		var v task.Task

		err := rows.Scan(&v.ID, &v.Title, &v.CreatedAt)
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
