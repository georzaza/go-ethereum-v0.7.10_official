package core

import (
	"fmt"
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/vm"
)

// Returns a string representation of a sequence of bytes that consist an evm bytecode.
// The opcodes are defined in vm/types.go. In case that we have a PUSHi opcode we expect
// the next i bytes to be the i items that we want to push to the stack.
//
// script: The evm bytecode. An example can be found here:
// https://rinkeby.etherscan.io/address/0x147b8eb97fd247d06c4006d269c90c1908fb5d54#code
// 
// Example: Passing the first series of bytes of the above link to this function as
//  
// script = []byte(0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15, 0x61, 0x00, 0x10, 0x57, 0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50, 0x60, 0x40, 0x51)
//
// will yield the following output:
//
// 0x60	0000: PUSH1
//
// 0x80	0001: 0x80
//
// 0x60	0002: PUSH1 (we got a PUSH1, so the next value is pushed onto the stack)
//
// 0x40	0003: 0x40 
//
// 0x52	0004: MSTORE
//
// 0x34	0005: CALLVALUE
//
// 0x80	0006: DUP1
//
// 0x15	0007: ISZERO
//
// 0x61	0008: PUSH2 (we got a PUSH2, so the next 2 values are pushed onto the stack)
//
// 0x00	0009: 0x00
//
// 0x10	0010: 0x10
//
// 0x57	0011: JUMPI
//
// 0x60	0012: PUSH1
//
// 0x00	0013: 0x00 
//
// 0x80	0014: DUP1
//
// 0xfd	0015: Missing opcode 0xfd
//
// 0x5b	0016: JUMPDEST
//
// 0x50	0017: POP
//
// 0x60	0018: PUSH1
//
// 0x40	0019: 0x40
//
// 0x51	0020: MLOAD
func Disassemble(script []byte) (asm []string) {
	pc := new(big.Int)
	for {
		if pc.Cmp(big.NewInt(int64(len(script)))) >= 0 {
			return
		}

		// Get the memory location of pc
		val := script[pc.Int64()]
		// Get the opcode (it must be an opcode!)
		op := vm.OpCode(val)

		asm = append(asm, fmt.Sprintf("%04v: %v", pc, op))

		switch op {
		case vm.PUSH1, vm.PUSH2, vm.PUSH3, vm.PUSH4, vm.PUSH5, vm.PUSH6, vm.PUSH7, vm.PUSH8,
			vm.PUSH9, vm.PUSH10, vm.PUSH11, vm.PUSH12, vm.PUSH13, vm.PUSH14, vm.PUSH15,
			vm.PUSH16, vm.PUSH17, vm.PUSH18, vm.PUSH19, vm.PUSH20, vm.PUSH21, vm.PUSH22,
			vm.PUSH23, vm.PUSH24, vm.PUSH25, vm.PUSH26, vm.PUSH27, vm.PUSH28, vm.PUSH29,
			vm.PUSH30, vm.PUSH31, vm.PUSH32:
			pc.Add(pc, ethutil.Big1)

			// For a PUSHi command, jump i positions of the bytecode to get to the next opcode. (i = a)
			a := int64(op) - int64(vm.PUSH1) + 1
			if int(pc.Int64()+a) > len(script) {
				return
			}

			data := script[pc.Int64() : pc.Int64()+a]
			if len(data) == 0 {
				data = []byte{0}
			}
			asm = append(asm, fmt.Sprintf("%04v: 0x%x", pc, data))

			pc.Add(pc, big.NewInt(a-1))
		}

		pc.Add(pc, ethutil.Big1)
	}

	return asm
}
