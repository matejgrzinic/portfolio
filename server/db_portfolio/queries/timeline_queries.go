package db_portfolio

import (
	"context"
	"fmt"
	"time"

	"github.com/matejgrzinic/portfolio/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbTimelineData struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

func GetUserTimelineQuery2(dba *db.DB, user string, timeframe string) (*[]DbTimelineData, error) {
	col := dba.Db.Collection("balance")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: 1}})
	findOptions.SetProjection(bson.M{"value": 1, "time": 1, "_id": 0})

	cursor, err := col.Find(context.TODO(), bson.D{{Key: "user", Value: user}}, findOptions)
	if err != nil {
		return nil, err
	}

	var timeline []DbTimelineData
	if err = cursor.All(context.TODO(), &timeline); err != nil {
		return nil, err
	}

	return &timeline, nil
}

func GetUserTimelineQuery(dba *db.DB, user string, timeframe string) (*[]DbTimelineData, error) {
	col := dba.Db.Collection("balance")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})
	findOptions.SetProjection(bson.M{"value": 1, "time": 1, "_id": 0})
	// findOptions.SetBatchSize(4)

	var timeStamp int64
	switch timeframe {
	case "hour":
		timeStamp = time.Now().Add(-time.Hour).Unix()
	case "day":
		timeStamp = time.Now().Add(-time.Hour * 24).Unix()
	case "week":
		timeStamp = time.Now().Add(-time.Hour * 24 * 7).Unix()
	case "month":
		timeStamp = time.Now().Add(-time.Hour * 24 * 30).Unix()
	case "year":
		timeStamp = time.Now().Add(-time.Hour * 24 * 365).Unix()
	case "all":
		timeStamp = 0
	default:
		return nil, fmt.Errorf("invalid timeframe: %s", timeframe)
	}

	selection := bson.D{{Key: "user", Value: user}, {Key: "time", Value: bson.M{"$gt": timeStamp}}}
	cursor, err := col.Find(context.TODO(), selection, findOptions)
	if err != nil {
		return nil, err
	}

	// queryLen, err := col.CountDocuments(context.TODO(), selection)
	// if err != nil {
	// 	return nil, err
	// }

	//n := int(math.Ceil(float64(queryLen) / float64(100)))
	var timeline []DbTimelineData
	for cursor.Next(context.TODO()) {
		var dtd DbTimelineData

		err = cursor.Decode(&dtd)
		if err != nil {
			return nil, err
		}

		timeline = append(timeline, dtd)
	}

	return &timeline, nil
}
