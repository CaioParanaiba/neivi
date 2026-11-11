package task

import "testing"

func TestMatchWords(t *testing.T) {
	tests := []struct {
		name  string
		title string
		words []string
		want  bool
	}{
		{"CA-09.4 acento e maiuscula", "revisar normalizacao", []string{"Normalização"}, true},
		{"CA-09.4 todas as palavras", "comprar cabo HDMI", []string{"cabo", "hdmi"}, true},
		{"CA-09.4 falta uma palavra", "comprar cabo HDMI", []string{"cabo", "azul"}, false},
		{"ordem nao importa", "comprar cabo HDMI", []string{"hdmi", "comprar"}, true},
		{"parte da palavra", "estudar sorting", []string{"sort"}, true},
		{"sem palavras passa tudo", "qualquer", nil, true},
		{"cedilha", "Fazer a lição", []string{"licao"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := MatchWords(tc.title, tc.words); got != tc.want {
				t.Errorf("MatchWords(%q, %v) = %v, queria %v", tc.title, tc.words, got, tc.want)
			}
		})
	}
}

func TestFilterIsEmpty(t *testing.T) {
	if !(Filter{Due: true, Desc: true}).IsEmpty() {
		t.Error("CA-09.8 due e desc nao escolhem tarefas: queria vazio")
	}
	if (Filter{Words: []string{"cabo"}}).IsEmpty() {
		t.Error("CA-09.8 busca escolhe tarefas: queria nao vazio")
	}
}
