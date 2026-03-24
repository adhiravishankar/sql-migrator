package sqlparse

import (
	"strings"

	"sql-migrator/internal/model"
)

// TryParseInsert parses INSERT ... VALUES ... into a model, or returns false.
func TryParseInsert(s string) (*model.Insert, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 6 || !strings.EqualFold(s[:6], "INSERT") {
		return nil, false
	}
	rest := strings.TrimSpace(s[6:])
	ins := &model.Insert{}
	switch {
	case len(rest) >= 10 && strings.EqualFold(rest[:10], "OR IGNORE "):
		ins.OrIgnore = true
		rest = strings.TrimSpace(rest[10:])
	case len(rest) >= 11 && strings.EqualFold(rest[:11], "OR REPLACE "):
		ins.OrReplace = true
		rest = strings.TrimSpace(rest[11:])
	case len(rest) >= 7 && strings.EqualFold(rest[:7], "IGNORE "):
		ins.OrIgnore = true
		rest = strings.TrimSpace(rest[7:])
	}
	const into = "INTO "
	if len(rest) < len(into) || !strings.EqualFold(rest[:len(into)], into) {
		return nil, false
	}
	rest = strings.TrimSpace(rest[len(into):])
	tid, rest2, ok := ParseLeadingIdent(rest)
	if !ok {
		return nil, false
	}
	ins.Table = tid
	rest2 = strings.TrimSpace(rest2)
	innerCols, rest3, ok := extractParenBody(rest2)
	if !ok {
		return nil, false
	}
	for _, p := range splitCommaList(innerCols) {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		cid, _, ok := ParseLeadingIdent(p)
		if !ok {
			return nil, false
		}
		ins.Columns = append(ins.Columns, cid)
	}
	rest3 = strings.TrimSpace(rest3)
	const values = "VALUES"
	if len(rest3) < len(values) || !strings.EqualFold(rest3[:len(values)], values) {
		return nil, false
	}
	ins.Values = strings.TrimSpace(rest3[len(values):])
	if ins.Values == "" {
		return nil, false
	}
	return ins, true
}
