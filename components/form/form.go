package form

import (
	. "github.com/cdvelop/tinystring"
)

// FormField defines a single field within a form.
type FormField struct {
	Type        string
	Name        string
	Placeholder string
	Required    bool
}

// FormConfig holds the configuration for a form.
type FormConfig struct {
	Action string
	Method string
	Fields []FormField
}

// Form implements HTMLRenderer, CSSRenderer, and JSRenderer interfaces.
type Form struct {
	Config FormConfig
}

// RenderHTML generates the HTML for the form.
func (f *Form) RenderHTML() string {
	var b = Convert()
	b.Write("<form class=\"contact-form\" action=\"")
	b.Write(Convert(f.Config.Action).EscapeAttr())
	b.Write("\" method=\"")
	b.Write(Convert(f.Config.Method).EscapeAttr())
	b.Write("\">\n")
	for _, field := range f.Config.Fields {
		if field.Type == "textarea" {
			b.Write("  <textarea name=\"")
			b.Write(Convert(field.Name).EscapeAttr())
			b.Write("\" placeholder=\"")
			b.Write(Convert(field.Placeholder).EscapeAttr())
			b.Write("\"")
			b.Write(boolAttr("required", field.Required))
			b.Write("></textarea>\n")
		} else {
			b.Write("<input type=\"")
			b.Write(Convert(field.Type).EscapeAttr())
			b.Write("\" name=\"")
			b.Write(Convert(field.Name).EscapeAttr())
			b.Write("\" placeholder=\"")
			b.Write(Convert(field.Placeholder).EscapeAttr())
			b.Write("\"")
			b.Write(boolAttr("required", field.Required))
			b.Write(">\n")
		}
	}
	b.Write("  <button type=\"submit\">Enviar Mensaje</button>\n")
	b.Write("</form>\n")
	return b.String()
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

// RenderJS returns the JavaScript for the form.
func (f *Form) RenderJS() string {
	return `// Simple form validation
document.addEventListener('DOMContentLoaded', function() {
    const forms = document.querySelectorAll('.contact-form');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            let isValid = true;
            const requiredFields = form.querySelectorAll('[required]');

            requiredFields.forEach(field => {
                if (!field.value.trim()) {
                    isValid = false;
                    // You might want to add a class to highlight the invalid field
                    field.style.borderColor = 'red';
                } else {
                    field.style.borderColor = ''; // Reset border color
                }
            });

            if (!isValid) {
                e.preventDefault(); // Stop form submission
                // You could also display a general error message to the user
                console.log('Please fill out all required fields.');
            }
        });
    });
});
`
}
