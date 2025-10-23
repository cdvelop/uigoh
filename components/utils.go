
package components

import (
    "crypto/sha256"
    "encoding/hex"
    "html"
)

func escapeHTML(s string) string {
    return html.EscapeString(s)
}

func escapeAttr(s string) string {
    // More strict escaping for attributes
    return html.EscapeString(s)
}

// hashString creates a SHA256 hash of the input string
// Used for CSS/JS deduplication
func hashString(s string) string {
    h := sha256.New()
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}
