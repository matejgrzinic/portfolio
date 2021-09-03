//+build integration

package portfolio

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
	"unsafe"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/stretchr/testify/assert"
)

func TestPrices1(t *testing.T) {
	ctx := &mockedCTX{db: db.NewDbAccess()}

	p := NewPortfolio(ctx)

	pp, err := p.AllPrices()

	n, err := getRealSizeOf(pp)
	fmt.Println("size:", int(unsafe.Sizeof(pp))*len(pp), len(pp), n, err)
	fmt.Println(err)

	assert.Fail(t, "kek")
}

func getRealSizeOf(v interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return 0, err
	}
	return b.Len(), nil
}
