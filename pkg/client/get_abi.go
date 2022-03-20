package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type ABIResponse struct {
	AccountName string                 `json:"account_name"`
	ABI         map[string]interface{} `json:"abi"`
}

func (c *RPCClient) GetABI(ctx context.Context, accountName string) (*ABIResponse, error) {
	u := c.cloneURL()
	u.Path = path.Join(u.Path, "v1", "chain", "get_abi")
	reqBytes, err := json.Marshal(map[string]string{"account_name": accountName})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}

	var abiResp ABIResponse
	err = json.NewDecoder(resp.Body).Decode(&abiResp)
	if err != nil {
		return nil, err
	}

	return &abiResp, nil
}
