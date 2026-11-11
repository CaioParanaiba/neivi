// Package store guarda e lê os dados do neivi no SQLite.
package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

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
        COALESCE(tasks.project_id, 0), COALESCE(tasks.done_at, ''),
        tasks.created_at, COALESCE(subjects.code, ''),
        COALESCE(projects.code, ''), COALESCE(subjects.name, ''),
        COALESCE(projects.name, '')
        FROM tasks
        LEFT JOIN subjects ON subjects.id = tasks.subject_id
        LEFT JOIN projects ON projects.id = tasks.project_id`

func scanTask(s scanner) (task.Task, error) {
	var v task.Task
	err := s.Scan(&v.ID, &v.Title, &v.Description,
		&v.Priority, &v.Due, &v.EstimateMin,
		&v.SubjectID, &v.ProjectID,
		&v.DoneAt, &v.CreatedAt,
		&v.SubjectCode, &v.ProjectCode,
		&v.SubjectName, &v.ProjectName)
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
	res, err := db.Exec(`INSERT INTO tasks (title, description, priority, due, estimate_min, subject_id, project_id)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.Title,
		nullString(t.Description),
		nullString(t.Priority),
		nullString(t.Due),
		nullInt(t.EstimateMin),
		nullInt(t.SubjectID),
		nullInt(t.ProjectID))
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
		estimate_min = ?, subject_id = ?, project_id = ? WHERE id = ?`,
		t.Title,
		nullString(t.Description),
		nullString(t.Priority),
		nullString(t.Due),
		nullInt(t.EstimateMin),
		nullInt(t.SubjectID),
		nullInt(t.ProjectID),
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
		return 0, fmt.Errorf("disciplina %s nao cadastrada (cadastrar: neivi subject add %s)", code, code)
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
		return 0, fmt.Errorf("projeto %s nao cadastrado (cadastrar: neivi project add %s)", code, code)
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

// DoneTask conclui a tarefa com a data de hoje (daysAgo = 0) ou de ontem (daysAgo = 1)
func DoneTask(db *sql.DB, id int64, daysAgo int) error {
	mod := fmt.Sprintf("-%d days", daysAgo)
	_, err := db.Exec("UPDATE tasks SET done_at = datetime('now', 'localtime', ?) WHERE id = ?", mod, id)
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

// CountDone conta as tarefas concluidas hoje (daysAgo = 0) ou ontem (daysAgo = 1)
func CountDone(db *sql.DB, daysAgo int) (int, error) {
	mod := fmt.Sprintf("-%d days", daysAgo)

	var n int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE date(done_at) = date('now', 'localtime', ?)", mod).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("contar concluidas: %w", err)
	}

	return n, nil
}

// FindDuplicate devolve o id de uma tarefa pendente com o mesmo titulo (RN-01),
// ignorando a tarefa exceptID (0 no add). Devolve 0 se nao tiver.
func FindDuplicate(db *sql.DB, title string, exceptID int64) (int64, error) {
	rows, err := db.Query("SELECT id, title FROM tasks WHERE done_at IS NULL AND id != ? ORDER BY id", exceptID)
	if err != nil {
		return 0, fmt.Errorf("procurar duplicada: %w", err)
	}

	defer rows.Close()

	want := task.NormalizeTitle(title)
	for rows.Next() {
		var id int64
		var t string

		err := rows.Scan(&id, &t)
		if err != nil {
			return 0, fmt.Errorf("ler tarefa: %w", err)
		}

		if task.NormalizeTitle(t) == want {
			return id, nil
		}
	}

	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("percorrer tarefas: %w", err)
	}

	return 0, nil
}

// orGroup junta as condicoes de um mesmo tipo de filtro com OR, entre
// parenteses, para depois entrar no WHERE com AND (RN-18)
func orGroup(parts []string) string {
	return "(" + strings.Join(parts, " OR ") + ")"
}

// ListTasksFilter devolve as tarefas que passam nos filtros que o banco sabe
// responder. A busca por palavras (sem acento) e a ordem ficam em Go, no task.
func ListTasksFilter(db *sql.DB, f task.Filter) ([]task.Task, error) {
	var where []string // cada item eh uma condicao; todas juntas com AND
	var args []any     // os valores dos ?, na mesma ordem em que aparecem

	// situacao: so pendentes, so concluidas ou as duas (all)
	if f.Done {
		where = append(where, "tasks.done_at IS NOT NULL")
	} else if !f.All {
		where = append(where, "tasks.done_at IS NULL")
	}

	// disciplina: @BD, @ e @- do mesmo tipo viram um grupo com OR
	var subj []string
	for _, code := range f.Subjects {
		subj = append(subj, "subjects.code = ?")
		args = append(args, code)
	}
	if f.AnySubject {
		subj = append(subj, "tasks.subject_id IS NOT NULL")
	}
	if f.NoSubject {
		subj = append(subj, "tasks.subject_id IS NULL")
	}
	if len(subj) > 0 {
		where = append(where, orGroup(subj))
	}

	// projeto: o mesmo jeito da disciplina
	var proj []string
	for _, code := range f.Projects {
		proj = append(proj, "projects.code = ?")
		args = append(args, code)
	}
	if f.AnyProject {
		proj = append(proj, "tasks.project_id IS NOT NULL")
	}
	if f.NoProject {
		proj = append(proj, "tasks.project_id IS NULL")
	}
	if len(proj) > 0 {
		where = append(where, orGroup(proj))
	}

	// prioridade: "-" eh sem prioridade (incompleta)
	var pri []string
	for _, p := range f.Priorities {
		if p == "-" {
			pri = append(pri, "tasks.priority IS NULL")
			continue
		}
		pri = append(pri, "tasks.priority = ?")
		args = append(args, p)
	}
	if len(pri) > 0 {
		where = append(where, orGroup(pri))
	}

	if f.New {
		where = append(where, "date(tasks.created_at) >= date('now', 'localtime', '-1 day')")
	}
	if f.Late {
		where = append(where, "tasks.due IS NOT NULL AND tasks.due < date('now', 'localtime')")
	}

	query := selectTasks
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY tasks.id" // ordem estavel; a ordem final eh feita em Go (RN-17)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listar tarefas: %w", err)
	}

	defer rows.Close()

	var lista []task.Task
	for rows.Next() {
		v, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("ler tarefa: %w", err)
		}

		// busca no titulo: em Go, porque o SQLite nao sabe ignorar acento
		if !task.MatchWords(v.Title, f.Words) {
			continue
		}

		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer tarefas: %w", err)
	}

	return lista, nil
}
