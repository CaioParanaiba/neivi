// Package codes confere o formato dos codigos de disciplina, projeto e habito (RN-26)
package codes

import (
	"fmt"
	"strconv"
	"strings"
)

// Check recusa codigo que comeca com numero (o todo numero se confundiria com o id)
// ou que eh nome de subcomando do grupo
func Check(code, group string, reserved []string) error {
	if _, err := strconv.Atoi(code); err == nil {
		return fmt.Errorf("o codigo %s nao pode ser so um numero (ele se confundiria com o id)", code)
	}
	if code != "" && code[0] >= '0' && code[0] <= '9' {
		return fmt.Errorf("o codigo %s nao pode comecar com numero", code)
	}
	for _, r := range reserved {
		if code == r {
			return fmt.Errorf("o codigo %s eh um comando do %s (%s); escolha outro", code, group, strings.Join(reserved, ", "))
		}
	}
	return nil
}
