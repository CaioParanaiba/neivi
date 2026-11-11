package todotxt

import (
	"fmt"
	"strings"

	"github.com/CaioParanaiba/neivi/internal/project"
	"github.com/CaioParanaiba/neivi/internal/subject"
)

// ProjectEntry eh o que veio no project add/edit. Como na tarefa: "" = nao veio,
// "-" = apagar (so no edit).
type ProjectEntry struct {
	Code        string // o primeiro argumento
	Name        string // texto solto
	NewCode     string // code:novo (so no edit)
	SubjectCode string // @BD
	GitHub      string // gh:usuario/repo ou a URL
	Dir         string // dir:~/Code/neivi
	Status      string // sit:pausado, ja em ingles (paused)
	Continuous  string // "sim", "nao" ou ""
	Description string // -- texto
}

// SubjectEntry eh o que veio no subject add/edit
type SubjectEntry struct {
	Code    string
	Name    string
	NewCode string
	Status  string // ja em ingles (active, closed)
}

// githubURL aceita usuario/repo (vira o link completo) ou a URL inteira
func githubURL(v string) (string, error) {
	if v == "-" || strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return v, nil
	}
	user, repo, found := strings.Cut(v, "/")
	if !found || user == "" || repo == "" || strings.Contains(repo, "/") {
		return "", fmt.Errorf("github %s invalido (use usuario/repo ou a URL)", v)
	}
	return "https://github.com/" + v, nil
}

// ParseProject le o project add (isEdit = false) ou o project edit (isEdit = true)
func ParseProject(args []string, isEdit bool) (ProjectEntry, error) {
	var e ProjectEntry
	if len(args) == 0 || args[0] == "" {
		return ProjectEntry{}, fmt.Errorf("falta o codigo do projeto (exemplo: neivi project add neivi)")
	}
	e.Code = args[0]

	var name, desc []string
	var inDesc bool
	for _, v := range args[1:] {
		if v == "--" {
			inDesc = !inDesc
			continue
		}
		if inDesc {
			desc = append(desc, v)
			continue
		}

		if strings.HasPrefix(v, "@") && len(v) > 1 {
			e.SubjectCode = v[1:]
			continue
		}

		key, value, found := strings.Cut(v, ":")
		if !found {
			name = append(name, v)
			continue
		}

		var err error
		switch key {
		case "gh":
			e.GitHub, err = githubURL(value)
		case "dir":
			e.Dir = value
		case "sit":
			if value == "-" {
				err = fmt.Errorf("a situacao nao pode ser apagada (use ativo, pausado, concluido ou cancelado)")
				break
			}
			e.Status, err = project.StatusFromPT(value)
		case "cont":
			if value != "sim" && value != "nao" {
				err = fmt.Errorf("cont:%s invalido (use cont:sim ou cont:nao)", value)
			}
			e.Continuous = value
		case "code":
			if !isEdit {
				err = fmt.Errorf("code: so vale no edit (no add, o codigo eh a primeira palavra)")
			}
			e.NewCode = value
		default:
			name = append(name, v) // outra palavra com ":" faz parte do nome (como no titulo, RN-04)
		}
		if err != nil {
			return ProjectEntry{}, err
		}
	}

	e.Name = strings.Join(name, " ")
	e.Description = strings.Join(desc, " ")

	if !isEdit && (e.Dir == "-" || e.GitHub == "-" || e.Description == "-" || e.SubjectCode == "-") {
		return ProjectEntry{}, fmt.Errorf(`"-" so vale no edit (para apagar um campo: neivi project edit %s dir:-)`, e.Code)
	}

	return e, nil
}

// IsEmpty diz se o edit nao trouxe nada para mudar
func (e ProjectEntry) IsEmpty() bool {
	return e.Name == "" && e.NewCode == "" && e.SubjectCode == "" && e.GitHub == "" &&
		e.Dir == "" && e.Status == "" && e.Continuous == "" && e.Description == ""
}

// ParseSubject le o subject add ou o subject edit: codigo, nome solto, sit: e code:
func ParseSubject(args []string, isEdit bool) (SubjectEntry, error) {
	var e SubjectEntry
	if len(args) == 0 || args[0] == "" {
		return SubjectEntry{}, fmt.Errorf("falta o codigo da disciplina (exemplo: neivi subject add BD)")
	}
	e.Code = args[0]

	var name []string
	for _, v := range args[1:] {
		key, value, found := strings.Cut(v, ":")
		if !found {
			name = append(name, v)
			continue
		}

		var err error
		switch key {
		case "sit":
			e.Status, err = subject.StatusFromPT(value)
		case "code":
			if !isEdit {
				err = fmt.Errorf("code: so vale no edit (no add, o codigo eh a primeira palavra)")
			}
			e.NewCode = value
		default:
			name = append(name, v)
		}
		if err != nil {
			return SubjectEntry{}, err
		}
	}

	e.Name = strings.Join(name, " ")
	return e, nil
}

// IsEmpty diz se o edit nao trouxe nada para mudar
func (e SubjectEntry) IsEmpty() bool {
	return e.Name == "" && e.NewCode == "" && e.Status == ""
}
