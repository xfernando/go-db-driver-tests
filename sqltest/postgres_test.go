package sqltest

import (
	"database/sql"
	"strconv"
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
		connectionString: "user=postgres password=root dbname=gosqltest sslmode=disable"}
	// github.com/jackc/pgx
	pgx = &postgresDB{
		driver:           "pgx",
		driverPkg:        "github.com/jackc/pgx",
		connectionString: "user=postgres password=root host=postgres dbname=gosqltest sslmode=disable"}
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
	t.Logf("Driver type: %T", db.Driver())

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
			//t.Logf("Postgres not running, waiting 60 seconds, try %d", i)
			<-time.After(60 * time.Second)
		}
	}
	if !Running("postgres", 5432) {
		t.Fatalf("skipping tests; Postgres not responding on postgres:5432 after 3 tries")
		return
	}

	t.Logf("%s revision: %s", pgx.driverPkg, gitRevision(t, pgx.driverPkg))
	t.Run("pgx: TXQuery", testPGXTxQuery)
	t.Run("pgx: Blobs", testPGXBlobs)
	t.Run("pgx: ManyQueryRow", testPGXPreparedStmt)
	t.Run("pgx: PreparedStmt", testPGXManyQueryRow)

	t.Logf("%s revision: %s", pq.driverPkg, gitRevision(t, pq.driverPkg))
	t.Run("pq: TXQuery", testPQTxQuery)
	t.Run("pq: Blobs", testPQBlobs)
	t.Run("pq: ManyQueryRow", testPQManyQueryRow)
	t.Run("pq: PreparedStmt", testPQPreparedStmt)

	t.Logf("%s revision: %s", gopg.driverPkg, gitRevision(t, gopg.driverPkg))
	t.Run("gopg: TXQuery", testGoPGTxQuery)
	t.Run("gopg: Blobs", testGoPGBlobs)
	t.Run("gopg: ManyQueryRow", testGoPGManyQueryRow)
	t.Run("gopg: PreparedStmt", testGoPGPreparedStmt)
}

func testPQTxQuery(t *testing.T)      { pq.RunTest(t, testTxQuery) }
func testPQBlobs(t *testing.T)        { pq.RunTest(t, testBlobs) }
func testPQManyQueryRow(t *testing.T) { pq.RunTest(t, testManyQueryRow) }
func testPQPreparedStmt(t *testing.T) { pq.RunTest(t, testPreparedStmt) }

func testPGXTxQuery(t *testing.T)      { pgx.RunTest(t, testTxQuery) }
func testPGXBlobs(t *testing.T)        { pgx.RunTest(t, testBlobs) }
func testPGXPreparedStmt(t *testing.T) { pgx.RunTest(t, testPreparedStmt) }
func testPGXManyQueryRow(t *testing.T) { pgx.RunTest(t, testManyQueryRow) }

func testGoPGTxQuery(t *testing.T)      { gopg.RunTest(t, testTxQuery) }
func testGoPGBlobs(t *testing.T)        { gopg.RunTest(t, testBlobs) }
func testGoPGManyQueryRow(t *testing.T) { gopg.RunTest(t, testManyQueryRow) }
func testGoPGPreparedStmt(t *testing.T) { gopg.RunTest(t, testPreparedStmt) }
