package repository

import(
	"github.com/procsy-tech/attorney/utils/logs"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	Repository interface{
		POARepository() POARepository
		}

	repositoryImpl struct {
		log logs.Logger
		stub shim.ChaincodeStubInterface
	}
)

func (rep *repositoryImpl)POARepository() POARepository{
	return NewPOARepositoryImpl(logs.WithTags(rep.log, "entity", "POA"), rep.stub)
}
func NewRepositoryImpl(
	log logs.Logger,
	stub shim.ChaincodeStubInterface,
) Repository {
	return &repositoryImpl{
		log: log,
		stub: stub,
	}
}
