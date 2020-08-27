package main

import (
	"context"
	"fmt"
	"log"
	"math"
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

type userData struct {
	Username     string `json:"username"`
	Passwordhash string `json:"passwordhash"`
	Created      int64  `json:"created"`
	Active       bool   `json:"active"`
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
	const n int = 10

	collectionName := "balances"

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})

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
		selection = bson.D{{Key: "username", Value: user}, {Key: "time", Value: bson.M{"$gt": time.Now().Add(-time.Hour * 24 * 7 * 30).Unix()}}}
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

	var tmpV []float64
	var tmpT []string

	for cur.Next(context.TODO()) {
		tmpV = append(tmpV, float64(int64(cur.Current.Lookup("value").Double()*100))/100.0)

		t := time.Unix(cur.Current.Lookup("time").Int64(), 0)
		tmpT = append(tmpT, t.UTC().Format("2006-01-02T15:04:05-0700")) // check for rfc2822 it looks better
		// output.Time = append(output.Time, t.Format("2006.01.02 15:04")) // shows warning in console
	}

	l := int(math.Ceil(float64(len(tmpV)) / float64(n)))

	for i := range tmpV {
		if i%l == 0 {
			output.Time = append(output.Time, tmpT[i])
			output.Value = append(output.Value, tmpV[i])
		}
	}

	output.Username = user
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

func updateUserPortfolioData(data *balanceData) {
	var totalSum float64
	for i, e := range data.Data {
		var typeSum float64
		if e.Type == "crypto" {
			for j, f := range e.Data {
				if newPrice, ok := latestPriceData.Crypto[f.Symbol]; ok {
					data.Data[i].Data[j].Price = newPrice * latestPriceData.USDRates.Rates.EUR
					data.Data[i].Data[j].Value = newPrice*f.Amount + 1
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

func updateUserPortfolio(user string, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := getUserLatestPortfolio(user)
	if err != nil {
		return
	}
	updateUserPortfolioData(data)
	insertUserPortfolio(user, data)
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
	exampleBalance(user.Username)
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
