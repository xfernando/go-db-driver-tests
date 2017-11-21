package sqltest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	_ "github.com/nakagami/firebirdsql"
)

var (
	nakgamiFirebirdDB = &firebirdDB{
		driver:           "firebirdsql",
		driverPkg:        "github.com/nakagami/firebirdsql",
		connectionString: "root:root@firebird/gosqltest"}
)

type firebirdDB database

func (f *firebirdDB) RunTest(t *testing.T, fn func(t *testing.T, db *sql.DB)) {
	db, err := sql.Open(f.driver, f.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer db.Close()

	// Drop all tables in the test database
	rows, err := db.Query("SELECT RDB$RELATION_NAME FROM RDB$RELATIONS WHERE (RDB$SYSTEM_FLAG <> 1 OR RDB$SYSTEM_FLAG IS NULL) AND RDB$VIEW_BLR IS NULL ORDER BY RDB$RELATION_NAME")
	if err != nil {
		t.Fatalf("failed to enumerate tables: %v", err)
	}
	for rows.Next() {
		var table string
		if rows.Scan(&table) == nil &&
			strings.HasPrefix(strings.ToLower(table), strings.ToLower(TablePrefix)) {
			mustExec(db, t, "DROP TABLE "+table)
		}
	}

	fn(t, db)
}

func TestFirebirdDrivers(t *testing.T) {
	for i := 0; i < 3; i++ {
		if !Running("firebird", 3050) {
			<-time.After(60 * time.Second)
		} else {
			break
		}
	}
	if !Running("firebird", 3050) {
		t.Fatalf("skipping tests; Firebird not responding on firebid:3050 after 3 tries")
		return
	}

	t.Logf("github.com/nakagami/firebirdsql revision: %s", gitRevision(t, nakgamiFirebirdDB.driverPkg))
	t.Run("github.com/nakagami/firebirdsql - Transaction", testFirebirdSQLTransaction)
	t.Run("github.com/nakagami/firebirdsql - Blobs", testFirebirdSQLBlobs)
	t.Run("github.com/nakagami/firebirdsql - InsertOnceReadOneThousandTimes", testFirebirdSQLInsertOnceReadOneThousandTimes)
	t.Run("github.com/nakagami/firebirdsql - ConcurrentPreparedReadWrites", testFirebirdSQLConcurrentPreparedReadWrites)
}

func testFirebirdSQLTransaction(t *testing.T) { nakgamiFirebirdDB.RunTest(t, testFirebirdTransaction) }
func testFirebirdSQLBlobs(t *testing.T)       { nakgamiFirebirdDB.RunTest(t, testFirebirdBlobs) }
func testFirebirdSQLInsertOnceReadOneThousandTimes(t *testing.T) {
	nakgamiFirebirdDB.RunTest(t, testFirebirdInsertOnceReadOneThousandTimes)
}
func testFirebirdSQLConcurrentPreparedReadWrites(t *testing.T) {
	nakgamiFirebirdDB.RunTest(t, testFirebirdConcurrentPreparedReadWrites)
}

func testFirebirdTransaction(t *testing.T, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	_, err = db.Exec("create table " + TablePrefix + "foo (id int primary key, name varchar(50))")
	if err != nil {
		t.Fatalf("cannot create table gosqltest_foo: %s", err)
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

func testFirebirdBlobs(t *testing.T, db *sql.DB) {
	var blob = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	mustExec(db, t, "create table gosqltest_foo (id int primary key, bar blob sub_type binary)")
	mustExec(db, t, "insert into gosqltest_foo (id, bar) values(?,?)", 0, blob)

	want := fmt.Sprintf("%x", blob)

	b := make([]byte, 16)
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

func testFirebirdInsertOnceReadOneThousandTimes(t *testing.T, db *sql.DB) {
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

func testFirebirdConcurrentPreparedReadWrites(t *testing.T, db *sql.DB) {
	mustExec(db, t, "CREATE TABLE gosqltest_t (count_f int)")
	sel, err := db.Prepare("SELECT count_f FROM " + TablePrefix + "t ORDER BY count_f DESC")
	if err != nil {
		t.Fatalf("prepare 1: %v", err)
	}
	ins, err := db.Prepare("INSERT INTO gosqltest_t (count_f) VALUES (?)")
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
