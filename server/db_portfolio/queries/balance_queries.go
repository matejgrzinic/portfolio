package db_portfolio

import (
	"context"
	"time"

	"github.com/matejgrzinic/portfolio/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbBalanceData struct {
	User  string           `json:"user"`
	Data  []DbCurrencyData `json:"data"`
	Value float64          `json:"value"`
	Time  int64            `json:"time"`
}

type DbCurrencyData struct {
	CurrencyType string  `json:"type"`
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Amount       float64 `json:"amount"`
	Value        float64 `json:"value"`
}

func GetUserBalance(dba *db.DB, user string) (*DbBalanceData, error) {
	col := dba.Db.Collection("balance")

	queryOptions := options.FindOne()
	queryOptions.SetSort(bson.D{{Key: "time", Value: -1}})

	balance := new(DbBalanceData)
	if err := col.FindOne(context.TODO(), bson.D{{Key: "user", Value: user}}, queryOptions).Decode(&balance); err != nil {
		return nil, err
	}

	return balance, nil
}

func SaveBalance(dba *db.DB, b *DbBalanceData) error {
	col := dba.Db.Collection("balance")

	b.Time = time.Now().Unix()

	if _, err := col.InsertOne(context.TODO(), b); err != nil {
		return err
	}

	return nil
}
