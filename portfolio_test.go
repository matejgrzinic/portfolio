package main

import (
	"testing"
)

func Test_databaseConnect(t *testing.T) {
	databaseSetup()
}

func Test_PriceAPIs(t *testing.T) {
	latestPriceData = new(priceData)
	updatePrice()

	if latestPriceData == nil {
		t.Error("latestPriceData is nil")
	}
	if len(latestPriceData.Rates["crypto"]) == 0 {
		t.Error("latestPriceData crypto is empty")
	}
	if len(latestPriceData.Rates["stock"]) == 0 {
		t.Error("latestPriceData stock is empty")
	}
	if len(latestPriceData.Rates["cash"]) == 0 {
		t.Error("latestPriceData cash is empty")
	}
	if latestPriceData.Rates["cash"]["EUR"] != 1 {
		t.Errorf("latestPriceData EUR value is not 1 and is %f", latestPriceData.Rates["cash"]["EUR"])
	}
	if latestPriceData.Rates["cash"]["USD"] >= latestPriceData.Rates["cash"]["EUR"] {
		t.Errorf("latestPriceData USD value is bigger or equad EUR and is %f", latestPriceData.Rates["cash"]["USD"])
	}
}
