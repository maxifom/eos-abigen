package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/maxifom/eos-abigen-go/pkg/base"
)

type AccountResponse struct {
	AccountName            string         `json:"account_name"`
	HeadBlockNum           int64          `json:"head_block_num"`
	HeadBlockTime          string         `json:"head_block_time"`
	Privileged             bool           `json:"privileged"`
	LastCodeUpdate         string         `json:"last_code_update"`
	Created                string         `json:"created"`
	RAMQuota               int64          `json:"ram_quota"`
	NetWeight              *base.BigInt   `json:"net_weight"`
	CPUWeight              *base.BigInt   `json:"cpu_weight"`
	NetLimit               AccountLimit   `json:"net_limit"`
	CPULimit               AccountLimit   `json:"cpu_limit"`
	RAMUsage               int64          `json:"ram_usage"`
	Permissions            []Permission   `json:"permissions"`
	TotalResources         TotalResources `json:"total_resources"`
	SelfDelegatedBandwidth interface{}    `json:"self_delegated_bandwidth"`
	RefundRequest          interface{}    `json:"refund_request"`
	VoterInfo              interface{}    `json:"voter_info"`
	RexInfo                interface{}    `json:"rex_info"`
}

type AccountLimit struct {
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
	Max       int64 `json:"max"`
}

type Permission struct {
	PermName     string       `json:"perm_name"`
	Parent       string       `json:"parent"`
	RequiredAuth RequiredAuth `json:"required_auth"`
}

type RequiredAuth struct {
	Threshold int64         `json:"threshold"`
	Keys      []Key         `json:"keys"`
	Accounts  []interface{} `json:"accounts"`
	Waits     []interface{} `json:"waits"`
}

type Key struct {
	Key    string `json:"key"`
	Weight int64  `json:"weight"`
}

type TotalResources struct {
	Owner     string `json:"owner"`
	NetWeight string `json:"net_weight"`
	CPUWeight string `json:"cpu_weight"`
	RAMBytes  int64  `json:"ram_bytes"`
}

func (c *RPCClient) GetAccount(ctx context.Context, accountName string) (*AccountResponse, error) {
	u := c.cloneURL()
	u.Path = path.Join(u.Path, "v1", "chain", "get_account")
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

	var accResp AccountResponse
	err = json.NewDecoder(resp.Body).Decode(&accResp)
	if err != nil {
		return nil, err
	}

	return &accResp, nil
}
