package currencies

import (
	"fmt"
	"log"
	"time"

	"github.com/matejgrzinic/portfolio/external"
)

func (c *Currencies) saveToDB() error {
	data := struct {
		Data map[string]external.CurrenciesDataMap `json:"data"`
		Time int64                                 `json:"time"`
	}{
		Time: time.Now().Unix(),
	}

	data.Data = make(map[string]external.CurrenciesDataMap)
	for key, value := range c.data {
		data.Data[key] = make(external.CurrenciesDataMap)
		for key2, value2 := range value.data {
			vCpy := &external.CurrencyData{Symbol: value2.Symbol, Price: value2.Price}
			err := c.convertToEur(vCpy, key)
			if err != nil {
				log.Println(err)
				return fmt.Errorf("save to db: %v", err)
			}
			data.Data[key][key2] = *vCpy
		}
	}

	err := c.db.InsertOne(
		"currencies",
		"price",
		&data,
	)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Currencies) saveToDbAndUpdateChanges() error { // TODO test but its useless?
	if err := c.saveToDB(); err != nil {
		return err
	}
	if err := c.updateChangesMap(); err != nil {
		return err
	}
	return nil
}
