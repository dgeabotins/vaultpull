package mask

import "strings"

const defaultMaskChar = "*"
const visibleChars = 4

// Masker redacts secret values for safe display in logs and output.
type Masker struct {
	revealChars int
	maskChar    string
}

// New returns a Masker with sensible defaults.
func New() *Masker {
	return &Masker{revealChars: visibleChars, maskChar: defaultMaskChar}
}

// Mask redacts all but the last n characters of value.
// Values shorter than or equal to revealChars are fully masked.
func (m *Masker) Mask(value string) string {
	if len(value) <= m.revealChars {
		return strings.Repeat(m.maskChar, 8)
	}
	visible := value[len(value)-m.revealChars:]
	return strings.Repeat(m.maskChar, 8) + visible
}

// MaskMap returns a copy of the map with all values masked.
func (m *Masker) MaskMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.Mask(v)
	}
	return out
}

// IsSafe reports whether the value appears already masked.
func IsSafe(value string) bool {
	return strings.HasPrefix(value, defaultMaskChar)
}
