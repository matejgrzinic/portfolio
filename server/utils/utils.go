package utils

import (
	"fmt"
	"reflect"
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
