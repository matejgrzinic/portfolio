package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	hash, err := HashPassword("unittest")
	assert.NoError(t, err)

	assert.True(t, DoesPasswordMatchHash("unittest", hash))
	assert.False(t, DoesPasswordMatchHash("unittest wrong", hash))
}
