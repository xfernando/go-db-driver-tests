package sqltest

import (
	"database/sql"
	"io/ioutil"
	"net"
	"strconv"
	"testing"
)

const TablePrefix = "gosqltest_"

type database struct {
	driver           string
	driverPkg        string
	connectionString string
}

type Tester interface {
	// Does test setup (delete old tables, opens and closes db)
	RunTest(*testing.T, func(t *testing.T, db *sql.DB))
}

func Running(host string, port int) bool {
	c, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err == nil {
		c.Close()
		return true
	}
	return false
}

func gitRevision(t *testing.T, pkg string) string {
	b, err := ioutil.ReadFile("/go/src/" + pkg + "/.git/refs/heads/master")
	if err != nil {
		t.Logf("Failed to get current git revision for package %s: %v", pkg, err)
		return ""
	}
	return string(b)
}

func mustExec(db *sql.DB, t *testing.T, sql string, args ...interface{}) sql.Result {
	res, err := db.Exec(sql, args...)
	if err != nil {
		t.Fatalf("Error running %q: %v", sql, err)
	}
	return res
}
