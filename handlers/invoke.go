package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Ocidemus/fabric-chaincode-simulator/chaincode"
	"github.com/Ocidemus/fabric-chaincode-simulator/models"
	"github.com/Ocidemus/fabric-chaincode-simulator/stub"
	"github.com/Ocidemus/fabric-chaincode-simulator/utils"
)


type InvokeRequest struct {
	Function string            `json:"function"`
	Args     map[string]string `json:"args"`
}

type InvokeHandler struct {
	stub     *stub.MockStub
	contract *chaincode.AssetContract
}

func NewInvokeHandler(s *stub.MockStub, c *chaincode.AssetContract) *InvokeHandler {
	return &InvokeHandler{stub: s, contract: c}
}

func (h *InvokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only POST is supported",errors.New("not supported"))
		return
	}

	var req InvokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid JSON",err)
		return
	}

	
	assetID := req.Args["id"]
	txID := stub.NewTxID(req.Function, assetID)
	h.stub.SetTxID(txID)

	tx := models.Transaction{
		TxID:      txID,
		Function:  req.Function,
		AssetID:   assetID,
		Timestamp: time.Now(),
	}

	
	var result interface{}
	var execErr error

	switch req.Function {

	case "CreateAsset":
		value := utils.ParseInt(req.Args["value"])
		execErr = h.contract.CreateAsset(h.stub, req.Args["id"], req.Args["owner"], value)
		tx.Type = models.TxCreate

	case "TransferAsset":
		execErr = h.contract.TransferAsset(h.stub, req.Args["id"], req.Args["newOwner"])
		tx.Type = models.TxTransfer

	case "UpdateValue":
		value := utils.ParseInt(req.Args["value"])
		execErr = h.contract.UpdateValue(h.stub, req.Args["id"], value)
		tx.Type = models.TxCreate 
	case "DeleteAsset":
		execErr = h.contract.DeleteAsset(h.stub, req.Args["id"])
		tx.Type = models.TxDelete

	case "GetAsset":
		result, execErr = h.contract.GetAsset(h.stub, req.Args["id"])
		tx.Type = models.TxRead

	case "AssetExists":
		result, execErr = h.contract.AssetExists(h.stub, req.Args["id"])
		tx.Type = models.TxRead

	default:
		utils.RespondWithError(w, http.StatusBadRequest, "unknown function: "+req.Function,errors.New("bad request"))
		return
	}

	
	tx.Success = execErr == nil
	if execErr != nil {
		tx.Error = execErr.Error()
	}
	h.stub.RecordTx(tx)

	if execErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest,"not excecutable" ,execErr)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"txId":   txID,
		"result": result,
	})
}
