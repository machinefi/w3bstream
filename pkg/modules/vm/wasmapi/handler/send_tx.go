package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
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
	TransactionID types.SFID `json:"transactionID,omitempty"`
	Hash          string     `json:"to,omitempty"`
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

	txResp, err := chainCli.SendTXWithOperator(h.chainConf, req.ChainID, req.ChainName, req.To, req.Value, req.Data, req.OperatorName, h.opPool, prj)
	if err != nil {
		l.Error(errors.Wrap(err, "send tx with operator failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	eventType := c.Request.Header.Get("eventType")
	transactionID, err := h.saveTx(&req, prj, eventType, txResp)
	if err != nil {
		l.Error(errors.Wrap(err, "save send tx resp failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &sendTxResp{Hash: txResp.Hash, TransactionID: transactionID})
}

func (h *Handler) saveTx(req *sendTxReq, prj *models.Project, eventType string, txResp *wasm.SendTxResp) (types.SFID, error) {
	id := h.sfid.MustGenSFID()
	m := &models.Transaction{
		RelTransaction: models.RelTransaction{TransactionID: id},
		RelProject:     models.RelProject{ProjectID: prj.ProjectID},
		TransactionInfo: models.TransactionInfo{
			ChainName: txResp.ChainName,
			Nonce:     txResp.Nonce,
			Hash:      txResp.Hash,
			Sender:    txResp.Sender,
			Receiver:  txResp.Receiver,
			Data:      txResp.Data,
			State:     enums.TRANSACTION_STATE__PENDING,
			EventType: eventType,
		},
	}
	return id, m.Create(h.mgrDB)
}
