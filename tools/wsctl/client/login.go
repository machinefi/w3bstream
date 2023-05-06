package client

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func (c *client) login() (string, error) {
	pvk, err := c.loadPrivateKey(c.Config().PrivateKey)
	if err != nil {
		return "", err
	}
	msg, err := prepareMessage(c.Address())
	if err != nil {
		return "", err
	}

	if _, err := siwe.ParseMessage(msg); err != nil {
		return "", err
	}

	sig, err := signMessage(msg, pvk)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(loginByAddress{msg, sig})
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/login/wallet", c.Config().Endpoint)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Panic(errors.Wrap(err, "failed to create login request"))
	}
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Panic(errors.Wrapf(err, "failed to login %s", url))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if !gjson.ValidBytes(respBody) {
		panic("invalid response")
	}

	ret := gjson.ParseBytes(respBody)

	return ret.Get("token").String(), nil
}

func prepareMessage(pbk string) (string, error) {
	msg, err := siwe.InitMessage("w3bstream.com",
		pbk,
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
