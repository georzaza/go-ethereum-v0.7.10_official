package ezp

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/georzaza/go-ethereum-v0.7.10_official/crypto"
	"github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"
	"github.com/georzaza/go-ethereum-v0.7.10_official/logger"
	"github.com/georzaza/go-ethereum-v0.7.10_official/pow"
	"github.com/obscuren/sha3"
)

var powlogger = logger.NewLogger("POW")

type EasyPow struct {
	hash     *big.Int
	HashRate int64
	turbo    bool
}

func New() *EasyPow {
	return &EasyPow{turbo: true}
}

func (pow *EasyPow) GetHashrate() int64 {
	return pow.HashRate
}

func (pow *EasyPow) Turbo(on bool) {
	pow.turbo = on
}

func (pow *EasyPow) Search(block pow.Block, stop <-chan struct{}) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hash := block.HashNoNonce()
	diff := block.Diff()
	i := int64(0)
	start := time.Now().UnixNano()
	t := time.Now()

	for {
		select {
		case <-stop:
			powlogger.Infoln("Breaking from mining")
			pow.HashRate = 0
			return nil
		default:
			i++

			if time.Since(t) > (1 * time.Second) {
				elapsed := time.Now().UnixNano() - start
				hashes := ((float64(1e9) / float64(elapsed)) * float64(i)) / 1000
				pow.HashRate = int64(hashes)
				powlogger.Infoln("Hashing @", pow.HashRate, "khash")

				t = time.Now()
			}

			sha := crypto.Sha3(big.NewInt(r.Int63()).Bytes())
			if pow.verify(hash, diff, sha) {
				return sha
			}
		}

		if !pow.turbo {
			time.Sleep(20 * time.Microsecond)
		}
	}

	return nil
}

func (pow *EasyPow) verify(hash []byte, diff *big.Int, nonce []byte) bool {
	sha := sha3.NewKeccak256()

	d := append(hash, nonce...)
	sha.Write(d)

	verification := new(big.Int).Div(ethutil.BigPow(2, 256), diff)
	res := ethutil.U256(ethutil.BigD(sha.Sum(nil)))

	return res.Cmp(verification) <= 0
}

func (pow *EasyPow) Verify(block pow.Block) bool {
	return pow.verify(block.HashNoNonce(), block.Diff(), block.N())
}
