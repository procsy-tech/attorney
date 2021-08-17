package repository

import (
	"encoding/json"

	"github.com/golang/protobuf/ptypes/timestamp"
	qr "github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	attorneyCollectionName = "attorney_pdc"
)

// MustWrapAsPrivateStub wraps stub with proxy common state methods to a private data
func MustWrapAsPrivateStub(stub shim.ChaincodeStubInterface /*, specification interface{}*/) shim.ChaincodeStubInterface {
	// switch specification { cases for selecting properly priv collection }
	return NewPrivateStubDecorator(attorneyCollectionName, stub)
}

// MustWrapAsPrivateStubWithHistory wraps stub with private data history
 func MustWrapAsPrivateStubWithHistory(stub shim.ChaincodeStubInterface /*, specification interface{}*/) shim.ChaincodeStubInterface {
 	// switch specification { cases for selecting properly priv collection }
 	return NewPrivateStubDecorator(attorneyCollectionName,
 		NewPrivateHistoryStubDecorator(attorneyCollectionName,
 			NewPrivateHistoryArrayAppendStrategy("", "_HIST"), stub))
}

type privateStubDecorator struct {
	shim.ChaincodeStubInterface
	collectionName string
}

func (s *privateStubDecorator) GetState(key string) ([]byte, error) {
	return s.ChaincodeStubInterface.GetPrivateData(s.collectionName, key)
}
func (s *privateStubDecorator) PutState(key string, value []byte) error {
	return s.ChaincodeStubInterface.PutPrivateData(s.collectionName, key, value)
}
func (s *privateStubDecorator) DelState(key string) error {
	return s.ChaincodeStubInterface.DelPrivateData(s.collectionName, key)
}
func (s *privateStubDecorator) SetStateValidationParameter(key string, ep []byte) error {
	return s.ChaincodeStubInterface.SetPrivateDataValidationParameter(s.collectionName, key, ep)
}
func (s *privateStubDecorator) GetStateValidationParameter(key string) ([]byte, error) {
	return s.ChaincodeStubInterface.GetPrivateDataValidationParameter(s.collectionName, key)
}
func (s *privateStubDecorator) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetPrivateDataByRange(s.collectionName, startKey, endKey)
}
func (s *privateStubDecorator) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	panic("not supported")
}
func (s *privateStubDecorator) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetPrivateDataByPartialCompositeKey(s.collectionName, objectType, keys)
}
func (s *privateStubDecorator) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	panic("not supported")
}
func (s *privateStubDecorator) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetPrivateDataQueryResult(s.collectionName, query)
}
func (s *privateStubDecorator) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	panic("not supported")
}
func (s *privateStubDecorator) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetHistoryForKey(key)
}

// NewPrivateStubDecorator decorates stub for using private data collection as a source.
func NewPrivateStubDecorator(collectionName string, stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &privateStubDecorator{
		ChaincodeStubInterface: stub,
		collectionName:         collectionName,
	}
}

type (
	// PrivateHistoryStrategy .
	PrivateHistoryStrategy interface {
		// Append .
		Append(stub shim.ChaincodeStubInterface, collection, key string, value []byte, isDelete bool) error
		// GetIterator .
		GetIterator(stub shim.ChaincodeStubInterface, collection, key string) (shim.HistoryQueryIteratorInterface, error)
	}

	privateHistoryArrayAppendStrategy struct {
		keysPrefix string
		keysSuffix string
	}

	privateHistoryArrayAppendIterator struct {
		hist []qr.KeyModification
		inx  int
	}

	privateHistoryStubDecorator struct {
		shim.ChaincodeStubInterface
		history    PrivateHistoryStrategy
		collection string
	}

	keyValueHistory struct {
		TxID      string               `json:"i,omitempty"`
		Value     string               `json:"v,omitempty"`
		Timestamp *timestamp.Timestamp `json:"t,omitempty"`
		IsDelete  bool                 `json:"d,omitempty"`
	}
)

// Strategies
func (*privateHistoryArrayAppendStrategy) tryToFindPreviousActualItem(hist []keyValueHistory, item *keyValueHistory) {
	for inx := len(hist) - 1; inx >= 0; inx-- {
		if hist[inx].Value != "" {
			item.Value = hist[inx].Value
		}
	}
}

func (a *privateHistoryArrayAppendStrategy) Append(stub shim.ChaincodeStubInterface, collection, key string, value []byte, isDelete bool) error {
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return err
	}

	data, err := stub.GetPrivateData(collection, a.keysPrefix+key+a.keysSuffix)
	if err != nil {
		return err
	}

	// @TODO: add marshal strategy

	hist := []keyValueHistory{}
	if data != nil {
		err = json.Unmarshal(data, &hist)
		if err != nil {
			return err
		}
	}

	newHistItem := keyValueHistory{
		TxID:      stub.GetTxID(),
		Timestamp: timestamp,
		IsDelete:  isDelete,
		Value:     string(value),
	}
	if isDelete {
		// @TODO: Is need here?
		a.tryToFindPreviousActualItem(hist, &newHistItem)
	}

	hist = append([]keyValueHistory{newHistItem}, hist...)

	data, err = json.Marshal(hist)
	if err != nil {
		return err
	}

	return stub.PutPrivateData(collection, a.keysPrefix+key+a.keysSuffix, data)
}

func (a *privateHistoryArrayAppendStrategy) GetIterator(stub shim.ChaincodeStubInterface, collection, key string) (shim.HistoryQueryIteratorInterface, error) {
	data, err := stub.GetPrivateData(collection, a.keysPrefix+key+a.keysSuffix)
	if err != nil {
		return nil, err
	}

	// @TODO: add marshal strategy

	hist := []qr.KeyModification{}
	if data != nil {
		rawHist := []keyValueHistory{}
		err = json.Unmarshal(data, &rawHist)
		if err != nil {
			return nil, err
		}
		for _, item := range rawHist {
			hist = append(hist, qr.KeyModification{
				TxId:      item.TxID,
				Value:     []byte(item.Value),
				Timestamp: item.Timestamp,
				IsDelete:  item.IsDelete,
			})
		}
	}

	return &privateHistoryArrayAppendIterator{hist, len(hist)}, nil
}

func (i *privateHistoryArrayAppendIterator) HasNext() bool {
	return i.inx > 0
}
func (i *privateHistoryArrayAppendIterator) Next() (*qr.KeyModification, error) {
	if !i.HasNext() {
		return nil, nil
	}
	i.inx--
	return &i.hist[i.inx], nil
}
func (i *privateHistoryArrayAppendIterator) Close() error {
	return nil
}

// NewPrivateHistoryArrayAppendStrategy .
func NewPrivateHistoryArrayAppendStrategy(keysPrefix, keysSuffix string) PrivateHistoryStrategy {
	return &privateHistoryArrayAppendStrategy{
		keysPrefix: keysPrefix,
		keysSuffix: keysSuffix,
	}
}

// Stub

func (s *privateHistoryStubDecorator) PutPrivateData(collection, key string, value []byte) error {
	err := s.history.Append(s.ChaincodeStubInterface, s.collection, key, value, false)
	if err != nil {
		return err
	}
	return s.ChaincodeStubInterface.PutPrivateData(collection, key, value)
}
func (s *privateHistoryStubDecorator) DelPrivateData(collection, key string) error {
	err := s.history.Append(s.ChaincodeStubInterface, s.collection, key, nil, true)
	if err != nil {
		return err
	}
	return s.ChaincodeStubInterface.DelPrivateData(s.collection, key)
}
func (s *privateHistoryStubDecorator) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return s.history.GetIterator(s.ChaincodeStubInterface, s.collection, key)
}

// NewPrivateHistoryStubDecorator decorates stub for using private data collection with history request.
func NewPrivateHistoryStubDecorator(collectionName string, histStrategy PrivateHistoryStrategy, stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &privateHistoryStubDecorator{
		collection: collectionName,
		history:    histStrategy,

		ChaincodeStubInterface: stub,
	}
}
