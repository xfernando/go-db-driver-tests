package sqltest

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"testing"
)

const TablePrefix = "gosqltest_"

type Tester interface {
	RunTest(*testing.T, func(Tester))
	SQLBlobParam(size int) string
	DB() *sql.DB
	T() *testing.T
	// q converts "?" characters to $1, $2, $n on postgres, :1, :2, :n on Oracle
	q(sql string) string
}

func Running(host string, port int) bool {
	c, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err == nil {
		c.Close()
		return true
	}
	return false
}

func mustExec(t Tester, sql string, args ...interface{}) sql.Result {
	res, err := t.DB().Exec(sql, args...)
	if err != nil {
		t.T().Fatalf("Error running %q: %v", sql, err)
	}
	return res
}

func testPreparedStmt(t Tester) {
	mustExec(t, "CREATE TABLE "+TablePrefix+"t (count INT)")
	sel, err := t.DB().Prepare("SELECT count FROM " + TablePrefix + "t ORDER BY count DESC")
	if err != nil {
		t.T().Fatalf("prepare 1: %v", err)
	}
	ins, err := t.DB().Prepare(t.q("INSERT INTO " + TablePrefix + "t (count) VALUES (?)"))
	if err != nil {
		t.T().Fatalf("prepare 2: %v", err)
	}

	for n := 1; n <= 3; n++ {
		if _, err := ins.Exec(n); err != nil {
			t.T().Fatalf("insert(%d) = %v", n, err)
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
					t.T().Errorf("Query: %v", err)
					return
				}
				if _, err := ins.Exec(rand.Intn(100)); err != nil {
					t.T().Errorf("Insert: %v", err)
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

func testTxQuery(t Tester) {
	tx, err := t.DB().Begin()
	if err != nil {
		t.T().Fatal(err)
	}
	defer tx.Rollback()

	_, err = t.DB().Exec("create table " + TablePrefix + "foo (id integer primary key, name varchar(50))")
	if err != nil {
		t.T().Logf("cannot drop table "+TablePrefix+"foo: %s", err)
	}

	_, err = tx.Exec(t.q("insert into "+TablePrefix+"foo (id, name) values(?,?)"), 1, "bob")
	if err != nil {
		t.T().Fatal(err)
	}

	r, err := tx.Query(t.q("select name from "+TablePrefix+"foo where id = ?"), 1)
	if err != nil {
		t.T().Fatal(err)
	}
	defer r.Close()

	if !r.Next() {
		if r.Err() != nil {
			t.T().Fatal(err)
		}
		t.T().Fatal("expected one rows")
	}

	var name string
	err = r.Scan(&name)
	if err != nil {
		t.T().Fatal(err)
	}
}

func testBlobs(t Tester) {
	var blob = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	mustExec(t, "create table "+TablePrefix+"foo (id integer primary key, bar "+t.SQLBlobParam(16)+")")
	mustExec(t, t.q("insert into "+TablePrefix+"foo (id, bar) values(?,?)"), 0, blob)

	want := fmt.Sprintf("%x", blob)

	b := make([]byte, 16)
	err := t.DB().QueryRow(t.q("select bar from "+TablePrefix+"foo where id = ?"), 0).Scan(&b)
	got := fmt.Sprintf("%x", b)
	if err != nil {
		t.T().Errorf("[]byte scan: %v", err)
	} else if got != want {
		t.T().Errorf("for []byte, got %q; want %q", got, want)
	}

	err = t.DB().QueryRow(t.q("select bar from "+TablePrefix+"foo where id = ?"), 0).Scan(&got)
	want = string(blob)
	if err != nil {
		t.T().Errorf("string scan: %v", err)
	} else if got != want {
		t.T().Errorf("for string, got %q; want %q", got, want)
	}
}

func testManyQueryRow(t Tester) {
	if testing.Short() {
		t.T().Logf("skipping in short mode")
		return
	}
	mustExec(t, "create table "+TablePrefix+"foo (id integer primary key, name varchar(50))")
	mustExec(t, t.q("insert into "+TablePrefix+"foo (id, name) values(?,?)"), 1, "bob")
	var name string
	for i := 0; i < 10000; i++ {
		err := t.DB().QueryRow(t.q("select name from "+TablePrefix+"foo where id = ?"), 1).Scan(&name)
		if err != nil || name != "bob" {
			t.T().Fatalf("on query %d: err=%v, name=%q", i, err, name)
		}
	}
}
