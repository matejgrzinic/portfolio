package main

import (
	"context"
	"log"
	"time"
)

func isValidTimeframe(timeframe string) bool {
	validTimeframes := []string{"day", "week", "month", "all"}
	for _, t := range validTimeframes {
		if timeframe == t {
			return true
		}
	}
	return false
}

func exampleBalance(username string) {
	data := &balanceData{
		Username: username,
		Time:     time.Now().Unix(),
		Value:    1234.123,
		Data: []struct {
			Type string "json:\"type\""
			Data []struct {
				Symbol string  "json:\"symbol\""
				Amount float64 "json:\"amount\""
				Price  float64 "json:\"price\""
				Value  float64 "json:\"value\""
			} "json:\"data\""
			Value float64 "json:\"value\""
		}{{Type: "crypto", Value: 1234.123, Data: []struct {
			Symbol string  "json:\"symbol\""
			Amount float64 "json:\"amount\""
			Price  float64 "json:\"price\""
			Value  float64 "json:\"value\""
		}{{Symbol: "BTC", Amount: 1.23, Price: 2323.23, Value: 1234.123}}}},
	}

	_, err := db.Collection("balances").InsertOne(context.TODO(), data)
	if err != nil {
		log.Println(err)
	}
}
