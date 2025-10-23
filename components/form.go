
package components

import (
	"fmt"
	"strings"
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

// Form implements the Component interface for a form.
type Form struct {
	Config FormConfig
}

// RenderHTML generates the HTML for the form.
func (f *Form) RenderHTML() string {
	var b strings.Builder
	fmt.Fprintf(&b, "<form class=\"contact-form\" action=\"%s\" method=\"%s\">\n", escapeAttr(f.Config.Action), escapeAttr(f.Config.Method))
	for _, field := range f.Config.Fields {
		if field.Type == "textarea" {
			fmt.Fprintf(&b, "  <textarea name=\"%s\" placeholder=\"%s\"%s></textarea>\n", escapeAttr(field.Name), escapeAttr(field.Placeholder), boolAttr("required", field.Required))
		} else {
			fmt.Fprintf(&b, "  <input type=\"%s\" name=\"%s\" placeholder=\"%s\"%s>\n", escapeAttr(field.Type), escapeAttr(field.Name), escapeAttr(field.Placeholder), boolAttr("required", field.Required))
		}
	}
	b.WriteString("  <button type=\"submit\">Enviar Mensaje</button>\n")
	b.WriteString("</form>\n")
	return b.String()
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
