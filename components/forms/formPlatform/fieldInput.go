package formPlatform

import (
	"reflect"

	. "github.com/cdvelop/tinystring"
)

func NewField(name string, params ...any) (*field, error) {
	f := &field{
		Name: Convert(name).SnakeLow().String(),
	}

	// Extract typed parameters first
	var remainingParams []any
	for _, param := range params {
		switch v := param.(type) {
		case uint32:
			f.Index = v
		case bool:
			f.Unique = v
		case *entity:
			f.ForeignKey = v
		case *reflect.StructField:
			f.reflectStructField = v
		default:
			remainingParams = append(remainingParams, param)
		}
	}
	params = remainingParams

	switch f.Name {
	case "date", "birth_date": // formato fecha: DD-MM-YYYY
		f.htmlName = "date"
		f.Title = `title="` + Translate(D.Format, D.Date, `: DD-MM-YYYY"`).String()
		f.ExtraValidation = G.Date.DateExists

	case "date_age":
		f.htmlName = "date"
		f.Title = `title="` + Translate(D.Format, D.Date, `: DD-MM-YYYY"`).String()

	case "day_word":
		f.htmlName = "date"
		f.DataSet = []map[string]string{{"spanish": ""}}

	case "email":
		f.htmlName = "mail"
		f.PlaceHolder = "ej: mi.correo@mail.com"
		f.permitted = permitted{Letters: true, Numbers: true, Characters: []rune{'@', '.', '_'}, Minimum: 0, Maximum: 0, ExtraValidation: func(s string) error {
			if Contains(s, "example") {
				return Err(D.Example, D.Email, D.Not, D.Allowed)
			}

			parts := Convert(s).Split("@")
			if len(parts) != 2 {
				return Err(D.Format, D.Email, D.Not, D.Valid)
			}

			return nil
		}}

	case "file_path":
		f.htmlName = "file"
		f.permitted = permitted{
			Letters: true, Tilde: false, Numbers: true, Characters: []rune{'\\', '/', '.', '_'}, Minimum: 1, Maximum: 100,
			StartWith: &permitted{Letters: true, Numbers: true, Characters: []rune{'.', '_', '/'}},
		}

	case "gender", "radio":
		f.htmlName = "radio"
		if f.Name == "gender" {
			f.options = []map[string]string{{"f": Translate(D.Female).String()}, {"m": Translate(D.Male).String()}}
		}

	case "hour":
		f.htmlName = "time"
		f.Title = Translate(D.Format, D.Hour, ": HH:MM").String()
		f.permitted = permitted{Numbers: true, Characters: []rune{':'}, Minimum: 5, Maximum: 5, TextNotAllowed: []string{"24:"}}

	case "id":
		f.htmlName = "hidden"
		f.allowSkipCompleted = true
		f.permitted = permitted{Letters: false, Numbers: true, Characters: []rune{'.'}, Minimum: 1, Maximum: 39}

	case "info":
		f.htmlName = "text"

	case "ip":
		f.htmlName = "text"
		f.Title = Translate(D.Example, ": 192.168.0.8").String()
		f.permitted = permitted{Letters: true, Numbers: true, Characters: []rune{'.', ':'}, Minimum: 7, Maximum: 39, ExtraValidation: func(value string) error {
			var ipV string

			if value == "0.0.0.0" {
				return Err(D.Example, "IP", D.Not, D.Allowed, ':', "0.0.0.0")
			}

			if Contains(value, ":") { //IPv6
				ipV = ":"
			} else if Contains(value, ".") { //IPv4
				ipV = "."
			}

			part := Convert(value).Split(ipV)

			if ipV == "." && len(part) != 4 {
				return Err(D.Format, "IPv4", D.Not, D.Valid)
			}

			if ipV == ":" && len(part) != 8 {
				return Err(D.Format, "IPv6", D.Not, D.Valid)
			}
			return nil
		}}

	case "month_day":
		f.htmlName = "text"
		f.permitted = permitted{Numbers: true, Minimum: 2, Maximum: 2, ExtraValidation: func(value string) error {
			_, err := G.Date.ValidateDay(value)
			return err
		}}

	case "name", "text":
		f.htmlName = "text"
		f.permitted = permitted{Letters: true, Tilde: false, Numbers: true, Characters: []rune{' ', '.', ',', '(', ')'}, Minimum: 2, Maximum: 100}

	case "number", "phone":
		f.htmlName = "number"
		f.permitted = permitted{Numbers: true, Minimum: 1, Maximum: 20}

		if f.Name == "phone" {
			f.Min = "7"
			f.Max = "11"
		}

	case "password":
		f.permitted = permitted{
			Letters: true, Tilde: true, Numbers: true,
			Characters: []rune{' ', '#', '%', '?', '.', ',', '-', '_'}, Minimum: 5, Maximum: 50}

	case "rut":
		f.htmlName = "text"
		f.Autocomplete = `autocomplete="off"`
		f.Title = "rut sin puntos y con guion ej.: 11222333-4"
		f.Class = []className{"rut"}

		f.permitted = permitted{Numbers: true, Letters: true, Minimum: 9, Maximum: 11, Characters: []rune{'-'},
			StartWith: &permitted{Numbers: true},
			ExtraValidation: func(value string) error {
				// Validar RUT chileno
				if !Contains(value, "-") {
					return Err(D.Hyphen, D.Not, D.Found)
				}
				data, onlyRun, err := runData(value)
				if err != nil {
					return err
				}

				if data[0][0:1] == "0" {
					return Err(D.Not, D.Begin, D.With, D.Digit, 0)
				}

				dv := G.Rut.DvRut(onlyRun)

				expectedDV := Convert(data[1]).ToLower().String()
				if dv != expectedDV {
					return Err(D.Digit, D.Checker, expectedDV, D.Not, D.Valid)
				}
				return nil
			},
		}

		if !f.hideTyping {
			f.PlaceHolder = "ej: 11222333-4"
		}

	case "select":
		// if reflectType != nil {
		// params = append(params, "structure="+reflectType.String())
		// }

	case "text_area":
		f.htmlName = "textarea"
		f.Rows = "3"
		f.Cols = "1"
		f.permitted = permitted{Letters: true, Tilde: true, Numbers: true, BreakLine: true, WhiteSpaces: true, Tabulation: true, Characters: []rune{'%', '+', '#', '-', '.', ',', ':', '(', ')'}, Minimum: 2, Maximum: 1000}

	case "text_number":
		f.htmlName = "text"
		f.permitted = permitted{Letters: true, Numbers: true, Characters: []rune{'_'}, Minimum: 5, Maximum: 20}

	case "text_number_code":
		f.htmlName = "tel"
		f.permitted = permitted{Letters: true, Numbers: true,
			Characters: []rune{'_', '-'}, Minimum: 2, Maximum: 15,
			StartWith: &permitted{Letters: true, Numbers: true}}

	case "text_only":
		f.htmlName = "text"
		f.permitted = permitted{Letters: true, Minimum: 3, Maximum: 50, Characters: []rune{' '}}

	case "text_search":
		f.htmlName = "search"
		f.permitted = permitted{Letters: true, Tilde: false, Numbers: true, Characters: []rune{'-', ' '}, Minimum: 2, Maximum: 20}
	default:
		return nil, Err(D.Field, ':', name, D.Not, D.Found, D.In, D.Dictionary)
	}

	f.SetPropertiesFromInputTag(params...)

	return f, nil
}
