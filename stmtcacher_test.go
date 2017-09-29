package stmtcacher

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCachingProxyPrepareSqlite(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	proxy := NewCachingProxy(db)
	query := "SELECT 1"

	proxy.Prepare(query)
}

func TestCachingWrapperPrepareSqlite(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	wrapper := NewCachingWrapper(db)
	query := "SELECT 1"

	wrapper.ExecPrepared(query)
	wrapper.Exec(query)
}
