package model

// Script is a neutral representation of a SQL file: ordered statements that
// the emitter can render for a target dialect.
type Script struct {
	Statements []Statement
}

// Statement is implemented by concrete statement types.
type Statement interface {
	stmt()
}

// Raw is a statement we did not structurally parse; it is passed through the
// legacy text transform.
type Raw struct {
	SQL string
}

func (*Raw) stmt() {}

// Ident is a table or column name with optional quoting in the source SQL.
type Ident struct {
	Name   string
	Quoted bool
}

// CreateTable is a parsed CREATE TABLE (subset: columns + optional table-level lines).
type CreateTable struct {
	IfNotExists bool
	Table       Ident
	Columns     []ColumnDef
	// TableLevel holds lines such as PRIMARY KEY (...) or FOREIGN KEY ... that
	// are not column definitions.
	TableLevel []string
}

func (*CreateTable) stmt() {}

// ColumnDef is one column definition inside CREATE TABLE (...).
type ColumnDef struct {
	Name Ident
	// Rest is everything after the column name (type, NOT NULL, DEFAULT, etc.).
	Rest string
}

// Insert is a single-row INSERT ... VALUES (...) statement.
type Insert struct {
	OrIgnore  bool
	OrReplace bool
	Table     Ident
	Columns   []Ident
	Values    string // text starting at the opening '(' of the first value tuple
}

func (*Insert) stmt() {}
