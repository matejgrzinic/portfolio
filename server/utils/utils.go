package utils

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CopyMap(m interface{}) (interface{}, error) {
	t := reflect.ValueOf(m)
	if t.Kind() != reflect.Map {
		return nil, fmt.Errorf("input is not a map")
	}

	cpy := reflect.MakeMap(reflect.TypeOf(m))
	for _, key := range t.MapKeys() {
		val := t.MapIndex(key)
		if val.Kind() == reflect.Map {
			innerMap, err := CopyMap(val.Interface())
			if err != nil {
				return nil, err
			}
			cpy.SetMapIndex(key, reflect.ValueOf(innerMap))
		} else {
			cpy.SetMapIndex(key, val)
		}
	}

	return cpy.Interface(), nil
}

func QueryRows(name string, col *mongo.Collection, filter interface{}, options *options.FindOptions, result interface{}, rowFunc func() error) error {
	ctx := context.TODO()
	cursor, err := col.Find(
		ctx,
		filter,
		options,
	)

	if err != nil {
		return fmt.Errorf("query rows [%s]: %v", name, err)
	}

	for cursor.Next(ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return fmt.Errorf("query rows [%s], decode row: %v", name, err)
		}

		if err = rowFunc(); err != nil {
			return fmt.Errorf("query rows [%s], rowFunc: %v", name, err)
		}
	}

	return nil
}
