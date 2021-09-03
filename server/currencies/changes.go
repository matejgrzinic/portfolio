package currencies

import (
	"fmt"
	"time"

	"github.com/matejgrzinic/portfolio/external"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type changesMap map[string]map[string]map[ChangeTimeframe]float64 // todo change to uppercase?
type ChangeTimeframe string

const (
	Hour  ChangeTimeframe = "hour"
	Day   ChangeTimeframe = "day"
	Week  ChangeTimeframe = "week"
	Month ChangeTimeframe = "month"
)

type CurrencyChanges map[ChangeTimeframe]float64

func (c *Currencies) getChangesForCurrency(currencyType, symbol string) (CurrencyChanges, error) {
	chgs, ok := c.changes[currencyType][symbol]
	if !ok {
		return nil, fmt.Errorf("get changes for [type: %v] [symbol: %v]", currencyType, symbol)
	}
	out := make(CurrencyChanges)
	out[Hour] = chgs[Hour]
	out[Day] = chgs[Day]
	out[Week] = chgs[Week]
	out[Month] = chgs[Month]
	return out, nil
}

type priceData map[string]map[string]external.CurrencyData
type dbPriceData map[string]priceData

func (c *Currencies) updateChangesMap() error {
	newChangesMap := make(changesMap)
	for cType, v1 := range c.data {
		newChangesMap[cType] = make(map[string]map[ChangeTimeframe]float64)
		for cSymbol := range v1.data {
			newChangesMap[cType][cSymbol] = make(map[ChangeTimeframe]float64)
		}
	}

	timeframeData := []struct {
		timeframe ChangeTimeframe
		interval  time.Duration
	}{
		{timeframe: Hour, interval: time.Hour},
		{timeframe: Day, interval: time.Hour * 24},
		{timeframe: Week, interval: time.Hour * 24 * 7},
		{timeframe: Month, interval: time.Hour * 24 * 7 * 30},
	}

	for _, tf := range timeframeData {
		data, err := c.priceDataIntervalAgo(tf.timeframe, tf.interval)
		if err != nil {
			return err
		}
		for cType, v1 := range data {
			for cSymbol, v2 := range v1 {
				if v2.Price != 0.0 {
					cur, err := c.GetCurrency(cType, cSymbol)
					if err != nil {
						return err
					}
					newChangesMap[cType][cSymbol][tf.timeframe] = cur.Price/v2.Price - 1.0
				}
			}
		}
	}

	c.changes = newChangesMap
	return nil
}

func (c *Currencies) priceDataIntervalAgo(timeframe ChangeTimeframe, interval time.Duration) (priceData, error) {
	filter := bson.D{{Key: "time", Value: bson.M{"$gt": time.Now().Add(-interval).Unix()}}}

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{"data": 1, "_id": 0})

	dbData := make(dbPriceData)
	err := c.db.QueryRow(
		fmt.Sprintf("price data [%v] interval ago", timeframe),
		"price",
		filter,
		findOptions,
		&dbData,
	)

	if err != nil {
		return nil, err
	}

	result, ok := dbData["data"]
	if !ok {
		return nil, fmt.Errorf("data field does not exist in dbData timeframe: %v", timeframe)
	}

	return result, nil
}
