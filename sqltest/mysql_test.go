package sqltest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

var (
	goMysqlDB = &mysqlDB{
		driver:           "mysql",
		driverPkg:        "github.com/go-sql-driver/mysql",
		connectionString: "root:root@tcp(mysql:3306)/gosqltest"}
	myMysqlDB = &mysqlDB{
		driver:           "mymysql",
		driverPkg:        "github.com/ziutek/mymysql",
		connectionString: "tcp:mysql:3306*gosqltest/root/root"}
)

type mysqlDB database

func (m *mysqlDB) RunTest(t *testing.T, fn func(t *testing.T, db *sql.DB)) {
	db, err := sql.Open(m.driver, m.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer db.Close()

	// Drop all tables in the test database
	rows, err := db.Query("SHOW TABLES")
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

func TestMySQLDrivers(t *testing.T) {
	for i := 0; i < 3; i++ {
		if !Running("mysql", 3306) {
			<-time.After(60 * time.Second)
		}
	}
	if !Running("mysql", 3306) {
		t.Fatalf("skipping tests; MySQL not responding on mysql:3306 after 3 tries")
		return
	}

	t.Logf("github.com/go-sql-driver/mysql revision: %s", gitRevision(t, goMysqlDB.driverPkg))
	t.Run("github.com/go-sql-driver/mysql: Transaction", testGoMySQLTransaction)
	t.Run("github.com/go-sql-driver/mysql: Blobs", testGoMySQLBlobs)
	t.Run("github.com/go-sql-driver/mysql: InsertOnceReadOneThousandTimes", testGoMySQLInsertOnceReadOneThousandTimes)
	t.Run("github.com/go-sql-driver/mysql: ConcurrentPreparedReadWrites", testGoMySQLConcurrentPreparedReadWrites)

	t.Logf("github.com/ziutek/mymysql revision: %s", gitRevision(t, myMysqlDB.driverPkg))
	t.Run("github.com/ziutek/mymysql: Transaction", testMyMySQLTransaction)
	t.Run("github.com/ziutek/mymysql: Blobs", testMyMySQLBlobs)
	t.Run("github.com/ziutek/mymysql: InsertOnceReadOneThousandTimes", testMyMySQLInsertOnceReadOneThousandTimes)
	t.Run("github.com/ziutek/mymysql: ConcurrentPreparedReadWrites", testMyMySQLConcurrentPreparedReadWrites)
}

func testGoMySQLTransaction(t *testing.T) { goMysqlDB.RunTest(t, testMysqlTransaction) }
func testGoMySQLBlobs(t *testing.T)       { goMysqlDB.RunTest(t, testMysqlBlobs) }
func testGoMySQLInsertOnceReadOneThousandTimes(t *testing.T) {
	goMysqlDB.RunTest(t, testMysqlInsertOnceReadOneThousandTimes)
}
func testGoMySQLConcurrentPreparedReadWrites(t *testing.T) {
	goMysqlDB.RunTest(t, testMysqlConcurrentPreparedReadWrites)
}

func testMyMySQLTransaction(t *testing.T) { myMysqlDB.RunTest(t, testMysqlTransaction) }
func testMyMySQLBlobs(t *testing.T)       { myMysqlDB.RunTest(t, testMysqlBlobs) }
func testMyMySQLInsertOnceReadOneThousandTimes(t *testing.T) {
	myMysqlDB.RunTest(t, testMysqlInsertOnceReadOneThousandTimes)
}
func testMyMySQLConcurrentPreparedReadWrites(t *testing.T) {
	myMysqlDB.RunTest(t, testMysqlConcurrentPreparedReadWrites)
}

func testMysqlTransaction(t *testing.T, db *sql.DB) {
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

func testMysqlBlobs(t *testing.T, db *sql.DB) {
	var blob = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	mustExec(db, t, fmt.Sprintf("create table gosqltest_foo (id integer primary key, bar VARBINARY(%d))", len(blob)))
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

func testMysqlInsertOnceReadOneThousandTimes(t *testing.T, db *sql.DB) {
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

func testMysqlConcurrentPreparedReadWrites(t *testing.T, db *sql.DB) {
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
