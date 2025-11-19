package formPlatform

import (
	. "github.com/cdvelop/tinystring"
)

// Method to dynamically generate the title
func (f *field) setDynamicTitle() {
	if f.Title != "" {
		return
	}

	var parts []string
	parts = append(parts, Translate(D.Allowed).String())

	if f.Letters {
		parts = append(parts, Translate(D.Letters).String())
	}

	if f.Numbers {
		parts = append(parts, Translate(D.Numbers).String())
	}

	if len(f.Characters) > 0 {
		var chars []string
		for _, char := range f.Characters {
			if char == ' ' {
				chars = append(chars, "␣")
			} else {
				chars = append(chars, string(char))
			}
		}
		parts = append(parts, Translate(D.Chars, chars).String())
	}

	if f.Minimum != 0 {
		parts = append(parts, Translate("min", f.Minimum).String())
	}

	if f.Maximum != 0 {
		parts = append(parts, Translate("max", f.Maximum).String())
	}

	f.Title = Convert(parts).Join().String()

}
