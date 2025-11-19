//go:build !wasm
// +build !wasm

package packagecard

import (
	_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the packagecard.
func (c *PackageCard) RenderCSS() string {
	return styleCss
}
