package util

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite3Conn(path string, dsn ...string) *gorm.DB {
	var DSN string
	if len(dsn) < 1 {
		DSN = fmt.Sprintf("file:%s?cache=private&_journal_mode=WAL", path)
	} else {
		DSN = fmt.Sprintf("file:%s?%s", path, dsn[0])
	}
	db, err := gorm.Open(sqlite.Open(DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		panic(err)
	}
	return db
}
