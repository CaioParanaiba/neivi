package task

import (
	"cmp"
	"slices"
)

// Cada funcao cmp* compara um campo e devolve -1 (a vem antes), 0 (empate)
// ou 1 (b vem antes). O cmp.Or pega o primeiro resultado que nao for empate:
// eh assim que os criterios ficam encadeados (RN-17).

// prioridade: A antes de D; sem prioridade por ultimo
func cmpPriority(a, b Task) int {
	if a.Priority == b.Priority {
		return 0
	}
	if a.Priority == "" {
		return 1
	}
	if b.Priority == "" {
		return -1
	}
	return cmp.Compare(a.Priority, b.Priority)
}

// prazo: mais perto primeiro; sem prazo por ultimo (ISO ordena como texto)
func cmpDue(a, b Task) int {
	if a.Due == b.Due {
		return 0
	}
	if a.Due == "" {
		return 1
	}
	if b.Due == "" {
		return -1
	}
	return cmp.Compare(a.Due, b.Due)
}

// EST: mais rapidas primeiro; sem EST por ultimo
func cmpEst(a, b Task) int {
	if a.EstimateMin == b.EstimateMin {
		return 0
	}
	if a.EstimateMin == 0 {
		return 1
	}
	if b.EstimateMin == 0 {
		return -1
	}
	return cmp.Compare(a.EstimateMin, b.EstimateMin)
}

func cmpID(a, b Task) int {
	return cmp.Compare(a.ID, b.ID)
}

// SortDefault eh a ordem padrao do ls: prioridade -> prazo -> EST -> id
func SortDefault(ts []Task) {
	slices.SortFunc(ts, func(a, b Task) int {
		return cmp.Or(cmpPriority(a, b), cmpDue(a, b), cmpEst(a, b), cmpID(a, b))
	})
}

// SortByDue eh a ordem do resumo diario e do ls due: prazo -> prioridade -> EST -> id
func SortByDue(ts []Task) {
	slices.SortFunc(ts, func(a, b Task) int {
		return cmp.Or(cmpDue(a, b), cmpPriority(a, b), cmpEst(a, b), cmpID(a, b))
	})
}

// SortByCreated eh a ordem do triage: mais antigas primeiro -> id
func SortByCreated(ts []Task) {
	slices.SortFunc(ts, func(a, b Task) int {
		return cmp.Or(cmp.Compare(a.CreatedAt, b.CreatedAt), cmpID(a, b))
	})
}

// SortDoneRecent eh a ordem das concluidas: mais recentes primeiro
func SortDoneRecent(ts []Task) {
	slices.SortFunc(ts, func(a, b Task) int {
		return cmp.Or(cmp.Compare(b.DoneAt, a.DoneAt), cmp.Compare(b.ID, a.ID))
	})
}
