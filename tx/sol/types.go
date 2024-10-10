package sol

import (
	"encoding/json"
)

type Commitment string

const (
	CommitmentFinalized Commitment = "finalized"
	CommitmentConfirmed Commitment = "confirmed"
	CommitmentProcessed Commitment = "processed"
)

// GetAccountInfoConfigEncoding is account's data encode format
type ConfigEncoding string

const (
	// GetAccountInfoConfigEncodingBase58 limited to Account data of less than 128 bytes
	ConfigEncodingBase58     ConfigEncoding = "base58"
	ConfigEncodingJsonParsed ConfigEncoding = "jsonParsed"
	ConfigEncodingBase64     ConfigEncoding = "base64"
	ConfigEncodingBase64Zstd ConfigEncoding = "base64+zstd"
)

type GeneralResponse struct {
	JsonRPC string         `json:"jsonrpc"`
	ID      uint64         `json:"id"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// GetAccountInfoConfig is an option config for `getAccountInfo`
type Cfg struct {
	Commitment Commitment     `json:"commitment,omitempty"`
	Encoding   ConfigEncoding `json:"encoding,omitempty"`
	DataSlice  *DataSlice     `json:"dataSlice,omitempty"`
}

// GetAccountInfoConfigDataSlice is a part of GetAccountInfoConfig
type DataSlice struct {
	Offset uint64 `json:"offset,omitempty"`
	Length uint64 `json:"length,omitempty"`
}

// GetBalanceResponse is a full raw rpc response of `getBalance`
type GetBalanceResponse struct {
	GeneralResponse
	Result GetBalanceResult `json:"result"`
}

// GetBalanceResult is a part of raw rpc response of `getBalance`
type GetBalanceResult struct {
	Context Context `json:"context"`
	Value   uint64  `json:"value"`
}

type Context struct {
	Slot uint64 `json:"slot"`
}

// GetSlotResponse is a full raw rpc response of `getSlot`
type GetSlotResponse struct {
	GeneralResponse
	Result uint64 `json:"result"`
}

// GetSlotResponse is a full raw rpc response of `getSlot`
type GetBlockHeightResponse struct {
	GeneralResponse
	Result uint64 `json:"result"`
}

// GetRecentBlockHashResponse is full raw response of `getRecentBlockhash`
type GetRecentBlockHashResponse struct {
	GeneralResponse
	Result GetRecentBlockHashResult `json:"result"`
}

// GetLatestBlockhash is full raw response of `getRecentBlockhash`
type GetLatestBlockHashResponse struct {
	GeneralResponse
	Result GetLatestBlockHashResult `json:"result"`
}

// GetLatestBlockhashResult is part of response of `getRecentBlockhash`
type GetLatestBlockHashResult struct {
	Context Context                       `json:"context"`
	Value   GetLatestBlockHashResultValue `json:"value"`
}

// GetRecentBlockHashResult is part of response of `getRecentBlockhash`
type GetRecentBlockHashResult struct {
	Context Context                       `json:"context"`
	Value   GetRecentBlockHashResultValue `json:"value"`
}

// GetRecentBlockHashResultValue is part of response of `getRecentBlockhash`
type GetRecentBlockHashResultValue struct {
	Blockhash     string        `json:"blockhash"`
	FeeCalculator FeeCalculator `json:"feeCalculator"`
}

// GetRecentBlockHashResultValue is part of response of `getRecentBlockhash`
type GetLatestBlockHashResultValue struct {
	Blockhash            string `json:"blockhash"`
	LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
}

// FeeCalculator is a list of fee
type FeeCalculator struct {
	LamportsPerSignature uint64 `json:"lamportsPerSignature"`
}

// GetAccountInfoResponse is a full raw rpc response of `getAccountInfo`
type GetAccountInfoResponse struct {
	GeneralResponse
	Result GetAccountInfoResult `json:"result"`
}

// GetAccountInfoResult is rpc result of `getAccountInfo`
type GetAccountInfoResult struct {
	Context Context                   `json:"context"`
	Value   GetAccountInfoResultValue `json:"value"`
}

// GetAccountInfoResultValue is rpc result of `getAccountInfo`
type GetAccountInfoResultValue struct {
	Lamports  uint64          `json:"lamports"`
	Owner     string          `json:"owner"`
	Excutable bool            `json:"excutable"`
	RentEpoch interface{}     `json:"rentEpoch"`
	Data      json.RawMessage `json:"data"`
}

type GetFeesResponse struct {
	GeneralResponse
	Result struct {
		Context `json:"context"`
		Value   struct {
			Blockhash     string `json:"blockhash"`
			FeeCalculator struct {
				LamportsPerSignature uint64 `json:"lamportsPerSignature"`
			} `json:"feeCalculator"`
			LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
			LastValidSlot        uint64 `json:"lastValidSlot"`
		} `json:"value"`
	} `json:"result"`
}

type SendTransactionResponse struct {
	GeneralResponse
	Result string `json:"result"`
}

type GetTokenAccountsByOwnerResponse struct {
	GeneralResponse
	Result struct {
		Context Context                              `json:"context"`
		Value   GetTokenAccountsByOwnerResponseValue `json:"value"`
	} `json:"result"`
}

type GetTokenAccountsByOwnerResponseValue []struct {
	Account struct {
		Data struct {
			Parsed struct {
				Info struct {
					IsNative    bool   `json:"isNative"`
					Mint        string `json:"mint"`
					Owner       string `json:"owner"`
					State       string `json:"state"`
					TokenAmount struct {
						Amount         string  `json:"amount"`
						Decimals       int     `json:"decimals"`
						UIAmount       float64 `json:"uiAmount"`
						UIAmountString string  `json:"uiAmountString"`
					} `json:"tokenAmount"`
				} `json:"info"`
				Type string `json:"type"`
			} `json:"parsed"`
			Program string `json:"program"`
			Space   int    `json:"space"`
		} `json:"data"`
		Executable bool        `json:"executable"`
		Lamports   int         `json:"lamports"`
		Owner      string      `json:"owner"`
		RentEpoch  interface{} `json:"rentEpoch"`
	} `json:"account"`
	Pubkey string `json:"pubkey"`
}

type GetTokenSupplyResponse struct {
	GeneralResponse
	Result struct {
		Context struct {
			APIVersion string `json:"apiVersion"`
			Slot       int    `json:"slot"`
		} `json:"context"`
		Value struct {
			Amount         string  `json:"amount"`
			Decimals       int     `json:"decimals"`
			UIAmount       float64 `json:"uiAmount"`
			UIAmountString string  `json:"uiAmountString"`
		} `json:"value"`
	} `json:"result"`
}
