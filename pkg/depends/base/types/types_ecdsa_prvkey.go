package types

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewEcdsaPrvKey() (EcdsaPrvKey, error) {
	sk, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	return EcdsaPrvKey(hex.EncodeToString(crypto.FromECDSA(sk))), nil
}

type EcdsaPrvKey string

func (p EcdsaPrvKey) String() string { return string(p) }

func (p EcdsaPrvKey) SecurityString() string { return "--------" }

func (p EcdsaPrvKey) Address() (string, error) {
	sk, err := crypto.HexToECDSA(p.String())
	if err != nil {
		return "", err
	}
	pk := crypto.PubkeyToAddress(sk.PublicKey)
	return pk.Hex(), nil
}

// Marshal to hex string
func (p EcdsaPrvKey) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// Unmarshal from hex string with strict mode
func (p *EcdsaPrvKey) UnmarshalText(d []byte) error {
	v := string(d)
	_, err := crypto.ToECDSA(common.FromHex(v))
	if err != nil {
		return err
	}
	*p = EcdsaPrvKey(v)
	return nil
}
