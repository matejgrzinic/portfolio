package db

import "go.mongodb.org/mongo-driver/mongo/options"

type MockedDB struct {
	QueryRowFunc  func(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error
	QueryRowsFunc func(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error
}

func (mdb *MockedDB) QueryRow(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error {
	return mdb.QueryRowFunc(name, col, filter, options, result)
}

func (mdb *MockedDB) QueryRows(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error {
	return mdb.QueryRowsFunc(name, col, filter, options, result, rowFunc)
}
