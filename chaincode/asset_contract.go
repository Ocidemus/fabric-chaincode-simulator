package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/Ocidemus/fabric-chaincode-simulator/models"
)


type StubInterface interface {
	PutState(key string, value []byte) error
	GetState(key string) ([]byte, error)
	DelState(key string) error
	GetTxID() string
}

type AssetContract struct{}


func (c *AssetContract) CreateAsset(stub StubInterface, id, owner string, value int) error {
	existing, err := stub.GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read world state: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("asset %q already exists", id)
	}

	asset := models.Asset{ID: id, Owner: owner, Value: value}
	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to serialise asset: %w", err)
	}

	return stub.PutState(id, data)
}


func (c *AssetContract) TransferAsset(stub StubInterface, id, newOwner string) error {
	asset, err := c.getAsset(stub, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to serialise asset: %w", err)
	}

	return stub.PutState(id, data)
}

func (c *AssetContract) UpdateValue(stub StubInterface, id string, newValue int) error {
	asset, err := c.getAsset(stub, id)
	if err != nil {
		return err
	}

	asset.Value = newValue
	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to serialise asset: %w", err)
	}

	return stub.PutState(id, data)
}

func (c *AssetContract) DeleteAsset(stub StubInterface, id string) error {
	existing, err := stub.GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read world state: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("asset %q does not exist", id)
	}
	return stub.DelState(id)
}

func (c *AssetContract) GetAsset(stub StubInterface, id string) (*models.Asset, error) {
	return c.getAsset(stub, id)
}

func (c *AssetContract) AssetExists(stub StubInterface, id string) (bool, error) {
	data, err := stub.GetState(id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}


func (c *AssetContract) getAsset(stub StubInterface, id string) (*models.Asset, error) {
	data, err := stub.GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read world state: %w", err)
	}
	if data == nil {
		return nil, fmt.Errorf("asset %q does not exist", id)
	}

	var asset models.Asset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, fmt.Errorf("failed to deserialise asset: %w", err)
	}
	return &asset, nil
}
