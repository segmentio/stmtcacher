package stmtcacher

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type DBStub struct {
	LastPrepareSql string
	PrepareCount   int
}

func (db *DBStub) Prepare(query string) (*sql.Stmt, error) {
	db.LastPrepareSql = query
	db.PrepareCount++
	return nil, nil
}

func (db *DBStub) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.Prepare(query)
}

func (db *DBStub) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (db *DBStub) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func (db *DBStub) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (db *DBStub) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func (db *DBStub) QueryRow(query string, args ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (db *DBStub) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func TestCachingProxyPrepare(t *testing.T) {
	db := &DBStub{}
	proxy := NewCachingProxy(db)
	query := "SELECT 1"

	proxy.Prepare(query)
	proxy.Prepare(query)

	assert.Equal(t, query, db.LastPrepareSql)
	assert.Equal(t, 1, db.PrepareCount)
}

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
