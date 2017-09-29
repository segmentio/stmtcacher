# stmtcacher

Prepared statement caching for go's sql.DB.

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/segmentio/stmtcacher)

## CachingWrapper

Wrapper around DB with additional functions that explicitly reuse cached
prepared statements. All sql.DB query functions are wrapped and available
with the suffix "Prepared", e.g. QueryPrepared, ExecPrepared

``` go
db, err := sql.Open("sqlite3", ":memory:")
if err != nil {
	// ...
}

wrapper := stmtcacher.NewCachingWrapper(db)

query := "SELECT * FROM table where id = ?"

// By default, sql.DB will perform 3 statements and roundtrips
db.Query(query, "foo")
// PREPARE 'SELECT * FROM table where name = ?'
// EXECUTE stmt1, "foo"
// CLOSE stmt1
db.Query(query, "bar")
// PREPARE 'SELECT * FROM table where name = ?'
// EXECUTE stmt2, "bar"
// CLOSE

// With stmtcacher, they're cached
wrapper.QueryPrepared(query, "foo")
// PREPARE 'SELECT * FROM table where name = ?'
// EXECUTE stmt1, "foo"
wrapper.QueryPrepared(query, "bar")
// EXECUTE stmt1, "bar"
wrapper.QueryPrepared(query, "baz")
// EXECUTE stmt1, "baz"
```

## CachingProxy

Caching proxy and wrapper around sql.DB. Exec, Query, QueryRow, ExecContext,
QueryContext, and QueryRowContext are each wrapped to cache and reuse
prepared statements

``` go
db, err := sql.Open("sqlite3", ":memory:")
if err != nil {
	// ...
}

proxy := stmtcacher.NewCachingProxy(db)

query := "SELECT * FROM table where name = ?"

// By default, sql.DB will perform 3 statements and roundtrips
db.Query(query, "foo")
// PREPARE 'SELECT * FROM table where name = ?'
// EXECUTE stmt1, "foo"
// CLOSE stmt1
db.Query(query, "bar")
// PREPARE 'SELECT * FROM table where name = ?'
// EXECUTE stmt2, "bar"
// CLOSE

// With stmtcacher, they're cached
proxy.Query(query, "foo")
// PREPARE 'SELECT * FROM table where name = ?'
// EXECUTE stmt1, "foo"
proxy.Query(query, "bar")
// EXECUTE stmt1, "bar"
proxy.Query(query, "baz")
// EXECUTE stmt1, "baz"
```

#### Comparison

Simplified and inspired by: https://github.com/Masterminds/squirrel

Compared to squirrel, this pkg has:
* Better error handling
* Doesn't enforce squirrel-specific types like RowScanner
* Easier use with other pkgs that wrap sql.DB
* Offers both a wrapper with additional fns, as well as caching proxy
