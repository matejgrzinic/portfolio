package portfolio

import (
	"github.com/matejgrzinic/portfolio/external"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Prices struct {
	Time int64                                       `json:"time"`
	Data map[string]map[string]external.CurrencyData `json:"data"`
}

func (p *Portfolio) AllPrices() ([]Prices, error) {
	data := make([]Prices, 0)
	var row Prices
	err := p.DB().QueryRows(
		"all prices",
		"price",
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "time", Value: 1}}),
		&row,
		func() error {
			data = append(data, row)
			return nil
		},
	)
	return data, err
}
