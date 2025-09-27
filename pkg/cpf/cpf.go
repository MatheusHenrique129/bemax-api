package cpf

import (
	"fmt"

	"github.com/klassmann/cpfcnpj"
)

func ValidateCPF(cpf string) error {
	if cpf == "" {
		return fmt.Errorf("CPF is required")
	}

	cpfClean := cpfcnpj.Clean(cpf)

	isValid := cpfcnpj.ValidateCPF(cpfClean)
	if !isValid {
		return fmt.Errorf("the CPF is invalid")
	}

	return nil
}

func FormatCPF(cpf string) string {
	newCPF := cpfcnpj.NewCPF(cpf)
	return newCPF.String()
}
