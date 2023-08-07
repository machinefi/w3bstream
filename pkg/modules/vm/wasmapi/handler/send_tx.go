package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type sendTxReq struct {
	ChainID      uint64          `json:"chainID,omitempty"`
	ChainName    enums.ChainName `json:"chainName,omitempty"`
	To           string          `json:"to,omitempty"`
	Value        string          `json:"value,omitempty"`
	Data         string          `json:"data"       binding:"required"`
	OperatorName string          `json:"operatorName,omitempty"`
}

type sendTxResp struct {
	Hash string `json:"to,omitempty"`
}

func (h *Handler) SendTx(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.SendTx")
	defer l.End()

	chainCli := wasm.MustChainClientFromContext(c.Request.Context())

	var req sendTxReq
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

	prj := types.MustProjectFromContext(c.Request.Context())

	l = l.WithValues("ProjectID", prj.ProjectID)

	if req.OperatorName == "" {
		req.OperatorName = operator.DefaultOperatorName
	}

	hash, err := chainCli.SendTXWithOperator(h.chainConf, req.ChainID, req.ChainName, req.To, req.Value, req.Data, req.OperatorName)
	if err != nil {
		l.Error(errors.Wrap(err, "send tx with operator failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &sendTxResp{Hash: hash})
}
