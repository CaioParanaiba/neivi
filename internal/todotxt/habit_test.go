package todotxt

import (
	"reflect"
	"testing"
	"time"
)

func TestParseHabit(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		isEdit  bool
		want    HabitEntry
		wantErr string
	}{
		{"CA-17.1 todo dia", []string{"digi", "Treino", "de", "digitacao"}, false, HabitEntry{Code: "digi", Name: "Treino de digitacao"}, ""},
		{"CA-17.1 meta", []string{"cf", "meta:3"}, false, HabitEntry{Code: "cf", Goal: 3}, ""},
		{"CA-17.1 por semana", []string{"lixo", "sem:3"}, false, HabitEntry{Code: "lixo", Weekly: 3}, ""},
		{"CA-17.1 dias", []string{"ingles", "dias:seg,qua,sex"}, false,
			HabitEntry{Code: "ingles", Days: []time.Weekday{time.Monday, time.Wednesday, time.Friday}}, ""},
		{"CA-21.3 arquivar", []string{"digi", "sit:arquivado"}, true, HabitEntry{Code: "digi", Status: "archived"}, ""},
		{"CA-21.1 trocar codigo", []string{"cf", "code:cf3"}, true, HabitEntry{Code: "cf", NewCode: "cf3"}, ""},

		{"CA-17.3 sem codigo", nil, false, HabitEntry{}, "falta o codigo do habito (exemplo: neivi habit add digi Treino de digitacao)"},
		{"CA-17.3 codigo numero", []string{"3", "x"}, false, HabitEntry{}, "o codigo 3 nao pode ser so um numero (ele se confundiria com o id)"},
		{"CA-17.3 codigo x2", []string{"x2"}, false, HabitEntry{}, "o codigo x2 nao pode ser x seguido de numero (ele se confundiria com a quantidade)"},
		{"CA-17.3 sem e dias", []string{"x", "sem:3", "dias:seg"}, false, HabitEntry{}, "use sem: ou dias:, nao os dois"},
		{"CA-17.3 sem e meta", []string{"x", "sem:3", "meta:2"}, false, HabitEntry{}, "no sem: o numero ja eh a meta da semana (use so sem:3)"},
		{"CA-17.3 meta 0", []string{"x", "meta:0"}, false, HabitEntry{}, "meta 0 invalida (use um numero maior que zero)"},
		{"CA-21.6 code numero no edit", []string{"cf", "code:3"}, true, HabitEntry{}, "o codigo 3 nao pode ser so um numero (ele se confundiria com o id)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseHabit(tc.in, tc.isEdit)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Errorf("erro = %v, queria %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseHabit(%v) = %+v, queria %+v", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseHabitRefs(t *testing.T) {
	tests := []struct {
		name      string
		in        []string
		want      []HabitRef
		wantOntem bool
		wantErr   string
	}{
		{"CA-18.1 codigos e quantidade", []string{"cf", "x2", "lixo", "neivi"},
			[]HabitRef{{Code: "cf", Qty: 2}, {Code: "lixo", Qty: 1}, {Code: "neivi", Qty: 1}}, false, ""},
		{"CA-18.2 pelo id", []string{"1"}, []HabitRef{{ID: 1, Qty: 1}}, false, ""},
		{"CA-18.3 ontem em qualquer lugar", []string{"ontem", "lixo", "neivi"},
			[]HabitRef{{Code: "lixo", Qty: 1}, {Code: "neivi", Qty: 1}}, true, ""},
		{"CA-18.8 x2 antes de habito", []string{"x2", "cf"}, nil, false, "x2 precisa vir depois de um habito (exemplo: neivi habit do cf x2)"},
		{"CA-18.8 x0", []string{"cf", "x0"}, nil, false, "quantidade x0 invalida (use x1, x2...)"},
		{"CA-18.8 vazio", nil, nil, false, "falta o habito (exemplo: neivi habit do cf x2 lixo)"},
		{"so ontem", []string{"ontem"}, nil, false, "falta o habito (exemplo: neivi habit do cf x2 lixo)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ontem, err := ParseHabitRefs(tc.in, "do")
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Errorf("erro = %v, queria %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) || ontem != tc.wantOntem {
				t.Errorf("ParseHabitRefs(%v) = %+v, ontem %v; queria %+v, ontem %v", tc.in, got, ontem, tc.want, tc.wantOntem)
			}
		})
	}
}
