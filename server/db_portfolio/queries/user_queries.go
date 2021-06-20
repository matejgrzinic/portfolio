package db_portfolio

import "github.com/matejgrzinic/portfolio/db"

type DbUserData struct {
	Name string `json:"name"`
}

func GetAllUsers(dba *db.DB) (*[]DbUserData, error) {
	return nil, nil
}
