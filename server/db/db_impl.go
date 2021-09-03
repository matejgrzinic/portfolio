package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBAccess struct {
	Db  *mongo.Database
	ctx context.Context
}

func (dba *DBAccess) QueryRow(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error {
	cursor := dba.Db.Collection(col).FindOne(
		dba.ctx,
		filter,
		options,
	)

	err := cursor.Decode(result)
	if err != nil {
		return fmt.Errorf("query row [%s], decode row: %v", name, err)
	}

	// TODO errNoRows check

	return nil
}

func (dba *DBAccess) QueryRows(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error {
	cursor, err := dba.Db.Collection(col).Find(
		dba.ctx,
		filter,
		options,
	)

	if err != nil {
		return fmt.Errorf("query rows [%s]: %v", name, err)
	}

	for cursor.Next(dba.ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return fmt.Errorf("query rows [%s], decode row: %v", name, err)
		}

		if err = rowFunc(); err != nil {
			return fmt.Errorf("query rows [%s], rowFunc: %v", name, err)
		}
	}

	// TODO errNoRows check

	return nil
}

func (dba *DBAccess) InsertOne(name string, col string, data interface{}) error {
	_, err := dba.Db.Collection(col).InsertOne(dba.ctx, data)
	if err != nil {
		return fmt.Errorf("insert one [%s]: %v", name, err)
	}

	// TODO errNoRows check

	return nil
}

func newDBAccess(dbName string) *DBAccess {
	clientOptions := options.Client().ApplyURI(os.Getenv("DB_URL"))
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

	mdb := client.Database(dbName)

	db := new(DBAccess)
	db.Db = mdb
	db.ctx = context.TODO()

	return db
}
