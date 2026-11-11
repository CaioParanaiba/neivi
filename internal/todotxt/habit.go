package todotxt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CaioParanaiba/neivi/internal/habit"
)

// HabitEntry eh o que veio no habit add/edit; zero/"" = nao veio
type HabitEntry struct {
	Code    string
	Name    string
	NewCode string
	Goal    int            // meta:3
	Weekly  int            // sem:3 (a meta da semana)
	Days    []time.Weekday // dias:seg,qua,sex
	Status  string         // ja em ingles (active, archived)
}

// ParseHabit le o habit add (isEdit = false) ou o habit edit (isEdit = true)
func ParseHabit(args []string, isEdit bool) (HabitEntry, error) {
	var e HabitEntry
	if len(args) == 0 || args[0] == "" {
		return HabitEntry{}, fmt.Errorf("falta o codigo do habito (exemplo: neivi habit add digi Treino de digitacao)")
	}
	e.Code = args[0]
	if !isEdit {
		err := habit.ValidCode(e.Code)
		if err != nil {
			return HabitEntry{}, err
		}
	}

	var name []string
	for _, v := range args[1:] {
		key, value, found := strings.Cut(v, ":")
		if !found {
			name = append(name, v)
			continue
		}

		var err error
		switch key {
		case "meta":
			e.Goal, err = positive("meta", value)
		case "sem":
			e.Weekly, err = positive("sem", value)
		case "dias":
			e.Days, err = habit.ParseDays(value)
		case "sit":
			e.Status, err = habit.StatusFromPT(value)
		case "code":
			if !isEdit {
				err = fmt.Errorf("code: so vale no edit (no add, o codigo eh a primeira palavra)")
				break
			}
			err = habit.ValidCode(value)
			e.NewCode = value
		default:
			name = append(name, v)
		}
		if err != nil {
			return HabitEntry{}, err
		}
	}

	// um habito tem uma periodicidade so (RN-23)
	if e.Weekly > 0 && len(e.Days) > 0 {
		return HabitEntry{}, fmt.Errorf("use sem: ou dias:, nao os dois")
	}
	if e.Weekly > 0 && e.Goal > 0 {
		return HabitEntry{}, fmt.Errorf("no sem: o numero ja eh a meta da semana (use so sem:%d)", e.Weekly)
	}

	e.Name = strings.Join(name, " ")
	return e, nil
}

// positive le um numero maior que zero (meta:3, sem:3)
func positive(key, value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%s %s invalida (use um numero maior que zero)", key, value)
	}
	return n, nil
}

// IsEmpty diz se o edit nao trouxe nada para mudar
func (e HabitEntry) IsEmpty() bool {
	return e.Name == "" && e.NewCode == "" && e.Goal == 0 && e.Weekly == 0 && len(e.Days) == 0 && e.Status == ""
}

// HabitRef eh um habito no habit do/undo: pelo codigo ou pelo id, com a quantidade
type HabitRef struct {
	Code string // vazio quando veio pelo id
	ID   int64
	Qty  int
}

// ParseHabitRefs le os argumentos do habit do/undo (RN-25): codigos ou ids,
// x2 depois de um habito e ontem em qualquer lugar. Le tudo antes de agir (RN-12).
func ParseHabitRefs(args []string, cmd string) ([]HabitRef, bool, error) {
	var refs []HabitRef
	var ontem bool

	for _, v := range args {
		switch {
		case v == "ontem":
			ontem = true
		case habit.IsQuantity(v):
			n, _ := strconv.Atoi(v[1:])
			if n <= 0 {
				return nil, false, fmt.Errorf("quantidade %s invalida (use x1, x2...)", v)
			}
			if len(refs) == 0 {
				return nil, false, fmt.Errorf("%s precisa vir depois de um habito (exemplo: neivi habit %s cf x2)", v, cmd)
			}
			refs[len(refs)-1].Qty = n
		default:
			id, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				refs = append(refs, HabitRef{ID: id, Qty: 1})
			} else {
				refs = append(refs, HabitRef{Code: v, Qty: 1})
			}
		}
	}

	if len(refs) == 0 {
		return nil, false, fmt.Errorf("falta o habito (exemplo: neivi habit %s cf x2 lixo)", cmd)
	}

	return refs, ontem, nil
}
