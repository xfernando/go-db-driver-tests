package sqltest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jbarham/gopgsqldriver"
	_ "github.com/lib/pq"
)

var (
	pq = &postgresDB{
		driver:           "postgres",
		driverPkg:        "github.com/lib/pq",
		connectionString: "user=postgres password=root host=postgres dbname=gosqltest sslmode=disable"}
	// we rename the driver during docker image construction, otherwise it conflicts with the previous one
	gopg = &postgresDB{driver: "gopgsql",
		driverPkg:        "github.com/jbarham/gopgsqldriver",
		connectionString: "user=postgres password=root host=postgres dbname=gosqltest sslmode=disable"}
	pgx = &postgresDB{
		driver:           "pgx",
		driverPkg:        "github.com/jackc/pgx",
		connectionString: "user=postgres password=root host=postgres dbname=gosqltest sslmode=disable"}
)

type postgresDB database

func (p *postgresDB) RunTest(t *testing.T, fn func(t *testing.T, db *sql.DB)) {
	db, err := sql.Open(p.driver, p.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	defer db.Close()

	// Drop all tables in the test database.
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_name LIKE '" +
		TablePrefix + "%' AND table_schema = 'public'")
	if err != nil {
		t.Fatalf("failed to enumerate tables: %v", err)
	}
	for rows.Next() {
		var table string
		if rows.Scan(&table) == nil {
			mustExec(db, t, "DROP TABLE "+table)
		}
	}

	fn(t, db)
}

func TestPostgresDrivers(t *testing.T) {
	for i := 0; i < 3; i++ {
		if !Running("postgres", 5432) {
			<-time.After(60 * time.Second)
		}
	}
	if !Running("postgres", 5432) {
		t.Fatalf("skipping tests; Postgres not responding on postgres:5432 after 3 tries")
		return
	}

	t.Logf("github.com/jackc/pgx revision: %s", gitRevision(t, pgx.driverPkg))
	t.Run("github.com/jackc/pgx: Transaction", testPGXTransaction)
	t.Run("github.com/jackc/pgx: Blobs", testPGXBlobs)
	t.Run("github.com/jackc/pgx: InsertOnceReadOneThousandTimes", testPGXInsertOnceReadOneThousandTimes)
	t.Run("github.com/jackc/pgx: ConcurrentPreparedReadWrites", testPGXConcurrentPreparedReadWrites)

	t.Logf("github.com/lib/pq revision: %s", gitRevision(t, pq.driverPkg))
	t.Run("github.com/lib/pq: Transaction", testPQTransaction)
	t.Run("github.com/lib/pq: Blobs", testPQBlobs)
	t.Run("github.com/lib/pq: InsertOnceReadOneThousandTimes", testPQInsertOnceReadOneThousandTimes)
	t.Run("github.com/lib/pq: ConcurrentPreparedReadWrites", testPQConcurrentPreparedReadWrites)

	t.Logf("github.com/jbarham/gopgsqldriver revision: %s", gitRevision(t, gopg.driverPkg))
	t.Run("github.com/jbarham/gopgsqldriver: Transaction", testGoPGTransaction)
	t.Run("github.com/jbarham/gopgsqldriver: Blobs", testGoPGBlobs)
	t.Run("github.com/jbarham/gopgsqldriver: InsertOnceReadOneThousandTimes", testGoPGInsertOnceReadOneThousandTimes)
	t.Run("github.com/jbarham/gopgsqldriver: ConcurrentPreparedReadWrites", testGoPGConcurrentPreparedReadWrites)
}

func testPQTransaction(t *testing.T) { pq.RunTest(t, testPostgresTransaction) }
func testPQBlobs(t *testing.T)       { pq.RunTest(t, testPostgresBlobs) }
func testPQInsertOnceReadOneThousandTimes(t *testing.T) {
	pq.RunTest(t, testPostgresInsertOnceReadOneThousandTimes)
}
func testPQConcurrentPreparedReadWrites(t *testing.T) {
	pq.RunTest(t, testPostgresConcurrentPreparedReadWrites)
}

func testPGXTransaction(t *testing.T) { pgx.RunTest(t, testPostgresTransaction) }
func testPGXBlobs(t *testing.T)       { pgx.RunTest(t, testPostgresBlobs) }
func testPGXInsertOnceReadOneThousandTimes(t *testing.T) {
	pgx.RunTest(t, testPostgresInsertOnceReadOneThousandTimes)
}
func testPGXConcurrentPreparedReadWrites(t *testing.T) {
	pgx.RunTest(t, testPostgresConcurrentPreparedReadWrites)
}

func testGoPGTransaction(t *testing.T) { gopg.RunTest(t, testPostgresTransaction) }
func testGoPGBlobs(t *testing.T)       { gopg.RunTest(t, testPostgresBlobs) }
func testGoPGInsertOnceReadOneThousandTimes(t *testing.T) {
	gopg.RunTest(t, testPostgresInsertOnceReadOneThousandTimes)
}
func testGoPGConcurrentPreparedReadWrites(t *testing.T) {
	gopg.RunTest(t, testPostgresConcurrentPreparedReadWrites)
}

func testPostgresTransaction(t *testing.T, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	_, err = db.Exec("create table gosqltest_foo (id integer primary key, name varchar(50))")
	if err != nil {
		t.Logf("cannot drop table gosqltest_foo: %s", err)
	}

	_, err = tx.Exec("insert into gosqltest_foo (id, name) values($1,$2)", 1, "bob")
	if err != nil {
		t.Fatal(err)
	}

	r, err := tx.Query("select name from gosqltest_foo where id = $1", 1)
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

func testPostgresBlobs(t *testing.T, db *sql.DB) {
	var blob = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	mustExec(db, t, "create table gosqltest_foo (id integer primary key, bar bytea)")
	mustExec(db, t, "insert into gosqltest_foo (id, bar) values($1,$2)", 0, blob)

	want := fmt.Sprintf("%x", blob)

	b := make([]byte, len(blob))
	err := db.QueryRow("select bar from gosqltest_foo where id = $1", 0).Scan(&b)
	got := fmt.Sprintf("%x", b)
	if err != nil {
		t.Errorf("[]byte scan: %v", err)
	} else if got != want {
		t.Errorf("for []byte, got %q; want %q", got, want)
	}

	err = db.QueryRow("select bar from gosqltest_foo where id = $1", 0).Scan(&got)
	want = string(blob)
	if err != nil {
		t.Errorf("string scan: %v", err)
	} else if got != want {
		t.Errorf("for string, got %q; want %q", got, want)
	}
}

func testPostgresInsertOnceReadOneThousandTimes(t *testing.T, db *sql.DB) {
	if testing.Short() {
		t.Logf("skipping in short mode")
		return
	}
	mustExec(db, t, "create table gosqltest_foo (id integer primary key, name varchar(50))")
	mustExec(db, t, "insert into gosqltest_foo (id, name) values($1,$2)", 1, "bob")
	var name string
	for i := 0; i < 10000; i++ {
		err := db.QueryRow("select name from gosqltest_foo where id = $1", 1).Scan(&name)
		if err != nil || name != "bob" {
			t.Fatalf("on query %d: err=%v, name=%q", i, err, name)
		}
	}
}

func testPostgresConcurrentPreparedReadWrites(t *testing.T, db *sql.DB) {
	mustExec(db, t, "CREATE TABLE gosqltest_t (count INT)")
	sel, err := db.Prepare("SELECT count FROM gosqltest_t ORDER BY count DESC")
	if err != nil {
		t.Fatalf("prepare 1: %v", err)
	}
	ins, err := db.Prepare("INSERT INTO gosqltest_t (count) VALUES ($1)")
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
