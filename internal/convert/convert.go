package convert

import (
	"fmt"
	"regexp"

	"sql-migrator/internal/dialect"
)

type rule struct {
	re   *regexp.Regexp
	repl string
}

func Convert(input string, from, to dialect.Dialect) (string, error) {
	if from == to {
		return input, nil
	}
	var out string
	switch {
	case from == dialect.SQLite && to == dialect.Postgres:
		out = apply(input, sqliteToPostgres()...)
	case from == dialect.Postgres && to == dialect.SQLite:
		out = apply(input, postgresToSQLite()...)
	case from == dialect.SQLite && to == dialect.MariaDB:
		out = apply(input, sqliteToMariaDB()...)
	case from == dialect.MariaDB && to == dialect.SQLite:
		out = apply(input, mariaDBToSQLite()...)
	case from == dialect.Postgres && to == dialect.MariaDB:
		out = apply(input, postgresToMariaDB()...)
	case from == dialect.MariaDB && to == dialect.Postgres:
		out = apply(input, mariaDBToPostgres()...)
	case from == dialect.SQLite && to == dialect.SQLServer:
		out = apply(input, sqliteToSQLServer()...)
	case from == dialect.SQLServer && to == dialect.SQLite:
		out = apply(input, sqlServerToSQLite()...)
	case from == dialect.Postgres && to == dialect.SQLServer:
		out = apply(input, postgresToSQLServer()...)
	case from == dialect.SQLServer && to == dialect.Postgres:
		out = apply(input, sqlServerToPostgres()...)
	case from == dialect.MariaDB && to == dialect.SQLServer:
		out = apply(input, mariaDBToSQLServer()...)
	case from == dialect.SQLServer && to == dialect.MariaDB:
		out = apply(input, sqlServerToMariaDB()...)
	default:
		return "", fmt.Errorf("unsupported conversion %s → %s", from, to)
	}
	return out, nil
}

func apply(s string, rules ...rule) string {
	for _, r := range rules {
		s = r.re.ReplaceAllString(s, r.repl)
	}
	return s
}

func mustRE(pattern string) *regexp.Regexp {
	return regexp.MustCompile(`(?is)` + pattern)
}

func sqliteToPostgres() []rule {
	return []rule{
		{mustRE(`\bINTEGER\s+PRIMARY\s+KEY\s+AUTOINCREMENT\b`), `SERIAL PRIMARY KEY`},
		{mustRE(`\bAUTOINCREMENT\b`), ``},
		{mustRE(`\bDATETIME\b`), `TIMESTAMP`},
		{mustRE(`\bBLOB\b`), `BYTEA`},
		{mustRE(`\bREAL\b`), `DOUBLE PRECISION`},
		// INSERT OR IGNORE → INSERT ... ON CONFLICT DO NOTHING (simplified single-table form)
		{mustRE(`\bINSERT\s+OR\s+IGNORE\s+INTO\b`), `INSERT INTO`},
		// Common function differences
		{mustRE(`\bIFNULL\s*\(`), `COALESCE(`},
		{mustRE(`\bstrftime\s*\(\s*'%Y-%m-%d\s+%H:%M:%S'\s*,\s*'now'\s*\)`), `CURRENT_TIMESTAMP`},
	}
}

func postgresToSQLite() []rule {
	return []rule{
		{mustRE(`\bBIGSERIAL\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bSERIAL\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bBIGSERIAL\b`), `INTEGER`},
		{mustRE(`\bSERIAL\b`), `INTEGER`},
		{mustRE(`\bBYTEA\b`), `BLOB`},
		{mustRE(`\bDOUBLE\s+PRECISION\b`), `REAL`},
		{mustRE(`\bTIMESTAMP\s+WITH\s+TIME\s+ZONE\b`), `DATETIME`},
		{mustRE(`\bTIMESTAMP\s+WITHOUT\s+TIME\s+ZONE\b`), `DATETIME`},
		{mustRE(`\bTIMESTAMP\b`), `DATETIME`},
		{mustRE(`\bBOOLEAN\b`), `INTEGER`},
	}
}

func sqliteToMariaDB() []rule {
	return []rule{
		{mustRE(`\bINTEGER\s+PRIMARY\s+KEY\s+AUTOINCREMENT\b`), `INT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bAUTOINCREMENT\b`), `AUTO_INCREMENT`},
		{mustRE(`\bDATETIME\b`), `DATETIME`},
		// SQLite allows BOOLEAN as NUMERIC; MariaDB prefers TINYINT(1) or BOOL
		{mustRE(`\bBOOLEAN\b`), `BOOL`},
	}
}

func mariaDBToSQLite() []rule {
	return []rule{
		{mustRE(`\bINT\s+NOT\s+NULL\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bINT\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bINTEGER\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bAUTO_INCREMENT\b`), `AUTOINCREMENT`},
		// Strip MySQL-specific table options (best-effort)
		{mustRE(`\s+ENGINE\s*=\s*\w+`), ``},
		{mustRE(`\s+DEFAULT\s+CHARSET\s*=\s*\w+`), ``},
		{mustRE(`\s+COLLATE\s*=\s*\w+`), ``},
		{mustRE(`\s+CHARSET\s*=\s*\w+`), ``},
		{mustRE(`\bUNSIGNED\b`), ``},
		{mustRE(`\bBOOL\b`), `BOOLEAN`},
	}
}

func postgresToMariaDB() []rule {
	return []rule{
		{mustRE(`\bBIGSERIAL\s+PRIMARY\s+KEY\b`), `BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bSERIAL\s+PRIMARY\s+KEY\b`), `INT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bBIGSERIAL\b`), `BIGINT NOT NULL AUTO_INCREMENT`},
		{mustRE(`\bSERIAL\b`), `INT NOT NULL AUTO_INCREMENT`},
		{mustRE(`\bBYTEA\b`), `LONGBLOB`},
		{mustRE(`\bDOUBLE\s+PRECISION\b`), `DOUBLE`},
		{mustRE(`\bTIMESTAMP\s+WITHOUT\s+TIME\s+ZONE\b`), `DATETIME`},
		{mustRE(`\bBOOLEAN\b`), `BOOL`},
	}
}

func mariaDBToPostgres() []rule {
	return []rule{
		{mustRE(`\bBIGINT\s+NOT\s+NULL\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `BIGSERIAL PRIMARY KEY`},
		{mustRE(`\bINT\s+NOT\s+NULL\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `SERIAL PRIMARY KEY`},
		{mustRE(`\bINT\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `SERIAL PRIMARY KEY`},
		{mustRE(`\bBIGINT\s+NOT\s+NULL\s+AUTO_INCREMENT\b`), `BIGSERIAL`},
		{mustRE(`\bINT\s+NOT\s+NULL\s+AUTO_INCREMENT\b`), `SERIAL`},
		{mustRE(`\bAUTO_INCREMENT\b`), ``},
		{mustRE(`\s+ENGINE\s*=\s*\w+`), ``},
		{mustRE(`\s+DEFAULT\s+CHARSET\s*=\s*\w+`), ``},
		{mustRE(`\s+COLLATE\s*=\s*\w+`), ``},
		{mustRE(`\s+CHARSET\s*=\s*\w+`), ``},
		{mustRE(`\bLONGBLOB\b`), `BYTEA`},
		{mustRE(`\bMEDIUMBLOB\b`), `BYTEA`},
		{mustRE(`\bTINYBLOB\b`), `BYTEA`},
		{mustRE(`\bBLOB\b`), `BYTEA`},
		{mustRE(`\bDOUBLE\b`), `DOUBLE PRECISION`},
		{mustRE(`\bBOOL\b`), `BOOLEAN`},
		{mustRE(`\bUNSIGNED\b`), ``},
	}
}

func sqliteToSQLServer() []rule {
	return []rule{
		{mustRE(`\bINTEGER\s+PRIMARY\s+KEY\s+AUTOINCREMENT\b`), `INT IDENTITY(1,1) PRIMARY KEY`},
		{mustRE(`\bAUTOINCREMENT\b`), `IDENTITY(1,1)`},
		{mustRE(`\bDATETIME\b`), `DATETIME2`},
		{mustRE(`\bBLOB\b`), `VARBINARY(MAX)`},
		{mustRE(`\bREAL\b`), `FLOAT`},
		{mustRE(`\bBOOLEAN\b`), `BIT`},
		{mustRE(`\bIFNULL\s*\(`), `ISNULL(`},
	}
}

func sqlServerToSQLite() []rule {
	return []rule{
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+PRIMARY\s+KEY\b`), `INTEGER PRIMARY KEY AUTOINCREMENT`},
		{mustRE(`\bIDENTITY\s*\(\s*1\s*,\s*1\s*\)\b`), ``},
		{mustRE(`\bVARBINARY\s*\(\s*MAX\s*\)`), `BLOB`},
		{mustRE(`\bVARBINARY\s*\(\s*\d+\s*\)`), `BLOB`},
		{mustRE(`\bDATETIME2\b`), `DATETIME`},
		{mustRE(`\bSMALLDATETIME\b`), `DATETIME`},
		{mustRE(`\bFLOAT\b`), `REAL`},
		{mustRE(`\bBIT\b`), `INTEGER`},
		{mustRE(`\bISNULL\s*\(`), `IFNULL(`},
		{mustRE(`\s+ON\s+\[PRIMARY\s*\]`), ``},
	}
}

func postgresToSQLServer() []rule {
	return []rule{
		{mustRE(`\bBIGSERIAL\s+PRIMARY\s+KEY\b`), `BIGINT IDENTITY(1,1) PRIMARY KEY`},
		{mustRE(`\bSERIAL\s+PRIMARY\s+KEY\b`), `INT IDENTITY(1,1) PRIMARY KEY`},
		{mustRE(`\bBIGSERIAL\b`), `BIGINT IDENTITY(1,1) NOT NULL`},
		{mustRE(`\bSERIAL\b`), `INT IDENTITY(1,1) NOT NULL`},
		{mustRE(`\bBYTEA\b`), `VARBINARY(MAX)`},
		{mustRE(`\bDOUBLE\s+PRECISION\b`), `FLOAT`},
		{mustRE(`\bTIMESTAMP\s+WITH\s+TIME\s+ZONE\b`), `DATETIME2`},
		{mustRE(`\bTIMESTAMP\s+WITHOUT\s+TIME\s+ZONE\b`), `DATETIME2`},
		{mustRE(`\bTIMESTAMP\b`), `DATETIME2`},
		{mustRE(`\bTEXT\b`), `VARCHAR(MAX)`},
		{mustRE(`\bBOOLEAN\b`), `BIT`},
	}
}

func sqlServerToPostgres() []rule {
	return []rule{
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+PRIMARY\s+KEY\b`), `BIGSERIAL PRIMARY KEY`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+PRIMARY\s+KEY\b`), `SERIAL PRIMARY KEY`},
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\s+PRIMARY\s+KEY\b`), `BIGSERIAL PRIMARY KEY`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\s+PRIMARY\s+KEY\b`), `SERIAL PRIMARY KEY`},
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\b`), `BIGSERIAL`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\b`), `SERIAL`},
		{mustRE(`\bIDENTITY\s*\(\s*1\s*,\s*1\s*\)\b`), ``},
		{mustRE(`\bVARBINARY\s*\(\s*MAX\s*\)`), `BYTEA`},
		{mustRE(`\bVARCHAR\s*\(\s*MAX\s*\)`), `TEXT`},
		{mustRE(`\bNVARCHAR\s*\(\s*MAX\s*\)`), `TEXT`},
		{mustRE(`\bDATETIME2\b`), `TIMESTAMP`},
		{mustRE(`\bSMALLDATETIME\b`), `TIMESTAMP`},
		{mustRE(`\bDATETIME\b`), `TIMESTAMP`},
		{mustRE(`\bFLOAT\b`), `DOUBLE PRECISION`},
		{mustRE(`\bBIT\b`), `BOOLEAN`},
		{mustRE(`\bISNULL\s*\(`), `COALESCE(`},
		{mustRE(`\bGETDATE\s*\(\s*\)`), `CURRENT_TIMESTAMP`},
		{mustRE(`\s+ON\s+\[PRIMARY\s*\]`), ``},
	}
}

func mariaDBToSQLServer() []rule {
	return []rule{
		{mustRE(`\bBIGINT\s+NOT\s+NULL\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `BIGINT IDENTITY(1,1) PRIMARY KEY`},
		{mustRE(`\bINT\s+NOT\s+NULL\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `INT IDENTITY(1,1) PRIMARY KEY`},
		{mustRE(`\bINT\s+AUTO_INCREMENT\s+PRIMARY\s+KEY\b`), `INT IDENTITY(1,1) PRIMARY KEY`},
		{mustRE(`\bBIGINT\s+NOT\s+NULL\s+AUTO_INCREMENT\b`), `BIGINT IDENTITY(1,1) NOT NULL`},
		{mustRE(`\bINT\s+NOT\s+NULL\s+AUTO_INCREMENT\b`), `INT IDENTITY(1,1) NOT NULL`},
		{mustRE(`\bAUTO_INCREMENT\b`), `IDENTITY(1,1)`},
		{mustRE(`\s+ENGINE\s*=\s*\w+`), ``},
		{mustRE(`\s+DEFAULT\s+CHARSET\s*=\s*\w+`), ``},
		{mustRE(`\s+COLLATE\s*=\s*\w+`), ``},
		{mustRE(`\s+CHARSET\s*=\s*\w+`), ``},
		{mustRE(`\bLONGBLOB\b`), `VARBINARY(MAX)`},
		{mustRE(`\bMEDIUMBLOB\b`), `VARBINARY(MAX)`},
		{mustRE(`\bTINYBLOB\b`), `VARBINARY(MAX)`},
		{mustRE(`\bBLOB\b`), `VARBINARY(MAX)`},
		{mustRE(`\bDOUBLE\b`), `FLOAT`},
		{mustRE(`\bBOOL\b`), `BIT`},
		{mustRE(`\bUNSIGNED\b`), ``},
	}
}

func sqlServerToMariaDB() []rule {
	return []rule{
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+PRIMARY\s+KEY\b`), `BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+PRIMARY\s+KEY\b`), `INT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\s+PRIMARY\s+KEY\b`), `BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\s+PRIMARY\s+KEY\b`), `INT NOT NULL AUTO_INCREMENT PRIMARY KEY`},
		{mustRE(`\bBIGINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\b`), `BIGINT NOT NULL AUTO_INCREMENT`},
		{mustRE(`\bINT\s+IDENTITY\s*\(\s*1\s*,\s*1\s*\)\s+NOT\s+NULL\b`), `INT NOT NULL AUTO_INCREMENT`},
		{mustRE(`\bIDENTITY\s*\(\s*1\s*,\s*1\s*\)\b`), `AUTO_INCREMENT`},
		{mustRE(`\bVARBINARY\s*\(\s*MAX\s*\)`), `LONGBLOB`},
		{mustRE(`\bFLOAT\b`), `DOUBLE`},
		{mustRE(`\bBIT\b`), `BOOL`},
		{mustRE(`\bISNULL\s*\(`), `IFNULL(`},
		{mustRE(`\s+ON\s+\[PRIMARY\s*\]`), ``},
	}
}
