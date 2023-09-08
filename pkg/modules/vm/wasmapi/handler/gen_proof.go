package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/vm/risc0vm"
	"github.com/machinefi/w3bstream/pkg/types"
)

func (h *Handler) GenRisc0Proof(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.GenRisc0Proof")
	defer l.End()

	var req risc0vm.CreateProofReq
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		l.Error(errors.Wrap(err, "decode http request failed"))
		c.JSON(http.StatusBadRequest, newErrResp(err))
		return
	}

	prj := types.MustProjectFromContext(c.Request.Context())

	l = l.WithValues("ProjectID", prj.ProjectID)

	if err := h.setAsync(c); err != nil {
		l.Error(err)
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GenRisc0ProofAsync(c *gin.Context) {
	l := types.MustLoggerFromContext(c.Request.Context())
	_, l = l.Start(c, "wasmapi.handler.GenRisc0ProofAsync")
	defer l.End()

	var req risc0vm.CreateProofReq
	c.ShouldBindJSON(&req)

	prj := types.MustProjectFromContext(c.Request.Context())

	l = l.WithValues("ProjectID", prj.ProjectID)

	rsp, err := risc0vm.CreateProof(c.Request.Context(), &req, h.risc0Conf.Endpoint, h.risc0Conf.CreateProofPath)
	if err != nil {
		l.Error(errors.Wrap(err, fmt.Sprintf("send risc0 server request %s failed", h.risc0Conf.CreateProofPath)))
		c.JSON(http.StatusInternalServerError, newErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &rsp)
}
