//go:build !wasm
// +build !wasm

package footer

import (
_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the footer.
func (c *Footer) RenderCSS() string {
return styleCss
}
