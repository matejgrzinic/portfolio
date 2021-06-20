package currencies

import (
	"fmt"
	"testing"

	"github.com/matejgrzinic/portfolio/db"
)

func Test_updateChangesMap(t *testing.T) {
	c := Currencies{dbAPI: db.NewDbAccess("portfolio2")}

	c.updateChangesMap()
	fmt.Println(c.changes)
}
