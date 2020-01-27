package funks_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uol/funks"
)

/**
* The util/collections library tests.
* @author rnojiri
**/

// TestSyncMapSize - tests the function
func TestSyncMapSize(t *testing.T) {

	m := sync.Map{}

	assert.Equal(t, 0, funks.GetSyncMapSize(&m), "expected 0")

	m.Store("a", 1)

	assert.Equal(t, 1, funks.GetSyncMapSize(&m), "expected 1")

	m.Store("b", 10)

	assert.Equal(t, 2, funks.GetSyncMapSize(&m), "expected 2")

	for i := 0; i < 100; i++ {
		m.Store(strconv.Itoa(i), i)
	}

	assert.Equal(t, 102, funks.GetSyncMapSize(&m), "expected 102")

	m.Delete("b")

	assert.Equal(t, 101, funks.GetSyncMapSize(&m), "expected 101")
}
