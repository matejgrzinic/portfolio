package db

import "go.mongodb.org/mongo-driver/mongo/options"

func NewMockedDB() *MockedDB {
	return new(MockedDB)
}

type MockedDB struct {
	QueryRowFunc  func(result interface{}) error
	QueryRowsFunc func(result interface{}, rowFunc func() error) error
	InsertOneFunc func(data interface{}) error
}

func (mdb *MockedDB) QueryRow(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error {
	return mdb.QueryRowFunc(result)
}

func (mdb *MockedDB) QueryRows(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error {
	return mdb.QueryRowsFunc(result, rowFunc)
}

func (mdb *MockedDB) InsertOne(name string, col string, data interface{}) error {
	return mdb.InsertOneFunc(data)
}

// type MockedDB struct {
// 	QueryRowFunc  func(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error
// 	QueryRowsFunc func(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error
// 	InsertOneFunc func(name string, col string, data interface{}) error
// }

// func (mdb *MockedDB) QueryRow(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error {
// 	return mdb.QueryRowFunc(name, col, filter, options, result)
// }

// func (mdb *MockedDB) QueryRows(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error {
// 	return mdb.QueryRowsFunc(name, col, filter, options, result, rowFunc)
// }

// func (mdb *MockedDB) InsertOne(name string, col string, data interface{}) error {
// 	return mdb.InsertOneFunc(name, col, data)
// }
