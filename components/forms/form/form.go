package form

import (
	. "github.com/cdvelop/tinystring"
)

// Field defines a single field within a form.
type Field struct {
	Type        string
	Name        string
	Placeholder string
	Required    bool
}

// Config holds the configuration for a form.
type Config struct {
	Action string
	Method string
	Fields []Field
}

// Form implements HTMLRenderer, CSSRenderer, and JSRenderer interfaces.
type Form struct {
	Config Config
}

// RenderHTML generates the HTML for the form.
func (f *Form) RenderHTML() string {
	// Build fields HTML
	fields := ""
	for _, field := range f.Config.Fields {
		name := Convert(field.Name).EscapeAttr()
		placeholder := Convert(field.Placeholder).EscapeAttr()
		if field.Type == "textarea" {
			req := boolAttr("required", field.Required)
			fields += Fmt("  <textarea name=\"%s\" placeholder=\"%s\"%s></textarea>\n", name, placeholder, req)
		} else {
			t := Convert(field.Type).EscapeAttr()
			req := boolAttr("required", field.Required)
			fields += Fmt("  <input type=\"%s\" name=\"%s\" placeholder=\"%s\"%s>\n", t, name, placeholder, req)
		}
	}

	action := Convert(f.Config.Action).EscapeAttr()
	method := Convert(f.Config.Method).EscapeAttr()

	tpl := `<form class="contact-form" action="%s" method="%s">
%s  <button type="submit">Enviar Mensaje</button>
</form>
`

	return Fmt(tpl, action, method, fields)
}

func boolAttr(attr string, val bool) string {
	if val {
		return " " + attr
	}
	return ""
}

// RenderCSS returns the CSS for the form.
func (f *Form) RenderCSS() string {
	return `.contact-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 500px;
}

.contact-form input,
.contact-form textarea {
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-family: inherit;
}

.contact-form button {
  padding: 0.75rem 1.5rem;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.contact-form button:hover {
  opacity: 0.9;
}
`
}
