 package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSliceUniq(t *testing.T) {
	s := []string{"foo3", "foo", "foo1", "foo2", "foo", "foo1", "foo2", "foo3"}
	ret := StringsDedupAndSort(s)
	assert.Equal(t, ret, []string{"foo", "foo1", "foo2", "foo3"})
}
