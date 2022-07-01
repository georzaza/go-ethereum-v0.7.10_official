package core

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/georzaza/go-ethereum-v0.7.10_official/core/types"
	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/event"
	"github.com/georzaza/go-ethereum-v0.7.10_official/logger"
	"github.com/georzaza/go-ethereum-v0.7.10_official/pow"
	"github.com/georzaza/go-ethereum-v0.7.10_official/pow/ezp"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
	"github.com/georzaza/go-ethereum-v0.7.10_official/wire"
)

// a Blocks logger channel.
var statelogger = logger.NewLogger("BLOCK")

// This interface is only implemented by a "Peer" object defined in the 'eth' package (file peer.go)
//
// Inbound(): Determines whether it's an inbound or outbound peer
//
// LastSend(): Last known message send time.
//
// LastPong(): Last received pong message
//
// Host(): the host.
//
// Port(): the port of the connection
//
// Version(): client identity
//
// PingTime(): Used to give some kind of pingtime to a node, not very accurate.
//
// Connected(): Flag for checking the peer's connectivity state
//
// Caps(): getter for the protocolCaps field of a Peer object.
type Peer interface {
	Inbound() bool
	LastSend() time.Time
	LastPong() int64
	Host() []byte
	Port() uint16
	Version() string
	PingTime() string
	Connected() *int32
	Caps() *ethutil.Value
}

// The 'EthManager' interface is only implemented by the 'Ethereum' object, defined in the package 'eth' (file ethereum.go)
//
// BlockManager: See type 'BlockManager' of the 'core' package.
//
// ChainManager: See type 'ChainManager' of the 'core' package.
//
// TxPool: See type 'TxPool' of the 'core' package.
//
// Broadcast: Used to broacast messages to all peers.
//
// IsMining: Returns whether the peer is mining.
//
// IsListening: Returns whether the peer is listening.
//
// Peers: Returns all connected peers.
//
// KeyManager: Key manager for the account(s) of the node.
//
// ClientIdentity: Used mainly for communication between peers.
//
// Db: The ethereum database. This should be a representation of the World State.
//
// EventMux: Used to dispatch events to registered nodes
type EthManager interface {
	BlockManager() *BlockManager
	ChainManager() *ChainManager
	TxPool() *TxPool
	Broadcast(msgType wire.MsgType, data []interface{})
	PeerCount() int
	IsMining() bool
	IsListening() bool
	Peers() *list.List
	KeyManager() *crypto.KeyManager
	ClientIdentity() wire.ClientIdentity
	Db() ethutil.Database
	EventMux() *event.TypeMux
}

// State manager for processing new blocks and managing the overall state
//
// mutex: Mutex for locking the block processor. Blocks can only be handled one at a time
//
// bc: Canonical block chain
//
// mem: non-persistent key/value memory storage
//
// Pow: Proof of work used for validating
//
// txpool: The transaction pool. See type 'TxPool' of the 'core' package.
//
// lastAttemptedBlock: The last attempted block is mainly used for debugging purposes.
//
// This does not have to be a valid block and will be set during 'Process' & canonical validation.
//
// events: Provides a way to subscribe to events.
//
// eventMux: event mutex, used to dispatch events to subscribers.
type BlockManager struct {
	mutex              sync.Mutex
	bc                 *ChainManager
	mem                map[string]*big.Int
	Pow                pow.PoW
	txpool             *TxPool
	lastAttemptedBlock *types.Block
	events             event.Subscription
	eventMux           *event.TypeMux
}

// Creates a new BlockManager object by initializing these fields of a BlockManager object type: mem, Pow, bc, eventMux, txpool.
// The Pow object will have it's 'turbo' field set to true when created. See the type 'EasyPow' of the 'ezp' package for more. (file pow/ezp/pow.go)
func NewBlockManager(txpool *TxPool, chainManager *ChainManager, eventMux *event.TypeMux) *BlockManager {
	sm := &BlockManager{
		mem:      make(map[string]*big.Int),
		Pow:      ezp.New(),
		bc:       chainManager,
		eventMux: eventMux,
		txpool:   txpool,
	}
	return sm
}

// Returns (receipts, nil) or (nil, error) if an IsGasLimitErr error has occured.
// This function together with the BlockManager.ApplyTransactions function form a recursion algorithm meant to apply all transactions
// of the current block to the world state. The main logic/application happens in the latter function. However, this function is the one
// to be called when we want to apply the transactions of a block. The TransitionState function sets
// the total gas pool (amount of gas left) for the coinbase address of the block before calling the ApplyTransactions function which
// in turn will call the TransitionState function again, and so on, until all transactions have been applied or an IsGasLimitErr error
// has occured. If no such an error has occured then the 'receipts' object returned by this function
// will hold the receipts that are the result of the application of the transactions of the 'block' param.
func (sm *BlockManager) TransitionState(statedb *state.StateDB, parent, block *types.Block) (receipts types.Receipts, err error) {
	coinbase := statedb.GetOrNewStateObject(block.Coinbase)
	coinbase.SetGasPool(block.CalcGasLimit(parent))
	receipts, _, _, _, err = sm.ApplyTransactions(coinbase, statedb, block, block.Transactions(), false)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

// This function will apply the transactions of a block to the world state and return the results as a tuple.
// It gets called by the BlockManager.TransitionState function, then calls the latter again, and so on, to form a recursion
// algorithm that will apply the transactions one by one. In case where an IsGasLimitErr error occurs during the application of any
// transaction, the process of the transactions stops.
//
// Returns: (receipts, handled, unhandled, erroneous, err)
//
// receipts: The receipts up to but not including any transaction that has caused an IsGasLimitErr error.
//
// handled: All transactions that were handled up to but not including any transaction that has caused an IsGasLimitErr error.
//
// unhandled: In case of an IsGasLimitErr error this object will contain all the transactions that were not applied (includes the transaction
// that caused the error). Otherwise, this object will be nil.
//
// erroneous: Any transactions that caused an error other than an IsGasLimitErr and/or an IsNonceErr errors.
//
// err: The err will be either an IsGasLimitErr error type or nil.
//
// A short description on what this function does follows.
//
// 1. Clear all state logs.
//
// 2. Get (or create a new) coinbase state object and call the TransitionState function.
//
// 3. The latter function will call this function again, forming the recursion.
//
// 4. If an error occured and is an IsGasLimitErr error then stop the process and set the 'unhandled' variable
// (to be returned later). If it is a IsNonceErr error, ignore it. If it is any other error, also ignore it,
// but append to the variable 'erroneous' (to be returned later) the transaction that caused that error.
//
// 5. Calculate the gas used so far and the current reward for the miner. Update the state.
//
// 6. Create the receipt of the current transaction and set the receipt's logs and Bloom field.
//
// 7. If the parameter 'transientProcess' is false, notify all subscribers about the transaction.
//
// 8. Append receipt, transaction to the 'receipts', 'handled' variables (to be returned later)
//
// 9. When the processing has ended, set the block's reward and totalUsedGas fields.
//
// 10.Return the results.
func (self *BlockManager) ApplyTransactions(coinbase *state.StateObject, state *state.StateDB, block *types.Block, txs types.Transactions, transientProcess bool) (types.Receipts, types.Transactions, types.Transactions, types.Transactions, error) {
	var (
		receipts           types.Receipts
		handled, unhandled types.Transactions
		erroneous          types.Transactions
		totalUsedGas       = big.NewInt(0)
		err                error
		cumulativeSum      = new(big.Int)
	)

done:
	for i, tx := range txs {
		// If we are mining this block and validating we want to set the logs back to 0
		state.EmptyLogs()

		txGas := new(big.Int).Set(tx.Gas())

		cb := state.GetStateObject(coinbase.Address())
		st := NewStateTransition(cb, tx, state, block)
		_, err = st.TransitionState()
		if err != nil {
			switch {
			case IsNonceErr(err):
				err = nil
				continue
			case IsGasLimitErr(err):
				unhandled = txs[i:]

				break done
			default:
				statelogger.Infoln(err)
				erroneous = append(erroneous, tx)
				err = nil
			}
		}

		txGas.Sub(txGas, st.gas)
		cumulativeSum.Add(cumulativeSum, new(big.Int).Mul(txGas, tx.GasPrice()))
		state.Update(txGas)

		cumulative := new(big.Int).Set(totalUsedGas.Add(totalUsedGas, txGas))
		receipt := types.NewReceipt(state.Root(), cumulative)
		receipt.SetLogs(state.Logs())
		receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
		chainlogger.Debugln(receipt)

		if !transientProcess {
			go self.eventMux.Post(TxPostEvent{tx})
		}

		receipts = append(receipts, receipt)
		handled = append(handled, tx)

		if ethutil.Config.Diff && ethutil.Config.DiffType == "all" {
			state.CreateOutputForDiff()
		}
	}

	block.Reward = cumulativeSum
	block.GasUsed = totalUsedGas

	return receipts, handled, unhandled, erroneous, err
}

// Processes a block. When successful, returns the return result of a call to the ProcessWithParent function.
//
// Otherwise, in case that the hash of the block or the hash of the parent of the block already exist in the ChainManager,
// returns the tuple (nil, nil KnownBlockError) or (nil, nil, ParentError) accordingly.
//
// Before calling the ProcessWithParent function, this function takes care of locking the BlockManager with a mutex and
// only after the ProcessWithParent function has returned the BlockManager is being unlocked.
func (sm *BlockManager) Process(block *types.Block) (td *big.Int, msgs state.Messages, err error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if sm.bc.HasBlock(block.Hash()) {
		return nil, nil, &KnownBlockError{block.Number, block.Hash()}
	}
	if !sm.bc.HasBlock(block.PrevHash) {
		return nil, nil, ParentError(block.PrevHash)
	}
	parent := sm.bc.GetBlock(block.PrevHash)
	return sm.ProcessWithParent(block, parent)
}

// The main process function of a block. Gets called by the function Process.
//
// Returns a tuple (td, messages, error) where td: total difficulty of the processed block, messages: messages generated
// by the application of the transactions of the block (and some extra messages like transactions regarding rewarding miners).
// If an error occured the returned tuple becomes (nil, nil, error). If no errors have occured the returned tuple becomes (td, messages, nil)
//
//
// 1. Saves a copy of the state. Also queues a reset of the state after the function's return based on that copy.
//
// 2. Validates the block with a call to the function ValidateBlock. If errors happened, returns
//
// 3. Calls TransitionState to attempt to do the state transition. If errors, returns.
//
// 4. Creates the bloom field of the receipts returned from step 2. If for some reason the bloom field is different
// from the bloom field of the provided block, it returns.
//
// 5. Validates the transactions and the receipts root hashes. If errors, returns.
//
// 6. Calls AccumelateRewards to calculate the miner rewards. If errors, returns.
//
// 7. Sets the state to 0 and makes a call to CalculateTD in order to calculate the total difficulty of the block. If errors, returns.
// If not, the last step is to remove the block's transactions from the BlockManager's txpool, sync the state db, cancel the queued
// state reset, send a message to the chainlogger channel and finally return the tuple (td, messages, nil).
func (sm *BlockManager) ProcessWithParent(block, parent *types.Block) (td *big.Int, messages state.Messages, err error) {
	sm.lastAttemptedBlock = block
	state := parent.State().Copy()
	defer state.Reset()

	// Block validation
	if err = sm.ValidateBlock(block, parent); err != nil {
		return
	}
	receipts, err := sm.TransitionState(state, parent, block)
	if err != nil {
		return
	}
	rbloom := types.CreateBloom(receipts)
	if bytes.Compare(rbloom, block.LogsBloom) != 0 {
		err = fmt.Errorf("unable to replicate block's bloom=%x", rbloom)
		return
	}
	txSha := types.DeriveSha(block.Transactions())
	if bytes.Compare(txSha, block.TxSha) != 0 {
		err = fmt.Errorf("validating transaction root. received=%x got=%x", block.TxSha, txSha)
		return
	}
	receiptSha := types.DeriveSha(receipts)
	if bytes.Compare(receiptSha, block.ReceiptSha) != 0 {
		fmt.Printf("%x\n", ethutil.Encode(receipts))
		err = fmt.Errorf("validating receipt root. received=%x got=%x", block.ReceiptSha, receiptSha)
		return
	}
	if err = sm.AccumelateRewards(state, block, parent); err != nil {
		return
	}
	state.Update(ethutil.Big0)
	if !block.State().Cmp(state) {
		err = fmt.Errorf("invalid merkle root. received=%x got=%x", block.Root(), state.Root())
		return
	}

	if td, ok := sm.CalculateTD(block); ok {
		state.Sync()
		messages := state.Manifest().Messages
		state.Manifest().Reset()
		chainlogger.Infof("Processed block #%d (%x...)\n", block.Number, block.Hash()[0:4])
		sm.txpool.RemoveSet(block.Transactions())
		return td, messages, nil
	} else {
		return nil, nil, errors.New("total diff failed")
	}
}

// Calculates the total difficulty for a given block.
//
// TD(genesis_block)=0 and TD(block)=TD(block.parent) + sum(u.difficulty for u in block.uncles) + block.difficulty
//
// Returns: If the calculated total difficulty is greater than the previous the tuple (total_difficulty, true) is returned.
// Otherwise, the tuple (nil, false) is returned.
func (sm *BlockManager) CalculateTD(block *types.Block) (*big.Int, bool) {
	uncleDiff := new(big.Int)
	for _, uncle := range block.Uncles {
		uncleDiff = uncleDiff.Add(uncleDiff, uncle.Difficulty)
	}

	td := new(big.Int)
	td = td.Add(sm.bc.Td(), uncleDiff)
	td = td.Add(td, block.Difficulty)

	if td.Cmp(sm.bc.Td()) > 0 {
		return td, true
	}

	return nil, false
}

// Validates the current block. Returns an error if the block was invalid,
// an uncle or anything that isn't on the current block chain.
// Validation validates easy over difficult (dagger takes longer time = difficult)
func (sm *BlockManager) ValidateBlock(block, parent *types.Block) error {
	expd := CalcDifficulty(block, parent)
	if expd.Cmp(block.Difficulty) < 0 {
		return fmt.Errorf("Difficulty check failed for block %v, %v", block.Difficulty, expd)
	}

	diff := block.Time - parent.Time
	if diff < 0 {
		return ValidationError("Block timestamp less then prev block %v (%v - %v)", diff, block.Time, sm.bc.CurrentBlock().Time)
	}

	// Verify the nonce of the block. Return an error if it's not valid
	if !sm.Pow.Verify(block /*block.HashNoNonce(), block.Difficulty, block.Nonce*/) {
		return ValidationError("Block's nonce is invalid (= %v)", ethutil.Bytes2Hex(block.Nonce))
	}

	return nil
}

// Calculates the reward of the miner. Returns an error if an error has occured during the
// validation process. If no errors have occured, nil is returned.
//
// More specifically an error is returned:
// 1. if the parent of any of the uncles of the 'block' is nil, or
//
// 2. if the (block) number of the parent of any of the uncles of the 'block' and the 'block' itself have a difference greater than 6, or
//
// 3. if the hash of any of the uncles of the param 'block' matches any of the uncles of the param 'parent'.
//
// 4. if the nonce of any of the uncles of the param 'block' is included in the nonce of the 'block'
//
// The reward to be appointed to the miner will be:
//
// If the 'block' has 1 uncle: r1 = BlockReward + BlockReward/32,
//
// If the 'block' has 2 uncles: r2 = r1 + r1/32, etc., where BlockReward = 1.5 Ether,( defined in the core package, file fees.go)
//
// Finally, a message is added to the state manifest regarding the value to be transferred to the miner address.
// This value will be the sum of the above calculated reward and the block.Reward.
func (sm *BlockManager) AccumelateRewards(statedb *state.StateDB, block, parent *types.Block) error {
	reward := new(big.Int).Set(BlockReward)

	knownUncles := ethutil.Set(parent.Uncles)
	nonces := ethutil.NewSet(block.Nonce)
	for _, uncle := range block.Uncles {
		if nonces.Include(uncle.Nonce) {
			// Error not unique
			return UncleError("Uncle not unique")
		}

		uncleParent := sm.bc.GetBlock(uncle.PrevHash)
		if uncleParent == nil {
			return UncleError(fmt.Sprintf("Uncle's parent unknown (%x)", uncle.PrevHash[0:4]))
		}

		if uncleParent.Number.Cmp(new(big.Int).Sub(parent.Number, big.NewInt(6))) < 0 {
			return UncleError("Uncle too old")
		}

		if knownUncles.Include(uncle.Hash()) {
			return UncleError("Uncle in chain")
		}

		nonces.Insert(uncle.Nonce)

		r := new(big.Int)
		r.Mul(BlockReward, big.NewInt(15)).Div(r, big.NewInt(16))

		uncleAccount := statedb.GetAccount(uncle.Coinbase)
		uncleAccount.AddAmount(r)

		reward.Add(reward, new(big.Int).Div(BlockReward, big.NewInt(32)))
	}

	// Get the account associated with the coinbase
	account := statedb.GetAccount(block.Coinbase)
	// Reward amount of ether to the coinbase address
	account.AddAmount(reward)

	statedb.Manifest().AddMessage(&state.Message{
		To:     block.Coinbase,
		Input:  nil,
		Origin: nil,
		Block:  block.Hash(), Timestamp: block.Time, Coinbase: block.Coinbase, Number: block.Number,
		Value: new(big.Int).Add(reward, block.Reward),
	})

	return nil
}

// Returns either the tuple (state.Manifest().Messages, nil) or (nil, error)
//
// If an error is returned it will be a ParentError regarding the parent of the 'block'
// (the error includes the hash of the parent of the 'block'). This error happens in the case where the
// the hash of the parent of the 'block' already exists.
//
// In essence, the state manifest's messages are the transactions that occured during the world state transition
// of the addition of a 'block'.
//
// To get those messages a simple trick is used: a deferred call on state.Reset() is queued and only then
// a call of the function TransitionState and following that a call on AccumelateRewards happen.
func (sm *BlockManager) GetMessages(block *types.Block) (messages []*state.Message, err error) {
	if !sm.bc.HasBlock(block.PrevHash) {
		return nil, ParentError(block.PrevHash)
	}

	sm.lastAttemptedBlock = block

	var (
		parent = sm.bc.GetBlock(block.PrevHash)
		state  = parent.State().Copy()
	)

	defer state.Reset()

	sm.TransitionState(state, parent, block)
	sm.AccumelateRewards(state, block, parent)

	return state.Manifest().Messages, nil
}
