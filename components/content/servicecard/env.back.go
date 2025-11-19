//go:build !wasm
// +build !wasm

package servicecard

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the servicecard.
func (c *ServiceCard) RenderCSS() string {
	return styleCss
}
