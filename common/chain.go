package common

import (
	"crypto/ecdsa"
	"fmt"

	hcchaincfg "github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/hcutil"
	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
)

// zero is deafult of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet type from bip44
const (
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types
	BTC       = Zero + 0
	LTC       = Zero + 2
	DOGE      = Zero + 3
	DASH      = Zero + 5
	Optimism  = Zero + 10
	ETH       = Zero + 60
	BCH       = Zero + 145
	TRX       = Zero + 195
	BSV       = Zero + 236
	Fantom    = Zero + 250
	ZKSYNC    = Zero + 324
	POLYGON   = Zero + 966
	ARBITRUM  = Zero + 9001
	OKChain   = Zero + 996
	BSC       = Zero + 714
	HECO      = Zero + 553
	Avalanche = Zero + 43114
	Apt       = Zero + 637
	SUI       = Zero + 784
	SOL       = Zero + 501
)

func SwitchChainAddress(ecdsaPk *ecdsa.PublicKey, chain string) (string, error) {
	return switchChainAddress(ecdsaPk, chain)
}

func switchChainAddress(ecdsaPk *ecdsa.PublicKey, chain string) (string, error) {
	var addressStr string
	switch chain {
	case "eth":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "bnb_bsc":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "ht_heco":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "trx":
		a := addr.PubkeyToAddress(*ecdsaPk)
		addressStr = a.String()
	case "btc":
		var xFieldVal btcec.FieldVal
		var yFieldVal btcec.FieldVal
		if overflow := xFieldVal.SetByteSlice(ecdsaPk.X.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.X.Bytes()) overflow: %x", ecdsaPk.X.Bytes())
			panic(err)
		}
		if overflow := yFieldVal.SetByteSlice(ecdsaPk.Y.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.Y.Bytes()) overflow: %x", ecdsaPk.Y.Bytes())
			panic(err)
		}
		btcecPubkey := btcec.NewPublicKey(&xFieldVal, &yFieldVal)
		param := &chaincfg.MainNetParams
		pkHash, err := btcutil.NewAddressPubKey(btcecPubkey.SerializeCompressed(), param)
		if err != nil {
			return "", err
		}
		addressStr = pkHash.EncodeAddress()
	case "btc_test":
		var xFieldVal btcec.FieldVal
		var yFieldVal btcec.FieldVal
		if overflow := xFieldVal.SetByteSlice(ecdsaPk.X.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.X.Bytes()) overflow: %x", ecdsaPk.X.Bytes())
			panic(err)
		}
		if overflow := yFieldVal.SetByteSlice(ecdsaPk.Y.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.Y.Bytes()) overflow: %x", ecdsaPk.Y.Bytes())
			panic(err)
		}
		btcecPubkey := btcec.NewPublicKey(&xFieldVal, &yFieldVal)
		pkHash, err := btcutil.NewAddressPubKey(btcecPubkey.SerializeCompressed(), &chaincfg.TestNet3Params)
		if err != nil {
			return "", err
		}
		addressStr = pkHash.EncodeAddress()
	case "ltc":
		var xFieldVal btcec.FieldVal
		var yFieldVal btcec.FieldVal
		if overflow := xFieldVal.SetByteSlice(ecdsaPk.X.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.X.Bytes()) overflow: %x", ecdsaPk.X.Bytes())
			panic(err)
		}
		if overflow := yFieldVal.SetByteSlice(ecdsaPk.Y.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.Y.Bytes()) overflow: %x", ecdsaPk.Y.Bytes())
			panic(err)
		}
		btcecPubkey := btcec.NewPublicKey(&xFieldVal, &yFieldVal)
		pkHash, err := btcutil.NewAddressPubKey(btcecPubkey.SerializeCompressed(), &LTCParams)
		if err != nil {
			return "", err
		}
		addressStr = pkHash.EncodeAddress()
		fmt.Printf("switchChainAddress, LTCParams.PubKeyHashAddrID %d \n", LTCParams.PubKeyHashAddrID)
	case "doge":
		var xFieldVal btcec.FieldVal
		var yFieldVal btcec.FieldVal
		if overflow := xFieldVal.SetByteSlice(ecdsaPk.X.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.X.Bytes()) overflow: %x", ecdsaPk.X.Bytes())
			panic(err)
		}
		if overflow := yFieldVal.SetByteSlice(ecdsaPk.Y.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.Y.Bytes()) overflow: %x", ecdsaPk.Y.Bytes())
			panic(err)
		}
		btcecPubkey := btcec.NewPublicKey(&xFieldVal, &yFieldVal)
		pkHash, err := btcutil.NewAddressPubKey(btcecPubkey.SerializeCompressed(), &DOGEParams)
		if err != nil {
			return "", err
		}
		addressStr = pkHash.EncodeAddress()
		fmt.Printf("switchChainAddress, DOGEParams.PubKeyHashAddrID %d \n", DOGEParams.PubKeyHashAddrID)
	case "usdt":
		var xFieldVal btcec.FieldVal
		var yFieldVal btcec.FieldVal
		if overflow := xFieldVal.SetByteSlice(ecdsaPk.X.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.X.Bytes()) overflow: %x", ecdsaPk.X.Bytes())
			panic(err)
		}
		if overflow := yFieldVal.SetByteSlice(ecdsaPk.Y.Bytes()); overflow {
			err := fmt.Errorf("xFieldVal.SetByteSlice(pk.Y.Bytes()) overflow: %x", ecdsaPk.Y.Bytes())
			panic(err)
		}
		btcecPubkey := btcec.NewPublicKey(&xFieldVal, &yFieldVal)
		param := &chaincfg.MainNetParams
		pkHash, err := btcutil.NewAddressPubKey(btcecPubkey.SerializeCompressed(), param)
		if err != nil {
			return "", err
		}
		addressStr = pkHash.EncodeAddress()
	case "hc":
		pubKey := ecrypto.CompressPubkey(ecdsaPk)
		pubKeyHash := hcutil.Hash160(pubKey)
		param := &hcchaincfg.MainNetParams
		addr, err := hcutil.NewAddressPubKeyHash(pubKeyHash,
			param, chainec.ECTypeSecp256k1)
		if err != nil {
			return "", err
		}
		addressStr = addr.EncodeAddress()
	case "bch":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &BCHParams)
		if err != nil {
			return "", err
		}
		fmt.Printf("switchChainAddress, BCHParams.PubKeyHashAddrID %d \n", BCHParams.PubKeyHashAddrID)
	case "dash":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &DASHParams)
		if err != nil {
			return "", err
		}
		fmt.Printf("switchChainAddress, DASHParams.PubKeyHashAddrID %d \n", DASHParams.PubKeyHashAddrID)
	case "dcr":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &DCRParams)
		if err != nil {
			return "", err
		}
		fmt.Printf("switchChainAddress, DCRParams.PubKeyHashAddrID %d \n", DCRParams.PubKeyHashAddrID)
	case "rvn":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &RVNParams)
		if err != nil {
			return "", err
		}
		fmt.Printf("switchChainAddress, RVNParams.PubKeyHashAddrID %d \n", RVNParams.PubKeyHashAddrID)
	case "okt":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "cmp":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "ftm":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "smartbch":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "eth_aurora":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "wemix":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "gdcc":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "eth_zksync":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "ethg":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "core":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "mbe":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "ethw":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "rei":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "eth_arbitrum":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "eth_optimism":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "movr":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "avax_c":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	case "matic_polygon":
		address := ecrypto.PubkeyToAddress(*ecdsaPk)
		addressStr = address.Hex()
	default:
		return "", fmt.Errorf("ecdsa, unsupport chain type for %s", chain)
	}
	return addressStr, nil
}

func makeBtcAddress(ecdsaPk *ecdsa.PublicKey, params *chaincfg.Params) (addressStr string, err error) {
	var xFieldVal btcec.FieldVal
	var yFieldVal btcec.FieldVal
	if overflow := xFieldVal.SetByteSlice(ecdsaPk.X.Bytes()); overflow {
		err := fmt.Errorf("xFieldVal.SetByteSlice(pk.X.Bytes()) overflow: %x", ecdsaPk.X.Bytes())
		panic(err)
	}
	if overflow := yFieldVal.SetByteSlice(ecdsaPk.Y.Bytes()); overflow {
		err := fmt.Errorf("xFieldVal.SetByteSlice(pk.Y.Bytes()) overflow: %x", ecdsaPk.Y.Bytes())
		panic(err)
	}
	btcecPubkey := btcec.NewPublicKey(&xFieldVal, &yFieldVal)
	pkHash, err := btcutil.NewAddressPubKey(btcecPubkey.SerializeCompressed(), params)
	if err != nil {
		return "", err
	}
	addressStr = pkHash.EncodeAddress()
	return addressStr, nil
}
