package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func isValidTimeframe(timeframe string) bool {
	validTimeframes := []string{"day", "week", "month", "all", "alldata"}
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

func twoDecimals(num float64) float64 {
	return float64(int64(num*100)) / 100.0
}

func twoDecimalsPercent(num float64) float64 {
	return float64(int64(num*10000)) / 100.0
}

func boolPointer(val bool) *bool {
	b := val
	return &b
}

func benchmark() {
	start := time.Now()
	collectionName := "balances"

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})
	findOptions.SetAllowDiskUse(true)
	findOptions.SetProjection(bson.M{"value": 1, "_id": 0})

	c := db.Collection(collectionName)

	selection := bson.D{{Key: "username", Value: "ace"}}

	cur, err := c.Find(context.TODO(), selection, findOptions)

	if err != nil {
		log.Println(err)
		return
	}

	var totalvalue float64

	for cur.Next(context.TODO()) {
		type valueType struct {
			Value float64 `json:"value"`
		}
		var curData valueType
		err := cur.Decode(&curData)

		if err != nil {
			log.Panicln(err)
			return
		}
		//totalvalue += curData.Value

	}

	fmt.Println("Normal", time.Since(start).Seconds(), "s, found", totalvalue, "records")
}

func benchmarkGo() {
	start := time.Now()
	collectionName := "balances"

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})
	findOptions.SetAllowDiskUse(true)
	findOptions.SetBatchSize(10000)

	c := db.Collection(collectionName)

	selection := bson.D{{Key: "username", Value: "ace"}}

	cur, err := c.Find(context.TODO(), selection, findOptions)

	if err != nil {
		log.Println(err)
		return
	}

	counter := 0
	var wg sync.WaitGroup

	for cur.Next(context.TODO()) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			var curData balanceData
			err := cur.Decode(&curData)

			if err != nil {
				log.Panicln(err)
				return
			}
			counter++
		}()
	}

	wg.Wait()

	fmt.Println("Goroutine", time.Since(start).Seconds(), "s, found", counter, "records")
}
