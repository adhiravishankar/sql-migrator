package convert

import "strings"

// doubleQuotedIdentToBackticks turns PostgreSQL/SQL standard double-quoted identifiers
// into MariaDB/MySQL backtick-quoted identifiers. Single-quoted string literals are
// copied verbatim so " characters inside strings are not treated as identifiers.
// Inside double quotes, PostgreSQL escapes " as "".
func doubleQuotedIdentToBackticks(sql string) string {
	var b strings.Builder
	b.Grow(len(sql) + 8)
	i := 0
	for i < len(sql) {
		switch sql[i] {
		case '\'':
			b.WriteByte(sql[i])
			i++
			for i < len(sql) {
				if sql[i] == '\'' {
					if i+1 < len(sql) && sql[i+1] == '\'' {
						b.WriteString("''")
						i += 2
						continue
					}
					b.WriteByte('\'')
					i++
					break
				}
				b.WriteByte(sql[i])
				i++
			}
		case '"':
			start := i
			i++
			var inner strings.Builder
			closed := false
			for i < len(sql) {
				if sql[i] == '"' {
					if i+1 < len(sql) && sql[i+1] == '"' {
						inner.WriteByte('"')
						i += 2
						continue
					}
					closed = true
					i++
					break
				}
				inner.WriteByte(sql[i])
				i++
			}
			if closed {
				b.WriteByte('`')
				b.WriteString(strings.ReplaceAll(inner.String(), "`", "``"))
				b.WriteByte('`')
			} else {
				b.WriteString(sql[start:])
				return b.String()
			}
		default:
			b.WriteByte(sql[i])
			i++
		}
	}
	return b.String()
}
