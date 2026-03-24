package sqlparse

import (
	"strings"

	"sql-migrator/internal/model"
)

// ParseLeadingIdent reads one SQL identifier starting at s (after trim). Returns
// the ident and the remainder of the string.
func ParseLeadingIdent(s string) (model.Ident, string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return model.Ident{}, "", false
	}
	switch s[0] {
	case '"':
		var inner strings.Builder
		i := 1
		for i < len(s) {
			if s[i] == '"' {
				if i+1 < len(s) && s[i+1] == '"' {
					inner.WriteByte('"')
					i += 2
					continue
				}
				i++
				return model.Ident{Name: inner.String(), Quoted: true}, strings.TrimSpace(s[i:]), true
			}
			inner.WriteByte(s[i])
			i++
		}
		return model.Ident{}, "", false
	case '`':
		var inner strings.Builder
		i := 1
		for i < len(s) {
			if s[i] == '`' {
				if i+1 < len(s) && s[i+1] == '`' {
					inner.WriteByte('`')
					i += 2
					continue
				}
				i++
				return model.Ident{Name: inner.String(), Quoted: true}, strings.TrimSpace(s[i:]), true
			}
			inner.WriteByte(s[i])
			i++
		}
		return model.Ident{}, "", false
	case '[':
		var inner strings.Builder
		i := 1
		for i < len(s) && s[i] != ']' {
			if s[i] == ']' && i+1 < len(s) && s[i+1] == ']' {
				inner.WriteByte(']')
				i += 2
				continue
			}
			inner.WriteByte(s[i])
			i++
		}
		if i >= len(s) || s[i] != ']' {
			return model.Ident{}, "", false
		}
		return model.Ident{Name: inner.String(), Quoted: true}, strings.TrimSpace(s[i+1:]), true
	default:
		j := 0
		for j < len(s) {
			ch := s[j]
			if ch == '.' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '$' {
				j++
				continue
			}
			break
		}
		if j == 0 {
			return model.Ident{}, "", false
		}
		name := s[:j]
		return model.Ident{Name: name, Quoted: false}, strings.TrimSpace(s[j:]), true
	}
}
