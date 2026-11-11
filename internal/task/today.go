package task

import "time"

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

// Summary sao os tres blocos do resumo diario (RN-16)
type Summary struct {
	Late      []Task
	Urgent    []Task
	Third     []Task
	Important bool // true: o terceiro bloco sao as A sem prazo; false: "sem prazo"
}

// IsEmpty diz se nao tem nada para hoje
func (s Summary) IsEmpty() bool {
	return len(s.Late) == 0 && len(s.Urgent) == 0 && len(s.Third) == 0
}

// Today monta o resumo diario a partir das tarefas pendentes (RN-16)
func Today(pending []Task, today time.Time) Summary {
	day := today.Format(isoDate)

	var s Summary
	var undated []Task // sem prazo e com prioridade (incompleta sem prazo nunca entra)
	for _, t := range pending {
		switch {
		case t.DoneAt != "":
			continue
		case t.Due != "" && t.Due < day:
			s.Late = append(s.Late, t)
		case t.Due == day:
			s.Urgent = append(s.Urgent, t)
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

	SortByDue(s.Late)
	SortByDue(s.Urgent)
	SortByDue(s.Third)

	return s
}
