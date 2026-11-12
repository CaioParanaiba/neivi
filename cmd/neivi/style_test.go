package main

import (
	"regexp"
	"strings"
	"testing"
)

var ansiCode = regexp.MustCompile("\033\\[[0-9;]*m")

// withColor liga as cores so durante o teste (o go test nao roda num terminal)
func withColor(t *testing.T) {
	t.Helper()
	old := useColor
	useColor = true
	t.Cleanup(func() { useColor = old })
}

func TestPaintOff(t *testing.T) {
	old := useColor
	useColor = false
	defer func() { useColor = old }()

	if got := paint(stSection, "urgentes (3)"); got != "urgentes (3)" {
		t.Errorf("sem cor, paint mudou o texto: %q", got)
	}
	if got := paintHelp(helps["add"]); got != helps["add"] {
		t.Errorf("sem cor, paintHelp mudou o texto")
	}
}

func TestPaintedLineAlign(t *testing.T) {
	withColor(t)

	var tb table
	cells := []string{"58", "B", "", "20m", "Lab 1 +p1 …"}
	tb.fit(header(false))
	tb.fit(cells)

	plain := tb.line(cells)
	painted := tb.paintedLine(cells, paintRow(cells, false))
	if !strings.Contains(painted, "\033[") {
		t.Fatalf("linha sem cor: %q", painted)
	}
	// sem os codigos, a linha pintada eh igual a de texto puro (mesmas colunas)
	if got := ansiCode.ReplaceAllString(painted, ""); got != plain {
		t.Errorf("desalinhou:\n%q\n%q", got, plain)
	}
}

func TestPaintHelp(t *testing.T) {
	withColor(t)

	// o texto continua o mesmo, so com codigos
	for name, text := range helps {
		if got := ansiCode.ReplaceAllString(paintHelp(text), ""); got != text {
			t.Errorf("help %s mudou o texto:\n%s", name, got)
		}
	}
	if got := ansiCode.ReplaceAllString(paintHelp(help), ""); got != help {
		t.Errorf("help geral mudou o texto:\n%s", got)
	}

	out := paintHelp(helps["add"])
	for _, want := range []string{
		paint(stSection, "uso:"),
		paint(stSection, "exemplos:"),
		paint(stSubject, "@ED"),           // marca no exemplo
		paint(priStyle("A"), "pri:A"),     // marca na tabela
		paint(stBold, "due:sex") + "  pr", // coluna separada por 2 espacos
		paint(stDim, "neivi") + " add",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("help add sem %q:\n%q", want, out)
		}
	}

	if !strings.Contains(paintHelp(help), paint(stBold, "habit")) {
		t.Errorf("help geral: habit sem negrito")
	}
}

func TestNoColorInPipe(t *testing.T) {
	c := newCLI(t)
	c.must("add", "entregar lista pri:A due:hoje")
	for _, args := range [][]string{{"today"}, {"ls"}, {"show", "1"}, {"help"}, {"help", "add"}} {
		out, _, _ := c.run("", args...)
		if strings.Contains(out, "\033") {
			t.Errorf("neivi %v pintou fora do terminal: %q", args, out)
		}
	}
}
