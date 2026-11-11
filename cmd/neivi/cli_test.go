package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// dia formata uma data relativa a hoje como o neivi mostra ("qua 11/11")
func dia(daysFromToday int) string {
	now := time.Now()
	return formatDue(now.AddDate(0, 0, daysFromToday).Format("2006-01-02"), now)
}

// entrada dd/mm/aa de uma data relativa a hoje
func entrada(daysFromToday int) string {
	return time.Now().AddDate(0, 0, daysFromToday).Format("02/01/06")
}

func TestToday(t *testing.T) {
	t.Run("CA-07.1 tres blocos, falta e contador", func(t *testing.T) {
		c := newCLI(t)
		c.must("subject", "add", "ED")
		c.must("subject", "add", "BD")
		c.must("add", "lista", "3", "@ED", "pri:B", "due:"+entrada(-2), "est:30m")     // 1
		c.must("add", "entregar", "feedback360", "@BD", "pri:A", "due:hoje", "est:1h") // 2
		c.must("add", "pagar", "boleto", "pri:C", "due:hoje", "est:30m")               // 3
		c.must("add", "revisar", "normalizacao", "@BD", "pri:A", "est:2h")             // 4
		c.must("add", "ligar", "para", "o", "dentista", "pri:A")                       // 5
		c.must("add", "estudar", "sorting", "pri:B")                                   // 6
		c.must("add", "comprar", "cabo", "pri:A", "due:"+entrada(5))                   // 7
		c.must("add", "feita")                                                         // 8
		c.must("done", "8")

		var tb table
		rows := [][]string{
			header(false),
			{"1", "B", dia(-2), "30m", "lista 3 @ED"},
			{"2", "A", dia(0), "1h", "entregar feedback360 @BD"},
			{"3", "C", dia(0), "30m", "pagar boleto"},
			{"4", "A", "", "2h", "revisar normalizacao @BD"},
			{"5", "A", "", "", "ligar para o dentista"},
		}
		for _, r := range rows {
			tb.fit(r)
		}
		h := tb.line(header(false))
		want := strings.Join([]string{
			"atrasadas (1)", h, tb.line(rows[1]), "",
			"urgentes (2)", h, tb.line(rows[2]), tb.line(rows[3]), "",
			"importantes (2)", h, tb.line(rows[4]), tb.line(rows[5]), "",
			"falta: ~4h (1 sem EST)",
			"hoje: 1 concluidas",
		}, "\n")

		out, errOut, _ := c.run("", "today")
		check(t, out, errOut, want, "")
	})

	t.Run("CA-07.2 sem A sem prazo vira sem prazo", func(t *testing.T) {
		c := newCLI(t)
		c.must("add", "estudar", "sorting", "pri:B", "est:1h")
		c.must("add", "algo", "pri:C")

		out, _, _ := c.run("", "today")
		if !strings.HasPrefix(out, "sem prazo (1)\n") || strings.Contains(out, "algo") {
			t.Errorf("queria so o bloco sem prazo com a B:\n%s", out)
		}
	})

	tests := []struct {
		name    string
		prep    [][]string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-07.6 dia livre", nil, []string{"today"},
			"nada para hoje, bom trabalho!\nhoje: 0 concluidas", ""},
		{"CA-07.7 dia livre com incompletas", [][]string{{"add", "ideia"}, {"add", "outra"}}, []string{"today"},
			"nada para hoje, bom trabalho!\nvoce tem 2 tarefas incompletas: que tal completar? (neivi triage)\nhoje: 0 concluidas", ""},
		{"CA-07.8 sem EST", [][]string{{"add", "pagar", "aluguel", "due:hoje"}}, []string{"today"},
			"", ""}, // conferido abaixo
		{"CA-07.9 argumento", nil, []string{"today", "xyz"},
			"", "erro: today nao recebe argumentos (exemplo: neivi today)"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range tc.prep {
				c.must(p...)
			}

			out, errOut, _ := c.run("", tc.in...)
			if strings.HasPrefix(tc.name, "CA-07.8") {
				// CA-07.3 e CA-07.8: incompleta com prazo hoje eh urgente; sem EST nenhuma
				if !strings.Contains(out, "urgentes (1)") || !strings.Contains(out, "falta: sem EST (1 tarefas)") {
					t.Errorf("saida:\n%s", out)
				}
				return
			}
			check(t, out, errOut, tc.wantOut, tc.wantErr)
		})
	}
}

func TestTriage(t *testing.T) {
	tests := []struct {
		name    string
		prep    [][]string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-08.1 so incompletas, mais antigas primeiro",
			[][]string{{"add", "comprar", "cabo", "HDMI"}, {"add", "pronta", "pri:B"}, {"add", "ligar"}},
			[]string{"triage"},
			"incompletas (2)\nID  P  PRAZO  EST  TAREFA\n1                  comprar cabo HDMI\n3                  ligar", ""},
		{"CA-08.3 concluida nao aparece",
			[][]string{{"add", "comprar", "cabo"}, {"add", "ligar"}, {"done", "1"}},
			[]string{"triage"},
			"incompletas (1)\nID  P  PRAZO  EST  TAREFA\n2                  ligar", ""},
		{"CA-08.4 nenhuma", [][]string{{"add", "pronta", "pri:A"}}, []string{"triage"},
			"nenhuma tarefa incompleta", ""},
		{"CA-08.5 argumento", nil, []string{"triage", "xyz"},
			"", "erro: triage nao recebe argumentos (exemplo: neivi triage)"},
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
}

func TestLs(t *testing.T) {
	// tarefas sem prazo, para a tabela nao depender do dia
	prep := [][]string{
		{"subject", "add", "BD"},
		{"subject", "add", "ED"},
		{"project", "add", "neivi"},
		{"add", "entregar", "feedback360", "@BD", "pri:A", "est:1h"}, // 1
		{"add", "revisar", "normalizacao", "@BD", "pri:A", "est:2h", "--", "slides", "4", "a", "9"}, // 2
		{"add", "parser", "+neivi", "pri:B"},          // 3
		{"add", "estudar", "sorting", "@ED", "pri:B"}, // 4
		{"add", "comprar", "cabo", "HDMI"},            // 5
		{"add", "ligar"},                              // 6
		{"done", "6"},
	}

	tests := []struct {
		name    string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-09.1 panorama",
			[]string{"ls"},
			expectTable(
				[]string{"1", "A", "", "1h", "entregar feedback360 @BD"},
				[]string{"2", "A", "", "2h", "revisar normalizacao @BD …"},
				[]string{"3", "B", "", "", "parser +neivi"},
				[]string{"4", "B", "", "", "estudar sorting @ED"},
				[]string{"5", "", "", "", "comprar cabo HDMI"},
			), ""},
		{"CA-09.5 ls done",
			[]string{"ls", "done"},
			"ID  P  PRAZO  EST  TAREFA\n6                  ligar", ""},
		{"CA-09.8 nada encontrado",
			[]string{"ls", "cabo", "azul"},
			"nenhuma tarefa com esses filtros (tente menos palavras ou filtros)", ""},
		{"CA-09.9 disciplina nao cadastrada",
			[]string{"ls", "@XX"},
			"", "erro: disciplina XX nao cadastrada (cadastrar: neivi subject add XX)"},
		{"CA-09.9 projeto nao cadastrado",
			[]string{"ls", "+XX"},
			"", "erro: projeto XX nao cadastrado (cadastrar: neivi project add XX)"},
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

	ids := func(out string) string {
		var got []string
		for _, l := range strings.Split(out, "\n")[1:] {
			if f := strings.Fields(l); len(f) > 0 && !strings.HasPrefix(l, " ") {
				got = append(got, f[0])
			}
		}
		return strings.Join(got, " ")
	}

	filters := []struct {
		name string
		in   []string
		want string
	}{
		{"CA-09.2 @BD", []string{"@BD"}, "1 2"},
		{"CA-09.2 @", []string{"@"}, "1 2 4"},
		{"CA-09.2 @-", []string{"@-"}, "3 5"},
		{"CA-09.2 +neivi", []string{"+neivi"}, "3"},
		{"CA-09.2 +", []string{"+"}, "3"},
		{"CA-09.2 +-", []string{"+-"}, "1 2 4 5"},
		{"CA-09.2 pri:A", []string{"pri:A"}, "1 2"},
		{"CA-09.2 pri:-", []string{"pri:-"}, "5"},
		{"CA-09.2 new", []string{"new"}, "1 2 3 4 5"},
		{"CA-09.3 @BD @ED", []string{"@BD", "@ED"}, "1 2 4"},
		{"CA-09.3 @BD pri:A", []string{"@BD", "pri:A"}, "1 2"},
		{"CA-09.4 busca com acento", []string{"Normalização"}, "2"},
		{"CA-09.4 todas as palavras", []string{"cabo", "hdmi"}, "5"},
		{"CA-09.5 ls all, concluidas no fim", []string{"all"}, "1 2 3 4 5 6"},
	}

	for _, tc := range filters {
		t.Run(tc.name, func(t *testing.T) {
			c := newCLI(t)
			for _, p := range prep {
				c.must(p...)
			}

			out, _, _ := c.run("", append([]string{"ls"}, tc.in...)...)
			if got := ids(out); got != tc.want {
				t.Errorf("ids = %q, queria %q\n%s", got, tc.want, out)
			}
		})
	}

	t.Run("CA-09.5 coluna X so no all", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range prep {
			c.must(p...)
		}

		out, _, _ := c.run("", "ls", "all")
		lines := strings.Split(out, "\n")
		if !strings.HasPrefix(lines[0], "ID  X  P") || !strings.HasPrefix(lines[len(lines)-1], "6   x") {
			t.Errorf("queria a coluna X e a concluida no fim:\n%s", out)
		}
	})

	t.Run("CA-09.6 ls desc", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range prep {
			c.must(p...)
		}

		out, _, _ := c.run("", "ls", "desc", "@BD")
		if strings.Contains(out, "…") || !strings.HasSuffix(out, "  slides 4 a 9") {
			t.Errorf("queria a descricao embaixo, sem o …:\n%s", out)
		}
	})
}

// expectTable monta a tabela esperada (cabecalho sem X + linhas) com o mesmo alinhamento do neivi
func expectTable(rows ...[]string) string {
	var tb table
	tb.fit(header(false))
	for _, r := range rows {
		tb.fit(r)
	}

	lines := []string{tb.line(header(false))}
	for _, r := range rows {
		lines = append(lines, tb.line(r))
	}
	return strings.Join(lines, "\n")
}

func TestShow(t *testing.T) {
	prep := [][]string{
		{"subject", "add", "BD", "Projeto", "de", "Banco", "de", "Dados"},
		{"add", "revisar", "normalizacao", "@BD", "pri:A", "due:" + entrada(2), "est:2h", "--", "slides", "4", "a", "9"}, // 1
		{"add", "comprar", "cabo"}, // 2
		{"add", "ligar"},           // 3
		{"done", "3"},
	}

	tests := []struct {
		name    string
		in      []string
		wantOut string
		wantErr string
	}{
		{"CA-10.1 todos os campos", []string{"show", "1"},
			"tarefa 1 · revisar normalizacao\n" +
				"  prioridade:  A\n" +
				"  prazo:       " + dia(2) + " (em 2 dias)\n" +
				"  EST:         2h\n" +
				"  disciplina:  BD (Projeto de Banco de Dados)\n" +
				"  projeto:     -\n" +
				"  criada:      " + dia(0) + " (nova)\n" +
				"  situacao:    pendente\n" +
				"\n" +
				"  slides 4 a 9", ""},
		{"CA-10.2 e CA-10.4 incompleta, campos vazios", []string{"show", "2"},
			"tarefa 2 · comprar cabo\n" +
				"  prioridade:  -\n" +
				"  prazo:       -\n" +
				"  EST:         -\n" +
				"  disciplina:  -\n" +
				"  projeto:     -\n" +
				"  criada:      " + dia(0) + " (nova)\n" +
				"  situacao:    incompleta", ""},
		{"CA-10.4 concluida", []string{"show", "3"},
			"tarefa 3 · ligar\n" +
				"  prioridade:  -\n" +
				"  prazo:       -\n" +
				"  EST:         -\n" +
				"  disciplina:  -\n" +
				"  projeto:     -\n" +
				"  criada:      " + dia(0) + " (nova)\n" +
				"  situacao:    concluida em " + dia(0), ""},
		{"CA-10.6 nao existe", []string{"show", "99"}, "", "erro: tarefa 99 nao existe (veja os ids no neivi ls)"},
		{"CA-10.6 id invalido", []string{"show", "abc"}, "", "erro: id abc invalido (use o numero que aparece no ls)"},
		{"CA-10.6 sem id", []string{"show"}, "", "erro: falta o id da tarefa (exemplo: neivi show 2)"},
		{"CA-10.6 mais de um id", []string{"show", "1", "2"}, "", "erro: show mostra uma tarefa por vez (exemplo: neivi show 2)"},
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

// runCases roda uma tabela de casos, cada um num banco novo com a mesma preparacao
type cliCase struct {
	name    string
	stdin   string
	in      []string
	wantOut string
	wantErr string
}

func runCases(t *testing.T, prep [][]string, tests []cliCase) {
	t.Helper()
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
}

func TestProject(t *testing.T) {
	prep := [][]string{
		{"subject", "add", "BD", "Projeto", "de", "Banco", "de", "Dados"},
		{"subject", "add", "ED"},
		{"subject", "edit", "ED", "sit:encerrada"},
	}

	runCases(t, prep, []cliCase{
		{name: "CA-11.1 todos os campos", in: []string{"project", "add", "neivi", "Meu", "sucessor", "gh:CaioParanaiba/neivi", "dir:" + os.TempDir(), "--", "em", "Go"},
			wantOut: "+ 1 · neivi · Meu sucessor · ativo"},
		{name: "CA-11.2 so o codigo", in: []string{"project", "add", "maratona"}, wantOut: "+ 1 · maratona · ativo"},
		{name: "CA-11.5 sem codigo", in: []string{"project", "add"}, wantErr: "erro: falta o codigo do projeto (exemplo: neivi project add neivi)"},
		{name: "CA-11.5 disciplina nao cadastrada", in: []string{"project", "add", "x", "@XX"}, wantErr: "erro: disciplina XX nao cadastrada (cadastrar: neivi subject add XX)"},
		{name: "CA-11.5 github invalido", in: []string{"project", "add", "x", "gh:neivi"}, wantErr: "erro: github neivi invalido (use usuario/repo ou a URL)"},
		{name: "CA-11.6 disciplina encerrada avisa", in: []string{"project", "add", "y", "@ED"},
			wantOut: "+ 1 · y · ativo", wantErr: "aviso: a disciplina ED esta encerrada (reativar: neivi subject edit ED sit:ativa)"},
		{name: "CA-11.7 pasta que nao existe avisa", in: []string{"project", "add", "z", "dir:/nao/existe"},
			wantOut: "+ 1 · z · ativo", wantErr: "aviso: a pasta /nao/existe nao existe"},
		{name: "CA-12.4 sem projetos", in: []string{"project", "ls"}, wantOut: "nenhum projeto (cadastrar: neivi project add codigo)"},
		{name: "CA-13.3 projeto nao cadastrado", in: []string{"project", "show", "xx"}, wantErr: "erro: projeto xx nao cadastrado (cadastrar: neivi project add xx)"},
		{name: "CA-13.3 sem codigo", in: []string{"project", "show"}, wantErr: "erro: falta o codigo do projeto (exemplo: neivi project show neivi)"},
	})

	full := append(prep, [][]string{
		{"project", "add", "neivi", "Meu", "sucessor", "gh:CaioParanaiba/neivi", "--", "em", "Go"}, // 1
		{"project", "add", "maratona", "cont:sim"},                                                 // 2
		{"project", "add", "feedback360", "Trabalho", "de", "BD", "@BD"},                           // 3
		{"project", "add", "ganbatte", "sit:pausado"},                                              // 4
		{"project", "add", "velho", "sit:cancelado"},                                               // 5
		{"add", "parser", "+neivi", "pri:B"},                                                       // 1
		{"add", "fatia", "2", "+neivi", "pri:A"},                                                   // 2
		{"add", "treino", "+maratona"},                                                             // 3
		{"add", "estudar", "joins", "+feedback360"},                                                // 4
		{"done", "1", "3"},
	}...)

	runCases(t, full, []cliCase{
		{name: "CA-12.1 lista: ativos com atividade primeiro, pausados no fim", in: []string{"project", "ls"},
			wantOut: expectGrid(
				[]string{"ID", "PROJETO", "SITUACAO", "DISC", "PROGRESSO", "PENDENTES", "ULTIMA"},
				[]string{"1", "neivi", "ativo", "-", "50%", "1", dia(0) + " (hoje)"},
				[]string{"2", "maratona", "ativo", "-", "continuo", "0", dia(0) + " (hoje)"},
				[]string{"3", "feedback360", "ativo", "BD", "0%", "1", "-"},
				[]string{"4", "ganbatte", "pausado", "-", "-", "0", "-"},
			)},
		{name: "CA-13.1 panorama", in: []string{"project", "show", "neivi"},
			wantOut: "projeto 1 · neivi · Meu sucessor\n" +
				"  situacao:    ativo\n" +
				"  disciplina:  -  (pessoal)\n" +
				"  github:      https://github.com/CaioParanaiba/neivi\n" +
				"  pasta:       -\n" +
				"  criado:      " + dia(0) + "\n" +
				"  progresso:   50% · 1 pendentes\n" +
				"  ultima:      " + dia(0) + " · 1 · parser  (hoje)\n" +
				"\n" +
				"  em Go\n" +
				"\n" +
				"pendentes (1)\n" +
				expectTable([]string{"2", "A", "", "", "fatia 2 +neivi"})},
		{name: "CA-14.1 mudar a situacao", in: []string{"project", "edit", "ganbatte", "sit:ativo"}, wantOut: "~ 4 · ganbatte · ativo"},
		{name: "CA-14.2 nome e codigo", in: []string{"project", "edit", "ganbatte", "code:gb", "Ganbatte", "antigo"}, wantOut: "~ 4 · gb · Ganbatte antigo · pausado"},
		{name: "CA-14.3 nada para mudar", in: []string{"project", "edit", "neivi"},
			wantErr: "erro: nada para mudar no projeto neivi (exemplo: neivi project edit neivi sit:pausado)"},
		{name: "CA-14.3 codigo que ja existe", in: []string{"project", "edit", "neivi", "code:maratona"}, wantErr: "erro: ja existe um projeto maratona"},
		{name: "CA-14.4 disciplina encerrada avisa", in: []string{"project", "edit", "neivi", "@ED"},
			wantOut: "~ 1 · neivi · Meu sucessor · ativo", wantErr: "aviso: a disciplina ED esta encerrada (reativar: neivi subject edit ED sit:ativa)"},
		{name: "CA-16.1 remover sem tarefas", stdin: "s\n", in: []string{"project", "rm", "ganbatte"},
			wantOut: "remover o projeto ganbatte (sem tarefas)? [s/N] ✗ ganbatte removido"},
		{name: "CA-16.2 projeto com tarefas recusa", in: []string{"project", "rm", "neivi"},
			wantErr: "erro: o projeto neivi tem 2 tarefas (mude as tarefas para outro com neivi edit, ou tire da lista: neivi project edit neivi sit:cancelado)"},
		{name: "CA-02.7 tarefa em projeto pausado avisa", in: []string{"add", "bug", "+ganbatte"},
			wantOut: "+ 5 · bug +ganbatte", wantErr: "aviso: o projeto ganbatte esta pausado (reativar: neivi project edit ganbatte sit:ativo)"},
		{name: "CA-02.7 tarefa em disciplina encerrada avisa", in: []string{"add", "rever", "@ED"},
			wantOut: "+ 5 · rever @ED", wantErr: "aviso: a disciplina ED esta encerrada (reativar: neivi subject edit ED sit:ativa)"},
		{name: "CA-03.8 edit para projeto pausado avisa", in: []string{"edit", "2", "+ganbatte"},
			wantOut: "~ 2 · A · fatia 2 +ganbatte", wantErr: "aviso: o projeto ganbatte esta pausado (reativar: neivi project edit ganbatte sit:ativo)"},
		{name: "CA-03.8 done nao avisa", in: []string{"done", "4"},
			wantOut: "✓ 4 · estudar joins @BD +feedback360\nhoje: 3 concluidas"},
	})

	t.Run("CA-12.3 all mostra os cancelados", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range full {
			c.must(p...)
		}
		out, _, _ := c.run("", "project", "ls", "all")
		if !strings.Contains(out, "velho") {
			t.Errorf("queria o cancelado no all:\n%s", out)
		}
	})

	t.Run("CA-10.5 disciplina herdada do projeto", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range full {
			c.must(p...)
		}
		out, _, _ := c.run("", "show", "4")
		if !strings.Contains(out, "disciplina:  BD (do projeto)") {
			t.Errorf("queria a disciplina herdada:\n%s", out)
		}
		out, _, _ = c.run("", "ls", "@BD")
		if !strings.Contains(out, "estudar joins @BD +feedback360") {
			t.Errorf("queria a herdada no ls @BD:\n%s", out)
		}
	})
}

func TestSubject(t *testing.T) {
	prep := [][]string{
		{"subject", "add", "BD", "Projeto", "de", "Banco", "de", "Dados"},
		{"subject", "add", "ED"},
		{"add", "lista", "@ED"},
		{"project", "add", "feedback360", "@BD"},
		{"add", "joins", "+feedback360"},
	}

	runCases(t, prep, []cliCase{
		{name: "CA-15.1 cadastrar", in: []string{"subject", "add", "POO"}, wantOut: "+ 3 · POO · ativa"},
		{name: "CA-15.2 so as ativas, pendentes pela efetiva", in: []string{"subject", "ls"},
			wantOut: expectGrid([]string{"ID", "DISC", "NOME", "PENDENTES"},
				[]string{"1", "BD", "Projeto de Banco de Dados", "1"},
				[]string{"2", "ED", "ED", "1"})},
		{name: "CA-15.3 encerrar", in: []string{"subject", "edit", "BD", "sit:encerrada"}, wantOut: "~ 1 · BD · Projeto de Banco de Dados · encerrada"},
		{name: "CA-15.4 codigo repetido", in: []string{"subject", "add", "BD"}, wantErr: "erro: ja existe uma disciplina BD (para mudar: neivi subject edit BD)"},
		{name: "CA-15.4 situacao invalida", in: []string{"subject", "edit", "BD", "sit:fechada"}, wantErr: "erro: situacao fechada invalida (use ativa ou encerrada)"},
		{name: "CA-15.4 nada para mudar", in: []string{"subject", "edit", "BD"},
			wantErr: "erro: nada para mudar na disciplina BD (exemplo: neivi subject edit BD sit:encerrada)"},
		{name: "CA-16.2 disciplina com tarefas recusa", in: []string{"subject", "rm", "ED"},
			wantErr: "erro: a disciplina ED tem 1 tarefas (mude as tarefas para outra com neivi edit, ou tire da lista: neivi subject edit ED sit:encerrada)"},
		{name: "CA-16.3 disciplina com projeto recusa", in: []string{"subject", "rm", "BD"},
			wantErr: "erro: a disciplina BD tem 1 projeto (feedback360)"},
	})

	t.Run("CA-15.2 all mostra as encerradas, depois", func(t *testing.T) {
		c := newCLI(t)
		for _, p := range prep {
			c.must(p...)
		}
		c.must("subject", "edit", "BD", "sit:encerrada")
		out, _, _ := c.run("", "subject", "ls", "all")
		lines := strings.Split(out, "\n")
		if len(lines) != 3 || !strings.HasPrefix(lines[1], "2   ED") || !strings.Contains(lines[2], "encerrada") {
			t.Errorf("queria a ativa e depois a encerrada:\n%s", out)
		}
	})

	t.Run("CA-16.1 remover sem nada ligado", func(t *testing.T) {
		c := newCLI(t)
		c.must("subject", "add", "BD2")
		out, _, _ := c.run("s\n", "subject", "rm", "BD2")
		if out != "remover a disciplina BD2 (sem tarefas)? [s/N] ✗ BD2 removida" {
			t.Errorf("saida: %q", out)
		}
		out, _, _ = c.run("\n", "subject", "add", "BD3")
		out, _, _ = c.run("n\n", "subject", "rm", "BD3")
		if out != "remover a disciplina BD3 (sem tarefas)? [s/N] BD3 mantida" {
			t.Errorf("saida: %q", out)
		}
	})
}

// expectGrid monta uma tabela qualquer (cabecalho e linhas) com o alinhamento do neivi
func expectGrid(head []string, rows ...[]string) string {
	var tb table
	tb.fit(head)
	for _, r := range rows {
		tb.fit(r)
	}
	lines := []string{tb.line(head)}
	for _, r := range rows {
		lines = append(lines, tb.line(r))
	}
	return strings.Join(lines, "\n")
}
