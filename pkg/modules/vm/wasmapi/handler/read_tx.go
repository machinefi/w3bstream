package handler

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/types"
)

type readTxReq struct {
	ChainID uint32 `json:"chainID"    binding:"required"`
	Hash    string `json:"hash"       binding:"required"`
}

type readTxResp struct {
	Transaction *ethtypes.Transaction `json:"transaction,omitempty"`
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

	l = l.WithValues("chain_id", req.ChainID)

	chainAddress, ok := h.ethCli.Clients[req.ChainID]
	if !ok {
		err := errors.New("blockchain not exist")
		l.Error(err)
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}

	client, err := ethclient.Dial(chainAddress)
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

	c.JSON(http.StatusOK, &readTxResp{Transaction: tx})
}
