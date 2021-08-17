package service

import (
	
		"github.com/procsy-tech/attorney/entity"
	

)

// POAService interface.
type POAService interface {
	Create(POA *entity.POA) (int64, error)
	ConfirmAttorney(ID string) error
	
}


