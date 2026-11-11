package todotxt

import (
	"reflect"
	"testing"

	"github.com/CaioParanaiba/neivi/internal/task"
)

func TestParseFilter(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    task.Filter
		wantErr bool
	}{
		{"CA-09.1 sem nada", nil, task.Filter{}, false},
		{"CA-09.2 disciplina", []string{"@BD"}, task.Filter{Subjects: []string{"BD"}}, false},
		{"CA-09.2 qualquer disciplina", []string{"@"}, task.Filter{AnySubject: true}, false},
		{"CA-09.2 sem disciplina", []string{"@-"}, task.Filter{NoSubject: true}, false},
		{"CA-09.2 projeto", []string{"+neivi"}, task.Filter{Projects: []string{"neivi"}}, false},
		{"CA-09.2 qualquer projeto", []string{"+"}, task.Filter{AnyProject: true}, false},
		{"CA-09.2 sem projeto", []string{"+-"}, task.Filter{NoProject: true}, false},
		{"CA-09.2 prioridade minuscula", []string{"pri:a"}, task.Filter{Priorities: []string{"A"}}, false},
		{"CA-08.2 sem prioridade", []string{"pri:-"}, task.Filter{Priorities: []string{"-"}}, false},
		{"CA-09.2 novas", []string{"new"}, task.Filter{New: true}, false},
		{"CA-09.3 mesmo tipo", []string{"@BD", "@ED"}, task.Filter{Subjects: []string{"BD", "ED"}}, false},
		{"CA-09.3 tipos diferentes", []string{"new", "@BD"}, task.Filter{New: true, Subjects: []string{"BD"}}, false},
		{"CA-09.4 palavras", []string{"cabo", "hdmi"}, task.Filter{Words: []string{"cabo", "hdmi"}}, false},
		{"CA-09.5 all", []string{"all"}, task.Filter{All: true}, false},
		{"CA-09.5 done", []string{"done"}, task.Filter{Done: true}, false},
		{"CA-09.6 desc", []string{"desc"}, task.Filter{Desc: true}, false},
		{"CA-09.7 due", []string{"due"}, task.Filter{Due: true}, false},
		{"late", []string{"late"}, task.Filter{Late: true}, false},

		{"prioridade invalida", []string{"pri:Z"}, task.Filter{}, true},
		{"all e done juntos", []string{"all", "done"}, task.Filter{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseFilter(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("erro = %v, queria erro: %v", err, tc.wantErr)
			}
			if !tc.wantErr && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseFilter(%v) = %+v, queria %+v", tc.in, got, tc.want)
			}
		})
	}
}
