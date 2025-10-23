
package tinystrings

import "strings"

// Builder is a wrapper around strings.Builder.
type Builder struct {
	b strings.Builder
}

// WriteString appends the contents of s to b's buffer.
func (b *Builder) WriteString(s string) (int, error) {
	return b.b.WriteString(s)
}

// String returns the accumulated string.
func (b *Builder) String() string {
	return b.b.String()
}

// ToLower returns s with all Unicode letters mapped to their lower case.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ReplaceAll returns a copy of the string s with all
// non-overlapping instances of old replaced by new.
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}
