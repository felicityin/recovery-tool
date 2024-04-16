package common

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"

	"recovery-tool/crypto"
	"recovery-tool/crypto/ckd"
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

type RecoveryData struct {
	HbcPrivKeys   []string
	HbcChainCodes []string
	UserPubKey    string
}

type RootKeys struct {
	HbcShare0 *RootKey
	HbcShare1 *RootKey
	UsrShare  *RootKey
}

type RootKey struct {
	PrivKey   *big.Int
	PubKey    *crypto.ECPoint
	ChainCode []byte
}

func DeriveChild(params *RootKeys, hdPath string) (*big.Int, string, error) {
	privKey, err := DerivePrivKey(params, hdPath)
	if err != nil {
		return nil, "", err
	}

	address, err := DeriveAddress(privKey, hdPath)
	if err != nil {
		return nil, "", err
	}

	return privKey, address, nil
}

func DerivePrivKey(params *RootKeys, hdPath string) (*big.Int, error) {
	hbcPrivKey0, err := deriveChildPrivKeys(params.HbcShare0, hdPath)
	if err != nil {
		return nil, err
	}

	hbcPrivKey1, err := deriveChildPrivKeys(params.HbcShare1, hdPath)
	if err != nil {
		return nil, err
	}

	userPrivKey, err := deriveChildPrivKeys(params.UsrShare, hdPath)
	if err != nil {
		return nil, err
	}

	privateKey := hbcPrivKey0

	privateKey.Add(privateKey, hbcPrivKey1)
	privateKey.Mod(privateKey, btcec.S256().Params().N)

	privateKey.Add(privateKey, userPrivKey)
	privateKey.Mod(privateKey, btcec.S256().Params().N)

	return privateKey, nil
}

func DeriveAddress(privKey *big.Int, hdPath string) (string, error) {
	pubECPoint := crypto.ScalarBaseMult(btcec.S256(), privKey)
	publicKey := &ecdsa.PublicKey{
		X:     big.NewInt(0).SetBytes(pubECPoint.X().Bytes()),
		Y:     big.NewInt(0).SetBytes(pubECPoint.Y().Bytes()),
		Curve: btcec.S256(),
	}

	hdPathSlices := strings.Split(hdPath, "/")
	chainIntStr := hdPathSlices[3]
	chainInt, err := strconv.Atoi(chainIntStr)
	if err != nil {
		return "", err
	}
	chainUint32 := uint32(chainInt)
	chain := SwitchChain(chainUint32)
	return switchChainAddress(publicKey, chain)
}

func SwitchChain(coinType uint32) string {
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
	case DASH:
		chain = "dash"
	case TRX:
		chain = "trx"
	case HECO:
		chain = "eth"
	case BSC:
		chain = "eth"
	case POLYGON:
		chain = "eth_arbitrum"
	case ARBITRUM:
		chain = "matic_polygon"
	default:
		panic("invalid chain type")
	}

	return chain
}

func deriveChildPrivKeys(key *RootKey, hdPath string) (*big.Int, error) {
	var buf [32]byte
	privKeyBytes := key.PrivKey.FillBytes(buf[:])

	childPrivateKeySlice, _, err := deriveChildPrivKey(privKeyBytes, key.ChainCode, key.PubKey, hdPath)
	if err != nil {
		return nil, err
	}
	privateKey := new(big.Int).SetBytes(childPrivateKeySlice[:])
	return privateKey, nil
}

func deriveChildPrivKey(prByte, chainCodeByte []byte, pubKey *crypto.ECPoint, path string) (childPrivKey [32]byte, childPubKey []byte, err error) {
	extendedKey := ckd.NewExtendKey(prByte, pubKey, pubKey, 0, 0, chainCodeByte)

	childPrivKey, childPubKey, err = ckd.DerivePrivateKeyForPath(extendedKey, path)
	if err != nil {
		return childPrivKey, nil, fmt.Errorf("derive child private err: %s", err.Error())
	}
	return childPrivKey, childPubKey, nil
}
