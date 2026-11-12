package main

// Cores da tela. Arquivo so de apresentacao: nenhuma regra do neivi depende dele.

import (
	"os"
	"regexp"
	"strings"
)

// colorOn liga as cores. Se um dia as cores derem problema, troque para false e
// recompile: a saida volta a ser texto puro, igual a de antes das cores.
const colorOn = true

// useColor diz se esta execucao pinta. Alem da constante, so pinta quando a
// saida eh um terminal (pipe, arquivo e os testes recebem texto puro) e quando
// NO_COLOR nao esta definida (convencao: NO_COLOR=1 neivi today).
var useColor = colorOn && os.Getenv("NO_COLOR") == "" && isTerminal(os.Stdout)

// isTerminal diz se o arquivo eh um terminal (dispositivo de caractere)
func isTerminal(f *os.File) bool {
	st, err := f.Stat()
	return err == nil && st.Mode()&os.ModeCharDevice != 0
}

// style eh o codigo ANSI de um papel na tela. O numero depois de 38;5; eh a cor
// da paleta do terminal (0 a 15); o tom real vem do terminal, que o wallust
// troca junto com o wallpaper. Por isso aqui nao tem nome de cor, so papel.
type style string

const (
	stSection style = "1;38;5;6" // titulo de bloco: urgentes (3), tarefas:
	stDim     style = "38;5;8"   // o que eh secundario: ID, cabecalho, dicas, …
	stSubject style = "1;38;5;4" // @disciplina
	stProject style = "1;38;5;2" // +projeto
	stDue     style = "38;5;4"   // prazo
	stEst     style = "38;5;1"   // EST
	stFooter  style = "38;5;3"   // falta: e hoje:
	stBold    style = "1"        // o que se digita, na ajuda
)

// priStyle eh o estilo de cada prioridade; D e sem prioridade ficam apagados
func priStyle(p string) style {
	switch p {
	case "A":
		return "1;38;5;5"
	case "B":
		return "38;5;3"
	case "C":
		return "38;5;1"
	}
	return stDim
}

// paint envolve o texto no estilo; sem cor (ou texto vazio) devolve igual
func paint(st style, s string) string {
	if !useColor || s == "" || st == "" {
		return s
	}
	return "\033[" + string(st) + "m" + s + "\033[0m"
}

// markStyle eh o estilo de uma marca (@BD, +neivi, pri:A); vazio se nao for marca
func markStyle(tok string) style {
	switch {
	case strings.HasPrefix(tok, "@"):
		return stSubject
	case strings.HasPrefix(tok, "+"):
		return stProject
	case len(tok) == 5 && strings.HasPrefix(tok, "pri:"):
		return priStyle(tok[4:])
	}
	return ""
}

// paintWords pinta as marcas de um texto; as outras palavras ficam com o estilo base
func paintWords(s string, base style) string {
	words := strings.Split(s, " ")
	for i, w := range words {
		if st := markStyle(w); st != "" {
			words[i] = paint(st, w)
		} else {
			words[i] = paint(base, w)
		}
	}
	return strings.Join(words, " ")
}

// paintTitle eh a coluna TAREFA com cor; mesmo texto do row (titulo @disc +proj …)
func paintTitle(title string) string {
	return paintWords(strings.ReplaceAll(title, " …", " "+paint(stDim, "…")), "")
}

// paintRow pinta as celulas de uma tarefa (as de row); as larguras continuam
// vindo do texto puro, no tb.paintedLine
func paintRow(cells []string, withX bool) []string {
	p := make([]string, len(cells))
	copy(p, cells)
	i := 0
	p[i] = paint(stDim, p[i]) // ID
	i++
	if withX {
		p[i] = paint(stDim, p[i])
		i++
	}
	p[i] = paint(priStyle(p[i]), p[i])
	p[i+1] = paint(stDue, p[i+1])
	p[i+2] = paint(stEst, p[i+2])
	p[i+3] = paintTitle(p[i+3])
	return p
}

var helpCols = regexp.MustCompile(` {3,}`)
var helpTwo = regexp.MustCompile(` {2,}`)

// paintHelp pinta o texto da ajuda linha a linha: titulos de secao, o que se
// digita em negrito, marcas com as cores do today e dicas apagadas
func paintHelp(text string) string {
	if !useColor {
		return text
	}
	lines := strings.Split(text, "\n")
	examples := false
	prevBlank := false
	for i, l := range lines {
		lines[i], examples = paintHelpLine(l, examples, prevBlank)
		prevBlank = l == ""
	}
	return strings.Join(lines, "\n")
}

// paintHelpLine pinta uma linha; examples diz se esta dentro de "exemplos:"
func paintHelpLine(l string, examples, prevBlank bool) (string, bool) {
	if l == "" {
		return l, examples
	}
	body := strings.TrimLeft(l, " ")
	pad := l[:len(l)-len(body)]

	if pad == "" {
		switch {
		case strings.HasPrefix(l, "uso:"):
			return paint(stSection, "uso:") + l[4:], false
		case strings.HasSuffix(l, ":"):
			return paint(stSection, l), strings.HasPrefix(l, "exemplo")
		}
		return paint(stDim, l), false // dica solta
	}

	// cabecalho de tabela: mais fundo e logo depois de linha vazia (help pri)
	if len(pad) >= 3 && prevBlank {
		return paint(stDim, l), examples
	}

	if examples && strings.HasPrefix(body, "neivi ") {
		return pad + paint(stDim, "neivi") + " " + paintWords(body[6:], ""), examples
	}

	// linha de prioridade: "A  definicao..." (help pri) e "-  apaga..." (help edit)
	if len(body) > 3 && strings.ContainsRune("ABCD-", rune(body[0])) && body[1:3] == "  " {
		return pad + paint(priStyle(body[:1]), body[:1]) + paintWords(body[1:], ""), examples
	}

	parts := splitHelpCols(body)
	if len(parts) == 1 {
		if strings.HasPrefix(body, "neivi ") {
			return pad + paintSyntax(body), examples
		}
		return l, examples // frase: sem cor
	}

	// colunas alternam: o que se digita (negrito) e a explicacao (normal)
	var b strings.Builder
	b.WriteString(pad)
	for k, p := range parts {
		switch {
		case k%2 == 1:
			b.WriteString(p) // espacos entre colunas
		case (k/2)%2 == 0:
			b.WriteString(paintSyntax(p))
		default:
			b.WriteString(p)
		}
	}
	return b.String(), examples
}

// splitHelpCols separa as colunas (3 ou mais espacos) mantendo os espacos nas
// posicoes impares; se a ultima coluna for "o que se digita" sem explicacao
// depois (due:sex  prazo...), ela eh separada nos 2 espacos
func splitHelpCols(s string) []string {
	var parts []string
	last := 0
	for _, m := range helpCols.FindAllStringIndex(s, -1) {
		parts = append(parts, s[last:m[0]], s[m[0]:m[1]])
		last = m[1]
	}
	parts = append(parts, s[last:])

	if (len(parts)/2)%2 == 0 { // ultima coluna eh de sintaxe
		tail := parts[len(parts)-1]
		if m := helpTwo.FindStringIndex(tail); m != nil {
			parts = append(parts[:len(parts)-1], tail[:m[0]], tail[m[0]:m[1]], tail[m[1]:])
		}
	}
	return parts
}

// paintSyntax pinta o que se digita: negrito, marcas na cor delas, "neivi" apagado
func paintSyntax(s string) string {
	if strings.HasPrefix(s, "neivi ") {
		return paint(stDim, "neivi") + " " + paintWords(s[6:], stBold)
	}
	return paintWords(s, stBold)
}
