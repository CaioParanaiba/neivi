package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/store"
	"github.com/CaioParanaiba/neivi/internal/subject"
	"github.com/CaioParanaiba/neivi/internal/task"
	"github.com/CaioParanaiba/neivi/internal/todotxt"
)

// fail imprime o erro e sai (o mesmo que os outros comandos fazem em cada if)
func fail(err error) {
	fmt.Fprintln(os.Stderr, "erro:", err)
	os.Exit(1)
}

// ---- avisos de inativo (RN-19) ----

// warnInactive avisa quando a tarefa ou o projeto foi ligado a algo que nao esta ativo
func warnInactive(db *sql.DB, subjectCode, projectCode string) {
	if subjectCode != "" && subjectCode != "-" {
		s, err := store.GetSubject(db, subjectCode)
		if err == nil && s.Status != "active" {
			fmt.Fprintf(os.Stderr, "aviso: a disciplina %s esta %s (reativar: neivi subject edit %s sit:ativa)\n",
				s.Code, subject.StatusPT(s.Status), s.Code)
		}
	}
	if projectCode != "" && projectCode != "-" {
		p, err := store.GetProject(db, projectCode)
		if err == nil && p.Status != "active" {
			fmt.Fprintf(os.Stderr, "aviso: o projeto %s esta %s (reativar: neivi project edit %s sit:ativo)\n",
				p.Code, project.StatusPT(p.Status), p.Code)
		}
	}
}

// warnDir avisa se a pasta do projeto nao existe (aceita o ~ do inicio)
func warnDir(dir string) {
	if dir == "" || dir == "-" {
		return
	}
	path := dir
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	_, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "aviso: a pasta %s nao existe\n", dir)
	}
}

// agoLabel diz ha quanto tempo foi: hoje, ha 1 dia, ha N dias
func agoLabel(date string, today time.Time) string {
	switch l := daysLabel(date, today); {
	case l == "hoje":
		return "hoje"
	case strings.HasPrefix(l, "atrasada "):
		return strings.TrimPrefix(l, "atrasada ")
	default:
		return l
	}
}

// confirmEntity eh a confirmacao curta: "+ 3 · neivi · Meu sucessor · ativo" (nome igual ao codigo nao se repete)
func confirmEntity(symbol string, id int64, code, name, status string) string {
	parts := []string{symbol + " " + strconv.FormatInt(id, 10), code}
	if name != "" && name != code {
		parts = append(parts, name)
	}
	return strings.Join(append(parts, status), " · ")
}

// askYes faz a pergunta [s/N]; so "s" confirma
func askYes(question string) bool {
	fmt.Print(question + " [s/N] ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.ToUpper(strings.TrimSpace(line)) == "S"
}

// ---- disciplinas ----

func runSubject(db *sql.DB, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "erro: falta o comando (exemplo: neivi subject ls; ou add, edit, rm)")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		runSubjectAdd(db, args[1:])
	case "ls":
		runSubjectLs(db, args[1:])
	case "edit":
		runSubjectEdit(db, args[1:])
	case "rm":
		runSubjectRm(db, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "erro: comando subject %s desconhecido (use add, ls, edit ou rm)\n", args[0])
		os.Exit(1)
	}
}

func runSubjectAdd(db *sql.DB, args []string) {
	e, err := todotxt.ParseSubject(args, false)
	if err != nil {
		fail(err)
	}

	_, err = store.GetSubject(db, e.Code)
	if err == nil {
		fail(fmt.Errorf("ja existe uma disciplina %s (para mudar: neivi subject edit %s)", e.Code, e.Code))
	}

	s := subject.Subject{Code: e.Code, Name: e.Name, Status: e.Status}
	if s.Name == "" {
		s.Name = s.Code
	}
	if s.Status == "" {
		s.Status = "active"
	}

	id, err := store.AddSubject(db, s)
	if err != nil {
		fail(err)
	}

	fmt.Println(confirmEntity("+", id, s.Code, s.Name, subject.StatusPT(s.Status)))
}

func runSubjectLs(db *sql.DB, args []string) {
	all := len(args) == 1 && args[0] == "all"
	if len(args) > 0 && !all {
		fail(fmt.Errorf("subject ls so aceita all (exemplo: neivi subject ls all)"))
	}

	lista, err := store.ListSubjects(db, all)
	if err != nil {
		fail(err)
	}

	if len(lista) == 0 {
		fmt.Println("nenhuma disciplina (cadastrar: neivi subject add codigo)")
		return
	}

	head := []string{"ID", "DISC", "NOME", "PENDENTES"}
	if all {
		head = []string{"ID", "DISC", "NOME", "SITUACAO", "PENDENTES"}
	}

	var tb table
	tb.fit(head)
	var rows [][]string
	for _, s := range lista {
		r := []string{strconv.FormatInt(s.ID, 10), s.Code, s.Name}
		if all {
			r = append(r, subject.StatusPT(s.Status))
		}
		r = append(r, strconv.Itoa(s.Pending))
		tb.fit(r)
		rows = append(rows, r)
	}

	fmt.Println(tb.line(head))
	for _, r := range rows {
		fmt.Println(tb.line(r))
	}
}

func runSubjectEdit(db *sql.DB, args []string) {
	e, err := todotxt.ParseSubject(args, true)
	if err != nil {
		fail(fmt.Errorf("%s", strings.Replace(err.Error(), "subject add BD", "subject edit BD sit:encerrada", 1)))
	}

	s, err := store.GetSubject(db, e.Code)
	if errors.Is(err, sql.ErrNoRows) {
		fail(fmt.Errorf("disciplina %s nao cadastrada (cadastrar: neivi subject add %s)", e.Code, e.Code))
	}
	if err != nil {
		fail(err)
	}

	if e.IsEmpty() {
		fail(fmt.Errorf("nada para mudar na disciplina %s (exemplo: neivi subject edit %s sit:encerrada)", s.Code, s.Code))
	}

	if e.NewCode != "" && e.NewCode != s.Code {
		_, err := store.GetSubject(db, e.NewCode)
		if err == nil {
			fail(fmt.Errorf("ja existe uma disciplina %s", e.NewCode))
		}
		s.Code = e.NewCode
	}
	if e.Name != "" {
		s.Name = e.Name
	}
	if e.Status != "" {
		s.Status = e.Status
	}

	err = store.UpdateSubject(db, s)
	if err != nil {
		fail(err)
	}

	fmt.Println(confirmEntity("~", s.ID, s.Code, s.Name, subject.StatusPT(s.Status)))
}

func runSubjectRm(db *sql.DB, args []string) {
	if len(args) != 1 {
		fail(fmt.Errorf("falta o codigo da disciplina (exemplo: neivi subject rm BD2)"))
	}

	s, err := store.GetSubject(db, args[0])
	if errors.Is(err, sql.ErrNoRows) {
		fail(fmt.Errorf("disciplina %s nao cadastrada (veja as disciplinas no neivi subject ls all)", args[0]))
	}
	if err != nil {
		fail(err)
	}

	// so sem nada ligado (RN-22)
	tasks, projects, err := store.SubjectLinks(db, s.ID)
	if err != nil {
		fail(err)
	}
	if tasks > 0 {
		fail(fmt.Errorf("a disciplina %s tem %d tarefas (mude as tarefas para outra com neivi edit, ou tire da lista: neivi subject edit %s sit:encerrada)",
			s.Code, tasks, s.Code))
	}
	if len(projects) > 0 {
		fail(fmt.Errorf("a disciplina %s tem %d projeto (%s)", s.Code, len(projects), strings.Join(projects, ", ")))
	}

	if !askYes(fmt.Sprintf("remover a disciplina %s (sem tarefas)?", s.Code)) {
		fmt.Printf("%s mantida\n", s.Code)
		return
	}

	err = store.RmSubject(db, s.ID)
	if err != nil {
		fail(err)
	}
	fmt.Printf("✗ %s removida\n", s.Code)
}

// ---- projetos ----

func runProject(db *sql.DB, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "erro: falta o comando (exemplo: neivi project ls; ou add, show, edit, rm)")
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		runProjectAdd(db, args[1:])
	case "ls":
		runProjectLs(db, args[1:])
	case "show":
		runProjectShow(db, args[1:])
	case "edit":
		runProjectEdit(db, args[1:])
	case "rm":
		runProjectRm(db, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "erro: comando project %s desconhecido (use add, ls, show, edit ou rm)\n", args[0])
		os.Exit(1)
	}
}

func runProjectAdd(db *sql.DB, args []string) {
	e, err := todotxt.ParseProject(args, false)
	if err != nil {
		fail(err)
	}

	_, err = store.GetProject(db, e.Code)
	if err == nil {
		fail(fmt.Errorf("ja existe um projeto %s (para mudar: neivi project edit %s)", e.Code, e.Code))
	}

	p := project.Project{
		Code: e.Code, Name: e.Name, Description: e.Description,
		GitHub: e.GitHub, Dir: e.Dir, Status: e.Status, Continuous: e.Continuous == "sim",
	}
	if p.Name == "" {
		p.Name = p.Code
	}
	if p.Status == "" {
		p.Status = "active"
	}
	if e.SubjectCode != "" {
		p.SubjectID, err = store.SubjectIDByCode(db, e.SubjectCode)
		if err != nil {
			fail(err)
		}
	}

	id, err := store.AddProject(db, p)
	if err != nil {
		fail(err)
	}

	fmt.Println(confirmEntity("+", id, p.Code, p.Name, project.StatusPT(p.Status)))
	warnInactive(db, e.SubjectCode, "")
	warnDir(p.Dir)
}

// progressLabel eh a coluna PROGRESSO: 71%, continuo ou - (RN-20)
func progressLabel(p project.Project) string {
	if p.Continuous {
		return "continuo"
	}
	pct, ok := p.Progress()
	if !ok {
		return "-"
	}
	return strconv.Itoa(pct) + "%"
}

func runProjectLs(db *sql.DB, args []string) {
	all := len(args) == 1 && args[0] == "all"
	if len(args) > 0 && !all {
		fail(fmt.Errorf("project ls so aceita all (exemplo: neivi project ls all)"))
	}

	lista, err := store.ListProjects(db, all)
	if err != nil {
		fail(err)
	}

	if len(lista) == 0 {
		fmt.Println("nenhum projeto (cadastrar: neivi project add codigo)")
		return
	}

	// grupos pela situacao; dentro, a atividade mais recente primeiro; sem atividade no fim
	slices.SortFunc(lista, func(a, b project.Project) int {
		if d := project.StatusOrder(a.Status) - project.StatusOrder(b.Status); d != 0 {
			return d
		}
		if (a.LastDoneAt == "") != (b.LastDoneAt == "") {
			if a.LastDoneAt == "" {
				return 1
			}
			return -1
		}
		if c := strings.Compare(b.LastDoneAt, a.LastDoneAt); c != 0 {
			return c
		}
		return int(a.ID - b.ID)
	})

	now := time.Now()
	head := []string{"ID", "PROJETO", "SITUACAO", "DISC", "PROGRESSO", "PENDENTES", "ULTIMA"}
	var tb table
	tb.fit(head)
	var rows [][]string
	for _, p := range lista {
		last := "-"
		if p.LastDoneAt != "" {
			day := p.LastDoneAt[:10]
			last = formatDue(day, now) + " (" + agoLabel(day, now) + ")"
		}
		disc := p.SubjectCode
		if disc == "" {
			disc = "-"
		}
		r := []string{strconv.FormatInt(p.ID, 10), p.Code, project.StatusPT(p.Status), disc,
			progressLabel(p), strconv.Itoa(p.Pending), last}
		tb.fit(r)
		rows = append(rows, r)
	}

	fmt.Println(tb.line(head))
	for _, r := range rows {
		fmt.Println(tb.line(r))
	}
}

// getProject busca pelo codigo e sai com a mensagem certa se nao existir
func getProject(db *sql.DB, args []string, example string) project.Project {
	if len(args) == 0 {
		fail(fmt.Errorf("falta o codigo do projeto (exemplo: %s)", example))
	}

	p, err := store.GetProject(db, args[0])
	if errors.Is(err, sql.ErrNoRows) {
		fail(fmt.Errorf("projeto %s nao cadastrado (cadastrar: neivi project add %s)", args[0], args[0]))
	}
	if err != nil {
		fail(err)
	}

	return p
}

func runProjectShow(db *sql.DB, args []string) {
	p := getProject(db, args, "neivi project show neivi")
	if len(args) > 1 {
		fail(fmt.Errorf("project show mostra um projeto por vez (exemplo: neivi project show %s)", p.Code))
	}

	now := time.Now()
	dash := func(s string) string {
		if s == "" {
			return "-"
		}
		return s
	}

	title := fmt.Sprintf("projeto %d · %s", p.ID, p.Code)
	if p.Name != p.Code {
		title += " · " + p.Name
	}
	fmt.Println(title)

	field := func(name, value string) {
		fmt.Printf("  %-12s %s\n", name+":", value)
	}

	disc := "-  (pessoal)"
	if p.SubjectCode != "" {
		disc = p.SubjectCode + " (academico)"
	}

	progress := "-  (sem tarefas)"
	switch {
	case p.Continuous:
		progress = fmt.Sprintf("continuo · %d pendentes", p.Pending)
	case p.Done+p.Pending > 0:
		pct, _ := p.Progress()
		progress = fmt.Sprintf("%d%% · %d pendentes", pct, p.Pending)
	}

	lastID, lastTitle, lastAt, err := store.LastDone(db, p.ID)
	if err != nil {
		fail(err)
	}
	last := "-"
	if lastID != 0 {
		day := lastAt[:10]
		last = fmt.Sprintf("%s · %d · %s  (%s)", formatDue(day, now), lastID, lastTitle, agoLabel(day, now))
	}

	field("situacao", project.StatusPT(p.Status))
	field("disciplina", disc)
	field("github", dash(p.GitHub))
	field("pasta", dash(p.Dir))
	field("criado", formatDue(p.CreatedAt[:10], now))
	field("progresso", progress)
	field("ultima", last)

	if p.Description != "" {
		fmt.Println()
		fmt.Println("  " + p.Description)
	}

	pending, err := store.ListTasksFilter(db, task.Filter{Projects: []string{p.Code}})
	if err != nil {
		fail(err)
	}

	fmt.Println()
	if len(pending) == 0 {
		fmt.Println("nenhuma tarefa pendente")
		return
	}
	task.SortDefault(pending)
	fmt.Printf("pendentes (%d)\n", len(pending))
	printTasks(pending, false, false, now)
}

// applyField aplica um campo de texto do edit: "" nao muda, "-" apaga, o resto troca
func applyField(field *string, v string) {
	if v == "-" {
		*field = ""
	} else if v != "" {
		*field = v
	}
}

func runProjectEdit(db *sql.DB, args []string) {
	e, err := todotxt.ParseProject(args, true)
	if err != nil {
		fail(fmt.Errorf("%s", strings.Replace(err.Error(), "project add neivi", "project edit neivi sit:pausado", 1)))
	}

	p := getProject(db, args, "neivi project edit neivi sit:pausado")

	if e.IsEmpty() {
		fail(fmt.Errorf("nada para mudar no projeto %s (exemplo: neivi project edit %s sit:pausado)", p.Code, p.Code))
	}

	if e.NewCode != "" && e.NewCode != p.Code {
		_, err := store.GetProject(db, e.NewCode)
		if err == nil {
			fail(fmt.Errorf("ja existe um projeto %s", e.NewCode))
		}
		p.Code = e.NewCode
	}

	applyField(&p.Name, e.Name)
	applyField(&p.Description, e.Description)
	applyField(&p.GitHub, e.GitHub)
	applyField(&p.Dir, e.Dir)
	if e.Status != "" {
		p.Status = e.Status
	}
	if e.Continuous != "" {
		p.Continuous = e.Continuous == "sim"
	}

	switch e.SubjectCode {
	case "":
	case "-":
		p.SubjectID = 0
	default:
		p.SubjectID, err = store.SubjectIDByCode(db, e.SubjectCode)
		if err != nil {
			fail(err)
		}
	}

	err = store.UpdateProject(db, p)
	if err != nil {
		fail(err)
	}

	fmt.Println(confirmEntity("~", p.ID, p.Code, p.Name, project.StatusPT(p.Status)))
	warnInactive(db, e.SubjectCode, "")
	warnDir(e.Dir)
}

func runProjectRm(db *sql.DB, args []string) {
	p := getProject(db, args, "neivi project rm codigo")

	// so sem tarefas (RN-22)
	if n := p.Done + p.Pending; n > 0 {
		fail(fmt.Errorf("o projeto %s tem %d tarefas (mude as tarefas para outro com neivi edit, ou tire da lista: neivi project edit %s sit:cancelado)",
			p.Code, n, p.Code))
	}

	if !askYes(fmt.Sprintf("remover o projeto %s (sem tarefas)?", p.Code)) {
		fmt.Printf("%s mantido\n", p.Code)
		return
	}

	err := store.RmProject(db, p.ID)
	if err != nil {
		fail(err)
	}
	fmt.Printf("✗ %s removido\n", p.Code)
}
