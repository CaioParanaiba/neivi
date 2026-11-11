package todotxt

import "testing"

func TestParseProject(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		isEdit  bool
		want    ProjectEntry
		wantErr string
	}{
		{"CA-11.1 todos os campos",
			[]string{"neivi", "Meu", "sucessor", "gh:CaioParanaiba/neivi", "dir:~/Code/neivi", "--", "em", "Go"}, false,
			ProjectEntry{Code: "neivi", Name: "Meu sucessor", GitHub: "https://github.com/CaioParanaiba/neivi", Dir: "~/Code/neivi", Description: "em Go"}, ""},
		{"CA-11.2 so o codigo", []string{"maratona"}, false, ProjectEntry{Code: "maratona"}, ""},
		{"CA-11.3 academico", []string{"feedback360", "Trabalho", "@BD"}, false,
			ProjectEntry{Code: "feedback360", Name: "Trabalho", SubjectCode: "BD"}, ""},
		{"CA-11.3 continuo", []string{"maratona", "cont:sim"}, false, ProjectEntry{Code: "maratona", Continuous: "sim"}, ""},
		{"CA-11.4 URL inteira", []string{"x", "gh:https://gitlab.com/a/b"}, false, ProjectEntry{Code: "x", GitHub: "https://gitlab.com/a/b"}, ""},
		{"CA-14.1 situacao", []string{"ganbatte", "sit:pausado"}, true, ProjectEntry{Code: "ganbatte", Status: "paused"}, ""},
		{"CA-14.2 apagar e trocar codigo", []string{"ganbatte", "dir:-", "gh:-", "@-", "code:gb"}, true,
			ProjectEntry{Code: "ganbatte", Dir: "-", GitHub: "-", SubjectCode: "-", NewCode: "gb"}, ""},
		{"dois pontos no nome", []string{"x", "parte:1"}, false, ProjectEntry{Code: "x", Name: "parte:1"}, ""},

		{"CA-11.5 sem codigo", nil, false, ProjectEntry{}, "falta o codigo do projeto (exemplo: neivi project add neivi)"},
		{"CA-11.5 github sem barra", []string{"x", "gh:neivi"}, false, ProjectEntry{}, "github neivi invalido (use usuario/repo ou a URL)"},
		{"CA-11.5 situacao invalida", []string{"x", "sit:parado"}, false, ProjectEntry{}, "situacao parado invalida (use ativo, pausado, concluido ou cancelado)"},
		{"CA-11.5 continuo invalido", []string{"x", "cont:talvez"}, false, ProjectEntry{}, "cont:talvez invalido (use cont:sim ou cont:nao)"},
		{"apagar no add", []string{"x", "dir:-"}, false, ProjectEntry{}, `"-" so vale no edit (para apagar um campo: neivi project edit x dir:-)`},
		{"code: no add", []string{"x", "code:y"}, false, ProjectEntry{}, "code: so vale no edit (no add, o codigo eh a primeira palavra)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseProject(tc.in, tc.isEdit)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Errorf("erro = %v, queria %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseProject(%v) = %+v, queria %+v", tc.in, got, tc.want)
			}
		})
	}

	t.Run("CA-14.3 edit vazio", func(t *testing.T) {
		e, _ := ParseProject([]string{"neivi"}, true)
		if !e.IsEmpty() {
			t.Error("queria vazio")
		}
	})
}

func TestParseSubject(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		isEdit  bool
		want    SubjectEntry
		wantErr string
	}{
		{"CA-15.1 codigo e nome", []string{"BD", "Projeto", "de", "Banco"}, false, SubjectEntry{Code: "BD", Name: "Projeto de Banco"}, ""},
		{"CA-15.3 encerrar", []string{"BD", "sit:encerrada"}, true, SubjectEntry{Code: "BD", Status: "closed"}, ""},
		{"CA-15.3 trocar codigo", []string{"BD", "code:PBD"}, true, SubjectEntry{Code: "BD", NewCode: "PBD"}, ""},
		{"CA-15.4 sem codigo", nil, false, SubjectEntry{}, "falta o codigo da disciplina (exemplo: neivi subject add BD)"},
		{"CA-15.4 situacao invalida", []string{"BD", "sit:fechada"}, true, SubjectEntry{}, "situacao fechada invalida (use ativa ou encerrada)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseSubject(tc.in, tc.isEdit)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Errorf("erro = %v, queria %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseSubject(%v) = %+v, queria %+v", tc.in, got, tc.want)
			}
		})
	}
}
