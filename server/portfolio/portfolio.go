package portfolio

import (
	"context"
	"log"
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/external"
)

type Portfolio struct {
	db  *db.Database
	pd  *external.PriceData
	ctx context.Context

	refreshInterval time.Duration
}

func NewPortfolio(db *db.Database, pd *external.PriceData) *Portfolio {
	p := Portfolio{db: db, ctx: context.TODO(), pd: pd}

	p.refreshInterval = time.Minute * 10

	go p.startRefresher()
	return &p
}

func (p *Portfolio) startRefresher() {
	time.Sleep(time.Second * 10) // todo. so they have time so values are different
	p.refreshAll()
	for range time.Tick(p.refreshInterval) {
		p.refreshAll()
	}
}

func (p *Portfolio) refreshAll() {
	users := []string{"Ace"}
	for _, user := range users {
		err := p.refreshUser(user)
		if err != nil {
			log.Println(err)
		}
	}
}

// func (p *Portfolio) getAllUsers() *[]DbUserData {

// }
