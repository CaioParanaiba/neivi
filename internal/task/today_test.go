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

func TestTomorrowEnd(t *testing.T) {
	tests := []struct {
		name string
		day  time.Time
		want string
	}{
		{"US-07 quarta ve a quinta", today, "2026-11-12"},
		{"US-07 quinta ve a sexta", today.AddDate(0, 0, 1), "2026-11-13"},
		{"US-07 sexta ve ate segunda", today.AddDate(0, 0, 2), "2026-11-16"},
		{"US-07 sabado ve ate segunda", today.AddDate(0, 0, 3), "2026-11-16"},
		{"US-07 domingo ve a segunda", today.AddDate(0, 0, 4), "2026-11-16"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := TomorrowEnd(tc.day).Format(isoDate); got != tc.want {
				t.Errorf("TomorrowEnd = %s, queria %s", got, tc.want)
			}
		})
	}
}

func TestTodayTomorrowAndProjects(t *testing.T) {
	in := []Task{
		{ID: 211, Priority: "B", Due: "2026-11-12", ProjectID: 4, ProjectCode: "bd"}, // amanha, mesmo de projeto
		{ID: 212, Priority: "C", Due: "2026-11-12", EstimateMin: 30},
		{ID: 300, Priority: "A", Due: "2026-11-13"}, // depois de amanha: fora
		{ID: 55, Priority: "B"},                                                      // avulsa: fica
		{ID: 58, Priority: "B", ProjectID: 1, ProjectCode: "p1"},                     // de projeto: so conta
		{ID: 59, Priority: "B", ProjectID: 1, ProjectCode: "p1"},                     // de projeto: so conta
		{ID: 215, Priority: "B", ProjectID: 3, ProjectCode: "mrt"},                   // de projeto: so conta
		{ID: 216, Priority: "C", ProjectID: 3, ProjectCode: "mrt"},                   // de projeto: so conta
		{ID: 217, ProjectID: 3, ProjectCode: "mrt"},                                  // incompleta: nem conta
		{ID: 218, Priority: "A", ProjectID: 3, ProjectCode: "mrt"},                   // A ganha do projeto
		{ID: 219, Priority: "A", Due: "2026-11-11", ProjectID: 1, ProjectCode: "p1"}, // prazo ganha do projeto
	}

	s := Today(in, today)
	if got := ids(s.Tomorrow); !slices.Equal(got, []int64{211, 212}) {
		t.Errorf("amanha = %v, queria [211 212]", got)
	}
	if got := ids(s.Urgent); !slices.Equal(got, []int64{219}) {
		t.Errorf("urgentes = %v, queria [219]", got)
	}
	if got := ids(s.Third); !slices.Equal(got, []int64{218}) || !s.Important {
		t.Errorf("terceiro bloco = %v (importantes %v), queria so a A de projeto", got, s.Important)
	}
	want := []ProjectCount{{"mrt", 2}, {"p1", 2}}
	if !slices.Equal(s.Hidden, want) || s.HiddenTotal() != 4 {
		t.Errorf("de projetos = %v (%d), queria %v (4)", s.Hidden, s.HiddenTotal(), want)
	}

	t.Run("sem A: o sem prazo fica so com as avulsas", func(t *testing.T) {
		s := Today(in[3:8], today)
		if got := ids(s.Third); !slices.Equal(got, []int64{55}) || s.Important {
			t.Errorf("sem prazo = %v, queria [55]", got)
		}
	})

	t.Run("so amanha: hoje continua vazio", func(t *testing.T) {
		s := Today([]Task{{ID: 1, Priority: "B", Due: "2026-11-12"}}, today)
		if !s.IsEmpty() || len(s.Tomorrow) != 1 {
			t.Errorf("queria hoje vazio e 1 amanha: %+v", s)
		}
	})
}
