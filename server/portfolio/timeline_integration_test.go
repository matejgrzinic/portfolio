//+build integration

package portfolio

import (
	"fmt"
	"testing"

	"github.com/matejgrzinic/portfolio/currencies"
	"github.com/matejgrzinic/portfolio/db"
	"github.com/stretchr/testify/assert"
)

func TestTimeline123(t *testing.T) {

	ctx := &mockedCTX{db: db.NewDbAccess()}
	ctx.currencies = currencies.NewCurrencies(ctx)
	p := NewPortfolio(ctx)

	b, err := p.UserTimeline(&User{Name: "Ace"}, "all")

	fmt.Println(b)
	fmt.Println(err)

	assert.Fail(t, "kek")
}
