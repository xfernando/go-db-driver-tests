package sqltest

import (
	"fmt"
	"testing"
)

var (
	oracle = &oracleDB{driver: "goracle", connectionString: "oracle://user:pass@db/"}
)

type oracleDB database

func (o *oracleDB) RunTest(t *testing.T, fn func(Tester)) {
	//TODO
}

func (o *oracleDB) SQLBlobParam(size int) string {
	return fmt.Sprintf("RAW(%d)", size)
}

func (o *oracleDB) q(sql string) string {
	// n := 0
	// return qrx.ReplaceAllStringFunc(sql, func(string) string {
	// 	n++
	// 	return ":" + strconv.Itoa(n)
	// })
	return ""
}
