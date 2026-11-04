package todotxt

import (
	"testing"
	"time"
)

// segunda-feira, 02/11/2026, à tarde (como viria do time.Now)
var today = time.Date(2026, 11, 2, 15, 30, 0, 0, time.Local)

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 15, 30, 0, 0, time.Local)
}

func TestNumberAddDay(t *testing.T) {
	tests := []struct {
		name    string
		day     time.Time
		weekday time.Weekday
		want    int
	}{
		{"seg pede seg", date(2026, 11, 2), time.Monday, 7},
		{"dom pede seg", date(2026, 11, 1), time.Monday, 1},
		{"sab pede seg", date(2026, 11, 7), time.Monday, 2},
		{"seg pede qua", date(2026, 11, 2), time.Wednesday, 2},
		{"seg pede dom", date(2026, 11, 2), time.Sunday, 6},
		{"sab pede dom", date(2026, 11, 7), time.Sunday, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := numberAddDay(tc.day, tc.weekday)
			if got != tc.want {
				t.Errorf("numberAddDay(%s, %s) = %d, queria %d",
					tc.day.Format("02/01"), tc.weekday, got, tc.want)
			}
		})
	}
}

func TestEstTime(t *testing.T) {
	tests := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"90", 90, false},
		{"30m", 30, false},
		{"1h", 60, false},
		{"1h15", 75, false},
		{"1h15m", 75, false},
		{"0h30", 30, false},
		{"1h60", 120, false},

		{"", 0, true},
		{"0", 0, true},
		{"0h", 0, true},
		{"0h0", 0, true},
		{"h", 0, true},
		{"m", 0, true},
		{"abc", 0, true},
		{"-5", 0, true},
		{"1h-5", 0, true},
		{"-1h", 0, true},
		{"1x", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := estTime(tc.in)

			if tc.wantErr {
				if err == nil {
					t.Errorf("estTime(%q): esperava erro, veio %d", tc.in, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("estTime(%q): erro inesperado: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("estTime(%q) = %d, queria %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestDueDate(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		today   time.Time
		want    string
		wantErr bool
	}{
		// atalhos
		{"hoje", "hoje", today, "2026-11-02", false},
		{"amanha", "amanha", today, "2026-11-03", false},
		{"seg numa seg", "seg", today, "2026-11-09", false},
		{"segunda", "segunda", today, "2026-11-09", false},
		{"ter", "ter", today, "2026-11-03", false},
		{"qua", "quarta", today, "2026-11-04", false},
		{"qui", "qui", today, "2026-11-05", false},
		{"sex", "sexta", today, "2026-11-06", false},
		{"sab", "sab", today, "2026-11-07", false},
		{"dom", "domingo", today, "2026-11-08", false},
		{"amanha na virada do ano", "amanha", date(2026, 12, 31), "2027-01-01", false},

		// com ano
		{"ano 2 digitos", "05/11/26", today, "2026-11-05", false},
		{"dia sem zero", "5/11/26", today, "2026-11-05", false},
		{"com traco", "05-11-26", today, "2026-11-05", false},
		{"ano 4 digitos", "05/11/2026", today, "2026-11-05", false},
		{"traco e 4 digitos", "5-11-2026", today, "2026-11-05", false},
		{"data passada com ano fica", "05/11/25", today, "2025-11-05", false},
		{"29/02 em bissexto", "29/02/28", today, "2028-02-29", false},

		// sem ano
		{"ainda vai chegar", "05/11", today, "2026-11-05", false},
		{"hoje sem ano", "02/11", today, "2026-11-02", false},
		{"ja passou", "01/11", today, "2027-11-01", false},
		{"janeiro em novembro", "5-1", today, "2027-01-05", false},
		{"janeiro em dezembro", "05/01", date(2026, 12, 20), "2027-01-05", false},
		{"29/02 com proximo bissexto", "29/02", date(2027, 11, 2), "2028-02-29", false},
		{"29/02 no proprio bissexto", "29/02", date(2028, 1, 10), "2028-02-29", false},

		// erros
		{"vazio", "", today, "", true},
		{"texto", "abc", today, "", true},
		{"31/02", "31/02/26", today, "", true},
		{"mes 13", "13/13/26", today, "", true},
		{"29/02 sem bissexto com ano", "29/02/27", today, "", true},
		{"29/02 sem bissexto sem ano", "29/02", today, "", true},
		{"barras demais", "1/2/3/4", today, "", true},
		{"ano com 3 digitos", "05/11/202", today, "", true},
		{"atalho com acento", "terça", today, "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := dueDate(tc.in, tc.today)

			if tc.wantErr {
				if err == nil {
					t.Errorf("dueDate(%q): esperava erro, veio %q", tc.in, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("dueDate(%q): erro inesperado: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("dueDate(%q, hoje %s) = %q, queria %q",
					tc.in, tc.today.Format("02/01/2006"), got, tc.want)
			}
		})
	}
}

func TestParseTask(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    Entry
		wantErr bool
	}{
		{
			"so titulo",
			[]string{"revisar", "normalizacao"},
			Entry{Title: "revisar normalizacao"},
			false,
		},
		{
			"exemplo da spec",
			[]string{"revisar", "normalizacao", "@BD", "pri:A", "due:05/11", "est:30m"},
			Entry{Title: "revisar normalizacao", SubjectCode: "BD", Priority: "A", Due: "2026-11-05", EstimateMin: 30},
			false,
		},
		{
			"marcas no meio do titulo",
			[]string{"escrever", "+neivi", "o", "parser", "est:1h"},
			Entry{Title: "escrever o parser", ProjectCode: "neivi", EstimateMin: 60},
			false,
		},
		{
			"prioridade minuscula",
			[]string{"tarefa", "pri:b"},
			Entry{Title: "tarefa", Priority: "B"},
			false,
		},
		{
			"extras na ordem",
			[]string{"aula", "de", "BD", "aula:3", "rec:1", "gh:42"},
			Entry{Title: "aula de BD", Extras: "aula:3 rec:1 gh:42"},
			false,
		},
		{
			"dois pontos no titulo",
			[]string{"verificar:", "reuniao", "10:30", "https://go.dev"},
			Entry{Title: "verificar: reuniao 10:30 https://go.dev"},
			false,
		},
		{
			"disciplina e projeto",
			[]string{"trabalho", "@ED", "+p1"},
			Entry{Title: "trabalho", SubjectCode: "ED", ProjectCode: "p1"},
			false,
		},

		{"sem argumentos", []string{}, Entry{}, true},
		{"so marcas", []string{"@BD", "pri:A"}, Entry{}, true},
		{"titulo so espaco", []string{" "}, Entry{}, true},
		{"palavra vazia", []string{"tarefa", ""}, Entry{}, true},
		{"arroba sozinho", []string{"tarefa", "@"}, Entry{}, true},
		{"mais sozinho", []string{"tarefa", "+"}, Entry{}, true},
		{"duas disciplinas", []string{"tarefa", "@BD", "@ED"}, Entry{}, true},
		{"dois projetos", []string{"tarefa", "+a", "+b"}, Entry{}, true},
		{"duas prioridades", []string{"tarefa", "pri:A", "pri:B"}, Entry{}, true},
		{"prioridade D", []string{"tarefa", "pri:D"}, Entry{}, true},
		{"prioridade longa", []string{"tarefa", "pri:AB"}, Entry{}, true},
		{"prioridade vazia", []string{"tarefa", "pri:"}, Entry{}, true},
		{"duas estimativas", []string{"tarefa", "est:30", "est:1h"}, Entry{}, true},
		{"estimativa invalida", []string{"tarefa", "est:abc"}, Entry{}, true},
		{"dois prazos", []string{"tarefa", "due:hoje", "due:amanha"}, Entry{}, true},
		{"prazo invalido", []string{"tarefa", "due:31/02/26"}, Entry{}, true},
		{"prazo vazio", []string{"tarefa", "due:"}, Entry{}, true},
		{"apagar prioridade no add", []string{"tarefa", "pri:-"}, Entry{}, true},
		{"apagar prazo no add", []string{"tarefa", "due:-"}, Entry{}, true},
		{"apagar estimativa no add", []string{"tarefa", "est:-"}, Entry{}, true},
		{"apagar disciplina no add", []string{"tarefa", "@-"}, Entry{}, true},
		{"apagar projeto no add", []string{"tarefa", "+-"}, Entry{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTask(tc.in, today)

			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseTask(%q): esperava erro, veio %+v", tc.in, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseTask(%q): erro inesperado: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("ParseTask(%q)\n got: %+v\nwant: %+v", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseEdit(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    Entry
		wantErr bool
	}{
		{"so prioridade", []string{"pri:B"}, Entry{Priority: "B"}, false},
		{"so titulo", []string{"Lista", "3", "de", "ED"}, Entry{Title: "Lista 3 de ED"}, false},
		{"disciplina e prazo", []string{"@BD", "due:sex"}, Entry{SubjectCode: "BD", Due: "2026-11-06"}, false},
		{"so estimativa", []string{"est:30"}, Entry{EstimateMin: 30}, false},
		{"so extra", []string{"rec:2"}, Entry{Extras: "rec:2"}, false},
		{"apagar prioridade", []string{"pri:-"}, Entry{Priority: "-"}, false},
		{"apagar prazo", []string{"due:-"}, Entry{Due: "-"}, false},
		{"apagar estimativa", []string{"est:-"}, Entry{EstimateMin: -1}, false},
		{"apagar disciplina", []string{"@-"}, Entry{SubjectCode: "-"}, false},
		{"apagar projeto", []string{"+-"}, Entry{ProjectCode: "-"}, false},
		{"apagar extra", []string{"rec:-"}, Entry{Extras: "rec:-"}, false},
		{"titulo comeca com traco", []string{"-1", "ponto"}, Entry{Title: "-1 ponto"}, false},
		{"titulo e apagar prazo", []string{"revisar", "due:-", "pri:A"}, Entry{Title: "revisar", Due: "-", Priority: "A"}, false},

		{"sem argumentos", []string{}, Entry{}, true},
		{"apagar titulo", []string{"-"}, Entry{}, true},
		{"prioridade invalida", []string{"pri:Z"}, Entry{}, true},
		{"prazo invalido", []string{"due:31/02/26"}, Entry{}, true},
		{"prioridade e apagar", []string{"pri:A", "pri:-"}, Entry{}, true},
		{"estimativa e apagar", []string{"est:30", "est:-"}, Entry{}, true},
		{"palavra vazia", []string{""}, Entry{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseEdit(tc.in, today)

			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseEdit(%q): esperava erro, veio %+v", tc.in, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseEdit(%q): erro inesperado: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("ParseEdit(%q)\n got: %+v\nwant: %+v", tc.in, got, tc.want)
			}
		})
	}
}
