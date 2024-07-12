package rpc

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// A Client is a Tron RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type DOTClient struct {
	BaseURL string
	// AccessToken string
	Debug  bool
	client *req.Req
}

// NewClient create new client to connect
func newClient(url, token string, debug bool) *DOTClient {
	c := DOTClient{
		BaseURL: url,
		// AccessToken: token,
		Debug: debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *DOTClient) Call(path string, param interface{}) (*gjson.Result, error) {

	if c == nil || c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	url := c.BaseURL + path
	//authHeader := req.Header{"Accept": "application/json"}
	fmt.Println("url = ", url)
	r, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	if c.Debug {

	}

	if r.Response().StatusCode != http.StatusOK {
		message := gjson.ParseBytes(r.Bytes()).String()
		message = fmt.Sprintf("[%s]%s", r.Response().Status, message)
		return nil, errors.New(message)
	}

	res := gjson.ParseBytes(r.Bytes())
	return &res, nil
}

type DOTBlock struct {
	/*
			{"account_number":51,"address":"DOT1jxfh2g85q3v0tdq56fnevx6xcxtcnhtsmcu64m",
			"balances":[]
			,"flags":0,"public_key":[3,86,224,165,128,56,154,111,210,204,
			145,205,82,92,109,90,77,128,84,175,112,223,23,72,78,88,103,143,159,87,74,11,77],
			"sequence":348312
		}
	*/
	AccountNumber int64 // 这里采用 BlockID
	Sequence      int64
	// Merkleroot        string
	//Confirmations     uint64
}

func (block *DOTBlock) GetAccountNumber() int64 {
	return block.AccountNumber
}
func (block *DOTBlock) GetSequence() int64 {
	return block.Sequence
}

// GetNowBlock Done!
// Function：Query the latest block
// 	demo: curl -X POST http://127.0.0.1:8090/wallet/getnowblock
// Parameters：None
// Return value：Latest block on full node
func GetDOTNowBlock(address string) (block *DOTBlock, err error) {
	client := newClient("wss://rpc.polkadot.io", "", true)
	path := "/block/"
	//path += address
	r, err := client.Call(path, nil)

	if err != nil {
		return nil, err
	}

	block = newBlock(r, false)

	return block, nil
}

func newBlock(json *gjson.Result, isTestnet bool) *DOTBlock {

	//header := gjson.Get(json.Raw, "block_header").Get("raw_data")
	// 解析json
	b := &DOTBlock{}
	fmt.Println("json.Raw = ", json.Raw)
	b.AccountNumber = gjson.Get(json.Raw, "account_number").Int()
	b.Sequence = gjson.Get(json.Raw, "sequence").Int()
	return b
}
