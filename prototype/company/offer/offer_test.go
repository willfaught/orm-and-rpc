// Manual

package offer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetterName(t *testing.T) {
	var m = Offer{}
	assert.Empty(t, m.BetterName())
	m.Name = "a"
	assert.Equal(t, "aa", m.BetterName())
	m.Name = "ab"
	assert.Equal(t, "abab", m.BetterName())
}
