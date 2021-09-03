package portfolio

type Portfolio struct {
	CTX
	API
}

func NewPortfolio(ctx CTX) *Portfolio {
	return &Portfolio{CTX: ctx}
}
