package task

import "testing"

func TestNormalizeTitle(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"CA-01.4 caixa diferente", "Comprar cabo HDMI", "comprar cabo hdmi"},
		{"CA-01.4 espaco a mais", "comprar  cabo   hdmi", "comprar cabo hdmi"},
		{"espaco nas pontas", "  comprar cabo ", "comprar cabo"},
		{"ja normalizado", "ligar", "ligar"},
		{"dois pontos fica", "reuniao 10:30", "reuniao 10:30"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeTitle(tc.in)
			if got != tc.want {
				t.Errorf("NormalizeTitle(%q) = %q, queria %q", tc.in, got, tc.want)
			}
		})
	}
}
