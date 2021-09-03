package portfolio

import "github.com/matejgrzinic/portfolio/currencies"

type CurrencyBalance struct {
	CurrencyType string                     `json:"currencytype"`
	Symbol       string                     `json:"symbol"`
	Price        float64                    `json:"price"`
	Amount       float64                    `json:"amount"`
	Value        float64                    `json:"value"`
	Changes      currencies.CurrencyChanges `json:"changes"`
}

func (p *Portfolio) UserBalance(user *User) ([]CurrencyBalance, error) {
	hashData := make(map[string]map[string]float64)

	ts, err := p.AllUserTransactions(user)
	if err != nil {
		return nil, err
	}

	addToHashDataLoop := func(tc []TransactionCurrency, multiplier float64) {
		for _, e := range tc {
			if _, ok := hashData[e.CurrencyType]; !ok {
				hashData[e.CurrencyType] = make(map[string]float64)
			}
			hashData[e.CurrencyType][e.Symbol] += e.Amount * multiplier
		}
	}

	for _, t := range ts {
		addToHashDataLoop(t.Gains, 1.0)
		addToHashDataLoop(t.Losses, -1.0)
	}

	result := make([]CurrencyBalance, 0)
	for currType, curs := range hashData {
		for symbol, amount := range curs {
			cwc, err := p.Currencies().GetCurrencyWithChanges(currType, symbol)
			if err != nil {
				return nil, err
			}
			cp := CurrencyBalance{
				CurrencyType: currType,
				Symbol:       symbol,
				Amount:       amount,
				Price:        cwc.Price,
				Value:        cwc.Price * amount,
				Changes:      cwc.Changes,
			}

			result = append(result, cp)
		}
	}

	return result, nil
}
