package sqltest

import (
	// so go get installs test dependencies when building the container
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/ziutek/mymysql/godrv"
)
