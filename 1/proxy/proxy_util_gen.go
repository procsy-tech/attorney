package proxy

import (
	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// FcnArgsAsTransientMap .
func FcnArgsAsTransientMap(fcn string, args ...interface{}) (map[string][]byte, error) {
	rawArgs := []interface{}{fcn}

	for _, arg := range args {
		if _, ok := arg.(string); ok {
			rawArgs = append(rawArgs, arg)
			continue
		}
		bytes, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}
		rawArgs = append(rawArgs, string(bytes))
	}

	data, err := json.Marshal(rawArgs)
	if err != nil {
		return nil, err
	}

	return map[string][]byte{
		"Args": data,
	}, nil
}

// MakeChaincodeTransMapRequest creates chaincode request with fcn and args passed through transient map
func MakeChaincodeTransMapRequest(ccid string,
	hint []*fab.ChaincodeCall, fcn string, args ...interface{}) (channel.Request, error) {

	tmap, err := FcnArgsAsTransientMap(fcn, args...)
	if err != nil {
		return channel.Request{}, err
	}

	return channel.Request{
		ChaincodeID: ccid,
		Fcn:         "*",
		Args:        [][]byte{},

		TransientMap:    tmap,
		InvocationChain: hint,
	}, nil
}
