package sqltest

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

var (
	// github.com/go-sql-driver/mysq
	goMysqlDB = &mysqlDB{
		driver:           "mysql",
		driverPkg:        "github.com/go-sql-driver/mysql",
		connectionString: "root:root@tcp(mysql:3306)/gosqltest"}
	// github.com/ziutek/mymysql/godrv
	myMysqlDB = &mysqlDB{
		driver:           "mymysql",
		driverPkg:        "github.com/ziutek/mymysql",
		connectionString: "tcp:mysql:3306*gosqltest/root/root"}
)

type mysqlDB database

func (m *mysqlDB) SQLBlobParam(size int) string {
	return fmt.Sprintf("VARBINARY(%d)", size)
}

func (m *mysqlDB) q(sql string) string {
	return sql
}

func (m *mysqlDB) DB() *sql.DB {
	return m.db
}

func (m *mysqlDB) T() *testing.T {
	return m.t
}

func (m *mysqlDB) RunTest(t *testing.T, fn func(Tester)) {
	m.t = t
	db, err := sql.Open(m.driver, m.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	m.db = db

	// Drop all tables in the test database
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		t.Fatalf("failed to enumerate tables: %v", err)
	}
	for rows.Next() {
		var table string
		if rows.Scan(&table) == nil &&
			strings.HasPrefix(strings.ToLower(table), strings.ToLower(TablePrefix)) {
			mustExec(m, "DROP TABLE "+table)
		}
	}

	fn(m)
}

func TestMySQLDrivers(t *testing.T) {
	for i := 0; i < 3; i++ {
		if !Running("mysql", 3306) {
			//t.Logf("Try %d: MySQL not running. Waiting 60 secs for container initialization.", i)
			<-time.After(60 * time.Second)
		}
	}
	if !Running("mysql", 3306) {
		t.Fatalf("skipping tests; MySQL not responding on mysql:3306 after 3 tries")
		return
	}

	t.Logf("%s revision: %s", goMysqlDB.driverPkg, gitRevision(t, goMysqlDB.driverPkg))
	t.Run("gomysql: TXQuery", testGoMySQLTxQuery)
	t.Run("gomysql: Blobs", testGoMySQLBlobs)
	t.Run("gomysql: ManyQueryRow", testGoMySQLManyQueryRow)
	t.Run("gomysql: PreparedStmt", testGoMySQLPreparedStmt)

	t.Logf("%s revision: %s", myMysqlDB.driverPkg, gitRevision(t, myMysqlDB.driverPkg))
	t.Run("mymysql: TXQuery", testMyMySQLTxQuery)
	t.Run("mymysql: Blobs", testMyMySQLBlobs)
	t.Run("mymysql: ManyQueryRow", testMyMySQLManyQueryRow)
	t.Run("mymysql: PreparedStmt", testMyMySQLPreparedStmt)
}

func testGoMySQLTxQuery(t *testing.T)      { goMysqlDB.RunTest(t, testTxQuery) }
func testGoMySQLBlobs(t *testing.T)        { goMysqlDB.RunTest(t, testBlobs) }
func testGoMySQLManyQueryRow(t *testing.T) { goMysqlDB.RunTest(t, testManyQueryRow) }
func testGoMySQLPreparedStmt(t *testing.T) { goMysqlDB.RunTest(t, testPreparedStmt) }

func testMyMySQLTxQuery(t *testing.T)      { myMysqlDB.RunTest(t, testTxQuery) }
func testMyMySQLBlobs(t *testing.T)        { myMysqlDB.RunTest(t, testBlobs) }
func testMyMySQLManyQueryRow(t *testing.T) { myMysqlDB.RunTest(t, testManyQueryRow) }
func testMyMySQLPreparedStmt(t *testing.T) { myMysqlDB.RunTest(t, testPreparedStmt) }
