// Package store guarda e lê os dados do neivi no SQLite.
package store

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
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
