package task

import (
	"slices"
	"testing"
	"time"
)

// quarta-feira, 11/11/2026
var today = time.Date(2026, 11, 11, 15, 30, 0, 0, time.Local)

func TestToday(t *testing.T) {
	// as pendentes do CA-07.1
	ca071 := []Task{
		{ID: 5, Priority: "B", Due: "2026-11-09", EstimateMin: 30},
		{ID: 3, Priority: "A", Due: "2026-11-11", EstimateMin: 60},
		{ID: 7, Priority: "C", Due: "2026-11-11", EstimateMin: 30},
		{ID: 2, Priority: "A", EstimateMin: 120},
		{ID: 8, Priority: "A"},
		{ID: 9, Priority: "B"},
		{ID: 10, Priority: "A", Due: "2026-11-13"},
	}

	tests := []struct {
		name          string
		in            []Task
		late, urgent  []int64
		third         []int64
		wantImportant bool
	}{
		{"CA-07.1 tres blocos", ca071, []int64{5}, []int64{3, 7}, []int64{2, 8}, true},
		{"CA-07.2 sem A sem prazo vira sem prazo", []Task{{ID: 9, Priority: "B"}, {ID: 12, Priority: "B"}, {ID: 13, Priority: "C"}},
			nil, nil, []int64{9, 12}, false},
		{"CA-07.3 incompleta com prazo hoje eh urgente", []Task{{ID: 13, Due: "2026-11-11"}},
			nil, []int64{13}, nil, false},
		{"CA-07.4 incompleta sem prazo nao aparece", []Task{{ID: 14}},
			nil, nil, nil, false},
		{"concluida nao aparece", []Task{{ID: 1, Due: "2026-11-11", DoneAt: "2026-11-11 08:00:00"}},
			nil, nil, nil, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := Today(tc.in, today)
			if !slices.Equal(ids(s.Late), tc.late) {
				t.Errorf("atrasadas = %v, queria %v", ids(s.Late), tc.late)
			}
			if !slices.Equal(ids(s.Urgent), tc.urgent) {
				t.Errorf("urgentes = %v, queria %v", ids(s.Urgent), tc.urgent)
			}
			if !slices.Equal(ids(s.Third), tc.third) {
				t.Errorf("terceiro bloco = %v, queria %v", ids(s.Third), tc.third)
			}
			if len(tc.third) > 0 && s.Important != tc.wantImportant {
				t.Errorf("importantes = %v, queria %v", s.Important, tc.wantImportant)
			}
		})
	}

	t.Run("CA-07.6 nada para hoje", func(t *testing.T) {
		if !Today([]Task{{ID: 14}}, today).IsEmpty() {
			t.Error("queria resumo vazio")
		}
	})
}

func TestIsNewIsLate(t *testing.T) {
	tests := []struct {
		name     string
		t        Task
		wantNew  bool
		wantLate bool
	}{
		{"D-11 criada hoje", Task{CreatedAt: "2026-11-11 08:00:00"}, true, false},
		{"D-11 criada ontem", Task{CreatedAt: "2026-11-10 23:00:00"}, true, false},
		{"D-11 criada anteontem", Task{CreatedAt: "2026-11-09 10:00:00"}, false, false},
		{"atrasada: prazo ontem", Task{CreatedAt: "2026-11-01 10:00:00", Due: "2026-11-10"}, false, true},
		{"prazo hoje nao eh atrasada", Task{CreatedAt: "2026-11-01 10:00:00", Due: "2026-11-11"}, false, false},
		{"concluida nao eh atrasada", Task{CreatedAt: "2026-11-01 10:00:00", Due: "2026-11-10", DoneAt: "2026-11-10 09:00:00"}, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsNew(tc.t, today); got != tc.wantNew {
				t.Errorf("IsNew = %v, queria %v", got, tc.wantNew)
			}
			if got := IsLate(tc.t, today); got != tc.wantLate {
				t.Errorf("IsLate = %v, queria %v", got, tc.wantLate)
			}
		})
	}
}
