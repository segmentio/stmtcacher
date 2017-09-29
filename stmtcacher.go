package stmtcacher

import (
	"context"
	"database/sql"
	"sync"
)

// Interface consisting of the necessary subset of sql.DB functions required for the wrapper and proxy
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Wrapper around DB with additional functions that explicitly reuse cached
// prepared statements. All sql.DB query functions are wrapped and available
// with the prefix "Prepared", e.g. PreparedQuery, PreparedExec
type CachingWrapper struct {
	proxy *CachingProxy
	DB
}

// Returns a CachingWrapper
func NewCachingWrapper(db DB) *CachingWrapper {
	proxy := NewCachingProxy(db)
	return &CachingWrapper{proxy: proxy, DB: db}
}

// sql.DB.Exec with caching and resuse of the generated prepared statement
func (wrapper *CachingWrapper) PreparedExec(query string, args ...interface{}) (res sql.Result, err error) {
	return wrapper.proxy.Exec(query, args...)
}

// sql.DB.Query with caching and resuse of the generated prepared statement
func (wrapper *CachingWrapper) PreparedQuery(query string, args ...interface{}) (rows *sql.Rows, err error) {
	return wrapper.proxy.Query(query, args...)
}

// sql.DB.QueryRow with caching and resuse of the generated prepared statement
func (wrapper *CachingWrapper) PreparedQueryRow(query string, args ...interface{}) *sql.Row {
	return wrapper.proxy.QueryRow(query, args...)
}

// sql.DB.ExecContext with caching and resuse of the generated prepared statement
func (wrapper *CachingWrapper) PreparedExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	return wrapper.proxy.ExecContext(ctx, query, args...)
}

// sql.DB.QueryContext with caching and resuse of the generated prepared statement
func (wrapper *CachingWrapper) PreparedQueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	return wrapper.proxy.QueryContext(ctx, query, args...)
}

// sql.DB.QueryRowContext with caching and resuse of the generated prepared statement
func (wrapper *CachingWrapper) PreparedQueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return wrapper.proxy.QueryRowContext(ctx, query, args...)
}

// Caching proxy and wrapper around sql.DB. Exec, Query, QueryRow, ExecContext,
// QueryContext, and QueryRowContext are each wrapped to cache and reuse
// prepared statements
type CachingProxy struct {
	cache map[string]*sql.Stmt
	mu    sync.Mutex
	DB
}

// Returns a CachingProxy
func NewCachingProxy(db DB) *CachingProxy {
	return &CachingProxy{cache: make(map[string]*sql.Stmt), DB: db}
}

// sql.DB.Prepare that checks an object cache before creating a new prepared statement
func (proxy *CachingProxy) Prepare(query string) (*sql.Stmt, error) {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	stmt, ok := proxy.cache[query]
	if ok {
		return stmt, nil
	}
	stmt, err := proxy.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	proxy.cache[query] = stmt
	return stmt, nil
}

// sql.DB.Exec with caching and resuse of the generated prepared statement
func (proxy *CachingProxy) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	stmt, err := proxy.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(args...)
}

// sql.DB.Query with caching and resuse of the generated prepared statement
func (proxy *CachingProxy) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := proxy.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Query(args...)
}

// sql.DB.QueryRow with caching and resuse of the generated prepared statement
func (proxy *CachingProxy) QueryRow(query string, args ...interface{}) *sql.Row {
	stmt, err := proxy.Prepare(query)
	if err != nil {
		return nil
	}
	return stmt.QueryRow(args...)
}

// sql.DB.PrepareContext that checks an object cache before creating a new prepared statement
func (proxy *CachingProxy) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	stmt, ok := proxy.cache[query]
	if ok {
		return stmt, nil
	}
	stmt, err := proxy.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	proxy.cache[query] = stmt
	return stmt, nil
}

// sql.DB.ExecContext with caching and resuse of the generated prepared statement
func (proxy *CachingProxy) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	stmt, err := proxy.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.ExecContext(ctx, args...)
}

// sql.DB.QueryContext with caching and resuse of the generated prepared statement
func (proxy *CachingProxy) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := proxy.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.QueryContext(ctx, args...)
}

// sql.DB.QueryRowContext with caching and resuse of the generated prepared statement
func (proxy *CachingProxy) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	stmt, err := proxy.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}
	return stmt.QueryRowContext(ctx, args...)
}
