package coll

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinEntries(t *testing.T) {
	result := JoinEntries(map[string]string{"a": "b"}, "=")
	assert.Equal(t, []string{"a=b"}, result)
}
