package types

import (
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
	"github.com/georzaza/go-ethereum-v0.7.10_official/wire"
)

// Any type that implements this interface must also implement the function Process which
// calculates and returns in that order: the total difficulty of the block as a bigInt, the
// messages of the block (aka the transactions) and/or an error.
//
// In case of an error the values returned by the function Process for the td and messages are nil.
type BlockProcessor interface {
	Process(*Block) (*big.Int, state.Messages, error)
}

// A Broadcaster is a type that can broadcast messages of a given type to a list of recipients.
type Broadcaster interface {
	Broadcast(wire.MsgType, []interface{})
}
