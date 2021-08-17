package entity


import (
	"fmt"
)


type POA struct{
    BlockchainID string
    State  POAState `json:"state"`
    DateFrom  string `json:"date_from"`
    DateTo  string `json:"date_to"`
    AuthorityINN  string `json:"authority_inn"`
    
}

type POASearchRequest struct{
    BlockchainID string
    State  *POAState `json:"state"`
    DateFrom  *string `json:"date_from"`
    DateTo  *string `json:"date_to"`
    AuthorityINN  *string `json:"authority_inn"`
    
}


type POAState string

const(
    POAStateCreated = "Created"
    POAStateSent = "Sent"
    POAStateReturned = "Returned"
    POAStateConfirmed = "Confirmed"
    POAStateRejected = "Rejected"
    
)


// SetStateCreated .
func (e *POA) SetStateCreated() error {
    e.State = POAStateCreated
    return nil
}

// SetStateSent .
func (e *POA) SetStateSent() error {if(e.State !=  POAStateCreated)||(e.State !=  POAStateReturned){
            return fmt.Errorf(" Order in state %s can not be set into 'Sent' (correct states: [Created Returned] )", e.State)
    }
    
    e.State = POAStateSent
    return nil
}

// SetStateReturned .
func (e *POA) SetStateReturned() error {if(e.State !=  POAStateSent){
            return fmt.Errorf(" Order in state %s can not be set into 'Returned' (correct states: [Sent] )", e.State)
    }
    
    e.State = POAStateReturned
    return nil
}

// SetStateConfirmed .
func (e *POA) SetStateConfirmed() error {if(e.State !=  POAStateSent){
            return fmt.Errorf(" Order in state %s can not be set into 'Confirmed' (correct states: [Sent] )", e.State)
    }
    
    e.State = POAStateConfirmed
    return nil
}

// SetStateRejected .
func (e *POA) SetStateRejected() error {if(e.State !=  POAStateSent){
            return fmt.Errorf(" Order in state %s can not be set into 'Rejected' (correct states: [Sent] )", e.State)
    }
    
    e.State = POAStateRejected
    return nil
}


