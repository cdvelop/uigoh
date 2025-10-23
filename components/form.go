
package components

import "strings"

type FormField struct {
    Type        string // text, email, textarea
    Name        string
    Placeholder string
    Required    bool
}

type FormConfig struct {
    Action string
    Method string
    Fields []FormField
}

type FormBuilder struct {
    Config FormConfig
}

func (f *FormBuilder) RenderHTML() string {
    var b strings.Builder

    b.WriteString("<form class=\"contact-form\" action=\"")
    b.WriteString(escapeAttr(f.Config.Action))
    b.WriteString("\" method=\"")
    b.WriteString(f.Config.Method)
    b.WriteString("\">\n")

    for _, field := range f.Config.Fields {
        if field.Type == "textarea" {
            b.WriteString("  <textarea name=\"")
            b.WriteString(escapeAttr(field.Name))
            b.WriteString("\" placeholder=\"")
            b.WriteString(escapeAttr(field.Placeholder))
            b.WriteString("\"")
            if field.Required {
                b.WriteString(" required")
            }
            b.WriteString("></textarea>\n")
        } else {
            b.WriteString("  <input type=\"")
            b.WriteString(field.Type)
            b.WriteString("\" name=\"")
            b.WriteString(escapeAttr(field.Name))
            b.WriteString("\" placeholder=\"")
            b.WriteString(escapeAttr(field.Placeholder))
            b.WriteString("\"")
            if field.Required {
                b.WriteString(" required")
            }
            b.WriteString(">\n")
        }
    }

    b.WriteString("  <button type=\"submit\">Enviar Mensaje</button>\n")
    b.WriteString("</form>\n")

    return b.String()
}

func (f *FormBuilder) RenderCSS() string {
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

func (f *FormBuilder) RenderJS() string {
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
