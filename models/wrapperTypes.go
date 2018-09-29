package models

import (
	"github.com/onrik/ethrpc"
)

type TxHash struct {
	Timestamp int
	Hash      string
}

type Tx struct {
	Timestamp int
	Tx        *ethrpc.Transaction
}
