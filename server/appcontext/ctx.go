package appcontext

import (
	"github.com/matejgrzinic/portfolio/currencies"
	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/portfolio"
)

type CTX interface {
	DB() db.API
	Currencies() currencies.API
	Portfolio() portfolio.API
}
