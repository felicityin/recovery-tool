package cmd

import (
	"archive/zip"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/btcsuite/btcd/btcec"
	ecies "github.com/ecies/go/v2"

	"recovery-tool/common"
	"recovery-tool/crypto"
)

const (
	// 81/WalletType/VaultIndex/CoinType/AddressIndex
	AssetWalletPath = "81/0/%d/%d/0"
	ApiWalletPath   = "81/1/0/%d/%d"
)

type RecoveryInput struct {
	ZipPath      string `yaml:"zip_path"`
	UserMnemonic string `yaml:"user_mnemonic"`
	EciesPrivKey string `yaml:"ecies_private_key"`
	RsaPrivKey   string `yaml:"coincover_private_key"`
	VaultCount   int    `yaml:"valut_count"`
	CoinType     []int  `yaml:"coin_type"`
}

type DeriveResult struct {
	VaultIndex int    `yaml:"vault_index"`
	CoinType   string `yaml:"coin_type"`
	Address    string `yaml:"address"`
	PrivKey    string `yaml:"private_key"`
}

type parsedParams struct {
	UserPrivKeyScalar *big.Int
	UserChainCode     []byte
	UserPubKey        string
	EciesPrivKey      *ecies.PrivateKey
	RsaPrivKey        *rsa.PrivateKey
}

func RecoverKeysCmd(paramsPath string, outputPath string) error {
	params := loadRecoveryParams(paramsPath)

	result, err := RecoverKeys(params)
	if err != nil {
		common.Logger.Errorf("derive keys failed")
		return err
	}

	if err := SaveResult(&result, outputPath); err != nil {
		common.Logger.Errorf("save result failed")
		return err
	}
	return nil
}

func RecoverKeys(params RecoveryInput) ([]*DeriveResult, error) {
	parsed, err := parseParams(params)
	if err != nil {
		return nil, err
	}

	hbcPrivs, err := findHbcPrivs(params.ZipPath, parsed.UserPubKey, parsed.EciesPrivKey, parsed.RsaPrivKey)
	if err != nil {
		common.Logger.Errorf("find hbc private info failed: %s", err)
		return nil, err
	}

	userPubKey := crypto.ScalarBaseMult(btcec.S256(), parsed.UserPrivKeyScalar)

	pubKey, err := hbcPrivs[0].PubKey.Add(hbcPrivs[1].PubKey)
	if err != nil {
		return nil, err
	}
	pubKey, err = pubKey.Add(userPubKey)
	if err != nil {
		return nil, err
	}

	privs := &common.RootKeys{
		HbcShare0: hbcPrivs[0],
		HbcShare1: hbcPrivs[1],
		UsrShare: &common.RootKey{
			PrivKey:   parsed.UserPrivKeyScalar,
			PubKey:    userPubKey,
			ChainCode: parsed.UserChainCode[:],
		},
		PubKey: pubKey,
	}

	keys, err := deriveChilds(params.VaultCount, params.CoinType, privs)
	if err != nil {
		common.Logger.Errorf("derive childs failed: %s", err)
		return nil, err
	}
	return keys, nil
}

func loadRecoveryParams(path string) RecoveryInput {
	bytess, err := ioutil.ReadFile(path)
	if err != nil {
		common.Logger.Errorf("load params error: %s", err.Error())
		panic(err)
	}

	var params RecoveryInput
	if err := yaml.UnmarshalStrict(bytess, &params); err != nil {
		common.Logger.Errorf("unmarshal params error: %s", err.Error())
		panic(err)
	}
	return params
}

func parseParams(params RecoveryInput) (*parsedParams, error) {
	userPrivKey, userChainCode, err := common.CalcMasterPriv(params.UserMnemonic)
	if err != nil {
		common.Logger.Errorf("calc user priv infos failed: %s", err)
		return nil, err
	}
	usrPrivKeyScalar := new(big.Int).SetBytes(userPrivKey[:])
	userPubKey := calcUserPubKey(usrPrivKeyScalar)
	common.Logger.Debugf("user pubkey: %s", userPubKey)

	eciesPrivKey, err := ecies.NewPrivateKeyFromHex(params.EciesPrivKey)
	common.Logger.Debugf("ecies privkey: %d", eciesPrivKey.D)
	if err != nil {
		common.Logger.Errorf("load ecies privkey failed: %s", err)
		return nil, err
	}

	rsaPrivKey, err := crypto.ParseRsaPrivKey(params.RsaPrivKey)
	common.Logger.Debugf("rsa privkey: %d, %d", rsaPrivKey.Primes[0], rsaPrivKey.Primes[1])
	if err != nil {
		common.Logger.Errorf("parse rsa privkey failed: %s", err)
		return nil, err
	}

	return &parsedParams{
		UserPrivKeyScalar: usrPrivKeyScalar,
		UserChainCode:     userChainCode[:],
		UserPubKey:        userPubKey,
		EciesPrivKey:      eciesPrivKey,
		RsaPrivKey:        rsaPrivKey,
	}, nil
}

func SaveResult(childs *[]*DeriveResult, outputPath string) error {
	yamlData, err := yaml.Marshal(childs)
	if err != nil {
		common.Logger.Errorf("yaml marshal result failed: %s", err)
		return err
	}
	err = ioutil.WriteFile(outputPath, yamlData, 0644)
	if err != nil {
		common.Logger.Errorf("unable to write data into the file")
		return err
	}
	return nil
}

func calcUserPubKey(privKey *big.Int) string {
	pubKey := crypto.ScalarBaseMult(crypto.S256(), privKey)
	pubKeyPoint := btcec.PublicKey{Curve: btcec.S256(), X: pubKey.X(), Y: pubKey.Y()}
	pubKeyBytes := pubKeyPoint.SerializeCompressed()
	return hex.EncodeToString(pubKeyBytes)
}

type encryptedTeam struct {
	HbcPrivKeys   []string `json:"hbc_private_keys"`
	HbcChainCodes []string `json:"hbc_chain_codes"`
	UserPubKey    string   `json:"user_pub_key"`
}

func findHbcPrivs(
	zipPath string,
	userPubKey string,
	eciesPrivKey *ecies.PrivateKey,
	rsaPrivKey *rsa.PrivateKey,
) ([]*common.RootKey, error) {
	zf, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	var result []*common.RootKey

	for _, file := range zf.File {
		fileBytes, err := common.ReadAll(file)
		if err != nil {
			return nil, err
		}

		encrypted := encryptedTeam{}
		err = json.Unmarshal(fileBytes, &encrypted)
		if err != nil {
			return nil, fmt.Errorf("unmarshal team failed: %s", err.Error())
		}

		decryptedUsrPubKey, err := decryptUserPubKey(encrypted.UserPubKey, eciesPrivKey, rsaPrivKey)
		common.Logger.Debugf("decrypted user pubkey: %s", decryptedUsrPubKey)
		if err != nil {
			return nil, err
		}

		if decryptedUsrPubKey == userPubKey {
			priv0, err := decryptHbcPriv(encrypted.HbcPrivKeys[0], encrypted.HbcChainCodes[0], eciesPrivKey, rsaPrivKey)
			if err != nil {
				return nil, err
			}
			priv1, err := decryptHbcPriv(encrypted.HbcPrivKeys[1], encrypted.HbcChainCodes[1], eciesPrivKey, rsaPrivKey)
			if err != nil {
				return nil, err
			}

			result = append(result, priv0)
			result = append(result, priv1)
			return result, nil
		}
	}
	return nil, fmt.Errorf("mnemonic and zip do not match")
}

func decryptUserPubKey(userPubKey string, eciesPrivKey *ecies.PrivateKey, rsaPrivKey *rsa.PrivateKey) (string, error) {
	userPubKeyBytes, err := hex.DecodeString(userPubKey)
	if err != nil {
		return "", fmt.Errorf("hex decode user pubkey error: %s", err.Error())
	}

	decryptedUsrPubKey, err := crypto.RsaDecryptOAEP(rsaPrivKey, userPubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("rsa decrypt user pubkey error: %s", err.Error())
	}

	decryptedUsrPubKey, err = ecies.Decrypt(eciesPrivKey, decryptedUsrPubKey)
	if err != nil {
		return "", fmt.Errorf("ecies decrypt user pubkey error: %s", err.Error())
	}

	return hex.EncodeToString(decryptedUsrPubKey), nil
}

func decryptHbcPriv(
	privKey, chainCode string,
	eciesPrivKey *ecies.PrivateKey,
	rsaPrivKey *rsa.PrivateKey,
) (*common.RootKey, error) {
	privKeyBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, fmt.Errorf("hex decode privkey failed: %s", err.Error())
	}
	chainCodeBytes, err := hex.DecodeString(chainCode)
	if err != nil {
		return nil, fmt.Errorf("hex decode chaincode failed: %s", err.Error())
	}

	decryptedPrivKey, err := crypto.RsaDecryptOAEP(rsaPrivKey, privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("rsa decode privkey failed: %s", err.Error())
	}
	decryptedChainCode, err := crypto.RsaDecryptOAEP(rsaPrivKey, chainCodeBytes)
	if err != nil {
		return nil, fmt.Errorf("rsa decode chaincode failed: %s", err.Error())
	}

	decryptedPrivKey, err = ecies.Decrypt(eciesPrivKey, decryptedPrivKey)
	if err != nil {
		return nil, fmt.Errorf("ecies decode privkey failed: %s", err.Error())
	}
	decryptedChainCode, err = ecies.Decrypt(eciesPrivKey, decryptedChainCode)
	if err != nil {
		return nil, fmt.Errorf("ecies decode chaincode failed: %s", err.Error())
	}

	privateKey := new(big.Int).SetBytes(decryptedPrivKey)

	return &common.RootKey{
		PrivKey:   privateKey,
		PubKey:    crypto.ScalarBaseMult(btcec.S256(), privateKey),
		ChainCode: decryptedChainCode,
	}, nil
}

func deriveChilds(vaultCount int, coinType []int, rootKeys *common.RootKeys) ([]*DeriveResult, error) {
	var deriveResult []*DeriveResult

	for vaultIndex := 0; vaultIndex < vaultCount; vaultIndex++ {
		for _, coin := range coinType {
			hdPath := fmt.Sprintf(AssetWalletPath, vaultIndex, coin) // Only support asset wallet for now
			privKey, address, err := common.DeriveChild(rootKeys, hdPath, coin)
			if err != nil {
				return nil, fmt.Errorf("derive child failed, err: %s", err.Error())
			}

			var buf [32]byte
			privKeyBytes := privKey.FillBytes(buf[:])

			deriveResult = append(deriveResult, &DeriveResult{
				VaultIndex: vaultIndex,
				CoinType:   common.SwitchChain(uint32(coin)),
				Address:    address,
				PrivKey:    hex.EncodeToString(privKeyBytes),
			})
		}
	}

	return deriveResult, nil
}
