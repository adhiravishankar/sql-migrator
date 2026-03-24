package sqlparse

import (
	"regexp"
	"strings"

	"sql-migrator/internal/model"
)

var reCreateTablePrefix = regexp.MustCompile(`(?is)^\s*CREATE\s+TABLE\s+`)

func isTableConstraint(line string) bool {
	u := strings.ToUpper(strings.TrimSpace(line))
	switch {
	case strings.HasPrefix(u, "PRIMARY KEY"),
		strings.HasPrefix(u, "FOREIGN KEY"),
		strings.HasPrefix(u, "CHECK "),
		strings.HasPrefix(u, "CONSTRAINT "):
		return true
	case strings.HasPrefix(u, "UNIQUE KEY"), strings.HasPrefix(u, "UNIQUE ("), strings.HasPrefix(u, "UNIQUE("):
		return true
	default:
		return false
	}
}

// TryParseCreateTable parses a CREATE TABLE statement into a model, or returns false.
func TryParseCreateTable(s string) (*model.CreateTable, bool) {
	loc := reCreateTablePrefix.FindStringIndex(s)
	if loc == nil {
		return nil, false
	}
	rest := strings.TrimSpace(s[loc[1]:])
	ct := &model.CreateTable{}
	const ifNotExists = "IF NOT EXISTS "
	if len(rest) >= len(ifNotExists) && strings.EqualFold(rest[:len(ifNotExists)], ifNotExists) {
		ct.IfNotExists = true
		rest = strings.TrimSpace(rest[len(ifNotExists):])
	}
	tid, rest2, ok := ParseLeadingIdent(rest)
	if !ok {
		return nil, false
	}
	ct.Table = tid
	rest2 = strings.TrimSpace(rest2)
	inner, after, ok := extractParenBody(rest2)
	if !ok {
		return nil, false
	}
	_ = after // trailing content after ); ignored for now
	parts := splitCommaList(inner)
	for _, p := range parts {
		line := strings.TrimSpace(p)
		if line == "" {
			continue
		}
		if isTableConstraint(line) {
			ct.TableLevel = append(ct.TableLevel, line)
			continue
		}
		colName, colRest, ok := ParseLeadingIdent(line)
		if !ok {
			return nil, false
		}
		ct.Columns = append(ct.Columns, model.ColumnDef{
			Name: colName,
			Rest: strings.TrimSpace(colRest),
		})
	}
	return ct, true
}
