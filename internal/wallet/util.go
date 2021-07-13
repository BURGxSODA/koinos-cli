package wallet

import (
	"encoding/base64"
	"fmt"

	types "github.com/koinos/koinos-types-golang"
	"github.com/shopspring/decimal"
	"github.com/ybbus/jsonrpc/v2"
)

// ContractStringToID converts a base64 contract id string to a contract id object
func ContractStringToID(s string) (*types.ContractIDType, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	cid := types.NewContractIDType()
	if err != nil {
		return cid, err
	}

	copy(cid[:], b)
	return cid, nil
}

// SatoshiToDecimal converts the given UInt64 value to a decimals with the given precision
func SatoshiToDecimal(balance int64, precision int) (*decimal.Decimal, error) {
	divisor, err := decimal.NewFromString(fmt.Sprintf("1e%d", precision))
	if err != nil {
		return nil, err
	}

	v := decimal.NewFromInt(balance).Div(divisor)
	return &v, nil
}

// KoinosRPCClient is a wrapper around the jsonrpc client
type KoinosRPCClient struct {
	client jsonrpc.RPCClient
}

// NewKoinosRPCClient creates a new koinos rpc client
func NewKoinosRPCClient(url string) *KoinosRPCClient {
	client := jsonrpc.NewClient(url)
	return &KoinosRPCClient{client: client}
}

// Call wraps the rpc client call and handles some of the boilerplate
func (c *KoinosRPCClient) Call(method string, params interface{}, returnType interface{}) error {
	// Make the rpc call
	resp, err := c.client.Call(method, params)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}

	// Fetch the contract response
	err = resp.GetObject(returnType)
	if err != nil {
		return err
	}

	return nil
}