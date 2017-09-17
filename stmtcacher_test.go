package stmtcacher

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: write tests
func TestStmtCacherPrepare(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cacher := NewStmtCacher(db)
	query := "SELECT 1"

	cacher.Prepare(query)
}
