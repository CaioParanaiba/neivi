package task

import "strings"

// Filter guarda os filtros do ls (RN-18). Listas do mesmo tipo sao "ou";
// tipos diferentes sao "e".
type Filter struct {
	Subjects   []string // codigos de disciplina (@BD)
	AnySubject bool     // @ sozinho: qualquer disciplina
	NoSubject  bool     // @-: sem disciplina
	Projects   []string // codigos de projeto (+neivi)
	AnyProject bool     // + sozinho
	NoProject  bool     // +-
	Priorities []string // A a D; "-" = sem prioridade
	New        bool     // criadas hoje ou ontem (D-11)
	Late       bool     // atrasadas
	All        bool     // pendentes e, depois, concluidas
	Done       bool     // so concluidas
	Due        bool     // ordena por prazo
	Desc       bool     // mostra a descricao
	Words      []string // busca no titulo (todas precisam aparecer)
}

// IsEmpty diz se nenhum filtro que escolhe tarefas foi usado
// (Due e Desc so mudam a forma de mostrar)
func (f Filter) IsEmpty() bool {
	return len(f.Subjects) == 0 && !f.AnySubject && !f.NoSubject &&
		len(f.Projects) == 0 && !f.AnyProject && !f.NoProject &&
		len(f.Priorities) == 0 && !f.New && !f.Late && !f.All && !f.Done &&
		len(f.Words) == 0
}

// troca as letras acentuadas do portugues pela letra sem acento
var accents = strings.NewReplacer(
	"á", "a", "à", "a", "â", "a", "ã", "a", "ä", "a",
	"é", "e", "è", "e", "ê", "e", "ë", "e",
	"í", "i", "ì", "i", "î", "i", "ï", "i",
	"ó", "o", "ò", "o", "ô", "o", "õ", "o", "ö", "o",
	"ú", "u", "ù", "u", "û", "u", "ü", "u",
	"ç", "c", "ñ", "n",
)

// Fold deixa o texto pronto para a busca: minusculo e sem acento
func Fold(s string) string {
	return accents.Replace(strings.ToLower(s))
}

// MatchWords diz se o titulo tem todas as palavras, ignorando maiusculas e acentos
func MatchWords(title string, words []string) bool {
	t := Fold(title)
	for _, w := range words {
		if !strings.Contains(t, Fold(w)) {
			return false
		}
	}
	return true
}
