package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "User_id"

// SimpleChaincode implements the fabric-contract-api-go programming model
type SimpleChaincode struct {
	contractapi.Contract
}

type Asset struct {
	DocType   string   `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ID        string   `json:"ID"`      //the field tags are needed to keep case from bouncing around
	Size      int      `json:"size"`
	Owner     string   `json:"owner"`
	Access    []string `json:"access"`
	S3Bucket  string   `json:"s3Bucket"`
	ObjectKey string   `json:"objectKey"`
	Name      string   `json:"name"`
	Content   string   `json:"content"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

func (t *SimpleChaincode) sha256Hash(ctx contractapi.TransactionContextInterface, input string) (string, error) {
	// Convert the input string to a byte slice (required by the SHA-256 function)
	inputBytes := []byte(input)

	// Create a new SHA-256 hash object
	hash := sha256.New()

	// Calculate the hash value
	_, err := hash.Write(inputBytes)
	if err != nil {
		return "", fmt.Errorf("error calculating SHA-256 hash: %w", err)
	}

	// Get the finalized hash result as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

func (t *SimpleChaincode) compareSHA256Hashes(hash1, hash2 string) (bool, error) {
	return hash1 == hash2, nil
}

func (t *SimpleChaincode) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string, size int, owner string, name string, content string) error {
	exists, err := t.AssetExists(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		return fmt.Errorf("asset already exists: %s", assetID)
	}

	encode, err := t.sha256Hash(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to calculate sha256: %v", err)
	}

	asset := &Asset{
		DocType: "asset",
		ID:      assetID,
		Size:    size,
		Owner:   owner,
		Access:  []string{owner},
		Content: encode,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(assetID, assetBytes)
	if err != nil {
		return err
	}

	OwnerNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Owner, asset.ID})
	if err != nil {
		return err
	}

	value := []byte{0x00}
	return ctx.GetStub().PutState(OwnerNameIndexKey, value)
}

func (t *SimpleChaincode) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string, expected string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", assetID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", assetID)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	encode, err := t.sha256Hash(ctx, expected)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate sha256: %v", err)
	}

	tamper, err := t.compareSHA256Hashes(asset.Content, encode)
	if err != nil {
		return nil, fmt.Errorf("failed to compare file contents %v", err)
	}
	if !tamper {
		return nil, fmt.Errorf("file has been tampered")
	}

	return &asset, nil
}

func (t *SimpleChaincode) ReadAssetchaincode(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", assetID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", assetID)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (t *SimpleChaincode) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	asset, err := t.ReadAssetchaincode(ctx, assetID)
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s: %v", assetID, err)
	}

	OwnerNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Owner, asset.ID})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(OwnerNameIndexKey)
}

// TransferAsset transfers an asset by setting a new owner name on the asset
func (t *SimpleChaincode) TransferAsset(ctx contractapi.TransactionContextInterface, assetID, newOwner string) error {
	asset, err := t.ReadAssetchaincode(ctx, assetID)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(assetID, assetBytes)
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Asset, error) {
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// QueryAssetsByOwner queries for assets based on the owners name.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Parameterized rich query
func (t *SimpleChaincode) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Asset, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"asset","owner":"%s"}}`, owner)
	return getQueryResultForQueryString(ctx, queryString)
}

// QueryAssets uses a query string to perform a query for assets.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryAssetsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Ad hoc rich query
func (t *SimpleChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// AssetExists returns true when asset with given ID exists in the ledger.
func (t *SimpleChaincode) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", assetID, err)
	}

	return assetBytes != nil, nil
}

// InitLedger creates the initial set of assets in the ledger.
func (t *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DocType: "asset", ID: "asset0", Size: 0, Owner: "own", Access: []string{"none"}, S3Bucket: "qwertyu", ObjectKey: "asdfghj", Name: "file0", Content: "Hello excited to get started !!"},
	}

	for _, asset := range assets {
		err := t.CreateAsset(ctx, asset.ID, asset.Size, asset.Owner, asset.Name, asset.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}

/*
// TransferAssetByColor will transfer assets of a given color to a certain new owner.
// Uses GetStateByPartialCompositeKey (range query) against color~name 'index'.
// Committing peers will re-execute range queries to guarantee that result sets are stable
// between endorsement time and commit time. The transaction is invalidated by the
// committing peers if the result set has changed between endorsement time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// Example: GetStateByPartialCompositeKey/RangeQuery
func (t *SimpleChaincode) TransferAssetByColor(ctx contractapi.TransactionContextInterface, color, newOwner string) error {
	// Execute a key range query on all keys starting with 'color'
	coloredAssetResultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index, []string{color})
	if err != nil {
		return err
	}
	defer coloredAssetResultsIterator.Close()

	for coloredAssetResultsIterator.HasNext() {
		responseRange, err := coloredAssetResultsIterator.Next()
		if err != nil {
			return err
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}

		if len(compositeKeyParts) > 1 {
			returnedAssetID := compositeKeyParts[1]
			asset, err := t.ReadAsset(ctx, returnedAssetID)
			if err != nil {
				return err
			}
			asset.Owner = newOwner
			assetBytes, err := json.Marshal(asset)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(returnedAssetID, assetBytes)
			if err != nil {
				return fmt.Errorf("transfer failed for asset %s: %v", returnedAssetID, err)
			}
		}
	}

	return nil
}
*/

/*
// GetAssetHistory returns the chain of custody for an asset since issuance.
func (t *SimpleChaincode) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
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
				ID: assetID,
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
*/
