package stmtcacher

import (
	"context"
	"database/sql"
	"sync"
)

type StmtCacher struct {
	cache map[string]*sql.Stmt
	mu    sync.Mutex
	*sql.DB
}

func NewStmtCacher(db *sql.DB) *StmtCacher {
	return &StmtCacher{cache: make(map[string]*sql.Stmt), DB: db}
}

func (sc *StmtCacher) Prepare(query string) (*sql.Stmt, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	stmt, ok := sc.cache[query]
	if ok {
		return stmt, nil
	}
	stmt, err := sc.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	sc.cache[query] = stmt
	return stmt, nil
}

func (sc *StmtCacher) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	stmt, err := sc.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(args...)
}

func (sc *StmtCacher) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := sc.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Query(args...)
}

func (sc *StmtCacher) QueryRow(query string, args ...interface{}) *sql.Row {
	stmt, err := sc.Prepare(query)
	if err != nil {
		return nil
	}
	return stmt.QueryRow(args...)
}

func (sc *StmtCacher) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	stmt, ok := sc.cache[query]
	if ok {
		return stmt, nil
	}
	stmt, err := sc.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	sc.cache[query] = stmt
	return stmt, nil
}

func (sc *StmtCacher) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	stmt, err := sc.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.ExecContext(ctx, args...)
}

func (sc *StmtCacher) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := sc.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.QueryContext(ctx, args...)
}

func (sc *StmtCacher) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	stmt, err := sc.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}
	return stmt.QueryRowContext(ctx, args...)
}
