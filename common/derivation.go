package common

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/decred/dcrd/dcrec/edwards/v2"

	"recovery-tool/crypto"
	"recovery-tool/crypto/ckd"
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
	PubKey    *crypto.ECPoint
}

type RootKey struct {
	PrivKey   *big.Int
	PubKey    *crypto.ECPoint
	ChainCode []byte
}

func DeriveChild(params *RootKeys, hdPath string, coin int) (*big.Int, string, error) {
	privKey, err := DerivePrivKey(params, hdPath, coin)
	if err != nil {
		return nil, "", err
	}

	address, err := DeriveAddress(privKey, hdPath, coin)
	if err != nil {
		return nil, "", err
	}

	return privKey, address, nil
}

func DerivePrivKey(params *RootKeys, hdPath string, coin int) (*big.Int, error) {
	hbcPrivKey0, n, err := deriveChildPrivKey(params.HbcShare0, hdPath, params.PubKey, coin)
	if err != nil {
		return nil, err
	}

	hbcPrivKey1, _, err := deriveChildPrivKey(params.HbcShare1, hdPath, params.PubKey, coin)
	if err != nil {
		return nil, err
	}

	userPrivKey, _, err := deriveChildPrivKey(params.UsrShare, hdPath, params.PubKey, coin)
	if err != nil {
		return nil, err
	}

	privateKey := hbcPrivKey0

	privateKey.Add(privateKey, hbcPrivKey1)
	privateKey.Mod(privateKey, n)

	privateKey.Add(privateKey, userPrivKey)
	privateKey.Mod(privateKey, n)

	return privateKey, nil
}

func DeriveAddress(privKey *big.Int, hdPath string, coin int) (string, error) {
	chain := GetChainName(uint32(coin))

	if isEddsaCoin(coin) {
		pubECPoint := crypto.ScalarBaseMult(edwards.Edwards(), privKey)
		publicKey := edwards.NewPublicKey(pubECPoint.X(), pubECPoint.Y())
		return SwitchEddsaChainAddress(publicKey, chain)
	} else {
		pubECPoint := crypto.ScalarBaseMult(btcec.S256(), privKey)
		publicKey := &ecdsa.PublicKey{
			X:     pubECPoint.X(),
			Y:     pubECPoint.Y(),
			Curve: btcec.S256(),
		}
		return SwitchEcdsaChainAddress(publicKey, chain)
	}
}

func GetChainName(coinType uint32) string {
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
		chain = "ht_heco"
	case BSC:
		chain = "bnb_bsc"
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
	default:
		panic("invalid chain type")
	}

	return chain
}

func deriveChildPrivKey(key *RootKey, hdPath string, deducePubKey *crypto.ECPoint, coin int) (*big.Int, *big.Int, error) {
	var buf [32]byte
	privKeyBytes := key.PrivKey.FillBytes(buf[:])

	var childPrivateKeySlice [32]byte
	var err error
	var n *big.Int

	if isEddsaCoin(coin) {
		n = edwards.Edwards().Params().N
		childPrivateKeySlice, _, err = DeriveEddsaChildPrivKey(privKeyBytes, key.ChainCode, key.PubKey, deducePubKey, hdPath)
	} else {
		n = btcec.S256().Params().N
		childPrivateKeySlice, _, err = DeriveEcdsaChildPrivKey(privKeyBytes, key.ChainCode, key.PubKey, hdPath)
	}
	if err != nil {
		return nil, nil, err
	}

	privateKey := new(big.Int).SetBytes(childPrivateKeySlice[:])
	return privateKey, n, nil
}

func DeriveEcdsaChildPrivKey(
	prByte, chainCodeByte []byte,
	pubKey *crypto.ECPoint,
	path string,
) (childPrivKey [32]byte, childPubKey []byte, err error) {
	extendedKey := ckd.NewExtendKey(prByte, pubKey, pubKey, 0, 0, chainCodeByte)

	childPrivKey, childPubKey, err = ckd.DerivePrivateKeyForPath(extendedKey, path)
	if err != nil {
		return childPrivKey, nil, fmt.Errorf("derive child private err: %s", err.Error())
	}
	return childPrivKey, childPubKey, nil
}

func DeriveEddsaChildPrivKey(
	prByte, chainCodeByte []byte,
	pubKey *crypto.ECPoint,
	deducePubKey *crypto.ECPoint,
	path string,
) (childPrivKey [32]byte, childPubKey []byte, err error) {
	extendedKey := ckd.NewExtendKeyD(prByte, pubKey, deducePubKey, 0, 0, chainCodeByte)

	childPrivKey, childPubKey, err = ckd.DerivePrivateKeyForPathD(extendedKey, path, pubKey.Curve())
	if err != nil {
		return childPrivKey, nil, fmt.Errorf("derive child private err: %s", err.Error())
	}
	return childPrivKey, childPubKey, nil
}

func isEddsaCoin(coin int) bool {
	return coin == 354 || coin == 501 || coin == 637
}
