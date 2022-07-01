package types

import (
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/trie"
)

// Any type that implements this interface can then call the function types.DeriveSha(list DerivableList)
// Len(): returns the number of the elements of the object that implements this interface.
// GetRlp(i int): returns the RLP-encoding of i-th element of the object that implements this interface.
type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}

// Takes as an argument a list of DerivableList objects, constructs a trie and returns the root hash of the trie.
func DeriveSha(list DerivableList) []byte {
	trie := trie.New(ethutil.Config.Db, "")
	for i := 0; i < list.Len(); i++ {
		trie.Update(string(ethutil.NewValue(i).Encode()), string(list.GetRlp(i)))
	}

	return trie.GetRoot()
}
