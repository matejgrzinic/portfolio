package portfolio

import (
	"log"

	db "github.com/matejgrzinic/portfolio/db/queries"
	"github.com/matejgrzinic/portfolio/external"
)

type DbBalanceData struct {
	User  string           `json:"user"`
	Data  []DbCurrencyData `json:"data"`
	Value float64          `json:"value"`
	Time  int64            `json:"time"`
}

type DbCurrencyData struct {
	CurrencyType string  `json:"type"`
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Amount       float64 `json:"amount"`
	Value        float64 `json:"value"`
}

func (p *Portfolio) refreshUser(user string) error {
	b, err := p.db.Query.GetUserBalance(user)
	if err != nil {
		return err
	}

	if err = p.RefreshBalance(b); err != nil {
		return err
	}

	if err = p.db.Query.SaveBalance(b); err != nil {
		return err
	}

	return nil
}

func (p *Portfolio) RefreshBalance(b *db.DbBalanceData) error {
	var totalValue float64
	for i := range b.Data {
		c := &b.Data[i]

		cur, err := p.pd.GetCurrency(c.CurrencyType, c.Symbol)
		if err != nil {
			log.Printf("refresh balance: %v\n", err)
			return err
		}

		c.Price = cur.Price
		c.Value = c.Price * c.Amount
		totalValue += c.Value
	}
	b.Value = totalValue

	return nil
}

// For /api/Balance
type BalanceReply struct {
	db.DbCurrencyData
	Change external.CurrencyChange
}

func (p *Portfolio) GetUserRefreshedBalance(user string) (*db.DbBalanceData, error) {
	balance, err := p.db.Query.GetUserBalance(user)
	if err != nil {
		return nil, err
	}

	if err = p.RefreshBalance(balance); err != nil {
		return nil, err
	}

	// for _, c := range balance.Data {
	// }

	return balance, nil
}
