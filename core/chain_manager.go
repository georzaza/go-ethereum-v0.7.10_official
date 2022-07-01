package core

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/georzaza/go-ethereum-v0.7.10_official/core/types"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/event"
	"github.com/georzaza/go-ethereum-v0.7.10_official/logger"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
)

// Logger channel of the chain
var chainlogger = logger.NewLogger("CHAIN")

// Sets the following accounts with a balance of 1606938044258990275541962092341162602522202 Ether for testing.
//
//51ba59315b3a95761d0863b05ccc7a7f54703d99
//
//e4157b34ea9615cfbde6b4fda419828124b70c78
//
//b9c015918bdaba24b4ff057a92a3873d6eb201be
//
//6c386a4b26f73c802f34673f7248bb118f97424a
//
//cd2a3d9f938e13cd947ec05abc7fe734df8dd826
//
//2ef47100e0787b915105fd5e3f4ff6752079d5cb
//
//e6716f9544a56c530d868e4bfbacb172315bdead
//
//1a26338f0d905e295fccb71fa9ea849ffa12aaf4
func AddTestNetFunds(block *types.Block) {
	for _, addr := range []string{
		"51ba59315b3a95761d0863b05ccc7a7f54703d99",
		"e4157b34ea9615cfbde6b4fda419828124b70c78",
		"b9c015918bdaba24b4ff057a92a3873d6eb201be",
		"6c386a4b26f73c802f34673f7248bb118f97424a",
		"cd2a3d9f938e13cd947ec05abc7fe734df8dd826",
		"2ef47100e0787b915105fd5e3f4ff6752079d5cb",
		"e6716f9544a56c530d868e4bfbacb172315bdead",
		"1a26338f0d905e295fccb71fa9ea849ffa12aaf4",
	} {
		codedAddr := ethutil.Hex2Bytes(addr)
		account := block.State().GetAccount(codedAddr)
		account.SetBalance(ethutil.Big("1606938044258990275541962092341162602522202993782792835301376")) //ethutil.BigPow(2, 200)
		block.State().UpdateStateObject(account)
	}
}

// Calculates the difficulty of a block and returns it.
// If the block was mined in less than 5 seconds, the difficulty of the block is increased by 1/1024th of the parent's
// difficulty. If the block was mined in more than 5 seconds, the difficulty is decreased by 1/1024th of
// the parent's difficulty.
func CalcDifficulty(block, parent *types.Block) *big.Int {
	diff := new(big.Int)

	adjust := new(big.Int).Rsh(parent.Difficulty, 10)
	if block.Time >= parent.Time+5 {
		diff.Sub(parent.Difficulty, adjust)
	} else {
		diff.Add(parent.Difficulty, adjust)
	}

	return diff
}

// ChainManager is mainly used for the creation of genesis or any other blocks
//
// processor: An interface. A neat way of calling the function Process of a BlockManager object.
//
// eventMux: Used to dispatch events to subscribers.
//
// genesisBlock: The special genesis block.
//
// mu: a mutex for the ChainManager object.
//
// td: the total difficulty. TD(genesis_block) = 0 and TD(B) = TD(B.parent) + sum(u.difficulty for u in B.uncles) + B.difficulty
//
// lastBlockNumber: the last block's number. (the last successfully inserted block on the chain)
//
// currentBlock: During the creation of a new block, the currentBlock will point to the parent of the block to be created.
//
// lastBlockHash: the last block's hash.
//
// transState: represents the world state.
type ChainManager struct {
	processor       types.BlockProcessor
	eventMux        *event.TypeMux
	genesisBlock    *types.Block
	mu              sync.RWMutex
	td              *big.Int
	lastBlockNumber uint64
	currentBlock    *types.Block
	lastBlockHash   []byte
	transState      *state.StateDB
}

// Returns the total difficulty.
//
// TD(genesis_block) = 0 and TD(B) = TD(B.parent) + sum(u.difficulty for u in B.uncles) + B.difficulty
func (self *ChainManager) Td() *big.Int {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.td
}

// Returns the last block number.
func (self *ChainManager) LastBlockNumber() uint64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.lastBlockNumber
}

// Returns the last block's hash
func (self *ChainManager) LastBlockHash() []byte {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.lastBlockHash
}

// Returns the current block todo
func (self *ChainManager) CurrentBlock() *types.Block {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.currentBlock
}

// Creates and returns a new ChainManager object by setting the genesisBlock and the eventMux field of the ChainManager.
func NewChainManager(mux *event.TypeMux) *ChainManager {
	bc := &ChainManager{}
	bc.genesisBlock = types.NewBlockFromBytes(ethutil.Encode(Genesis))
	bc.eventMux = mux

	bc.setLastBlock()

	bc.transState = bc.State().Copy()

	return bc
}

// Sets the processor field of the ChainManager object
func (self *ChainManager) SetProcessor(proc types.BlockProcessor) {
	self.processor = proc
}

// Returns the world state as 'seen' by the current block.
func (self *ChainManager) State() *state.StateDB {
	return self.CurrentBlock().State()
}

// Returns the world state.
func (self *ChainManager) TransState() *state.StateDB {
	return self.transState
}

// An inner function, used by the ChainManager 'constructor' function that sets the last block of the ChainManager.
// If the chain has 0 blocks so far, it makes a call to the function AddTestNetFunds(genesisBlock).
func (bc *ChainManager) setLastBlock() {
	data, _ := ethutil.Config.Db.Get([]byte("LastBlock"))
	if len(data) != 0 {
		// Prep genesis
		AddTestNetFunds(bc.genesisBlock)

		block := types.NewBlockFromBytes(data)
		bc.currentBlock = block
		bc.lastBlockHash = block.Hash()
		bc.lastBlockNumber = block.Number.Uint64()

		// Set the last know difficulty (might be 0x0 as initial value, Genesis)
		bc.td = ethutil.BigD(ethutil.Config.Db.LastKnownTD())
	} else {
		bc.Reset()
	}

	chainlogger.Infof("Last block (#%d) %x\n", bc.lastBlockNumber, bc.currentBlock.Hash())
}

// Creates a new Block by making a call to the function CreateBlock of the types package, sets
// it's difficulty, number and gaslimit and returns it.
func (bc *ChainManager) NewBlock(coinbase []byte) *types.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var root interface{}
	hash := ZeroHash256

	if bc.CurrentBlock != nil {
		root = bc.currentBlock.Root()
		hash = bc.lastBlockHash
	}

	block := types.CreateBlock(
		root,
		hash,
		coinbase,
		ethutil.BigPow(2, 32),
		nil,
		"")

	parent := bc.currentBlock
	if parent != nil {
		block.Difficulty = CalcDifficulty(block, parent)
		block.Number = new(big.Int).Add(bc.currentBlock.Number, ethutil.Big1)
		block.GasLimit = block.CalcGasLimit(bc.currentBlock)

	}

	return block
}

// Resets the chain to the point where the chain will only contain the genesis block. This includes the call on the
// function AddTestNetFunds.
func (bc *ChainManager) Reset() {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	AddTestNetFunds(bc.genesisBlock)
	bc.genesisBlock.Trie().Sync()
	bc.write(bc.genesisBlock)
	bc.insert(bc.genesisBlock)
	bc.currentBlock = bc.genesisBlock

	bc.setTotalDifficulty(ethutil.Big("0"))

	// Set the last know difficulty (might be 0x0 as initial value, Genesis)
	bc.td = ethutil.BigD(ethutil.Config.Db.LastKnownTD())
}

// Returns the RLP-encoding of all blocks of the chain.
func (self *ChainManager) Export() []byte {
	self.mu.RLock()
	defer self.mu.RUnlock()

	chainlogger.Infof("exporting %v blocks...\n", self.currentBlock.Number)

	blocks := make([]*types.Block, int(self.currentBlock.Number.Int64())+1)
	for block := self.currentBlock; block != nil; block = self.GetBlock(block.PrevHash) {
		blocks[block.Number.Int64()] = block
	}

	return ethutil.Encode(blocks)
}

// Inner function, used to insert a block on the chain. What actually gets inserted into the chain is
// the block's rlp-encoding.
func (bc *ChainManager) insert(block *types.Block) {
	encodedBlock := block.RlpEncode()
	ethutil.Config.Db.Put([]byte("LastBlock"), encodedBlock)
	bc.currentBlock = block
	bc.lastBlockHash = block.Hash()
}

// Inner function, used to write a block to the database
func (bc *ChainManager) write(block *types.Block) {
	bc.writeBlockInfo(block)

	encodedBlock := block.RlpEncode()
	ethutil.Config.Db.Put(block.Hash(), encodedBlock)
}

// Returns the genesis block.
func (bc *ChainManager) Genesis() *types.Block {
	return bc.genesisBlock
}

// Returns whether the given hash param matches a block's hash present in the chain.
func (bc *ChainManager) HasBlock(hash []byte) bool {
	data, _ := ethutil.Config.Db.Get(hash)
	return len(data) != 0
}

// Returns a list of hashes of the chain starting from the genesis hash and up to, but not including, the 'max' block.
func (self *ChainManager) GetChainHashesFromHash(hash []byte, max uint64) (chain [][]byte) {
	block := self.GetBlock(hash)
	if block == nil {
		return
	}

	// XXX Could be optimised by using a different database which only holds hashes (i.e., linked list)
	for i := uint64(0); i < max; i++ {
		chain = append(chain, block.Hash())

		if block.Number.Cmp(ethutil.Big0) <= 0 {
			break
		}

		block = self.GetBlock(block.PrevHash)
	}

	return
}

// Returns the block of the chain that has the given hash.
func (self *ChainManager) GetBlock(hash []byte) *types.Block {
	data, _ := ethutil.Config.Db.Get(hash)
	if len(data) == 0 {
		return nil
	}

	return types.NewBlockFromBytes(data)
}

// Returns the block of the chain that has the given num .
func (self *ChainManager) GetBlockByNumber(num uint64) *types.Block {
	self.mu.RLock()
	defer self.mu.RUnlock()

	block := self.currentBlock
	for ; block != nil; block = self.GetBlock(block.PrevHash) {
		if block.Number.Uint64() == num {
			break
		}
	}

	if block != nil && block.Number.Uint64() == 0 && num != 0 {
		return nil
	}

	return block
}

// Sets the total difficulty of the ChainManager object. Also stores that information on the database.
func (bc *ChainManager) setTotalDifficulty(td *big.Int) {
	ethutil.Config.Db.Put([]byte("LTD"), td.Bytes())
	bc.td = td
}

// Calculates the total difficulty of the ChainManager and returns it in a tuple (td, nil). If an error
// occured, then the tuple (nil, error) is returned.
//
// TD(genesis_block) = 0 and TD(B) = TD(B.parent) + sum(u.difficulty for u in B.uncles) + B.difficulty
func (self *ChainManager) CalcTotalDiff(block *types.Block) (*big.Int, error) {
	parent := self.GetBlock(block.PrevHash)
	if parent == nil {
		return nil, fmt.Errorf("Unable to calculate total diff without known parent %x", block.PrevHash)
	}

	parentTd := parent.BlockInfo().TD

	uncleDiff := new(big.Int)
	for _, uncle := range block.Uncles {
		uncleDiff = uncleDiff.Add(uncleDiff, uncle.Difficulty)
	}

	td := new(big.Int)
	td = td.Add(parentTd, uncleDiff)
	td = td.Add(td, block.Difficulty)

	return td, nil
}

// Returns the block's BlockInfo object representation. See types.BlockInfo for more.
func (bc *ChainManager) BlockInfo(block *types.Block) types.BlockInfo {
	bi := types.BlockInfo{}
	data, _ := ethutil.Config.Db.Get(append(block.Hash(), []byte("Info")...))
	bi.RlpDecode(data)

	return bi
}

// Inner function for writing extra non-essential block info to the database.
func (bc *ChainManager) writeBlockInfo(block *types.Block) {
	bc.lastBlockNumber++
	bi := types.BlockInfo{Number: bc.lastBlockNumber, Hash: block.Hash(), Parent: block.PrevHash, TD: bc.td}

	// For now we use the block hash with the words "info" appended as key
	ethutil.Config.Db.Put(append(block.Hash(), []byte("Info")...), bi.RlpEncode())
}

// Sends a stop message to the chain logger channel if and only if the currentBlock field
// of the ChainManager is not nil.
func (bc *ChainManager) Stop() {
	if bc.CurrentBlock != nil {
		chainlogger.Infoln("Stopped")
	}
}

// This function iterates over the blocks in the chain param and does the following:
//
// 1. It calls the `Process` method of the `BlockProcessor` interface.
//
// 2. writes the block to the database.
//
// 4. sets the total difficulty of the block.
//
// 5. inserts the block into the chain.
//
// 6. posts a `NewBlockEvent` to the event mux.
//
// 7. posts the messages to the event mux.
//
// Returns: either nil for success or an error.
func (self *ChainManager) InsertChain(chain types.Blocks) error {
	for _, block := range chain {
		td, messages, err := self.processor.Process(block)
		if err != nil {
			if IsKnownBlockErr(err) {
				continue
			}

			chainlogger.Infof("block #%v process failed (%x)\n", block.Number, block.Hash()[:4])
			chainlogger.Infoln(block)
			chainlogger.Infoln(err)
			return err
		}

		self.mu.Lock()
		{
			self.write(block)
			if td.Cmp(self.td) > 0 {
				if block.Number.Cmp(new(big.Int).Add(self.currentBlock.Number, ethutil.Big1)) < 0 {
					chainlogger.Infof("Split detected. New head #%v (%x), was #%v (%x)\n", block.Number, block.Hash()[:4], self.currentBlock.Number, self.currentBlock.Hash()[:4])
				}

				self.setTotalDifficulty(td)
				self.insert(block)
				self.transState = self.currentBlock.State().Copy()
			}

		}
		self.mu.Unlock()

		self.eventMux.Post(NewBlockEvent{block})
		self.eventMux.Post(messages)
	}

	return nil
}
