package core

import (
	"fmt"
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/core/types"
	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
	"github.com/georzaza/go-ethereum-v0.7.10_official/vm"
)

/*
 * The State transitioning model
 *
 * A state transition is a change made when a transaction is applied to the current world state
 * The state transitioning model does all all the necessary work to work out a valid new state root.
 * 1) Nonce handling
 * 2) Pre pay / buy gas of the coinbase (miner)
 * 3) Create a new state object if the recipient is \0*32, aka contract creation.
 * 4) Value transfer
 * == If contract creation ==
 * 4a) Attempt to run transaction data
 * 4b) If valid, use result as code for the new state object
 * == end ==
 * 5) Run Script section
 * 6) Derive new state root
 */
type StateTransition struct {
	coinbase, receiver []byte
	msg                Message
	gas, gasPrice      *big.Int
	initialGas         *big.Int
	value              *big.Int
	data               []byte
	state              *state.StateDB
	block              *types.Block
	cb, rec, sen 			*state.StateObject
	Env 							vm.Environment
}

// A Message represents a Transaction.
type Message interface {
	Hash() 			[]byte
	From() 			[]byte
	To() 				[]byte
	GasPrice() 	*big.Int
	Gas() 			*big.Int
	Value() 		*big.Int
	Nonce() 		uint64
	Data() 			[]byte
}

// creates and returns a new address based on a msg sender and nonce fields.
func AddressFromMessage(msg Message) []byte {
	// Generate a new address
	return crypto.Sha3(ethutil.NewValue([]interface{}{msg.From(), msg.Nonce()}).Encode())[12:]
}

// returns whether the _msg_ is a contract creation aka whether the msg's/transaction's
// recipient is 0.
func MessageCreatesContract(msg Message) bool {
	return len(msg.To()) == 0
}

// returns the amount of Wei based on the msg's/transaction's gas and gasPrice fields.
// gasValue = gas * gasPrice
func MessageGasValue(msg Message) *big.Int {
	return new(big.Int).Mul(msg.Gas(), msg.GasPrice())
}

// creates and returns a todo StateTransition object.
// The fields _gas_ and initialGas are set to 0
// The fields rec, sen and Env are set to nil
// The field coinbase is set to the address of the _coinbase_ param.
// The field cb is set to the param _coinbase_.
// All other fields are set to the corresponding params.
func NewStateTransition(coinbase *state.StateObject, msg Message, state *state.StateDB, block *types.Block) *StateTransition {
	return &StateTransition{coinbase.Address(), msg.To(), msg, new(big.Int), new(big.Int).Set(msg.GasPrice()), new(big.Int), msg.Value(), msg.Data(), state, block, coinbase, nil, nil, nil}
}

// is a getter method for the Env field of a StateTransition object.
// If the Env field of the caller is nil, a new Env is created by
// calling the function todo NewEnv defined in the file todo vm_env.go
func (self *StateTransition) VmEnv() vm.Environment {
	if self.Env == nil {
		self.Env = NewEnv(self.state, self.msg, self.block)
	}

	return self.Env
}

// returns the miner's account as a StateObject. If the miner does
// not exist in the current world state, it's created.
func (self *StateTransition) Coinbase() *state.StateObject {
	return self.state.GetOrNewStateObject(self.coinbase)
}

// returns the _from_ field of the caller as a StateObject. If
// _from_ does not exist in the current world state, it's created.
func (self *StateTransition) From() *state.StateObject {
	return self.state.GetOrNewStateObject(self.msg.From())
}

// returns the _to_ field of the caller as a StateObject. If
// _to_ does not exist in the current world state, it's created.
// This function will return nil in the case where the msg is about
// a contract creation (aka if _to_ is 0)
func (self *StateTransition) To() *state.StateObject {
	if self.msg != nil && MessageCreatesContract(self.msg) {
		return nil
	}
	return self.state.GetOrNewStateObject(self.msg.To())
}

// attempts to use _amount_ gas of the caller's gas. If the caller's
// gas is less than the _amount_ provided, an OutOfGasError is returned.
// Otherwise, nil is returned for success. In case of success, the new
// gas of the caller will become newGas = prevGas - amount.
func (self *StateTransition) UseGas(amount *big.Int) error {
	if self.gas.Cmp(amount) < 0 {
		return OutOfGasError()
	}
	self.gas.Sub(self.gas, amount)

	return nil
}

// adds _amount_ gas to the caller's gas.
func (self *StateTransition) AddGas(amount *big.Int) {
	self.gas.Add(self.gas, amount)
}

// attempts to reward the miner with the gas of the transaction.
// If the sender's balance is less than the calculated gas in Wei
// (gas*gasPrice of caller), an error is returned.
// Buying the gas does not directly happen in this function. Instead,
// the BuyGas function of the miner StateObject is called through
// this function. If the latter does not return an error, this function
// will increase the _gas_ field of the caller, set the caller's
// _initialGas_ field and decrease the sender's balance by an amount
// of _gas_ *_gasPrice_
func (self *StateTransition) BuyGas() error {
	var err error

	sender := self.From()
	if sender.Balance().Cmp(MessageGasValue(self.msg)) < 0 {
		return fmt.Errorf("insufficient ETH for gas (%x). Req %v, has %v", sender.Address()[:4], MessageGasValue(self.msg), sender.Balance())
	}

	coinbase := self.Coinbase()
	err = coinbase.BuyGas(self.msg.Gas(), self.msg.GasPrice())
	if err != nil {
		return err
	}

	self.AddGas(self.msg.Gas())
	self.initialGas.Set(self.msg.Gas())
	sender.SubAmount(MessageGasValue(self.msg))

	return nil
}

// is an inner function, used by todo and does 2 things:
// 1. Checks whether the caller's msg sender nonce is the same as the caller's msg nonce.
// If not, it returns an error.
// 2. Calls todo BuyGas in order to reward the miner. If todo BuyGas returns an error,
// this function returns that error.
// If everything went well, this function returns nil.
func (self *StateTransition) preCheck() (err error) {
	var (
		msg    = self.msg
		sender = self.From()
	)

	// Make sure this transaction's nonce is correct
	if sender.Nonce != msg.Nonce() {
		return NonceError(msg.Nonce(), sender.Nonce)
	}

	// Pre-pay gas / Buy gas of the coinbase account
	if err = self.BuyGas(); err != nil {
		return err
	}

	return nil
}

// Attempts to alter the world state by applying the msg/transaction of the caller.
// 1. Calls todo preCheck for nonce validation and to reward the miner.
// 2. Schedules a gas refund todo
// 3. increases nonce of the msg sender.
// 4. uses the TxGas. (defined as Gtransaction in the Ethereum Yellow Paper)
//    TxGas is a constant value set to 500 Wei. (see the GasTx variable defined in vm/common.go)
//    On the Ethereum Yellow Paper this value is set to 21000 Wei.
// 5. uses the GasData. This is the gas that must be payed for every byte
//    of the log field of the msg. On the Ethereum Yellow Paper this value is set to 8 Wei per byte.
// 6. If the msg is about a contract creation (msg recipient is 0) this function makes a call
//    to the todo MakeContract function which creates the contract as a StateObject, then through
//    the caller's _vm_ field

func (self *StateTransition) TransitionState() (ret []byte, err error) {
	statelogger.Debugf("(~) %x\n", self.msg.Hash())

	// XXX Transactions after this point are considered valid.
	if err = self.preCheck(); err != nil {
		return
	}

	var (
		msg    = self.msg
		sender = self.From()
	)

	defer self.RefundGas()

	// Increment the nonce for the next transaction
	sender.Nonce += 1

	// Transaction gas
	if err = self.UseGas(vm.GasTx); err != nil {
		return
	}

	// Pay data gas
	var dgas int64
	for _, byt := range self.data {
		if byt != 0 {
			dgas += vm.GasData.Int64()
		} else {
			dgas += 1 // This is 1/5. If GasData changes this fails
		}
	}
	if err = self.UseGas(big.NewInt(dgas)); err != nil {
		return
	}

	vmenv := self.VmEnv()
	var ref vm.ClosureRef
	if MessageCreatesContract(msg) {
		contract := MakeContract(msg, self.state)
		ret, err, ref = vmenv.Create(sender, contract.Address(), self.msg.Data(), self.gas, self.gasPrice, self.value)
		if err == nil {
			dataGas := big.NewInt(int64(len(ret)))
			dataGas.Mul(dataGas, vm.GasCreateByte)
			if err = self.UseGas(dataGas); err == nil {
				ref.SetCode(ret)
			}
		}
	} else {
		ret, err = vmenv.Call(self.From(), self.To().Address(), self.msg.Data(), self.gas, self.gasPrice, self.value)
	}

	if err != nil {
		self.UseGas(self.gas)
	}

	return
}

// Converts an transaction in to a state object
func MakeContract(msg Message, state *state.StateDB) *state.StateObject {
	addr := AddressFromMessage(msg)

	contract := state.GetOrNewStateObject(addr)
	contract.InitCode = msg.Data()

	return contract
}

func (self *StateTransition) RefundGas() {
	coinbase, sender := self.Coinbase(), self.From()
	// Return remaining gas
	remaining := new(big.Int).Mul(self.gas, self.msg.GasPrice())
	sender.AddAmount(remaining)

	uhalf := new(big.Int).Div(self.GasUsed(), ethutil.Big2)
	for addr, ref := range self.state.Refunds() {
		refund := ethutil.BigMin(uhalf, ref)
		self.gas.Add(self.gas, refund)
		self.state.AddBalance([]byte(addr), refund.Mul(refund, self.msg.GasPrice()))
	}

	coinbase.RefundGas(self.gas, self.msg.GasPrice())
}

func (self *StateTransition) GasUsed() *big.Int {
	return new(big.Int).Sub(self.initialGas, self.gas)
}
