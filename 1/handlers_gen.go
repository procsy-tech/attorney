package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/procsy-tech/attorney/dto"
	"github.com/procsy-tech/attorney/registry"
)
// Create .
func (chaincode *attorneyChaincode) Create(svcFactory registry.ServiceLocator, args []string) ([]byte, error) {
	payload := args[0]
	if len(payload) == 0 {
		return nil, errors.New("empty request payload")
	}

	data := []byte(payload)
	var request dto.CreateRequest

	err := json.Unmarshal(data, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request payload: %s", err)
	}
	
	result, err := svcFactory.POAService().Create(request.POA)
	if err != nil{
		chaincode.logger.Infof("error invoking method Create: %s", err)
	}
	response := dto.CreateResponse{
	
    	Result: result,
    }
	if err != nil{
		response.Error = err.Error()
	}
	resultData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return resultData, nil
}
// ConfirmAttorney .
func (chaincode *attorneyChaincode) ConfirmAttorney(svcFactory registry.ServiceLocator, args []string) ([]byte, error) {
	payload := args[0]
	if len(payload) == 0 {
		return nil, errors.New("empty request payload")
	}

	data := []byte(payload)
	var request dto.ConfirmAttorneyRequest

	err := json.Unmarshal(data, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request payload: %s", err)
	}
	
	err = svcFactory.POAService().ConfirmAttorney(request.ID)
	if err != nil{
		chaincode.logger.Infof("error invoking method ConfirmAttorney: %s", err)
	}
	response := dto.ConfirmAttorneyResponse{
	}
	if err != nil{
		response.Error = err.Error()
	}
	resultData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return resultData, nil
}



