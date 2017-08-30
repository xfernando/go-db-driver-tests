package sqltest

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

var (
	// github.com/lib/pq
	pq = &postgresDB{driver: "postgres", connectionString: "user=postgres password=root host=postgres dbname=gosqltest sslmode=disable"}
	// github.com/jbarham/gopgsqldriver
	// not going to test this now, it registers as postgres, conflicting with the previous driver
	gopgsql = &postgresDB{driver: "postgres", connectionString: "user=postgres password=root dbname=gosqltest sslmode=disable"}
	// github.com/jackc/pgx
	pgx = &postgresDB{driver: "pgx", connectionString: "user=postgres password=root host=postgres port=5432 database=gosqltest sslmode=disable"}
)

type postgresDB database

func (p *postgresDB) DB() *sql.DB {
	return p.db
}

func (p *postgresDB) T() *testing.T {
	return p.t
}

func (p *postgresDB) RunTest(t *testing.T, fn func(Tester)) {
	p.t = t
	db, err := sql.Open(p.driver, p.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	p.db = db

	// Drop all tables in the test database.
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_name LIKE '" +
		TablePrefix + "%' AND table_schema = 'public'")
	if err != nil {
		t.Fatalf("failed to enumerate tables: %v", err)
	}
	for rows.Next() {
		var table string
		if rows.Scan(&table) == nil {
			mustExec(p, "DROP TABLE "+table)
		}
	}

	fn(p)
}

func (p *postgresDB) SQLBlobParam(size int) string {
	return "bytea"
}

func (p *postgresDB) q(sql string) string {
	n := 0
	return qrx.ReplaceAllStringFunc(sql, func(string) string {
		n++
		return "$" + strconv.Itoa(n)
	})
}

func TestPostgresDrivers(t *testing.T) {
	for i := 0; i < 3; i++ {
		if !Running("postgres", 5432) {
			t.Logf("Postgres not running, waiting 60 seconds, try %d", i)
			<-time.After(60 * time.Second)
		}
	}
	if !Running("postgres", 5432) {
		t.Fatalf("skipping tests; Postgres not responding on postgres:5432 after 3 tries")
		return
	}

	t.Run("pgx: TXQuery", testPGXTxQuery)
	t.Run("pgx: Blobs", testPGXBlobs)
	t.Run("pgx: ManyQueryRow", testPGXPreparedStmt)
	t.Run("pgx: PreparedStmt", testPGXManyQueryRow)

	t.Run("pq: TXQuery", testPQTxQuery)
	t.Run("pq: Blobs", testPQBlobs)
	t.Run("pq: ManyQueryRow", testPQManyQueryRow)
	t.Run("pq: PreparedStmt", testPQPreparedStmt)
}

func testPQTxQuery(t *testing.T)      { pq.RunTest(t, testTxQuery) }
func testPQBlobs(t *testing.T)        { pq.RunTest(t, testBlobs) }
func testPQManyQueryRow(t *testing.T) { pq.RunTest(t, testManyQueryRow) }
func testPQPreparedStmt(t *testing.T) { pq.RunTest(t, testPreparedStmt) }

func testPGXTxQuery(t *testing.T)      { pgx.RunTest(t, testTxQuery) }
func testPGXBlobs(t *testing.T)        { pgx.RunTest(t, testBlobs) }
func testPGXPreparedStmt(t *testing.T) { pgx.RunTest(t, testPreparedStmt) }
func testPGXManyQueryRow(t *testing.T) { pgx.RunTest(t, testManyQueryRow) }
