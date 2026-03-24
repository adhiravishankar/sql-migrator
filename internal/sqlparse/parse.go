package sqlparse

import "sql-migrator/internal/model"

// Parse splits SQL into statements and parses supported forms into structured models.
// Unrecognized statements are kept as model.Raw.
func Parse(sql string) *model.Script {
	stmts := SplitStatements(sql)
	var out []model.Statement
	for _, s := range stmts {
		if ct, ok := TryParseCreateTable(s); ok {
			out = append(out, ct)
			continue
		}
		if ins, ok := TryParseInsert(s); ok {
			out = append(out, ins)
			continue
		}
		out = append(out, &model.Raw{SQL: s})
	}
	return &model.Script{Statements: out}
}
