
package stub

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Ocidemus/fabric-chaincode-simulator/models"
)


type MockStub struct {
	mu      sync.RWMutex
	state   map[string][]byte       
	history map[string][]models.Transaction 
	txID    string
}


func NewMockStub() *MockStub {
	return &MockStub{
		state:   make(map[string][]byte),
		history: make(map[string][]models.Transaction),
	}
}


func (s *MockStub) PutState(key string, value []byte) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state[key] = value
	return nil
}


func (s *MockStub) GetState(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.state[key]
	if !ok {
		return nil, nil 
	}
	return val, nil
}


func (s *MockStub) DelState(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state[key]; !ok {
		return fmt.Errorf("key %q does not exist", key)
	}
	delete(s.state, key)
	return nil
}


func (s *MockStub) SetTxID(txID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.txID = txID
}


func (s *MockStub) GetTxID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.txID
}


func (s *MockStub) RecordTx(tx models.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history[tx.AssetID] = append(s.history[tx.AssetID], tx)
}

func (s *MockStub) GetHistoryForAsset(assetID string) []models.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.history[assetID]
}


func (s *MockStub) GetAllTransactions() []models.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []models.Transaction
	for _, txs := range s.history {
		all = append(all, txs...)
	}
	return all
}


func NewTxID(function, assetID string) string {
	return fmt.Sprintf("txn-%s-%s-%d", function, assetID, time.Now().UnixNano())
}


func MarshalAsset(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalAsset(data []byte, dst interface{}) error {
	return json.Unmarshal(data, dst)
}
