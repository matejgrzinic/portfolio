package db

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type API interface {
	QueryRow(name string, col string, filter interface{}, options *options.FindOneOptions, result interface{}) error
	QueryRows(name string, col string, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error
}
