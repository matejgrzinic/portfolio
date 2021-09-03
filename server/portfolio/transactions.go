package portfolio

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionCurrency struct {
	CurrencyType string  `json:"type"`
	Symbol       string  `json:"symbol"`
	Amount       float64 `json:"amount,string"`
}

type Transaction struct {
	User   string                `json:"user"`
	Gains  []TransactionCurrency `json:"gains"`
	Losses []TransactionCurrency `json:"losses"`
	Time   int64                 `json:"time"`
	Note   string                `json:"note"`
}

func (p *Portfolio) AllUserTransactions(user *User) ([]Transaction, error) {
	data := make([]Transaction, 0)
	var row Transaction
	err := p.DB().QueryRows(
		"all user transactions",
		"transactions",
		bson.M{"user": user.Name},
		options.Find().SetSort(bson.D{{Key: "time", Value: 1}}),
		&row,
		func() error {
			var cpyT Transaction
			copier.CopyWithOption(&cpyT, &row, copier.Option{DeepCopy: true})
			data = append(data, cpyT)
			return nil
		},
	)
	return data, err
}

func (p *Portfolio) InsertTransaction(t *Transaction) error {
	err := p.DB().InsertOne(
		"insert transaction",
		"transactions",
		t,
	)
	return err
}
