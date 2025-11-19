//go:build !wasm
// +build !wasm

package banner

import (
_ "embed"
)

//go:embed style.css
var styleCss string

// RenderCSS returns the CSS for the banner.
func (c *Banner) RenderCSS() string {
return styleCss
}
