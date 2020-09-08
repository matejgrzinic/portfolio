package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type crypto []struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type usdrates struct {
	Rates struct {
		CAD float64 `json:"CAD"`
		HKD float64 `json:"HKD"`
		ISK float64 `json:"ISK"`
		PHP float64 `json:"PHP"`
		DKK float64 `json:"DKK"`
		HUF float64 `json:"HUF"`
		CZK float64 `json:"CZK"`
		GBP float64 `json:"GBP"`
		RON float64 `json:"RON"`
		SEK float64 `json:"SEK"`
		IDR float64 `json:"IDR"`
		INR float64 `json:"INR"`
		BRL float64 `json:"BRL"`
		RUB float64 `json:"RUB"`
		HRK float64 `json:"HRK"`
		JPY float64 `json:"JPY"`
		THB float64 `json:"THB"`
		CHF float64 `json:"CHF"`
		EUR float64 `json:"EUR"`
		MYR float64 `json:"MYR"`
		BGN float64 `json:"BGN"`
		TRY float64 `json:"TRY"`
		CNY float64 `json:"CNY"`
		NOK float64 `json:"NOK"`
		NZD float64 `json:"NZD"`
		ZAR float64 `json:"ZAR"`
		USD float64 `json:"USD"`
		MXN float64 `json:"MXN"`
		SGD float64 `json:"SGD"`
		AUD float64 `json:"AUD"`
		ILS float64 `json:"ILS"`
		KRW float64 `json:"KRW"`
		PLN float64 `json:"PLN"`
	} `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}

type priceData struct {
	Rates map[string]map[string]float64
	mux   sync.Mutex
}

var latestPriceData *priceData

func startUpdatePriceInterval() {
	latestPriceData = new(priceData)
	latestPriceData.Rates = make(map[string]map[string]float64)
	latestPriceData.Rates["crypto"] = make(map[string]float64)
	latestPriceData.Rates["cash"] = make(map[string]float64)
	latestPriceData.Rates["stock"] = make(map[string]float64)

	updatePrice()
	for range time.Tick(time.Minute) {
		updatePrice()
	}
}

func updatePrice() {
	latestPriceData.mux.Lock()
	defer latestPriceData.mux.Unlock()

	var wg sync.WaitGroup

	wg.Add(1)
	go getCryptoPrices(&wg)
	wg.Add(1)
	go getUSD(&wg)
	wg.Wait()

	insertPrice()
}

func getCryptoPrices(wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var c crypto
	err = json.Unmarshal(data, &c)

	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range c {
		if strings.HasSuffix(v.Symbol, "USDT") {
			fValue, err := strconv.ParseFloat(v.Price, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			latestPriceData.Rates["crypto"][v.Symbol[:len(v.Symbol)-4]] = fValue
		}
	}
}

func getUSD(wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get("https://api.exchangeratesapi.io/latest?base=USD")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var c usdrates
	err = json.Unmarshal(data, &c)

	if err != nil {
		log.Println(err)
		return
	}

	v := reflect.ValueOf(c.Rates)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		latestPriceData.Rates["cash"][typeOfS.Field(i).Name] = v.Field(i).Float()
	}
}
