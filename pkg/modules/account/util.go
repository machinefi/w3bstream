package account

import (
	"crypto/rand"
	"math/big"
	"sync"
)

var (
	max  *big.Int
	once sync.Once
)

func GenerateNonce() (string, error) {
	once.Do(func() {
		max = new(big.Int)
		max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	})
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return n.Text(10), nil
}
