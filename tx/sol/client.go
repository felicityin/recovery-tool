package sol

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"recovery-tool/common"
	"recovery-tool/common/code"
	"time"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

// Call will return body of response. if http code beyond 200~300, the error also returns.
func (c *Client) Call(ctx context.Context, params ...interface{}) ([]byte, error) {
	// prepare payload
	j, err := preparePayload(params)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare payload, err: %v", err)
	}

	// prepare request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(j))
	if err != nil {
		return nil, fmt.Errorf("failed to do http.NewRequestWithContext, err: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// do request
	httpclient := &http.Client{
		Timeout: 120 * time.Second,
	}

	res, err := httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request, err: %v", err)
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, err: %v", err)
	}

	// check response code
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return body, fmt.Errorf("get status code: %v", res.StatusCode)
	}

	return body, nil
}

type jsonRpcRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      uint64        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

func preparePayload(params []interface{}) ([]byte, error) {
	// prepare payload
	j, err := json.Marshal(jsonRpcRequest{
		JsonRpc: "2.0",
		Id:      1,
		Method:  params[0].(string),
		Params:  params[1:],
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (c *Client) processRpcCall(body []byte, rpcErr error, res interface{}) error {
	if rpcErr != nil {
		common.Logger.Errorf("rpc: call error, err: %v", rpcErr)
		return code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
	}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("rpc: failed to json decode body, err: %v", err)
	}
	return nil
}

func CheckRpcResult(res GeneralResponse, err error) error {
	if err != nil {
		return err
	}
	if res.Error != nil {
		errRes, err := json.Marshal(res.Error)
		if err != nil {
			return fmt.Errorf("rpc response error: %v", res.Error)
		}
		return fmt.Errorf("rpc response error: %v", string(errRes))
	}
	return nil
}
