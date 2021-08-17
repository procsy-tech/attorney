package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/kbkontrakt/hlfabric-ccdevkit/debug"
	"github.com/kbkontrakt/hlfabric-ccdevkit/logs"
	"github.com/kbkontrakt/hlfabric-ccdevkit/utils"
	"github.com/procsy-tech/attorney/api"
	"github.com/procsy-tech/attorney/registry"
)

const (
	chaincodeVersion      = "0.1.0"
	attorneyCollectionName = "attorneys"
)

// Config .
type Config struct {
	Version     string `json:"version"`
	ChaincodeID string `json:"chaincode_id"`
}

type attorneyChaincode struct {
	logger    logs.Logger
}

func NewattorneyChaincode() *attorneyChaincode {
	chaincode := attorneyChaincode{
		logger:    logs.WithTags(shim.NewLogger("main"), "module", "attorneycc"),
	}

	return &chaincode
}

// Init .
func (chaincode *attorneyChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger := chaincode.logger

	fn, args := stub.GetFunctionAndParameters()
	logger = logs.WithTags(logger, "method", "Init", "fn", fn)
	if logger.IsEnabledFor(shim.LogDebug) {
		logger.Debugf("Call with args [%v]", args)
	} else {
		logger.Info("Call")
	}

	return shim.Success(nil)
}

// Invoke .
func (chaincode *attorneyChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	logger := chaincode.logger
	fn, args, err := utils.GetFnArgsOrFromTransientMap(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	logger = logs.WithTags(logger, "method", "Invoke", "fn", fn)

	defer func() {
		if err := recover(); err != nil {
			response = shim.Error("Internal error was occurred, please see log for more details")
			logger.Criticalf("Panic was occurred with [%+v] and stacktrace=[%s]", err, debug.GetStacktrace(false))
		}

		if logger.IsEnabledFor(shim.LogDebug) {
			logger.Debugf("Call with args [%v] and response is status=[%d] msg=[%s] payload=[%s]", args,
				response.GetStatus(),
				response.GetMessage(),
				response.GetPayload())
		} else {
			logger.Infof("Call %d", response.GetStatus())
		}
	}()

	return chaincode.handleByRoute(stub, fn, args)
}

func (chaincode *attorneyChaincode) handleByRoute(stub shim.ChaincodeStubInterface, fn string, args []string) peer.Response {
	svcFactory := registry.NewServiceLocatorImpl(stub)

	var err error
	var payload []byte

	switch fn {
		case api.Create:
        payload, err = chaincode.Create(svcFactory, args)
    case api.ConfirmAttorney:
        payload, err = chaincode.ConfirmAttorney(svcFactory, args)
    

	case "_debug":
		payload, err = debug.Invoke(stub, args)
	default:
		return shim.Error("unsupported function")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(payload)
}

func main() {
	chaincode := NewattorneyChaincode()

	shim.SetupChaincodeLogging()

	chaincode.logger.Info("Start")

	if err := shim.Start(chaincode); err != nil {
		fmt.Printf("Error starting attorneyChaincode: %s", err)
	}
}
