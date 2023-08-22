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

type sendUserOpReq struct {
	ChainID      uint64          `json:"chainID,omitempty"`
	ChainName    enums.ChainName `json:"chainName,omitempty"`
	OperatorName string          `json:"operatorName,omitempty"`
	PayMasterKey string          `json:"payMasterKey,omitempty"`
	Data         string          `json:"data"       binding:"required"`
}

type sendUserOpResp struct {
	Hash string `json:"to,omitempty"`
}

func (h *Handler) SendUserOp(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.SendUserOp")
	defer l.End()

	chainCli := wasm.MustChainClientFromContext(c.Request.Context())

	var req sendUserOpReq
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

	hash, err := chainCli.SendUserOpWithOperator(h.chainConf, req.ChainID, req.ChainName, req.Data, req.OperatorName)
	if err != nil {
		l.Error(errors.Wrap(err, "send user operation with operator failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &sendUserOpResp{Hash: hash})
}
