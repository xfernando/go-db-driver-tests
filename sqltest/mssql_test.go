package sqltest

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/minus5/gofreetds"
)

var (
	gofreetds = &mssqlDB{
		driver:           "mssql",
		driverPkg:        "github.com/minus5/gofreetds",
		connectionString: "Server=mssql;Database=msdb;User Id=sa;Password=Gosqldbr00t!"}
	mssqldb = &mssqlDB{
		driver:           "denisenkom_mssql",
		driverPkg:        "github.com/denisenkom/go-mssqldb",
		connectionString: "Server=mssql;Database=msdb;User Id=sa;Password=Gosqldbr00t!"}
)

type mssqlDB database

func (m *mssqlDB) SQLBlobParam(size int) string {
	return fmt.Sprintf("VARBINARY(%d)", size)
}

func (m *mssqlDB) q(sql string) string {
	return sql
}

func (m *mssqlDB) DB() *sql.DB {
	return m.db
}

func (m *mssqlDB) T() *testing.T {
	return m.t
}

func (m *mssqlDB) RunTest(t *testing.T, fn func(Tester)) {
	m.t = t
	db, err := sql.Open(m.driver, m.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	m.db = db
	t.Logf("Driver type: %T", db.Driver())

	// Drop all tables in the test database
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_name LIKE '" +
		TablePrefix + "%'")
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

func TestMsSQLDrivers(t *testing.T) {
	for i := 0; i < 3; i++ {
		if !Running("mssql", 1433) {
			<-time.After(60 * time.Second)
		}
	}
	if !Running("mssql", 1433) {
		t.Fatalf("skipping tests; Microsoft SQL Server not responding on mssql:1433 after 3 tries")
		return
	}

	t.Logf("%s revision: %s", gofreetds.driverPkg, gitRevision(t, gofreetds.driverPkg))
	t.Run("gofreetds: TXQuery", testGofreetdsTxQuery)
	t.Run("gofreetds: Blobs", testGofreetdsBlobs)
	t.Run("gofreetds: ManyQueryRow", testGofreetdsManyQueryRow)
	t.Run("gofreetds: PreparedStmt", testGofreetdsPreparedStmt)

	t.Logf("%s revision: %s", mssqldb.driverPkg, gitRevision(t, mssqldb.driverPkg))
	t.Run("mssqldb: TXQuery", testMssqldbTxQuery)
	t.Run("mssqldb: Blobs", testMssqldbBlobs)
	t.Run("mssqldb: ManyQueryRow", testMssqldbManyQueryRow)
	t.Run("mssqldb: PreparedStmt", testMssqldbPreparedStmt)
}

func testGofreetdsTxQuery(t *testing.T)      { gofreetds.RunTest(t, testTxQuery) }
func testGofreetdsBlobs(t *testing.T)        { gofreetds.RunTest(t, testBlobs) }
func testGofreetdsManyQueryRow(t *testing.T) { gofreetds.RunTest(t, testManyQueryRow) }
func testGofreetdsPreparedStmt(t *testing.T) { gofreetds.RunTest(t, testPreparedStmt) }

func testMssqldbTxQuery(t *testing.T)      { mssqldb.RunTest(t, testTxQuery) }
func testMssqldbBlobs(t *testing.T)        { mssqldb.RunTest(t, testBlobs) }
func testMssqldbManyQueryRow(t *testing.T) { mssqldb.RunTest(t, testManyQueryRow) }
func testMssqldbPreparedStmt(t *testing.T) { mssqldb.RunTest(t, testPreparedStmt) }
