package convert

import (
	"strings"
	"testing"

	"sql-migrator/internal/dialect"
)

func TestDoubleQuotedIdentToBackticks_TableAndColumns(t *testing.T) {
	in := `CREATE TABLE "vehicle_series" (
  "region" INTEGER NOT NULL,
  PRIMARY KEY ("region", "transit_hub")
);`
	out := doubleQuotedIdentToBackticks(in)
	if !strings.Contains(out, "`vehicle_series`") || !strings.Contains(out, "`region`") {
		t.Fatalf("expected backtick identifiers, got:\n%s", out)
	}
	if strings.Contains(out, `"vehicle_series"`) {
		t.Fatalf("should not keep double quotes: %s", out)
	}
}

func TestDoubleQuotedIdentToBackticks_SkipsInsideSingleQuotedStrings(t *testing.T) {
	in := `INSERT INTO t VALUES ('say "hello"');`
	out := doubleQuotedIdentToBackticks(in)
	if !strings.Contains(out, `'say "hello"'`) {
		t.Fatalf("string literal must stay intact: %s", out)
	}
}

func TestDoubleQuotedIdentToBackticks_PostgresEscapedQuoteInIdent(t *testing.T) {
	in := `CREATE TABLE "a""b" (id INT);`
	out := doubleQuotedIdentToBackticks(in)
	if !strings.Contains(out, "`a\"b`") {
		t.Fatalf("expected embedded quote in identifier, got: %s", out)
	}
}

func TestConvert_PostgresToMariaDB_BacktickIdentifiers(t *testing.T) {
	in := `CREATE TABLE "vehicle_series" (series INT, vehicle INT);`
	out, err := Convert(in, dialect.Postgres, dialect.MariaDB)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, `"vehicle_series"`) || !strings.Contains(out, "`vehicle_series`") {
		t.Fatalf("unexpected: %s", out)
	}
}
