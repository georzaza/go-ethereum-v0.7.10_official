package core

import (
	"bytes"
	"math"
	"math/big"

	"github.com/georzaza/go-ethereum-v0.7.10_official/core/types"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/state"
)

// Sometimes the coin value of the output is higher than what the user wishes to pay.
// In this case, the client generates a new Ethereum address, and sends the difference
// back to this address. This is known as "account change" or just "change".

// Change addresses in Ethereum are a privacy feature, not a bug.
// The idea is users would choose a different change address for each transaction.
// Someone analyzing the blockchain would then not be able to tell which address
// was the receivers, and which was the change.
// Ethereum developers decided the privacy potential of change addresses were not
// worth the additional complexity they required.
// One of the methods to increase users’ privacy is coin mixing or tumbling. This technique
// provides k-anonymity or plausible deniability. The idea is that k users deposit 1 coin each and then
// in the course of a coin shuffling protocol either a centralized trusted third party or a smart contract
// mixes the coins and redistributes them to designated fresh public keys. This powerful technique
// gives users superior privacy and anonimity since their new received coins cannot be linked to them
// Possibly, they foresaw on-chain trustless mixers as a viable and more scalable alternative.
//
// Source: https://ethereum.stackexchange.com/questions/581/why-are-there-no-change-addresses
//
// A cryptocurrency tumbler or cryptocurrency mixing service is a service that mixes potentially identifiable
// or "tainted" cryptocurrency funds with others, so as to obscure the trail back to the fund's original source.
// This is usually done by pooling together source funds from multiple inputs for a large and random period of time,
// and then spitting them back out to destination addresses. As all the funds are lumped together and then
// distributed at random times, it is very difficult to trace exact coins.
//
// Source: https://en.wikipedia.org/wiki/Cryptocurrency_tumbler
//
// There is even an event emitted on MetaMask for such cases, see:
// https://docs.metamask.io/guide/ethereum-provider.html#accountschanged

type AccountChange struct {
	Address, StateAddress []byte
}

// Filtering interface
type Filter struct {
	eth             EthManager
	earliest        int64
	latest          int64
	skip            int
	from, to        [][]byte
	max             int
	Altered         []AccountChange
	BlockCallback   func(*types.Block)
	MessageCallback func(state.Messages)
}

// Create a new filter which uses a bloom filter on blocks
// to figure out whether a particular block is interesting or not.
func NewFilter(eth EthManager) *Filter {
	return &Filter{eth: eth}
}

func (self *Filter) AddAltered(address, stateAddress []byte) {
	self.Altered = append(self.Altered, AccountChange{address, stateAddress})
}

// Set the earliest and latest block for filtering.
// -1 = latest block (i.e., the current block)
// hash = particular hash from-to
func (self *Filter) SetEarliestBlock(earliest int64) {
	self.earliest = earliest
}

func (self *Filter) SetLatestBlock(latest int64) {
	self.latest = latest
}

func (self *Filter) SetFrom(addr [][]byte) {
	self.from = addr
}

func (self *Filter) AddFrom(addr []byte) {
	self.from = append(self.from, addr)
}

func (self *Filter) SetTo(addr [][]byte) {
	self.to = addr
}

func (self *Filter) AddTo(addr []byte) {
	self.to = append(self.to, addr)
}

func (self *Filter) SetMax(max int) {
	self.max = max
}

func (self *Filter) SetSkip(skip int) {
	self.skip = skip
}

// Run filters messages with the current parameters set
func (self *Filter) Find() []*state.Message {
	var earliestBlockNo uint64 = uint64(self.earliest)
	if self.earliest == -1 {
		earliestBlockNo = self.eth.ChainManager().CurrentBlock().Number.Uint64()
	}
	var latestBlockNo uint64 = uint64(self.latest)
	if self.latest == -1 {
		latestBlockNo = self.eth.ChainManager().CurrentBlock().Number.Uint64()
	}

	var (
		messages []*state.Message
		block    = self.eth.ChainManager().GetBlockByNumber(latestBlockNo)
		quit     bool
	)
	for i := 0; !quit && block != nil; i++ {
		// Quit on latest
		switch {
		case block.Number.Uint64() == earliestBlockNo, block.Number.Uint64() == 0:
			quit = true
		case self.max <= len(messages):
			break
		}

		// Use bloom filtering to see if this block is interesting given the
		// current parameters
		if self.bloomFilter(block) {
			// Get the messages of the block
			msgs, err := self.eth.BlockManager().GetMessages(block)
			if err != nil {
				chainlogger.Warnln("err: filter get messages ", err)

				break
			}

			messages = append(messages, self.FilterMessages(msgs)...)
		}

		block = self.eth.ChainManager().GetBlock(block.PrevHash)
	}

	skip := int(math.Min(float64(len(messages)), float64(self.skip)))

	return messages[skip:]
}

func includes(addresses [][]byte, a []byte) (found bool) {
	for _, addr := range addresses {
		if bytes.Compare(addr, a) == 0 {
			return true
		}
	}

	return
}

func (self *Filter) FilterMessages(msgs []*state.Message) []*state.Message {
	var messages []*state.Message

	// Filter the messages for interesting stuff
	for _, message := range msgs {
		if len(self.to) > 0 && !includes(self.to, message.To) {
			continue
		}

		if len(self.from) > 0 && !includes(self.from, message.From) {
			continue
		}

		var match bool
		if len(self.Altered) == 0 {
			match = true
		}

		for _, accountChange := range self.Altered {
			if len(accountChange.Address) > 0 && bytes.Compare(message.To, accountChange.Address) != 0 {
				continue
			}

			if len(accountChange.StateAddress) > 0 && !includes(message.ChangedAddresses, accountChange.StateAddress) {
				continue
			}

			match = true
			break
		}

		if !match {
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

/*
 * Returns either true or false based on if both "fromIncluded" and "toIncluded" are true or not.
 */
func (self *Filter) bloomFilter(block *types.Block) bool {
	var fromIncluded, toIncluded bool

	// if Filter.from has been set (bytes array)
	if len(self.from) > 0 {

		// for every byte of Filter.from
		for _, from := range self.from {

			// if ...
			if types.BloomLookup(block.LogsBloom, from) || bytes.Equal(block.Coinbase, from) {
				fromIncluded = true
				break
			}
		}
	} else {
		fromIncluded = true
	}

	if len(self.to) > 0 {
		for _, to := range self.to {

			// if type.BloomLookup( (1+to) & (2^256 - 1) )  || ...
			if types.BloomLookup(block.LogsBloom, ethutil.U256(new(big.Int).Add(ethutil.Big1, ethutil.BigD(to))).Bytes()) || bytes.Equal(block.Coinbase, to) {
				toIncluded = true
				break
			}
		}
	} else {
		toIncluded = true
	}

	return fromIncluded && toIncluded
}
