# stmtcacher

Prepared statement caching for go. Works as a proxy for sql.DB,
caching and reusing all prepared statements.

Simplified and based on:
https://github.com/Masterminds/squirrel/blob/v1/stmtcacher.go

Compared to squirrel, this pkg has:
* Better error handling
* Doesn't enforce squirrel-specific types like RowScanner
* Easier use with other pkgs that wrap sql.DB

``` go
db, err := sql.Open("sqlite3", ":memory:")
if err != nil {
	// ...
}

proxy := stmtcacher.NewStmtCacher(db)

query := "SELECT * FROM table where id = ?"

// By default, sql.DB will perform 3 statements and roundtrips
db.Query(query, 1)
// PREPARE
// EXECUTE
// CLOSE
db.Query(query, 1)
// PREPARE
// EXECUTE
// CLOSE

// With stmtcacher, you reuse the statements once prepared
proxy.Query(query, 1)
// PREPARE
// EXECUTE
proxy.Query(query, 1)
// EXECUTE
proxy.Query(query, 1)
// EXECUTE
```
