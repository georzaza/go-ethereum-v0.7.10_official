package core

import (
	"bytes"
	"container/list"
	"fmt"
	"math/big"
	"sync"

	"github.com/georzaza/go-ethereum-v0.7.10_official/core/types"
	"github.com/georzaza/go-ethereum-v0.7.10_official/event"
	"github.com/georzaza/go-ethereum-v0.7.10_official/logger"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
	"github.com/georzaza/go-ethereum-v0.7.10_official/wire"
)

var txplogger = logger.NewLogger("TXP")

// Used to initialize the queueChan field of a TxPool
// (which is used as a queue channel to reading and writing transactions)
const txPoolQueueSize = 50

// Although defined, this type is never used in this go-ethereum version.
type TxPoolHook chan *types.Transaction

// Although defined, this variable is never used in this go-ethereum version.
const (
	minGasPrice = 1000000
)

// Although defined, this variable is never used in this go-ethereum version.
var MinGasPrice = big.NewInt(10000000000000)

// The only use of a TxMsgTy type is as a field of a TxMsg type.
// Although it's not clear how this type is supposed to be used, since
// there are no other references to it in the whole codebase, we could
// make a guess based on the fact that there are 3 types of transactions
// in ethereum:
// Regular transactions: a transaction from one wallet to another.
// Contract deployment transactions: a transaction without a 'to' address,
// where the data field is used for the contract code.
// Execution of a contract: a transaction that interacts with a
// deployed smart contract. In this case, 'to' address is the smart contract address.
// TxMsgTy may had been used to represent the above types
// of transactions.
type TxMsgTy byte

// todo TxMsg represents the type of the channel of
// the _subscribers_ field of a TxPool object. However, that field
// is never actually used in the whole codebase of this go-ethereum version.
type TxMsg struct {
	Tx   *types.Transaction
	Type TxMsgTy
}

// todo EachTx is used as a means to iterate over the _pool_ list of transactions.
func EachTx(pool *list.List, it func(*types.Transaction, *list.Element) bool) {
	for e := pool.Front(); e != nil; e = e.Next() {
		if it(e.Value.(*types.Transaction), e) {
			break
		}
	}
}

// todo FindTx searches in the caller's transactions for _finder_ and returns
// either the matching transaction if found, or nil. This is a neat way of searching
// for transactions that match the criteria defined from the _finder_ param. For example,
// todo Add uses the hash of a transaction as a searching criterion.
func FindTx(pool *list.List, finder func(*types.Transaction, *list.Element) bool) *types.Transaction {
	for e := pool.Front(); e != nil; e = e.Next() {
		if tx, ok := e.Value.(*types.Transaction); ok {
			if finder(tx, e) {
				return tx
			}
		}
	}

	return nil
}

// The todo TxProcessor interface, although defined, is not implemented
// at all in the whole codebase of this go-ethereum version.
type TxProcessor interface {
	ProcessTransaction(tx *types.Transaction)
}

// TxPool is a thread safe transaction pool handler. In order to
// guarantee a non blocking pool the _queueChan_ is used which can be
// independently read without needing access to the actual pool. If the
// pool is being drained or synced for whatever reason, the transactions
// will simply queue up and be handled when the mutex is freed.
// mutex: a mutex for accessing the Tx pool.
// queueChan: Queueing channel for reading and writing incoming transactions to
// quit: Quiting channel (quitting is equivalent to emptying the TxPool)
// pool: The actual pool, aka the list of transactions.
// SecondaryProcessor: This field is actually never used as the todo TxProcessor interface is not implemented.
// subscribers: Although defined, this channel is never used.
// broadcaster: used to broadcast messages to all connected peers.
// chainManager: the chain to which the TxPool object is attached to.
// eventMux: used to dispatch events to subscribers.
type TxPool struct {
	mutex              sync.Mutex
	queueChan          chan *types.Transaction
	quit               chan bool
	pool               *list.List
	SecondaryProcessor TxProcessor
	subscribers        []chan TxMsg
	broadcaster        types.Broadcaster
	chainManager       *ChainManager
	eventMux           *event.TypeMux
}

// todo NewTxPool creates a new todo TxPool object and sets it's fields.
// TxPool.pool will be empty.
// TxPool.queueChain wil be set to a Transaction channel with a txPoolQueueSize size.
// TxPool.quit will be set to a boolean channel.
// TxPool.chainManager will be assigned the param _chainManager_
// TxPool.eventMux will be assigned the param _eventMux_
// TxPool.broadcaster will be assigned the param _broadcaster_
// All other fields of the todo TxPool object that gets created are not set by NewTxPool.
func NewTxPool(chainManager *ChainManager, broadcaster types.Broadcaster, eventMux *event.TypeMux) *TxPool {
	return &TxPool{
		pool:         list.New(),
		queueChan:    make(chan *types.Transaction, txPoolQueueSize),
		quit:         make(chan bool),
		chainManager: chainManager,
		eventMux:     eventMux,
		broadcaster:  broadcaster,
	}
}

// todo addTransaction is an inner function used to add the
// todo Transaction _tx_ to the end of the TxPool. Also, a message
// is broadcasted to all peers which contains the rlp-encodable fields
// _tx_ (See todo RlpData)
// todo locked.
func (pool *TxPool) addTransaction(tx *types.Transaction) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	pool.pool.PushBack(tx)

	// Broadcast the transaction to the rest of the peers
	pool.broadcaster.Broadcast(wire.MsgTxTy, []interface{}{tx.RlpData()})
}

// todo ValidateTransaction validates the _tx_ todo Transaction.
// Returns either an error if _tx_ can not be validated or nil.
// These are the cases where _tx_ is not validated:
// 1. For some reason, the currentBlock field of the chainManager field
// of the caller is nil. (aka the chain is empty)
// 2. The recipient field (_to_) of _tx_ is is either nil or != 20 bytes.
// This means that the validation of contract creation transactions,
// for which  the recipient (_to_) is set to nil will always return an error.
// 3. The _v_ field of _tx_ is neither 28 nor 27. (See todo Transaction object)
// 4. The sender account of _tx_ does not have enough Ether to send to the recipient of _tx_
func (pool *TxPool) ValidateTransaction(tx *types.Transaction) error {
	// Get the last block so we can retrieve the sender and receiver from
	// the merkle trie
	block := pool.chainManager.CurrentBlock
	// Something has gone horribly wrong if this happens
	if block == nil {
		return fmt.Errorf("No last block on the block chain")
	}

	if len(tx.To()) != 0 && len(tx.To()) != 20 {
		return fmt.Errorf("Invalid recipient. len = %d", len(tx.To()))
	}

	v, _, _ := tx.Curve()
	if v > 28 || v < 27 {
		return fmt.Errorf("tx.v != (28 || 27)")
	}

	// Get the sender
	sender := pool.chainManager.State().GetAccount(tx.Sender())

	totAmount := new(big.Int).Set(tx.Value())
	// Make sure there's enough in the sender's account. Having insufficient
	// funds won't invalidate this transaction but simple ignores it.
	if sender.Balance().Cmp(totAmount) < 0 {
		return fmt.Errorf("Insufficient amount in sender's (%x) account", tx.From())
	}

	// Increment the nonce making each tx valid only once to prevent replay
	// attacks

	return nil
}

// todo Add is the function to be called for adding a todo Transaction to the todo TxPool caller.
// Returns either an error on not successfully adding _tx_ or nil for success.
// If _tx_ was added, a message is posted to the subscribed peers, containing the
// _tx_ from, to, value and hash fields.
// An error is returned in any of these cases:
// 1. _tx_'s hash already exists in the todo TxPool caller, aka the transaction
// to be added is already part of the caller.
// 2. _tx_ validation returned an error when calling todo ValidateTransaction.
// If no errors are produced from steps 1. and 2. todo Add makes a call to the inner
// function todo addTransaction to add _tx_ to the current todo TxPool.
func (self *TxPool) Add(tx *types.Transaction) error {
	hash := tx.Hash()
	foundTx := FindTx(self.pool, func(tx *types.Transaction, e *list.Element) bool {
		return bytes.Compare(tx.Hash(), hash) == 0
	})

	if foundTx != nil {
		return fmt.Errorf("Known transaction (%x)", hash[0:4])
	}

	err := self.ValidateTransaction(tx)
	if err != nil {
		return err
	}

	self.addTransaction(tx)

	txplogger.Debugf("(t) %x => %x (%v) %x\n", tx.From()[:4], tx.To()[:4], tx.Value, tx.Hash())

	// Notify the subscribers
	go self.eventMux.Post(TxPreEvent{tx})

	return nil
}

// todo Size returns the number of Transactions of the caller.
func (self *TxPool) Size() int {
	return self.pool.Len()
}

// todo CurrentTransactions returns the transactions of the todo TxPool caller as a slice.
func (pool *TxPool) CurrentTransactions() []*types.Transaction {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	txList := make([]*types.Transaction, pool.pool.Len())
	i := 0
	for e := pool.pool.Front(); e != nil; e = e.Next() {
		tx := e.Value.(*types.Transaction)

		txList[i] = tx

		i++
	}

	return txList
}

// todo RemoveInvalid removed all transactions from the caller for which either:
// 1. the transaction returns an error when validated through the todo ValidateTransaction function, or
// 2. the transaction sender's nonce field is >= to the transaction's nonce field.
func (pool *TxPool) RemoveInvalid(state *state.StateDB) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for e := pool.pool.Front(); e != nil; e = e.Next() {
		tx := e.Value.(*types.Transaction)
		sender := state.GetAccount(tx.Sender())
		err := pool.ValidateTransaction(tx)
		if err != nil || sender.Nonce >= tx.Nonce() {
			pool.pool.Remove(e)
		}
	}
}

// todo RemoveSet takes as an argument a set of transactions _txs_ and
// removes from the caller's transactions set those that match the ones from _txs_.
// Looping over the transactions of the caller happens through the todo EachTx function.
func (self *TxPool) RemoveSet(txs types.Transactions) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, tx := range txs {
		EachTx(self.pool, func(t *types.Transaction, element *list.Element) bool {
			if t == tx {
				self.pool.Remove(element)
				return true // To stop the loop
			}
			return false
		})
	}
}

// todo Flush resets the caller's transactions list to an empty list.
func (pool *TxPool) Flush() []*types.Transaction {
	txList := pool.CurrentTransactions()

	// Recreate a new list all together
	// XXX Is this the fastest way?
	pool.pool = list.New()

	return txList
}

// Although defined, this function does not contain any executable code in this go-ethereum version.
func (pool *TxPool) Start() {
	//go pool.queueHandler()
}

// todo Stop makes a call on todo Flush to empty the caller's transactions list
// and then sends the message "Stopped" to the todo txplogger channel.
func (pool *TxPool) Stop() {
	pool.Flush()

	txplogger.Infoln("Stopped")
}
