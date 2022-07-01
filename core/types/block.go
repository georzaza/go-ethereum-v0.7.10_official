package types

//nothing

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
	"github.com/georzaza/go-ethereum-v0.7.10_official/trie"
)

// An object used to represent a Block's main info.
//
// Number: The number of the block
//
// Hash: The hash of the block
//
// Parent: The parent of the block.
//
// TD: used by the package core to store the total difficulty of the chain up to and including this block.
//
// TD(Block) = TD(Block.parent) + Block.difficulty + sum(u.difficulty for u in Block.uncles)
type BlockInfo struct {
	Number uint64
	Hash   []byte
	Parent []byte
	TD     *big.Int
}

//
// Only a BlockInfo object can call this function.
//
// Param: data: Should be the bytes representation of the RLP-encoding of the caller.
//
// Description: Sets the caller's fields to the RLP-decoded parts of the 'data' parameter. To do so,
// the parameter `data` is converted to an ethutil.Value object and the RLP-decoding operation happens on the Value object.
//
// More information about the ethutil.Value object can be found in the file ethutil/README.md
func (bi *BlockInfo) RlpDecode(data []byte) {
	decoder := ethutil.NewValueFromBytes(data)

	bi.Number = decoder.Get(0).Uint()
	bi.Hash = decoder.Get(1).Bytes()
	bi.Parent = decoder.Get(2).Bytes()
	bi.TD = decoder.Get(3).BigInt()
}

// Returns the rlp-encoded object of a BlockInfo object by calling the function Encode defined in ethutil/rlp.go
func (bi *BlockInfo) RlpEncode() []byte {
	return ethutil.Encode([]interface{}{bi.Number, bi.Hash, bi.Parent, bi.TD})
}

type Blocks []*Block

// Returns Blocks as a set where the elements of the set are the hashes of the Blocks.
func (self Blocks) AsSet() ethutil.UniqueSet {
	set := make(ethutil.UniqueSet)
	for _, block := range self {
		set.Insert(block.Hash())
	}

	return set
}

// Used for sorting Blocks
type BlockBy func(b1, b2 *Block) bool

// Sorts Blocks using a blockSorter object
func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

// This data type is used for sorting Blocks.
type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

// Returns the number of the Block objects that a blockSorter object consists of.
func (self blockSorter) Len() int { return len(self.blocks) }

// Swaps two Blocks of a blockSorter object.
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}

// May be called only on blockSorter objects. For 2 Blocks i, j, returns true if Block i is
// less than j. To determine whether i is less than j, a BlockBy object is used.
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

// Returns true if the block number of b1 is less than b2, false otherwise.
func Number(b1, b2 *Block) bool { return b1.Number.Cmp(b2.Number) < 0 }

// PrevHash: The Keccak 256-bit hash of the parent block’s header, in its entirety;
//
// Uncles : The uncles of the Block.
//
// UncleSha: The Keccak 256-bit hash of the uncles of the Block.
//
// CoinBase: The address of the beneficiary.
//
// beneficiary: The 160-bit address to which all fees collected from the successful mining of this block are transferred.
//
// state: The Keccak 256-bit hash of the root node of the state trie, after all transactions are executed
//
// Difficulty : Difficulty of the current block.
//
// Time: Creation time of the Block.
//
// Number: The number of the Block. This is equal to the number of ancestor blocks.
//
// GasLimit: The maximum gas limit all the transactions inside this Block are allowed to consume.
//
// GasUsed: The total gas used by the transactions of the Block.
//
// Extra: An arbitrary byte array containing data relevant to this block. According to the Ethereum
// yellow paper this should be 32 bytes or less.
//
// Nonce: The Block nonce. This is what miners keep changing to compute a solution to PoW.
//
// transactions: List of transactions and/or contracts to be created included in this Block.
//
// receipts: The receipts of the transactions
//
// TxSha: The Keccak 256-bit hash of the root node of the transactions trie of the block
//
// ReceiptSha: The Keccak 256-bit hash of the root node of the receipts trie of the block.
//
// LogsBloom: The Bloom filter composed from indexable information (logger address and log topics)
// contained in each log entry from the receipt of each transaction in the transactions list
//
// Reward: The reward of the beneficiary (miner)
type Block struct {
	PrevHash          ethutil.Bytes
	Uncles            Blocks
	UncleSha          []byte
	Coinbase          []byte
	state             *state.StateDB
	Difficulty        *big.Int
	Time              int64
	Number            *big.Int
	GasLimit          *big.Int
	GasUsed           *big.Int
	Extra             string
	Nonce             ethutil.Bytes
	transactions      Transactions
	receipts          Receipts
	TxSha, ReceiptSha []byte
	LogsBloom         []byte
	Reward            *big.Int
}

// Creates a new Block from raw bytes.
// These bytes should be the bytes representation of result of the RLP-encoding of a Block.
func NewBlockFromBytes(raw []byte) *Block {
	block := &Block{}
	block.RlpDecode(raw)

	return block
}

// Creates a new Block from the rlpValue object.
// Only the block's header, transactions and uncles are derived from the rlpValue object.
func NewBlockFromRlpValue(rlpValue *ethutil.Value) *Block {
	block := &Block{}
	block.RlpValueDecode(rlpValue)

	return block
}

// Creates a Block. See the Block data structure for an explanation of the fields used.
//
// The root parameter is the root of the block's state trie.
//
// The Block created will be created at the current Unix time.
//
// Notice that there are no transactions and receipts passed as arguments to this function.
func CreateBlock(root interface{},
	prevHash []byte,
	base []byte,
	Difficulty *big.Int,
	Nonce []byte,
	extra string) *Block {

	block := &Block{
		PrevHash:   prevHash,
		Coinbase:   base,
		Difficulty: Difficulty,
		Nonce:      Nonce,
		Time:       time.Now().Unix(),
		Extra:      extra,
		UncleSha:   nil,
		GasUsed:    new(big.Int),
		GasLimit:   new(big.Int),
	}
	block.SetUncles([]*Block{})

	block.state = state.New(trie.New(ethutil.Config.Db, root))

	return block
}

// Returns the block's hash. To do so, the block's header is first converted to an ethutil.Value object
// and then the Encode function is called upon the latter.
func (block *Block) Hash() ethutil.Bytes {
	return crypto.Sha3(ethutil.NewValue(block.header()).Encode())
	//return crypto.Sha3(block.Value().Encode())
}

// Returns the hash of an object that is almost the same as a Block. The differences are:
//
// 1. The object to be hashed contains only the uncles hash.
//
// 2. The object to be hashed contains the root of the state and not eh state as a whole.
//
// 3. The object to be hashed contains only the TxSha and not the receipts.
//
// 4. The object to be hashed contains only the ReceiptSha and not the transactions.
//
// 5. The object to be hashed does not contain the Reward of the miner.
//
// Note: The object to be hashed if appended with the Nonce field of the Block will comprise the Block's header.
func (block *Block) HashNoNonce() []byte {
	return crypto.Sha3(ethutil.Encode(block.miningHeader()))
}

// Returns the state of the Block. (not just the root)
func (block *Block) State() *state.StateDB {
	return block.state
}

// Returns the transactions of the block.
func (block *Block) Transactions() Transactions {
	return block.transactions
}

// Calculates the gas limit.
//
// If the Block passed as a parameter is the genesis block the gas limit is set to 10^6.
//
// Otherwise the gas limit will be ~= 1023 * parent.GasLimit + parent.GasUsed*6/5
//
// The minimum gas limit is set to 125000.
func (block *Block) CalcGasLimit(parent *Block) *big.Int {
	if block.Number.Cmp(big.NewInt(0)) == 0 {
		return ethutil.BigPow(10, 6)
	}

	// ((1024-1) * parent.gasLimit + (gasUsed * 6 / 5)) / 1024

	previous := new(big.Int).Mul(big.NewInt(1024-1), parent.GasLimit)
	current := new(big.Rat).Mul(new(big.Rat).SetInt(parent.GasUsed), big.NewRat(6, 5))
	curInt := new(big.Int).Div(current.Num(), current.Denom())

	result := new(big.Int).Add(previous, curInt)
	result.Div(result, big.NewInt(1024))

	min := big.NewInt(125000)

	return ethutil.BigMax(min, result)
}

// Returns the BlockInfo representation of a Block.
func (block *Block) BlockInfo() BlockInfo {
	bi := BlockInfo{}
	data, _ := ethutil.Config.Db.Get(append(block.Hash(), []byte("Info")...))
	bi.RlpDecode(data)

	return bi
}

// Searches for a transaction in the current Block based on the transactions hash
// and the hash parameter and returns it if it exists. Otherwise returns nil.
func (self *Block) GetTransaction(hash []byte) *Transaction {
	for _, tx := range self.transactions {
		if bytes.Compare(tx.Hash(), hash) == 0 {
			return tx
		}
	}

	return nil
}

// Sync the block's state and contract respectively.
// For more, see the state package of this go-ethereum version.
func (block *Block) Sync() {
	block.state.Sync()
}

// Resets the state to nil.
func (block *Block) Undo() {
	// Sync the block state itself
	block.state.Reset()
}

// Returns the receipts as a string.
// This function is used to derive the ReceiptSha of the receipts.
func (block *Block) rlpReceipts() interface{} {
	// Marshal the transactions of this block
	encR := make([]interface{}, len(block.receipts))
	for i, r := range block.receipts {
		// Cast it to a string (safe)
		encR[i] = r.RlpData()
	}

	return encR
}

// Returns the uncles as a string. This function is used to derive the UncleSha of the uncles.
func (block *Block) rlpUncles() interface{} {
	// Marshal the transactions of this block
	uncles := make([]interface{}, len(block.Uncles))
	for i, uncle := range block.Uncles {
		// Cast it to a string (safe)
		uncles[i] = uncle.header()
	}

	return uncles
}

// Sets the uncles of a Block to the parameter 'uncles'. This function is also called by Block.CreateBlock.
// Also sets the UncleSha of the Block based on the provided parameter. To do so, an inner function called rlpUncles() is used.
func (block *Block) SetUncles(uncles []*Block) {
	block.Uncles = uncles
	block.UncleSha = crypto.Sha3(ethutil.Encode(block.rlpUncles()))
}

// Sets the receipts of a Block to the parameter 'receipts'.
// Also calculates and sets the LogsBloom field which is derived by the receipts.
func (self *Block) SetReceipts(receipts Receipts) {
	self.receipts = receipts
	self.ReceiptSha = DeriveSha(receipts)
	self.LogsBloom = CreateBloom(receipts)
}

// Sets the transactions of a Block to the parameter 'transactions'.
// Also calculates and sets the TxSha field which is derived by the transactions.
func (self *Block) SetTransactions(txs Transactions) {
	self.transactions = txs
	self.TxSha = DeriveSha(txs)
}

// Casts a block to an ethutil.Value object containing the header,
// the transactions and the uncles of the block and then returns it.
func (block *Block) Value() *ethutil.Value {
	return ethutil.NewValue([]interface{}{block.header(), block.transactions, block.rlpUncles()})
}

// Calls Block.Value() function on the current Block and then rlp-encodes the Block.
// The rlp-encodable fields of a Block are it's header, transactions and uncles.
func (block *Block) RlpEncode() []byte {
	return block.Value().Encode()
}

// RLP-decodes a Block. To do so, a new ethutil.Value object is created from the 'data' parameter
// and then the function Block.RlpValueDecode(data) is called on the current Block.
// The 'data' parameter represents the rlp-encodable fields of a Block and can be obtained by using the function Block.RlpData().
func (block *Block) RlpDecode(data []byte) {
	rlpValue := ethutil.NewValueFromBytes(data)
	block.RlpValueDecode(rlpValue)
}

// RLP-decodes the rlp-encodable fields of a Block, aka the header, transactions and uncles.
// parameter decoder: The above fields of the Block after having been cast to an ethutil.Value object.
func (block *Block) RlpValueDecode(decoder *ethutil.Value) {
	block.setHeader(decoder.Get(0))

	// Tx list might be empty if this is an uncle. Uncles only have their
	// header set.
	if decoder.Get(1).IsNil() == false { // Yes explicitness
		txs := decoder.Get(1)
		block.transactions = make(Transactions, txs.Len())
		for i := 0; i < txs.Len(); i++ {
			block.transactions[i] = NewTransactionFromValue(txs.Get(i))
		}
	}

	if decoder.Get(2).IsNil() == false { // Yes explicitness
		uncles := decoder.Get(2)
		block.Uncles = make([]*Block, uncles.Len())
		for i := 0; i < uncles.Len(); i++ {
			block.Uncles[i] = NewUncleBlockFromValue(uncles.Get(i))
		}
	}
}

// Sets the header of the block given an ethutil.Value object that contains that information.
func (self *Block) setHeader(header *ethutil.Value) {
	self.PrevHash = header.Get(0).Bytes()
	self.UncleSha = header.Get(1).Bytes()
	self.Coinbase = header.Get(2).Bytes()
	self.state = state.New(trie.New(ethutil.Config.Db, header.Get(3).Val))
	self.TxSha = header.Get(4).Bytes()
	self.ReceiptSha = header.Get(5).Bytes()
	self.LogsBloom = header.Get(6).Bytes()
	self.Difficulty = header.Get(7).BigInt()
	self.Number = header.Get(8).BigInt()
	self.GasLimit = header.Get(9).BigInt()
	self.GasUsed = header.Get(10).BigInt()
	self.Time = int64(header.Get(11).BigInt().Uint64())
	self.Extra = header.Get(12).Str()
	self.Nonce = header.Get(13).Bytes()
}

// Creates and sets the uncle of a block based on an ethutil.Value object that contains that information.
//
// 'header' parameter: an ethutil.Value object containing the header of the uncle block.
//
// The header of the uncle Block (and any Block) contains those fields: PrevHash, UncleSha, Coinbase, state, TxSha, ReceiptSha,
// LogsBloom, Difficulty, Number, GasLimit, GasUsed, Time, Extra, Nonce. See also the function Block.HashNoNonce() for details about these fields.
func NewUncleBlockFromValue(header *ethutil.Value) *Block {
	block := &Block{}
	block.setHeader(header)

	return block
}

// Returns the block's state trie
func (block *Block) Trie() *trie.Trie {
	return block.state.Trie
}

// Returns the block's state root
func (block *Block) Root() interface{} {
	return block.state.Root()
}

// Returns the block difficulty.
func (block *Block) Diff() *big.Int {
	return block.Difficulty
}

// Returns the block receipts.
func (self *Block) Receipts() []*Receipt {
	return self.receipts
}

// Returns an object that is almost the same as a Block. The differences are:
// 1. The Block struct also contains the uncles as a slice of Blocks and not only the uncles hash.
// 2. The Block struct contains the state object as a whole and not the root of the state.
// 3. The Block struct contains the receipts and not only the TxSha
// 4. The Block struct contains the transactions and not only the ReceiptSha
// 5. The Block struct contains the Reward to be given to the miner.
// The object returned, if appended with the Nonce field of the Block is the Block's header.
func (block *Block) miningHeader() []interface{} {
	return []interface{}{
		// Sha of the previous block
		block.PrevHash,
		// Sha of uncles
		block.UncleSha,
		// Coinbase address
		block.Coinbase,
		// root state
		block.Root(),
		// tx root
		block.TxSha,
		// Sha of tx
		block.ReceiptSha,
		// Bloom
		block.LogsBloom,
		// Current block Difficulty
		block.Difficulty,
		// The block number
		block.Number,
		// Block upper gas bound
		block.GasLimit,
		// Block gas used
		block.GasUsed,
		// Time the block was found?
		block.Time,
		// Extra data
		block.Extra,
	}
}

// Returns the header of the Block. See the function Block.miningHeader() for a detailed explanation.
func (block *Block) header() []interface{} {
	return append(block.miningHeader(), block.Nonce)
}

// Returns the string representation of a Block object.
func (block *Block) String() string {
	return fmt.Sprintf(`
	BLOCK(%x): Size: %v
	PrevHash:   %x
	UncleSha:   %x
	Coinbase:   %x
	Root:       %x
	TxSha       %x
	ReceiptSha: %x
	Bloom:      %x
	Difficulty: %v
	Number:     %v
	MaxLimit:   %v
	GasUsed:    %v
	Time:       %v
	Extra:      %v
	Nonce:      %x
	NumTx:      %v
`,
		block.Hash(),
		block.Size(),
		block.PrevHash,
		block.UncleSha,
		block.Coinbase,
		block.Root(),
		block.TxSha,
		block.ReceiptSha,
		block.LogsBloom,
		block.Difficulty,
		block.Number,
		block.GasLimit,
		block.GasUsed,
		block.Time,
		block.Extra,
		block.Nonce,
		len(block.transactions),
	)
}

// Returns a float64 object representing the size of storage needed to save this block's rlp-encoding.
// The rlp-encodable fields of a Block are it's header, transactions and uncles.
func (self *Block) Size() ethutil.StorageSize {
	return ethutil.StorageSize(len(self.RlpEncode()))
}

// Returns an object containing the fields of the caller (Block) that can be rlp-encoded.
func (self *Block) RlpData() interface{} {
	return []interface{}{self.header(), self.transactions, self.rlpUncles()}
}

// Returns the Nonce field of the caller (Block)
func (self *Block) N() []byte { return self.Nonce }
