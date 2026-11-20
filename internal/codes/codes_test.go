package codes

import "testing"

func TestCheck(t *testing.T) {
	reserved := []string{"add", "ls"}
	tests := []struct {
		code    string
		wantErr bool
	}{
		{"BD", false}, {"ARQ1", false}, {"p1", false}, {"bd-feedback360", false}, {"lss", false},
		{"2", true}, {"15", true}, {"2BD", true}, {"1arq", true}, {"ls", true}, {"add", true},
	}
	for _, tc := range tests {
		if err := Check(tc.code, "teste", reserved); (err != nil) != tc.wantErr {
			t.Errorf("RN-26 Check(%q) = %v, queria erro: %v", tc.code, err, tc.wantErr)
		}
	}
}
