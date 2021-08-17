package service

import (
	"errors"
	
		"github.com/procsy-tech/attorney/entity"
	
	"github.com/procsy-tech/attorney/repository"
	"github.com/procsy-tech/attorney/utils/logs"
)

func NewPOAServiceImpl(
	log logs.Logger,
	rep repository.Repository,
) POAService {
	return &POAServiceImpl{
		log,
		rep,
	}
}

type POAServiceImpl struct {
	log logs.Logger
	rep repository.Repository
}

// Create .
func (svc *POAServiceImpl) Create(POA *entity.POA) (int64, error) {
    // implement method logic .
	return 0, errors.New("not implemented")
}
// ConfirmAttorney .
func (svc *POAServiceImpl) ConfirmAttorney(ID string) error {
    // implement method logic .
	return errors.New("not implemented")
}

