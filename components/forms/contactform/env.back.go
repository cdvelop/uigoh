//go:build !wasm
// +build !wasm

package contactform

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the contact form.
func (c *ContactForm) RenderCSS() string {
	return styleCss
}
