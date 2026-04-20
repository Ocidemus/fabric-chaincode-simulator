package models

import "time"

type TxType string

const (
	TxCreate   TxType = "CREATE"
	TxTransfer TxType = "TRANSFER"
	TxDelete   TxType = "DELETE"
	TxRead     TxType = "READ"
)

type Transaction struct {
	TxID      string    `json:"txId"`
	Function  string    `json:"function"`
	AssetID   string    `json:"assetId"`
	Type      TxType    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}
