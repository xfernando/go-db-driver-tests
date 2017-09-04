package sqltest

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/gwenn/gosqlite"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/mxk/go-sqlite/sqlite3"
)

var (
	gwennGosqlite = &sqliteDB{
		driver:    "gwenn_sqlite3",
		driverPkg: "github.com/gwenn/gosqlite"}
	mattnGosqlite3 = &sqliteDB{
		driver:    "mattn_sqlite3",
		driverPkg: "github.com/mattn/go-sqlite3"}
	mkxSqlite3 = &sqliteDB{
		driver:    "mxk_sqlite3",
		driverPkg: "github.com/mxk/go-sqlite"}
)

type sqliteDB database

func (s *sqliteDB) SQLBlobParam(size int) string {
	return fmt.Sprintf("blob[%d]", size)
}

func (s *sqliteDB) q(sql string) string {
	return sql
}

func (s *sqliteDB) DB() *sql.DB {
	return s.db
}

func (s *sqliteDB) T() *testing.T {
	return s.t
}

func (s *sqliteDB) RunTest(t *testing.T, fn func(Tester)) {
	s.t = t
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	db, err := sql.Open(s.driver, filepath.Join(tempDir, "foo.db"))
	s.db = db
	t.Logf("Driver type: %T", db.Driver())

	if err != nil {
		t.Fatalf("foo.db open fail: %v", err)
	}
	fn(s)
}

func TestSqliteDrivers(t *testing.T) {
	t.Logf("%s revision: %s", gwennGosqlite.driverPkg, gitRevision(t, gwennGosqlite.driverPkg))
	t.Run("gwenn_sqlite3: TXQuery", testGwennGoSqliteSQLTxQuery)
	t.Run("gwenn_sqlite3: Blobs", testGwennGoSqliteBlobs)
	t.Run("gwenn_sqlite3: ManyQueryRow", testGwennGoSqliteManyQueryRow)
	t.Run("gwenn_sqlite3: PreparedStmt", testGwennGoSqlitePreparedStmt)

	// t.Logf("%s revision: %s", mattnGosqlite3.driverPkg, gitRevision(t, mattnGosqlite3.driverPkg))
	// t.Run("mattn_gosqlite3: TXQuery", testMattnGoSqliteTxQuery)
	// t.Run("mattn_gosqlite3: Blobs", testMattnGoSqliteBlobs)
	// t.Run("mattn_gosqlite3: ManyQueryRow", testMattnGoSqliteManyQueryRow)
	// t.Run("mattn_gosqlite3: PreparedStmt", testMattnGoSqlitePreparedStmt)

	t.Logf("%s revision: %s", mkxSqlite3.driverPkg, gitRevision(t, mkxSqlite3.driverPkg))
	t.Run("mkx_sqlite3: TXQuery", testMkxSqliteSQLTxQuery)
	t.Run("mkx_sqlite3: Blobs", testMkxSqliteSQLBlobs)
	t.Run("mkx_sqlite3: ManyQueryRow", testMkxSqliteSQLManyQueryRow)
	t.Run("mkx_sqlite3: PreparedStmt", testMkxSqliteSQLPreparedStmt)
}

func testGwennGoSqliteSQLTxQuery(t *testing.T)   { gwennGosqlite.RunTest(t, testTxQuery) }
func testGwennGoSqliteBlobs(t *testing.T)        { gwennGosqlite.RunTest(t, testBlobs) }
func testGwennGoSqliteManyQueryRow(t *testing.T) { gwennGosqlite.RunTest(t, testManyQueryRow) }
func testGwennGoSqlitePreparedStmt(t *testing.T) { gwennGosqlite.RunTest(t, testPreparedStmt) }

func testMattnGoSqliteTxQuery(t *testing.T)      { mattnGosqlite3.RunTest(t, testTxQuery) }
func testMattnGoSqliteBlobs(t *testing.T)        { mattnGosqlite3.RunTest(t, testBlobs) }
func testMattnGoSqliteManyQueryRow(t *testing.T) { mattnGosqlite3.RunTest(t, testManyQueryRow) }
func testMattnGoSqlitePreparedStmt(t *testing.T) { mattnGosqlite3.RunTest(t, testPreparedStmt) }

func testMkxSqliteSQLTxQuery(t *testing.T)      { mkxSqlite3.RunTest(t, testTxQuery) }
func testMkxSqliteSQLBlobs(t *testing.T)        { mkxSqlite3.RunTest(t, testBlobs) }
func testMkxSqliteSQLManyQueryRow(t *testing.T) { mkxSqlite3.RunTest(t, testManyQueryRow) }
func testMkxSqliteSQLPreparedStmt(t *testing.T) { mkxSqlite3.RunTest(t, testPreparedStmt) }
