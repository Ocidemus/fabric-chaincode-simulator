package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Ocidemus/fabric-chaincode-simulator/stub"
	"github.com/Ocidemus/fabric-chaincode-simulator/utils"
)


type HistoryHandler struct {
	stub *stub.MockStub
}

func NewHistoryHandler(s *stub.MockStub) *HistoryHandler {
	return &HistoryHandler{stub: s}
}

func (h *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only GET is supported",errors.New("not supported"))
		return
	}

	assetID := strings.TrimPrefix(r.URL.Path, "/history/")
	if assetID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing asset ID in path",errors.New("data not found"))
		return
	}

	history := h.stub.GetHistoryForAsset(assetID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"assetId": assetID,
		"history": history,
		"count":   len(history),
	})
}

type AllTxHandler struct {
	stub *stub.MockStub
}

func NewAllTxHandler(s *stub.MockStub) *AllTxHandler {
	return &AllTxHandler{stub: s}
}

func (h *AllTxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only GET is supported",errors.New("not supported"))
		return
	}
	txs := h.stub.GetAllTransactions()
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"transactions": txs,
		"count":        len(txs),
	})
}
