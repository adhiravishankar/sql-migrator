package dialect

import (
	"fmt"
	"strings"
)

// Dialect names supported for conversion.
type Dialect int

const (
	Unknown Dialect = iota
	SQLite
	Postgres
	MariaDB
	SQLServer
)

func (d Dialect) String() string {
	switch d {
	case SQLite:
		return "sqlite"
	case Postgres:
		return "postgres"
	case MariaDB:
		return "mariadb"
	case SQLServer:
		return "sqlserver"
	default:
		return "unknown"
	}
}

// Parse accepts common aliases (e.g. postgresql, mysql, mariadb).
func Parse(s string) (Dialect, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "sqlite", "sqlite3":
		return SQLite, nil
	case "postgres", "postgresql", "pg":
		return Postgres, nil
	case "mariadb", "mysql":
		return MariaDB, nil
	case "sqlserver", "mssql", "tsql", "t-sql":
		return SQLServer, nil
	default:
		return Unknown, fmt.Errorf("unknown dialect %q (use sqlite, postgres, mariadb or mysql, or sqlserver)", s)
	}
}
