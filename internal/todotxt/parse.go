package todotxt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parse(args []string, today time.Time) (Entry, error) {
	var t Entry
	var err error

	var title, desc []string
	var inDesc bool
	for _, v := range args {
		if len(v) == 0 {
			return Entry{}, fmt.Errorf("sem entrada")
		}

		//valida se os argumentos sao da descricao ou nao
		if v == "--" {
			inDesc = !inDesc
			continue
		}

		if inDesc {
			desc = append(desc, v)
			continue
		}

		switch v[0] {
		case '@':
			if t.SubjectCode != "" {
				return Entry{}, fmt.Errorf("multiplas disciplinas")
			}
			if len(v) == 1 {
				return Entry{}, fmt.Errorf("disciplina vazia")
			}

			t.SubjectCode = v[1:]
		case '+':
			if t.ProjectCode != "" {
				return Entry{}, fmt.Errorf("multiplos projetos")
			}
			if len(v) == 1 {
				return Entry{}, fmt.Errorf("projeto vazio")
			}

			t.ProjectCode = v[1:]
		default:
			key, value, found := strings.Cut(v, ":")
			if found {
				switch key {
				case "pri":
					if t.Priority != "" {
						return Entry{}, fmt.Errorf("multiplas prioridades")
					}
					if value == "" {
						return Entry{}, fmt.Errorf("prioridade vazia")
					}
					p := strings.ToUpper(value)
					if len(value) > 1 || (p != "A" && p != "B" && p != "C" && p != "-") {
						return Entry{}, fmt.Errorf("prioridade %s invalida", value)
					}

					t.Priority = p
				case "est":
					if t.EstimateMin != 0 {
						return Entry{}, fmt.Errorf("multiplas estimativas")
					}

					var est int64

					if value == "-" || value == "0" {
						est = -1
					} else {
						est, err = estTime(value)
						if err != nil {
							return Entry{}, fmt.Errorf("tempo invalido: %w", err)
						}
					}

					t.EstimateMin = est
				case "due":
					var day string
					if t.Due != "" {
						return Entry{}, fmt.Errorf("multiplas prazos")
					}
					if value == "-" {
						day = "-"
					} else {
						day, err = dueDate(value, today)
						if err != nil {
							return Entry{}, fmt.Errorf("data: %w", err)
						}
					}

					t.Due = day

				case "rec", "esf", "top", "rep", "aula", "gh", "av", "prova":
					// chaves do Ganbatte, guardadas como vieram
					if t.Extras != "" {
						t.Extras += " "
					}
					t.Extras += v
				default:
					title = append(title, v)
				}
			} else {
				title = append(title, v)
			}
		}

	}

	t.Title = strings.Join(title, " ")
	t.Description = strings.Join(desc, " ")

	return t, nil
}

func ParseTask(args []string, today time.Time) (Entry, error) {
	t, err := parse(args, today)
	if err != nil {
		return Entry{}, err
	}

	//obrigatorio
	if len(strings.TrimSpace(t.Title)) == 0 {
		return Entry{}, fmt.Errorf("tarefa sem titulo")
	}
	//opcionais, verificar "-" e -1 (de remocao)
	if t.hasClear() {
		return Entry{}, fmt.Errorf(`"-" so vale no edit`)
	}

	return t, nil
}

func ParseEdit(args []string, today time.Time) (Entry, error) {
	t, err := parse(args, today)
	if err != nil {
		return Entry{}, err
	}

	//titulo nao remove
	if t.Title == "-" {
		return Entry{}, fmt.Errorf("titulo nao pode ser removido")
	}

	if t.isEmpty() {
		return Entry{}, fmt.Errorf("edicao vazia")
	}

	return t, nil
}

func estTime(t string) (int64, error) {
	t = strings.TrimSuffix(t, "m")
	h, m, found := strings.Cut(t, "h")

	var total int64

	if found {
		n, err := strconv.Atoi(h)
		if err != nil {
			return 0, fmt.Errorf("horas invalidadas")
		}
		if n < 0 {
			return 0, fmt.Errorf("horas negativas")
		}
		total += int64(n) * 60

		if m == "" {
			n = 0
		} else {
			n, err = strconv.Atoi(m)
			if err != nil {
				return 0, fmt.Errorf("minutos invalidadas")
			}
			if n < 0 {
				return 0, fmt.Errorf("minutos negativas")
			}
		}
		total += int64(n)
	} else {
		n, err := strconv.Atoi(t)
		if err != nil {
			return 0, fmt.Errorf("minutos invalidadas")
		}
		if n < 0 {
			return 0, fmt.Errorf("minutos negativas")
		}
		total += int64(n)
	}

	if total == 0 {
		return 0, fmt.Errorf("tempo nulo")
	}
	return total, nil
}

func dueDate(d string, today time.Time) (string, error) {
	const isoDate = "2006-01-02"

	var date string
	day := today
	switch d {
	case "hoje":
		date = day.Format(isoDate)
	case "amanha":
		day = day.AddDate(0, 0, 1)
		date = day.Format(isoDate)
	case "seg", "segunda":
		nday := numberAddDay(day, time.Monday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	case "ter", "terca":
		nday := numberAddDay(day, time.Tuesday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	case "qua", "quarta":
		nday := numberAddDay(day, time.Wednesday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	case "qui", "quinta":
		nday := numberAddDay(day, time.Thursday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	case "sex", "sexta":
		nday := numberAddDay(day, time.Friday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	case "sab", "sabado":
		nday := numberAddDay(day, time.Saturday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	case "dom", "domingo":
		nday := numberAddDay(day, time.Sunday)
		day = day.AddDate(0, 0, nday)
		date = day.Format(isoDate)
	default:
		d = strings.ReplaceAll(d, "-", "/")
		n := strings.Count(d, "/")

		if n == 2 {
			//com ano
			parsed, err := time.ParseInLocation("2/1/06", d, today.Location())
			if err != nil {
				parsed, err = time.ParseInLocation("2/1/2006", d, today.Location())
				if err != nil {
					return "", fmt.Errorf("data %s invalida", d)
				}
			}
			return parsed.Format(isoDate), nil
		} else if n == 1 {
			//sem ano
			parsed, err := time.ParseInLocation("2/1", d, today.Location())
			if err != nil {
				return "", fmt.Errorf("data %s invalida", d)
			}

			// o Parse devolve o ano 0: usa este ano, ou o próximo se o dia já passou.
			// Decide o ano antes de montar a data, senão o 29/02 transborda para 01/03.
			year := today.Year()
			if parsed.Month() < today.Month() || (parsed.Month() == today.Month() && parsed.Day() < today.Day()) {
				year++
			}

			next := time.Date(year, parsed.Month(), parsed.Day(), 0, 0, 0, 0, today.Location())
			if next.Day() != parsed.Day() {
				// só acontece com 29/02 em ano que não é bissexto
				return "", fmt.Errorf("data %s nao existe em %d", d, year)
			}
			return next.Format(isoDate), nil
		} else {
			//erro
			return "", fmt.Errorf("data %s invalida", d)
		}

	}

	return date, nil
}

func numberAddDay(day time.Time, weekday time.Weekday) int {
	nday := int(day.Weekday())
	nday = (7 - nday + int(weekday)) % 7
	if nday == 0 {
		nday = 7
	}

	return nday
}
