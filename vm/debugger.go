package vm

import "github.com/georzaza/go-ethereum-v0.7.10_official/state"

type Debugger interface {
	BreakHook(step int, op OpCode, mem *Memory, stack *Stack, object *state.StateObject) bool
	StepHook(step int, op OpCode, mem *Memory, stack *Stack, object *state.StateObject) bool
	BreakPoints() []int64
	SetCode(byteCode []byte)
}
