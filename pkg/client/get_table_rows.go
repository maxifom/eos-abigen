package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"strconv"
)

type RowsRequest struct {
	Code          string `json:"code"`
	Table         string `json:"table"`
	Scope         string `json:"scope"`
	IndexPosition string `json:"index_position"`
	KeyType       string `json:"key_type"`
	UpperBound    string `json:"upper_bound"`
	LowerBound    string `json:"lower_bound"`
	JSON          bool   `json:"json"`
	Limit         int64  `json:"limit"`
	Reverse       bool   `json:"reverse"`
	ShowPayer     bool   `json:"show_payer"`
	logResponse   bool
}

type RequestOption func(*RowsRequest)

func Code(code string) RequestOption {
	return func(r *RowsRequest) {
		r.Code = code
	}
}

func Table(table string) RequestOption {
	return func(r *RowsRequest) {
		r.Table = table
	}
}

func Scope(scope string) RequestOption {
	return func(r *RowsRequest) {
		r.Scope = scope
	}
}

func IndexPosition(indexPos int64) RequestOption {
	return func(r *RowsRequest) {
		r.IndexPosition = strconv.FormatInt(indexPos, 10)
	}
}

func Limit(limit int64) RequestOption {
	return func(r *RowsRequest) {
		r.Limit = limit
	}
}

func KeyType(keyType string) RequestOption {
	return func(r *RowsRequest) {
		r.KeyType = keyType
	}
}

func Bounds(bound string) RequestOption {
	return func(r *RowsRequest) {
		r.UpperBound = bound
		r.LowerBound = bound
	}
}

func UpperBound(bound string) RequestOption {
	return func(r *RowsRequest) {
		r.UpperBound = bound
	}
}

func LowerBound(bound string) RequestOption {
	return func(r *RowsRequest) {
		r.LowerBound = bound
	}
}

func Reverse(reverse bool) RequestOption {
	return func(r *RowsRequest) {
		r.Reverse = reverse
	}
}

func ShowPayer(show bool) RequestOption {
	return func(r *RowsRequest) {
		r.ShowPayer = show
	}
}

func LogResponse(log bool) RequestOption {
	return func(r *RowsRequest) {
		r.logResponse = log
	}
}

func (c *RPCClient) GetTableRows(ctx context.Context, output interface{}, options ...RequestOption) error {
	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return fmt.Errorf("output expects pointer to rows resulting array")
	}

	request := &RowsRequest{}
	for _, o := range options {
		o(request)
	}

	u := c.cloneURL()
	u.Path = path.Join(u.Path, "v1", "chain", "get_table_rows")

	request.JSON = true
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}

	if request.logResponse {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(b))
		return json.NewDecoder(bytes.NewReader(b)).Decode(output)
	}

	return json.NewDecoder(resp.Body).Decode(output)
}
