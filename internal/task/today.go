package task

import (
	"cmp"
	"slices"
	"time"
)

const isoDate = "2006-01-02"

// IsNew diz se a tarefa foi criada hoje ou ontem (D-11)
func IsNew(t Task, today time.Time) bool {
	if len(t.CreatedAt) < 10 {
		return false
	}
	yesterday := today.AddDate(0, 0, -1).Format(isoDate)
	return t.CreatedAt[:10] >= yesterday
}

// IsLate diz se a tarefa pendente tem prazo antes de hoje (glossario)
func IsLate(t Task, today time.Time) bool {
	return t.DoneAt == "" && t.Due != "" && t.Due < today.Format(isoDate)
}

// Summary sao os blocos do resumo diario (RN-16)
type Summary struct {
	Late      []Task
	Urgent    []Task
	Tomorrow  []Task // prazo amanha; na sexta, no sabado e no domingo, ate segunda
	Third     []Task
	Important bool // true: o terceiro bloco sao as A sem prazo; false: "sem prazo"

	// sem prazo de projeto ficam fora do terceiro bloco (menos as A):
	// so a contagem por codigo, na ordem de quem tem mais
	Hidden []ProjectCount
}

// ProjectCount eh uma parte da linha "+ 7 de projetos (p1: 4, mrt: 2, neivi: 1)"
type ProjectCount struct {
	Code string
	N    int
}

// HiddenTotal soma as tarefas de projeto que ficaram fora
func (s Summary) HiddenTotal() int {
	n := 0
	for _, c := range s.Hidden {
		n += c.N
	}
	return n
}

// IsEmpty diz se nao tem nada para hoje (amanha e as de projeto nao contam)
func (s Summary) IsEmpty() bool {
	return len(s.Late) == 0 && len(s.Urgent) == 0 && len(s.Third) == 0
}

// TomorrowEnd eh o ultimo dia do bloco amanha: o dia seguinte, ou a segunda
// quando hoje eh sexta, sabado ou domingo (a entrega de segunda aparece no fim de semana)
func TomorrowEnd(today time.Time) time.Time {
	switch today.Weekday() {
	case time.Friday:
		return today.AddDate(0, 0, 3)
	case time.Saturday:
		return today.AddDate(0, 0, 2)
	}
	return today.AddDate(0, 0, 1)
}

// Today monta o resumo diario a partir das tarefas pendentes (RN-16)
func Today(pending []Task, today time.Time) Summary {
	day := today.Format(isoDate)
	end := TomorrowEnd(today).Format(isoDate)

	var s Summary
	var undated []Task // sem prazo e com prioridade (incompleta sem prazo nunca entra)
	hidden := map[string]int{}
	for _, t := range pending {
		switch {
		case t.DoneAt != "":
			continue
		case t.Due != "" && t.Due < day:
			s.Late = append(s.Late, t)
		case t.Due == day:
			s.Urgent = append(s.Urgent, t)
		case t.Due != "" && t.Due <= end:
			s.Tomorrow = append(s.Tomorrow, t)
		case t.Due == "" && t.Priority != "" && t.ProjectID != 0 && t.Priority != "A":
			hidden[t.ProjectCode]++ // o projeto tem horario proprio; a A ganha do projeto
		case t.Due == "" && t.Priority != "":
			undated = append(undated, t)
		}
	}

	// terceiro bloco: as A sem prazo; se nao tiver, o nivel mais alto que existir
	best := ""
	for _, t := range undated {
		if best == "" || t.Priority < best {
			best = t.Priority
		}
	}
	for _, t := range undated {
		if t.Priority == best {
			s.Third = append(s.Third, t)
		}
	}
	s.Important = best == "A"

	for code, n := range hidden {
		s.Hidden = append(s.Hidden, ProjectCount{code, n})
	}
	slices.SortFunc(s.Hidden, func(a, b ProjectCount) int {
		return cmp.Or(cmp.Compare(b.N, a.N), cmp.Compare(a.Code, b.Code))
	})

	SortByDue(s.Late)
	SortByDue(s.Urgent)
	SortByDue(s.Tomorrow)
	SortByDue(s.Third)

	return s
}
