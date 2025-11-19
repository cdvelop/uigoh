//go:build !wasm
// +build !wasm

package doctorcard

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the doctorcard.
func (c *DoctorCard) RenderCSS() string {
	return styleCss
}
