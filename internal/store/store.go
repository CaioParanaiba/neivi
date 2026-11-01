// Package store guarda e lê os dados do neivi no SQLite.
package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

//const path_sql = "~/Code/neivi/migrations/001_init.sql"

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}
	return db, nil
}
