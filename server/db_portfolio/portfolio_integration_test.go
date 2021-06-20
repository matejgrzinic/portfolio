package db_portfolio

import (
	"testing"

	"github.com/matejgrzinic/portfolio/db"

	"github.com/stretchr/testify/assert"
)

type mockedCTX struct {
	db *db.DB
}

func (m *mockedCTX) DB() *db.DB {
	return m.db
}

func TestA(t *testing.T) {
	mockedCTX := &mockedCTX{db: db.NewDbAccess("portfolio2")}
	p := NewPortfolio(mockedCTX)

	data, err := p.GetUserBalance("Ace")
	assert.Nil(t, err)
	assert.Equal(t, "abc", data)
}
