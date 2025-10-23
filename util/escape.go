
package util

import "html"

// EscapeHTML escapes HTML special characters.
func EscapeHTML(s string) string {
	return html.EscapeString(s)
}

// EscapeAttr escapes HTML attribute special characters.
func EscapeAttr(s string) string {
	return html.EscapeString(s)
}
