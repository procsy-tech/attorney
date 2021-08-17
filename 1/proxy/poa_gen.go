package proxy

import(
	"github.com/procsy-tech/attorney/dto"
	"github.com/procsy-tech/attorney/entity"
    "github.com/procsy-tech/attorney/api"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"encoding/json"
	"fmt"
	"errors"
)

type POAService struct {
	channelClient   *channel.Client
}


func (svc *POAService) Create(POA *entity.POA) (int64, error){
	ccRequest,err := MakeChaincodeTransMapRequest("attorney", []*fab.ChaincodeCall{
			{ID: "attorney"},
		}, api.Create, POA)
	if err != nil{
		return 0,  fmt.Errorf("error creating ccRequest: %s", err)
	}
	var ccResponse channel.Response
	
		ccResponse, err = svc.channelClient.Execute(ccRequest, channel.WithRetry(retry.DefaultChannelOpts))
		if err != nil {
			return 0,  fmt.Errorf("failed to execute: %s", err)
		}
	

	if ccResponse.ChaincodeStatus != 200 {
		return 0,  errors.New(string(ccResponse.Payload))
	}

	var response dto.CreateResponse
	err = json.Unmarshal(ccResponse.Payload, &response)
	if err != nil {
		return 0,  fmt.Errorf("failed to parse response payload: %s", err)
	}

	if len(response.Error) != 0{
		return 0,  errors.New(response.Error)
	}

	
    	return response.Result, nil
	}

func (svc *POAService) ConfirmAttorney(ID string) error{
	ccRequest,err := MakeChaincodeTransMapRequest("attorney", []*fab.ChaincodeCall{
			{ID: "attorney"},
		}, api.ConfirmAttorney, ID)
	if err != nil{
		return  fmt.Errorf("error creating ccRequest: %s", err)
	}
	var ccResponse channel.Response
	
		ccResponse, err = svc.channelClient.Execute(ccRequest, channel.WithRetry(retry.DefaultChannelOpts))
		if err != nil {
			return  fmt.Errorf("failed to execute: %s", err)
		}
	

	if ccResponse.ChaincodeStatus != 200 {
		return  errors.New(string(ccResponse.Payload))
	}

	var response dto.ConfirmAttorneyResponse
	err = json.Unmarshal(ccResponse.Payload, &response)
	if err != nil {
		return  fmt.Errorf("failed to parse response payload: %s", err)
	}

	if len(response.Error) != 0{
		return  errors.New(response.Error)
	}

	
	return nil
    }


func NewPOAService(
	chanProv context.ChannelProvider,
) (*POAService, error) {
	channelClient, err := channel.New(chanProv)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel client: %s", err)
	}
	return &POAService{
		channelClient: channelClient,
	}, nil
}
