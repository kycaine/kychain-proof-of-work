package blockchain

import (
	"strings"
)

func (b *Block) Mine(difficulty int) {
	prefix := strings.Repeat("0", difficulty)
	for {
		b.Nonce++
		b.Hash = b.CalculateHash()
		if b.Hash[:difficulty] == prefix {
			break
		}
	}
}
