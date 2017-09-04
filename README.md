<img src="img/hugging-docker.png" width="140px" height="195px">[1]

## go-db-driver-tests

This project is based on Brad Fitz's [earlier work](https://github.com/bradfitz/go-sql-test) and aims to provide a testing suite for go [database drivers
](https://github.com/golang/go/wiki/SQLDrivers) using docker to start the database servers needed to run all the tests.

Progress:

MySQL drivers tested: 2 out of 2

Postgres drivers tested: 3 out of 3 (1 failing)

SQLite drivers tested: 2 out of 3

## Usage

We use `docker-compose` to build the image with the test code and start all the database containers needed for the tests to run:

```bash
docker-compose build
docker-compose up -d
```

Then, you can see the output of the test run with:
```bash
docker-compose logs -f godbtests
```

When the tests are completed, run the following to stop the database containers:
```bash
docker-compose down
```

## Output of latest run

```
godbtests    | === RUN   TestMySQLDrivers
godbtests    | === RUN   TestMySQLDrivers/gomysql:_TXQuery
godbtests    | === RUN   TestMySQLDrivers/gomysql:_Blobs
godbtests    | === RUN   TestMySQLDrivers/gomysql:_ManyQueryRow
godbtests    | === RUN   TestMySQLDrivers/gomysql:_PreparedStmt
godbtests    | === RUN   TestMySQLDrivers/mymysql:_TXQuery
godbtests    | === RUN   TestMySQLDrivers/mymysql:_Blobs
godbtests    | === RUN   TestMySQLDrivers/mymysql:_ManyQueryRow
godbtests    | === RUN   TestMySQLDrivers/mymysql:_PreparedStmt
godbtests    | --- PASS: TestMySQLDrivers (10.46s)
godbtests    | 	mysql_test.go:80: github.com/go-sql-driver/mysql revision: 26471af196a17ee75a22e6481b5a5897fb16b081
godbtests    |     --- PASS: TestMySQLDrivers/gomysql:_TXQuery (0.42s)
godbtests    |     	mysql_test.go:50: Driver type: *mysql.MySQLDriver
godbtests    |     --- PASS: TestMySQLDrivers/gomysql:_Blobs (0.56s)
godbtests    |     	mysql_test.go:50: Driver type: *mysql.MySQLDriver
godbtests    |     --- PASS: TestMySQLDrivers/gomysql:_ManyQueryRow (2.82s)
godbtests    |     	mysql_test.go:50: Driver type: *mysql.MySQLDriver
godbtests    |     --- PASS: TestMySQLDrivers/gomysql:_PreparedStmt (1.39s)
godbtests    |     	mysql_test.go:50: Driver type: *mysql.MySQLDriver
godbtests    | 	mysql_test.go:86: github.com/ziutek/mymysql revision: 1d19cbf98d83564cc561192ae7d7183d795f7ac7
godbtests    |     --- PASS: TestMySQLDrivers/mymysql:_TXQuery (0.49s)
godbtests    |     	mysql_test.go:50: Driver type: *godrv.Driver
godbtests    |     --- PASS: TestMySQLDrivers/mymysql:_Blobs (0.52s)
godbtests    |     	mysql_test.go:50: Driver type: *godrv.Driver
godbtests    |     --- PASS: TestMySQLDrivers/mymysql:_ManyQueryRow (2.55s)
godbtests    |     	mysql_test.go:50: Driver type: *godrv.Driver
godbtests    |     --- PASS: TestMySQLDrivers/mymysql:_PreparedStmt (1.71s)
godbtests    |     	mysql_test.go:50: Driver type: *godrv.Driver
godbtests    | === RUN   TestPostgresDrivers
godbtests    | === RUN   TestPostgresDrivers/pgx:_TXQuery
godbtests    | === RUN   TestPostgresDrivers/pgx:_Blobs
godbtests    | === RUN   TestPostgresDrivers/pgx:_ManyQueryRow
godbtests    | === RUN   TestPostgresDrivers/pgx:_PreparedStmt
godbtests    | === RUN   TestPostgresDrivers/pq:_TXQuery
godbtests    | === RUN   TestPostgresDrivers/pq:_Blobs
godbtests    | === RUN   TestPostgresDrivers/pq:_ManyQueryRow
godbtests    | === RUN   TestPostgresDrivers/pq:_PreparedStmt
godbtests    | === RUN   TestPostgresDrivers/gopg:_TXQuery
godbtests    | === RUN   TestPostgresDrivers/gopg:_Blobs
godbtests    | === RUN   TestPostgresDrivers/gopg:_ManyQueryRow
godbtests    | === RUN   TestPostgresDrivers/gopg:_PreparedStmt
godbtests    | --- FAIL: TestPostgresDrivers (9.04s)
godbtests    | 	postgres_test.go:89: github.com/jackc/pgx revision: 9c8ef1acddff5f99e8cab446c0fd30a413ba69ab
godbtests    |     --- PASS: TestPostgresDrivers/pgx:_TXQuery (0.09s)
godbtests    |     	postgres_test.go:47: Driver type: *stdlib.Driver
godbtests    |     --- PASS: TestPostgresDrivers/pgx:_Blobs (0.14s)
godbtests    |     	postgres_test.go:47: Driver type: *stdlib.Driver
godbtests    |     --- PASS: TestPostgresDrivers/pgx:_ManyQueryRow (0.32s)
godbtests    |     	postgres_test.go:47: Driver type: *stdlib.Driver
godbtests    |     --- PASS: TestPostgresDrivers/pgx:_PreparedStmt (3.85s)
godbtests    |     	postgres_test.go:47: Driver type: *stdlib.Driver
godbtests    | 	postgres_test.go:95: github.com/lib/pq revision: e42267488fe361b9dc034be7a6bffef5b195bceb
godbtests    |     --- PASS: TestPostgresDrivers/pq:_TXQuery (0.10s)
godbtests    |     	postgres_test.go:47: Driver type: *pq.Driver
godbtests    |     --- PASS: TestPostgresDrivers/pq:_Blobs (0.23s)
godbtests    |     	postgres_test.go:47: Driver type: *pq.Driver
godbtests    |     --- PASS: TestPostgresDrivers/pq:_ManyQueryRow (3.98s)
godbtests    |     	postgres_test.go:47: Driver type: *pq.Driver
godbtests    |     --- PASS: TestPostgresDrivers/pq:_PreparedStmt (0.33s)
godbtests    |     	postgres_test.go:47: Driver type: *pq.Driver
godbtests    | 	postgres_test.go:101: github.com/jbarham/gopgsqldriver revision: f8287ee9bfe224aa4a7edcd73815ecbe69db7f68
godbtests    |     --- FAIL: TestPostgresDrivers/gopg:_TXQuery (0.00s)
godbtests    |     	postgres_test.go:47: Driver type: *pgsqldriver.postgresDriver
godbtests    |     	postgres_test.go:53: failed to enumerate tables: conn error:could not connect to server: No such file or directory
godbtests    |     			Is the server running locally and accepting
godbtests    |     			connections on Unix domain socket "/var/run/postgresql/.s.PGSQL.5432"?
godbtests    |     --- FAIL: TestPostgresDrivers/gopg:_Blobs (0.00s)
godbtests    |     	postgres_test.go:47: Driver type: *pgsqldriver.postgresDriver
godbtests    |     	postgres_test.go:53: failed to enumerate tables: conn error:could not connect to server: No such file or directory
godbtests    |     			Is the server running locally and accepting
godbtests    |     			connections on Unix domain socket "/var/run/postgresql/.s.PGSQL.5432"?
godbtests    |     --- FAIL: TestPostgresDrivers/gopg:_ManyQueryRow (0.00s)
godbtests    |     	postgres_test.go:47: Driver type: *pgsqldriver.postgresDriver
godbtests    |     	postgres_test.go:53: failed to enumerate tables: conn error:could not connect to server: No such file or directory
godbtests    |     			Is the server running locally and accepting
godbtests    |     			connections on Unix domain socket "/var/run/postgresql/.s.PGSQL.5432"?
godbtests    |     --- FAIL: TestPostgresDrivers/gopg:_PreparedStmt (0.00s)
godbtests    |     	postgres_test.go:47: Driver type: *pgsqldriver.postgresDriver
godbtests    |     	postgres_test.go:53: failed to enumerate tables: conn error:could not connect to server: No such file or directory
godbtests    |     			Is the server running locally and accepting
godbtests    |     			connections on Unix domain socket "/var/run/postgresql/.s.PGSQL.5432"?
godbtests    | === RUN   TestSqliteDrivers
godbtests    | === RUN   TestSqliteDrivers/gwenn_sqlite3:_TXQuery
godbtests    | === RUN   TestSqliteDrivers/gwenn_sqlite3:_Blobs
godbtests    | === RUN   TestSqliteDrivers/gwenn_sqlite3:_ManyQueryRow
godbtests    | === RUN   TestSqliteDrivers/gwenn_sqlite3:_PreparedStmt
godbtests    | === RUN   TestSqliteDrivers/mkx_sqlite3:_TXQuery
godbtests    | === RUN   TestSqliteDrivers/mkx_sqlite3:_Blobs
godbtests    | === RUN   TestSqliteDrivers/mkx_sqlite3:_ManyQueryRow
godbtests    | === RUN   TestSqliteDrivers/mkx_sqlite3:_PreparedStmt
godbtests    | --- FAIL: TestSqliteDrivers (24.66s)
godbtests    | 	sqlite_test.go:64: github.com/gwenn/gosqlite revision: 9d694b0a6a3946d24b601e669ab4a2e6117504c2
godbtests    |     --- PASS: TestSqliteDrivers/gwenn_sqlite3:_TXQuery (0.17s)
godbtests    |     	sqlite_test.go:55: Driver type: *sqlite.impl
godbtests    |     --- PASS: TestSqliteDrivers/gwenn_sqlite3:_Blobs (0.20s)
godbtests    |     	sqlite_test.go:55: Driver type: *sqlite.impl
godbtests    |     --- PASS: TestSqliteDrivers/gwenn_sqlite3:_ManyQueryRow (0.50s)
godbtests    |     	sqlite_test.go:55: Driver type: *sqlite.impl
godbtests    |     --- PASS: TestSqliteDrivers/gwenn_sqlite3:_PreparedStmt (13.55s)
godbtests    |     	sqlite_test.go:55: Driver type: *sqlite.impl
godbtests    | 	sqlite_test.go:76: github.com/mxk/go-sqlite revision: 167da9432e1f4602e95ea67b67051cfa34412e3f
godbtests    |     --- PASS: TestSqliteDrivers/mkx_sqlite3:_TXQuery (0.12s)
godbtests    |     	sqlite_test.go:55: Driver type: sqlite3.Driver
godbtests    |     --- PASS: TestSqliteDrivers/mkx_sqlite3:_Blobs (0.27s)
godbtests    |     	sqlite_test.go:55: Driver type: sqlite3.Driver
godbtests    |     --- PASS: TestSqliteDrivers/mkx_sqlite3:_ManyQueryRow (0.64s)
godbtests    |     	sqlite_test.go:55: Driver type: sqlite3.Driver
godbtests    |     --- FAIL: TestSqliteDrivers/mkx_sqlite3:_PreparedStmt (9.21s)
godbtests    |     	sqlite_test.go:55: Driver type: sqlite3.Driver
godbtests    |     	sql_test.go:90: Insert: sqlite3: database is locked [5]
godbtests    |     	sql_test.go:90: Insert: sqlite3: database is locked [5]
godbtests    |     	sql_test.go:90: Insert: sqlite3: database is locked [5]
godbtests    | FAIL
godbtests    | exit status 1
godbtests    | FAIL	app/sqltest	44.192s
```

[1] Image copied from https://github.com/egonelbre/gophers
