package project

import "testing"

func TestProgress(t *testing.T) {
	tests := []struct {
		name   string
		p      Project
		want   int
		wantOk bool
	}{
		{"CA-12.2 feitas sobre o total", Project{Done: 12, Pending: 5}, 71, true},
		{"CA-12.2 nenhuma feita", Project{Done: 0, Pending: 6}, 0, true},
		{"CA-12.2 tudo feito nao conclui", Project{Done: 4}, 100, true},
		{"CA-12.2 sem tarefas", Project{}, 0, false},
		{"CA-12.2 continuo nao tem %", Project{Continuous: true, Done: 3, Pending: 12}, 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.p.Progress()
			if got != tc.want || ok != tc.wantOk {
				t.Errorf("Progress() = %d, %v; queria %d, %v", got, ok, tc.want, tc.wantOk)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	for en, pt := range statusPT {
		got, err := StatusFromPT(pt)
		if err != nil || got != en || StatusPT(en) != pt {
			t.Errorf("situacao %s/%s nao vai e volta", en, pt)
		}
	}
	_, err := StatusFromPT("parado")
	want := "situacao parado invalida (use ativo, pausado, concluido ou cancelado)"
	if err == nil || err.Error() != want {
		t.Errorf("CA-11.5 erro = %v, queria %q", err, want)
	}
}
