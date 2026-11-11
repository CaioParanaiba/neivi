package main

import (
	"testing"
	"time"

	"github.com/CaioParanaiba/neivi/internal/task"
)

// terca-feira, 10/11/2026: o "hoje" dos criterios de aceitacao
var today = time.Date(2026, 11, 10, 15, 30, 0, 0, time.Local)

func TestFormatDue(t *testing.T) {
	tests := []struct {
		name string
		due  string
		want string
	}{
		{"CA-02.1 sexta deste ano", "2026-11-13", "sex 13/11"},
		{"CA-02.3 proxima terca", "2026-11-17", "ter 17/11"},
		{"CA-02.4 ano seguinte mostra o ano", "2027-11-09", "ter 09/11/27"},
		{"CA-02.5 ontem, deste ano", "2026-11-09", "seg 09/11"},
		{"CA-03.1 segunda", "2026-11-16", "seg 16/11"},
		{"CA-03.2 amanha", "2026-11-11", "qua 11/11"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatDue(tc.due, today)
			if got != tc.want {
				t.Errorf("formatDue(%q) = %q, queria %q", tc.due, got, tc.want)
			}
		})
	}
}

func TestFormatLine(t *testing.T) {
	lista := task.Task{
		ID: 1, Title: "entregar lista", Priority: "A", Due: "2026-11-13",
		EstimateMin: 60, SubjectCode: "ED", Description: "exercicios 3 a 7",
	}

	tests := []struct {
		name   string
		symbol string
		t      task.Task
		want   string
	}{
		{"CA-01.1 so titulo", "+", task.Task{ID: 1, Title: "comprar cabo HDMI"}, "+ 1 · comprar cabo HDMI"},
		{"CA-01.6 dois pontos no titulo", "+", task.Task{ID: 1, Title: "reuniao 10:30"}, "+ 1 · reuniao 10:30"},
		{"CA-02.1 todos os campos", "+", lista, "+ 1 · A · sex 13/11 · 1h · entregar lista @ED …"},
		{"CA-02.2 so alguns campos", "+", task.Task{ID: 1, Title: "revisar slide", Priority: "B", SubjectCode: "BD"}, "+ 1 · B · revisar slide @BD"},
		{"CA-02.4 prazo de outro ano", "+", task.Task{ID: 1, Title: "c", Due: "2027-11-09"}, "+ 1 · ter 09/11/27 · c"},
		{"projeto", "+", task.Task{ID: 2, Title: "parser", ProjectCode: "neivi"}, "+ 2 · parser +neivi"},
		{"CA-03.1 edicao", "~", task.Task{ID: 4, Title: "entregar lista", Priority: "A", Due: "2026-11-16", EstimateMin: 60, SubjectCode: "ED", Description: "x"}, "~ 4 · A · seg 16/11 · 1h · entregar lista @ED …"},
		{"CA-04.1 conclusao", "✓", task.Task{ID: 3, Title: "entregar lista", Priority: "A", SubjectCode: "ED"}, "✓ 3 · A · entregar lista @ED"},
		{"CA-05.1 reabrir", "↺", task.Task{ID: 5, Title: "comprar cabo HDMI", Priority: "B"}, "↺ 5 · B · comprar cabo HDMI"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatLine(tc.symbol, tc.t, today)
			if got != tc.want {
				t.Errorf("formatLine = %q, queria %q", got, tc.want)
			}
		})
	}
}

func TestParseIDs(t *testing.T) {
	tests := []struct {
		name       string
		in         []string
		allowOntem bool
		wantIDs    int
		wantOntem  bool
		wantErr    string
	}{
		{"CA-04.1 varios ids", []string{"3", "5", "8"}, true, 3, false, ""},
		{"CA-04.2 com ontem", []string{"3", "5", "ontem"}, true, 2, true, ""},
		{"ontem em qualquer posicao", []string{"ontem", "3"}, true, 1, true, ""},
		{"CA-04.5 palavra errada recusa tudo", []string{"3", "otem"}, true, 0, false, `"otem" nao eh um id nem "ontem" (exemplo: neivi done 3 5 ontem)`},
		{"CA-04.5 so ontem, sem id", []string{"ontem"}, true, 0, false, "falta o id da tarefa (exemplo: neivi done 3 5 ontem)"},
		{"CA-05.4 reopen nao aceita ontem", []string{"5", "ontem"}, false, 0, false, `"ontem" nao eh um id (exemplo: neivi reopen 5)`},
		{"CA-05.4 sem id", []string{}, false, 0, false, "falta o id da tarefa (exemplo: neivi reopen 5)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			example := "neivi reopen 5"
			if tc.allowOntem {
				example = "neivi done 3 5 ontem"
			}

			ids, ontem, err := parseIDs(tc.in, tc.allowOntem, example)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Errorf("erro = %v, queria %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if len(ids) != tc.wantIDs || ontem != tc.wantOntem {
				t.Errorf("parseIDs(%v) = %v, ontem %v; queria %d ids, ontem %v", tc.in, ids, ontem, tc.wantIDs, tc.wantOntem)
			}
		})
	}
}
