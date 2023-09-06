package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type sendTxReq struct {
	ChainName    enums.ChainName `json:"chainName"                 binding:"required"`
	To           string          `json:"to,omitempty"`
	Value        string          `json:"value,omitempty"`
	Data         string          `json:"data"                      binding:"required"`
	OperatorName string          `json:"operatorName,omitempty"`
}

type sendTxResp struct {
	TransactionID types.SFID `json:"transactionID,omitempty"`
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

	chain, ok := h.chainConf.Chains[req.ChainName]
	if !ok {
		err := errors.New("blockchain not exist")
		l.Error(err)
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}

	prj := types.MustProjectFromContext(c.Request.Context())
	l = l.WithValues("ProjectID", prj.ProjectID)

	eventType := c.Request.Header.Get("eventType")

	id := h.sfid.MustGenSFID()
	l = l.WithValues("TransactionID", id)

	m := &models.Transaction{
		RelTransaction: models.RelTransaction{TransactionID: id},
		RelProject:     models.RelProject{ProjectID: prj.ProjectID},
		TransactionInfo: models.TransactionInfo{
			ChainName:    chain.Name,
			State:        enums.TRANSACTION_STATE__INIT,
			EventType:    eventType,
			Receiver:     req.To,
			Value:        req.Value,
			Data:         req.Data,
			OperatorName: req.OperatorName,
		},
	}
	if err := m.Create(h.mgrDB); err != nil {
		l.Error(errors.Wrap(err, "create transaction db failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &sendTxResp{TransactionID: id})
}
