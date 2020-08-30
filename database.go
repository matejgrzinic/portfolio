package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type graphData struct {
	Time     []string
	Value    []float64
	TimeUnit string
	Username string
}

type currencyData struct {
	Currency    string
	Symbol      string
	Amount      float64
	Price       float64
	Value       float64
	HourChange  float64
	DayChange   float64
	WeekChange  float64
	MonthChange float64
}

type userData struct {
	Username     string `json:"username"`
	Passwordhash string `json:"passwordhash"`
	Created      int64  `json:"created"`
	Active       bool   `json:"active"`
	Started      bool   `json:"started"`
}

type balanceData struct {
	Data []struct {
		Type string `json:"type"`
		Data []struct {
			Symbol string  `json:"symbol"`
			Amount float64 `json:"amount"`
			Price  float64 `json:"price"`
			Value  float64 `json:"value"`
		} `json:"data"`
		Value float64 `json:"value"`
	} `json:"data"`
	Value    float64 `json:"value"`
	Time     int64   `json:"time"`
	Username string  `json:"username"`
}

func setup() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	const dbName = "test"
	db := client.Database(dbName)

	db.CreateCollection(context.TODO(), "users")
	db.CreateCollection(context.TODO(), "balances")
	db.CreateCollection(context.TODO(), "transactions")

	return db
}

func countall() {
	data := db.Collection("portfolio")
	fmt.Println(data.EstimatedDocumentCount(context.TODO()))
}

func getUserTimeframeData(user string, timeframe string) graphData {
	var output graphData

	const n int = 50

	collectionName := "balances"

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})
	findOptions.AllowDiskUse = boolPointer(true)
	findOptions.SetProjection(bson.M{"value": 1, "time": 1, "_id": 0})

	c := db.Collection(collectionName)

	var selection bson.D

	switch timeframe {
	case "day":
		selection = bson.D{{Key: "username", Value: user}, {Key: "time", Value: bson.M{"$gt": time.Now().Add(-time.Hour * 24).Unix()}}}
		output.TimeUnit = "hour"
		break
	case "week":
		selection = bson.D{{Key: "username", Value: user}, {Key: "time", Value: bson.M{"$gt": time.Now().Add(-time.Hour * 24 * 7).Unix()}}}
		output.TimeUnit = "day"
		break
	case "month":
		selection = bson.D{{Key: "username", Value: user}, {Key: "time", Value: bson.M{"$gt": time.Now().Add(-time.Hour * 24 * 30).Unix()}}}
		output.TimeUnit = "week"
		break
	case "all":
		selection = bson.D{{Key: "username", Value: user}}
		output.TimeUnit = "month"
		break
	}

	cur, err := c.Find(context.TODO(), selection, findOptions)

	if err != nil {
		log.Println(err)
		return output // ?
	}

	type valueType struct {
		Value float64 `json:"value"`
		Time  int64   `json:"time"`
	}

	var queryLen int64
	if timeframe != "all" {
		queryLen, err = c.CountDocuments(context.TODO(), selection)
	} else {
		queryLen, err = c.EstimatedDocumentCount(context.TODO(), nil)

	}

	if err != nil {
		log.Println(err)
		return output // ?
	}

	l := int(math.Ceil(float64(queryLen) / float64(n)))
	counter := 0

	for cur.Next(context.TODO()) {
		if counter%l == 0 {
			var curData valueType
			err := cur.Decode(&curData)

			if err != nil {
				log.Panicln(err)
				continue
			}

			output.Value = append(output.Value, twoDecimals(curData.Value))
			t := time.Unix(curData.Time, 0)
			output.Time = append(output.Time, t.UTC().Format("2006-01-02T15:04:05-0700")) // check for rfc2822 it looks better
			// output.Time = append(output.Time, t.Format("2006.01.02 15:04")) // shows warning in console
		}
		counter++
	}

	output.Username = user
	return output
}

func getUserDisplayValues(user string) []currencyData {
	var output []currencyData

	collectionName := "balances"

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})
	findOptions.SetAllowDiskUse(true)

	c := db.Collection(collectionName)

	selection := bson.D{{Key: "username", Value: user}, {Key: "time", Value: bson.M{"$gt": time.Now().Add(-time.Hour * 24 * 30).Unix()}}}

	cur, err := c.Find(context.TODO(), selection, findOptions)

	if err != nil {
		log.Println(err)
		return output // ?
	}

	currencyDisplay := map[string]*currencyData{}

	first := true
	var timeFirst int64

	const hourSeconds int64 = 60 * 60
	const daySeconds int64 = 60 * 60 * 24
	const weekSeconds int64 = 60 * 60 * 24 * 7

	for cur.Next(context.TODO()) {
		var curData balanceData
		err := cur.Decode(&curData)

		if err != nil {
			log.Panicln(err)
			continue
		}

		if first {
			for _, e := range curData.Data {
				for _, f := range e.Data {
					currencyDisplay[f.Symbol] = &currencyData{
						Currency: e.Type,
						Symbol:   f.Symbol,
						Price:    f.Price,
						Amount:   f.Amount,
						Value:    f.Value,
					}
				}
			}
			first = false
			timeFirst = time.Now().Unix()

		} else {
			for _, e := range curData.Data {
				for _, f := range e.Data {
					if fCur, ok := currencyDisplay[f.Symbol]; ok {
						fCur.MonthChange = fCur.Price/f.Price - 1
						timeDiff := timeFirst - curData.Time

						if timeDiff < weekSeconds {
							fCur.WeekChange = fCur.Price/f.Price - 1

							if timeDiff < daySeconds {
								fCur.DayChange = fCur.Price/f.Price - 1

								if timeDiff < hourSeconds {
									fCur.HourChange = fCur.Price/f.Price - 1
								}
							}
						}
					}
				}
			}
		}

	}

	for _, value := range currencyDisplay {
		value.Price = twoDecimals(value.Price)
		value.Amount = twoDecimals(value.Amount)
		value.Value = twoDecimals(value.Value)
		value.HourChange = twoDecimalsPercent(value.HourChange)
		value.DayChange = twoDecimalsPercent(value.DayChange)
		value.WeekChange = twoDecimalsPercent(value.WeekChange)
		value.MonthChange = twoDecimalsPercent(value.MonthChange)
		output = append(output, *value)
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Value > output[j].Value
	})

	return output
}

func getUserLatestPortfolio(user string) (*balanceData, error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})

	c := db.Collection("balances")
	cur := c.FindOne(context.TODO(), bson.D{{Key: "username", Value: user}}, findOptions)

	var data *balanceData
	err := cur.Decode(&data)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}

func updateUserPortfolio(user string, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := getUserLatestPortfolio(user)
	if err != nil {
		return
	}
	updateUserPortfolioData(data)
	insertUserPortfolio(user, data)
}

func updateUserPortfolioData(data *balanceData) {
	var totalSum float64
	for i, e := range data.Data {
		var typeSum float64
		if e.Type == "crypto" {
			for j, f := range e.Data {
				if newPrice, ok := latestPriceData.Crypto[f.Symbol]; ok {
					data.Data[i].Data[j].Price = newPrice * latestPriceData.USDRates.Rates.EUR
					data.Data[i].Data[j].Value = data.Data[i].Data[j].Price*f.Amount + 1
				}
				typeSum += data.Data[i].Data[j].Value
			}
			data.Data[i].Value = typeSum
		}
		totalSum += data.Data[i].Value
	}
	data.Value = totalSum
	data.Time = time.Now().Unix()
}

func insertUserPortfolio(user string, data *balanceData) {
	_, err := db.Collection("balances").InsertOne(context.TODO(), data)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("Updated user", user)
}

func insertNewUser(user *userData) {
	_, err := db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
	}
	//exampleBalance(user.Username)
}

func getUserData(user string) (*userData, error) {
	findOptions := options.FindOne()

	c := db.Collection("users")
	cur := c.FindOne(context.TODO(), bson.D{{Key: "username", Value: user}}, findOptions)

	var data *userData
	err := cur.Decode(&data)

	if err != nil {
		log.Println(err)
		return data, err
	}

	return data, nil
}

func startUpdateUserPortfolioInterval() {
	for range time.Tick(time.Second) {
		if len(latestPriceData.Crypto) > 0 {
			updateLoop()
			break
		}
	}

	for range time.Tick(time.Minute) {
		updateLoop()
	}
}

func getActiveUsers() []string {
	var usernames []string

	findOptions := options.Find()

	ctx := context.Background()

	c := db.Collection("users")
	cur, err := c.Find(ctx, bson.D{{Key: "active", Value: true}}, findOptions)
	defer cur.Close(ctx)

	if err != nil {
		log.Println(err)
	}

	for cur.Next(ctx) {
		s := cur.Current.Lookup("username").StringValue()
		usernames = append(usernames, s)
	}

	return usernames
}

func getAllUsers() []string {
	var usernames []string

	findOptions := options.Find()

	ctx := context.Background()

	c := db.Collection("users")
	cur, err := c.Find(ctx, bson.D{{}}, findOptions)
	defer cur.Close(ctx)

	if err != nil {
		log.Println(err)
	}

	for cur.Next(ctx) {
		s := cur.Current.Lookup("username").StringValue()
		usernames = append(usernames, s)
	}

	return usernames
}

func updateLoop() {
	users := getActiveUsers()
	// start := time.Now()
	var wg sync.WaitGroup
	for _, c := range users {
		wg.Add(1)
		go updateUserPortfolio(c, &wg)
	}
	wg.Wait()
	// fmt.Println("This round took ", time.Since(start).Microseconds(), "microseconds")
}
