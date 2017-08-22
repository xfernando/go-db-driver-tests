package sqltest

import (
	"database/sql"
	"math/rand"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

const TablePrefix = "gosqltest_"

type mysqlDB struct {
	driver           string
	connectionString string
	container        string
	once             sync.Once
	running          bool
}

type pqDB struct {
}

type oracleDB struct {
}

type params struct {
	dbType Tester
	*testing.T
	*sql.DB
}

type Tester interface {
	RunTest(*testing.T, func(params))
}

var (
	goMysqlDB Tester = &mysqlDB{driver: "mysql", container: "mysql", connectionString: "root:root@/gosqltest"}
	myMysqlDB Tester = &mysqlDB{driver: "mymysql", container: "mysql", connectionString: "gosqltest/root/root"}
	pq        Tester = &pqDB{}
	oracle    Tester = &oracleDB{}
)

func (m *mysqlDB) Running() bool {
	m.once.Do(func() {
		c, err := net.Dial("tcp", "localhost:3306")
		if err == nil {
			m.running = true
			c.Close()
		}
	})
	return m.running
}

func startContainer(t *testing.T, c string) {
	cmd := exec.Command("docker-compose", "up", "-d", c)
	if err := cmd.Start(); err != nil {
		t.Fatalf("could not start %s using docker-compose: %v", c, err)
	}
}

func stopContainer(t *testing.T, c string) {
	cmd := exec.Command("docker-compose", "down", c)
	if err := cmd.Start(); err != nil {
		t.Fatalf("could not stop %s using docker-compose: %v", c, err)
	}
}

func (m *mysqlDB) RunTest(t *testing.T, fn func(params)) {
	//	startContainer(t, m.container)
	//	defer stopContainer(t, m.container)
	// wait 60 seconds for db to start
	// <-time.After(60 * time.Second)

	if !m.Running() {
		t.Logf("skipping test; no MySQL running on localhost:3306")
		return
	}
	db, err := sql.Open(m.driver, m.connectionString)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}

	params := params{m, t, db}

	// Drop all tables in the test database.
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		t.Fatalf("failed to enumerate tables: %v", err)
	}
	for rows.Next() {
		var table string
		if rows.Scan(&table) == nil &&
			strings.HasPrefix(strings.ToLower(table), strings.ToLower(TablePrefix)) {
			params.mustExec("DROP TABLE " + table)
		}
	}

	fn(params)
}

func (p *pqDB) RunTest(t *testing.T, fn func(params)) {
	//TODO
}

func (o *oracleDB) RunTest(t *testing.T, fn func(params)) {
	//TODO
}

func (t params) mustExec(sql string, args ...interface{}) sql.Result {
	res, err := t.DB.Exec(sql, args...)
	if err != nil {
		t.Fatalf("Error running %q: %v", sql, err)
	}
	return res
}

func testPreparedStmt(t params) {
	t.mustExec("CREATE TABLE " + TablePrefix + "t (count INT)")
	sel, err := t.Prepare("SELECT count FROM " + TablePrefix + "t ORDER BY count DESC")
	if err != nil {
		t.Fatalf("prepare 1: %v", err)
	}
	ins, err := t.Prepare(t.q("INSERT INTO " + TablePrefix + "t (count) VALUES (?)"))
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

var qrx = regexp.MustCompile(`\?`)

// q converts "?" characters to $1, $2, $n on postgres, :1, :2, :n on Oracle
func (t params) q(sql string) string {
	var pref string
	switch t.dbType {
	case pq:
		pref = "$"
	case oracle:
		pref = ":"
	default:
		return sql
	}
	n := 0
	return qrx.ReplaceAllStringFunc(sql, func(string) string {
		n++
		return pref + strconv.Itoa(n)
	})
}

func testTxQuery(t params) {
	tx, err := t.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	_, err = t.DB.Exec("create table " + TablePrefix + "foo (id integer primary key, name varchar(50))")
	if err != nil {
		t.Logf("cannot drop table "+TablePrefix+"foo: %s", err)
	}

	_, err = tx.Exec(t.q("insert into "+TablePrefix+"foo (id, name) values(?,?)"), 1, "bob")
	if err != nil {
		t.Fatal(err)
	}

	r, err := tx.Query(t.q("select name from "+TablePrefix+"foo where id = ?"), 1)
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

func TestTxQuery_GoMySQL(t *testing.T) { goMysqlDB.RunTest(t, testTxQuery) }
func TestTxQuery_MyMySQL(t *testing.T) { myMysqlDB.RunTest(t, testTxQuery) }
