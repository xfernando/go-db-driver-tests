package sqltest

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/gwenn/gosqlite"
	// _ "github.com/mattn/go-sqlite3"
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

func (s *sqliteDB) RunTest(t *testing.T, fn func(t *testing.T, db *sql.DB)) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	db, err := sql.Open(s.driver, filepath.Join(tempDir, "foo.db"))
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	fn(t, db)
}

func TestSqliteDrivers(t *testing.T) {
	t.Logf("github.com/gwenn/gosqlite revision: %s", gitRevision(t, gwennGosqlite.driverPkg))
	t.Run("github.com/gwenn/gosqlite: Transaction", testGwennGoSqliteTransaction)
	t.Run("github.com/gwenn/gosqlite: Blobs", testGwennGoSqliteBlobs)
	t.Run("github.com/gwenn/gosqlite: InsertOnceReadOneThousandTimes", testGwennGoSqliteInsertOnceReadOneThousandTimes)
	t.Run("github.com/gwenn/gosqlite: ConcurrentPreparedReadWrites", testGwennGoSqliteConcurrentPreparedReadWrites)

	// t.Logf("github.com/mattn/go-sqlite3 revision: %s", gitRevision(t, mattnGosqlite3.driverPkg))
	// t.Run("github.com/mattn/go-sqlite3: TXQuery", testMattnGoSqliteTransaction)
	// t.Run("github.com/mattn/go-sqlite3: Blobs", testMattnGoSqliteBlobs)
	// t.Run("github.com/mattn/go-sqlite3: InsertOnceReadOneThousandTimes", testMattnGoSqliteInsertOnceReadOneThousandTimes)
	// t.Run("github.com/mattn/go-sqlite3: ConcurrentPreparedReadWrites", testMattnGoSqliteConcurrentPreparedReadWrites)

	t.Logf("github.com/mxk/go-sqlite/sqlite3 revision: %s", gitRevision(t, mkxSqlite3.driverPkg))
	t.Run("github.com/mxk/go-sqlite/sqlite3: Transaction", testMkxSqliteTransaction)
	t.Run("github.com/mxk/go-sqlite/sqlite3: Blobs", testMkxSqliteBlobs)
	t.Run("github.com/mxk/go-sqlite/sqlite3: InsertOnceReadOneThousandTimes", testMkxSqliteInsertOnceReadOneThousandTimes)
	t.Run("github.com/mxk/go-sqlite/sqlite3: ConcurrentPreparedReadWrites", testMkxSqliteConcurrentPreparedReadWrites)
}

func testGwennGoSqliteTransaction(t *testing.T) { gwennGosqlite.RunTest(t, testSqliteTransaction) }
func testGwennGoSqliteBlobs(t *testing.T)       { gwennGosqlite.RunTest(t, testSqliteBlobs) }
func testGwennGoSqliteInsertOnceReadOneThousandTimes(t *testing.T) {
	gwennGosqlite.RunTest(t, testSqliteInsertOnceReadOneThousandTimes)
}
func testGwennGoSqliteConcurrentPreparedReadWrites(t *testing.T) {
	gwennGosqlite.RunTest(t, testSqliteConcurrentPreparedReadWrites)
}

func testMattnGoSqliteTransaction(t *testing.T) { mattnGosqlite3.RunTest(t, testSqliteTransaction) }
func testMattnGoSqliteBlobs(t *testing.T)       { mattnGosqlite3.RunTest(t, testSqliteBlobs) }
func testMattnGoSqliteInsertOnceReadOneThousandTimes(t *testing.T) {
	mattnGosqlite3.RunTest(t, testSqliteInsertOnceReadOneThousandTimes)
}
func testMattnGoSqliteConcurrentPreparedReadWrites(t *testing.T) {
	mattnGosqlite3.RunTest(t, testSqliteConcurrentPreparedReadWrites)
}

func testMkxSqliteTransaction(t *testing.T) { mkxSqlite3.RunTest(t, testSqliteTransaction) }
func testMkxSqliteBlobs(t *testing.T)       { mkxSqlite3.RunTest(t, testSqliteBlobs) }
func testMkxSqliteInsertOnceReadOneThousandTimes(t *testing.T) {
	mkxSqlite3.RunTest(t, testSqliteInsertOnceReadOneThousandTimes)
}
func testMkxSqliteConcurrentPreparedReadWrites(t *testing.T) {
	mkxSqlite3.RunTest(t, testSqliteConcurrentPreparedReadWrites)
}

func testSqliteTransaction(t *testing.T, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	_, err = db.Exec("create table gosqltest_foo (id integer primary key, name varchar(50))")
	if err != nil {
		t.Logf("cannot drop table gosqltest_foo: %s", err)
	}

	_, err = tx.Exec("insert into gosqltest_foo (id, name) values(?,?)", 1, "bob")
	if err != nil {
		t.Fatal(err)
	}

	r, err := tx.Query("select name from gosqltest_foo where id = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	if !r.Next() {
		if r.Err() != nil {
			t.Fatal(err)
		}
		t.Fatal("expected one rows")
	}

	var name string
	err = r.Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
}

func testSqliteBlobs(t *testing.T, db *sql.DB) {
	var blob = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	mustExec(db, t, fmt.Sprintf("create table gosqltest_foo (id integer primary key, bar blob[%d])", len(blob)))
	mustExec(db, t, "insert into gosqltest_foo (id, bar) values(?,?)", 0, blob)

	want := fmt.Sprintf("%x", blob)

	b := make([]byte, len(blob))
	err := db.QueryRow("select bar from gosqltest_foo where id = ?", 0).Scan(&b)
	got := fmt.Sprintf("%x", b)
	if err != nil {
		t.Errorf("[]byte scan: %v", err)
	} else if got != want {
		t.Errorf("for []byte, got %q; want %q", got, want)
	}

	err = db.QueryRow("select bar from gosqltest_foo where id = ?", 0).Scan(&got)
	want = string(blob)
	if err != nil {
		t.Errorf("string scan: %v", err)
	} else if got != want {
		t.Errorf("for string, got %q; want %q", got, want)
	}
}

func testSqliteInsertOnceReadOneThousandTimes(t *testing.T, db *sql.DB) {
	if testing.Short() {
		t.Logf("skipping in short mode")
		return
	}
	mustExec(db, t, "create table gosqltest_foo (id integer primary key, name varchar(50))")
	mustExec(db, t, "insert into gosqltest_foo (id, name) values(?,?)", 1, "bob")
	var name string
	for i := 0; i < 10000; i++ {
		err := db.QueryRow("select name from gosqltest_foo where id = ?", 1).Scan(&name)
		if err != nil || name != "bob" {
			t.Fatalf("on query %d: err=%v, name=%q", i, err, name)
		}
	}
}

func testSqliteConcurrentPreparedReadWrites(t *testing.T, db *sql.DB) {
	mustExec(db, t, "CREATE TABLE gosqltest_t (count INT)")
	sel, err := db.Prepare("SELECT count FROM gosqltest_t ORDER BY count DESC")
	if err != nil {
		t.Fatalf("prepare 1: %v", err)
	}
	ins, err := db.Prepare("INSERT INTO gosqltest_t (count) VALUES (?)")
	if err != nil {
		t.Fatalf("prepare 2: %v", err)
	}

	for n := 1; n <= 3; n++ {
		if _, err := ins.Exec(n); err != nil {
			t.Fatalf("insert(%d) = %v", n, err)
		}
	}

	const nRuns = 10
	ch := make(chan bool)
	for i := 0; i < nRuns; i++ {
		go func() {
			defer func() {
				ch <- true
			}()
			for j := 0; j < 10; j++ {
				count := 0
				if err := sel.QueryRow().Scan(&count); err != nil && err != sql.ErrNoRows {
					t.Errorf("Query: %v", err)
					return
				}
				if _, err := ins.Exec(rand.Intn(100)); err != nil {
					t.Errorf("Insert: %v", err)
					return
				}
			}
		}()
	}
	for i := 0; i < nRuns; i++ {
		<-ch
	}
}
