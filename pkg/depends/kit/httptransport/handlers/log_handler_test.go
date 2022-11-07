package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/mock"
)

func ExampleLogHandler() {
	ctx := log.WithLogger(context.Background(), log.Std())

	var handle http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(20 * time.Millisecond)

		switch req.Method {
		case http.MethodGet:
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(`{"status":"ok"}`))
		case http.MethodPost:
			rw.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte(`{"key":"StatusBadRequest","msg":"something wrong"}`))
		case http.MethodPut:
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(`{"key":"StatusInternalServerError","msg":"internal server error"}`))
		}
	}

	handler := LogHandler()(handle).(*loggerHandler)

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPost} {
		req, _ := http.NewRequestWithContext(ctx, method, "/", nil)
		handler.ServeHTTP(mock.NewMockResponseWriter(), req)
	}
	// Output:
}
