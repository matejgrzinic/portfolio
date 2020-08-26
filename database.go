package main

import (
	"context"
	"fmt"
	"log"
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

type dbData struct {
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

func price(user string) graphData {
	var output graphData

	collectionName := "balance"

	findOptions := options.Find()
	findOptions.SetLimit(int64(60))
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})

	c := db.Collection(collectionName)
	cur, err := c.Find(context.TODO(), bson.D{{Key: "username", Value: user}}, findOptions)

	if err != nil {
		log.Println(err)
		return output // ?
	}

	for cur.Next(context.TODO()) {
		output.Value = append(output.Value, float64(int64(cur.Current.Lookup("value").Double()*100))/100.0)

		t := time.Unix(cur.Current.Lookup("time").Int64(), 0)                         // remove /1000 when updating db with unix timestamp
		output.Time = append(output.Time, t.UTC().Format("2006-01-02T15:04:05-0700")) // check for rfc2822 it looks better
		// output.Time = append(output.Time, t.Format("2006.01.02 15:04")) // shows warning in console
	}

	output.TimeUnit = "hour"
	output.Username = user
	return output
}

func getUserLatestPortfolio(user string) (*dbData, error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})

	c := db.Collection("balance")
	cur := c.FindOne(context.TODO(), bson.D{{Key: "username", Value: user}}, findOptions)

	var data *dbData
	err := cur.Decode(&data)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}

func updateUserPortfolioData(data *dbData) {
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

func insertUserPortfolio(user string, data *dbData) {
	_, err := db.Collection("balance").InsertOne(context.TODO(), data) // todo change to balances
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("Updated user", user)
}

func insertNewUser(user *userData) {
	fmt.Println("REEEEEEEEEEEEEEE")
	_, err := db.Collection("users").InsertOne(context.TODO(), user) // todo change to balances
	if err != nil {
		log.Println(err)
	}
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
