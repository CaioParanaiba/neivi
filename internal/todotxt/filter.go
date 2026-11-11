package todotxt

import (
	"fmt"
	"strings"

	"github.com/CaioParanaiba/neivi/internal/task"
)

// ParseFilter le os argumentos do ls (RN-18). Marcas escolhem, palavras
// reservadas mudam o que mostrar, e o resto eh busca no titulo.
func ParseFilter(args []string) (task.Filter, error) {
	var f task.Filter

	for _, v := range args {
		switch v {
		case "@":
			f.AnySubject = true
			continue
		case "@-":
			f.NoSubject = true
			continue
		case "+":
			f.AnyProject = true
			continue
		case "+-":
			f.NoProject = true
			continue
		case "all":
			f.All = true
			continue
		case "done":
			f.Done = true
			continue
		case "late":
			f.Late = true
			continue
		case "new":
			f.New = true
			continue
		case "due":
			f.Due = true
			continue
		case "desc":
			f.Desc = true
			continue
		}

		switch {
		case strings.HasPrefix(v, "@"):
			f.Subjects = append(f.Subjects, v[1:])
		case strings.HasPrefix(v, "+"):
			f.Projects = append(f.Projects, v[1:])
		case strings.HasPrefix(v, "pri:"):
			p := strings.ToUpper(v[4:])
			if p != "A" && p != "B" && p != "C" && p != "D" && p != "-" {
				return task.Filter{}, fmt.Errorf("prioridade %s invalida (use A, B, C, D ou - para sem prioridade)", v[4:])
			}
			f.Priorities = append(f.Priorities, p)
		default:
			f.Words = append(f.Words, v)
		}
	}

	if f.All && f.Done {
		return task.Filter{}, fmt.Errorf("use all ou done, nao os dois (all = tudo; done = so concluidas)")
	}

	return f, nil
}
