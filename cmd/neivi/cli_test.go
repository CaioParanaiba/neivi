package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// bin eh o neivi compilado uma vez para todos os testes deste arquivo
var bin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "neivi-teste")
	if err != nil {
		fmt.Println("criar pasta temporaria:", err)
		os.Exit(1)
	}

	bin = filepath.Join(dir, "neivi")
	out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput()
	if err != nil {
		fmt.Printf("compilar o neivi: %v\n%s", err, out)
		os.Exit(1)
	}

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// cli roda o neivi sempre no mesmo banco temporario
type cli struct {
	t  *testing.T
	db string
}

func newCLI(t *testing.T) *cli {
	return &cli{t: t, db: filepath.Join(t.TempDir(), "teste.db")}
}

// run roda um comando e devolve a saida, os erros/avisos e o codigo de saida
func (c *cli) run(stdin string, args ...string) (string, string, int) {
	c.t.Helper()

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "NEIVI_DB="+c.db)
	cmd.Stdin = strings.NewReader(stdin)

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	code := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		c.t.Fatalf("rodar neivi %v: %v", args, err)
	}

	return strings.TrimSpace(out.String()), strings.TrimSpace(errOut.String()), code
}

// must roda um comando que precisa dar certo (preparacao do teste)
func (c *cli) must(args ...string) {
	c.t.Helper()

	_, errOut, code := c.run("", args...)
	if code != 0 {
		c.t.Fatalf("neivi %v falhou: %s", args, errOut)
	}
}

// check compara a saida e os erros/avisos com o esperado
func check(t *testing.T, gotOut, gotErr, wantOut, wantErr string) {
	t.Helper()

	if gotOut != wantOut {
		t.Errorf("saida:\n%s\nqueria:\n%s", gotOut, wantOut)
	}
	if gotErr != wantErr {
		t.Errorf("erro/aviso:\n%s\nqueria:\n%s", gotErr, wantErr)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		prep    [][]string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-01.1 so o titulo", nil,
			[]string{"add", "comprar", "cabo", "HDMI"}, "+ 1 · comprar cabo HDMI", ""},
		{"CA-01.2 sem texto", nil,
			[]string{"add"}, "", "erro: tarefa sem titulo"},
		{"CA-01.3 so marcas", nil,
			[]string{"add", "pri:A"}, "", "erro: tarefa sem titulo"},
		{"CA-01.4 duplicada avisa", [][]string{{"add", "comprar", "cabo", "HDMI"}},
			[]string{"add", "Comprar", "cabo", "hdmi", "pri:B"}, "+ 2 · B · Comprar cabo hdmi", "aviso: ja existe a tarefa 1 com esse titulo"},
		{"CA-01.5 duplicada concluida nao avisa", [][]string{{"add", "comprar", "cabo", "HDMI"}, {"done", "1"}},
			[]string{"add", "comprar", "cabo", "HDMI"}, "+ 2 · comprar cabo HDMI", ""},
		{"CA-01.6 dois pontos no titulo", nil,
			[]string{"add", "reuniao", "10:30"}, "+ 1 · reuniao 10:30", ""},
		{"CA-02.2 so alguns campos", [][]string{{"subject", "add", "BD"}},
			[]string{"add", "revisar", "slide", "@BD", "pri:B"}, "+ 1 · B · revisar slide @BD", ""},
		{"CA-02.6 disciplina nao cadastrada", nil,
			[]string{"add", "lista", "@XX", "pri:A"}, "", "erro: disciplina XX nao cadastrada (cadastrar: neivi subject add XX)"},
		{"CA-02.6 prioridade invalida", nil,
			[]string{"add", "lista", "pri:E"}, "", "erro: prioridade E invalida (use A, B, C ou D)"},
		{"CA-02.6 data invalida", nil,
			[]string{"add", "lista", "due:32/11"}, "", "erro: data 32/11 invalida (use dd/mm, dd/mm/aa, hoje, amanha ou seg...dom)"},
		{"CA-02.6 estimativa invalida", nil,
			[]string{"add", "lista", "est:abc"}, "", "erro: estimativa abc invalida (use 30m, 1h, 1h15 ou 90)"},
		{"CA-02.6 duas prioridades", nil,
			[]string{"add", "lista", "pri:A", "pri:B"}, "", "erro: mais de uma prioridade (pri:A e pri:B)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range tc.prep {
				c.must(p...)
			}

			out, errOut, _ := c.run("", tc.in...)
			check(t, out, errOut, tc.wantOut, tc.wantErr)
		})
	}

	t.Run("CA-02.6 erro nao cria nada", func(t *testing.T) {
		c := newCLI(t)
		c.run("", "add", "lista", "pri:E")

		out, _, _ := c.run("", "add", "ligar")
		if out != "+ 1 · ligar" {
			t.Errorf("depois do erro: %q, queria a tarefa 1 (nada criado antes)", out)
		}
	})
}

func TestEdit(t *testing.T) {
	// tarefa 1 sem prazo, para nao depender do dia de hoje
	prep := [][]string{
		{"subject", "add", "ED"},
		{"add", "entregar", "lista", "@ED", "pri:A", "est:1h", "--", "exercicios", "3", "a", "7"},
	}

	tests := []struct {
		name    string
		extra   [][]string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-03.3 troca o titulo", nil,
			[]string{"edit", "1", "entregar", "lista", "2"}, "~ 1 · A · 1h · entregar lista 2 @ED …", ""},
		{"CA-03.4 apagar prioridade", nil,
			[]string{"edit", "1", "pri:-"}, "~ 1 · 1h · entregar lista @ED …", ""},
		{"CA-03.4 apagar disciplina", nil,
			[]string{"edit", "1", "@-"}, "~ 1 · A · 1h · entregar lista …", ""},
		{"CA-03.4 apagar descricao", nil,
			[]string{"edit", "1", "--", "-"}, "~ 1 · A · 1h · entregar lista @ED", ""},
		{"CA-03.5 tarefa nao existe", nil,
			[]string{"edit", "99", "pri:B"}, "", "erro: tarefa 99 nao existe (veja os ids no neivi ls)"},
		{"CA-03.5 nada para mudar", nil,
			[]string{"edit", "1"}, "", "erro: nada para mudar na tarefa 1 (exemplo: neivi edit 1 pri:B due:sex)"},
		{"CA-03.5 sem id", nil,
			[]string{"edit"}, "", "erro: falta o id da tarefa (exemplo: neivi edit 1 pri:B)"},
		{"CA-03.5 id invalido", nil,
			[]string{"edit", "abc", "pri:B"}, "", "erro: id abc invalido (use o numero que aparece no ls)"},
		{"CA-03.5 marca invalida", nil,
			[]string{"edit", "1", "due:32/11"}, "", "erro: data 32/11 invalida (use dd/mm, dd/mm/aa, hoje, amanha ou seg...dom)"},
		{"CA-03.6 tarefa concluida avisa", [][]string{{"add", "comprar", "cabo"}, {"done", "2"}},
			[]string{"edit", "2", "pri:C"}, "~ 2 · C · comprar cabo", "aviso: a tarefa 2 esta concluida"},
		{"CA-03.7 titulo duplicado avisa", [][]string{{"add", "comprar", "cabo", "HDMI"}},
			[]string{"edit", "1", "comprar", "cabo", "hdmi"}, "~ 1 · A · 1h · comprar cabo hdmi @ED …", "aviso: ja existe a tarefa 2 com esse titulo"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range append(prep, tc.extra...) {
				c.must(p...)
			}

			out, errOut, _ := c.run("", tc.in...)
			check(t, out, errOut, tc.wantOut, tc.wantErr)
		})
	}

	t.Run("CA-03.5 marca invalida nao muda nada", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range prep {
			c.must(p...)
		}
		c.run("", "edit", "1", "pri:B", "due:32/11")

		out, _, _ := c.run("", "edit", "1", "est:2h")
		if out != "~ 1 · A · 2h · entregar lista @ED …" {
			t.Errorf("depois do erro: %q, queria a prioridade A ainda", out)
		}
	})
}

func TestDone(t *testing.T) {
	prep := [][]string{
		{"subject", "add", "ED"},
		{"add", "ligar"}, // 1
		{"done", "1"},    // 1 concluida hoje
		{"add", "entregar", "lista", "@ED", "pri:A"}, // 2
		{"add", "comprar", "cabo", "HDMI", "pri:B"},  // 3
		{"add", "revisar", "slide"},                  // 4
	}

	tests := []struct {
		name    string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-04.1 varias e contador",
			[]string{"done", "2", "3", "4"},
			"✓ 2 · A · entregar lista @ED\n✓ 3 · B · comprar cabo HDMI\n✓ 4 · revisar slide\nhoje: 4 concluidas", ""},
		{"CA-04.2 ontem",
			[]string{"done", "2", "3", "ontem"},
			"✓ 2 · A · entregar lista @ED (ontem)\n✓ 3 · B · comprar cabo HDMI (ontem)\nontem: 2 concluidas", ""},
		{"CA-04.4 id que nao existe segue",
			[]string{"done", "2", "99"},
			"✓ 2 · A · entregar lista @ED\nhoje: 2 concluidas", "aviso: tarefa 99 nao existe (veja os ids no neivi ls)"},
		{"CA-04.5 palavra errada recusa tudo",
			[]string{"done", "2", "otem"},
			"", `erro: "otem" nao eh um id nem "ontem" (exemplo: neivi done 3 5 ontem)`},
		{"CA-04.5 palavra no lugar do id",
			[]string{"done", "abc"},
			"", `erro: "abc" nao eh um id nem "ontem" (exemplo: neivi done 3 5 ontem)`},
		{"CA-04.5 sem id",
			[]string{"done"},
			"", "erro: falta o id da tarefa (exemplo: neivi done 3 5)"},
		{"CA-04.6 ja concluida avisa",
			[]string{"done", "1", "2"},
			"✓ 2 · A · entregar lista @ED\nhoje: 2 concluidas", "aviso: a tarefa 1 ja estava concluida"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range prep {
				c.must(p...)
			}

			out, errOut, _ := c.run("", tc.in...)
			check(t, out, errOut, tc.wantOut, tc.wantErr)
		})
	}

	t.Run("CA-04.5 recusa nao conclui nenhuma", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range prep {
			c.must(p...)
		}
		c.run("", "done", "2", "otem")

		out, _, _ := c.run("", "done", "2")
		if out != "✓ 2 · A · entregar lista @ED\nhoje: 2 concluidas" {
			t.Errorf("a 2 ja tinha sido concluida pelo comando recusado: %q", out)
		}
	})
}

func TestReopen(t *testing.T) {
	prep := [][]string{
		{"add", "ligar"}, // 1 pendente
		{"add", "comprar", "cabo", "HDMI", "pri:B"}, // 2
		{"add", "revisar", "slide"},                 // 3
		{"done", "2", "3"},
	}

	tests := []struct {
		name    string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-05.1 reabre", []string{"reopen", "2"}, "↺ 2 · B · comprar cabo HDMI", ""},
		{"CA-05.2 varias", []string{"reopen", "2", "3"}, "↺ 2 · B · comprar cabo HDMI\n↺ 3 · revisar slide", ""},
		{"CA-05.3 nao estava concluida", []string{"reopen", "1"}, "", "aviso: a tarefa 1 nao estava concluida"},
		{"CA-05.4 id que nao existe segue", []string{"reopen", "2", "99"}, "↺ 2 · B · comprar cabo HDMI", "aviso: tarefa 99 nao existe (veja os ids no neivi ls)"},
		{"CA-05.4 palavra no lugar do id", []string{"reopen", "abc"}, "", `erro: "abc" nao eh um id (exemplo: neivi reopen 5)`},
		{"CA-05.4 sem id", []string{"reopen"}, "", "erro: falta o id da tarefa (exemplo: neivi reopen 5)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range prep {
				c.must(p...)
			}

			out, errOut, _ := c.run("", tc.in...)
			check(t, out, errOut, tc.wantOut, tc.wantErr)
		})
	}
}

func TestRm(t *testing.T) {
	prep := [][]string{
		{"subject", "add", "ED"},
		{"add", "ligar"}, // 1
		{"done", "1"},    // 1 concluida
		{"add", "entregar", "lista", "@ED", "pri:A"}, // 2
		{"add", "comprar", "cabo", "HDMI", "pri:B"},  // 3
	}

	tests := []struct {
		name    string
		stdin   string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-06.1 remove com s", "s\n", []string{"rm", "2"},
			"remover 2 · A · entregar lista @ED? [s/N] ✗ 2 removida", ""},
		{"CA-06.2 N mantem", "N\n", []string{"rm", "2"},
			"remover 2 · A · entregar lista @ED? [s/N] 2 mantida", ""},
		{"CA-06.2 Enter mantem", "\n", []string{"rm", "2"},
			"remover 2 · A · entregar lista @ED? [s/N] 2 mantida", ""},
		{"CA-06.2 outra resposta mantem", "sim\n", []string{"rm", "2"},
			"remover 2 · A · entregar lista @ED? [s/N] 2 mantida", ""},
		{"CA-06.3 uma pergunta por tarefa", "s\n\n", []string{"rm", "2", "3"},
			"remover 2 · A · entregar lista @ED? [s/N] ✗ 2 removida\nremover 3 · B · comprar cabo HDMI? [s/N] 3 mantida", ""},
		{"CA-06.4 concluida avisa na pergunta", "\n", []string{"rm", "1"},
			"remover 1 · ligar (concluida, sai do historico)? [s/N] 1 mantida", ""},
		{"CA-06.5 id que nao existe segue", "\n", []string{"rm", "2", "99"},
			"remover 2 · A · entregar lista @ED? [s/N] 2 mantida", "aviso: tarefa 99 nao existe (veja os ids no neivi ls)"},
		{"CA-06.5 palavra no lugar do id", "", []string{"rm", "abc"},
			"", `erro: "abc" nao eh um id (exemplo: neivi rm 3)`},
		{"CA-06.5 sem id", "", []string{"rm"},
			"", "erro: falta o id da tarefa (exemplo: neivi rm 3)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range prep {
				c.must(p...)
			}

			out, errOut, _ := c.run(tc.stdin, tc.in...)
			check(t, out, errOut, tc.wantOut, tc.wantErr)
		})
	}

	t.Run("CA-06.1 a tarefa some", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range prep {
			c.must(p...)
		}
		c.run("s\n", "rm", "2")

		_, errOut, _ := c.run("", "done", "2")
		if errOut != "aviso: tarefa 2 nao existe (veja os ids no neivi ls)" {
			t.Errorf("depois do rm: %q, queria que a 2 nao existisse", errOut)
		}
	})
}
