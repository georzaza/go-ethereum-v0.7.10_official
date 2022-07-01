package types

import (
	"fmt"
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/obscuren/secp256k1-go"
)

// Returns true if the provided 'addr' has a length of 0, false otherwise.
func IsContractAddr(addr []byte) bool {
	return len(addr) == 0
}

// Definition of a transaction object.
//
// Nonce: A value equal to the number of transactions that originated from this address or,
// in the case of accounts with associated code, the number of contract-creations made by this account.
//
// recipient: The 160-bit address of the transaction's recipient or, for a contract creation transaction, nil.
//
// value: This is equal to the number of Wei to be transferred to the transaction's recipient or,
// in the case of contract creation, as an endowment to the newly created account.
//
// gas: The gas limit equal to the maximum amount of gas that should be used when executing this transaction.
//
// gasPrice: This is equal to the number of Wei to be paid per unit of gas for all computation costs
// incurred as a result of the execution of this transaction.
//
// data: An unlimited size byte array specifying the input data of the transaction.
// The first 4 bytes of this field specify which function to call when
// the transaction will execute by using the hash of the function's name to be
// called and it's arguments. The rest of the data field are the arguments passed
// to the function. If the data field is empty, it means a transaction is for a
// payment and not an execution of the contract.
//
// v, r, s: Values corresponding to the ECDSA digital signature of the transaction and used to determine the originating Externally Owned Account.
type Transaction struct {
	nonce     uint64
	recipient []byte
	value     *big.Int
	gas       *big.Int
	gasPrice  *big.Int
	data      []byte
	v         byte
	r, s      []byte
}

// Creates and returns a new Transaction which represents the creation of a contract. The Transaction's 'recipient' field
// will be set to nil and the Transaction's 'data' field will be set equal to the 'script' parameter. All other Transaction
// fields will be set to the corresponing parameters passed to this function ( except of course for the fields nonce,v,r,s which are not set)
func NewContractCreationTx(value, gas, gasPrice *big.Int, script []byte) *Transaction {
	return &Transaction{recipient: nil, value: value, gas: gas, gasPrice: gasPrice, data: script}
}

// Creates and returns a new Transaction. All Transaction fields will be set to the corresponing parameters passed
// to this function ( except of course for the fields nonce,v,r,s which are not set)
func NewTransactionMessage(to []byte, value, gas, gasPrice *big.Int, data []byte) *Transaction {
	return &Transaction{recipient: to, value: value, gasPrice: gasPrice, gas: gas, data: data}
}

// Creates and returns a new Transaction based on the 'data' parameter. The latter should be a valid rlp-encoding of a
// Transaction object. (The 'data' parameter is converted to an ethutil.Value object first and then gets decoded)
func NewTransactionFromBytes(data []byte) *Transaction {
	tx := &Transaction{}
	tx.RlpDecode(data)

	return tx
}

// Creates and returns a new Transaction based on the 'val' parameter. The latter should be a valid rlp-encoding of
// a Transaction object that has been cast to an ethutil.Value object.
func NewTransactionFromValue(val *ethutil.Value) *Transaction {
	tx := &Transaction{}
	tx.RlpValueDecode(val)

	return tx
}

// Returns the hash of the caller. The hash will be the Kecchak-256 hash of these Transaction fields: nonce, gasPrice,
// gas, recipient, value, data.
func (tx *Transaction) Hash() []byte {
	data := []interface{}{tx.nonce, tx.gasPrice, tx.gas, tx.recipient, tx.value, tx.data}

	return crypto.Sha3(ethutil.NewValue(data).Encode())
}

// Returns the data field of the caller.
func (self *Transaction) Data() []byte {
	return self.data
}

// Returns the gas field of the caller.
func (self *Transaction) Gas() *big.Int {
	return self.gas
}

// Returns the gasPrice field of the caller.
func (self *Transaction) GasPrice() *big.Int {
	return self.gasPrice
}

// Returns the value field of the caller.
func (self *Transaction) Value() *big.Int {
	return self.value
}

// Returns the nonce field of the caller.
func (self *Transaction) Nonce() uint64 {
	return self.nonce
}

// Sets the caller's nonce field equal to the 'nonce' parameter.
func (self *Transaction) SetNonce(nonce uint64) {
	self.nonce = nonce
}

// Returns the sender of the transaction.
func (self *Transaction) From() []byte {
	return self.Sender()
}

// Returns the recipient field of the caller.
func (self *Transaction) To() []byte {
	return self.recipient
}

// Returns v, r, s of the Transaction. r and s are left-padded to 32 bytes. For more, see the function ethutil.leftPadBytes
func (tx *Transaction) Curve() (v byte, r []byte, s []byte) {
	v = tx.v
	r = ethutil.LeftPadBytes(tx.r, 32)
	s = ethutil.LeftPadBytes(tx.s, 32)

	return
}

// Calculates and returns the signature of the hash of a transaction with the param key.
// The signature used to be obtained through the github.com/obscuren/secp256k1-go package repo but this package no longer exists.
// One may now use the function Sign from https://github.com/ethereum/go-ethereum/blob/master/crypto/secp256k1/secp256.go
func (tx *Transaction) Signature(key []byte) []byte {
	hash := tx.Hash()

	sig, _ := secp256k1.Sign(hash, key)

	return sig
}

// Retrieves and returns the public key of the transaction. The public key used to be obtained through the
// github.com/obscuren/secp256k1-go package repo but this package no longer exists.
// One may now use the function RecoverPubKey from https://github.com/ethereum/go-ethereum/blob/master/crypto/secp256k1/secp256.go
func (tx *Transaction) PublicKey() []byte {
	hash := tx.Hash()

	v, r, s := tx.Curve()

	sig := append(r, s...) // sig = r appended by s
	sig = append(sig, v-27)

	pubkey, _ := secp256k1.RecoverPubkey(hash, sig)

	return pubkey
}

// Returns the sender of the transaction. To do so, the public key of the transaction is retrieved first through the function
// Transaction.PublicKey(). If the public key passes validation then the last 12 bytes of the public key are returned (aka the sender address)
func (tx *Transaction) Sender() []byte {
	pubkey := tx.PublicKey()
	if len(pubkey) != 0 && pubkey[0] != 4 {
		return nil
	}
	return crypto.Sha3(pubkey[1:])[12:]
}

// Signes the transaction. To do so, the function Transaction.Signature is called, which makes use of a non-existent package. Refer to
// the Transaction.Signature for more. After the transaction has been signed the r,s,v fields of the Transaction are set
// to the signatures appropriate fields.
//
// r is set to the first 32 bytes of the signature.
//
// s is set to the 32-th up to and including the 63-th bytes of the signature.
//
// v is set to the sum of the last bit of the transaction and the number 27.
func (tx *Transaction) Sign(privk []byte) error {

	sig := tx.Signature(privk)

	tx.r = sig[:32]
	tx.s = sig[32:64]
	tx.v = sig[64] + 27

	return nil
}

// Returns the rlp-encodable fields of the caller (all fields of a Transaction object).
func (tx *Transaction) RlpData() interface{} {
	data := []interface{}{tx.nonce, tx.gasPrice, tx.gas, tx.recipient, tx.value, tx.data}

	return append(data, tx.v, new(big.Int).SetBytes(tx.r).Bytes(), new(big.Int).SetBytes(tx.s).Bytes())
}

// Casts and returns the caller to an ethutil.Value object.
func (tx *Transaction) RlpValue() *ethutil.Value {
	return ethutil.NewValue(tx.RlpData())
}

// Returns the rlp-encoding of the cast of the caller object to an ethutil.Value object. The cast is done
// through the Transaction.RlpValue() function)
func (tx *Transaction) RlpEncode() []byte {
	return tx.RlpValue().Encode()
}

// Sets the caller's fields equal to the rlp-decoding of the 'data' parameter. The 'data' parameter must
// be a valid bytes representation of a Transaction that has been cast to an ethutil.Value object.
func (tx *Transaction) RlpDecode(data []byte) {
	tx.RlpValueDecode(ethutil.NewValueFromBytes(data))
}

// Sets the caller's fields equal to the rlp-decoding of the 'decoder'.
// The 'decoder' must be Transaction cast to an ethutil.Value object.
func (tx *Transaction) RlpValueDecode(decoder *ethutil.Value) {
	tx.nonce = decoder.Get(0).Uint()
	tx.gasPrice = decoder.Get(1).BigInt()
	tx.gas = decoder.Get(2).BigInt()
	tx.recipient = decoder.Get(3).Bytes()
	tx.value = decoder.Get(4).BigInt()
	tx.data = decoder.Get(5).Bytes()
	tx.v = byte(decoder.Get(6).Uint())

	tx.r = decoder.Get(7).Bytes()
	tx.s = decoder.Get(8).Bytes()
}

// Returns the string representation of the caller.
func (tx *Transaction) String() string {
	return fmt.Sprintf(`
	TX(%x)
	Contract: %v
	From:     %x
	To:       %x
	Nonce:    %v
	GasPrice: %v
	Gas:      %v
	Value:    %v
	Data:     0x%x
	V:        0x%x
	R:        0x%x
	S:        0x%x
	Hex:      %x
	`,
		tx.Hash(),
		len(tx.recipient) == 0,
		tx.Sender(),
		tx.recipient,
		tx.nonce,
		tx.gasPrice,
		tx.gas,
		tx.value,
		tx.data,
		tx.v,
		tx.r,
		tx.s,
		ethutil.Encode(tx),
	)
}

// Transaction slice type for basic sorting
type Transactions []*Transaction

// Returns the rlp-encodable fields of all the Transaction objects that the caller consists of.
func (self Transactions) RlpData() interface{} {
	// Marshal the transactions of this block
	enc := make([]interface{}, len(self))
	for i, tx := range self {
		// Cast it to a string (safe)
		enc[i] = tx.RlpData()
	}

	return enc
}

// Returns the number of the Transaction objects that the caller consists of.
// This function, along with the GetRlp function are implemented so as that a Transaction
// object can call the function DeriveSha (defined in the file types/derive_sha.go) which
// constructs a trie and returns the hash root of that trie.
func (s Transactions) Len() int { return len(s) }

// Swaps two Transaction objects.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Returns the rlp-encoding of the i-th Transaction of the caller.
// This function, along with the Len function are implemented so as that a Transaction
// object can call the function DeriveSha (defined in the file types/derive_sha.go) which
// constructs a trie and returns the hash root of that trie.
func (s Transactions) GetRlp(i int) []byte { return ethutil.Rlp(s[i]) }

// Data type used for performing comparison operations on Transaction objects.
type TxByNonce struct{ Transactions }

// Returns true if the i-th Transaction of a list of Transactions is less than the j-th Transaction.
// This is determined only by comparing the two Transactions' nonces.
func (s TxByNonce) Less(i, j int) bool {
	return s.Transactions[i].nonce < s.Transactions[j].nonce
}
