package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
)

// Represents the receipts field of a Transaction. The receipt contains certain information regarding the
// execution of the transaction. Each receipt is placed on a trie.
//
// PostState: the state after the execution of the transaction
//
// CumulativeGasUsed: The cumulative gas used up to and including the current transaction of a Block.
//
// Bloom: A hash (2048 bits)
//
// logs: A series of log entries.  Each entry is a tuple of the logger’s address, a
// possibly empty series of 32-byte log topics, and some number of bytes of data. See state/log.go for more.
type Receipt struct {
	PostState         []byte
	CumulativeGasUsed *big.Int
	Bloom             []byte
	logs              state.Logs
}

// Creates a new Receipt object provided the 'root' of the receipts trie and the cumulative gas used.
func NewReceipt(root []byte, cumalativeGasUsed *big.Int) *Receipt {
	return &Receipt{PostState: ethutil.CopyBytes(root), CumulativeGasUsed: cumalativeGasUsed}
}

// Creates a new Receipt object from an rlp-encoded ethutil.Value object 'val'.
func NewRecieptFromValue(val *ethutil.Value) *Receipt {
	r := &Receipt{}
	r.RlpValueDecode(val)

	return r
}

// Sets the Receipt logs field equal to the 'logs' parameter.
func (self *Receipt) SetLogs(logs state.Logs) {
	self.logs = logs
}

// Sets the caller's receipt fields to the rlp-decoded fields of the decoder parameter.
func (self *Receipt) RlpValueDecode(decoder *ethutil.Value) {
	self.PostState = decoder.Get(0).Bytes()
	self.CumulativeGasUsed = decoder.Get(1).BigInt()
	self.Bloom = decoder.Get(2).Bytes()

	it := decoder.Get(3).NewIterator()
	for it.Next() {
		self.logs = append(self.logs, state.NewLogFromValue(it.Value()))
	}
}

// Returns the rlp-encodable fields of the caller.
// The only difference between what this function returns and the caller's fields is
// that this function returns the result of the RlpData() function called on the logs field.
func (self *Receipt) RlpData() interface{} {
	return []interface{}{self.PostState, self.CumulativeGasUsed, self.Bloom, self.logs.RlpData()}
}

// Rlp-encoded the rlp-encodable fields of a Receipt object. These fields are obtained through the
// function Receipt.RlpData().
func (self *Receipt) RlpEncode() []byte {
	return ethutil.Encode(self.RlpData())
}

// Returns true if the caller Receipt is the same as the 'other', false otherwise.
// Two Receipt objects are the same if their PostState fields are equal.
func (self *Receipt) Cmp(other *Receipt) bool {
	if bytes.Compare(self.PostState, other.PostState) != 0 {
		return false
	}

	return true
}

// Returns the string representation of the caller.
func (self *Receipt) String() string {
	return fmt.Sprintf("receipt{med=%x cgas=%v bloom=%x logs=%v}", self.PostState, self.CumulativeGasUsed, self.Bloom, self.logs)
}

// An array of Receipts.
type Receipts []*Receipt

// Returns the rlp-encodable fields of the caller, aka an array that contains the result
// of each of the Receipt.RlpData() calls on each Receipt that the caller consists of.
func (self Receipts) RlpData() interface{} {
	data := make([]interface{}, len(self))
	for i, receipt := range self {
		data[i] = receipt.RlpData()
	}

	return data
}

// Returns the rlp-encoding of the caller, aka rlp-encodes each Receipt object that the caller consists of and returns
// the result.
func (self Receipts) RlpEncode() []byte {
	return ethutil.Encode(self.RlpData())
}

// Returns the number of Receipts that the caller consists of.
// This function, along with the GetRlp function are implemented so as that a Receipt
// object can call the function DeriveSha (defined in the file types/derive_sha.go) which
// constructs a trie and returns the hash root of that trie.
func (self Receipts) Len() int { return len(self) }

// Returns the rlp-encoding of the i-th Receipt of the caller.
// This function, along with the Len function are implemented so as that a Receipt
// object can call the function DeriveSha (defined in the file types/derive_sha.go) which
// constructs a trie and returns the hash root of that trie.
func (self Receipts) GetRlp(i int) []byte { return ethutil.Rlp(self[i]) }
