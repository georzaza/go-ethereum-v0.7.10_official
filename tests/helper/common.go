package helper

import "github.com/georzaza/go-ethereum-v0.7.10_official/ethutil"

func FromHex(h string) []byte {
	if ethutil.IsHex(h) {
		h = h[2:]
	}

	return ethutil.Hex2Bytes(h)
}
