package sqlparse

import "strings"

// splitCommaList splits s on commas that are not inside parentheses or string/identifier quotes.
func splitCommaList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var parts []string
	var b strings.Builder
	depth := 0
	inSingle := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inSingle:
			inSingle = true
			b.WriteByte(c)
		case c == '\'' && inSingle:
			if i+1 < len(s) && s[i+1] == '\'' {
				b.WriteString("''")
				i++
				continue
			}
			inSingle = false
			b.WriteByte(c)
		case inSingle:
			b.WriteByte(c)
		case c == '"':
			b.WriteByte(c)
			i++
			for i < len(s) {
				b.WriteByte(s[i])
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						b.WriteByte(s[i+1])
						i += 2
						continue
					}
					break
				}
				i++
			}
		case c == '`':
			b.WriteByte(c)
			i++
			for i < len(s) {
				b.WriteByte(s[i])
				if s[i] == '`' {
					if i+1 < len(s) && s[i+1] == '`' {
						b.WriteByte(s[i+1])
						i += 2
						continue
					}
					break
				}
				i++
			}
		case c == '(':
			depth++
			b.WriteByte(c)
		case c == ')':
			depth--
			b.WriteByte(c)
		case c == ',' && depth == 0:
			parts = append(parts, strings.TrimSpace(b.String()))
			b.Reset()
		default:
			b.WriteByte(c)
		}
	}
	if b.Len() > 0 {
		parts = append(parts, strings.TrimSpace(b.String()))
	}
	return parts
}
