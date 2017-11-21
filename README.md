<img src="img/hugging-docker.png" width="140px" height="195px">[1]

## go-db-driver-tests

This project is based on Brad Fitz's [earlier work](https://github.com/bradfitz/go-sql-test) and aims to provide a testing suite for go [database drivers
](https://github.com/golang/go/wiki/SQLDrivers) using docker to start the database servers needed to run all the tests.

Progress:

MySQL drivers tested: 2 out of 2 (passing all tests)

Postgres drivers tested: 3 out of 3 (passing all tests)

SQLite drivers tested: 2 out of 3 (passing all tests)

Microsoft SQL Server drivers: 2 out 2 (passing all tests)

Firebird drivers: 1 out of 1 (failing concurrent read/write test)


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
godbtests    | === RUN   TestFirebirdDrivers
godbtests    | === RUN   TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_Transaction
godbtests    | === RUN   TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_Blobs
godbtests    | === RUN   TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_ConcurrentPreparedReadWrites
godbtests    | --- FAIL: TestFirebirdDrivers (5.73s)
godbtests    |  firebird_test.go:59: github.com/nakagami/firebirdsql revision: 8e47fbdec3e05de2cd21de138dba93e46a6336ee
godbtests    |     --- PASS: TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_Transaction (0.11s)
godbtests    |     --- PASS: TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_Blobs (0.09s)
godbtests    |     --- PASS: TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_InsertOnceReadOneThousandTimes (5.38s)
godbtests    |     --- FAIL: TestFirebirdDrivers/github.com/nakagami/firebirdsql_-_ConcurrentPreparedReadWrites (0.15s)
godbtests    |          firebird_test.go:180: Query: Error op_response:0
godbtests    |          firebird_test.go:180: Query: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    |          firebird_test.go:184: Insert: Error op_response:0
godbtests    | === RUN   TestMsSQLDrivers
godbtests    | === RUN   TestMsSQLDrivers/github.com/minus5/gofreetds:_Transaction
godbtests    | === RUN   TestMsSQLDrivers/github.com/minus5/gofreetds:_Blobs
godbtests    | === RUN   TestMsSQLDrivers/github.com/minus5/gofreetds:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestMsSQLDrivers/github.com/minus5/gofreetds:_ConcurrentPreparedReadWrites
godbtests    | === RUN   TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_Transaction
godbtests    | === RUN   TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_Blobs
godbtests    | === RUN   TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_ConcurrentPreparedReadWrites
godbtests    | --- PASS: TestMsSQLDrivers (5.05s)
godbtests    |  mssql_test.go:63: github.com/minus5/gofreetds revision: d9ec41707cd222494a1dddc7ffc128d014879053
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/minus5/gofreetds:_Transaction (0.04s)
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/minus5/gofreetds:_Blobs (0.04s)
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/minus5/gofreetds:_InsertOnceReadOneThousandTimes (2.59s)
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/minus5/gofreetds:_ConcurrentPreparedReadWrites (0.14s)
godbtests    |  mssql_test.go:69: github.com/denisenkom/go-mssqldb revision: 88555645b640cc621e32f8693d7586a1aa1575f4
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_Transaction (0.02s)
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_Blobs (0.02s)
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_InsertOnceReadOneThousandTimes (2.11s)
godbtests    |     --- PASS: TestMsSQLDrivers/github.com/denisenkom/go-mssqldb:_ConcurrentPreparedReadWrites (0.09s)
godbtests    | === RUN   TestMySQLDrivers
godbtests    | === RUN   TestMySQLDrivers/github.com/go-sql-driver/mysql:_Transaction
godbtests    | === RUN   TestMySQLDrivers/github.com/go-sql-driver/mysql:_Blobs
godbtests    | === RUN   TestMySQLDrivers/github.com/go-sql-driver/mysql:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestMySQLDrivers/github.com/go-sql-driver/mysql:_ConcurrentPreparedReadWrites
godbtests    | === RUN   TestMySQLDrivers/github.com/ziutek/mymysql:_Transaction
godbtests    | === RUN   TestMySQLDrivers/github.com/ziutek/mymysql:_Blobs
godbtests    | === RUN   TestMySQLDrivers/github.com/ziutek/mymysql:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestMySQLDrivers/github.com/ziutek/mymysql:_ConcurrentPreparedReadWrites
godbtests    | --- PASS: TestMySQLDrivers (4.10s)
godbtests    |  mysql_test.go:63: github.com/go-sql-driver/mysql revision: cd4cb909ce1a31435164be29bf3682031f61539a
godbtests    |     --- PASS: TestMySQLDrivers/github.com/go-sql-driver/mysql:_Transaction (0.02s)
godbtests    |     --- PASS: TestMySQLDrivers/github.com/go-sql-driver/mysql:_Blobs (0.01s)
godbtests    |     --- PASS: TestMySQLDrivers/github.com/go-sql-driver/mysql:_InsertOnceReadOneThousandTimes (2.52s)
godbtests    |     --- PASS: TestMySQLDrivers/github.com/go-sql-driver/mysql:_ConcurrentPreparedReadWrites (0.05s)
godbtests    |  mysql_test.go:69: github.com/ziutek/mymysql revision: 1d19cbf98d83564cc561192ae7d7183d795f7ac7
godbtests    |     --- PASS: TestMySQLDrivers/github.com/ziutek/mymysql:_Transaction (0.02s)
godbtests    |     --- PASS: TestMySQLDrivers/github.com/ziutek/mymysql:_Blobs (0.01s)
godbtests    |     --- PASS: TestMySQLDrivers/github.com/ziutek/mymysql:_InsertOnceReadOneThousandTimes (1.42s)
godbtests    |     --- PASS: TestMySQLDrivers/github.com/ziutek/mymysql:_ConcurrentPreparedReadWrites (0.05s)
godbtests    | === RUN   TestPostgresDrivers
godbtests    | === RUN   TestPostgresDrivers/github.com/jackc/pgx:_Transaction
godbtests    | === RUN   TestPostgresDrivers/github.com/jackc/pgx:_Blobs
godbtests    | === RUN   TestPostgresDrivers/github.com/jackc/pgx:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestPostgresDrivers/github.com/jackc/pgx:_ConcurrentPreparedReadWrites
godbtests    | === RUN   TestPostgresDrivers/github.com/lib/pq:_Transaction
godbtests    | === RUN   TestPostgresDrivers/github.com/lib/pq:_Blobs
godbtests    | === RUN   TestPostgresDrivers/github.com/lib/pq:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestPostgresDrivers/github.com/lib/pq:_ConcurrentPreparedReadWrites
godbtests    | === RUN   TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_Transaction
godbtests    | === RUN   TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_Blobs
godbtests    | === RUN   TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_ConcurrentPreparedReadWrites
godbtests    | --- PASS: TestPostgresDrivers (7.67s)
godbtests    |  postgres_test.go:66: github.com/jackc/pgx revision: 152dbffa4ac70ffd3c4efd1b7160d88ae8c21250
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jackc/pgx:_Transaction (0.02s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jackc/pgx:_Blobs (0.03s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jackc/pgx:_InsertOnceReadOneThousandTimes (2.45s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jackc/pgx:_ConcurrentPreparedReadWrites (0.09s)
godbtests    |  postgres_test.go:72: github.com/lib/pq revision: 8c6ee72f3e6bcb1542298dd5f76cb74af9742cec
godbtests    |     --- PASS: TestPostgresDrivers/github.com/lib/pq:_Transaction (0.01s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/lib/pq:_Blobs (0.02s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/lib/pq:_InsertOnceReadOneThousandTimes (2.23s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/lib/pq:_ConcurrentPreparedReadWrites (0.06s)
godbtests    |  postgres_test.go:78: github.com/jbarham/gopgsqldriver revision: f8287ee9bfe224aa4a7edcd73815ecbe69db7f68
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_Transaction (0.02s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_Blobs (0.02s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_InsertOnceReadOneThousandTimes (2.65s)
godbtests    |     --- PASS: TestPostgresDrivers/github.com/jbarham/gopgsqldriver:_ConcurrentPreparedReadWrites (0.07s)
godbtests    | === RUN   TestSqliteDrivers
godbtests    | === RUN   TestSqliteDrivers/github.com/gwenn/gosqlite:_Transaction
godbtests    | === RUN   TestSqliteDrivers/github.com/gwenn/gosqlite:_Blobs
godbtests    | === RUN   TestSqliteDrivers/github.com/gwenn/gosqlite:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestSqliteDrivers/github.com/gwenn/gosqlite:_ConcurrentPreparedReadWrites
godbtests    | === RUN   TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_Transaction
godbtests    | === RUN   TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_Blobs
godbtests    | === RUN   TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_InsertOnceReadOneThousandTimes
godbtests    | === RUN   TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_ConcurrentPreparedReadWrites
godbtests    | --- PASS: TestSqliteDrivers (1.09s)
godbtests    |  sqlite_test.go:49: github.com/gwenn/gosqlite revision: dd5964ffecc22120fa56ee8845e513e2a1cd08bb
godbtests    |     --- PASS: TestSqliteDrivers/github.com/gwenn/gosqlite:_Transaction (0.01s)
godbtests    |     --- PASS: TestSqliteDrivers/github.com/gwenn/gosqlite:_Blobs (0.01s)
godbtests    |     --- PASS: TestSqliteDrivers/github.com/gwenn/gosqlite:_InsertOnceReadOneThousandTimes (0.14s)
godbtests    |     --- PASS: TestSqliteDrivers/github.com/gwenn/gosqlite:_ConcurrentPreparedReadWrites (0.41s)
godbtests    |  sqlite_test.go:61: github.com/mxk/go-sqlite/sqlite3 revision: 167da9432e1f4602e95ea67b67051cfa34412e3f
godbtests    |     --- PASS: TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_Transacion (0.00s)
godbtests    |     --- PASS: TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_Blobs (0.00s)
godbtests    |     --- PASS: TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_InsertOnceReadOneThousandTimes (0.22s)
godbtests    |     --- PASS: TestSqliteDrivers/github.com/mxk/go-sqlite/sqlite3:_ConcurrentPreparedReadWrites (0.30s)
godbtests    | FAIL
godbtests    | exit status 1
godbtests    | FAIL     app/sqltest     23.651s
```

[1] Image copied from https://github.com/egonelbre/gophers
