package convert

import (
	"sql-migrator/internal/dialect"
	"sql-migrator/internal/sqlparse"
)

// Convert transforms SQL from one dialect to another. It parses the input into
// an intermediate model (CREATE TABLE and INSERT where possible), then emits
// target SQL; other statements use the legacy regex pipeline.
func Convert(input string, from, to dialect.Dialect) (string, error) {
	if from == to {
		return input, nil
	}
	script := sqlparse.Parse(input)
	return emitScript(script, from, to)
}
