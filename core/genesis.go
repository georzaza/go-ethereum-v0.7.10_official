package core

import (
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
)

/*
 * This is the special genesis block.
 */

// Used to set the parent's hash of the genesis block
var ZeroHash256 = make([]byte, 32)
// Used to set the coinbase of the genesis block
var ZeroHash160 = make([]byte, 20)
// Used to set the bloom field of the genesis block
var ZeroHash512 = make([]byte, 64)
// Used to set the root state of the genesis block
var EmptyShaList = crypto.Sha3(ethutil.Encode([]interface{}{}))
// Used to set the tx root and receipt root of the genesis block
var EmptyListRoot = crypto.Sha3(ethutil.Encode(""))


// ZeroHash256: Previous hash (none)
//
// EmptyShaList: Empty uncles
//
// ZeroHash160: Coinbase
//
// EmptyShaList: Root state
// 
// EmptyListRoot: tx root
//
// EmptyListRoot: receipt root
//
// ZeroHash512: bloom field
// 
// big.NewInt(131072): difficulty
//
// ethutil.Big0: number of block
// 
// big.NewInt(1000000): Block upper gas bound
// 
// ethutil.Big0: Block gas used
//
// ethutil.Big0: Time field
// 
// nil: Extra field.
// 
// crypto.Sha3(big.NewInt(42).Bytes()): Nonce field
var GenesisHeader = []interface{}{
	ZeroHash256,
	EmptyShaList,
	ZeroHash160,
	EmptyShaList,
	EmptyListRoot,
	EmptyListRoot,
	ZeroHash512,
	big.NewInt(131072),
	ethutil.Big0,
	big.NewInt(1000000),
	ethutil.Big0,
	ethutil.Big0,
	nil,
	crypto.Sha3(big.NewInt(42).Bytes()),
}

// Used to be able to RLP-encode the genesis block. To rlp-encode a block we provide it's  
// header, transactions and uncles. The genesis block does not have any transactions and uncles, thus
// the use of the 2 'empty' objects.
var Genesis = []interface{}{GenesisHeader, []interface{}{}, []interface{}{}}
