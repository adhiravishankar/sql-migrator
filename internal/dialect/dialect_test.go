package dialect

import "testing"

func TestParse_MySQLAliasMariaDB(t *testing.T) {
	d, err := Parse("mysql")
	if err != nil || d != MariaDB {
		t.Fatalf("Parse(mysql) = %v, %v (want MariaDB)", d, err)
	}
}
