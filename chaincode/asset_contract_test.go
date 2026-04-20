package chaincode

import (
	"testing"

	"github.com/Ocidemus/fabric-chaincode-simulator/stub"
)
func TestCreateAsset(t *testing.T) {
    stub := stub.NewMockStub()
    contract := &AssetContract{}
    
    err := contract.CreateAsset(stub, "asset1", "Alice", 500)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    
    asset, err := contract.GetAsset(stub, "asset1")
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if asset.Owner != "Alice" {
        t.Errorf("expected owner Alice, got %s", asset.Owner)
    }
}