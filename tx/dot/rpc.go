package dot

import (
	"context"
	"fmt"
	"recovery-tool/common"
)

func (c *Client) AccountsBalanceInfo(address string) (AccountsBalanceInfoResponse, error) {
	res := struct {
		GeneralResponse
		AccountsBalanceInfoResponse
	}{}
	err := c.Get(context.Background(), fmt.Sprintf("%v/accounts/%v/balance-info", c.endpoint, address), nil, &res)
	err = CheckRpcResult(res.GeneralResponse, err)

	if err != nil {
		return AccountsBalanceInfoResponse{}, err
	}
	return res.AccountsBalanceInfoResponse, nil
}

func (c *Client) BlocksHead() (BlocksHeadResponse, error) {
	res := struct {
		GeneralResponse
		BlocksHeadResponse
	}{}
	err := c.Get(context.Background(), fmt.Sprintf("%v/blocks/head", c.endpoint), nil, &res)
	err = CheckRpcResult(res.GeneralResponse, err)

	if err != nil {
		return BlocksHeadResponse{}, err
	}
	return res.BlocksHeadResponse, nil
}

func (c *Client) RuntimeSpec() (RuntimeSpecResponse, error) {
	res := struct {
		GeneralResponse
		RuntimeSpecResponse
	}{}
	err := c.Get(context.Background(), fmt.Sprintf("%v/runtime/spec", c.endpoint), nil, &res)
	err = CheckRpcResult(res.GeneralResponse, err)

	if err != nil {
		return RuntimeSpecResponse{}, err
	}
	return res.RuntimeSpecResponse, nil
}

func (c *Client) Transaction(raw string) (TransactionResponse, error) {
	res := struct {
		GeneralResponse
		TransactionResponse
	}{}
	err := c.Request(context.Background(), "transaction",
		map[string]string{
			"tx": raw,
		}, &res)
	err = CheckRpcResult(res.GeneralResponse, err)
	if err != nil {
		common.Logger.Errorf("send tx err: %s", err.Error())
		return TransactionResponse{}, err
	}
	return res.TransactionResponse, nil
}
