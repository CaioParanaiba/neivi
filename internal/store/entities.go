package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/subject"
)

// ---- disciplinas ----

// a contagem de pendentes usa a disciplina efetiva (RN-14): as tarefas com o @
// dela e as dos projetos dela que nao tem @ proprio
const selectSubjects = `SELECT s.id, s.code, s.name, s.status,
        (SELECT COUNT(*) FROM tasks t LEFT JOIN projects p ON p.id = t.project_id
         WHERE t.done_at IS NULL AND COALESCE(t.subject_id, p.subject_id) = s.id)
        FROM subjects s`

func scanSubject(sc scanner) (subject.Subject, error) {
	var v subject.Subject
	err := sc.Scan(&v.ID, &v.Code, &v.Name, &v.Status, &v.Pending)
	return v, err
}

// AddSubject cadastra uma disciplina e devolve o id criado.
// O code precisa ser unico; repetido devolve erro.
func AddSubject(db *sql.DB, s subject.Subject) (int64, error) {
	if s.Status == "" {
		s.Status = "active"
	}

	res, err := db.Exec("INSERT INTO subjects (code, name, status) VALUES (?, ?, ?)", s.Code, s.Name, s.Status)
	if err != nil {
		return 0, fmt.Errorf("inserir disciplina: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("retornar id: %w", err)
	}

	return id, nil
}

// GetSubject busca a disciplina pelo codigo; nao existe = sql.ErrNoRows
func GetSubject(db *sql.DB, code string) (subject.Subject, error) {
	return scanSubject(db.QueryRow(selectSubjects+" WHERE s.code = ?", code))
}

// ListSubjects devolve as disciplinas: so as ativas, ou todas (all), as ativas primeiro
func ListSubjects(db *sql.DB, all bool) ([]subject.Subject, error) {
	query := selectSubjects
	if !all {
		query += " WHERE s.status = 'active'"
	}
	query += " ORDER BY s.status = 'closed', s.id"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar disciplinas: %w", err)
	}

	defer rows.Close()

	var lista []subject.Subject
	for rows.Next() {
		v, err := scanSubject(rows)
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

// UpdateSubject grava codigo, nome e situacao
func UpdateSubject(db *sql.DB, s subject.Subject) error {
	_, err := db.Exec("UPDATE subjects SET code = ?, name = ?, status = ? WHERE id = ?", s.Code, s.Name, s.Status, s.ID)
	if err != nil {
		return fmt.Errorf("atualizar disciplina %s: %w", s.Code, err)
	}
	return nil
}

// SubjectLinks conta o que esta ligado a disciplina (RN-22): tarefas com o @ dela
// e os codigos dos projetos dela
func SubjectLinks(db *sql.DB, id int64) (int, []string, error) {
	var tasks int
	err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE subject_id = ?", id).Scan(&tasks)
	if err != nil {
		return 0, nil, fmt.Errorf("contar tarefas: %w", err)
	}

	rows, err := db.Query("SELECT code FROM projects WHERE subject_id = ? ORDER BY id", id)
	if err != nil {
		return 0, nil, fmt.Errorf("listar projetos: %w", err)
	}

	defer rows.Close()

	var projects []string
	for rows.Next() {
		var code string
		err := rows.Scan(&code)
		if err != nil {
			return 0, nil, fmt.Errorf("ler projeto: %w", err)
		}
		projects = append(projects, code)
	}

	return tasks, projects, rows.Err()
}

// RmSubject apaga a disciplina (o cli confere antes que nada esta ligado)
func RmSubject(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM subjects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deletar disciplina %d: %w", id, err)
	}
	return nil
}

// ---- projetos ----

// progresso e ultima atividade sao contados aqui, nunca guardados (RN-20, D-17)
const selectProjects = `SELECT p.id, p.code, p.name, COALESCE(p.description, ''),
        COALESCE(p.github, ''), COALESCE(p.dir, ''), p.status, p.continuous,
        COALESCE(p.subject_id, 0), COALESCE(s.code, ''), COALESCE(s.name, ''),
        p.created_at,
        (SELECT COUNT(*) FROM tasks t WHERE t.project_id = p.id AND t.done_at IS NOT NULL),
        (SELECT COUNT(*) FROM tasks t WHERE t.project_id = p.id AND t.done_at IS NULL),
        COALESCE((SELECT MAX(t.done_at) FROM tasks t WHERE t.project_id = p.id), '')
        FROM projects p
        LEFT JOIN subjects s ON s.id = p.subject_id`

func scanProject(sc scanner) (project.Project, error) {
	var v project.Project
	err := sc.Scan(&v.ID, &v.Code, &v.Name, &v.Description,
		&v.GitHub, &v.Dir, &v.Status, &v.Continuous,
		&v.SubjectID, &v.SubjectCode, &v.SubjectName,
		&v.CreatedAt,
		&v.Done, &v.Pending, &v.LastDoneAt)
	return v, err
}

// AddProject cadastra um projeto e devolve o id criado.
// O code precisa ser unico; repetido devolve erro.
func AddProject(db *sql.DB, p project.Project) (int64, error) {
	if p.Status == "" {
		p.Status = "active"
	}

	res, err := db.Exec(`INSERT INTO projects (code, name, description, github, dir, status, continuous, subject_id)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Code, p.Name, nullString(p.Description), nullString(p.GitHub), nullString(p.Dir),
		p.Status, p.Continuous, nullInt(p.SubjectID))
	if err != nil {
		return 0, fmt.Errorf("inserir projeto: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("retornar id: %w", err)
	}

	return id, nil
}

// GetProject busca o projeto pelo codigo; nao existe = sql.ErrNoRows
func GetProject(db *sql.DB, code string) (project.Project, error) {
	return scanProject(db.QueryRow(selectProjects+" WHERE p.code = ?", code))
}

// ListProjects devolve os projetos: ativos e pausados, ou todos (all).
// A ordem dos grupos e da ultima atividade fica no cli.
func ListProjects(db *sql.DB, all bool) ([]project.Project, error) {
	query := selectProjects
	if !all {
		query += " WHERE p.status IN ('active', 'paused')"
	}
	query += " ORDER BY p.id"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar projetos: %w", err)
	}

	defer rows.Close()

	var lista []project.Project
	for rows.Next() {
		v, err := scanProject(rows)
		if err != nil {
			return nil, fmt.Errorf("ler projeto: %w", err)
		}
		lista = append(lista, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("percorrer projetos: %w", err)
	}

	return lista, nil
}

// UpdateProject grava todos os campos editaveis
func UpdateProject(db *sql.DB, p project.Project) error {
	_, err := db.Exec(`UPDATE projects SET code = ?, name = ?, description = ?, github = ?, dir = ?,
        status = ?, continuous = ?, subject_id = ? WHERE id = ?`,
		p.Code, p.Name, nullString(p.Description), nullString(p.GitHub), nullString(p.Dir),
		p.Status, p.Continuous, nullInt(p.SubjectID), p.ID)
	if err != nil {
		return fmt.Errorf("atualizar projeto %s: %w", p.Code, err)
	}
	return nil
}

// LastDone devolve a ultima tarefa concluida do projeto (id, titulo, data); id 0 se nao tiver
func LastDone(db *sql.DB, projectID int64) (int64, string, string, error) {
	var id int64
	var title, doneAt string
	err := db.QueryRow(`SELECT id, title, done_at FROM tasks
        WHERE project_id = ? AND done_at IS NOT NULL
        ORDER BY done_at DESC, id DESC LIMIT 1`, projectID).Scan(&id, &title, &doneAt)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", "", nil
	}
	if err != nil {
		return 0, "", "", fmt.Errorf("ultima tarefa: %w", err)
	}
	return id, title, doneAt, nil
}

// RmProject apaga o projeto (o cli confere antes que nao tem tarefas)
func RmProject(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deletar projeto %d: %w", id, err)
	}
	return nil
}
