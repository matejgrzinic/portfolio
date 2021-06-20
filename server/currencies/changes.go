package currencies

import (
	"fmt"
	"time"

	"github.com/matejgrzinic/portfolio/external"
	"github.com/matejgrzinic/portfolio/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Currencies) updateChangesMap() error {
	filter := bson.D{{Key: "time", Value: bson.M{"$gt": time.Now().Add(-time.Hour * 24 * 30).Unix()}}}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}})
	findOptions.SetProjection(bson.M{"cryptocurrency": 1, "fiat": 1, "_id": 0})

	type rowType map[string]map[string]external.CurrencyData
	var initialData rowType
	rowData := make(rowType)

	newChangesMap := make(changesMap)

	err := c.dbAPI.QueryRows(
		"update changes map",
		"price",
		filter,
		findOptions,
		rowData,
		func() error {
			if initialData == nil {
				cpy, err2 := utils.CopyMap(rowData)
				if err2 != nil {
					return err2
				}
				initialData = cpy.(rowType)
			}

			for key, value := range rowData {
				if newChangesMap[key] == nil {
					newChangesMap[key] = make(map[string]map[string]float64)
				}
				for key2, value2 := range value {
					if newChangesMap[key][key2] == nil {
						newChangesMap[key][key2] = make(map[string]float64)
					}
					newChangesMap[key][key2]["hour"] = value2.Price
				}
			}
			return nil
		},
	)

	if err != nil {
		return err
	}

	for key, value := range newChangesMap {
		for key2, value2 := range value {
			for time, value3 := range value2 {
				fmt.Println(newChangesMap[key][key2][time])
				fmt.Println(initialData[key][key2].Price)
				newChangesMap[key][key2][time] = initialData[key][key2].Price/value3 - 1.0
				fmt.Println(newChangesMap[key][key2][time])
			}
		}
	}

	c.changes = newChangesMap
	return nil
}
