package chaincode

import (
	"encoding/json"
	"fmt"
    "time"
	"log"
	"strconv"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
/*type Asset struct {
	AppraisedValue int    `json:"AppraisedValue"`
	Color          string `json:"Color"`
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Size           int    `json:"Size"`
}*/
type Asset struct {
	DealerID       string    `json:"DealerID"`
	MSISDN         int    `json:"MSISDN"`
	MPIN           int    `json:"MPIN"`
	Balance        float64    `json:"Balance"`
	Status         string `json:"Status"`
	TransAmount    float64    `json:"TransAmount"`
	TransType      string `json:"TransType"`
	Remarks        string `json:"Remarks"`
}

type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DealerID: "Asset1", MSISDN: 10, MPIN: 1024, Balance: 100000.9, Status: "active", TransAmount:0.0 ,TransType: "" ,Remarks:"initial deposit"},
		{DealerID: "Asset2", MSISDN:11, MPIN: 2048, Balance: 56000.6, Status: "active", TransAmount:0.0 ,TransType: "",Remarks:"intial deposit"},
		{DealerID: "Asset3", MSISDN: 12, MPIN: 0507, Balance: 230000.8, Status: "active", TransAmount:0.0 ,TransType: "" ,Remarks:"initial deposit"},
		{DealerID: "Asset4", MSISDN: 13, MPIN: 7557, Balance: 576000.1, Status: "active", TransAmount:0.0 ,TransType: "",Remarks:"intial deposit"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.DealerID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
                                
	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerid string, msisdn int, mpin int, balance float64, status string, transamount float64, transtype string, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerid)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", dealerid)
	}

	asset := Asset{
		DealerID:       dealerid,
		MSISDN:         msisdn,
		MPIN:           mpin,
		Balance:        balance,
		Status:         status,
        TransAmount:    transamount,
		TransType:      transtype,
		Remarks:        remarks,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerid, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerid string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", dealerid)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerid string, msisdn int, mpin int, balance float64, status string, transamount float64, transtype string, remarks string) error {
	exists, err := s.AssetExists(ctx, dealerid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dealerid)
	}

	// overwriting original asset with new asset
	asset := Asset{
		DealerID:       dealerid,
		MSISDN:         msisdn,
		MPIN:           mpin,
		Balance:        balance,
		Status:         status,
        TransAmount:    transamount,
		TransType:      transtype,
		Remarks:        remarks,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dealerid, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, dealerid string) error {
	exists, err := s.AssetExists(ctx, dealerid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dealerid)
	}

	return ctx.GetStub().DelState(dealerid)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dealerid string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
/*func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}*/

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, dealerid string, newOwner int) (string, error) {
	asset, err := s.ReadAsset(ctx, dealerid)
	if err != nil {
		return "", err
	}

	oldOwner := asset.MPIN
	asset.MPIN = newOwner

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(dealerid, assetJSON)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(oldOwner), nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, dealerid string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", dealerid)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerid)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				DealerID: dealerid,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}