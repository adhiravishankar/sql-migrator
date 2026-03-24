package sqlparse

import "strings"

// extractParenBody returns the text inside the first balanced (...) where s begins with '('.
func extractParenBody(s string) (inner string, rest string, ok bool) {
	s = strings.TrimSpace(s)
	if len(s) == 0 || s[0] != '(' {
		return "", "", false
	}
	depth := 0
	inSingle := false
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inSingle {
			if c == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
			}
			continue
		}
		switch c {
		case '\'':
			inSingle = true
		case '"':
			i++
			for i < len(s) {
				if s[i] == '"' {
					if i+1 < len(s) && s[i+1] == '"' {
						i += 2
						continue
					}
					break
				}
				i++
			}
		case '`':
			i++
			for i < len(s) {
				if s[i] == '`' {
					if i+1 < len(s) && s[i+1] == '`' {
						i += 2
						continue
					}
					break
				}
				i++
			}
		case '(':
			depth++
			if depth == 1 {
				start = i + 1
			}
		case ')':
			depth--
			if depth == 0 {
				return s[start:i], s[i+1:], true
			}
		}
	}
	return "", "", false
}
