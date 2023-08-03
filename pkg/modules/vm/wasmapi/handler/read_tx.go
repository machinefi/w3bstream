package handler

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
)

type readTxReq struct {
	ChainID   uint32          `json:"chainID"`
	ChainName enums.ChainName `json:"chainName"`
	Hash      string          `json:"hash"       binding:"required"`
}

type readEthTxResp struct {
	Transaction *ethtypes.Transaction `json:"transaction,omitempty"`
}

type readSolanaTxResp struct {
	Transaction *rpc.GetTransactionResult `json:"result,omitempty"`
}

func (h *Handler) ReadTx(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.ReadTx")
	defer l.End()

	var req readTxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Error(errors.Wrap(err, "decode http request failed"))
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}
	if req.ChainID == 0 && req.ChainName == "" {
		err := errors.New("missing chain param")
		l.Error(err)
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}

	l = l.WithValues("chain_id", req.ChainID, "chain_name", req.ChainName)

	chain, ok := h.chainConf.GetChain(uint64(req.ChainID), req.ChainName)
	if !ok {
		err := errors.New("blockchain not exist")
		l.Error(err)
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}

	var resp any

	switch {
	case chain.ChainID != 0:
		client, err := ethclient.Dial(chain.Endpoint)
		if err != nil {
			l.Error(errors.Wrap(err, "dial chain address failed"))
			c.JSON(http.StatusInternalServerError, newErrResp(err))
			return
		}

		tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(req.Hash))
		if err != nil {
			l.Error(errors.Wrap(err, "query transaction failed"))
			c.JSON(http.StatusInternalServerError, newErrResp(err))
			return
		}
		resp = &readEthTxResp{Transaction: tx}

	case chain.Name == enums.SOLANA_DEVNET || chain.Name == enums.SOLANA_TESTNET || chain.Name == enums.SOLANA_MAINNET_BETA:
		client := rpc.New(chain.Endpoint)
		txSig, err := solana.SignatureFromBase58(req.Hash)
		if err != nil {
			l.Error(errors.Wrap(err, "decode http request hash failed"))
			c.JSON(http.StatusBadRequest, newErrResp(err))
			return
		}

		tx, err := client.GetTransaction(context.Background(), txSig, nil)
		if err != nil {
			l.Error(errors.Wrap(err, "query transaction failed"))
			c.JSON(http.StatusInternalServerError, newErrResp(err))
			return
		}
		resp = &readSolanaTxResp{Transaction: tx}

	default:
		err := errors.New("server error")
		l.Error(err)
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
