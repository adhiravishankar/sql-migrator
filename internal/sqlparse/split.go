package sqlparse

import "strings"

// SplitStatements splits SQL on semicolons outside of string literals and
// quoted identifiers ('...', "...", `...`).
func SplitStatements(sql string) []string {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil
	}
	var out []string
	var b strings.Builder
	inSingle := false
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		switch {
		case c == '\'' && !inSingle:
			inSingle = true
			b.WriteByte(c)
		case c == '\'' && inSingle:
			if i+1 < len(sql) && sql[i+1] == '\'' {
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
			for i < len(sql) {
				if sql[i] == '"' {
					if i+1 < len(sql) && sql[i+1] == '"' {
						b.WriteString(`""`)
						i += 2
						continue
					}
					b.WriteByte('"')
					break
				}
				b.WriteByte(sql[i])
				i++
			}
		case c == '`':
			b.WriteByte(c)
			i++
			for i < len(sql) {
				if sql[i] == '`' {
					if i+1 < len(sql) && sql[i+1] == '`' {
						b.WriteString("``")
						i += 2
						continue
					}
					b.WriteByte('`')
					break
				}
				b.WriteByte(sql[i])
				i++
			}
		case c == ';':
			s := strings.TrimSpace(b.String())
			if s != "" {
				out = append(out, s)
			}
			b.Reset()
		default:
			b.WriteByte(c)
		}
	}
	if tail := strings.TrimSpace(b.String()); tail != "" {
		out = append(out, tail)
	}
	return out
}
