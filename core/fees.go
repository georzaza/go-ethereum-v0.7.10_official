package core

import (
	"math/big"
)

// Initial block reward for miners. See the function BlockManager.AccumelateRewards in the 'core' package (file block_manager.go)
var BlockReward *big.Int = big.NewInt(1.5e+18)
