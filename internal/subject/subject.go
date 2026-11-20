// Package subject define o tipo das disciplinas
package subject

import (
	"fmt"

	"github.com/CaioParanaiba/neivi/internal/codes"
)

// Reserved sao os subcomandos do subject, ja com o show do atalho futuro (RN-26)
var Reserved = []string{"add", "ls", "show", "edit", "rm"}

// ValidCode confere o codigo (RN-26)
func ValidCode(code string) error {
	return codes.Check(code, "subject", Reserved)
}

type Subject struct {
	ID      int64
	Code    string
	Name    string
	Status  string // active ou closed (no banco, em ingles)
	Pending int    // tarefas pendentes, pela disciplina efetiva (RN-14)
}

// a situacao no banco eh em ingles; na tela e na entrada, em portugues
var statusPT = map[string]string{"active": "ativa", "closed": "encerrada"}

// StatusPT mostra a situacao em portugues
func StatusPT(status string) string {
	return statusPT[status]
}

// StatusFromPT le a situacao digitada (sit:ativa, sit:encerrada)
func StatusFromPT(s string) (string, error) {
	for en, pt := range statusPT {
		if s == pt {
			return en, nil
		}
	}
	return "", fmt.Errorf("situacao %s invalida (use ativa ou encerrada)", s)
}
