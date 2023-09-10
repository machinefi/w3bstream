package handler

import (
	"context"
	"net/http"
	"time"

	solclient "github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type sendTxReq struct {
	ChainName    enums.ChainName `json:"chainName"                 binding:"required"`
	To           string          `json:"to,omitempty"`
	Value        string          `json:"value,omitempty"`
	Data         string          `json:"data"                      binding:"required"`
	OperatorName string          `json:"operatorName,omitempty"`
}

type sendTxResp struct {
	TransactionID types.SFID             `json:"transactionID,omitempty"`
	State         enums.TransactionState `json:"state,omitempty"`
	Timeout       bool                   `json:"timeout,omitempty"`
}

func (h *Handler) SendTx(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.SendTx")
	defer l.End()

	var req sendTxReq
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
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

	c.Request.Header.Add("TransactionID", id.String())
	if err := h.setAsync(c); err != nil {
		l.Error(err)
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &sendTxResp{TransactionID: id})
}

func (h *Handler) SendTxAsync(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.SendTxAsync")
	defer l.End()

	chainCli := wasm.MustChainClientFromContext(c.Request.Context())

	var req sendTxReq
	var id types.SFID
	c.ShouldBindBodyWith(&req, binding.JSON)
	id.UnmarshalText([]byte(c.Request.Header.Get("TransactionID")))

	prj := types.MustProjectFromContext(c.Request.Context())
	if req.OperatorName == "" {
		req.OperatorName = operator.DefaultOperatorName
	}

	l = l.WithValues("ProjectID", prj.ProjectID).WithValues("TransactionID", id)
	var state enums.TransactionState
	txResp, err := chainCli.SendTXWithOperator(h.chainConf, 0, req.ChainName, req.To, req.Value, req.Data, req.OperatorName, h.opPool, prj)
	if err != nil {
		state = enums.TRANSACTION_STATE__FAILED
		l.Error(errors.Wrap(err, "send tx with operator failed"))
	} else {
		state = enums.TRANSACTION_STATE__PENDING
	}

	txInfo := models.TransactionInfo{
		State: state,
	}
	if txResp != nil {
		txInfo.Sender = txResp.Sender
		txInfo.Hash = txResp.Hash
		txInfo.Nonce = txResp.Nonce
	}

	m := &models.Transaction{
		RelTransaction:  models.RelTransaction{TransactionID: id},
		TransactionInfo: txInfo,
	}
	if err := m.UpdateByTransactionID(h.mgrDB); err != nil {
		l.Error(errors.Wrap(err, "update transaction db failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	if state == enums.TRANSACTION_STATE__FAILED {
		c.JSON(http.StatusOK, &sendTxResp{TransactionID: id, State: state})
		return
	}

	if err := h.setAsyncAdvance(c, "/system/send_tx/async/state", 10*time.Second); err != nil {
		l.Error(err)
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) SendTxAsyncStateCheck(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.SendTxAsyncStateCheck")
	defer l.End()

	var req sendTxReq
	var id types.SFID
	c.ShouldBindBodyWith(&req, binding.JSON)
	id.UnmarshalText([]byte(c.Request.Header.Get("TransactionID")))
	l = l.WithValues("TransactionID", id)

	m := &models.Transaction{
		RelTransaction: models.RelTransaction{TransactionID: id},
	}
	if err := m.FetchByTransactionID(h.mgrDB); err != nil {
		l.Error(errors.Wrap(err, "fetch by transaction id failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	if time.Since(m.CreatedAt.Time) > 2*time.Hour {
		l.Error(errors.New("transaction timeout"))
		c.JSON(http.StatusOK, &sendTxResp{TransactionID: id, Timeout: true})
		return
	}

	chain := h.chainConf.Chains[m.ChainName]

	state := m.State
	var err error
	if chain.IsSolana() {
		state, err = h.getSolanaState(chain, m.Hash)
	} else {
		state, err = h.getEthState(chain, m.Hash)
	}
	if err != nil {
		l.Error(errors.Wrap(err, "get transaction state failed"))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	if state != m.State {
		m.State = state
		if err := m.UpdateByTransactionID(h.mgrDB); err != nil {
			l.Error(errors.Wrap(err, "update db failed"))
			c.JSON(http.StatusInternalServerError, newErrResp(err))
			return
		}
	}

	if state == enums.TRANSACTION_STATE__FAILED || state == enums.TRANSACTION_STATE__CONFIRMED {
		c.JSON(http.StatusOK, &sendTxResp{TransactionID: id, State: state})
		return
	}

	if err := h.setAsyncAdvance(c, "/system/send_tx/async/state", 10*time.Second); err != nil {
		l.Error(err)
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) getEthState(chain *types.Chain, hash string) (enums.TransactionState, error) {
	client, err := ethclient.Dial(chain.Endpoint)
	if err != nil {
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "dial chain failed")
	}
	nh := common.HexToHash(hash)

	_, p, err := client.TransactionByHash(context.Background(), nh)
	if err != nil {
		if err == ethereum.NotFound {
			return enums.TRANSACTION_STATE__FAILED, nil
		} else {
			return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "get transaction by hash failed")
		}
	} else {
		if p {
			return enums.TRANSACTION_STATE__PENDING, nil
		}
	}

	receipt, err := client.TransactionReceipt(context.Background(), nh)
	if err != nil {
		if err == ethereum.NotFound {
			return enums.TRANSACTION_STATE__IN_BLOCK, nil
		}
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "get transaction receipt failed")
	}
	if receipt.Status == 0 {
		return enums.TRANSACTION_STATE__FAILED, nil
	}
	return enums.TRANSACTION_STATE__CONFIRMED, nil
}

func (h *Handler) getSolanaState(chain *types.Chain, hash string) (enums.TransactionState, error) {
	cli := solclient.NewClient(chain.Endpoint)
	status, err := cli.GetSignatureStatus(context.Background(), hash)
	if err != nil {
		return enums.TRANSACTION_STATE_UNKNOWN, errors.Wrap(err, "query solana transaction failed")
	}
	switch *status.ConfirmationStatus {
	case rpc.CommitmentProcessed:
		return enums.TRANSACTION_STATE__PENDING, nil
	case rpc.CommitmentConfirmed:
		return enums.TRANSACTION_STATE__IN_BLOCK, nil
	case rpc.CommitmentFinalized:
		return enums.TRANSACTION_STATE__CONFIRMED, nil
	}
	return enums.TRANSACTION_STATE_UNKNOWN, errors.New("get solana transaction status failed")
}
