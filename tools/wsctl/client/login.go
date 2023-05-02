package client

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spruceid/siwe-go"
	"github.com/tidwall/gjson"
)

type (
	loginByAddress struct {
		Message   string `json:"message"`   // Message siwe serialized message
		Signature string `json:"signature"` // Signature should have '0x' prefix
	}
)

var pvk1 = "3b3a4ccb94b92b43af8e3987d181d340a754e6b4168811f3a80bdc7e6edbcda4"

func (c *client) login() error {
	// pvk1 := c.Config().PrivateKey

	pvk, err := loadPrivateKey(pvk1)
	if err != nil {
		return err
	}

	msg, err := prepareMessage(pvk.PublicKey)
	if err != nil {
		return err
	}

	if _, err := siwe.ParseMessage(msg); err != nil {
		return err
	}

	sig, err := signMessage(msg, pvk)
	if err != nil {
		return err
	}

	body, err := json.Marshal(loginByAddress{msg, sig})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/login/wallet", c.cfg.Endpoint)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		c.logger.Panic(errors.Wrap(err, "failed to create login request"))
	}
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		c.logger.Panic(errors.Wrapf(err, "failed to login %s", url))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))
	if !gjson.ValidBytes(respBody) {
		panic("invalid response")
	}

	ret := gjson.ParseBytes(respBody)
	tt := ret.Get("token").String()
	fmt.Println("token acquired", tt)

	c.token.Store(tt)
	return nil
}

func loadPrivateKey(pvkStr string) (*ecdsa.PrivateKey, error) {
	var pvk *ecdsa.PrivateKey
	if len(pvkStr) > 0 {
		fmt.Println("loaded private key from the config.")
		pvk = crypto.ToECDSAUnsafe(common.FromHex(pvkStr))
	} else {
		fmt.Println("no private key is found in the config; a new one is randomly generated.")
		var err error
		if pvk, err = crypto.GenerateKey(); err != nil {
			return nil, err
		}
	}
	return pvk, nil
}

func prepareMessage(pbk ecdsa.PublicKey) (string, error) {
	msg, err := siwe.InitMessage("w3bstream.com",
		crypto.PubkeyToAddress(pbk).String(),
		"https://w3bstream.com",
		uniuri.NewLen(16),
		nil,
	)
	if err != nil {
		return "", err
	}
	return msg.String(), nil
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	sign := signHash([]byte(message))
	signature, err := crypto.Sign(sign.Bytes(), privateKey)

	if err != nil {
		return "", err
	}

	signature[64] += 27
	return "0x" + hex.EncodeToString(signature), nil
}

func signHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}
