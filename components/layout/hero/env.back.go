//go:build !wasm
// +build !wasm

package hero

import (
_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the hero section.
func (h *Hero) RenderCSS() string {
return styleCss
}
