package sqltest

import (
	// so go get installs test dependencies when building the container
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gwenn/gosqlite"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jbarham/gopgsqldriver"
	_ "github.com/lib/pq"
	// "github.com/mattn/go-sqlite3"
	_ "github.com/minus5/gofreetds"
	_ "github.com/mxk/go-sqlite/sqlite3"
	_ "github.com/nakagami/firebirdsql"
	_ "github.com/ziutek/mymysql/godrv"
)
