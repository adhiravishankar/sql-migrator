package convert

import (
	"strings"
	"testing"

	"sql-migrator/internal/dialect"
)

func TestConvert_SQLiteToPostgres(t *testing.T) {
	in := `CREATE TABLE t (id INTEGER PRIMARY KEY AUTOINCREMENT, data BLOB);`
	out, err := Convert(in, dialect.SQLite, dialect.Postgres)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "SERIAL PRIMARY KEY") || !strings.Contains(out, "BYTEA") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestConvert_SameDialect(t *testing.T) {
	in := `SELECT 1;`
	out, err := Convert(in, dialect.Postgres, dialect.Postgres)
	if err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Fatalf("expected unchanged, got %q", out)
	}
}

func TestConvert_MariaDBToPostgres_StripsEngine(t *testing.T) {
	in := `CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	out, err := Convert(in, dialect.MariaDB, dialect.Postgres)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "ENGINE") || strings.Contains(out, "CHARSET") {
		t.Fatalf("expected engine/charset stripped: %s", out)
	}
}

func TestConvert_PostgresToSQLServer(t *testing.T) {
	in := `CREATE TABLE t (id SERIAL PRIMARY KEY, body TEXT);`
	out, err := Convert(in, dialect.Postgres, dialect.SQLServer)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "IDENTITY(1,1)") || !strings.Contains(out, "VARCHAR(MAX)") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestConvert_SQLServerToSQLite(t *testing.T) {
	in := `CREATE TABLE t (id INT IDENTITY(1,1) PRIMARY KEY, data VARBINARY(MAX)) ON [PRIMARY];`
	out, err := Convert(in, dialect.SQLServer, dialect.SQLite)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "INTEGER PRIMARY KEY AUTOINCREMENT") || !strings.Contains(out, "BLOB") {
		t.Fatalf("unexpected output: %s", out)
	}
	if strings.Contains(out, "ON [PRIMARY]") {
		t.Fatalf("expected ON [PRIMARY] stripped: %s", out)
	}
}

func TestDialect_ParseSQLServer(t *testing.T) {
	for _, s := range []string{"sqlserver", "mssql", "tsql", "T-SQL"} {
		d, err := dialect.Parse(s)
		if err != nil || d != dialect.SQLServer {
			t.Fatalf("Parse(%q) = %v, %v", s, d, err)
		}
	}
}
