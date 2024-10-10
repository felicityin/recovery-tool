package sol

import "context"

func (c *Client) GetBalance(ctx context.Context, base58Addr string) (uint64, error) {
	var res GetBalanceResponse
	body, rpcErr := c.Call(ctx, "getBalance", base58Addr)
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return 0, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result.Value, nil
}

// GetSlot returns the SOL balance
func (c *Client) GetSlot(ctx context.Context) (GetSlotResponse, error) {
	var res GetSlotResponse
	body, rpcErr := c.Call(ctx, "getSlot", map[string]string{"commitment": string(CommitmentFinalized)})
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetSlotResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetSlotResponse{}, err
	}
	return res, nil
}

func (c *Client) GetBlockHeight(ctx context.Context) (GetBlockHeightResponse, error) {
	var res GetBlockHeightResponse
	body, rpcErr := c.Call(ctx, "getBlockHeight", map[string]string{"commitment": string(CommitmentFinalized)})
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetBlockHeightResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetBlockHeightResponse{}, err
	}
	return res, nil
}

func (c *Client) GetRecentBlockhash(ctx context.Context) (GetRecentBlockHashResponse, error) {
	var res GetRecentBlockHashResponse
	body, rpcErr := c.Call(ctx, "getRecentBlockhash", map[string]string{"commitment": string(CommitmentFinalized)})
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetRecentBlockHashResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetRecentBlockHashResponse{}, err
	}
	return res, nil
}

func (c *Client) GetLatestBlockHash(ctx context.Context) (GetLatestBlockHashResponse, error) {
	var res GetLatestBlockHashResponse
	body, rpcErr := c.Call(ctx, "getLatestBlockhash", map[string]string{"commitment": string(CommitmentFinalized)})
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetLatestBlockHashResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetLatestBlockHashResponse{}, err
	}
	return res, nil
}

// GetAccountInfo returns all information associated with the account of provided Pubkey
func (c *Client) GetAccountInfoWithCfg(ctx context.Context, base58Addr string, cfg Cfg) (GetAccountInfoResponse, error) {
	var res GetAccountInfoResponse
	body, rpcErr := c.Call(ctx, "getAccountInfo", base58Addr, cfg)
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetAccountInfoResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetAccountInfoResponse{}, err
	}
	return res, nil
}

func (c *Client) GetFeesWithCfg(ctx context.Context, cfg Cfg) (GetFeesResponse, error) {
	var res GetFeesResponse
	body, rpcErr := c.Call(ctx, "getFees", cfg)
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetFeesResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetFeesResponse{}, err
	}
	return res, nil
}

func (c *Client) SendTransaction(ctx context.Context, raw string) (SendTransactionResponse, error) {
	var res SendTransactionResponse
	body, rpcErr := c.Call(ctx, "sendTransaction", raw, map[string]interface{}{"skipPreflight": false})
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return SendTransactionResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return SendTransactionResponse{}, err
	}
	return res, nil
}

func (c *Client) GetTokenAccountsByOwnerWithCfg(ctx context.Context, address, mint string, cfg Cfg) (GetTokenAccountsByOwnerResponse, error) {
	var res GetTokenAccountsByOwnerResponse
	body, rpcErr := c.Call(ctx, "getTokenAccountsByOwner", address, map[string]string{"mint": mint}, cfg)
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetTokenAccountsByOwnerResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetTokenAccountsByOwnerResponse{}, err
	}
	return res, nil
}

func (c *Client) GetTokenSupply(ctx context.Context, contractAddress string) (GetTokenSupplyResponse, error) {
	var res GetTokenSupplyResponse
	body, rpcErr := c.Call(ctx, "getTokenSupply", contractAddress)
	err := c.processRpcCall(body, rpcErr, &res)
	if err != nil {
		return GetTokenSupplyResponse{}, err
	}
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetTokenSupplyResponse{}, err
	}
	return res, nil
}
