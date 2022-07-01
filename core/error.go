package core

import (
	"fmt"
	"math/big"
)

// Parent error. In case a parent is unknown this error will be thrown
// by the block manager
type ParentErr struct {
	Message string
}

// Returns the error message of the caller.
func (err *ParentErr) Error() string {
	return err.Message
}

// Creates a ParentError object by setting it's message to a message that includes the 'hash' and returns it.
func ParentError(hash []byte) error {
	return &ParentErr{Message: fmt.Sprintf("Block's parent unkown %x", hash)}
}

// Returns whether 'err' is a ParentErr error.
func IsParentErr(err error) bool {
	_, ok := err.(*ParentErr)

	return ok
}

// Uncle error. This error is thrown only from the function BlockManager.AccumelateRewards defined in the 
// 'core' package (file block_manager.go). See that function for more.
type UncleErr struct {
	Message string
}

// Returns the error message of an UncleErr error.
func (err *UncleErr) Error() string {
	return err.Message
}

// Creates an UncleErr error by setting it's message to 'str' and returns it.
func UncleError(str string) error {
	return &UncleErr{Message: str}
}

// Returns whether 'err' is an UncleErr error.
func IsUncleErr(err error) bool {
	_, ok := err.(*UncleErr)

	return ok
}

// Block validation error. If any validation fails, this error will be thrown
type ValidationErr struct {
	Message string
}

// Returns the error message of a ValidationErr error.
func (err *ValidationErr) Error() string {
	return err.Message
}

// Creates a ValidationErr error by setting it's message and returns it.
func ValidationError(format string, v ...interface{}) *ValidationErr {
	return &ValidationErr{Message: fmt.Sprintf(format, v...)}
}

// Returns whether 'err' is a ValidationErr error.
func IsValidationErr(err error) bool {
	_, ok := err.(*ValidationErr)

	return ok
}

// Happens when the total gas left for the coinbase address is less than the gas to be bought.
type GasLimitErr struct {
	Message string
	Is, Max *big.Int
}

// Returns whether 'err' is a GasLimitErr error.
func IsGasLimitErr(err error) bool {
	_, ok := err.(*GasLimitErr)

	return ok
}

// Returns the error message of a GasLimitErr error.
func (err *GasLimitErr) Error() string {
	return err.Message
}

// Creates and returns a GasLimitError given the total gas left to be bought by a coinbase address and the actual gas.
func GasLimitError(is, max *big.Int) *GasLimitErr {
	return &GasLimitErr{Message: fmt.Sprintf("GasLimit error. Max %s, transaction would take it to %s", max, is), Is: is, Max: max}
}

// Happens when a transaction's nonce is incorrect. 
type NonceErr struct {
	Message string
	Is, Exp uint64
}

// Returns the error message of a NonceErr error.
func (err *NonceErr) Error() string {
	return err.Message
}

// Creates and returns a NonceError given the transaction's nonce and the nonce of the sender of the transaction. 
func NonceError(is, exp uint64) *NonceErr {
	return &NonceErr{Message: fmt.Sprintf("Nonce err. Is %d, expected %d", is, exp), Is: is, Exp: exp}
}

// Returns whether 'err' is a NonceErr error.
func IsNonceErr(err error) bool {
	_, ok := err.(*NonceErr)

	return ok
}

// Happens when the gas provided runs out before the state transition happens.
type OutOfGasErr struct {
	Message string
}

// Creates and returns an OutOfGasError error.
func OutOfGasError() *OutOfGasErr {
	return &OutOfGasErr{Message: "Out of gas"}
}

// Returns the error message of an OutOfGasError error.
func (self *OutOfGasErr) Error() string {
	return self.Message
}

// Returns whether 'err' is an OutOfGasErr error.
func IsOutOfGasErr(err error) bool {
	_, ok := err.(*OutOfGasErr)

	return ok
}

// Defined, but not used. Meant to be used when there is a total difficulty error, a < b.
type TDError struct {
	a, b *big.Int
}

// Creates and returns a TDError error.
func (self *TDError) Error() string {
	return fmt.Sprintf("incoming chain has a lower or equal TD (%v <= %v)", self.a, self.b)
}

// Returns whether 'e' is a TDError error.
func IsTDError(e error) bool {
	_, ok := e.(*TDError)
	return ok
}

// Happens when there is already a block in the chain with the same hash. The number field is the number of the existing block.
type KnownBlockError struct {
	number *big.Int
	hash   []byte
}

// Creates and returns a KnownBlockError error.
func (self *KnownBlockError) Error() string {
	return fmt.Sprintf("block %v already known (%x)", self.number, self.hash[0:4])
}

// Returns whether 'e' is a KnownBlockErr error.
func IsKnownBlockErr(e error) bool {
	_, ok := e.(*KnownBlockError)
	return ok
}
