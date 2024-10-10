package common

import (
	"crypto/ecdsa"
	"fmt"

	hcchaincfg "github.com/HcashOrg/hcd/chaincfg"
	"github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/hcutil"
	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/portto/aptos-go-sdk/crypto"
	"github.com/portto/aptos-go-sdk/models"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"golang.org/x/crypto/blake2b"
)

const (
	SignKindEcdsa   = 0 //ecdsa
	SignKindEddsa   = 1 //eddsa 只有keygen和sign
	SignKindSchnorr = 2 //schnorr
)

// zero is deafult of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// coin type from bip44
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
	Dot       = Zero + 354
	SUI       = Zero + 784
	SOL       = Zero + 501
	DOT       = Zero + 354
	TON       = Zero + 607
)

type Option struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

// chain name
const (
	BitcoinChain     = "Bitcoin"
	EthereumChain    = "Ethereum"
	TronChain        = "Tron"
	BSCChain         = "BSC"
	BitcoinCashChain = "Bitcoin Cash"
	DogeChain        = "Doge"
	LitecoinChain    = "Litecoin"
	HecoChain        = "Heco"
	PolygonChain     = "Polygon"
	ArbitrumChain    = "Arbitrum"
	PolkadotChain    = "Polkadot"
	AptostChain      = "Aptos"
	SolanaChain      = "Solana"
	BaseChain        = "Base Chain"
	TonChain         = "TON"
)

var ChainList = []*Option{
	{
		Name: BitcoinChain,
		Val:  BitcoinChain,
	},
	{
		Name: EthereumChain,
		Val:  EthereumChain,
	},
	{
		Name: TronChain,
		Val:  TronChain,
	},
	{
		Name: BSCChain,
		Val:  BSCChain,
	},
	{
		Name: BitcoinCashChain,
		Val:  BitcoinCashChain,
	},
	{
		Name: DogeChain,
		Val:  DogeChain,
	},
	{
		Name: LitecoinChain,
		Val:  LitecoinChain,
	},
	{
		Name: HecoChain,
		Val:  HecoChain,
	},
	{
		Name: PolygonChain,
		Val:  PolygonChain,
	},
	{
		Name: ArbitrumChain,
		Val:  ArbitrumChain,
	},
	{
		Name: PolkadotChain,
		Val:  PolkadotChain,
	},
	{
		Name: AptostChain,
		Val:  AptostChain,
	},
	{
		Name: SolanaChain,
		Val:  SolanaChain,
	},
	{
		Name: BaseChain,
		Val:  BaseChain,
	},
	{
		Name: TonChain,
		Val:  TonChain,
	},
}

type CoinInfo struct {
	SignKind int    // 币种对应签名算法
	CoinType uint32 // 币种使用的coin type
}

var ChainInfos = map[string]CoinInfo{
	BitcoinChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: BTC,
	},
	LitecoinChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: LTC,
	},
	DogeChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: DOGE,
	},
	EthereumChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: ETH,
	},
	BitcoinCashChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: BCH,
	},
	BSCChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: ETH,
	},
	HecoChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: ETH,
	},
	TronChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: TRX,
	},
	PolygonChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: ETH,
	},
	ArbitrumChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: ETH,
	},
	BaseChain: CoinInfo{
		SignKind: SignKindEcdsa,
		CoinType: ETH,
	},

	// eddsa
	AptostChain: CoinInfo{
		SignKind: SignKindEddsa,
		CoinType: Apt,
	},
	SolanaChain: CoinInfo{
		SignKind: SignKindEddsa,
		CoinType: SOL,
	},
	PolkadotChain: CoinInfo{
		SignKind: SignKindEddsa,
		CoinType: DOT,
	},
	TonChain: CoinInfo{
		SignKind: SignKindEddsa,
		CoinType: TON,
	},
}

var ss58Prefix = []byte("SS58PRE")
var DOTNetWorkByteMap = map[string]byte{
	"DOT": 0x00,
	"KSM": 0x00,
}

func SwitchCoin(coinType uint32) string {
	var chain string

	switch coinType + Zero {
	case BTC:
		chain = "btc"
	case LTC:
		chain = "ltc"
	case DOGE:
		chain = "doge"
	case ETH:
		chain = "eth"
	case BCH:
		chain = "bch"
	//case DASH:
	//	chain = "dash"
	case TRX:
		chain = "trx"
	case HECO:
		chain = "eth"
	case BSC:
		chain = "eth"
	case POLYGON:
		chain = "eth"
	case ARBITRUM:
		chain = "eth"
	case SOL:
		chain = "sol"
	case Apt:
		chain = "apt"
	case Dot:
		chain = "dot"
	case TON:
		chain = "ton"
	default:
		panic("invalid chain type")
	}

	return chain
}

func SwitchEcdsaChainAddress(ecdsaPk *ecdsa.PublicKey, chain string) (string, error) {
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
	case "dash":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &DASHParams)
		if err != nil {
			return "", err
		}
	case "dcr":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &DCRParams)
		if err != nil {
			return "", err
		}
	case "rvn":
		var err error
		addressStr, err = makeBtcAddress(ecdsaPk, &RVNParams)
		if err != nil {
			return "", err
		}
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

func SwitchEddsaChainAddress(publicKey *edwards.PublicKey, chain string) (addressStr string, err error) {
	switch chain {
	case "sol":
		addressStr = base58.Encode(publicKey.Serialize())
	case "apt":
		accountAddress := models.AccountAddress{}
		accountAddress = crypto.SingleSignerAuthKey(publicKey.Serialize())
		addressStr = accountAddress.PrefixZeroTrimmedHex()
	case "dot":
		addressStr, err = DOTPublicKeyToAddress(publicKey.Serialize(), DOTNetWorkByteMap["DOT"])
		if err != nil {
			return "", err
		}
	case "ton":
		address, err := wallet.AddressFromPubKey(publicKey.Serialize(), wallet.V3, wallet.DefaultSubwallet)
		if err != nil {
			return "", err
		}
		return address.String(), nil
	default:
		return "", fmt.Errorf("eddsa unsupport chain type: %s", chain)
	}
	return addressStr, nil
}

func DOTPublicKeyToAddress(pub []byte, network byte) (string, error) {
	enc := append([]byte{network}, pub...)
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return "", err
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return "", err
	}
	checksum := hasher.Sum(nil)

	return base58.Encode(append(enc, checksum[:2]...)), nil
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
