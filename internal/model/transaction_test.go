package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTransactionType(t *testing.T) {
	val, _ := TransactionIn.Value()
	assert.Equal(t, "in", val)
	var txn TransactionType
	txn.Scan("out")
	assert.Equal(t, TransactionOut, txn)
}
