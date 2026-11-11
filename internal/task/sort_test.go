package task

import (
	"slices"
	"testing"
)

func ids(ts []Task) []int64 {
	var out []int64
	for _, t := range ts {
		out = append(out, t.ID)
	}
	return out
}

// as pendentes do CA-09.1 (hoje = 11/11/2026)
func lsTasks() []Task {
	return []Task{
		{ID: 4, Title: "comprar cabo HDMI"},
		{ID: 9, Title: "estudar sorting", Priority: "B"},
		{ID: 11, Title: "parser", Priority: "B", Due: "2026-11-13"},
		{ID: 2, Title: "revisar normalizacao", Priority: "A", EstimateMin: 120},
		{ID: 3, Title: "entregar parte do feedback360", Priority: "A", Due: "2026-11-11", EstimateMin: 60},
	}
}

func TestSorts(t *testing.T) {
	tests := []struct {
		name string
		sort func([]Task)
		in   []Task
		want []int64
	}{
		{"CA-09.1 ordem padrao", SortDefault, lsTasks(), []int64{3, 2, 11, 9, 4}},
		{"CA-09.7 ordem por prazo", SortByDue, lsTasks(), []int64{3, 11, 2, 9, 4}},
		{"RN-17 EST desempata, sem EST por ultimo", SortDefault,
			[]Task{{ID: 1, Priority: "A"}, {ID: 2, Priority: "A", EstimateMin: 60}, {ID: 3, Priority: "A", EstimateMin: 30}},
			[]int64{3, 2, 1}},
		{"RN-17 id desempata no fim", SortDefault,
			[]Task{{ID: 7, Priority: "B"}, {ID: 5, Priority: "B"}}, []int64{5, 7}},
		{"CA-08.1 triage, mais antigas primeiro", SortByCreated,
			[]Task{{ID: 12, CreatedAt: "2026-11-11 08:00:00"}, {ID: 4, CreatedAt: "2026-11-05 10:00:00"}, {ID: 9, CreatedAt: "2026-11-10 09:00:00"}},
			[]int64{4, 9, 12}},
		{"CA-09.5 concluidas, mais recentes primeiro", SortDoneRecent,
			[]Task{{ID: 1, DoneAt: "2026-11-09 10:00:00"}, {ID: 2, DoneAt: "2026-11-11 08:00:00"}, {ID: 3, DoneAt: "2026-11-10 12:00:00"}},
			[]int64{2, 3, 1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.sort(tc.in)
			got := ids(tc.in)
			if !slices.Equal(got, tc.want) {
				t.Errorf("ordem = %v, queria %v", got, tc.want)
			}
		})
	}
}
