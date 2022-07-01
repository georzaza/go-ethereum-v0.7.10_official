package ar

import (
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/trie"
)

type Block interface {
	Trie() *trie.Trie
	Diff() *big.Int
}
