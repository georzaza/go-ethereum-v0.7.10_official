<style>
# SOME MD VIEWERS WILL NOT UNDERSTAND THIS SEGMENT. THIS DOES NOT AFFECT THE VIEW OF THIS DOCUMENT.
body {
	margin: auto;
	max-width: 700px;
}

a {
	text-decoration: none;
}

a.file {
	color: brown;
	font-size: 20px;
}

a.link{
	color: cornflowerblue;
}
a.struct {
	color: seagreen;
}

a.interface {
	color: blueviolet;
}

a.function {
	color: lightseagreen;
}
</style>


<body>

<br>

## go-ethereum v0.7.10 - `core` Module - Source Code Analysis
In this project I dive into the core module of the official go-ethereumv0.7.10 and try  
to give a thorough explanation on all code that's included in this module.

The code can alternatively be explored through:  
https://pkg.go.dev/github.com/georzaza/go-ethereum-v0.7.10_official  
The above link contains the official code along with comments that I've written on  
in order to help myself get through with this project.

Any inner links in this file point just below the definition/signature of  
a data structure/function. Remember to scroll up just a bit after clicking  
on a link.  

If you are interested in analyzing any of the other modules you might find useful  
the folder `georzaza_scripts` and/or the file `georzaza_output`, which examine file  
similarity across all modules of this go-ethereum version.

<br><br>

## 1. Description

The core module is the implementation of the basic logic of Ethereum.  
What is shown below is a list of all files the `core` contains with a 
short description.  
With very few exceptions, each file implements a ... core idea of Ethereum.

<br>

- <a class="link" href="#block.go">core/types/block.go</a>  
	Block type definition along with functions on how to create a block,   
	set it's header, find a transaction in a block, set it's uncles and more.

<br>

- <a class="link" href="#bloom9.go">core/types/bloom9.go</a>  
	Bloom definition. The bloom is used in order to save space while storing  
	a transaction receipt's log fields. 

<br>

- <a class="link" href="#common.go">core/types/common.go</a>  
	Contains the definition of 2 interfaces, namely BlockProcessor and Broadcaster  
	which are implemented by the ChainManager and TxPool accordingly. 

<br>

- <a class="link" href="#derive_sha.go">core/types/derive_sha.go</a>  
	Provides a way for any kind of object to create a trie based on it's contents .

<br>

- <a class="link" href="#receipt.go">core/types/receipt.go</a>  
	Contains the definition of a receipt and functions which may be used on a receipt object.

<br>

- <a class="link" href="#transaction.go">core/types/transaction.go</a>  
	Contains the definition of a transaction. It's important to note that in this file      
	functions are defined which either sign a transaction or get it's public key.

<br>

- <a class="link" href="#asm.go">core/asm.go</a>  
	EVM bytecode to opcode disassembler.

<br>

- <a class="link" href="#block_manager.go">core/block_manager.go</a>  
	Aside from the BlockManager, this file also contains the definition of a very  
	important interface, 'EthManager', which is implemented by the 'Ethereum' object,  
	defined in the file ethereum.go. The BlockManager is responsible for:  
	1. Applying the transactions of a block. 
	2. Validating the block.  
	3. Calculating the new total difficulty and the miner's reward.

<br>

- <a class="link" href="#chain_manager.go">core/chain_manager.go</a>
	The ChainManager as it's name suggests is responsible for:
	1. Adding a block to the chain. Note however that the BlockManager is the one  
		responsible for applying the transactions. After adding a block, that block   
		will also be saved in the main db.
	2. Reseting the chain. This means that through the ChainManager we can reset  
		the chain to the point where it will include only the genesis block.
	3. Providing ways to determine whether a block is in the chain or not.
	5. Calculating the total difficulty. 
	6. Inserting a chain of blocks to the current chain. (function 
	<a class="link" href="#ChainManager.InsertChain">InsertChain</a>)  
	This file also includes the definition of the function 
	<a class="link" href="#CalcDifficulty">CalcDifficulty</a> which calculates  
	the difficulty of a block (which is different than the total difficulty of the chain).

<br>

- <a class="link" href="#dagger.go">core/dagger.go</a>  
	The Dagger consensus model. Dagger is a memory-hard proof of work model   
	based on directed acyclic graphs. Memory hardness refers to the concept that  
	computing a valid proof of work should require not only a large number of   
	computations, but also a large amount of memory.

<br>

- <a class="link" href="#error.go">core/error.go</a>  
	Contains the definition of error types that might come up when either   
	validating a block, applying it's transactions or inserting it into the chain.

<br>

- <a class="link" href="#events.go">core/events.go</a>  
	Containts the definition of events fired when:
	1. a transaction is added to the transaction pool
	2. a transaction has been processed
	3. a block has been imported to the chain.

<br>

- <a class="link" href="#go</">execution.go</a>  
	Contains the definition of an Execution type alongside the 
	most basic (inner) function  
	through which all code execution happens.
	The EVM creates and calls an Execution's   
	object functions in order to execute code.

<br>

- <a class="link" href="#fees.go">core/fees.go</a>  
	Contains the definition of the <a class="link" href="#BlockReward">BlockReward</a> 
	variable.

<br>

- <a class="link" href="#filter.go">core/filter.go</a>  
	todo

<br>

- <a class="link" href="#genesis.go">core/genesis.go</a>  
	Contains the definition of the genesis block.

<br>

- <a class="link" href="#simple_pow.go">core/simple_pow.go</a>  
	Does not contain any code.

<br>

- <a class="link" href="#state_transition.go">core/state_transition.go</a>  
	Contains the implementation of the state transition that comes up
	through the execution of a msg/transaction.   
	In short, it can be defined as follows:
	
	1. Check if the transaction is well-formed (ie. has the right number of values), the signature is valid, 
		and the nonce matches the nonce in the sender's account. If not, return an error.
	2. Calculate the transaction fee as initialGas * gasPrice, and determine the sending address from the 
		signature. Subtract the fee from the sender's account balance and increment the sender's nonce. If 	
		there is not enough balance to spend, return an error.
	3. Initialize gas = initialGas, and take off a certain quantity of gas per byte to pay for the bytes in 
		the transaction.
	4. Transfer the transaction value from the sender's account to the receiving account. If the receiving 
		account does not yet exist, create it. If the receiving account is a contract, run the contract's 
		code either to completion or until the execution runs out of gas.
	5. If the value transfer failed because the sender did not have enough money, or the code execution ran 
		out of gas, revert all state changes except the payment of the fees, and add the fees to the miner's 
		account.
	6. Otherwise, refund the fees for all remaining gas to the sender, and send the fees paid for gas 
		consumed to the miner.

<br>

- <a class="link" href="#transaction_pool.go">core/transaction_pool.go</a>  
	Contains the definition of <a class="link" href="#TxPool">TxPool</a>. A TxPool object
	is used to group transactions.  
	Validation of transactions happens through this object's methods.

<br>

- <a class="link" href="#vm_env.go">core/vm_env.go</a>  
	Definition of the EVM. The EVM implements the Environment interface,   
	defined in the file vm/environment.go and works closely with the 
	<a class="link" href="#execution.go">Execution</a>  
	model in order to handle smart contract deployment and code execution.   
	The below image shows in detail how most of the above concepts are glued together
	with the EVM:	
	
	<br>

	![hello](https://cypherpunks-core.github.io/ethereumbook/images/evm-architecture.png)
<br>


<br><br>

## 1. Package <a class="file" name="types"> types </a>


###  1.1 <a class="file" name="block.go"> block.go </a>

#### 1.1.1 Data structures

```
type BlockInfo struct {
	Number uint64
	Hash   []byte
	Parent []byte
	TD     *big.Int
}
```
A 'compact' way of representing a block is by the use of the 
<a class="struct" name="BlockInfo">BlockInfo</a> type
- **Number**: The number of the block
- **Hash**: The hash of the block
- **Parent**: The parent of the block.
- **TD**: used by the package core to store the total difficulty 
	of the chain up to and including a block.  
	More specifically:   
	$Block.TD = Block.parent.TD + Block.difficulty + \sum_{u\;\epsilon\;Uncles} {u.difficulty}$<br><br>

```
type Block struct {
	PrevHash          ethutil.Bytes
	Uncles            Blocks
	UncleSha          []byte
	Coinbase          []byte
	state             *state.StateDB
	Difficulty        *big.Int
	Time              int64
	Number            *big.Int
	GasLimit          *big.Int
	GasUsed           *big.Int
	Extra             string
	Nonce             ethutil.Bytes
	transactions      Transactions
	receipts          Receipts
	TxSha, ReceiptSha []byte
	LogsBloom         []byte
	Reward            *big.Int
}
```
A <a class="struct" name="Block">Block</a> 
type is the main type used to represent a block.
- **PrevHash**: The Keccak 256-bit hash of the parent block’s header
- **Uncles** : The uncles of the Block.
- **UncleSha**: The Keccak 256-bit hash of the uncles of the Block.
- **CoinBase**: The address of the beneficiary.
- **beneficiary**: The 160-bit address to which all fees collected from   
	the successful mining of this block are transferred to (miner's address).
- **state**: The Keccak 256-bit hash of the root node of the state trie,   
	after all transactions are executed
- **Difficulty** : Difficulty of the current block.
- **Time**: Creation time of the Block.
- **Number**: The number of the Block. This is equal to the number of 
	ancestor blocks.
- **GasLimit**: The maximum gas limit all the transactions inside this
	Block are allowed to consume.
- **GasUsed**: The total gas used by the transactions of the Block.
- **Extra**: An arbitrary byte array containing data relevant to this block.  
	According to the Ethereum yellow paper this should be 32 bytes or less.
- **Nonce**: The Block nonce. This is what miners keep changing to compute
	a solution to PoW.
- **transactions**: List of transactions and/or contracts to be created 
	included in this Block.
- **receipts**: The receipts of the transactions
- **TxSha**: The Keccak 256-bit hash of the root node of the transactions 
	trie of the block
- **ReceiptSha**: The Keccak 256-bit hash of the root node of the receipts 
	trie of the block.
- **LogsBloom**: The Bloom filter composed from indexable information 
	(logger address and log  
	topics) contained in each log entry from 
	the receipt of each transaction in the transactions list.
- **Reward**: The reward of the beneficiary (miner)


Notes about uncles: "In case of proof-of-work mining, there are many miners   
trying to mine the same set of transactions at the same time. Since the block   
mining time is very short (about 15 sec. in case of ethereum) there is a   
possibility,that more than one blocks are mined within a very short interval.  
The block mined first is added to the main chain but the effort of miner who  
mined the other block in not simply let off. These competing blocks are called  
orphaned blocks.

According to ethereum beige paper, " An uncle is a block whose parent is equal  
to the current block’s parent’s parent." The purpose of uncles is to help reward  
miners for including these orphaned blocks. The uncles that miners include must  
be “valid,” meaning within the sixth generation or smaller of the present block.  
After six children, stale orphaned blocks can no longer be referenced.

Uncle blocks receive a smaller reward than a full block. Nonetheless, there’s   
still some incentive for miners to include these orphaned blocks and reap a reward."

Source: https://medium.com/@preethikasireddy/how-does-ethereum-work-anyway-22d1df506369
<br><br>

```
type Blocks []*Block
```
<a class="struct" name="Blocks">Blocks</a> grouping type. <br><br>

```
type BlockBy func(b1, b2 *Block) bool

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}
```
A <a class="struct" name="BlockBy">BlockBy</a> 
type along with 
a <a class="struct" name="blockSorter">blockSorter</a> 
type are used as a sorting mechanism  
(which makes use of the sort.Sort function: 
https://pkg.go.dev/sort#Sort ).  
Note that the blockSorter type is not exported (does not start with a capital 
letter),  
so it's only visible to the current file.

- **blocks**: the blocks to be sorted (can be any number of Blocks).
- **by**: The sorting function to be used between two Blocks of the blockSorter.   
	Although internally a BlockBy object is assigned to this field, note   
	that the definition allows any function that is matching the function   
	signature of 'by') to be assigned to this field.

<br>

#### 1.1.2 Functions


`func (bi *BlockInfo) RlpDecode(data []byte)`  
<a class="function" name="BlockInfo.RlpDecode">RlpDecode</a>
sets the caller's 
(<a class="link" href="#BlockInfo">BlockInfo</a>) 
fields to the RLP-decoded parts of the 'data' parameter.  
To do so, _data_ is converted to an ethutil.Value object and the RLP-decoding operation happens on the Value object.<br><br>


`func (bi *BlockInfo) RlpEncode() []byte`  
<a class="function" name="BlockInfo.RlpEncode"> RlpEncode</a>
returns the rlp-encoding of a <a class="link" href="#BlockInfo"> BlockInfo</a> 
object (by calling the function Encode defined in ethutil/rlp.go)<br><br>


`func (self Blocks) AsSet() ethutil.UniqueSet`  
<a class="function" name="Blocks.AsSet"> AsSet</a>
returns the Blocks that the caller consists of as an ethutil.UniqueSet object.   
The elements of the returned set will be the hashes of these Blocks (only a 
<a class="link" href="#Blocks">Blocks</a>  object can call this function).<br><br>


`func (self BlockBy) Sort(blocks Blocks)`  
<a class="function" name="BlockBy.Sort">Sort </a> 
sorts Blocks.   
This function creates a 
<a class="link" href="#blockSorter">blockSorter</a> object, 
assigns the caller (<a class="link" href="#BlockBy">BlockBy</a>)
object to the blockSorter's 'by' field   
and the param 'blocks' object to the 
blockSorter's 'blocks' field, and then calls the sort.Sort function 
(https://pkg.go.dev/sort#Sort) by passing to it the created blockSorter object.  
The sort.Sort function makes use of Len, Less, Swap functions to sort the objects,
all defined below.<br><br>


`func (self blockSorter) Len() int`  
<a class="function" name="blockSorter.Len">Len</a> 
returns the number of Blocks that a 
<a class="link" href="#blockSorter">blockSorter</a> 
object consists of. <br><br>

`func (self blockSorter) Swap(i, j int)`  
<a class="function" name="blockSorter.Swap">Swap</a> 
swaps two Blocks of a 
<a class="link" href="#blockSorter"> blockSorter</a> 
object.<br><br>

`func (self blockSorter) Less(i, j int)`  
<a class="function" name="blockSorter.Less">Less</a>. For any 2 Blocks
that the caller (<a class="link" href="#blockSorter">blockSorter</a>) consists of,
denoted as i, j, this function returns  
true if Block i is less than j. To determine
whether i is less than j, a 
<a class="link" href="#BlockBy">BlockBy</a> object is used.<br><br>

`func Number(b1, b2 *Block) bool`  
<a class="function" name="Number">Number</a>
returns true if the block number of b1 is less than b2, false otherwise.<br><br>

`func NewBlockFromBytes(raw []byte) *Block`  
<a class="function" name="NewBlockFromBytes">NewBlockFromBytes</a>
creates a new Block from the 'raw' bytes param and returns it.   
**These bytes however should be the bytes representation of the result 
of the RLP-encoding of a <a class="link" href="#Block">Block</a>.** <br><br>

`func NewBlockFromRlpValue(rlpValue *ethutil.Value) *Block`  
<a class="function" name="NewBlockFromRlpValue">NewBlockFromRlpValue </a>
creates a new Block from the 'rlpValue' param.  
To do so, this function calls 
the <a class="link" href="#Block.RlpValueDecode">Block.RlpValueDecode</a>
function.  As described in the latter,  
only the Block's header, transactions
and uncles fields are derived from the rlpValue object.<br><br>

`func CreateBlock(root interface{}, prevHash []byte, base []byte, Difficulty *big.Int, Nonce []byte, extra string) *Block `  
<a class="function" name="CreateBlock">CreateBlock</a> 
creates a Block and returns it. See the <a class="link" href="#Block">Block</a> 
data structure for an explanation of the fields used.  
The root parameter is the root of the Block's state trie. 
The Block created will be created at the current Unix time.  
Notice that there are no transactions and receipts passed as arguments
to this function, aka these are not set. <br><br>

`func (block *Block) Hash() ethutil.Bytes`  
<a class="function" name="Block.Hash">Hash</a>
returns the block's hash. To do so, these are the steps followed: 
1. block's header is first cast to an ethutil.Value object
2. the Encode function is called upon the casted object.
3. the keccak-256 hash of the returned object of the step 2 is returned.<br><br>

`func (block *Block) HashNoNonce() []byte`  
<a class="function" name="Block.HashNoNonce">HashNoNonce</a>
returns the hash of an object that is almost the same as a Block. The differences are:
1. The object to be hashed contains only the uncles hash.
2. The object to be hashed contains the root of the state and not the state as a whole.
3. The object to be hashed contains only the TxSha and not the receipts.
4. The object to be hashed contains only the ReceiptSha and not the transactions.
5. The object to be hashed does not contain the Reward of the miner.

Note: The object to be hashed if appended with the Nonce field of the Block 
will comprise the Block's header.<br><br>

`func (block *Block) State() *state.StateDB`  
<a class="function" name="Block.State">State</a>
returns the state of the Block. (not just the root)<br><br>

`func (block *Block) Transactions() Transactions`  
<a class="function" name="Block.Transactions">Transactions</a>
returns the transactions of the block.<br><br>

`func (block *Block) CalcGasLimit(parent *Block) *big.Int`  
<a class="function" name="Block.CalcGasLimit">CalcGasLimit</a>
calculates the gas limit.  
If the Block passed as a parameter (aka the _parent_)
is the genesis block the gas limit is set to $10^6$.  
Otherwise the gas limit will be:  
$1023 \cdot parent.GasLimit + parent.GasUsed \cdot \frac{6}{5}$.  
The minimum gas limit is set to $125000$ Wei.<br><br>

`func (block *Block) BlockInfo() BlockInfo`  
<a class="function" name="Block.BlockInfo"> BlockInfo </a>
returns the <a class="link" href="#BlockInfo">BlockInfo</a>
representation of a Block.<br><br>

`func (self *Block) GetTransaction(hash []byte) *Transaction`  
<a class="function" name="Block.GetTransaction"> GetTransaction </a>
searches for a transaction in the current Block by comparing the hash of all   
the transactions in the block to the 'hash' parameter.  
Returns the block's matching transaction, or nil if there is no match.<br><br>

`func (block *Block) Sync()`  
<a class="function" name="Block.Sync"> Sync </a>
syncs the block's state and contract respectively. For more, see the state package.<br><br>

`func (block *Block) Undo()`  
<a class="function" name="Block.Undo"> Undo </a>
resets the block's state to nil.<br><br>

`func (block *Block) rlpReceipts() interface{}`  
<a class="function" name="Block.rlpReceipts"> rlpReceipts </a>
is an inner function that returns the receipts as a string.  
This function is used to derive the ReceiptSha of the receipts.  
Although defined, this function is never used in this go-ethereum version.<br><br>

`func (block *Block) rlpUncles() interface{}`  
<a class="function" name="Block.rlpUncles"> rlpUncles </a> returns the uncles 
as a string.  
This function is used to derive the UncleSha of the uncles.  
Gets called by the 
<a class="link" href="#Block.SetUncles">Block.SetUncles</a>, 
<a class="link" href="#Block.Value">Block.Value</a> and 
<a class="link" href="#Block.RlpData"> Block.RlpData</a> functions.<br><br>

`func (block *Block) SetUncles(uncles []*Block)`  
<a class="function" name="Block.SetUncles"> SetUncles </a>
sets the uncles of Block to the parameter 'uncles'.  
This function is also 
called by the <a class="link" href="#CreateBlock">CreateBlock</a> 
function. Also sets the UncleSha of  
the Block based on the provided parameter.
To do so, the inner function <a class="link" href="#Block.rlpUncles">rlpUncles</a> 
is called.<br><br>

`func (self *Block) SetReceipts(receipts Receipts)`  
<a class="function" name="Block.SetReceipts">SetReceipts</a>
sets the receipts of this block to the parameter 'receipts'.  
Also calculates and sets the LogsBloom field of the block which 
is derived by the receipts.<br><br>

`func (self *Block) SetTransactions(txs Transactions)`  
<a class="function" name="Block.SetTransactions"> SetTransactions </a>
sets the transactions of a Block to the parameter 'transactions'.  
Also calculates and sets the TxSha field which is derived by the transactions.<br><br>

`func (block *Block) Value() *ethutil.Value`  
<a class="function" name="Block.Value"> Value </a>
casts a block to an ethutil.Value object containing the header,  
the transactions and the uncles of the block and then returns it.<br><br>

`func (block *Block) RlpEncode() []byte`  
<a class="function" name="Block.RlpEncode"> RlpEncode </a>
calls the 
<a class="link" href="#Block.Value"> Value</a> 
function on the current Block and then rlp-encodes the Block.  
The rlp-encodable fields of a Block are it's header, transactions and uncles.<br><br>

`func (block *Block) RlpDecode(data []byte)`  
<a class="function" name="Block.RlpDecode"> RlpDecode </a>
RLP-decodes a Block.  
To do so, a new ethutil.Value object is created from the
_data_ param and then  
the function <a class="link" href="#Block.RlpValueDecode"> RlpValueDecode(data)</a> 
is called on the current Block.
- **data** represents the rlp-encodable fields of **this** block and can be   
	obtained through the function 
<a class="link" href="#Block.RlpData"> RlpData</a>.<br><br>

`func (block *Block) RlpValueDecode(decoder *ethutil.Value)`  
<a class="function" name="Block.RlpValueDecode"> RlpValueDecode </a>
RLP-decodes the rlp-encodable fields of a Block, aka the header, transactions and uncles.
- **decoder**: The above fields of **this** Block after having been cast to an ethutil.Value object.<br><br>

`func (self *Block) setHeader(header *ethutil.Value)`  
<a class="function" name="Block.setHeader"> setHeader </a>
is an inner function that sets the header of this block,  
given an ethutil.Value object that contains that information.<br><br>

`func NewUncleBlockFromValue(header *ethutil.Value) *Block`  
<a class="function" name="NewUncleBlockFromValue"> NewUncleBlockFromValue </a>
creates and sets the uncle of a block based   
on an ethutil.Value object that contains that information.
- **header**: an ethutil.Value object containing the header of the uncle block.  
	The header of the uncle Block (and any Block) contains those fields:   
	_PrevHash, UncleSha, Coinbase, state, TxSha, ReceiptSha, LogsBloom,  
	Difficulty, Number, GasLimit, GasUsed, Time, Extra, Nonce_.   
	See also the <a class="link" href="#Block.HashNoNonce">HashNoNonce</a>
	function for details about these fields.<br><br>

`func (block *Block) Trie() *trie.Trie`  
<a class="function" name="Block.Trie"> Trie </a> 
returns the block's state trie<br><br>

`func (block *Block) Root() interface{}`  
<a class="function" name="Block.Root"> Root </a> 
returns the block's state root<br><br>

`func (block *Block) Diff() *big.Int`  
<a class="function" name="Block.Diff"> Diff </a> 
returns the block's difficulty.<br><br>

`func (self *Block) Receipts() []*Receipt`  
<a class="function" name="Block.Receipts"> Receipts </a> 
returns the block's receipts.<br><br>

`func (block *Block) miningHeader() []interface{}`  
<a class="function" name="Block.miningHeader"> miningHeader </a> 
returns an object that is almost the same as a Block. The differences are:
1. The Block struct also contains the uncles as a slice of Blocks and not only the uncles hash.
2. The Block struct contains the state object as a whole and not the root of the state.
3. The Block struct contains the receipts and not only the TxSha
4. The Block struct contains the transactions and not only the ReceiptSha
5. The Block struct contains the Reward to be given to the miner.  
The object returned, if appended with the Nonce field of the Block is the Block's header.<br><br>

`func (block *Block) header() []interface{}`  
<a class="function" name="Block.header"> header </a> 
returns the header of the Block. See the function 
<a class="link" href="#Block.miningHeader">miningHeader</a> 
for a detailed explanation.<br><br>

`func (block *Block) String() string`  
<a class="function" name="Block.String"> String </a> 
returns the string representation of the block.<br><br>

`func (self *Block) Size() ethutil.StorageSize`  
<a class="function" name="Block.Size"> Size </a> 
returns a float64 object representing the size of storage needed to save this
block's rlp-encoding.   
The rlp-encodable fields of a Block are it's header, transactions and uncles.<br><br>

`func (self *Block) RlpData() interface{}`  
<a class="function" name="Block.RlpData"> RlpData </a> 
returns an object containing this block's rlp-encodable fields.   
The rlp-encodable fields of a Block are it's header, transactions and uncles.   
(this is also the order in which the returned object contains them).<br><br>


`func (self *Block) N() []byte`  
<a class="function" name="Block.N"> N </a>
returns the Nonce field of the block.

<br>

###  1.2 <a class="file" name="common.go"> common.go </a>

Only contains the definition of 2 interfaces.<br>

```
type BlockProcessor interface {
	Process(*Block) (*big.Int, state.Messages, error)
}
```
Any type that implements the 
<a class="interface" name="BlockProcessor">BlockProcessor</a> interface must 
also implement the `Process`  
function which is called upon a Block and calculates
and returns in that order:  
the total difficulty of the block, the messages of the
block (aka the transactions) and/or an error.  
In case of an error, the values 
returned by the function Process for the td and messages are nil.  
In case of no errors, the error value returned by the Process function is nil.  
Only a <a class="link" href="#ChainManager">ChainManager</a> 
object implements this interface.<br><br><br>

```
type Broadcaster interface {
	Broadcast(wire.MsgType, []interface{})
}
```
Any object that implements the <a class="interface" name="Broadcaster">Broadcaster</a> interface must also implement the function 'Broadcast',  
which is used to broadcast messages of a given type to a list of recipients. 
The type of the message to   
be broadcasted, wire.MsgType, is a single byte and is defined in the package
wire (file wire/messaging.go).  
Only a <a class="link" href="#TxPool">TxPool</a>
object implements this interface.

<br>

###  1.3 <a class="file" name="derive_sha.go"> derive_sha.go </a>

<br>

```
type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}
```
Any object that implements the 
<a class="interface" name="DerivableList">DerivableList</a> interface can then
call the function <a class="link" href="#DeriveSha">DeriveSha</a>.  
Examples of objects that implement this interface are the 
<a class="link" href="#Receipt">Receipt</a> and the 
<a class="link" href="#Transaction">Transaction</a> objects.
- **Len()**: returns the number of the elements of the object 
	that implements this interface.
- **GetRlp(i int)**: returns the RLP-encoding of i-th element of the object 
	that implements this interface.<br><br>


`func DeriveSha(list DerivableList) []byte`  
<a class="function" name="DeriveSha">DerivableList</a> contructs a trie based on 
the _list_ parameter and returns the root hash of the trie.  
Examples of objects that make use of this function are the 
<a class="link" href="#Receipt">Receipt</a> and 
<a class="link" href="#Transaction">Transaction</a> objects.<br><br>

<br>

###  1.4 <a class="file" name="transaction.go">transaction.go </a>

#### 1.4.1 Data structures

```
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
```

Definition of a <a class="struct" name="Transaction"> Transaction </a> object.
A Transaction object implements the <a class="link" href="#Message">Message</a> 
interface.

- **Nonce**: A value equal to the number of transactions that originated from   
	this address or, in the case of accounts with associated code, the number   
	of contract-creations made by this account.
- **recipient**: The 160-bit address of the transaction's recipient or,  
	for a contract creation transaction, nil.
- **value**: This is equal to the number of Wei to be transferred to the   
	transaction's recipient or, in the case of contract creation, as an  
	endowment to the newly created account.
- **gas**: The gas limit equal to the maximum amount of gas that should be  
	used when executing this transaction.
- **gasPrice**: This is equal to the number of Wei to be paid per unit of gas for  
	all computation costs incurred as a result of the execution of this
	transaction.
- **data**: An unlimited size byte array specifying the input data of the
	transaction.   
	The first 4 bytes of this field specify which function to call
	when the transaction  
	will execute by using the hash of the function's name
	to be called and it's arguments.  
	The rest of the data field are the arguments
	passed to the function.  
	If the data field is empty, it means a transaction is for a payment and not 
	an execution of the contract.
- **v**, **r**, **s**: Values corresponding to the ECDSA digital signature of 
	the transaction and used to determine   
	the originating Externally Owned Account.<br><br>

```
type Transactions []*Transaction
```
The <a class="struct" name="Transactions">Transactions</a>
data type is defined as a means to 
1. sort Transactions and
2. derive the root hash of a trie consisting of a list of transactions. <br><br>

```
type TxByNonce struct{ Transactions }
```
the <a class="struct" name="TxByNonce">TxByNonce</a> data type is used
for to perform comparison operations on Transaction objects.<br><br>

<br>

#### 1.4.2 Functions<br><br>

`func IsContractAddr(addr []byte) bool`  
<a class="function" name="Transaction.IsContractAddr">IsContractAddr</a>
returns true if the provided 'addr' has a length of 0, false otherwise.<br><br>

`func NewContractCreationTx(value, gas, gasPrice *big.Int, script []byte) *Transaction`  
<a class="function" name="Transaction.NewContractCreationTx">NewContractCreationTx</a>
creates and returns a new Transaction which represents the creation of a contract.  
The Transaction's 'recipient' field will be set to nil and the Transaction's 
'data' field will be set to the 'script' parameter.  
All other Transaction fields will be set to the corresponing parameters 
passed to this function   
(except of course for the fields nonce,v,r,s which are not set)<br><br>

`func NewTransactionMessage(to []byte, value, gas, gasPrice *big.Int, data []byte) *Transaction`  
<a class="function" name="Transaction.NewTransactionMessage">NewTransactionMessage</a>
creates and returns a new Transaction.  
All Transaction fields will be set to the corresponing parameters passed to this function   
(except of course for the fields nonce,v,r,s which are not set)<br><br>

`func NewTransactionFromBytes(data []byte) *Transaction`  
<a class="function" name="Transaction.NewTransactionFromBytes">NewTransactionFromBytes</a>
creates and returns a new Transaction based on the 'data' parameter.  
The latter should be a valid rlp-encoding of the Transaction object to be created.  
(The 'data' param is converted to an ethutil.Value object first and then gets decoded)<br><br>

`func NewTransactionFromValue(val *ethutil.Value) *Transaction`  
<a class="function" name="Transaction.NewTransactionFromValue">NewTransactionFromValue</a>
creates and returns a new Transaction based on the 'val' parameter.   
The latter should be a valid rlp-encoding - of the Transaction object to be created
\- that has been cast   
to an ethutil.Value object.<br><br>

`func (tx *Transaction) Hash() []byte`  
<a class="function" name="Transaction.Hash">Hash</a>
returns the hash of the caller. The hash will be the Kecchak-256 hash of 
these Transaction fields:  
nonce, gasPrice, gas, recipient, value, data.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (self *Transaction) Data() []byte`  
<a class="function" name="Transaction.Data">Data</a> 
returns the 'data' field of the caller.<br><br>

`func (self *Transaction) Gas() *big.Int`  
<a class="function" name="Transaction.Gas">Gas</a> 
returns the 'gas' field of the caller.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (self *Transaction) GasPrice() *big.Int`  
<a class="function" name="Transaction.GasPrice">GasPrice</a> 
returns the 'gasPrice' field of the caller.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (self *Transaction) Value() *big.Int`  
<a class="function" name="Transaction.Value">Value</a>
returns the 'value' field of the caller.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (self *Transaction) Nonce() uint64`  
<a class="function" name="Transaction.Nonce">Nonce</a>
returns the 'nonce' field of the caller.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (self *Transaction) SetNonce(nonce uint64)`  
<a class="function" name="Transaction.SetNonce">SetNonce</a>
assigns the 'nonce' param to the caller's nonce field.<br><br>

`func (self *Transaction) From() []byte`  
<a class="function" name="Transaction.From">From</a>
returns the 'sender' of the caller.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (self *Transaction) To() []byte`  
<a class="function" name="Transaction.To">To</a>
returns the 'recipient' field of the caller.  
This function is implemented as part of the 
<a class="link" href="#Message">Message</a> 
interface<br><br>

`func (tx *Transaction) Curve() (v byte, r []byte, s []byte)`  
<a class="function" name="Transaction.Curve">Curve</a>
returns v, r, s of the Transaction. r and s are left-padded to 32 bytes.   
For more, see the function ethutil.leftPadBytes<br><br>

`func (tx *Transaction) Signature(key []byte) []byte`  
<a class="function" name="Transaction.Signature">Signature</a>
calculates and returns the signature of the hash of a transaction with the param key.  
The signature used to be obtained through the 
github.com/obscuren/secp256k1-go package repo   
but this package no longer exists. One may now use the function Sign from:  
https://github.com/ethereum/go-ethereum/blob/master/crypto/secp256k1/secp256.go<br><br>

`func (tx *Transaction) PublicKey() []byte`  
<a class="function" name="Transaction.PublicKey">PublicKey</a>
retrieves and returns the public key of the transaction.  
The public key used to be obtained through the github.com/obscuren/secp256k1-go  
package repo but this package no longer exists. 
One may now use the function RecoverPubKey from:  
https://github.com/ethereum/go-ethereum/blob/master/crypto/secp256k1/secp256.go<br><br>

`func (tx *Transaction) Sender() []byte`  
<a class="function" name="Transaction.Sender">Sender</a>
returns the sender of the transaction. To do so, the public key of the transaction
is first   
retrieved through the function <a class="link" href="#Transaction.PublicKey">PublicKey</a>.
If the public key passes validation then the last 12  
bytes of the public key 
are returned (aka the sender address).<br><br>

`func (tx *Transaction) Sign(privk []byte) error`  
<a class="function" name="Transaction.Sign">Sign</a> signes the transaction.  
To do so, the function <a class="link" href="#Transaction.Signature">Signature</a>
is called (which makes use of a non-existent package, refer to it for more).  
After the transaction has been signed the r,s,v fields of the Transaction are set
to the signatures appropriate fields.  
r is set to the first 32 bytes of thTransaction.e signature.  
s is set to the 32-th up to and including the 63-th bytes of the signature.  
v is set to the sum of the last bit of the transaction and the number 27.<br><br>

`func (tx *Transaction) RlpData() interface{}`  
<a class="function" name="Transaction.RlpData">RlpData</a>
returns the rlp-encodable fields of the caller 
(all fields of a <a class="link" href="#Transaction">Transaction</a> object).<br><br>

`func (tx *Transaction) RlpValue() *ethutil.Value`  
<a class="function" name="Transaction.RlpValue">RlpValue</a>
casts and returns the caller to an ethutil.Value object.<br><br>

`func (tx *Transaction) RlpEncode() []byte`  
<a class="function" name="Transaction.RlpEncode">RlpEncode</a> returns the
rlp-encoding of the cast of the caller object to an ethutil.Value object.   
The cast is done through the 
<a class="link" href="#Transaction.RlpValue">RlpValue</a> function.<br><br>

`func (tx *Transaction) RlpDecode(data []byte)`  
<a class="function" name="Transaction.RlpDecode">RlpDecode</a>
sets the caller's fields equal to the rlp-decoding of the 'data' parameter.   
The _data_ param must be a valid bytes representation of a 
<a class="link" href="#Transaction">Transaction</a> 
that has been cast to an ethutil.Value object.<br><br>

`func (tx *Transaction) RlpValueDecode(decoder *ethutil.Value)`  
<a class="function" name="Transaction.RlpValueDecode">RlpValueDecode</a>
sets the caller's fields equal to the rlp-decoding of the 'decoder'.   
The _decoder_ must be a valid bytes representation of a 
<a class="link" href="#Transaction">Transaction</a>
that has been cast to an ethutil.Value object.<br><br>

`func (tx *Transaction) String() string`  
<a class="function" name="Transaction.String">String</a>
returns the string representation of the caller.<br><br>

`func (self Transactions) RlpData() interface{}`  
<a class="function" name="Transaction.RlpData">RlpData</a>
returns the rlp-encodable fields of all the 
<a class="link" href="#Transaction">Transaction</a> objects 
that the caller consists of.<br><br>

`func (s Transactions) Len() int`  
<a class="function" name="Transaction.Len">Len</a>
returns the number of the Transaction objects that the caller consists of.  
This function, along with the <a class="link" href="#Transaction.GetRlp">GetRlp</a> function 
are implemented as parts of the <a class="link" href="#DerivableList">DerivableList</a> interface  
in order to be able to call the <a class="link" href="#DeriveSha">DeriveSha</a> function on 
a Transactions object, which would return the hash  
of the root of the trie consisted of
the Transaction objects that are included in the Transactions object.<br><br>

`func (s Transactions) GetRlp(i int) []byte`  
<a class="function" name="Transaction.GetRlp">GetRlp</a>
returns the rlp-encoding of the i-th Transaction of the caller.
This function, along with the <a class="link" href="#Transaction.Len">Len</a> function  
are implemented as parts of the <a class="link" href="#DerivableList">DerivableList</a> interface
in order to be able to call the <a class="link" href="#DeriveSha">DeriveSha</a> function on   
a Transactions object, which would return the hash of the root of the trie consisted of
the Transaction objects  
that are included in the Transactions object.<br><br>

`func (s TxByNonce) Less(i, j int) bool`  
<a class="function" name="Transaction.Less">Less</a>
returns true if the i-th Transaction of a list of Transactions is less
than the j-th Transaction.  
This is determined by only comparing the two
Transactions' nonces.<br><br>

<br>

###  1.5 <a class="file" name="receipt.go">receipt.go </a>

#### 1.5.1 Data structures

<br>

```
type Receipt struct {
	PostState         []byte
	CumulativeGasUsed *big.Int
	Bloom             []byte
	logs              state.Logs
}
```
The <a class="struct" name="Receipt">Receipt</a> data structure represents
the receipts field of a <a class="link" href="#Transaction">Transaction</a>.   
The receipt contains certain information regarding the execution of the transaction  
and are mainly used for providing logs and events regarding the transaction.  
Each receipt is placed on a trie.
- **PostState**: the world state after the execution of the transaction 
	this receipt refers to.
- **CumulativeGasUsed**: The cumulative gas used up to and including the
	current transaction of a Block.
- **Bloom**: A 'bloom' is composed from the log entries of the receipt. 
	This field is the 2048-bit hash of that bloom.  
	(See file <a class="link" href="#bloom9.go">bloom9.go</a>)
- **logs**: A series of log entries. Each entry is a tuple of the logger’s address, a
	possibly empty series of 32-byte log   
	topics, and some number of bytes of data. 
	The logs are created through the execution of a transaction.<br><br><br>


```
type Receipts []*Receipt
```
The <a class="struct" name="Receipts">Receipts</a>
data type is defined as a means to 
1. sort Receipt objects and
2. derive the root hash of a trie consisting of a list of receipts. <br><br>

<br>

#### 1.5.2 Functions

<br>

`func NewReceipt(root []byte, cumalativeGasUsed *big.Int) *Receipt`  
<a class="function" name="NewReceipt">NewReceipt</a> creates a new 
<a class="link" href="#Receipt">Receipt</a> provided the 'root' of 
the receipts trie and the cumulative gas used, and returns it.<br><br>

`func NewRecieptFromValue(val *ethutil.Value) *Receipt`  
<a class="function" name="NewRecieptFromValue">NewRecieptFromValue</a> creates 
a new <a class="link" href="#Receipt">Receipt</a> from 'val' and returns it.   
The val param should be the bytes representation of the rlp-encoded
Receipt object to be created. <br><br>

`func (self *Receipt) SetLogs(logs state.Logs)`  
<a class="function" name="Receipt.SetLogs">SetLogs</a> 
sets the caller's logs field equal to the 'logs' parameter.<br><br>

`func (self *Receipt) RlpValueDecode(decoder *ethutil.Value)`  
<a class="function" name="Receipt.RlpValueDecode">RlpValueDecode</a> 
sets the caller's (<a class="link" href="#Receipt">Receipt</a>)
fields to the rlp-decoded fields of the decoder parameter.<br><br>

`func (self *Receipt) RlpData() interface{}`  
<a class="function" name="Receipt.RlpData">RlpData</a> 
returns the rlp-encodable fields of the caller. The only difference between what
this function returns   
and the caller's fields is that instead of the logs field
this function returns the result of the function call logs.RlpData().<br><br>

`func (self *Receipt) RlpEncode() []byte`  
<a class="function" name="Receipt.RlpEncode">RlpEncode</a> 
Returns the RLP-encoding of the rlp-encodable fields of a Receipt object.   
These fields are obtained through the function 
<a class="link" href="#Receipts.RlpData">RlpData</a><br><br>

`func (self *Receipt) Cmp(other *Receipt) bool`  
<a class="function" name="Receipt.Cmp">Cmp</a>
returns true if the caller Receipt is the same as the 'other', false otherwise.  
Two Receipt objects are the same if their PostState fields are equal.<br><br>

`func (self *Receipt) String() string`  
<a class="function" name="Receipt.String">String</a>
returns the string representation of the caller.<br><br>

`func (self Receipts) RlpData() interface{}`  
<a class="function" name="Receipts.RlpData">RlpData</a> returns the rlp-encodable
fields of the caller, aka an array that is populated with the result of the   
consecutive calls to <a class="link" href="#Receipt.RlpData">RlpData</a> function 
for each of the Receipt objects that the Receipts caller object consists of.<br><br>

`func (self Receipts) RlpEncode() []byte`  
<a class="function" name="Receipts.RlpEncode">RlpEncode</a> returns the rlp-encoding
of the caller, aka rlp-encodes each Receipt object that the caller consists  
of and returns the result.<br><br>

`func (self Receipts) Len() int`  
<a class="function" name="Receipts.Len">Len</a> returns the number of Receipts that
the caller consists of.  
This function, along with the <a class="link" href="#Receipts.GetRlp">GetRlp</a> 
function are implemented  as parts of the 
<a class="link" href="#DerivableList">DerivableList</a>  
interface in order to be able to call the <a class="link" href="#DeriveSha">DeriveSha</a>
function on a Receipts object, which would  
return the hash of the root of the trie
consisted of the Receipt objects that are included in the  
Receipts object.<br><br>

`func (self Receipts) GetRlp(i int) []byte`  
<a class="function" name="Receipts.GetRlp">GetRlp</a> returns the rlp-encoding of
the i-th Receipt of the caller.  
This function, along with the 
<a class="link" href="#Receipts.Len">Len</a> function are implemented as parts of the 
<a class="link" href="#DerivableList">DerivableList</a>  
interfacein order to be able to call the <a class="link" href="#DeriveSha">DeriveSha</a>
function on a Receipts object, which would  
return the hash of the root of the trie consisted of the Receipt objects 
that are included in the  
Receipts object.<br><br>

<br>


###  1.6 <a class="file" name="bloom9.go">bloom9.go </a>

This file includes functions regarding the logsBloom field of a Block.  
Every block contains many transactions. Each transaction has a receipt. Each receipt  
has  many logs. Each log consists of 3 fields: the address, the topics and the data.  
The LogsBloom will be the lower 11 bits of the modulo operation by 2048 of the expression   
kecchak-256(address and logs of all the receipts of the block's transactions).  

"Events in the ethereum system must be easily searched for, so that applications can   
filter and display events, including historical ones, without undue overhead. At the  
same time, storage space is expensive, so we don't want to store a lot of duplicate data  
\- such as the list of transactions, and the logs they generate. The logs bloom filter  
exists to resolve this.

When a block is generated or verified, the address of any logging contract,  
and all the indexed fields from the logs generated by executing those transactions  
are added to a bloom filter, which is included in the block header. The actual logs  
are not included in the block data, to save space.

When an application wants to find all the log entries from a given contract, or with  
specific indexed fields (or both), the node can quickly scan over the header of each  
block, checking the bloom filter to see if it may contain relevant logs. If it does,  
the node re-executes the transactions from that block, regenerating the logs, and   
returning the relevant ones to the application."

Source:  
https://ethereum.stackexchange.com/questions/3418/how-does-ethereum-make-use-of-bloom-filters

A thorough explanation on bloom filters:  
https://brilliant.org/wiki/bloom-filter/

A complete example of the bloom filters on Ethereum:  
https://medium.com/coinmonks/ethereum-under-the-hood-part-8-blocks-2-5fba93293213<br>

<br>

#### 1.6.1 Functions

`func CreateBloom(receipts Receipts) []byte`  
<a class="function" name="CreateBloom">CreateBloom</a> creates the LogsBloom field 
of the given 'receipts', aka the LogsBloom field of a
<a class="link" href="#Block">Block</a> and returns it.  
For each of the receipts, a logical OR operation happens between the so-far 
result and the result of a call to the   
<a class="link" href="#LogsBloom">LogsBloom</a> 
function call on the current
receipt. Finally, the overall result is left-padded by 64 bytes and returned.<br><br>

`func LogsBloom(logs state.Logs) *big.Int`  
<a class="function" name="LogsBloom">LogsBloom</a> returns the LogBloom of the param
logs, aka the LogBloom of a receipt, as a big.Int.  
Steps to do so are as follows:
1. Iterate over the logs field of the receipt.
2. Create an array for every log. The array's first element will be the address of the log.   
All other elements will be the topics of the log.
3. Initialize the result to be returned to 0 and iterate over that array. 
4. result = (previous_result) | (lower 11 bits of kecchak-256(current_element_of_array)  %  2048)  
where the | operator is a logical OR operator.<br><br>

`func bloom9(b []byte) *big.Int`  
<a class="function" name="bloom9">bloom9</a> is an inner function that 
returns the lower 11 bits of b % 2048.<br><br>

`func BloomLookup(bin, topic []byte) bool`  
<a class="function" name="BloomLookup">BloomLookup</a> checks whether the bin bloom
has it's bit field of the 'topic' set to 1, aka the bloom **may** contain the 'topic'.  
This function is only called by the 
<a class="link" href="#bloomFilter">bloomFilter</a> function defined in 
the file core/filter.go. Even if the bit is set to 1 that  
does not mean that the bin bloom will contain the topic. A bloom filter is like a 
hash table where instead of storing the values 
of the keys, a bit set to 1 indicates that there is a value with the specified key.
That means that multiple values with different 
keys are mapped to the same bit field of the bloom filter. 
Only if that bit field is 0 we know that the bloom does **not** 
contain that value pair.<br><br>

<br>

## 2. Package <a class="file" name="core"> core </a>

###  2.1 <a class="file" name="simple_pow.go">simple_pow.go </a>

This file does not contain any code.<br><br>

### 2.2 <a class="file" name="fees.go">fees.go </a>

Only contains the definition of the following variable:  

`var BlockReward *big.Int = big.NewInt(1.5e+18)`  
The <a class="struct" name="BlockReward">BlockReward</a>
represents the initial block reward for miners, set to 1.5 Ether.  
Only the function 
<a class="link" href="#BlockManager.AccumelateRewards">AccumelateRewards</a> 
defined in the file <a class="link" href="#block_manager.go">block_manager.go</a>
uses this variable.<br><br>

<br>

### 2.3 <a class="file" name="genesis.go">genesis.go </a>


#### 2.3.1 Data Structures

`var ZeroHash256 = make([]byte, 32)`  
Used to set the parent's hash of the genesis block.<br><br>

`var ZeroHash160 = make([]byte, 20)`  
Used to set the coinbase of the genesis block.<br><br>

`var ZeroHash512 = make([]byte, 64)`  
Used to set the bloom field of the genesis block.<br><br>

`var EmptyShaList = crypto.Sha3(ethutil.Encode([]interface{}{}))`  
Used to set the root state of the genesis block.<br><br>

`var EmptyListRoot = crypto.Sha3(ethutil.Encode(""))`  
Used to set the tx root and receipt root of the genesis block.<br><br>


```
var GenesisHeader = []interface{}{
	ZeroHash256,
	EmptyShaList,
	ZeroHash160,
	EmptyShaList,
	EmptyListRoot,
	EmptyListRoot,
	ZeroHash512,
	big.NewInt(131072),
	ethutil.Big0,
	big.NewInt(1000000),
	ethutil.Big0,
	ethutil.Big0,
	nil,
	crypto.Sha3(big.NewInt(42).Bytes()),
}
```
This is the special genesis block. The fields set are the same as those of a 
<a class="link" href="#Block">Block</a>.
- **ZeroHash256**: Previous hash (none)
- **EmptyShaList**: Empty uncles
- **ZeroHash160**: Coinbase, (aka miner's address)
- **EmptyShaList**: Root state
- **EmptyListRoot**: tx root
- **EmptyListRoot**: receipt root
- **ZeroHash512**: bloom field
- **big.NewInt(131072)**: difficulty
- **ethutil.Big0**: block number
- **big.NewInt(1000000)**: Block upper gas bound
- **ethutil.Big0**: Block gas used
- **ethutil.Big0**: Time field
- **nil**: Extra field.
- **crypto.Sha3(big.NewInt(42).Bytes())**: Nonce field<br><br>

`var Genesis = []interface{}{GenesisHeader, []interface{}{}, []interface{}{}}`  
<a class="struct" name="Genesis">Genesis</a> 
is defined so as to be able to RLP-encode the genesis block.  
To rlp-encode any Block we provide it's header, transactions and uncles.   
The genesis block does not have any transactions and uncles, thus the use 
of the 2 empty objects.  
This variable is only used by the function
<a class="link" href="#NewChainManager">NewChainManager</a> defined in the file
<a class="link" href="#chain_manager.go">chain_manager.go</a>.  
For example, the function 
<a class="link" href="#NewChainManager">NewChainManager</a> creates a new <a class="link" href="#ChainManager">ChainManager</a> object for which there  
is a field that represents the genesis block. To set that block, 
the NewChainManager function uses  
the ethutil.Encode(Genesis) function.<br><br>

<br>

### 2.4 <a class="file" name="block_manager.go">block_manager.go</a>

#### 2.4.1 Data Structures

`var statelogger = logger.NewLogger("BLOCK")`  
<a class="struct" name="statelogger">statelogger</a> is a channel 
used to log messages regarding Blocks processing.<br><br>

```
type Peer interface {
	Inbound() bool
	LastSend() time.Time
	LastPong() int64
	Host() []byte
	Port() uint16
	Version() string
	PingTime() string
	Connected() *int32
	Caps() *ethutil.Value
}
```
This interface is only implemented by a Peer object defined in the 'eth' package (file peer.go)
- **Inbound()**: Determines whether it's an inbound or outbound peer
- **LastSend()**: Last known message send time.
- **LastPong()**: Last received pong message
- **Host()**: the host.
- **Port()**: the port of the connection
- **Version()**: client identity
- **PingTime()**: Used to give some kind of pingtime to a node, not very accurate.
- **Connected()**: Flag for checking the peer's connectivity state
- **Caps()**: getter for the protocolCaps field of a Peer object.<br><br>

```
type EthManager interface {
	BlockManager() *BlockManager
	ChainManager() *ChainManager
	TxPool() *TxPool
	Broadcast(msgType wire.MsgType, data []interface{})
	PeerCount() int
	IsMining() bool
	IsListening() bool
	Peers() *list.List
	KeyManager() *crypto.KeyManager
	ClientIdentity() wire.ClientIdentity
	Db() ethutil.Database
	EventMux() *event.TypeMux
}
```
The <a class="interface" name="EthManager">EthManager</a> 
interface is only implemented by the 
'Ethereum' object, defined in the package 'eth' (file ethereum.go)
- **BlockManager**: Getter for the <a class="link" href="#BlockManager">
BlockManager</a> field of the Ethereum object
- **ChainManager**: Getter for the <a class="link" href="#ChainManager">
ChainManager</a> field of the Ethereum object
- **TxPool**: Getter for the <a class="link" href="#TxPool">
TxPool</a> field of the Ethereum object.
- **Broadcast**: Used to broacast the 'data' msg to all peers. 
- **IsMining**: Returns whether the peer is mining.
- **IsListening**: Returns whether the peer is listening.
- **Peers**: Returns all connected peers.
- **KeyManager**: Getter for the KeyManager field of the Ethereum object
- **ClientIdentity**: Getter for the ClientIdentity field of the Ethereum object.  
	That field is mainly used for communication between peers.
- **Db**: Getter for the ethereum database of the Ethereum object.  
That field represents the World State.
- **EventMux**: Getter for the EventMux field of the Ethereum object.   
That field is mainly used to dispatch events to registered nodes.<br><br>

```
type BlockManager struct {
	mutex              sync.Mutex
	bc                 *ChainManager
	mem                map[string]*big.Int
	Pow                pow.PoW
	txpool             *TxPool
	lastAttemptedBlock *types.Block
	events             event.Subscription
	eventMux           *event.TypeMux
}
```
The <a class="struct" name="BlockManager">BlockManager</a> data
structure is essentially a state manager for processing new blocks 
and managing the overall state.
- **mutex**: Mutex for locking the block processor as Blocks can 
only be handled one at a time
- **bc**: Canonical block chain
- **mem**: non-persistent key/value memory storage
- **Pow**: Proof of work used for validating
- **txpool**: The transaction pool. See <a class="link" href="#TxPool">TxPool</a>.
- **lastAttemptedBlock**: The last attempted block is mainly used for debugging purposes.  
	This does not have to be a valid block and will be set during 'Process' & canonical validation.
- **events**: Provides a way to subscribe to events.
- **eventMux**: event mutex, used to dispatch events to subscribers.<br><br>

<br>

#### 2.4.2 Functions

`func NewBlockManager(txpool *TxPool, chainManager *ChainManager, eventMux *event.TypeMux) *BlockManager`  
<a class="function" name="NewBlockManager">NewBlockManager</a> creates a new 
BlockManager object by initializing these fields of a BlockManager object type:   
_mem, Pow, bc, eventMux, txpool_. The Pow object will have it's _turbo_ field 
set to true when created.  
See the type 'EasyPow' of the 'ezp' package for more (file pow/ezp/pow.go)<br><br>

`func (sm *BlockManager) TransitionState(statedb *state.StateDB, parent, block *types.Block) (receipts types.Receipts, err error)`  
<a class="function" name="BlockManager.TransitionState">TransitionState</a>
returns (receipts, nil) or (nil, error) if an 
<a class="link" href="#GasLimitErr">GasLimitErr</a> error has occured.  
This function along with the 
<a class="link" href="#BlockManager.ApplyTransactions">ApplyTransactions</a> 
function form a recursion algorithm meant to  
apply all transactions of the current block to the world state. The main logic/
application happens  
in the latter function. However, this function is the one
to be called when we want to apply the  
transactions of a block. The 
<a class="link" href="#BlockManager.TransitionState">TransitionState</a> 
function sets the total gas pool (amount of gas left)  
for the coinbase address of the block before calling the 
<a class="link" href="#BlockManager.ApplyTransactions">ApplyTransactions</a> 
function which in turn  
will call the <a class="link" href="#BlockManager.TransitionState">TransitionState</a> 
function again, and so on, until all transactions have been applied or a  
<a class="link" href="#GasLimitErr">GasLimitErr</a> error has occured. 
If no such an error has occured then the 'receipts' object returned by  
this function will hold the receipts that are
the result of the application of the transactions of the _block_ param.<br><br>

```
func (self *BlockManager) ApplyTransactions(
		coinbase *state.StateObject, 
		state *state.StateDB, 
		block *types.Block, 
		txs types.Transactions, 
		transientProcess bool
	)
	(
		types.Receipts, 
		types.Transactions, 
		types.Transactions, 
		types.Transactions, 
		error
	)
```
The
<a class="function" name="BlockManager.ApplyTransactions">ApplyTransactions</a>
function will apply the transactions of a block to the world state and return
the results as a tuple.   
It gets called by the 
<a class="link" href="#BlockManager.TransitionState">TransitionState</a> 
function, then calls the latter again, and so on, to form a recursion
algorithm that will apply the transactions one by one. In case where a
<a class="link" href="#GasLimitErr">GasLimitErr</a>
error occurs during the application of any transaction, the processing of the transactions stops.

Returns: (receipts, handled, unhandled, erroneous, err)

- **receipts**: The receipts up to but not including any transaction that has 
	caused a <a class="link" href="#GasLimitErr">GasLimitErr</a> error.
- **handled**: All transactions that were handled up to but not including 
	any transaction that has caused a 
	<a class="link" href="#GasLimitErr">GasLimitErr</a> error.
- **unhandled**: In case of a <a class="link" href="#GasLimitErr">GasLimitErr</a>
	error this object will contain all the transactions that were not applied   
	(includes the transaction that caused the error). Otherwise, this object will be nil.
- **erroneous**: Any transactions that caused an error other than a
	<a class="link" href="#GasLimitErr">GasLimitErr</a>
	and/or a <a class="link" href="#NonceErr">NonceErr</a> errors.
- **err**: The err will be either a 
<a class="link" href="#GasLimitErr">GasLimitErr</a> error or nil.

A short description on what this function does follows.
1. Clear all state logs.
2. Get (or create a new) coinbase state object and call the 
	<a class="link" href="#BlockManager.TransitionState">TransitionState</a> func.
3. The latter function will call this function again, forming the recursion.
4. If an error occured and is a 
	<a class="link" href="#GasLimitErr">GasLimitErr</a>
	error then stop the process and set the _unhandled_ variable
	(to be returned later). If it is a <a class="link" href="#NonceErr">NonceErr</a> error, 
	ignore it. If it is any other error, also ignore it, but also append to the variable _erroneous_   
	(to be returned 
	later) the transaction that caused that error.

5. Calculate the gas used so far and the current reward for the miner. 
	Update the state.
6. Create the receipt of the current transaction and set the receipt's 
	logs and Bloom field.
7. If the parameter _transientProcess_ is false, notify all subscribers 
	about the transaction.
8. Append receipt, transaction to the _receipts_, _handled_ variables 
	(to be returned later)
9. When the processing has ended, set the block's reward and totalUsedGas fields.
9. Return the results.<br><br>


`func (sm *BlockManager) Process(block *types.Block) (td *big.Int, msgs state.Messages, err error)`  
<a class="function" name="BlockManager.Process">Process</a>
processes a block. When successful, returns the return result of a call to the 
<a class="link" href="#BlockManager.ProcessWithParent">ProcessWithParent</a> 
function.  
Otherwise, in case that the hash of the block or the hash of the parent of the 
block already exist in the ChainManager,  
returns the tuple (nil, nil KnownBlockError) or (nil, nil, ParentError) accordingly.
Before calling the 
<a class="link" href="#BlockManager.ProcessWithParent">ProcessWithParent</a>  
function, <a class="link" href="#BlockManager.Process">Process</a>
takes care of locking the BlockManager with a mutex and only after the
<a class="link" href="#BlockManager.ProcessWithParent">ProcessWithParent</a>  
function has returned the <a class="link" href="#BlockManager">BlockManager</a> 
is unlocked.<br><br>


`func (sm *BlockManager) ProcessWithParent(block, parent *types.Block) (td *big.Int, messages state.Messages, err error)`  
<a class="function" name="BlockManager.ProcessWithParent">ProcessWithParent</a>
is the main process function of a block. Gets called by the function
<a class="link" href="#BlockManager.Process">Process</a>.

Returns a tuple (td, messages, error) where 
- **td**: total difficulty of the processed _block__
- **messages**: messages generated by the application of the transactions of the _block__  
	(and some extra messages like transactions regarding rewarding miners).
- If an error occured the returned tuple becomes (nil, nil, error).  
	If no errors have occured the returned tuple becomes (td, messages, nil).

Below follows a short description of what this function does:
1. Saves a copy of the state. Also queues a reset of the state after the function's  
	return based on that copy. (use of golang's _defer_)
2. Validates the _block__ with a call to the function 
	<a class="link" href="#BlockManager.ValidateBlock">ValidateBlock</a>. 
	If errors happened, returns.
3. Calls TransitionState to attempt to do the state transition. 
	If errors, returns.
4. Creates the bloom field of the receipts returned from step 2. 
	If for some reason the bloom  
	field is different from the bloom field of the provided _block__, it returns.
5. Validates the transactions and the receipts root hashes. 
	If errors, returns.
6. Calls <a class="link" href="#BlockManager.AccumelateRewards">AccumelateRewards</a>
	to calculate the miner rewards. If errors, returns.
7. Sets the state to 0 and makes a call to CalculateTD in order to 
	calculate the total difficulty of the _block__.  
	If errors, returns. If not, the last step is to remove the _block__'s transactions from 
	the BlockManager's txpool,  
	sync the state db, cancel the queued state reset, send a message to 
	the chainlogger channel and finally  
	return the tuple (td, messages, nil).<br><br>


`func (sm *BlockManager) CalculateTD(block *types.Block) (*big.Int, bool)`  
<a class="function" name="BlockManager.CalculateTD">BlockManager.CalculateTD</a>
calculates the total difficulty for a given block. If the calculated total difficulty  
is greater than the previous, the tuple (total_difficulty, true) is returned. 
Otherwise, the tuple (nil, false) is returned.  
TD(genesis_block)=0 and TD(_block_)=TD(_block_.parent) + sum(u.difficulty for u in _block_.uncles) + _block_.difficulty.<br><br>

`func (sm *BlockManager) ValidateBlock(block, parent *types.Block) error`  
<a class="function" name="BlockManager.ValidateBlock">BlockManager.ValidateBlock</a>
validates the current _block_.  
Returns an error if the _block_ was invalid, an uncle or 
anything that isn't on the current block chain.  
Validation validates easy over difficult (dagger takes longer time = difficult)<br><br>

`func (sm *BlockManager) AccumelateRewards(statedb *state.StateDB, block, parent *types.Block) error`  
<a class="function" name="BlockManager.AccumelateRewards">BlockManager.AccumelateRewards</a>
calculates the reward of the miner.  
Returns an error if an error has occured during the
validation process.  
If no errors have occured, nil is returned. More specifically an error is returned if:
1. the _parent_ of any of the uncles of the _block_ is nil, or
2. the (block) number of the _parent_ of any of the uncles of the _block_ and  
	the _block_ itself have a difference greater than 6, or
3. the hash of any of the uncles of the param _block_ matches any of the 
	uncles of the param _parent_.
4. the nonce of any of the uncles of the param _block_ is included in the nonce 
	of the _block_

The reward to be appointed to the miner will be:
- If the _block_ has 1 uncle:  
	r1 = <a class="link" href="#BlockReward">BlockReward</a> + <a class="link" href="#BlockReward">BlockReward</a>/32,
- If the _block_ has 2 uncles:  
	r2 = r1 + r1/32, etc.

Finally, a message is added to the _state_ manifest regarding the value to 
be transferred to the miner address.  
This value will be the sum of the above 
calculated reward and the _block_.Reward field.<br><br>

`func (sm *BlockManager) GetMessages(block *types.Block) (messages []*state.Message, err error)`  
<a class="function" name="BlockManager.GetMessages">BlockManager.GetMessages</a>
returns either the tuple (_state_.Manifest().Messages, nil) or (nil, error).

If an error is returned it will be a 
<a class="link" href="#ParentError">ParentError</a> regarding the parent of the _block_
(the error includes the hash of the parent of the _block_). This error happens in the case where the the hash of the parent of the _block_ already exists. In essence, the _state_ manifest's messages 
are the transactions that occured during the world state transition of the addition of a _block_. To get those messages a simple trick is used:  a deferred call on _state_.Reset() is queued and only then a call of the function 
<a class="link" href="#BlockManager.TransitionState">TransitionState</a> and following that a call on 
<a class="link" href="#BlockManager.AccumelateRewards">AccumelateRewards</a> is made.<br><br>

<br>

### 2.5 <a class="file" name="chain_manager.go">chain_manager.go </a>

#### 2.5.1 Data Structures

`var chainlogger = logger.NewLogger("CHAIN")`  
<a class="struct" name="chainlogger">chainlogger</a> is a channel 
used to log messages regarding the chain.<br><br>

```
type ChainManager struct {
	processor       types.BlockProcessor
	eventMux        *event.TypeMux
	genesisBlock    *types.Block
	mu              sync.RWMutex
	td              *big.Int
	lastBlockNumber uint64
	currentBlock    *types.Block
	lastBlockHash   []byte
	transState      *state.StateDB
}
```
The <a class="struct" name="ChainManager">ChainManager</a> is mainly 
used for the creation of the
<a class="link" href="#Genesis">Genesis</a> or any other blocks.
- **processor**: An interface. A neat way of calling the 
	function Process of a BlockManager object.
- **eventMux**: Used to dispatch events to subscribers.
- **genesisBlock**: The special genesis block.
- **mu**: a mutex for the ChainManager object.
- **td**: the total difficulty.  
	TD(genesis_block) = 0 and TD(B) = TD(B.parent) + sum(u.difficulty for u in B.uncles) + B.difficulty
- **lastBlockNumber**: the last block's number. 
	(the last successfully inserted block on the chain)
- **currentBlock**: During the creation of a new block, the 
	currentBlock will point to the parent  
	of the block to be created.
- **lastBlockHash**: the last block's hash.
- **transState**: represents the world state.<br><br>

<br>

#### 2.5.2 Functions

`func AddTestNetFunds(block *types.Block)`  
<a class="function" name="AddTestNetFunds">AddTestNetFunds</a>
sets the following accounts with a balance of  
1606938044258990275541962092341162602522202 Ether for testing:

- 1ba59315b3a95761d0863b05ccc7a7f54703d99
- 4157b34ea9615cfbde6b4fda419828124b70c78
- 9c015918bdaba24b4ff057a92a3873d6eb201be
- c386a4b26f73c802f34673f7248bb118f97424a
- d2a3d9f938e13cd947ec05abc7fe734df8dd826
- ef47100e0787b915105fd5e3f4ff6752079d5cb
- 6716f9544a56c530d868e4bfbacb172315bdead
- a26338f0d905e295fccb71fa9ea849ffa12aaf4

This function gets only called indirectly from a 
<a class="link" href="#ChainManager">ChainManager</a> object, and 
more specifically from the functions 
<a class="link" href="#ChainManager.Reset">Reset</a> and 
<a class="link" href="#ChainManager.setLastBlock">setLastBlock</a>.  
The latter only calls this function if the chain has 0 blocks so far. 
As one might have guessed already, when this functions gets called  
the _block_ parameter is always set to the 
<a class="link" href="#Genesis">Genesis</a> variable.<br><br>

`func CalcDifficulty(block, parent *types.Block) *big.Int`  
<a class="function" name="CalcDifficulty">CalcDifficulty</a>
calculates the difficulty of a _block_ and returns it. If the _block_ 
was mined in less than 5 seconds, the difficulty of the block   
is increased by 1/1024th of the _parent_'s difficulty. 
If the _block_ was mined in more than 5 seconds, the difficulty is
decreased by 1/1024th  
of the _parent_'s difficulty.<br><br>

`func (self *ChainManager) Td() *big.Int`  
<a class="function" name="ChainManager.Td">Td</a> 
returns the total difficulty.  
TD(genesis_block) = 0 and TD(B) = TD(B.parent) + sum(u.difficulty for u in B.uncles) + B.difficulty<br><br>

`func (self *ChainManager) LastBlockNumber() uint64`  
<a class="function" name="ChainManager.LastBlockNumber">LastBlockNumber</a>
returns the last block number.<br><br>

`func (self *ChainManager) LastBlockHash() []byte`  
<a class="function" name="ChainManager.LastBlockHash">LastBlockHash</a>
returns the last block's hash.<br><br>

`func (self *ChainManager) CurrentBlock() *types.Block`  
<a class="function" name="ChainManager.CurrentBlock">CurrentBlock</a>
returns the current block.<br><br>


`func NewChainManager(mux *event.TypeMux) *ChainManager`  
<a class="function" name="NewChainManager">NewChainManager</a> creates 
and returns a new ChainManager object by setting the genesis block  
and the eventMux field of the ChainManager. When creating the 
<a class="link" href="#Genesis">Genesis</a> block - or any other block -  
the RLP-encoding of that block is used and this is the exact purpose of
the <a class="link" href="#Genesis">Genesis</a> variable:  
it represents the RLP-encodable fields of the genesis block.<br><br>


`func (self *ChainManager) SetProcessor(proc types.BlockProcessor)`  
<a class="function" name="ChainManager.SetProcessor">SetProcessor</a>
sets the processor field of the caller.<br><br>


`func (self *ChainManager) State() *state.StateDB`  
<a class="function" name="ChainManager.State">State</a>
returns the world state. 
Access on the current state happens through the  
CurrentBlock field of the ChainManager object that 
made the call to this function.<br><br>


`func (self *ChainManager) TransState() *state.StateDB`  
<a class="function" name="ChainManager.TransState">TransState</a>
returns the transState field of the caller.<br><br>


`func (bc *ChainManager) setLastBlock()`  
<a class="function" name="ChainManager.setLastBlock">setLastBlock</a> 
is an inner function, that gets called by the 
<a class="link" href="#NewChainManager">NewChainManager</a> function.  
If the chain has 0 blocks so far, 
<a class="link" href="#ChainManager.setLastBlock">setLastBlock</a>
makes a call to the function 
<a class="link" href="#AddTestNetFunds">AddTestNetFunds</a>   
and also sets the currentBlock, lastBlockHash, lastBlockNumber 
and td fields of the  ChainManager.   
Otherwise it makes a call to the 
<a class="link" href="#ChainManager.Reset">Reset</a> function. 
In all cases, a message is sent to the chainlogger  
channel which logs the lastBlockNumber and currentBlock.Hash 
fields of the ChainManager  
after the processing has happened.<br><br>


`func (bc *ChainManager) NewBlock(coinbase []byte) *types.Block`  
<a class="function" name="ChainManager.NewBlock">NewBlock</a>
provides a way of creating a new block through a ChainManager 
object. It creates a new  
block by making a call to the function 
<a class="link" href="#CreateBlock">CreateBlock</a>,
sets the created block's difficulty, number and  
gaslimit and returns the block. _coinbase_ is the block's 
beneficiary address.  The ChainManager remains  
locked starting with the call of this function 
and until it has returned.<br><br>


`func (bc *ChainManager) Reset()`  
<a class="function" name="ChainManager.Reset">Reset</a>
resets the chain to the point where the chain will only contain
the genesis block. This includes  
the call of the function 
<a class="link" href="#AddTestNetFunds">AddTestNetFunds</a>. 
The
ChainManager remains locked starting with the call of  
this function and until it has returned.<br><br>


`func (self *ChainManager) Export() []byte {`  
<a class="function" name="ChainManager.Export">Export</a>
returns the RLP-encoding of all blocks of the chain. A message is 
sent to the <a class="link" href="#chainlogger">chainlogger</a>  
channel containing the number of the currentBlock block field of 
the ChainManager.  
The ChainManager remains locked starting with the  call of this 
function and until it has returned.<br><br>


`func (bc *ChainManager) insert(block *types.Block)`  
<a class="function" name="ChainManager.insert">insert</a> is an
inner function, used to insert a block on the chain.  
What actually gets inserted into the chain is the block's rlp-encoding.  
This function is called by the 
<a class="link" href="#ChainManager.InsertChain">InsertChain</a> and the 
<a class="link" href="#ChainManager.Reset">Reset</a> functions.<br><br>

`func (bc *ChainManager) write(block *types.Block)`  
<a class="function" name="ChainManager.write">write</a> is an
inner function, used to write a block to the database.  
Gets called by the 
<a class="link" href="#ChainManager.InsertChain">InsertChain</a> and the 
<a class="link" href="#ChainManager.Reset">Reset</a> functions.<br><br>

`func (bc *ChainManager) Genesis() *types.Block`  
<a class="function" name="ChainManager.Genesis">Genesis</a> returns
the genesis block of the chain.<br><br>

`func (bc *ChainManager) HasBlock(hash []byte) bool`  
<a class="function" name="ChainManager.HasBlock">HasBlock</a> returns 
whether the given _hash_ param matches the hash of a block of the chain.<br><br>

`func (self *ChainManager) GetChainHashesFromHash(hash []byte, max uint64) (chain [][]byte)`  
<a class="function" name="ChainManager.GetChainHashesFromHash">GetChainHashesFromHash</a>
returns a list of hashes of the chain starting from the genesis   
hash and up to, and including, the _max_ block number. 
The _hash_ and the _max_ params  
should match the hash and the number of a specific block.<br><br>

`func (self *ChainManager) GetBlock(hash []byte) *types.Block`  
<a class="function" name="ChainManager.GetBlock">GetBlock</a>
returns the block of the chain that has the given _hash_.<br><br>

`func (self *ChainManager) GetBlockByNumber(num uint64) *types.Block`  
<a class="function" name="ChainManager.GetBlockByNumber">GetBlockByNumber</a>
returns the block of the chain that has the given _num_.<br><br>

`func (bc *ChainManager) setTotalDifficulty(td *big.Int)`  
<a class="function" name="ChainManager.setTotalDifficulty">setTotalDifficulty</a>
sets the total difficulty of the ChainManager object. 
Also stores that information on the database.<br><br>

`func (self *ChainManager) CalcTotalDiff(block *types.Block) (*big.Int, error)`  
<a class="function" name="ChainManager.CalcTotalDiff">CalcTotalDiff</a>
calculates the total difficulty of the ChainManager and returns 
it in a tuple (td, nil).  
If an error occured, then the tuple (nil, error) is returned.  
TD(genesis_block) = 0 and TD(B) = TD(B.parent) + sum(u.difficulty for u in B.uncles) + B.difficulty.<br><br>

`func (bc *ChainManager) BlockInfo(block *types.Block) types.BlockInfo`  
<a class="function" name="ChainManager.BlockInfo">BlockInfo</a>
returns the _block_'s <a class="link" href="#BlockInfo">BlockInfo</a> 
object representation.<br><br>

`func (bc *ChainManager) writeBlockInfo(block *types.Block)`  
<a class="function" name="ChainManager.writeBlockInfo">writeBlockInfo</a>
is an inner function for writing the _block_'s 
<a class="function" name="ChainManager.BlockInfo">BlockInfo</a>
representation to the database.  
This is extra information that normally should not be saved in
the db.<br><br>

`func (bc *ChainManager) Stop()`  
<a class="function" name="ChainManager.Stop">Stop</a>
sends a stop message to the chain logger channel if and only if the currentBlock  
field of the ChainManager caller object is not nil.<br><br>

`func (self *ChainManager) InsertChain(chain types.Blocks) error`    
<a class="function" name="ChainManager.InsertChain">InsertChain</a>  
iterates over the blocks in the _chain_ param and does the following:  
1. It calls the `Process` method of the `BlockProcessor` interface.  
2. writes the block to the database.  
4. sets the total difficulty of the block.  
5. inserts the block into the chain.  
6. posts a NewBlockEvent to the event mux.  
7. posts the messages to the event mux.  
Returns: either nil for success or an error.<br><br>

<br>

### 2.6 <a class="file" name="transaction_pool.go">transaction_pool.go</a>

#### 2.6.1 Data Structures

`var txplogger = logger.NewLogger("TXP")`  
<a class="struct" name="statelogger">txplogger</a> is a channel 
used to log messages regarding Blocks processing.<br><br>


`const txPoolQueueSize = 50`  
Used to initialize the _queueChan_ field of a 
<a class="link" href="#TxPool">TxPool</a>
(which is used as a queue channel to reading and writing transactions)<br><br>


`type TxPoolHook chan *types.Transaction`  
Although defined, this type is never used in this go-ethereum version.<br><br>

```
const (
	minGasPrice = 1000000
)
```
Although defined, this variable is never used in this go-ethereum version.<br><br>


`var MinGasPrice = big.NewInt(10000000000000)`  
Although defined, this variable is never used in this go-ethereum version.<br><br>


`type TxMsgTy byte`  
The only use of a 
<a class="struct" name="TxMsgTy">TxMsgTy</a> type is as a field of a 
<a class="link" href="#TxMsg">TxMsg</a> type.
Although it's not clear how this type is supposed to be used,  
since there are no other references to it in the whole codebase,
we could make a guess based on the fact that there are  
3 types of transactions in ethereum:  
1. Regular transactions: a transaction from one wallet to another.
2. Contract deployment transactions: a transaction without a _to_ address,
where the data field is used for the contract code.
3. Execution of a contract: a transaction that interacts with a
deployed smart contract.  
In this case, _to_ address is the smart contract address.  

**TxMsgTy** may had been used to represent the above types
of transactions.<br><br>



```
type TxMsg struct {
	Tx   *types.Transaction
	Type TxMsgTy
}
```
<a class="struct" name="TxMsg">TxMsg</a> represents the type of the channel of
the _subscribers_ field of a TxPool object.  
However, that field is never actually used in the whole 
codebase of this go-ethereum version.<br><br>


```
type TxProcessor interface {
	ProcessTransaction(tx *types.Transaction)
}
```
The 
<a class="interface" name="TxProcessor">TxProcessor</a> interface, 
although defined, is not implemented at all in the whole codebase   
of this go-ethereum version.<br><br>


```
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
```
TxPool is a thread safe transaction pool handler. 
In order to guarantee a non blocking pool the _queueChan_  
is used which can be independently read without needing access 
to the actual pool. If the pool is being drained  
or synced for whatever reason, the transactions
will simply queue up and be handled when the mutex is freed.  
**mutex**: a mutex for accessing the Tx pool.  
**queueChan**: Queueing channel for reading and writing incoming transactions to.  
**quit**: Quiting channel (quitting is equivalent to emptying the TxPool)  
**pool**: The actual pool, aka the list of transactions.  
**SecondaryProcessor**: **This field is actually never used as 
	the <a class="link" href="#TxProcessor">TxProcessor</a> interface is not implemented**.  
**subscribers**: **Although defined, this channel is never used**.  
**broadcaster**: used to broadcast messages to all connected peers.  
**chainManager**: the chain to which the TxPool object is attached to.  
**eventMux**: used to dispatch events to subscribers. <br><br>

#### 2.6.2 Functions

`func EachTx(pool *list.List, it func(*types.Transaction, *list.Element) bool)`  
<a class="function" name="EachTx">EachTx</a> 
is used as a means to iterate over the _pool_ list of transactions.<br><br>

`func FindTx(pool *list.List, finder func(*types.Transaction, *list.Element) bool) *types.Transaction`  
<a class="function" name="FindTx">FindTx</a> 
searches in the caller's transactions for _finder_ and returns
either the matching transaction if found, or nil.  
This is a neat way of searching for transactions that match the 
criteria defined from the _finder_ param. For example,  
<a class="link" href="#TxPool.Add">Add</a> 
uses the hash of a transaction as a searching criterion.<br><br>

`func NewTxPool(chainManager *ChainManager, broadcaster types.Broadcaster, eventMux *event.TypeMux) *TxPool`  
<a class="function" name="NewTxPool">NewTxPool</a> 
creates a new 
<a class="link" href="#TxPool">TxPool</a> object and sets it's fields.
- TxPool.pool will be set to an empty list.
- TxPool.queueChain wil be set to a Transaction channel with a txPoolQueueSize size.
- TxPool.quit will be set to a boolean channel.
- TxPool.chainManager will be assigned the param _chainManager_.
- TxPool.eventMux will be assigned the param _eventMux_.
- TxPool.broadcaster will be assigned the param _broadcaster_.  

All other fields of the 
<a class="link" href="#TxPool">TxPool</a> object that gets created are not set by NewTxPool.<br><br>


`func (pool *TxPool) addTransaction(tx *types.Transaction)`  
<a class="function" name="TxPool.addTransaction">addTransaction</a> 
is an inner function used to add the
<a class="link" href="#Transaction">Transaction</a> 
_tx_ to the end of the TxPool. This function also  
broadcasts a msg to all peers which contains 
the rlp-encodable fields of the _tx_ 
(See <a class="link" href="#Transaction.RlpData">RlpData</a>). 
The <a class="link" href="#TxPool">TxPool</a>  
remains locked starting with the call of this function 
and until it has returned.<br><br>


`func (pool *TxPool) ValidateTransaction(tx *types.Transaction) error`  
<a class="function" name="TxPool.ValidateTransaction">ValidateTransaction</a> 
validates the _tx_ 
<a class="link" href="#Transaction">Transaction</a>.
Returns either an error if _tx_ can not be validated or nil.  
These are the cases where _tx_ is not validated:
1. For some reason, the currentBlock field of the chainManager field
of the caller is nil. (aka the chain is empty)
2. The recipient field (_to_) of _tx_ is is either nil or != 20 bytes.
This means that trying to validate contract creation  
transactions, (for which the recipient _to_ is set to nil)
will always return an error.
3. The _v_ field of _tx_ is neither 28 nor 27. (See 
<a class="link" href="#Transaction">Transaction</a>)
4. The sender account of _tx_ does not have enough Ether to send to the recipient of _tx_.<br><br>


`func (self *TxPool) Add(tx *types.Transaction) error`  
<a class="function" name="TxPool.Add">Add</a> 
is the function to be called for adding a 
<a class="link" href="#Transaction">Transaction</a> to the 
<a class="link" href="#TxPool">TxPool</a> caller.
Returns either an error on not successfully  
adding _tx_ or nil for success. If _tx_ was added, 
a message is posted to the subscribed peers, 
containing the _tx_ from, to, value  
and hash fields. An error is returned in any of these cases:
1. _tx_'s hash already exists in the 
<a class="link" href="#TxPool">TxPool</a> caller, aka the transaction
to be added is already part of the caller.
2. _tx_ validation returned an error when calling 
<a class="link" href="#TxPool.ValidateTransaction">ValidateTransaction</a>.

If no errors are produced from steps 1 and 2,
<a class="link" href="#TxPool.Add">Add</a> makes a call to the inner
function 
<a class="link" href="#TxPool.addTransaction">addTransaction</a> to add _tx_ to the current  
<a class="link" href="#TxPool">TxPool</a>.<br><br>


`func (self *TxPool) Size() int`  
<a class="function" name="TxPool.Size">Size</a> 
returns the number of Transactions of the caller.<br><br>


`func (pool *TxPool) CurrentTransactions() []*types.Transaction`  
<a class="function" name="TxPool.CurrentTransactions">CurrentTransactions</a> 
returns the transactions of the 
<a class="link" href="#TxPool">TxPool</a> caller as a slice.<br><br>


`func (pool *TxPool) RemoveInvalid(state *state.StateDB)`  
<a class="function" name="TxPool.RemoveInvalid">RemoveInvalid</a> 
removed all transactions from the caller for which either:
1. the transaction returns an error when validated through the 
<a class="link" href="#TxPool.ValidateTransaction">ValidateTransaction</a> function, or
2. the transaction sender's nonce field is >= to the transaction's nonce field.<br><br>


`func (self *TxPool) RemoveSet(txs types.Transactions)`  
<a class="function" name="TxPool.RemoveSet">RemoveSet</a> 
takes as an argument a set of transactions _txs_ and
removes from the caller's transactions set those that  
match the ones from _txs_. Looping over the transactions 
of the caller happens through the 
<a class="link" href="#EachTx">EachTx</a> function.<br><br>


`func (pool *TxPool) Flush() []*types.Transaction`  
<a class="function" name="TxPool.Flush">Flush</a> 
resets the caller's transactions list to an empty list.<br><br>


`func (pool *TxPool) Start()`  
Although defined, this function does not contain any executable code in this go-ethereum version.<br><br>


`func (pool *TxPool) Stop()`  
<a class="function" name="TxPool.Stop">Stop</a> 
makes a call on 
<a class="link" href="#TxPool.Flush">Flush</a> 
to empty the caller's transactions list  
and then sends the message  "Stopped" to the 
<a class="link" href="#txplogger">txplogger</a> channel.<br><br>

<br>

### 2.7 <a class="file" name="asm.go">asm.go </a>

The file contains only one function:

`func Disassemble(script []byte) (asm []string)`  
<a class="function" name="Disassemble">Disassemble</a> returns a 
string representation of a sequence of bytes that consist an evm bytecode.  
The opcodes are defined in vm/types.go. 
In case that we have a PUSHi opcode we expect the next i  
bytes to be the i items that we want to push to the stack.  
_script_: The evm bytecode.  

An example can be found here:  
https://rinkeby.etherscan.io/address/0x147b8eb97fd247d06c4006d269c90c1908fb5d54#code

**Example**: Passing the first series of bytes of the above link to this function as  
```
script := []byte(
					0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 
					0x80, 0x15, 0x61, 0x00, 0x10, 0x57, 
					0x60, 0x00, 0x80, 0xfd, 0x5b, 0x50
				)
```
will yield the following output which is presented in a more human-readable way.  
The first column is only there to show match the bytes of the input with the results.

0x60 --> 0000: PUSH1  
0x80 --> 0001: 0x80  
0x60 --> 0002: PUSH1 (next value is pushed onto the stack)  
0x40 --> 0003: 0x40  
0x52 --> 0004: MSTORE  
0x34 --> 0005: CALLVALUE  
0x80 --> 0006: DUP1  
0x15 --> 0007: ISZERO  
0x61 --> 0008: PUSH2 (next 2 values are pushed onto the stack)   
0x00 --> 0009: 0x00  
0x10 --> 0010: 0x10  
0x57 --> 0011: JUMPI  
0x60 --> 0012: PUSH1  
0x00 --> 0013: 0x00  
0x80 --> 0014: DUP1  
0xfd --> 0015: Missing opcode 0xfd  
0x5b --> 0016: JUMPDEST  
0x50 --> 0017: POP<br><br>

<br>

### 2.8 <a class="file" name="error.go">error.go </a>

#### 2.8.1 Data structures


```
type UncleErr struct {
	Message string
}
```
<a class="struct" name="UncleErr">UncleErr</a> error. 
This error is thrown only from the function 
<a class="link" href="#BlockManager.AccumelateRewards">AccumelateRewards</a>.<br><br>


```
type ParentErr struct {
	Message string
}
```
<a class="struct" name="ParentErr">ParentErr</a> error. 
In case a parent is unknown this error will be thrown by the block manager.<br><br>


```
type ValidationErr struct {
	Message string
}
```
<a class="struct" name="ValidationErr">ValidationErr</a>
is a block validation error. 
If any validation fails, this error will be thrown.<br><br>


```
type GasLimitErr struct {
	Message string
	Is, Max *big.Int
}
```
A <a class="struct" name="GasLimitErr">GasLimitErr</a>
happens when the total gas left for the coinbase address 
is less than the gas to be bought.<br><br>


```
type NonceErr struct {
	Message string
	Is, Exp uint64
}
```
A <a class="struct" name="NonceErr">NonceErr</a>
happens when a transaction's nonce is incorrect. <br><br>


```
type OutOfGasErr struct {
	Message string
}
```
A <a class="struct" name="OutOfGasErr">OutOfGasErr</a>
happens when the gas provided runs out before the state transition happens.<br><br>


```
type TDError struct {
	a, b *big.Int
}
```
<a class="struct" name="TDError">TDError</a> is defined, but 
never used in this go-ethereum version. 
Meant to be used as a total difficulty error.<br><br>


```
type KnownBlockError struct {
	number *big.Int
	hash   []byte
}
```
A <a class="struct" name="KnownBlockError">KnownBlockError</a>
happens when there is already a block in the chain with the same hash.  
The number field is the number of the existing block.

<br>

#### 2.8.2 Functions

`func (self *TDError) Error() string`  
<a class="function" name="TDError.Error">TDError.Error</a>
creates and returns a TDError error.<br><br>

`func (self *KnownBlockError) Error() string`  
<a class="function" name="KnownBlockError.Error">KnownBlockError.Error</a>
creates and returns a KnownBlockError error.<br><br>

`func ParentError(hash []byte) error`  
<a class="function" name="ParentError">ParentError</a>
creates a ParentError object by setting it's message to a message  
that includes the 'hash' and returns it.<br><br>

`func UncleError(str string) error`  
<a class="function" name="UncleError">UncleError</a>
creates an UncleErr error by setting it's message to 'str' and returns it.<br><br>

`func ValidationError(format string, v ...interface{}) *ValidationErr`  
<a class="function" name="ValidationError">ValidationError</a>
creates a ValidationErr error by setting it's message and returns it.<br><br>

`func GasLimitError(is, max *big.Int) *GasLimitErr`  
<a class="function" name="GasLimitError">GasLimitError</a>
creates and returns a GasLimitError given the total gas left  
to be bought by a coinbase address and the actual gas.<br><br>

`func NonceError(is, exp uint64) *NonceErr`  
<a class="function" name="NonceError">NonceError</a>
creates and returns a NonceError given the transaction's nonce  
and the nonce of the sender of the transaction.<br><br> 

`func OutOfGasError() *OutOfGasErr`  
<a class="function" name="OutOfGasError">OutOfGasError</a>
creates and returns an OutOfGasError error.<br><br>

`func IsParentErr(err error) bool`  
<a class="function" name="IsParentErr">IsParentErr</a>
returns whether 'err' is a ParentErr error.<br><br>

`func IsUncleErr(err error) bool`  
<a class="function" name="IsUncleErr">IsUncleErr</a>
returns whether 'err' is an UncleErr error.<br><br>

`func IsValidationErr(err error) bool`  
<a class="function" name="IsValidationErr">IsValidationErr</a>
returns whether 'err' is a ValidationErr error.<br><br>

`func IsGasLimitErr(err error) bool`  
<a class="function" name="IsGasLimitErr">IsGasLimitErr</a>
returns whether 'err' is a GasLimitErr error.<br><br>

`func IsNonceErr(err error) bool`  
<a class="function" name="IsNonceErr">IsNonceErr</a>
returns whether 'err' is a NonceErr error.<br><br>

`func IsOutOfGasErr(err error) bool`  
<a class="function" name="IsOutOfGasErr">IsOutOfGasErr</a>
returns whether 'err' is an OutOfGasErr error.<br><br>

`func IsTDError(e error) bool`  
<a class="function" name="IsTDError">IsTDError</a>
returns whether 'e' is a TDError error.<br><br>

`func IsKnownBlockErr(e error) bool`  
<a class="function" name="IsKnownBlockErr">IsKnownBlockErr</a>
returns whether 'e' is a KnownBlockErr error.<br><br>

`func (err *ParentErr) Error() string`  
<a class="function" name="ParentErr.Error">ParentErr.Error</a>
returns the error message of the caller.<br><br>

`func (err *UncleErr) Error() string`  
<a class="function" name="UncleErr.Error">UncleErr.Error</a>
returns the error message of an UncleErr error.<br><br>

`func (err *ValidationErr) Error() string`  
<a class="function" name="ValidationErr.Error">ValidationErr.Error</a>
returns the error message of a ValidationErr error.<br><br>

`func (err *GasLimitErr) Error() string`  
<a class="function" name="GasLimitErr.Error">GasLimitErr.Error</a>
returns the error message of a GasLimitErr error.<br><br>

`func (err *NonceErr) Error() string`  
returns the error message of a NonceErr error.<br><br>

`func (self *OutOfGasErr) Error() string`  
<a class="function" name="OutOfGasErr.Error">OutOfGasErr.Error</a>
returns the error message of an OutOfGasError error.<br><br>

<br>

### 2.9 <a class="file" name="event.go">event.go </a>

#### 2.9.1 Data Structures

`type TxPreEvent struct{ Tx *types.Transaction }`  
A <a class="struct" name="TxPreEvent">TxPreEvent</a>
is posted when a transaction enters the transaction pool.<br><br>

`type TxPostEvent struct{ Tx *types.Transaction }`  
A <a class="struct" name="TxPostEvent">TxPostEvent</a>
is posted when a transaction has been processed.<br><br>

`type NewBlockEvent struct{ Block *types.Block }`  
A <a class="struct" name="NewBlockEvent">NewBlockEvent</a>
is posted when a block has been imported.<br><br>

<br>


### 2.10 <a class="file" name="dagger.go">dagger.go </a>

Dagger is the consensus algorithm.

Dagger was proven to be vulnerable to shared memory hardware acceleration:  
https://bitslog.com/2014/01/17/ethereum-dagger-pow-is-flawed/

The consensus algorithm for Ethereum 1.0 is Ethash (PoW algorithm):  
https://eth.wiki/en/concepts/ethash/ethash

<br> 

Essentially, the <a class="function" name="Dagger">Dagger</a>
algorithm works by creating a directed acyclic graph with ten levels  
including the root and a total of 225 - 1 values. In levels 1 through 8, the value of each node  
depends on three nodes in the level above it, and the number of nodes in each level is eight  
times larger than in the previous. In level 9, the value of each node depends on 16 of its  
parents, and the level is only twice as large as the previous; the purpose of this is to make the  
natural time-memory tradeoff attack be artificially costly to implement at the first level, so that  
it would not be a viable strategy to implement any time-memory tradeoff optimizations at all.  
Finally, the algorithm uses the underlying data, combined with a nonce, to pseudorandomly select  
eight bottom-level nodes in the graph, and computes the hash of all of these nodes put together.  
If the miner finds a nonce such that this resulting hash is below 2^256 divided by the difficulty  
parameter, the result is a valid proof of work.

Here is the breakdown of the algorithm that the below functions implement:  
- **D**: block header
- **N**: nonce 
- **||**: string concatenation operator.

```
if L==9
	dependsOn(L)=16
else
	dependsOn(L)=3

Node(D,xn,0,0)=D
Node(D,xn,L,i) =
	for k in [0...dependsOn(L)-1]
		p[k] = sha256(D || xn || L || i || k) mod 8^(L-1)
		sha256(node(L-1,p[0]) || node(L-1,p[1]) ... || node(L-1,p[dependsOn(L)-1]))

eval(D,N) =
	for k in [0...3]
		p[k] = 	sha256( D  || floor(n / 2^26) || i || k ) mod 8^8 * 2 
	
		sha256(	node( D,floor(n / 2^26),9,p[0])  ||  node(D,floor(n / 2^26),9,p[1]) ... 
			... || node(D,floor(n / 2^26),9,p[3])))

Objective: find k such that eval(D,k) < 2^256 / diff
```

<br>

#### 2.10.1 Functions

<br>

`func (dag *Dagger) Find(obj *big.Int, resChan chan int64)`  

`func (dag *Dagger) Search(hash, diff *big.Int) *big.Int`  

`func (dag *Dagger) Verify(hash, diff, nonce *big.Int) bool`  

`func DaggerVerify(hash, diff, nonce *big.Int) bool`  

`func (dag *Dagger) Node(L uint64, i uint64) *big.Int`  

`func Sum(sha hash.Hash) []byte`  

`func (dag *Dagger) Eval(N *big.Int) *big.Int`  <br><br>

<br>


### 2.11 <a class="file" name="state_transition.go">state_transition.go </a>


#### 2.11.1 Data Structures



```
type StateTransition struct {
	coinbase, receiver []byte
	msg                Message
	gas, gasPrice      *big.Int
	initialGas         *big.Int
	value              *big.Int
	data               []byte
	state              *state.StateDB
	block              *types.Block

	cb, rec, sen *state.StateObject

	Env vm.Environment
}
```
A <a class="struct" name="StateTransition">StateTransition</a>. 
A state transition is a change made when a transaction is applied  
to the current world state. The state transitioning model does all  
the necessary work to work out a valid new state root:
1. Nonce handling  
2. Pre pay / buy gas of the coinbase (miner)
3. Create a new state object if the recipient is 0, aka contract creation.
4. Value transfer  
== If contract creation ==  
4a. Attempt to run transaction data  
4b. If valid, use result as  code for the new state object  
== end ==
5) Run Script section
6) Derive the new state root<br><br>

- **coinbase**: The miner state object.
- **receiver**: The receiver state object (in case of Ether transfer)
- **msg**: Representation of the transaction to be applied.
- **gas**: The gas limit equal to the maximum amount of gas that should be  
	used when executing the msg/transaction. 
- **gasPrice**: This is equal to the number of Wei to be paid per unit of gas for  
	all computation costs incurred as a result of the execution of the msg/transaction.
- **initialGas**: The gas that has been supplied for the execution of the msg/transaction.
- **value**: The amount of Wei to be transferred.
- **data**: An unlimited size byte array specifying the input data of the
	msg/transaction.   
	The first 4 bytes of this field specify which function to call
	when the msg/transaction  
	will execute by using the hash of the function's name
	to be called and it's arguments.  
	The rest of the data field are the arguments
	passed to the function.  
	If the data field is empty, it means a msg/transaction is for a payment and not 
	an execution of the contract.
- **state**: The current world state before the execution of the msg/transaction.
- **block**: The block (to be added to the chain) that contains the _msg_ to be executed.
- **cb**: The miner's address.
- **rec**: The receiver's address.
- **sen**: The sender's address.<br><br><br>

```
type Message interface {
	Hash() 			[]byte
	From() 			[]byte
	To() 			[]byte
	GasPrice() 		*big.Int
	Gas() 			*big.Int
	Value() 		*big.Int
	Nonce() 		uint64
	Data() 			[]byte
}
```
A <a class="interface" name="Message">Message</a>
represents a transaction. All functions of this interface are  
getter functions for the relevant fields of the msg/transaction.<br><br>

<br>

#### 2.11.2 Functions


`func AddressFromMessage(msg Message) []byte`  
<a class="function" name="AddressFromMessage">AddressFromMessage</a>
creates and returns a new address based on the _msg sender_ and _nonce_ fields.<br><br>


`func MessageCreatesContract(msg Message) bool`  
<a class="function" name="MessageCreatesContract">MessageCreatesContract</a>
returns whether the _msg_ is a contract creation aka whether the _msg
recipient_ is 0.<br><br>

`func MessageGasValue(msg Message) *big.Int`  
<a class="function" name="MessageGasValue">MessageGasValue</a>
returns the amount of Wei based on the _msg gas_ and _gasPrice_ fields.  
gasValue = gas * gasPrice<br><br>

`func NewStateTransition(coinbase *state.StateObject, msg Message, state *state.StateDB, block *types.Block) *StateTransition`  
<a class="function" name="NewStateTransition*types.Block) *StateTransition`  ">NewStateTransition</a>
creates and returns a <a class="link" href="#StateTransition">StateTransition</a> object.  
The fields _gas_ and initialGas are set to 0.  
The fields rec, sen and Env are set to nil.  
The field coinbase is set to the address of the _coinbase_ param.  
The field cb is set to the param _coinbase_.  
All other fields are set to the corresponding params.<br><br>

`func (self *StateTransition) VmEnv() vm.Environment`  
<a class="function" name="StateTransition.VmEnv">VmEnv</a>
is a getter method for the _Env_ field of a 
<a class="link" href="#StateTransition">StateTransition</a> object.   
If the Env field of the caller is nil, a new Env is created by
calling the function <a class="link" href="#	">NewEnv</a>.<br><br>

`func (self *StateTransition) Coinbase() *state.StateObject`  
<a class="function" name="StateTransition.Coinbase">Coinbase</a>
returns the miner of the msg's block. If the miner does
not exist in the current world state, it's created.<br><br>

`func (self *StateTransition) From() *state.StateObject`  
<a class="function" name="StateTransition.From">From</a>
returns the _from_ field of the msg. If
_from_ does not exist in the current world state, it's created.<br><br>

`func (self *StateTransition) To() *state.StateObject`  
<a class="function" name="StateTransition.To">To</a>
returns the _to_ field of the msg. If
_to_ does not exist in the current world state, it's created.  
This function will return nil in the case where the msg is about
a contract creation (aka if _to_ is 0)<br><br>

`func (self *StateTransition) UseGas(amount *big.Int) error`  
<a class="function" name="StateTransition.UseGas">UseGas</a>
attempts to use _amount_ gas of the caller's gas. If the caller's
gas is less than the _amount_ provided,   
an <a class="link" href="#OutOfGasError">OutOfGasError</a> is returned.
Otherwise, nil is returned for success. In case of success, the new  
gas of the caller will become:  
newGas = prevGas - amount.<br><br>

`func (self *StateTransition) AddGas(amount *big.Int)`  
<a class="function" name="StateTransition.AddGas">AddGas</a>
adds _amount_ gas to the caller's gas. Helps in keeping track of the 
altogether used gas of a list of transactions.<br><br>

`func (self *StateTransition) BuyGas() error`  
<a class="function" name="StateTransition.BuyGas">BuyGas</a>
attempts to reward the miner with the gas of the transaction.
If the sender's balance is less than the calculated gas in Wei
(gas*gasPrice of caller), an error is returned.
Buying the gas does not directly happen in this function. Instead,
the _BuyGas_ function of the miner (StateObject) is called through
this function. If the latter does not return an error, this function
will increase the _gas_ field of the caller, set the caller's
_initialGas_ field and decrease the sender's balance by an amount
of _gas_ *_gasPrice_<br><br>

`func (self *StateTransition) preCheck() (err error)`  
<a class="function" name="StateTransition.preCheck">preCheck</a>
is an inner function, used by the 
<a class="link" href="#StateTransition.TransitionState">TransitionState</a> 
function and does 2 things:
1. Checks whether the caller's msg sender nonce is the same as the caller's msg nonce.
	If not, it returns an error.
2. Calls 
<a class="link" href="#StateTransition.BuyGas">BuyGas</a> in order to reward the miner.
	If <a class="link" href="#StateTransition.BuyGas">BuyGas</a> returns an error,
	this function returns that error.  
	If everything went well, this function returns nil.<br><br>

`func (self *StateTransition) TransitionState() (ret []byte, err error)`  
<a class="function" name="StateTransition.TransitionState">TransitionState</a>
attempts to alter the world state by applying the msg/transaction of the caller.
1. Calls <a class="link" href="#StateTransition.preCheck">preCheck</a> 
	for nonce validation and to reward the miner.
2. Schedules a gas refund (to return the unused gas)
3. increases nonce of the msg sender.
4. uses the TxGas. (defined as Gtransaction in the Ethereum Yellow Paper)
   TxGas is a constant value set to 500 Wei. (see the GasTx variable defined in vm/common.go)
   On the Ethereum Yellow Paper this value is set to 21000 Wei.
5. uses the GasData. This is the gas that must be payed for every byte
   of the log field of the msg. On the Ethereum Yellow Paper this value is set to 8 Wei per byte.
6. If the msg is about a contract creation (msg recipient is 0) this function makes a call
   to the <a class="link" href="#MakeContract">MakeContract</a>, uses the appropriate gas 	needed and sets the code of the to-be-deployed contract.
7. If the msg is not a contract creation msg, a call to the vm's 
	<a class="link" href="#VMEnv.Call">Call</a> function is made. The vm takes care of executing the transaction and may return an error.
8. In case where none of the above operations returned an error, a call on 
	<a class="link" href="StateTransition.UseGas">UseGas</a> is made.
9. As usual, in case of any errors the returned tuple will be (nil, error) or in case 
	of no errors, the error returned will be nil and the _ret_ field of the tuple will
	contain _the returned result of the vm_ of steps 6 or 7.<br><br>

`func MakeContract(msg Message, state *state.StateDB) *state.StateObject`  
<a class="function" name="MakeContract">MakeContract</a>
converts an transaction in to a state object. The stateObject that this function 
creates and returns may then be used (as a param) to actually deploy the contract through the
vm <a class="link" href="#VMEnv.Create">Create</a> function.<br><br>

`func (self *StateTransition) RefundGas()`  
<a class="function" name="StateTransition.RefundGas">RefundGas</a>
takes care of refunding gas any remaining gas in case of a successful msg/transaction execution.<br><br>

`func (self *StateTransition) GasUsed() *big.Int`  
<a class="function" name="StateTransition.GasUsed">GasUsed</a>
returns how much gas has been used by the execution of the msg/transaction<br><br>

<br>


### 2.12 <a class="file" name="vm_env.go">vm_env.go </a>


#### 2.12.1 Data Structures

```
type VMEnv struct {
	state 	*state.StateDB
	block 	*types.Block
	msg   	Message
	depth 	int
}
```

The <a class="struct" name="VMEnv">VMEnv</a> represents the Ethereum Virtual Machine.
- **state**: The current world state.
- **block**: The block to which the msg/transaction to be executed belongs to.
- **msg**: The msg/transaction to be executed.
- **depth**: The EVM operates as a stack machine. This variable is the maximum number
	of items the stack can hold.<br><br>


#### 2.12.2 Functions

`func NewEnv(state *state.StateDB, msg Message, block *types.Block) *VMEnv`  
creates and returns a new <a class="link" href="#VMEnv">VMEnv</a> object.<br><br>

The <a class="link" href="#VMEnv">VMEnv</a> implements the _Environment_ interface. 
For that reason, the following functions are implemented:<br><br>

`func (self *VMEnv) Origin() []byte        { return self.msg.From() }`  
<a class="function" name="VMEnv.Origin">Origin</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the origin of the msg to be executed.<br>

`func (self *VMEnv) BlockNumber() *big.Int { return self.block.Number }`  
<a class="function" name="VMEnv.BlockNumber">BlockNumber</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the block number of the block to which the msg to executed belongs to.<br>

`func (self *VMEnv) PrevHash() []byte      { return self.block.PrevHash }`  
<a class="function" name="VMEnv.PrevHash">PrevHash</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the hash of the parent of the block that contains the msg to be executed.<br>

`func (self *VMEnv) Coinbase() []byte      { return self.block.Coinbase }`  
<a class="function" name="VMEnv.Coinbase">Coinbase</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the miner's address of the block that contains the msg to be executed.<br>

`func (self *VMEnv) Time() int64           { return self.block.Time }`  
<a class="function" name="VMEnv.Time">Time</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the time that the block that contains the msg to be executed will be/is created.<br>

`func (self *VMEnv) Difficulty() *big.Int  { return self.block.Difficulty }`  
<a class="function" name="VMEnv.Difficulty">Difficulty</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the difficulty of the block that contains the msg to be executed.<br>

`func (self *VMEnv) BlockHash() []byte     { return self.block.Hash() }`  
<a class="function" name="VMEnv.BlockHash">BlockHash</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the hash of the block that contains the msg to be executed.<br>

`func (self *VMEnv) Value() *big.Int       { return self.msg.Value() }`  
<a class="function" name="VMEnv.Value">Value</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the value in Wei of the msg to be executed.<br>

`func (self *VMEnv) State() *state.StateDB { return self.state }`  
<a class="function" name="VMEnv.State">State</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the current world state.<br>

`func (self *VMEnv) GasLimit() *big.Int    { return self.block.GasLimit }`  
<a class="function" name="VMEnv.GasLimit">GasLimit</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the gasLimit field of the block that contains the msg to be executed.<br>

`func (self *VMEnv) Depth() int            { return self.depth }`  
<a class="function" name="VMEnv.Depth">Depth</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which returns the maximum number of items of the EVM stack.<br>

`func (self *VMEnv) SetDepth(i int)        { self.depth = i }`  
<a class="function" name="VMEnv.SetDepth">SetDepth</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which sets the maximum number of items of the EVM stack.<br>

`func (self *VMEnv) AddLog(log state.Log)  { self.state.AddLog(log) }`  
<a class="function" name="VMEnv.AddLog">AddLog</a> 
is defined as part of the implementation of the _Environment_ interface for a VMEnv object,  
which adds a log to the world state through the function AddLog of a StateDB object
(defined in the file state/state.go).<br><br>


`func (self *VMEnv) Transfer(from, to vm.Account, amount *big.Int) error`  
<a class="function" name="VMEnv.Transfer">Transfer</a> makes use of the generic function 
Transfer defined in the file vm/environment.go to  
make the transfer of _amount_ Wei from the _from_ to the _to_ accounts.
Returns nil on success or an error if there is not sufficient balance in the _from_
account.<br><br>


`func (self *VMEnv) vm(addr, data []byte, gas, price, value *big.Int) *Execution`  
<a class="function" name="VMEnv.vm">vm</a>
is an inner function that creates and returns a new 
<a class="link" href="#Execution">Execution</a> object.<br><br>

`func (self *VMEnv) Call(me vm.ClosureRef, addr, data []byte, gas, price, value *big.Int) ([]byte, error)`  
<a class="function" name="VMEnv.Call">Call</a>
is defined as part of the implementation of the Environment interface for a VMEnv object,
which executes the contract associated at _addr_ with the given _data_. 
This happens through a call on the inner function 
<a class="function" name="Execution.exec">exec</a>.<br><br>

`func (self *VMEnv) CallCode(me vm.ClosureRef, addr, data []byte, gas, price, value *big.Int) ([]byte, error)`  
<a class="function" name="VMEnv.CallCode">CallCode</a>
is defined as part of the implementation of the Environment interface for a VMEnv object,
which executes the contract associated at _addr_ with the given _data_.
This happens through a call on the inner function 
<a class="function" name="Execution.exec">exec</a>. The only 
diffference between <a class="link" href="VMEnv.Call">Call</a> and CallCode
is that the latter executes the given address' 
code with the caller as context (aka CallCode is used when a function calls 
another function with the 1st function being the caller).<br><br>

`func (self *VMEnv) Create(me vm.ClosureRef, addr, data []byte, gas, price, value *big.Int) ([]byte, error, vm.ClosureRef)`  
<a class="function" name="VMEnv.Create">Create</a>
is defined as part of the implementation of the Environment interface for a VMEnv object,
which creates a contract through the function 
<a class="function" name="Execution.Create">Create</a.<br><br>

<br>

### 2.13 <a class="file" name="execution.go">execution.go </a>

#### 2.13.1 Data Structures

```
type Execution struct {
	env               vm.Environment
	address, input    []byte
	Gas, price, value *big.Int
	SkipTransfer      bool
}
```
The EVM always creates an 
<a class="struct" name="Execution">Execution</a>  object
and then calls upon it's functions in order to execute any transaction.
In other words, this is the core object through which any code on the blockchain is executed.

- **env**: The Execution object needs to be "attached" to an EVM.
- **address**: The address of the Execution. This should be a contract's address as it 
	refers to an address which contains the code to be executed.
- **input**: The code to be executed. This param should be of the form
	{hash, params}, where hash is the hash of the function's signature to be executed
	and params are 32-byte words which are passed as params to the function. This is where
	the stack comes into play. 
- **Gas**: The gas that has been supplied for the the msg/transaction.
- **price**: The gas price.
- **value**: Amount of Wei to be transferred throught the execution of the transaction.
- **SkipTransfer**: Used only for testing.<br><br> 

<br>

#### 2.13.2 Functions


`func NewExecution(env vm.Environment, address, input []byte, gas, gasPrice, value *big.Int) *Execution`  
<a class="struct" name="NewExecution">NewExecution</a> is a contructor for an
<a class="link" href="#Execution">Execution</a> object. This function gets 
mainly called from the EVM to create an Execution object through which (the EVM)
will call either <a class="link" href="#Execution.Call">Call</a> or 
<a class="link" href="#Execution.Create">Create</a> in order to execute code or
create a contract accordingly.<br><br>


`func (self *Execution) Addr() []byte`  
<a class="function" name="Execution.Addr">Addr</a>
returns the address of the <a class="struct" name="Execution">Execution</a>.<br><br>


`func (self *Execution) Call(codeAddr []byte, caller vm.ClosureRef) ([]byte, error)`  
<a class="function" name="Execution.Call">Call</a>
is the function to be called when a msg/transaction is not about contract creation.
It does nothing more than retrieving the code of the _codeAddr_ (to be executed)
and makes a call to the inner function _exec_ which takes care of the execution.
This function gets called from both 
<a class="link" href="#VMEnv.CallCode">CallCode</a> and
<a class="link" href="#VMEnv.Call">Call</a> functions of the EVM.<br><br>


`func (self *Execution) exec(code, contextAddr []byte, caller vm.ClosureRef) (ret []byte, err error)`  
<a class="function" name="Execution.exec">exec</a>
is the most important function of the whole core package and one of the most important
ones of the whole go-ethereum codebase. First of all, it's an inner function, called by 
<a class="link" href="#Execution.Call">Call</a> and/or 
<a class="link" href="#Execution.Create">Create</a>. 
exec executes the _code_ from the contract _codeAddr_. 
It handles any necessary value transfer required and takes
the necessary steps to create accounts and reverses the state in case of an
execution error or failed value transfer. <br><br>

- **caller**: In cases where eg. a function A
	calls function B, the _caller_ param will point to the the contract's address 
	that made the call to B.
- **contextAddr**: The address of the contract that includes the **code** to be
	executed. <br><br>


`func (self *Execution) Create(caller vm.ClosureRef) (ret []byte, err error, account *state.StateObject)`  
<a class="function" name="Execution.Create">Create</a>
creates a new contract and returns it's bytecode and the account created as a StateObject.<br><br>

<br>


</body>
