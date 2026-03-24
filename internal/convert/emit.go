package convert

import (
	"fmt"
	"strings"

	"sql-migrator/internal/dialect"
	"sql-migrator/internal/model"
)

func emitScript(script *model.Script, from, to dialect.Dialect) (string, error) {
	var b strings.Builder
	for i, st := range script.Statements {
		var s string
		var err error
		switch st := st.(type) {
		case *model.Raw:
			s, err = legacyApply(st.SQL, from, to)
		case *model.CreateTable:
			s, err = emitCreateTable(st, from, to)
		case *model.Insert:
			s, err = emitInsert(st, from, to)
		default:
			return "", fmt.Errorf("unknown statement type %T", st)
		}
		if err != nil {
			return "", err
		}
		b.WriteString(s)
		if i < len(script.Statements)-1 {
			b.WriteString(";\n")
		} else {
			b.WriteString(";")
		}
	}
	return b.String(), nil
}

func emitCreateTable(ct *model.CreateTable, from, to dialect.Dialect) (string, error) {
	var b strings.Builder
	b.WriteString("CREATE TABLE ")
	if ct.IfNotExists {
		b.WriteString("IF NOT EXISTS ")
	}
	b.WriteString(formatIdent(ct.Table, to))
	b.WriteString(" (\n")
	for i, col := range ct.Columns {
		if i > 0 {
			b.WriteString(",\n")
		}
		b.WriteString("  ")
		b.WriteString(formatIdent(col.Name, to))
		b.WriteString(" ")
		rest, err := legacyApply(col.Rest, from, to)
		if err != nil {
			return "", err
		}
		b.WriteString(rest)
	}
	for _, line := range ct.TableLevel {
		b.WriteString(",\n  ")
		lineOut, err := legacyApply(line, from, to)
		if err != nil {
			return "", err
		}
		b.WriteString(lineOut)
	}
	b.WriteString("\n)")
	s := b.String()
	if to == dialect.MariaDB && (from == dialect.Postgres || from == dialect.SQLite) {
		s = doubleQuotedIdentToBackticks(s)
	}
	return s, nil
}

func emitInsert(ins *model.Insert, from, to dialect.Dialect) (string, error) {
	s := buildInsertSQL(ins, from)
	return legacyApply(s, from, to)
}

func buildInsertSQL(ins *model.Insert, from dialect.Dialect) string {
	if ins.OrReplace && from == dialect.MariaDB {
		var b strings.Builder
		b.WriteString("REPLACE INTO ")
		b.WriteString(formatIdent(ins.Table, from))
		b.WriteString(" (")
		for i, c := range ins.Columns {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatIdent(c, from))
		}
		b.WriteString(") VALUES ")
		b.WriteString(ins.Values)
		return b.String()
	}
	var b strings.Builder
	b.WriteString("INSERT ")
	if ins.OrIgnore {
		switch from {
		case dialect.SQLite, dialect.SQLServer:
			b.WriteString("OR IGNORE ")
		case dialect.MariaDB:
			b.WriteString("IGNORE ")
		}
	} else if ins.OrReplace {
		switch from {
		case dialect.SQLite:
			b.WriteString("OR REPLACE ")
		}
	}
	b.WriteString("INTO ")
	b.WriteString(formatIdent(ins.Table, from))
	b.WriteString(" (")
	for i, c := range ins.Columns {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(formatIdent(c, from))
	}
	b.WriteString(") VALUES ")
	b.WriteString(ins.Values)
	return b.String()
}

func formatIdent(id model.Ident, d dialect.Dialect) string {
	if !id.Quoted {
		return id.Name
	}
	switch d {
	case dialect.Postgres, dialect.SQLite:
		return `"` + escapeDouble(id.Name) + `"`
	case dialect.MariaDB:
		return "`" + escapeBacktick(id.Name) + "`"
	case dialect.SQLServer:
		return "[" + strings.ReplaceAll(id.Name, "]", "]]") + "]"
	default:
		return id.Name
	}
}

func escapeDouble(s string) string { return strings.ReplaceAll(s, `"`, `""`) }

func escapeBacktick(s string) string { return strings.ReplaceAll(s, "`", "``") }
