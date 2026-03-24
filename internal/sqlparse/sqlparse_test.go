package sqlparse

import (
	"strings"
	"testing"

	"sql-migrator/internal/model"
)

func TestSplitStatements(t *testing.T) {
	sql := `SELECT 'a;b'; SELECT 1;`
	parts := SplitStatements(sql)
	if len(parts) != 2 || !strings.Contains(parts[0], "'a;b'") {
		t.Fatalf("got %v", parts)
	}
}

func TestParse_CreateTableAndInsert(t *testing.T) {
	sql := `CREATE TABLE t (id INT NOT NULL);
INSERT INTO t (id) VALUES (1);`
	script := Parse(sql)
	if len(script.Statements) != 2 {
		t.Fatalf("len=%d", len(script.Statements))
	}
	if _, ok := script.Statements[0].(*model.CreateTable); !ok {
		t.Fatalf("want CreateTable, got %T", script.Statements[0])
	}
	if _, ok := script.Statements[1].(*model.Insert); !ok {
		t.Fatalf("want Insert, got %T", script.Statements[1])
	}
}
