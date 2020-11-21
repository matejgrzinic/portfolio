package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transaction struct {
	Type         string  `json:"type"`
	CurrencyType string  `json:"currency-type"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
	Price        float64 `json:"price"`
	Value        float64 `json:"value"`
	Description  string  `json:"description"`
	Time         int64   `json:"time"`
	User         string  `json:"username"`
}

type trade struct {
	Gain transaction `json:"gain"`
	Loss transaction `json:"loss"`
	Time int64       `json:"time"`
	User string      `json:"username"`
}

func isValidTimeframe(timeframe string) bool {
	validTimeframes := []string{"day", "week", "month", "all"}
	for _, t := range validTimeframes {
		if timeframe == t {
			return true
		}
	}
	return false
}

func isValidTable(table string) bool {
	validTables := []string{"portfolio", "gain", "loss"}
	for _, t := range validTables {
		if table == t {
			return true
		}
	}
	return false
}

func isValidCurrecyType(currencyType string) bool {
	viableCurrencyTypes := []string{"crypto", "cash", "stock"}
	for _, t := range viableCurrencyTypes {
		if currencyType == t {
			return true
		}
	}
	return false
}

func isValidTransactionType(currencyType string) bool {
	viableCurrencyTypes := []string{"gain", "loss"}
	for _, t := range viableCurrencyTypes {
		if currencyType == t {
			return true
		}
	}
	return false
}

func isValidTransaction(parameters map[string]string, username string) error {
	tType := parameters["type"]
	currencyType := parameters["currency-type"]
	currency := parameters["currency"]
	amount := parameters["amount"]

	validType := map[string]bool{"gain": true, "loss": true}
	validCurrencyType := map[string]bool{"cash": true, "crypto": true, "stock": true}

	if _, ok := validType[tType]; !ok {
		return fmt.Errorf("invalid type")
	}

	if _, ok := validCurrencyType[currencyType]; !ok {
		return fmt.Errorf("invalid currency type")
	}

	if _, ok := latestPriceData.Rates[currencyType][currency]; !ok {
		return fmt.Errorf("invalid currency")
	}

	f, err := strconv.ParseFloat(amount, 64)
	if err != nil || f <= 0 {
		return fmt.Errorf("invalid amount")
	}

	return nil
}

func createNewTransaction(parameters map[string]string, username string) *transaction {
	tType := parameters["type"]
	currencyType := parameters["currency-type"]
	currency := parameters["currency"]
	amount, _ := strconv.ParseFloat(parameters["amount"], 64)
	description := parameters["description"]

	return &transaction{
		Type:         tType,
		CurrencyType: currencyType,
		Currency:     currency,
		Amount:       amount,
		Price:        latestPriceData.Rates[currencyType][currency],
		Value:        latestPriceData.Rates[currencyType][currency] * amount,
		Description:  description,
		Time:         time.Now().Unix(),
		User:         username,
	}
}

func getAllCurrencies(currencyType string) []string {
	data := []string{}
	for k := range latestPriceData.Rates[currencyType] {
		data = append(data, k)
	}
	sort.Strings(data)
	return data
}

func parseTradeParameters(parameters map[string]string) (map[string]string, map[string]string) {
	gain := map[string]string{"type": "gain", "currency-type": parameters["buy-type"], "currency": parameters["buy-currency"], "amount": parameters["buy-amount"]}
	loss := map[string]string{"type": "loss", "currency-type": parameters["sell-type"], "currency": parameters["sell-currency"], "amount": parameters["sell-amount"]}
	return gain, loss
}

func getUserCurrencies(user string, currencyType string) []string {
	output := []string{}
	userData, err := getUserLatestPortfolio(user)
	if err != nil {
		log.Println("error getting user data while getting his currencies. getUserCurrencies")
		return output
	}

	for _, cType := range userData.Data {
		if cType.Type == currencyType {
			for _, c := range cType.Data {
				output = append(output, c.Symbol)
			}
			break
		}
	}
	return output
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

func replyReturnMessage(w *http.ResponseWriter, status string, message string) {
	reply := &struct {
		Status  string
		Message string
	}{Status: status, Message: message}

	replyJSON, err := json.Marshal(reply)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(*w, string(replyJSON))
}
