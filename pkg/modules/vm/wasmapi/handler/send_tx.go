package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/projectoperator"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type sendTxReq struct {
	ChainID      uint32 `json:"chainID"    binding:"required"`
	To           string `json:"to"         binding:"required"`
	Value        string `json:"value"      binding:"required"`
	Data         string `json:"data"       binding:"required"`
	OperatorName string `json:"operatorName,omitempty"`
}

type sendTxResp struct {
	Hash string `json:"to,omitempty"`
}

func (h *Handler) SendTx(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.SendTx")
	defer l.End()

	var req sendTxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Error(errors.Wrap(err, "decode http request failed"))
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}

	prj := types.MustProjectFromContext(c.Request.Context())

	l = l.WithValues("ProjectID", prj.ProjectID)

	ctx := contextx.WithContextCompose(
		types.WithLoggerContext(l),
		types.WithMgrDBExecutorContext(h.mgrDB),
		types.WithETHClientConfigContext(h.ethCli),
	)(context.Background())

	prjOp, err := projectoperator.GetByProject(ctx, prj.ProjectID)
	if err != nil && err != status.ProjectOperatorNotFound {
		l.Error(errors.Wrap(err, "fetch project operator failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}
	accOp, err := operator.ListByCond(ctx, &operator.CondArgs{AccountID: prj.RelAccount.AccountID})
	if err != nil {
		l.Error(errors.Wrap(err, "fetch operators failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	chainClient := wasm.NewChainClient(ctx, prj, accOp, prjOp)

	if req.OperatorName == "" {
		req.OperatorName = operator.DefaultOperatorName
	}

	hash, err := chainClient.SendTXWithOperator(req.ChainID, req.To, req.Value, req.Data, req.OperatorName)
	if err != nil {
		l.Error(errors.Wrap(err, "send tx with operator failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &sendTxResp{Hash: hash})
}
