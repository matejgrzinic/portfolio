package portfolio

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimelineData struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

func (p *Portfolio) UserTimeline(user *User, timeframe string) ([]TimelineData, error) {
	validTimeframe := map[string]int64{
		"day":   time.Now().Add(-time.Hour * 24).Unix(),
		"week":  time.Now().Add(-time.Hour * 24 * 7).Unix(),
		"month": time.Now().Add(-time.Hour * 24 * 30).Unix(),
		"all":   0,
	}

	timeframeInt, ok := validTimeframe[timeframe]
	if !ok {
		return nil, fmt.Errorf("invalid timeframe: %v", timeframe)
	}

	// contains user Currency -> amount in hashmap
	hashData := make(map[string]map[string]float64)

	ts, err := p.AllUserTransactions(user)
	if err != nil {
		return nil, err
	}

	addToHashDataLoop := func(tc []TransactionCurrency, multiplier float64) {
		for _, e := range tc {
			if _, ok := hashData[e.CurrencyType]; !ok {
				hashData[e.CurrencyType] = make(map[string]float64)
			}
			hashData[e.CurrencyType][e.Symbol] += e.Amount * multiplier
		}
	}

	updateHashMap := func(t Transaction) {
		addToHashDataLoop(t.Gains, 1.0)
		addToHashDataLoop(t.Losses, -1.0)
	}
	transactionIndex := 0

	result := make([]TimelineData, 0)

	var row Prices
	err = p.DB().QueryRows( // TODO move to prices.go ?
		"all prices",
		"price",
		bson.M{"time": bson.M{"$gt": timeframeInt}},
		options.Find().SetSort(bson.D{{Key: "time", Value: 1}}),
		&row,
		func() error {
			for transactionIndex < len(ts) && ts[transactionIndex].Time < row.Time {
				updateHashMap(ts[transactionIndex])
				transactionIndex++
			}
			if len(hashData) > 0 && timeframeInt < row.Time {
				var value float64

				for currType, curs := range hashData {
					for symbol, amount := range curs {
						if rowCurr, ok := row.Data[currType][symbol]; ok {
							value += rowCurr.Price * amount
						} else {
							log.Printf("row data does not contain: %v %v", currType, symbol)
							return nil
						}
					}
				}

				td := TimelineData{Time: row.Time, Value: value}
				result = append(result, td)
			}
			return nil
		},
	)

	if len(result) > 1000 {
		eachN := len(result) / 1000
		for i := 1; i <= 1000; i++ {
			result[i] = result[i*eachN]
		}
		result = result[:1000]
	}

	if len(hashData) > 0 {
		var value float64
		for currType, curs := range hashData {
			for symbol, amount := range curs {
				cRn, err := p.Currencies().GetCurrency(currType, symbol)
				if err != nil {
					return nil, err
				}
				value += cRn.Price * amount
			}
		}
		td := TimelineData{Time: row.Time, Value: value}
		result = append(result, td)
	}

	fmt.Println(len(result))
	return result, err
}
