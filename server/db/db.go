package db

import "os"

type DB struct {
	API
}

func NewDbAccess() *DB {
	db := new(DB)
	db.API = newDBAccess(os.Getenv("DB_NAME"))
	return db
}
