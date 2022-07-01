package types

import (
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
)

// Creates the LogsBloom of a block.
// That is, the union of the LogsBloom of all the receipts in the block.
// Each block has many receipts. Each receipt has many logs (each log
// has 3 fields: the address, the topics and the data). The LogsBloom of
// a receipt will be the lower 11 bits of the modulo operation by 2048
// of the kecchak(address and logs of the receipt) modulo 2048.
func CreateBloom(receipts Receipts) []byte {
	bin := new(big.Int)
	for _, receipt := range receipts {
		bin.Or(bin, LogsBloom(receipt.logs))
	}

	return ethutil.LeftPadBytes(bin.Bytes(), 64)
}

// Returns the LogBloom of a receipt.
func LogsBloom(logs state.Logs) *big.Int {
	bin := new(big.Int)
	for _, log := range logs {
		data := make([][]byte, len(log.Topics())+1)
		data[0] = log.Address()

		for i, topic := range log.Topics() {
			data[i+1] = topic
		}

		for _, b := range data {
			bin.Or(bin, ethutil.BigD(bloom9(crypto.Sha3(b)).Bytes()))
		}
	}

	return bin
}

// todo improve
// Get the lower 11 bits of kecchak(b) % 2048
func bloom9(b []byte) *big.Int {
	r := new(big.Int)
	for _, i := range []int{0, 2, 4} {
		t := big.NewInt(1)
		b := uint(b[i+1]) + 256*(uint(b[i])&1)
		r.Or(r, t.Lsh(t, b))
	}

	return r
}

// This function is called by the function "bloomFilter" defined in the file core/filter.go
//
func BloomLookup(bin, topic []byte) bool {
	bloom := ethutil.BigD(bin)
	cmp := bloom9(crypto.Sha3(topic))

	return bloom.And(bloom, cmp).Cmp(cmp) == 0
}
