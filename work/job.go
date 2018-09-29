package work

import (
	"github.com/onrik/ethrpc"
)

type Job struct {
	BlockHeight int
	TxHash      string
	Tx          *ethrpc.Transaction
	Timestamp   int
}
