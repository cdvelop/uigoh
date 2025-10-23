
package components

import (
	"crypto/sha256"
	"encoding/hex"
	"html"
)

func escapeAttr(s string) string {
	return html.EscapeString(s)
}

func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func boolAttr(attr string, val bool) string {
	if val {
		return " " + attr
	}
	return ""
}
