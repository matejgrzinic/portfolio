package db

type DB struct {
	API
}

func NewDbAccess(dbName string) *DB {
	db := new(DB)
	db.API = newDBAccess(dbName)
	return db
}
