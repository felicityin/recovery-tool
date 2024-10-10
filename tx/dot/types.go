package dot

type GeneralResponse struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
	Error   string      `json:"error,omitempty"`
	Cause   string      `json:"cause,omitempty"`
	Stack   string      `json:"stack,omitempty"`
}

type AccountsBalanceInfoResponse struct {
	At          At            `json:"at"`
	Nonce       string        `json:"nonce"`
	TokenSymbol string        `json:"tokenSymbol"`
	Free        string        `json:"free"`
	Reserved    string        `json:"reserved"`
	MiscFrozen  string        `json:"miscFrozen"`
	FeeFrozen   string        `json:"feeFrozen"`
	Locks       []interface{} `json:"locks"`
}

type BlocksHeadResponse struct {
	Number string `json:"number"`
	Hash   string `json:"hash"`
}

type RuntimeSpecResponse struct {
	At                 At         `json:"at"`
	AuthoringVersion   string     `json:"authoringVersion"`
	TransactionVersion string     `json:"transactionVersion"`
	ImplVersion        string     `json:"implVersion"`
	SpecName           string     `json:"specName"`
	SpecVersion        string     `json:"specVersion"`
	ChainType          ChainType  `json:"chainType"`
	Properties         Properties `json:"properties"`
}

type TransactionResponse struct {
	Hash string `json:"hash"`
}

type ChainType struct {
	Live interface{} `json:"live"`
}

type Properties struct {
	Ss58Format    string   `json:"ss58Format"`
	TokenDecimals []string `json:"tokenDecimals"`
	TokenSymbol   []string `json:"tokenSymbol"`
}

type At struct {
	Hash   string `json:"hash"`
	Height string `json:"height"`
}
