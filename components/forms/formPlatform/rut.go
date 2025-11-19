package formPlatform

import (
	. "github.com/cdvelop/tinystring"
)

func runData(runIn string) (data []string, onlyRun int, err error) {

	if len(runIn) < 3 {
		return nil, 0, Err(D.Value, D.Empty)
	}

	// Separar número y dígito verificador
	data = Convert(runIn).Split("-")
	if len(data) != 2 {
		return nil, 0, Err(D.Format, D.Not, D.Valid)
	}

	// Validar caracteres del número
	if !isDigits(data[0]) {
		return nil, 0, Err(D.Chars, D.Not, D.Allowed, D.In, D.Numbers)
	}

	// Validar dígito verificador
	dv := Convert(data[1]).ToLower().String()
	if len(dv) != 1 || (dv != "k" && !isDigits(dv)) {
		return nil, 0, Err(D.Digit, D.Checker, dv, D.Not, D.Valid)
	}

	// Convertir número a entero
	onlyRun, err = Convert(data[0]).Int()
	if err != nil {
		return nil, 0, Err(D.Numbers, D.Not, D.Valid)
	}

	return data, onlyRun, nil
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
