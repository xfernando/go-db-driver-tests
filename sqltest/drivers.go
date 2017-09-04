package sqltest

import (
	// so go get installs test dependencies when building the container
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gwenn/gosqlite"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jbarham/gopgsqldriver"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/mxk/go-sqlite/sqlite3"
	_ "github.com/ziutek/mymysql/godrv"
)
