package habit

import (
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 15, 30, 0, 0, time.Local)
}

// quarta-feira, 11/11/2026
var wed = date(2026, 11, 11)

var (
	digi   = Habit{Code: "digi", Goal: 1, Kind: Daily, Status: "active"}
	cf     = Habit{Code: "cf", Goal: 3, Kind: Daily, Status: "active"}
	lixo   = Habit{Code: "lixo", Goal: 3, Kind: Weekly, Status: "active"}
	ingles = Habit{Code: "ingles", Goal: 1, Kind: Weekdays, Days: []time.Weekday{time.Monday, time.Wednesday, time.Friday}, Status: "active"}
)

func TestWeekStart(t *testing.T) {
	tests := []struct {
		name string
		day  time.Time
		want string
	}{
		{"CA-19.3 quarta volta ao domingo", wed, "2026-11-08"},
		{"CA-19.3 domingo eh o proprio dia", date(2026, 11, 8), "2026-11-08"},
		{"CA-19.3 sabado ainda eh a mesma semana", date(2026, 11, 14), "2026-11-08"},
		{"virada de mes", date(2026, 12, 2), "2026-11-29"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := WeekStart(tc.day).Format(isoDate); got != tc.want {
				t.Errorf("WeekStart(%s) = %s, queria %s", tc.day.Format(isoDate), got, tc.want)
			}
		})
	}
}

func TestPeriod(t *testing.T) {
	from, to := lixo.Period(wed)
	if from != "2026-11-08" || to != "2026-11-14" {
		t.Errorf("CA-19.3 semana = %s a %s, queria 2026-11-08 a 2026-11-14", from, to)
	}
	from, to = cf.Period(wed)
	if from != "2026-11-11" || to != "2026-11-11" {
		t.Errorf("dia = %s a %s, queria so 2026-11-11", from, to)
	}
}

func TestMissing(t *testing.T) {
	tests := []struct {
		name string
		h    Habit
		done int
		day  time.Time
		want bool
	}{
		{"RN-23 diario nao feito", digi, 0, wed, true},
		{"RN-23 diario feito", digi, 1, wed, false},
		{"CA-18.1 quantidade pela metade", cf, 2, wed, true},
		{"CA-18.4 passou da meta", cf, 5, wed, false},
		{"CA-19.3 semanal ate bater", lixo, 2, wed, true},
		{"CA-19.3 semanal bateu", lixo, 3, wed, false},
		{"CA-19.4 dia dele", ingles, 0, wed, true},
		{"CA-19.4 nao eh dia dele", ingles, 0, date(2026, 11, 10), false},
		{"CA-21.3 arquivado nunca falta", Habit{Goal: 1, Kind: Daily, Status: "archived"}, 0, wed, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.h.Missing(tc.done, tc.day); got != tc.want {
				t.Errorf("Missing = %v, queria %v", got, tc.want)
			}
		})
	}
}

func TestLabels(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"CA-17.1 meta todo dia", digi.MetaLabel(), "todo dia"},
		{"CA-17.1 meta 3 todo dia", cf.MetaLabel(), "3 todo dia"},
		{"CA-17.1 meta por semana", lixo.MetaLabel(), "3 por semana"},
		{"CA-17.1 meta dias", ingles.MetaLabel(), "seg, qua, sex"},
		{"CA-19.2 feito 0/1", digi.DoneLabel(0, wed), "0/1"},
		{"CA-19.2 feito na semana", lixo.DoneLabel(2, wed), "2/3 na semana"},
		{"CA-19.2 bateu", cf.DoneLabel(3, wed), "✓"},
		{"CA-19.2 nao eh dia", ingles.DoneLabel(0, date(2026, 11, 10)), "-"},
		{"CA-18.1 confirma quantidade", cf.Confirm(2, false), "2/3 hoje"},
		{"CA-18.1 confirma semana", lixo.Confirm(2, false), "2/3 nesta semana"},
		{"CA-18.1 confirma feito", digi.Confirm(1, false), "feito hoje"},
		{"CA-18.3 confirma ontem", digi.Confirm(1, true), "feito ontem"},
		{"CA-18.4 passou da meta", cf.Confirm(5, false), "5/3 hoje"},
		{"CA-18.6 nome do dia", DayName(date(2026, 11, 10)), "terca"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("= %q, queria %q", tc.got, tc.want)
			}
		})
	}
}

func TestValidCode(t *testing.T) {
	tests := []struct {
		code    string
		wantErr bool
	}{
		{"cf", false}, {"cf3", false}, {"x", false}, {"xadrez", false},
		{"3", true}, {"x2", true}, {"x10", true},
	}

	for _, tc := range tests {
		if err := ValidCode(tc.code); (err != nil) != tc.wantErr {
			t.Errorf("CA-17.3 ValidCode(%q) = %v, queria erro: %v", tc.code, err, tc.wantErr)
		}
	}
}

func TestDays(t *testing.T) {
	days, err := ParseDays("seg,qua,sex")
	if err != nil || DaysValue(days) != "1,3,5" {
		t.Errorf("CA-17.1 dias = %v %v, queria 1,3,5", days, err)
	}
	if got := DaysFromValue("1,3,5"); len(got) != 3 || got[2] != time.Friday {
		t.Errorf("volta do banco = %v", got)
	}
	if _, err := ParseDays("seg,xyz"); err == nil || err.Error() != "dia xyz invalido (use dom, seg, ter, qua, qui, sex, sab)" {
		t.Errorf("CA-17.3 erro = %v", err)
	}
}
