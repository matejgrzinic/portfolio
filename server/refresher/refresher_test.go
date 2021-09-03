package refresher

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartRefresher(t *testing.T) {
	ec := make(chan error, 1)
	sc := make(chan struct{})

	c := 0
	StartRefresher(
		ec,
		sc,
		time.Millisecond*100, // TODO somehow test duration
		func() error {
			c++
			return fmt.Errorf("unittest")
		},
	)

	err := <-ec
	assert.EqualError(t, err, "unittest")
	assert.Equal(t, 1, c)

	err = <-ec
	sc <- struct{}{}
	assert.EqualError(t, err, "unittest")
	assert.Equal(t, 2, c)
}
