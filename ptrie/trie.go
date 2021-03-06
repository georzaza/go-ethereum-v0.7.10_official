package ptrie

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"

	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/trie"
)

func ParanoiaCheck(t1 *Trie, backend Backend) (bool, *Trie) {
	t2 := New(nil, backend)

	it := t1.Iterator()
	for it.Next() {
		t2.Update(it.Key, it.Value)
	}

	return bytes.Compare(t2.Hash(), t1.Hash()) == 0, t2
}

type Trie struct {
	mu       sync.Mutex
	root     Node
	roothash []byte
	cache    *Cache

	revisions *list.List
}

func New(root []byte, backend Backend) *Trie {
	trie := &Trie{}
	trie.revisions = list.New()
	trie.roothash = root
	trie.cache = NewCache(backend)

	if root != nil {
		value := ethutil.NewValueFromBytes(trie.cache.Get(root))
		trie.root = trie.mknode(value)
	}

	return trie
}

func (self *Trie) Iterator() *Iterator {
	return NewIterator(self)
}

// Legacy support
func (self *Trie) Root() []byte { return self.Hash() }
func (self *Trie) Hash() []byte {
	var hash []byte
	if self.root != nil {
		//hash = self.root.Hash().([]byte)
		t := self.root.Hash()
		if byts, ok := t.([]byte); ok {
			hash = byts
		} else {
			hash = crypto.Sha3(ethutil.Encode(self.root.RlpData()))
		}
	} else {
		hash = crypto.Sha3(ethutil.Encode(""))
	}

	if !bytes.Equal(hash, self.roothash) {
		self.revisions.PushBack(self.roothash)
		self.roothash = hash
	}

	return hash
}
func (self *Trie) Commit() {
	// Hash first
	self.Hash()

	self.cache.Flush()
}

// Reset should only be called if the trie has been hashed
func (self *Trie) Reset() {
	self.cache.Reset()

	revision := self.revisions.Remove(self.revisions.Back()).([]byte)
	self.roothash = revision
	value := ethutil.NewValueFromBytes(self.cache.Get(self.roothash))
	self.root = self.mknode(value)
}

func (self *Trie) UpdateString(key, value string) Node { return self.Update([]byte(key), []byte(value)) }
func (self *Trie) Update(key, value []byte) Node {
	self.mu.Lock()
	defer self.mu.Unlock()

	k := trie.CompactHexDecode(string(key))

	if len(value) != 0 {
		self.root = self.insert(self.root, k, &ValueNode{self, value})
	} else {
		self.root = self.delete(self.root, k)
	}

	return self.root
}

func (self *Trie) GetString(key string) []byte { return self.Get([]byte(key)) }
func (self *Trie) Get(key []byte) []byte {
	self.mu.Lock()
	defer self.mu.Unlock()

	k := trie.CompactHexDecode(string(key))

	n := self.get(self.root, k)
	if n != nil {
		return n.(*ValueNode).Val()
	}

	return nil
}

func (self *Trie) DeleteString(key string) Node { return self.Delete([]byte(key)) }
func (self *Trie) Delete(key []byte) Node {
	self.mu.Lock()
	defer self.mu.Unlock()

	k := trie.CompactHexDecode(string(key))
	self.root = self.delete(self.root, k)

	return self.root
}

func (self *Trie) insert(node Node, key []byte, value Node) Node {
	if len(key) == 0 {
		return value
	}

	if node == nil {
		return NewShortNode(self, key, value)
	}

	switch node := node.(type) {
	case *ShortNode:
		k := node.Key()
		cnode := node.Value()
		if bytes.Equal(k, key) {
			return NewShortNode(self, key, value)
		}

		var n Node
		matchlength := trie.MatchingNibbleLength(key, k)
		if matchlength == len(k) {
			n = self.insert(cnode, key[matchlength:], value)
		} else {
			pnode := self.insert(nil, k[matchlength+1:], cnode)
			nnode := self.insert(nil, key[matchlength+1:], value)
			fulln := NewFullNode(self)
			fulln.set(k[matchlength], pnode)
			fulln.set(key[matchlength], nnode)
			n = fulln
		}
		if matchlength == 0 {
			return n
		}

		return NewShortNode(self, key[:matchlength], n)

	case *FullNode:
		cpy := node.Copy().(*FullNode)
		cpy.set(key[0], self.insert(node.branch(key[0]), key[1:], value))

		return cpy

	default:
		panic("Invalid node")
	}
}

func (self *Trie) get(node Node, key []byte) Node {
	if len(key) == 0 {
		return node
	}

	if node == nil {
		return nil
	}

	switch node := node.(type) {
	case *ShortNode:
		k := node.Key()
		cnode := node.Value()

		if len(key) >= len(k) && bytes.Equal(k, key[:len(k)]) {
			return self.get(cnode, key[len(k):])
		}

		return nil
	case *FullNode:
		return self.get(node.branch(key[0]), key[1:])
	default:
		panic(fmt.Sprintf("%T: invalid node: %v", node, node))
	}
}

func (self *Trie) delete(node Node, key []byte) Node {
	if len(key) == 0 {
		return nil
	}

	switch node := node.(type) {
	case *ShortNode:
		k := node.Key()
		cnode := node.Value()
		if bytes.Equal(key, k) {
			return nil
		} else if bytes.Equal(key[:len(k)], k) {
			child := self.delete(cnode, key[len(k):])

			var n Node
			switch child := child.(type) {
			case *ShortNode:
				nkey := append(k, child.Key()...)
				n = NewShortNode(self, nkey, child.Value())
			case *FullNode:
				n = NewShortNode(self, node.key, child)
			}

			return n
		} else {
			return node
		}

	case *FullNode:
		n := node.Copy().(*FullNode)
		n.set(key[0], self.delete(n.branch(key[0]), key[1:]))

		pos := -1
		for i := 0; i < 17; i++ {
			if n.branch(byte(i)) != nil {
				if pos == -1 {
					pos = i
				} else {
					pos = -2
				}
			}
		}

		var nnode Node
		if pos == 16 {
			nnode = NewShortNode(self, []byte{16}, n.branch(byte(pos)))
		} else if pos >= 0 {
			cnode := n.branch(byte(pos))
			switch cnode := cnode.(type) {
			case *ShortNode:
				// Stitch keys
				k := append([]byte{byte(pos)}, cnode.Key()...)
				nnode = NewShortNode(self, k, cnode.Value())
			case *FullNode:
				nnode = NewShortNode(self, []byte{byte(pos)}, n.branch(byte(pos)))
			}
		} else {
			nnode = n
		}

		return nnode

	default:
		panic("Invalid node")
	}
}

// casting functions and cache storing
func (self *Trie) mknode(value *ethutil.Value) Node {
	l := value.Len()
	switch l {
	case 2:
		return NewShortNode(self, trie.CompactDecode(string(value.Get(0).Bytes())), self.mknode(value.Get(1)))
	case 17:
		fnode := NewFullNode(self)
		for i := 0; i < l; i++ {
			fnode.set(byte(i), self.mknode(value.Get(i)))
		}
		return fnode
	case 32:
		return &HashNode{value.Bytes()}
	default:
		return &ValueNode{self, value.Bytes()}
	}
}

func (self *Trie) trans(node Node) Node {
	switch node := node.(type) {
	case *HashNode:
		value := ethutil.NewValueFromBytes(self.cache.Get(node.key))
		return self.mknode(value)
	default:
		return node
	}
}

func (self *Trie) store(node Node) interface{} {
	data := ethutil.Encode(node)
	if len(data) >= 32 {
		key := crypto.Sha3(data)
		self.cache.Put(key, data)

		return key
	}

	return node.RlpData()
}
