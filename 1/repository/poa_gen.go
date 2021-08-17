package repository

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/procsy-tech/attorney/entity"
	"github.com/procsy-tech/attorney/utils/logs"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	ErrPOANotFound = errors.New("POA not found")
)

type (
	POARepository interface {
        New(*entity.POA) (string, error)
        GetByBlockchainID(string) (*entity.POA, error)
		Update(*entity.POA) error
        DeleteByBlockchainID(string) error
        HistoryByBlockchainID(string) ([]entity.POA, error)
        FindItem(string) (*entity.POA, error)
        Find(*entity.POASearchRequest) ([]entity.POA, error)
        List() ([]entity.POA, error)
	}

	POARepositoryImpl struct {
		log logs.Logger
		stub shim.ChaincodeStubInterface
	}
)

func (rep *POARepositoryImpl) New(e *entity.POA) (string, error) {
	log := logs.WithTags(rep.log, "method", "New")

	document := NewPOADocument(e)

	log.Infof("created entity POA with id %s", document.BlockchainID)

	data, err := json.Marshal(document)
	if err != nil {
		return "", err
	}

	err = rep.stub.PutState(document.BlockchainID, data)
	if err != nil {
		return "", err
	}

	return document.BlockchainID, nil
}

func (rep *POARepositoryImpl) GetByBlockchainID(blockchainID string) (*entity.POA, error) {
	log := logs.WithTags(rep.log, "method", "GetByBlockchainID")
	
	log.Infof("searching entity by id %s", blockchainID)

	data, err := rep.stub.GetState(blockchainID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrPOANotFound
	}

	document := new(POADocument)

	err = json.Unmarshal(data, document)
	if err != nil {
		return nil, err
	}
	
	if document.Type != POADocumentType {
		return nil, fmt.Errorf("wrong document type: %s", document.Type)
	}

	return &document.POA, nil
}

func (rep *POARepositoryImpl) Update(e *entity.POA) error {
	log := logs.WithTags(rep.log, "method", "Update")
	
	log.Infof("updating entity with id %s", e.BlockchainID)

	document := POADocument{
			Document{
				Type: POADocumentType,
			},
			*e,
		}

	data, err := json.Marshal(document)
	if err != nil {
		return err
	}

	err = rep.stub.PutState(document.BlockchainID, data)
	if err != nil {
		return err
	}

	return nil
}

func (rep *POARepositoryImpl) DeleteByBlockchainID(blockchainID string) error {
	log := logs.WithTags(rep.log, "method", "DeleteByBlockchainID")
	
	log.Infof("deleting entity with id %s", blockchainID)
	
	err := rep.stub.DelState(blockchainID)
	if err != nil {
		return err
	}

	return nil
}

func (rep *POARepositoryImpl) HistoryByBlockchainID(blockchainID string) ([]entity.POA, error) {
	log := logs.WithTags(rep.log, "method", "HistoryByBlockchainID")
	
	log.Infof("searching entity history by id %s", blockchainID)
	
	iterator, err := rep.stub.GetHistoryForKey(blockchainID)
	if err != nil {
		return nil, errors.New("failed to excute query: " + err.Error())
	}

	defer iterator.Close()

	var entities []entity.POA

	for iterator.HasNext() {
		entry, err := iterator.Next()
		if err != nil {
			return nil, errors.New("failed to get next entry: " + err.Error())
		}

		var document POADocument
		err = json.Unmarshal(entry.Value, &document)
		if err != nil {
			return nil, err
		}
		if document.Type != POADocumentType {
			return nil, fmt.Errorf("wrong document type: %s", document.Type)
		}

		entities = append(entities, document.POA)
	}

	return entities, nil

}

func (rep *POARepositoryImpl) List() ([]entity.POA, error) {
	log := logs.WithTags(rep.log, "method", "List")
	
	log.Infof("getting all POA entities")

	query := fmt.Sprintf(`{"selector":{"type":"%s"}}`, POADocumentType)

	iterator, err := rep.stub.GetQueryResult(query)
	if err != nil {
		return nil, errors.New("failed to excute query: " + err.Error())
	}

	defer iterator.Close()

	var entities []entity.POA

	for iterator.HasNext() {
		entry, err := iterator.Next()
		if err != nil {
			return nil, errors.New("failed to get next entry: " + err.Error())
		}

		var document POADocument
		err = json.Unmarshal(entry.Value, &document)
		if err != nil {
			return nil, err
		}
		if document.Type != POADocumentType {
			return nil, fmt.Errorf("wrong document type: %s", document.Type)
		}

		entities = append(entities, document.POA)
	}

	return entities, nil
}

func (rep *POARepositoryImpl) FindItem(query string) (*entity.POA, error) {
	log := logs.WithTags(rep.log, "method", "FindItem")
	
	log.Infof("finding entity item by query")
	
	iterator, err := rep.stub.GetQueryResult(query)
	if err != nil {
		return nil, errors.New("failed to excute query: " + err.Error())
	}

	defer iterator.Close()

	var entity *entity.POA

	for iterator.HasNext() {
		entry, err := iterator.Next()
		if err != nil {
			return nil, errors.New("failed to get next entry: " + err.Error())
		}

		var document POADocument
		err = json.Unmarshal(entry.Value, &document)
		if err != nil {
			return nil, err
		}
		if document.Type != POADocumentType {
			return nil, fmt.Errorf("wrong document type: %s", document.Type)
		}

		entity = &document.POA
	}

	if entity == nil {
		return nil, ErrPOANotFound
	}

	return entity, nil
}

func (rep *POARepositoryImpl) Find(req *entity.POASearchRequest) ([]entity.POA, error) {
	log := logs.WithTags(rep.log, "method", "FindItem")
	
	log.Infof("finding entity item by search request %+v", req)
	
	querySelector := map[string]interface{}{"type": POADocumentType}
	
	if req.State != nil{
		querySelector["State"] = *req.State
	}
    if req.DateFrom != nil{
		querySelector["DateFrom"] = *req.DateFrom
	}
    if req.DateTo != nil{
		querySelector["DateTo"] = *req.DateTo
	}
    if req.AuthorityINN != nil{
		querySelector["AuthorityINN"] = *req.AuthorityINN
	}
    

	query, err := json.Marshal(querySelector)
	if err != nil {
		return nil, fmt.Errorf("failed to format query: %s", err)
	}

	iterator, err := rep.stub.GetQueryResult(string(query))
	if err != nil {
		return nil, errors.New("failed to excute query: " + err.Error())
	}

	defer iterator.Close()

	var entities []entity.POA

	for iterator.HasNext() {
		entry, err := iterator.Next()
		if err != nil {
			return nil, errors.New("failed to get next entry: " + err.Error())
		}

		var document POADocument
		err = json.Unmarshal(entry.Value, &document)
		if err != nil {
			return nil, err
		}
		if document.Type != POADocumentType {
			return nil, fmt.Errorf("wrong document type: %s", document.Type)
		}

		entities = append(entities, document.POA)
	}

	return entities, nil
}

func NewPOARepositoryImpl(
	log logs.Logger,
	stub shim.ChaincodeStubInterface,
) POARepository {
	return &POARepositoryImpl{
		log: log,
		
		stub: stub,
    	}
}
