package registry

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/procsy-tech/attorney/service"
	"github.com/procsy-tech/attorney/repository"
	"github.com/procsy-tech/attorney/utils/logs"
)

type (
	// ServiceLocator .
	ServiceLocator interface {
			POAService() service.POAService

		Logger() logs.Logger
		Repository() repository.Repository
	}

	serviceLocatorImpl struct {
		stub shim.ChaincodeStubInterface
	}
)

var (
	POAServiceLog   = shim.NewLogger("POAService")
)
func (sl *serviceLocatorImpl) POAService() service.POAService {
	return service.NewPOAServiceImpl(
		POAServiceLog,
		sl.Repository(),
		)
}

func (sl *serviceLocatorImpl) Logger() logs.Logger {
	return shim.NewLogger("attorney")
}

func (sl *serviceLocatorImpl) Repository() repository.Repository {
	return repository.NewRepositoryImpl(shim.NewLogger("Repository"), sl.stub)
}

func NewServiceLocatorImpl(stub shim.ChaincodeStubInterface) ServiceLocator {
	return &serviceLocatorImpl{stub}
}
