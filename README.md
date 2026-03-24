# sql-migrator

A small command-line tool that reads a SQL file written for one database dialect and emits a **best-effort** translation for another dialect. Supported dialects are **SQLite**, **PostgreSQL**, **MariaDB** (MySQL-compatible), and **Microsoft SQL Server**.

Conversion is implemented with ordered pattern rules, not a full SQL parser. It works well for common DDL (especially `CREATE TABLE`) and simple statements; stored procedures, dialect-specific functions, and complex scripts may need manual follow-up.

## Requirements

- [Go](https://go.dev/dl/) (see `go.mod` for the required version).

## Build

```bash
go build -o sql-migrator .
```

If the build fails with a VCS stamping error in a directory that is not a Git repository, use:

```bash
go build -buildvcs=false -o sql-migrator .
```

## Usage

```text
sql-migrator -input <file.sql> -from <dialect> -to <dialect> [-output <out.sql>]
```

| Flag | Description |
|------|-------------|
| `-input` | Path to the source `.sql` file (required). |
| `-from` | Source dialect (required). |
| `-to` | Target dialect (required). |
| `-output` | Path to write the converted SQL. If omitted, the result is printed to **stdout**. |

### Dialect names and aliases

| Dialect | Accepted names |
|---------|----------------|
| SQLite | `sqlite`, `sqlite3` |
| PostgreSQL | `postgres`, `postgresql`, `pg` |
| MariaDB / MySQL | `mariadb`, `mysql` |
| SQL Server | `sqlserver`, `mssql`, `tsql`, `t-sql` |

### Examples

Convert SQLite-style DDL to PostgreSQL and print to the terminal:

```bash
./sql-migrator -input schema.sql -from sqlite -to postgres
```

Write MariaDB output to a file:

```bash
./sql-migrator -input source.sql -from postgres -to mariadb -output out.sql
```

SQL Server to SQLite:

```bash
./sql-migrator -input tsql.sql -from sqlserver -to sqlite -output sqlite.sql
```

## Tests

```bash
go test ./...
```

## License

See [license.md](license.md).

## Contributing

See [contributing.md](contributing.md).
