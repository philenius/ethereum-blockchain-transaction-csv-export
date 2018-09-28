package main

import (
	"github.com/onrik/ethrpc"
)

type TxHash struct {
	timestamp int
	hash      string
}

type Tx struct {
	timestamp int
	tx        *ethrpc.Transaction
}
