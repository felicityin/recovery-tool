package ton

import (
	"context"
	"encoding/base64"
	"fmt"
	url2 "net/url"
	"strings"

	"recovery-tool/common"
)

func (this Client) Getblockchaininfo() (block *MasterBlockChainInfo, err error) {
	var chainInfo MasterBlockChainInfo
	url := fmt.Sprintf("%s/api/v3/masterchainInfo", this.endpoint)
	err = this.GetRequest(context.TODO(), url, &chainInfo)
	if err != nil {
		return nil, err
	}
	return &chainInfo, nil
}

func (this Client) Balance(addr string) (account *Account, err error) {
	var walletInfo WalletInfo
	addr = url2.QueryEscape(addr)
	url := fmt.Sprintf("%s/api/v3/wallet?address=%s", this.endpoint, addr)
	err = this.GetRequest(context.TODO(), url, &walletInfo)
	//fmt.Printf("url: %+v\n", url)
	//fmt.Printf("walletInfo:%+v\n", walletInfo)
	//fmt.Printf("error:%+v\n", err)
	if err != nil {
		return nil, err
	}
	isActive := false
	status := walletInfo.Status
	if strings.EqualFold(status, "active") {
		isActive = true
	}
	return &Account{
		Nonce:    walletInfo.Seqno,
		IsActive: isActive,
		Status:   strings.ToUpper(walletInfo.Status),
		Balance:  walletInfo.Balance,
	}, nil
}

func (this Client) JettonBalance(addr string, tokenAddr string, contractAddress string) (account *JettonWalletsRes, err error) {
	var walletInfo JettonWalletsRes
	ownerAddress := url2.QueryEscape(addr)
	address := url2.QueryEscape(tokenAddr)
	jettonAddress := url2.QueryEscape(contractAddress)
	//https://toncenter.com/api/v3/jetton/wallets?address=EQBIQoZJHCaRw-MRuDPDEkY1x33-g99Dvf4fQu-xXh4eNOEQ&owner_address=EQAW-1_rm44ppdD6qzcSSyZDAZH-KwldLeXmb2uTH6-WSkG0&jetton_address=EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs&limit=128&offset=0
	url := fmt.Sprintf("%s/api/v3/jetton/wallets?address=%s&owner_address=%s&jetton_address=%s", this.endpoint, address, ownerAddress, jettonAddress)
	err = this.GetRequest(context.TODO(), url, &walletInfo)
	if err != nil {
		return nil, err
	}
	if walletInfo.Error != "" {
		return nil, fmt.Errorf("%v", walletInfo.Error)
	}
	return &walletInfo, nil
}

func (this Client) Sequence(address string) (nonce uint64, err error) {
	account, err := this.Balance(address)
	if err != nil {
		return 0, err
	}
	return account.Nonce, nil
}

func (this Client) BroadCast(rawTX string) (txRes *SendTxResult, err error) {
	rawTXByte, err := base64.StdEncoding.DecodeString(rawTX)
	if err != nil {
		return nil, fmt.Errorf("DecodeString failed: %+v", err)
	}
	url := fmt.Sprintf("%s/api/v3/message", this.endpoint)
	var params SendTxParams
	var res SendTxResult
	params.Boc = string(rawTXByte)
	if err = this.PostRequest(context.TODO(), url, params, &res); err != nil {
		common.Logger.Errorf("post err: %s, res: %+v", err.Error(), res)
		return nil, err
	}
	return &res, nil
}
