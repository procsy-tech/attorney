package dto

import (
	"github.com/procsy-tech/attorney/entity"
)


type CreateRequest struct{
    
    POA *entity.POA `json:"poa"`
    }

type ConfirmAttorneyRequest struct{
    
    ID string `json:"id"`
    }


type CreateResponse struct{
    
    Result int64 `json:"result"`
    Error string `json:"error"`
}

type ConfirmAttorneyResponse struct{
    Error string `json:"error"`
}


type POASearchRequest struct{
    State  *entity.POAState `json:"state"`
    DateFrom  *string `json:"date_from"`
    DateTo  *string `json:"date_to"`
    AuthorityINN  *string `json:"authority_inn"`
    
}