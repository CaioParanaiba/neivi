// Package habit define os habitos e as regras de data deles (RN-23)
package habit

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const isoDate = "2006-01-02"

// Kind eh a periodicidade (D-08); no banco, freq_kind
const (
	Daily    = "daily"          // todo dia
	Weekdays = "weekdays"       // dias especificos
	Weekly   = "times_per_week" // N por semana
)

type Habit struct {
	ID        int64
	Code      string
	Name      string
	Goal      int // meta do dia (daily, weekdays) ou da semana (Weekly)
	Kind      string
	Days      []time.Weekday // so no Weekdays
	Status    string         // active ou archived
	CreatedAt string
}

// dias da semana na ordem do time.Weekday (domingo = 0)
var dayShort = [...]string{"dom", "seg", "ter", "qua", "qui", "sex", "sab"}
var dayLong = [...]string{"domingo", "segunda", "terca", "quarta", "quinta", "sexta", "sabado"}

// ParseDays le "seg,qua,sex" (o dias: do cadastro)
func ParseDays(s string) ([]time.Weekday, error) {
	var days []time.Weekday
	for _, d := range strings.Split(s, ",") {
		found := false
		for i, name := range dayShort {
			if d == name {
				days = append(days, time.Weekday(i))
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("dia %s invalido (use dom, seg, ter, qua, qui, sex, sab)", d)
		}
	}
	return days, nil
}

// DaysValue guarda os dias no banco como "1,3,5"; DaysFromValue faz o caminho de volta
func DaysValue(days []time.Weekday) string {
	var parts []string
	for _, d := range days {
		parts = append(parts, strconv.Itoa(int(d)))
	}
	return strings.Join(parts, ",")
}

func DaysFromValue(v string) []time.Weekday {
	var days []time.Weekday
	for _, p := range strings.Split(v, ",") {
		n, err := strconv.Atoi(p)
		if err == nil && n >= 0 && n <= 6 {
			days = append(days, time.Weekday(n))
		}
	}
	return days
}

// ValidCode confere o codigo (RN-25): nunca so numero, nunca x seguido de numero
func ValidCode(code string) error {
	if _, err := strconv.Atoi(code); err == nil {
		return fmt.Errorf("o codigo %s nao pode ser so um numero (ele se confundiria com o id)", code)
	}
	if IsQuantity(code) {
		return fmt.Errorf("o codigo %s nao pode ser x seguido de numero (ele se confundiria com a quantidade)", code)
	}
	return nil
}

// IsQuantity diz se a palavra eh uma quantidade (x2, x10)
func IsQuantity(s string) bool {
	if len(s) < 2 || s[0] != 'x' {
		return false
	}
	_, err := strconv.Atoi(s[1:])
	return err == nil
}

// StatusPT e StatusFromPT traduzem a situacao (banco em ingles, tela em portugues)
func StatusPT(status string) string {
	if status == "archived" {
		return "arquivado"
	}
	return "ativo"
}

func StatusFromPT(s string) (string, error) {
	switch s {
	case "ativo":
		return "active", nil
	case "arquivado":
		return "archived", nil
	}
	return "", fmt.Errorf("situacao %s invalida (use ativo ou arquivado)", s)
}

// WeekStart eh o domingo da semana do dia (D-19: semana de domingo a sabado)
func WeekStart(day time.Time) time.Time {
	d := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	return d.AddDate(0, 0, -int(d.Weekday()))
}

// IsDay diz se o habito vale nesse dia (so os de dias especificos tem dia fora)
func (h Habit) IsDay(day time.Time) bool {
	if h.Kind != Weekdays {
		return true
	}
	for _, d := range h.Days {
		if d == day.Weekday() {
			return true
		}
	}
	return false
}

// Period eh o intervalo (de, ate, em ISO) em que a meta vale: o dia, ou a semana do dia
func (h Habit) Period(day time.Time) (string, string) {
	if h.Kind == Weekly {
		start := WeekStart(day)
		return start.Format(isoDate), start.AddDate(0, 0, 6).Format(isoDate)
	}
	d := day.Format(isoDate)
	return d, d
}

// Missing diz se o habito "falta hoje" (RN-23): ativo, eh dia dele e nao bateu a meta do periodo
func (h Habit) Missing(done int, day time.Time) bool {
	return h.Status == "active" && h.IsDay(day) && done < h.Goal
}

// MetaLabel eh a coluna META: "todo dia", "3 todo dia", "3 por semana", "seg, qua, sex"
func (h Habit) MetaLabel() string {
	switch h.Kind {
	case Weekly:
		return fmt.Sprintf("%d por semana", h.Goal)
	case Weekdays:
		var names []string
		for _, d := range h.Days {
			names = append(names, dayShort[d])
		}
		label := strings.Join(names, ", ")
		if h.Goal > 1 {
			label = fmt.Sprintf("%d %s", h.Goal, label)
		}
		return label
	}
	if h.Goal > 1 {
		return fmt.Sprintf("%d todo dia", h.Goal)
	}
	return "todo dia"
}

// DoneLabel eh a coluna FEITO: "0/1", "2/3", "2/3 na semana", "✓" (bateu) ou "-" (nao eh dia)
func (h Habit) DoneLabel(done int, day time.Time) string {
	if !h.IsDay(day) {
		return "-"
	}
	if done >= h.Goal {
		return "✓"
	}
	label := fmt.Sprintf("%d/%d", done, h.Goal)
	if h.Kind == Weekly {
		label += " na semana"
	}
	return label
}

// Confirm eh o fim da confirmacao do habit do/undo: "feito hoje", "2/3 hoje", "2/3 nesta semana"
func (h Habit) Confirm(done int, ontem bool) string {
	when := "hoje"
	if ontem {
		when = "ontem"
	}
	if h.Kind == Weekly {
		return fmt.Sprintf("%d/%d nesta semana", done, h.Goal)
	}
	if h.Goal == 1 && done == 1 {
		return "feito " + when
	}
	return fmt.Sprintf("%d/%d %s", done, h.Goal, when)
}

// DayName eh o nome do dia para o aviso: "terca"
func DayName(day time.Time) string {
	return dayLong[day.Weekday()]
}
